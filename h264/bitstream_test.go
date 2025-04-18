package h264_test

var _ = bitstreamBytes
var _ = fuStartBytes
var _ = accessUnitDelimiterBytes
var _ = spsBytes
var _ = ppsBytes
var _ = seiBytes

var bitstreamBytes = [...]byte{
	0x00, 0x00, 0x01, 0xE0, 0x00, 0x00, 0x84, 0x80, 0x05, 0x21, 0x00, 0x37, 0x77, 0x41, 0x00, 0x00,
	0x00, 0x01, 0x09, 0xF0, 0x00, 0x00, 0x00, 0x01, 0x67, 0x42, 0xC0, 0x1E, 0xB9, 0x08, 0x0D, 0x0F,
	0xFC, 0x98, 0x08, 0x80, 0x00, 0x01, 0xF4, 0x80, 0x00, 0x5D, 0xC0, 0x70, 0x30, 0x01, 0x48, 0x20,
	0x02, 0x90, 0x77, 0xBD, 0xC0, 0x7C, 0x22, 0x11, 0x92, 0x00, 0x00, 0x00, 0x01, 0x68, 0xDB, 0x8F,
	0x20, 0x00, 0x00, 0x01, 0x06, 0x00, 0x07, 0x80, 0xAF, 0xC8, 0x00, 0xAF, 0xC8, 0x40, 0x80, 0x00,
	0x00, 0x01, 0x06, 0x01, 0x07, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x04, 0x80, 0x00,
	0x00, 0x01, 0x06, 0x05, 0x48, 0x8F, 0xBB, 0x6C, 0x74, 0x7C, 0x3E, 0x4F, 0x78, 0x9F, 0x07, 0x8C,
	0xB3, 0x5D, 0x3C, 0x17, 0x7E, 0x45, 0x6C, 0x65, 0x6D, 0x65, 0x6E, 0x74, 0x61, 0x6C, 0x20, 0x56,
	0x69, 0x64, 0x65, 0x6F, 0x20, 0x45, 0x6E, 0x67, 0x69, 0x6E, 0x65, 0x28, 0x74, 0x6D, 0x29, 0x20,
	0x77, 0x77, 0x77, 0x2E, 0x65, 0x6C, 0x65, 0x6D, 0x65, 0x6E, 0x74, 0x61, 0x6C, 0x74, 0x65, 0x63,
	0x68, 0x6E, 0x6F, 0x6C, 0x6F, 0x67, 0x69, 0x65, 0x73, 0x2E, 0x63, 0x6F, 0x6D, 0x80, 0x00, 0x00,
}

var fuStartBytes = [...]byte{
	0xE0, 0x00, 0x00, 0x84, 0x80, 0x05, 0x21, 0x00, 0x37, 0x77, 0x41,
}

var accessUnitDelimiterBytes = [...]byte{
	0x09, 0xF0,
}

var spsBytes = [...]byte{
	0x67, 0x42, 0xC0, 0x1E, 0xB9, 0x08, 0x0D, 0x0F, 0xFC, 0x98, 0x08, 0x80, 0x00, 0x01, 0xF4, 0x80,
	0x00, 0x5D, 0xC0, 0x70, 0x30, 0x01, 0x48, 0x20, 0x02, 0x90, 0x77, 0xBD, 0xC0, 0x7C, 0x22, 0x11,
	0x92,
}

var ppsBytes = [...]byte{
	0x68, 0xDB, 0x8F, 0x20,
}

var seiBytes = [...][]byte{
	{
		0x06, 0x00, 0x07, 0x80, 0xAF, 0xC8, 0x00, 0xAF, 0xC8, 0x40,
	}, {
		0x06, 0x01, 0x07, 0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x04, 0x80,
	}, {
		0x06, 0x05, 0x48, 0x8F, 0xBB, 0x6C, 0x74, 0x7C, 0x3E, 0x4F, 0x78, 0x9F, 0x07, 0x8C, 0xB3, 0x5D,
		0x3C, 0x17, 0x7E, 0x45, 0x6C, 0x65, 0x6D, 0x65, 0x6E, 0x74, 0x61, 0x6C, 0x20, 0x56, 0x69, 0x64,
		0x65, 0x6F, 0x20, 0x45, 0x6E, 0x67, 0x69, 0x6E, 0x65, 0x28, 0x74, 0x6D, 0x29, 0x20, 0x77, 0x77,
		0x77, 0x2E, 0x65, 0x6C, 0x65, 0x6D, 0x65, 0x6E, 0x74, 0x61, 0x6C, 0x74, 0x65, 0x63, 0x68, 0x6E,
		0x6F, 0x6C, 0x6F, 0x67, 0x69, 0x65, 0x73, 0x2E, 0x63, 0x6F, 0x6D, 0x80,
	},
}
