// Package motec_ld_parser provides functionality for writing MoTeC LD (Logged Data) files.
//
// MoTeC LD files are binary files used by MoTeC data acquisition systems to store
// telemetry data from racing vehicles. This package supports creating and writing
// LD files with multiple channels of different data types (float32, int16, int32).
//
// Basic usage:
//
//	file := motec_ld_parser.File{
//	    Time:    time.Now(),
//	    Driver:  "Driver Name",
//	    Vehicle: "Vehicle Name",
//	}
//	channel := &motec_ld_parser.Channel[float32]{
//	    Frequency: 100,
//	    Name:      "Speed",
//	    Data:      &[]float32{0, 10, 20},
//	}
//	file.AddChannels(channel)
//	file.Write(fileDescriptor)
package motec_ld_parser

import (
	"bytes"
	"encoding/binary"
	"os"
	"time"

	"github.com/riccardotornesello/motecldparser/ldfile"
)

/*
	|---------------|
	|	HEAD		|
	|---------------| <- EVENT_POINTER
	|	EVENT		|
	|---------------| <- VENUE_POINTER
	|	VENUE		|
	|---------------| <- VEHICLE_POINTER
	|	VEHICLE		|
	|---------------| <- CHANNELS_META_POINTER
	|	CHANNEL H	|
	|	CHANNEL H	|
	|---------------| <- CHANNELS_DATA_POINTER
	|	CHANNEL D	|
	|	CHANNEL D	|
	|---------------|
*/

// File represents a MoTeC LD file containing telemetry data and metadata.
//
// The File structure holds all the information needed to create a complete MoTeC LD file,
// including session metadata (time, driver, venue), event details, vehicle information,
// and a collection of data channels.
//
// All string fields are limited in length when written to the binary file format:
//   - Driver, Vehicle, Venue: max 64 bytes
//   - ShortComment: max 64 bytes
//   - EventName, EventSession: max 64 bytes
//   - EventComment: max 1024 bytes
//   - VehicleId: max 64 bytes
//   - VehicleType, VehicleComment: max 32 bytes
type File struct {
	Time         time.Time // Timestamp of when the data was logged
	Driver       string    // Name of the driver
	Vehicle      string    // Vehicle identifier or name
	Venue        string    // Track or venue name
	ShortComment string    // Brief description or notes

	EventName    string // Name of the event (e.g., "Grand Prix")
	EventSession string // Session identifier (e.g., "Q1", "Race", "Practice")
	EventComment string // Detailed event description or notes

	VehicleId      string // Unique vehicle identifier
	VehicleWeight  uint32 // Vehicle weight in kilograms
	VehicleType    string // Vehicle type or class (e.g., "GT3", "Formula")
	VehicleComment string // Additional vehicle notes

	Channels []interface{} // Collection of Channel pointers (use AddChannels to add)
}

// Channel represents a single data channel in a MoTeC LD file.
//
// A channel contains a series of measurements sampled at a specific frequency.
// The type parameter T specifies the data type and must be one of:
//   - float32: for floating-point values (speed, temperature, etc.)
//   - int16: for 16-bit integer values
//   - int32: for 32-bit integer values
//
// String field limits when written to the binary file:
//   - Name: max 32 bytes
//   - ShortName: max 8 bytes
//   - Unit: max 12 bytes
//
// Example:
//
//	speedChannel := &Channel[float32]{
//	    Frequency: 100,      // 100 Hz sampling rate
//	    Name:      "Speed",
//	    ShortName: "SPD",
//	    Unit:      "km/h",
//	    Data:      &[]float32{0.0, 10.5, 25.3},
//	}
type Channel[T float32 | int16 | int32] struct {
	Frequency uint16 // Sampling frequency in Hz
	Name      string // Full channel name
	ShortName string // Abbreviated name (displayed in compact views)
	Unit      string // Unit of measurement (e.g., "km/h", "rpm", "°C")
	Data      *[]T   // Pointer to the data array
}

// Write writes the complete MoTeC LD file to the provided file descriptor.
//
// This method serializes all file metadata, event information, vehicle details,
// and channel data into the MoTeC LD binary format and writes it to the file.
// The file must be opened for writing before calling this method.
//
// The method handles:
//   - Computing all internal pointers for the binary structure
//   - Writing the file header with session metadata
//   - Writing event, venue, and vehicle information blocks
//   - Writing channel metadata and data for all channels
//
// Example:
//
//	fd, err := os.Create("telemetry.ld")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer fd.Close()
//	file.Write(fd)
//
// Note: This method does not return an error. Any write errors will cause a panic.
// Consider wrapping file operations in appropriate error handling.
func (f *File) Write(fd *os.File) {
	// Calculate pointers
	headerSize := uintptr(binary.Size(ldfile.LdFileHead{}))
	eventSize := uintptr(binary.Size(ldfile.LdFileEvent{}))
	venueSize := uintptr(binary.Size(ldfile.LdFileVenue{}))
	vehicleSize := uintptr(binary.Size(ldfile.LdFileVehicle{}))
	channelMetaSize := uintptr(binary.Size(ldfile.LdFileChannelMeta{}))

	eventPointer := headerSize
	venuePointer := eventPointer + eventSize
	vehiclePointer := venuePointer + venueSize
	channelsMetaPointer := vehiclePointer + vehicleSize
	channelsDataPointer := channelsMetaPointer + channelMetaSize*uintptr(len(f.Channels))

	// Create the file header
	head := ldfile.LdFileHead{
		LDMarker:         0x40,
		Unknown1:         1,
		Unknown2:         0x4240,
		Unknown3:         0xF,
		Unknown4:         0xADB0,
		DeviceSerial:     0x1F44,
		DeviceType:       [8]byte{'A', 'D', 'L', 0, 0, 0, 0, 0},
		DeviceVersion:    420,
		EnableProLogging: 0xC81A4,
		ChannelsCount:    uint32(len(f.Channels)),

		EventPointer:        uint32(eventPointer),
		ChannelsMetaPointer: uint32(channelsMetaPointer),
		ChannelsDataPointer: uint32(channelsDataPointer),
	}

	date := f.Time.Format("02/01/2006")
	hour := f.Time.Format("15:04:05")
	copy(head.Date[:], date)
	copy(head.Time[:], hour)

	copy(head.Driver[:], f.Driver)
	copy(head.Vehicle[:], f.Vehicle)
	copy(head.Venue[:], f.Venue)
	copy(head.ShortComment[:], f.ShortComment)

	// Create the Event
	event := ldfile.LdFileEvent{
		VenuePointer: uint16(venuePointer),
	}

	copy(event.Name[:], f.EventName)
	copy(event.Session[:], f.EventSession)
	copy(event.Comment[:], f.EventComment)

	// Create the Venue
	venue := ldfile.LdFileVenue{
		VehiclePointer: uint16(vehiclePointer),
	}

	copy(venue.Name[:], f.Venue)

	// Create the Vehicle
	vehicle := ldfile.LdFileVehicle{
		Weight: f.VehicleWeight,
	}

	copy(vehicle.Id[:], f.VehicleId)
	copy(vehicle.Type[:], f.VehicleType)
	copy(vehicle.Comment[:], f.VehicleComment)

	// Write to file
	binary.Write(fd, binary.LittleEndian, head)

	fd.Seek(int64(eventPointer), 0)
	binary.Write(fd, binary.LittleEndian, event)

	fd.Seek(int64(venuePointer), 0)
	binary.Write(fd, binary.LittleEndian, venue)

	fd.Seek(int64(vehiclePointer), 0)
	binary.Write(fd, binary.LittleEndian, vehicle)

	// Write channels
	currentDataPointer := channelsDataPointer
	for i, channel := range f.Channels {
		switch any(channel).(type) {
		case *Channel[float32]:
			currentDataPointer = channel.(*Channel[float32]).Write(fd, uint16(i), head.ChannelsCount, channelsMetaPointer, currentDataPointer)
			break
		case *Channel[int16]:
			currentDataPointer = channel.(*Channel[int16]).Write(fd, uint16(i), head.ChannelsCount, channelsMetaPointer, currentDataPointer)
			break
		case *Channel[int32]:
			currentDataPointer = channel.(*Channel[int32]).Write(fd, uint16(i), head.ChannelsCount, channelsMetaPointer, currentDataPointer)
			break
		}
	}
}

// AddChannels adds one or more channels to the file.
//
// Channels must be pointers to Channel instances with appropriate type parameters.
// Multiple channels can be added in a single call.
//
// Example:
//
//	speedChannel := &Channel[float32]{...}
//	rpmChannel := &Channel[int16]{...}
//	file.AddChannels(speedChannel, rpmChannel)
func (f *File) AddChannels(channels ...interface{}) {
	f.Channels = append(f.Channels, channels...)
}

// Write writes a single channel's metadata and data to the file.
//
// This method is called internally by File.Write for each channel.
// It serializes the channel's metadata (name, unit, frequency, etc.) and
// binary data to the appropriate locations in the file.
//
// Parameters:
//   - fd: the file descriptor to write to
//   - n: the channel index (0-based)
//   - channelsCount: total number of channels in the file
//   - channelsMetaPointer: file offset where channel metadata begins
//   - currentDataPointer: file offset where this channel's data should be written
//
// Returns the file offset for the next channel's data.
//
// This method should not typically be called directly by users.
func (c *Channel[T]) Write(
	fd *os.File,
	n uint16,
	channelsCount uint32,
	channelsMetaPointer uintptr,
	currentDataPointer uintptr,
) uintptr {
	var dataType ldfile.DataType
	var previousMetaPointer uintptr = 0
	var nextMetaPointer uintptr = 0

	switch any(c).(type) {
	case *Channel[float32]:
		dataType = ldfile.DataTypeFloat32
		break
	case *Channel[int16]:
		dataType = ldfile.DataTypeInt16
		break
	case *Channel[int32]:
		dataType = ldfile.DataTypeInt32
		break
	}

	if n > 0 {
		previousMetaPointer = channelsMetaPointer + uintptr(binary.Size(ldfile.LdFileChannelMeta{}))*(uintptr(n-1))
	}

	if n < uint16(channelsCount-1) {
		nextMetaPointer = channelsMetaPointer + uintptr(binary.Size(ldfile.LdFileChannelMeta{}))*(uintptr(n+1))
	}

	currentMetaPointer := channelsMetaPointer + uintptr(binary.Size(ldfile.LdFileChannelMeta{}))*uintptr(n)

	channelMeta := ldfile.LdFileChannelMeta{
		PreviousMetaPointer: uint32(previousMetaPointer),
		NextMetaPointer:     uint32(nextMetaPointer),
		DataPointer:         uint32(currentDataPointer),
		DataLength:          uint32(len(*c.Data)),
		ChannelId:           0x2EE1 + n,
		DataType:            dataType.DataType,
		DataTypeLength:      dataType.DataTypeLength,
		Frequency:           c.Frequency,
		Shift:               0,
		Mul:                 1,
		Scale:               1,
		DecPlaces:           0,
	}

	copy(channelMeta.Name[:], c.Name)
	copy(channelMeta.ShortName[:], c.ShortName)
	copy(channelMeta.Unit[:], c.Unit)

	// Convert data to binary slice
	binaryDataWriter := new(bytes.Buffer)
	binary.Write(binaryDataWriter, binary.LittleEndian, c.Data)
	binaryData := binaryDataWriter.Bytes()

	// Write to file
	fd.Seek(int64(currentMetaPointer), 0)
	binary.Write(fd, binary.LittleEndian, channelMeta)

	fd.Seek(int64(currentDataPointer), 0)
	binary.Write(fd, binary.LittleEndian, binaryData)

	// Return next data pointer
	nextDataPointer := currentDataPointer + uintptr(len(binaryData))
	return nextDataPointer
}

// AddData appends a single data point to the channel.
//
// This is a convenience method for adding data incrementally rather than
// providing all data at once.
//
// Example:
//
//	channel := &Channel[float32]{
//	    Name: "Temperature",
//	    Data: &[]float32{},
//	}
//	channel.AddData(20.5)
//	channel.AddData(21.3)
//	channel.AddData(22.1)
func (c *Channel[T]) AddData(data T) {
	*c.Data = append(*c.Data, data)
}
