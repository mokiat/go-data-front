/*
This package provides APIs through which one can parse Wavefront OBJ resources.
The OBJ file format is quite common and is used to store 3D model data.

The parsers provided by this library do not support the full OBJ spec, which
is quite lengthy, but rather only the most essential and common aspects.
*/
package obj

import "io"

/*
An implementation of this interface should be passed to the constructor
of the `Scanner` interface.

Method of this interface will be called during scanning.
*/
type ScannerHandler interface {

	// This method is called when a comment section has been parsed.
	OnComment(comment string) error

	// This method is called when a material library dependency is parsed.
	OnMaterialLibrary(path string) error

	// This method is called when a material reference has been parsed.
	OnMaterialReference(name string) error

	// This method is called when a vertex is about to be parsed.
	OnVertexStart() error

	// This method is called when the X coordinate of a vertex is parsed.
	OnVertexX(x float32) error

	// This method is called when the Y coordinate of a vertex is parsed.
	OnVertexY(y float32) error

	// This method is called when the Z coordinate of a vertex is parsed.
	OnVertexZ(z float32) error

	// This method is called when the W coordinate of a vertex is parsed.
	OnVertexW(w float32) error

	// This method is called when the parsing of a given vertex has finished.
	OnVertexEnd() error

	// This method is called when a texture coordinate is about to be parsed.
	OnTexCoordStart() error

	// This method is called when the U coordinate of a texture coordinate is parsed.
	OnTexCoordU(u float32) error

	//  This method is called when the V coordinate of a texture coordinate is parsed.
	OnTexCoordV(v float32) error

	//  This method is called when the W coordinate of a texture coordinate is parsed.
	OnTexCoordW(w float32) error

	// This method is called when a texture coordinate has been fully parsed.
	OnTexCoordEnd() error

	// This method is called when a normal has been parsed.
	OnNormal(x, y, z float32) error

	// This method is called when a new object declaration is parsed.
	OnObject(name string) error

	// This method is called when a new face is about to be parsed.
	OnFaceStart() error

	// This method is called when a new coord reference is about to be parsed.
	OnCoordReferenceStart() error

	// This method is called when a vertex index reference is parsed.
	OnVertexIndex(index int) error

	// This method is called when a texture coordinate index reference is parsed.
	OnTexCoordIndex(index int) error

	// This method is called when a normal index reference is parsed.
	OnNormalIndex(index int) error

	// This method is called when the parsing of a coord reference has finished.
	OnCoordReferenceEnd() error

	// This method is called when a face has been fully parsed.
	OnFaceEnd() error
}

/*
The Scanner interface represents an event-based parser. The model data
is iterated - an interesting element at a time - and an event is thrown.

Events are sent to the user via method invocations on the ScannerHandler
interface that is previously specified.
*/
type Scanner interface {

	// Reads an OBJ resource from the specified io.Reader stream.
	// An error is returned should parsing fail for some reason.
	Scan(io.Reader) error
}

/*
Creates a new Scanner using the specified ScannerHandler for callback
handling.
*/
func NewScanner(handler ScannerHandler) Scanner {
	return &scanner{
		handler: handler,
	}
}
