package ldfile

// DataType represents the data type encoding for channel data in MoTeC LD files.
//
// The DataType field specifies the type of data (float, int16, int32), and
// DataTypeLength specifies the size in bytes of each data point.
type DataType struct {
	DataType       uint16 // Type identifier (0x07 for float, 0x03 for int16, 0x05 for int32)
	DataTypeLength uint16 // Size in bytes (2 or 4)
}

// Predefined data type constants for use in channel metadata.
var (
	DataTypeFloat16 = DataType{0x07, 2} // 16-bit floating point (2 bytes)
	DataTypeFloat32 = DataType{0x07, 4} // 32-bit floating point (4 bytes)
	DataTypeInt16   = DataType{0x03, 2} // 16-bit signed integer (2 bytes)
	DataTypeInt32   = DataType{0x05, 4} // 32-bit signed integer (4 bytes)
)
