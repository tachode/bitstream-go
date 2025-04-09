package bits_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tachode/bitstream-go/bits"
)

func TestNewItuDecoder(t *testing.T) {
	data := []byte{0xFF, 0x00, 0xAA}
	decoder := bits.NewItuDecoder(data)

	assert.NotNil(t, decoder)
	assert.Nil(t, decoder.Error())
	assert.NotNil(t, decoder.Log())
}

func TestItuDecoder_Reset(t *testing.T) {
	data := []byte{0xFF, 0x00, 0xAA}
	decoder := bits.NewItuDecoder(data)

	err := decoder.Reset(data)
	assert.Nil(t, err)
	assert.Nil(t, decoder.Error())
}

func TestItuDecoder_SetValueAndValue(t *testing.T) {
	decoder := bits.NewItuDecoder(nil)

	decoder.SetValue("testField", 42)
	value := decoder.Value("testField")

	assert.Equal(t, 42, value)
}

type MockReader struct{ BitsRead int }

func (m *MockReader) ReadBits(bits int) (uint64, error) { m.BitsRead += bits; return 12, nil }
func (m *MockReader) NextBits(bits int) (uint64, error) { return 0, nil }
func (m *MockReader) ByteAligned() bool                 { return m.BitsRead%8 == 0 }
func (m *MockReader) MoreRbspData() bool                { return false }
func (m *MockReader) Align() error                      { return nil }
func (m *MockReader) Read([]byte) (int, error)          { return 0, nil }

func TestItuDecoder_SetValueLength(t *testing.T) {
	st := struct {
		A uint `descriptor:"u(v)"`
	}{}
	mockReader := &MockReader{}
	ituReader := &bits.ItuReader{Reader: mockReader}
	decoder := bits.NewItuDecoder(ituReader)

	decoder.SetValueLength("A", 6)
	err := decoder.Decode(&st, "A")
	require.NoError(t, err)

	assert.EqualValues(t, 12, st.A)
	assert.Equal(t, 6, mockReader.BitsRead)
}

func TestItuDecoder_ErrorHandling(t *testing.T) {
	decoder := bits.NewItuDecoder(nil)

	err := decoder.Reset(123) // Invalid input type
	assert.NotNil(t, err)
	assert.NotNil(t, decoder.Error())
}

func TestItuDecoder_ByteAligned(t *testing.T) {
	data := []byte{0xFF, 0x00, 0xAA}
	decoder := bits.NewItuDecoder(data)

	assert.True(t, decoder.ByteAligned())
	data = []byte{0xFF, 0x00, 0xAA, 0xF0}
	decoder = bits.NewItuDecoder(data)

	st := struct {
		A uint `descriptor:"u(6)"`
	}{}
	decoder.Decode(&st, "A")

	assert.False(t, decoder.ByteAligned())
}

func TestItuDecoder_MoreRbspData(t *testing.T) {
	data := []byte{0xFF, 0x00, 0xAA}
	decoder := bits.NewItuDecoder(data)

	assert.True(t, decoder.MoreRbspData())
}

func TestItuDecoder_NextBits(t *testing.T) {
	data := []byte{0xFF, 0x00, 0xAA}
	decoder := bits.NewItuDecoder(data)

	bits := decoder.NextBits(8)
	assert.Equal(t, uint64(0xFF), bits)
}

func TestItuDecoder_Decode(t *testing.T) {
	type TestStruct struct {
		Field1 uint `descriptor:"u(8)"`
		Field2 uint `descriptor:"u(8)"`
	}

	data := []byte{0x12, 0x34}
	decoder := bits.NewItuDecoder(data)

	var testStruct TestStruct
	err := decoder.Decode(&testStruct, "Field1")
	require.NoError(t, err)

	assert.Equal(t, uint(0x12), testStruct.Field1)
}

func TestItuDecoder_DecodeWithSliceFields(t *testing.T) {
	type TestStruct struct {
		Field1 []uint `descriptor:"u(8)"`
		Field2 []uint `descriptor:"u(4)"`
	}

	data := []byte{0x12, 0x34, 0x56, 0x78}
	decoder := bits.NewItuDecoder(data)

	var testStruct TestStruct

	err := decoder.DecodeIndex(&testStruct, "Field1", 0)
	require.NoError(t, err)
	err = decoder.DecodeIndex(&testStruct, "Field1", 1)
	require.NoError(t, err)
	assert.Equal(t, []uint{0x12, 0x34}, testStruct.Field1)

	err = decoder.DecodeIndex(&testStruct, "Field2", 0)
	require.NoError(t, err)
	err = decoder.DecodeIndex(&testStruct, "Field2", 1)
	require.NoError(t, err)
	assert.Equal(t, []uint{0x05, 0x06}, testStruct.Field2)
}

func TestItuDecoder_DecodeWithMultidimensionalSliceFields(t *testing.T) {
	type TestStruct struct {
		Field1 [][][]uint `descriptor:"u(8)"`
		Field2 [][][]uint `descriptor:"u(4)"`
	}

	data := []byte{0x12, 0x34, 0x56, 0x78}
	decoder := bits.NewItuDecoder(data)

	var testStruct TestStruct

	err := decoder.DecodeIndex(&testStruct, "Field1", 0, 1, 1)
	require.NoError(t, err)
	err = decoder.DecodeIndex(&testStruct, "Field1", 1, 2, 3)
	require.NoError(t, err)
	assert.EqualValues(t, 0x12, testStruct.Field1[0][1][1])
	assert.EqualValues(t, 0x34, testStruct.Field1[1][2][3])
	assert.EqualValues(t, 0x12, decoder.Value("Field1[0][1][1]"))
	assert.EqualValues(t, 0x34, decoder.Value("Field1[1][2][3]"))

	err = decoder.DecodeIndex(&testStruct, "Field2", 0, 2, 1)
	require.NoError(t, err)
	err = decoder.DecodeIndex(&testStruct, "Field2", 1, 1, 3)
	require.NoError(t, err)
	assert.EqualValues(t, 0x05, testStruct.Field2[0][2][1])
	assert.EqualValues(t, 0x06, testStruct.Field2[1][1][3])
	assert.EqualValues(t, 0x05, decoder.Value("Field2[0][2][1]"))
	assert.EqualValues(t, 0x06, decoder.Value("Field2[1][1][3]"))

	assert.Nil(t, decoder.Value("Field2[1][1]"))

	// Note: the semantics of the base name of a slice might change
	// in the future
	assert.Nil(t, decoder.Value("Field2"))

	t.Logf("Log:\n%s", strings.Join(decoder.Log(), "\n"))
}
