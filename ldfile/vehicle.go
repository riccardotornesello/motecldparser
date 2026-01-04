package ldfile

// LdFileVehicle represents the vehicle information structure in a MoTeC LD file.
//
// This structure contains detailed information about the vehicle, including
// its unique identifier, weight, type/class, and additional comments.
type LdFileVehicle struct {
	Id      [64]byte
	_       [128]byte
	Weight  uint32
	Type    [32]byte
	Comment [32]byte
}
