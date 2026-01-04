package ldfile

// LdFileEvent represents the event information structure in a MoTeC LD file.
//
// This structure contains metadata about the logging event, including the
// event name, session identifier, and descriptive comments. It also includes
// a pointer to the venue information structure.
type LdFileEvent struct {
	Name         [64]byte
	Session      [64]byte
	Comment      [1024]byte
	VenuePointer uint16
}
