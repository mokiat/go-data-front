package common

import "io"

// Event represents an event that has occurred during scanning.
// This is the mechanism through which the Scanner API returns
// key MTL resource elements to the user.
//
// This interface is subclassed and users should cast this interface
// to one of them in order to extract meaningful data.
//
// Example event types:
//   - CommentEvent
//   - obj.ObjectEvent
//   - mtl.MaterialEvent
type Event interface{}

// CommentEvent indicates that a comment section has been
// scanned through.
type CommentEvent struct {

	// Comment holds the comment text that has been scanned.
	Comment string
}

// EventHandler function is passed to the Scanner by the API user
// in order to receive scanning events.
//
// Implementations of this function should do a type switch on the
// event in order to determine the exact type and read meaningful
// data from it.
//
// Implementations can return an error in order to stop any further
// processing of the Wavefront file.
type EventHandler func(event Event) error

// Scanner represents an event-based parser for Wavefront resources.
//
// Implementations of this interface would usually scan through the
// Wavefront resource and act on each special element that is
// detected (e.g. when a normal or material declaration is parsed)
type Scanner interface {

	// Scan performs a scan through the Wavefront resource that is
	// provided through the io.Reader.
	//
	// The EventHandler parameter is used by the implementation to
	// pass scanning events back to the user for processing.
	//
	// An error is returned should parsing fail for some reason or
	// if the user returns an error via the EventHandler.
	Scan(io.Reader, EventHandler) error
}
