package bits

import (
	"fmt"
	"reflect"
)

type ItuDecoder struct {
	reader *ItuReader
	err    error
}

func NewItuDecoder(r *ItuReader) *ItuDecoder {
	return &ItuDecoder{reader: r}
}

func (d *ItuDecoder) Decode(v any, fieldName string) error {
	return d.DecodeRange(v, fieldName, fieldName)
}

func (d *ItuDecoder) DecodeRange(v any, start string, end string) error {
	if d.err != nil {
		return d.err
	}
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		d.err = fmt.Errorf("expected a pointer to a struct, got %T", v)
		return d.err
	}
	typ := val.Elem().Type()

	startField, ok := typ.FieldByName(start)
	if !ok {
		d.err = fmt.Errorf("start field %s not found in struct", start)
		return d.err
	}
	endField, ok := typ.FieldByName(end)
	if !ok {
		d.err = fmt.Errorf("end field %s not found in struct", start)
		return d.err
	}

	for fieldIndex := startField.Index[0]; fieldIndex <= endField.Index[0]; fieldIndex++ {
		structField := typ.FieldByIndex([]int{fieldIndex})
		structVal := val.FieldByIndex([]int{fieldIndex})
		descriptor := structField.Tag.Get("descriptor")
		// We just skip fields without descriptors
		if descriptor == "" {
			continue
		}
		err := d.load(structField.Name, structVal, descriptor)
		if err != nil {
			d.err = err
			return d.err
		}

	}
	return nil
}

func (d *ItuDecoder) DecodeIndex(v any, fieldName string, index int) error {
	if d.err != nil {
		return d.err
	}
	if d.err != nil {
		return d.err
	}
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		d.err = fmt.Errorf("expected a pointer to a struct, got %T", v)
		return d.err
	}
	typ := val.Elem().Type()

	structField, ok := typ.FieldByName(fieldName)
	if !ok {
		d.err = fmt.Errorf("field %s not found in struct", fieldName)
		return d.err
	}
	descriptor := structField.Tag.Get("descriptor")
	structVal := val.FieldByName(fieldName)
	if structVal.Kind() != reflect.Slice {
		d.err = fmt.Errorf("field %s is not a slice", fieldName)
		return d.err
	}
	if structVal.IsZero() {
		structVal.Set(reflect.MakeSlice(structVal.Type().Elem(), index, index+16))
	}
	if index+1 > structVal.Cap() {
		// We don't want to do this too often, so we're going to increase by
		// the sum of the index and the length -- meaning we'll always at least triple
		// the capacity.
		structVal.Grow(index + structVal.Len())
	}
	if index+1 > structVal.Len() {
		structVal.SetLen(index)
	}
	err := d.load(fmt.Sprintf("%v[%v]", structField.Name, index), structVal.Index(index), descriptor)
	if err != nil {
		d.err = err
		return d.err
	}
	return nil
}

func (d *ItuDecoder) Error() error {
	return d.err
}

func (d *ItuDecoder) load(name string, val reflect.Value, descriptor string) error {
	if !val.CanSet() {
		return fmt.Errorf("field %v cannot be set", name)
	}
	descriptorType, descriptorLength, fixedValue, err := parseDescriptor(descriptor)
	if err != nil {
		return fmt.Errorf("could not decode field %v: %w", name, err)
	}
	switch descriptorType {
	case "ae", "ce", "i", "me", "st", "te":
		return fmt.Errorf("field %v: descriptor type %v is not yet supported", name, descriptorType)
	case "u", "b", "f":
		v, err := d.reader.U(descriptorLength)
		if err != nil {
			return fmt.Errorf("field %v: %w", name, err)
		}
		if val.CanUint() {
			val.SetUint(v)
		} else if val.Kind() == reflect.Bool {
			val.SetBool(v != 0)
		} else {
			return fmt.Errorf("field %v: cannot store value of descriptor type %v in %v", name, descriptor, val.Kind())
		}
		if descriptorType == "f" {
			if v != fixedValue {
				return fmt.Errorf("field %v: value %v does not match expected value %v", name, v, fixedValue)
			}
		}
	case "ue":
		v, err := d.reader.UE()
		if err != nil {
			return fmt.Errorf("field %v: %w", name, err)
		}
		if val.CanUint() {
			val.SetUint(v)
		} else {
			return fmt.Errorf("field %v: cannot store value of descriptor type %v in %v", name, descriptor, val.Kind())
		}

	case "se":
		v, err := d.reader.SE()
		if err != nil {
			return fmt.Errorf("field %v: %w", name, err)
		}
		if val.CanInt() {
			val.SetInt(v)
		} else {
			return fmt.Errorf("field %v: cannot store value of descriptor type %v in %v", name, descriptor, val.Kind())
		}
	default:
		return fmt.Errorf("field %v: descriptor type %v is invalid", name, descriptorType)
	}
	return nil
}

func parseDescriptor(descriptor string) (typ string, length int, fixedValue uint64, err error) {
	_, err = fmt.Sscanf(descriptor, "%[a-z](%d)=%d", &typ, &length, &fixedValue)
	if err == nil {
		return typ, length, fixedValue, nil
	}

	_, err = fmt.Sscanf(descriptor, "%[a-z](%d)", &typ, &length)
	if err == nil {
		return typ, length, 0, nil
	}

	_, err = fmt.Sscanf(descriptor, "%[a-z](v)", &typ)
	if err == nil {
		return typ, 0, 0, nil
	}

	return "", 0, 0, fmt.Errorf("invalid descriptor format: %s", descriptor)
}
