/*
obj package provides APIs through which one can parse Wavefront OBJ resources.
The OBJ file format is quite common and is used to store 3D model data.

The parsers provided by this library do not support the full OBJ spec, which
is quite lengthy, but rather only the most essential and common aspects.
*/
package obj

import "io"

/*
ScannerHandler interface needs to be implemented by users and passed
to the Scanner construtor so that the user may receive parsing events.
*/
type ScannerHandler interface {

	// OnComment is called when a comment section has been parsed.
	OnComment(comment string) error

	// OnMaterialLibrary is called when a material library dependency is parsed.
	OnMaterialLibrary(path string) error

	// OnMaterialReference is called when a material reference has been parsed.
	OnMaterialReference(name string) error

	// OnVertexStart is called when a vertex is about to be parsed.
	OnVertexStart() error

	// OnVertexX is called when the X coordinate of a vertex is parsed.
	OnVertexX(x float32) error

	// OnVertexY is called when the Y coordinate of a vertex is parsed.
	OnVertexY(y float32) error

	// OnVertexZ is called when the Z coordinate of a vertex is parsed.
	OnVertexZ(z float32) error

	// OnVertexW is called when the W coordinate of a vertex is parsed.
	OnVertexW(w float32) error

	// OnVertexEnd is called when the parsing of a given vertex has finished.
	OnVertexEnd() error

	// OnTexCoordStart is called when a texture coordinate is about to be parsed.
	OnTexCoordStart() error

	// OnTexCoordU is called when the U coordinate of a texture coordinate is parsed.
	OnTexCoordU(u float32) error

	// OnTexCoordV is called when the V coordinate of a texture coordinate is parsed.
	OnTexCoordV(v float32) error

	// OnTexCoordW is called when the W coordinate of a texture coordinate is parsed.
	OnTexCoordW(w float32) error

	// OnTexCoordEnd is called when a texture coordinate has been fully parsed.
	OnTexCoordEnd() error

	// OnNormal is called when a normal has been parsed.
	OnNormal(x, y, z float32) error

	// OnObject is called when a new object declaration is parsed.
	OnObject(name string) error

	// OnFaceStart is called when a new face is about to be parsed.
	OnFaceStart() error

	// OnCoordReferenceStart is called when a new coord reference is about to be parsed.
	OnCoordReferenceStart() error

	// OnVertexIndex is called when a vertex index reference is parsed.
	OnVertexIndex(index int) error

	// OnTexCoordIndex is called when a texture coordinate index reference is parsed.
	OnTexCoordIndex(index int) error

	// OnNormalIndex is called when a normal index reference is parsed.
	OnNormalIndex(index int) error

	// OnCoordReferenceEnd is called when the parsing of a coord reference has finished.
	OnCoordReferenceEnd() error

	// OnFaceEnd is called when a face has been fully parsed.
	OnFaceEnd() error
}

/*
Scanner interface represents an event-based parser. The model data
is iterated - an interesting element at a time - and an event is thrown.

Events are sent to the user via method invocations on the ScannerHandler
interface that is previously specified.
*/
type Scanner interface {

	// Scan parses an OBJ resource from the specified io.Reader stream.
	// An error is returned should parsing fail for some reason.
	Scan(io.Reader) error
}

/*
NewScanner creates a new Scanner using the specified ScannerHandler
for callback handling.
*/
func NewScanner(handler ScannerHandler) Scanner {
	return &scanner{
		handler: handler,
	}
}
