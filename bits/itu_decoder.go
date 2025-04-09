package bits

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

type ItuDecoder struct {
	reader      *ItuReader
	value       map[string]any
	valueSource map[string]string
	valueLength map[string]int
	err         error
	decodingLog []string
}

// NewItuDecoder creates a new ItuDecoder instance.
// It takes an io.Reader, an *ItuReader, or a byte
// slice as input.
func NewItuDecoder(in any) *ItuDecoder {
	d := &ItuDecoder{
		value:       make(map[string]any),
		valueSource: make(map[string]string),
		valueLength: make(map[string]int),
		decodingLog: make([]string, 0),
	}
	d.log("Creating new ItuDecoder")
	d.Reset(in)
	return d
}

func (d *ItuDecoder) Log() []string {
	return d.decodingLog
}

func (d *ItuDecoder) SetValueLength(name string, length int) {
	d.valueLength[name] = length
	d.log(fmt.Sprintf("Setting length of field '%s' to %d", name, length))
}

func (d *ItuDecoder) SetValue(name string, value any) {
	d.value[name] = value
	pc, _, _, ok := runtime.Caller(1)
	funcName := "SetValue()"
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	d.valueSource[name] = funcName + "()"
	d.log(fmt.Sprintf("Setting field '%s' to %v in %v()", name, value, funcName))
}

// Reset resets the ItuDecoder instance with a new input,
// but retains all previously decoded variable values.
// It takes an io.Reader, an *ItuReader, or a byte
// slice as input.
func (d *ItuDecoder) Reset(in any) error {
	// In order to implement several H.264 functions, parsing of RBSP payloads
	// needs to be buffered. The functions that require this are more_rbsp_data()
	// more_rbsp_trailing_data(), and next_bits(). To support, this, any Reader
	// passed to the ItuDecoder is converted to a bufio.Reader if it isn't already.

	if reader, ok := in.(*ItuReader); ok {
		d.reader = reader
	} else if buf, ok := in.([]byte); ok {
		reader := &ReadBuffer{Reader: bufio.NewReader(bytes.NewReader(buf))}
		d.reader = &ItuReader{Reader: reader}
	} else if reader, ok := in.(io.Reader); ok {
		if _, ok := reader.(*bufio.Reader); !ok {
			reader = bufio.NewReader(reader)
		}
		d.reader = &ItuReader{Reader: &ReadBuffer{Reader: reader}}
	} else {
		return d.setError(fmt.Errorf("unsupported type %T", in))
	}
	d.err = nil
	d.decodingLog = make([]string, 0)
	d.log(fmt.Sprintf("Resetting decoder with %T", in))
	return nil
}

func (d *ItuDecoder) Value(name string) any {
	val, ok := d.value[name]
	if !ok {
		d.log(fmt.Sprintf("Attempt to read unset variable: '%s'", name))
		return nil
	}
	from := d.valueSource[name]

	d.log(fmt.Sprintf("Reading variable '%s' = %v (from %v)", name, val, from))
	return val
}

func (d *ItuDecoder) Read(b []byte) (n int, err error) {
	n, err = d.reader.Read(b)
	if err != nil {
		d.setError(err)
	}
	return n, err
}

func (d *ItuDecoder) Decode(v any, fieldName string) error {
	return d.DecodeRange(v, fieldName, fieldName)
}

func (d *ItuDecoder) DecodeRange(v any, start string, end string) error {
	if d.err != nil {
		return d.err
	}
	structName := fmt.Sprintf("%T", v)
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return d.setError(fmt.Errorf("expected a pointer to a struct, got %T", v))
	}
	val = val.Elem()
	typ := val.Type()

	startField, ok := typ.FieldByName(start)
	if !ok {
		return d.setError(fmt.Errorf("start field %s not found in struct", start))
	}
	endField, ok := typ.FieldByName(end)
	if !ok {
		return d.setError(fmt.Errorf("end field %s not found in struct", start))
	}

	for fieldIndex := startField.Index[0]; fieldIndex <= endField.Index[0]; fieldIndex++ {
		structField := typ.FieldByIndex([]int{fieldIndex})
		structVal := val.FieldByIndex([]int{fieldIndex})
		descriptor := structField.Tag.Get("descriptor")
		// We just skip fields without descriptors
		if descriptor == "" {
			continue
		}
		err := d.load(structName, structField.Name, structVal, descriptor)
		if err != nil {
			return d.setError(err)
		}

	}
	return nil
}

// TODO -- implement multidimensional arrays
func (d *ItuDecoder) DecodeIndex(v any, fieldName string, index int, subindex ...int) error {
	if d.err != nil {
		return d.err
	}
	val := reflect.ValueOf(v)
	structName := fmt.Sprintf("%T", v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return d.setError(fmt.Errorf("expected a pointer to a struct, got %T", v))
	}
	val = val.Elem()
	typ := val.Type()

	structField, ok := typ.FieldByName(fieldName)
	if !ok {
		return d.setError(fmt.Errorf("field %s not found in struct", fieldName))
	}
	descriptor := structField.Tag.Get("descriptor")
	structVal := val.FieldByName(fieldName)
	return d.decodeIndex(structName, fieldName, structVal, descriptor, index, subindex...)
}

func (d *ItuDecoder) decodeIndex(structName string, fieldName string, structVal reflect.Value, descriptor string, index int, subindex ...int) error {
	if structVal.Kind() != reflect.Slice {
		return d.setError(fmt.Errorf("field %s is not a slice", fieldName))
	}
	if structVal.IsZero() {
		structVal.Set(reflect.MakeSlice(structVal.Type(), index, index+16))
	}
	if index+1 > structVal.Cap() {
		// We don't want to do this too often, so we're going to increase by
		// the sum of the index and the length -- meaning we'll always at least triple
		// the capacity.
		structVal.Grow(index + structVal.Len())
	}
	if index+1 > structVal.Len() {
		structVal.SetLen(index + 1)
	}

	if len(subindex) > 0 {
		if structVal.Index(index).Kind() != reflect.Slice {
			return d.setError(fmt.Errorf("field %s[%d] is not a slice", fieldName, index))
		}
		structVal = structVal.Index(index)
		fieldName = fmt.Sprintf("%s[%d]", fieldName, index)
		index = subindex[0]
		subindex = subindex[1:]
		return d.decodeIndex(structName, fieldName, structVal, descriptor, index, subindex...)
	} else {
		err := d.load(structName, fmt.Sprintf("%v[%v]", fieldName, index), structVal.Index(index), descriptor)
		if err != nil {
			return d.setError(fmt.Errorf("field %s[%d]: %w", fieldName, index, err))
		}
	}

	return nil
}

func (d *ItuDecoder) Error() error {
	return d.err
}

func (d *ItuDecoder) load(structName string, name string, val reflect.Value, descriptor string) error {
	var bitsRead int
	if strings.Contains(structName, ".") {
		structName = structName[strings.LastIndex(structName, ".")+1:]
	}
	if !val.CanSet() {
		return fmt.Errorf("field %v cannot be set", name)
	}
	descriptorType, descriptorLength, fixedValue, err := parseDescriptor(descriptor)
	if err != nil {
		return fmt.Errorf("could not decode field %v: %w", name, err)
	}

	// For fixed-length descriptors, "variable-length" descriptors
	// need to have their actual length set by the caller
	switch descriptorType {
	case "u", "b", "f", "i":
		if descriptorLength == 0 {
			var ok bool
			baseName := name
			if strings.Contains(name, "[") {
				baseName = name[:strings.Index(name, "[")]
			}
			descriptorLength, ok = d.valueLength[baseName]
			if !ok {
				return d.setError(fmt.Errorf("field %v: descriptor length is not set for %s", name, baseName))
			}
		}
	}

	d.valueSource[name] = structName
	switch descriptorType {
	case "ae", "ce", "me", "st", "te":
		return fmt.Errorf("field %v: descriptor type %v is not yet supported", name, descriptorType)
	case "i":
		v, n, err := d.reader.I(descriptorLength)
		bitsRead = n
		if err != nil {
			return fmt.Errorf("field %v: %w", name, err)
		}
		if val.CanInt() {
			val.SetInt(v)
			d.value[name] = v
		} else {
			return fmt.Errorf("field %v: cannot store value of descriptor type %v in %v", name, descriptor, val.Kind())
		}
	case "u", "b", "f":
		v, n, err := d.reader.U(descriptorLength)
		bitsRead = n
		if err != nil {
			return fmt.Errorf("field %v: %w", name, err)
		}
		if val.CanUint() {
			val.SetUint(v)
			d.value[name] = v
		} else if val.Kind() == reflect.Bool {
			val.SetBool(v != 0)
			d.value[name] = (v != 0)
		} else {
			return fmt.Errorf("field %v: cannot store value of descriptor type %v in %v", name, descriptor, val.Kind())
		}
		if descriptorType == "f" {
			if v != fixedValue {
				return fmt.Errorf("field %v: value %v does not match expected value %v with descriptor %v", name, v, fixedValue, descriptor)
			}
		}
	case "ue":
		v, n, err := d.reader.UE()
		bitsRead = n
		if err != nil {
			return fmt.Errorf("field %v: %w", name, err)
		}
		if val.CanUint() {
			val.SetUint(v)
			d.value[name] = v
		} else {
			return fmt.Errorf("field %v: cannot store value of descriptor type %v in %v", name, descriptor, val.Kind())
		}

	case "se":
		v, n, err := d.reader.SE()
		bitsRead = n
		if err != nil {
			return fmt.Errorf("field %v: %w", name, err)
		}
		if val.CanInt() {
			val.SetInt(v)
			d.value[name] = v
		} else {
			return fmt.Errorf("field %v: cannot store value of descriptor type %v in %v", name, descriptor, val.Kind())
		}
	default:
		return d.setError(fmt.Errorf("field %v: descriptor type %v is invalid", name, descriptorType))
	}
	bitsWord := "bits"
	if bitsRead == 1 {
		bitsWord = "bit"
	}
	d.log(fmt.Sprintf("Setting %s.%s = %v (%s ⇒ %v %s)", structName, name, val, descriptor, bitsRead, bitsWord))
	return nil
}

func (d *ItuDecoder) ByteAligned() bool {
	if d.err != nil {
		return true
	}
	return d.reader.ByteAligned()
}

func (d *ItuDecoder) MoreRbspData() bool {
	if d.err != nil {
		return false
	}
	return d.reader.MoreRbspData()
}

func (d *ItuDecoder) NextBits(bits int) uint64 {
	if d.err != nil {
		return 0
	}
	val, _ := d.reader.NextBits(bits)
	return val
}

func parseDescriptor(descriptor string) (typ string, length int, fixedValue uint64, err error) {
	f := func(c rune) bool { return !unicode.IsLetter(c) && !unicode.IsNumber(c) }
	fields := strings.FieldsFunc(descriptor, f)

	if len(fields) < 2 || len(fields) > 3 {
		err = fmt.Errorf("invalid descriptor format: %s", descriptor)
		return
	}
	typ = fields[0]
	length, _ = strconv.Atoi(fields[1])
	if len(fields) == 3 {
		fixedValue, _ = strconv.ParseUint(fields[2], 10, 0)
	}
	return
}

func (d *ItuDecoder) setError(err error) error {
	d.err = err
	pc, file, line, ok := runtime.Caller(1)
	funcName := ""
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	if strings.Contains(funcName, ".") {
		funcName = funcName[strings.LastIndex(funcName, ".")+1:]
	}
	if strings.Contains(file, "/") {
		file = file[strings.LastIndex(file, "/")+1:]
	}

	if err != nil {
		d.log(fmt.Sprintf("Error in %v (%v:%v): %v", funcName, file, line, err))
	}
	return err
}

func (d *ItuDecoder) log(msg string) {
	// fmt.Println(msg)
	d.decodingLog = append(d.decodingLog, msg)
}
