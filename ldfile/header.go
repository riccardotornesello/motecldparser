// Package ldfile defines the binary structures for MoTeC LD file format.
//
// This package contains the low-level data structures that represent the binary
// layout of MoTeC LD files. These structures are used internally by the main
// package to serialize data into the correct format.
package ldfile

// LdFileHead represents the header structure of a MoTeC LD file.
//
// This structure defines the binary layout of the file header, which includes
// pointers to other sections of the file, device information, session metadata,
// and channel count. All multi-byte fields use little-endian byte order.
//
// The header contains several unknown fields that are required for compatibility
// with MoTeC software but whose exact purpose is not documented.
type LdFileHead struct {
	LDMarker            uint32 // 0x40
	_                   [4]byte
	ChannelsMetaPointer uint32
	ChannelsDataPointer uint32
	_                   [20]byte
	EventPointer        uint32
	_                   [24]byte
	Unknown1            uint16  // 1
	Unknown2            uint16  // 0x4240
	Unknown3            uint16  // 0xF
	DeviceSerial        uint32  // 0x1F44
	DeviceType          [8]byte // "ADL"
	DeviceVersion       uint16  // 420
	Unknown4            uint16  // 0xADB0
	ChannelsCount       uint32
	_                   [4]byte
	Date                [16]byte // "dd/MM/yyyy"
	_                   [16]byte
	Time                [16]byte // "HH:mm:ss"
	_                   [16]byte
	Driver              [64]byte
	Vehicle             [64]byte
	_                   [64]byte
	Venue               [64]byte
	_                   [64]byte
	_                   [1024]byte
	EnableProLogging    uint32 // 0xC81A4
	_                   [66]byte
	ShortComment        [64]byte
	_                   [126]byte
}
