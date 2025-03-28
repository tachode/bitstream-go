package bits

import (
	"fmt"
	"math"
	"reflect"
)

type ItuWriter struct {
	Writer
}

// This is simply an alias so that the methods on ItuWriter match the syntax in the
// H.264 spec
func (r *ItuWriter) U(v any, bits int) error {
	return r.WriteBits(v, bits)
}

// Unsigned exponential Golomb encoding
func (w *ItuWriter) UE(v any) error {

	val := reflect.ValueOf(v)
	var value uint64
	switch {
	case val.CanUint():
		value = uint64(val.Uint())
	default:
		return fmt.Errorf("UE(%v): cannot convert %T into a uint64", v, v)
	}

	// Calculate the number of bits required
	bits := int(math.Log2(float64(value+1))) + 1

	// Write leading zeros
	w.WriteBits(uint64(0), bits-1)

	// Write the value
	w.WriteBits(uint64(value+1), bits)

	return nil
}

// Signed exponential Golomb encoding
func (w *ItuWriter) SE(v any) error {

	val := reflect.ValueOf(v)
	var ival int64
	switch {
	case val.CanInt():
		ival = val.Int()
	default:
		return fmt.Errorf("SE(%v): cannot convert %T into an int64", v, v)
	}

	// Map signed value to unsigned value
	var uval uint64
	if ival > 0 {
		uval = uint64(2*ival - 1)
	} else {
		uval = uint64(-2 * ival)
	}

	// Use UE to encode the mapped value
	return w.UE(uval)
}
