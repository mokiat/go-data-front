package common

import "io"

// Scanner represents an event-based parser for Wavefront resources.
//
// Implementations of this interface would usually scan through the
// Wavefront resource and act on each special element that is
// detected (e.g. when a normal declaration is parsed)
//
// This interface does not indicate the exact way through which the
// user will be notified of such events and that should be checked
// in the documentation of the actual implementation.
// Most common, there will be an interface that needs to be implemented
// by the user and passed to the constructor of the Scanner. A method
// of the interface would be invoked for each interesting element that
// is parsed.
type Scanner interface {

	// Performs a scan through the Wavefront resource that is
	// provided through the io.Reader.
	//
	// An error is returned should parsing fail for some reason.
	Scan(io.Reader) error
}
