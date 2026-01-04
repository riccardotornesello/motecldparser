package ldfile

// LdFileVenue represents the venue (track/location) information structure.
//
// This structure contains the venue name and a pointer to the vehicle
// information structure. The venue typically refers to the racing circuit
// or location where the data was logged.
type LdFileVenue struct {
	Name           [64]byte
	_              [1034]byte
	VehiclePointer uint16
}
