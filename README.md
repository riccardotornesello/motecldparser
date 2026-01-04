# motec-ld-go

A Go library for reading and writing MoTeC LD (Logged Data) files.

## Overview

This library provides a simple and efficient way to create and write MoTeC LD files in Go. MoTeC LD files are a binary format used by MoTeC data acquisition systems to store telemetry data from racing vehicles, including channel data such as speed, RPM, temperatures, and other sensor measurements.

## Features

- Write MoTeC LD files with full metadata support
- Support for multiple data types (float32, int16, int32)
- Type-safe channel definitions using Go generics
- Configurable channel properties (frequency, units, names)
- Event, venue, and vehicle metadata support

## Installation

```bash
go get github.com/riccardotornesello/motecldparser
```

## Usage

### Basic Example

```go
package main

import (
    "os"
    "time"
    motec "github.com/riccardotornesello/motecldparser"
)

func main() {
    // Create a new LD file
    file := motec.File{
        Time:         time.Now(),
        Driver:       "John Doe",
        Vehicle:      "Race Car #42",
        Venue:        "Silverstone Circuit",
        ShortComment: "Qualifying session",
        
        EventName:    "Grand Prix",
        EventSession: "Q1",
        EventComment: "Qualifying session 1",
        
        VehicleId:      "CAR42",
        VehicleWeight:  1450,
        VehicleType:    "GT3",
        VehicleComment: "Race spec",
    }
    
    // Create channels with data
    speed := &motec.Channel[float32]{
        Frequency: 100,
        Name:      "Speed",
        ShortName: "SPD",
        Unit:      "km/h",
        Data:      &[]float32{0.0, 10.5, 25.3, 50.2, 75.8},
    }
    
    rpm := &motec.Channel[int16]{
        Frequency: 100,
        Name:      "Engine RPM",
        ShortName: "RPM",
        Unit:      "rpm",
        Data:      &[]int16{800, 1200, 2500, 4000, 6500},
    }
    
    // Add channels to file
    file.AddChannels(speed, rpm)
    
    // Write to file
    fd, err := os.Create("output.ld")
    if err != nil {
        panic(err)
    }
    defer fd.Close()
    
    file.Write(fd)
}
```

### Adding Data Incrementally

You can add data to channels incrementally:

```go
channel := &motec.Channel[float32]{
    Frequency: 100,
    Name:      "Temperature",
    ShortName: "TEMP",
    Unit:      "°C",
    Data:      &[]float32{},
}

// Add data points one at a time
channel.AddData(20.5)
channel.AddData(21.3)
channel.AddData(22.1)
```

## API Reference

### File Structure

The `File` struct contains all metadata and channels for the MoTeC LD file:

```go
type File struct {
    Time         time.Time  // Timestamp of the data logging session
    Driver       string     // Driver name
    Vehicle      string     // Vehicle identifier
    Venue        string     // Venue/track name
    ShortComment string     // Short description
    
    EventName    string     // Event name
    EventSession string     // Session identifier (e.g., "Q1", "Race")
    EventComment string     // Event description
    
    VehicleId      string   // Vehicle ID
    VehicleWeight  uint32   // Vehicle weight in kg
    VehicleType    string   // Vehicle type/class
    VehicleComment string   // Vehicle notes
    
    Channels []interface{}  // List of channels
}
```

### Channel Structure

Channels are type-safe and support three data types:

```go
type Channel[T float32 | int16 | int32] struct {
    Frequency uint16  // Sampling frequency in Hz
    Name      string  // Full channel name
    ShortName string  // Abbreviated name (max 8 characters)
    Unit      string  // Unit of measurement
    Data      *[]T    // Pointer to data array
}
```

### Methods

#### File.Write

```go
func (f *File) Write(fd *os.File)
```

Writes the complete MoTeC LD file to the provided file descriptor.

#### File.AddChannels

```go
func (f *File) AddChannels(channels ...interface{})
```

Adds one or more channels to the file.

#### Channel.AddData

```go
func (c *Channel[T]) AddData(data T)
```

Appends a single data point to the channel.

## File Format

The library writes MoTeC LD files with the following structure:

```
|---------------|
|    HEADER     |
|---------------| <- EVENT_POINTER
|    EVENT      |
|---------------| <- VENUE_POINTER
|    VENUE      |
|---------------| <- VEHICLE_POINTER
|    VEHICLE    |
|---------------| <- CHANNELS_META_POINTER
|  CHANNEL HDR  |
|  CHANNEL HDR  |
|---------------| <- CHANNELS_DATA_POINTER
|  CHANNEL DATA |
|  CHANNEL DATA |
|---------------|
```

## Supported Data Types

- `float32` - 32-bit floating point values
- `int16` - 16-bit signed integers
- `int32` - 32-bit signed integers

## Status and Contributions

This library is production-ready for writing MoTeC LD files. Reading functionality is not currently implemented.

Contributions are welcome! Please feel free to:
- Report issues
- Suggest enhancements
- Submit pull requests for new features or improvements

## License

This project is available under the MIT License unless otherwise specified.

## Acknowledgments

This library implements the MoTeC LD file format specification for data logging and telemetry storage.
