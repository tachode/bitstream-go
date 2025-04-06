package bits

import "io"

type ItuReader struct {
	Reader
}

// This is simply an alias so that the methods on ItuReader match the syntax in the
// H.264 spec
func (r *ItuReader) U(bits int) (val uint64, n int, err error) {
	val, err = r.ReadBits(bits)
	return val, bits, err
}

func (r *ItuReader) UE() (val uint64, n int, err error) {
	leadingZeroBits := 0

	for {
		bit, err := r.ReadBits(1)
		n += 1
		if err != nil {
			return 0, n, err
		}
		if bit == 1 {
			break
		}
		leadingZeroBits++
	}

	if leadingZeroBits != 0 {
		val, err = r.ReadBits(leadingZeroBits)
		n += leadingZeroBits
		if err != nil {
			return 0, n, err
		}
	}
	// Need to put the leading bit back onto the value
	// (since we read it while we were counting zero bits)
	val |= 1 << (leadingZeroBits)

	return val - 1, n, nil
}

func (r *ItuReader) SE() (val int64, n int, err error) {
	ueVal, n, err := r.UE()
	if err != nil {
		return 0, n, err
	}

	if ueVal == 0 {
		return 0, n, nil
	}

	if ueVal%2 == 0 {
		return -int64(ueVal / 2), n, nil
	}
	return int64((ueVal + 1) / 2), n, nil
}

func (r *ItuReader) MoreDataInByteStream() bool {
	_, err := r.NextBits(8)
	return err != io.EOF
}
