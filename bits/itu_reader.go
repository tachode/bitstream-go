package bits

type ItuReader struct {
	Reader
}

// This is simply an alias so that the methods on ItuReader match the syntax in the
// H.264 spec
func (r *ItuReader) U(bits int) (uint64, error) {
	return r.ReadBits(bits)
}

func (r *ItuReader) UE() (uint64, error) {
	leadingZeroBits := 0

	var val uint64
	var err error

	for {
		bit, err := r.ReadBits(1)
		if err != nil {
			return 0, err
		}
		if bit == 1 {
			break
		}
		leadingZeroBits++
	}

	if leadingZeroBits != 0 {
		val, err = r.ReadBits(leadingZeroBits)
		if err != nil {
			return 0, err
		}
	}
	// Need to put the leading bit back onto the value
	// (since we read it while we were counting zero bits)
	val |= 1 << (leadingZeroBits)

	return val - 1, nil
}

func (r *ItuReader) SE() (int64, error) {
	ueVal, err := r.UE()
	if err != nil {
		return 0, err
	}

	if ueVal == 0 {
		return 0, nil
	}

	if ueVal%2 == 0 {
		return -int64(ueVal / 2), nil
	}
	return int64((ueVal + 1) / 2), nil
}
