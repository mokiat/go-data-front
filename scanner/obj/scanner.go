package obj

import (
	"errors"
	"io"

	"github.com/momchil-atanasov/go-data-front/common"
)

// MaterialLibraryEvent indicates that a material library (MTL resource)
// dependency declaration (`mtllib`) has been scanned.
type MaterialLibraryEvent struct {

	// FilePath holds the file location of the material library.
	FilePath string
}

// VertexEvent indicates that a vertex declaration (`v`) has
// been scanned.
//
// This event combines 1D, 2D, 3D, and 4D vertex declarations.
// If the original declaration did not specify one of the dimensions
// then it will be automatically defaulted by the parser (all to 0.0
// except for W which will default to 1.0)
//
// In general you should know the dimension of the entity that
// you are parsing. If not, you can always check whether all the
// vertices for a given object have the same dimension
// (e.g. all vertices have Z and W equal to 0.0 and 1.0 respectively
// which would mean 2D).
type VertexEvent struct {

	// X defines the X coordinate of this vertex.
	X float64

	// Y defines the Y coordinate of this vertex.
	Y float64

	// Z defines the Z coordinate of this vertex.
	Z float64

	// W defines the W coordinate of this vertex.
	W float64
}

// TexCoordEvent indicates that a texture coordinate declaration (`vt`) has
// been scanned.
//
// This event combines 1D, 2D, and 3D coordinate declarations.
// If the original declaration did not specify one of the dimensions
// then it will be automatically defaulted to 0.0.
//
// You can make a good guess on what the original dimension was by
// evaluating the material that is used with the texture coordinates.
//
// Another option is to check all the texture coordinates for a given
// objet. If they all have their last components defaulted, then its
// likely the coordinates were of lower dimension.
// (e.g. all texture coordinates have their W equal to 0.0 which would
// mean a 2D texture coordinate set)
type TexCoordEvent struct {

	// U defines the U coordinate of this vertex.
	U float64

	// V defines the V coordinate of this vertex.
	V float64

	// W defines the W coordinate of this vertex.
	W float64
}

// NormalEvent indicates that a normal declaration (`vn`) has
// been scanned.
type NormalEvent struct {

	// X defines the X coordinate of this normal.
	X float64

	// Y defines the Y coordinate of this normal.
	Y float64

	// Z defines the Z coordinate of this normal.
	Z float64
}

// ObjectEvent indicates that an object declaration (`o`) has
// been scanned.
type ObjectEvent struct {

	// ObjectName holds the name of the declared object
	ObjectName string
}

// MaterialReferenceEvent indicates that a material reference
// declaration (`usemtl`) has been scanned.
type MaterialReferenceEvent struct {

	// MaterialName holds the name of the material that should be
	// used for the rendering of entities that follow.
	MaterialName string
}

// FaceStartEvent indicates that a face declaration (`f`) is being
// scanned.
//
// Events that follow will concern that specific face up until the
// point a FaceEndEvent is thrown.
type FaceStartEvent struct {
}

// FaceEndEvent indicates that the scanning of the face declaration
// has completed.
type FaceEndEvent struct {
}

// ReferenceSetStartEvent indicates that a reference set (e.g. 1/2/3)
// is being parsed as part of face scanning.
//
// Events that follow will contain reference data, until a
// ReferenceSetEndEvent is received.
type ReferenceSetStartEvent struct {
}

// ReferenceSetEndEvent indicates that the scanning of the reference
// set has completed.
type ReferenceSetEndEvent struct {
}

// VertexReferenceEvent indicates that a vertex reference has been
// scanned.
type VertexReferenceEvent struct {

	// VertexIndex holds the index of the vertex that is referenced
	VertexIndex int64
}

// TexCoordReferenceEvent indicates that a texture coordinate reference
// has been scanned.
type TexCoordReferenceEvent struct {

	// TexCoordIndex holds the index of the texture coordinate that
	// is referenced.
	TexCoordIndex int64
}

// NormalReferenceEvent indicates that a normal reference has been
// scanned.
type NormalReferenceEvent struct {

	// NormalIndex holds the index of the normal that is referenced.
	NormalIndex int64
}

// NewScanner creates a new Scanner object that can scan through
// Wavefront OBJ resources.
func NewScanner() common.Scanner {
	return &scanner{}
}

type scanner struct {
}

func (s *scanner) Scan(reader io.Reader, handler common.EventHandler) error {
	lineScanner := common.NewLineScanner(reader)

	var err error
	for lineScanner.Scan() {
		line := lineScanner.Line()
		switch {
		case line.IsBlank():
			break
		case line.IsComment():
			err = s.processComment(line, handler)
			break
		case line.IsCommand():
			err = s.processCommand(line, handler)
			break
		}
		if err != nil {
			return err
		}
	}

	if lineScanner.Err() != nil {
		return lineScanner.Err()
	}
	return nil
}

func (s *scanner) processComment(line common.Line, handler common.EventHandler) error {
	event := common.CommentEvent{
		Comment: line.Comment(),
	}
	return handler(event)
}

func (s *scanner) processCommand(line common.Line, handler common.EventHandler) error {
	switch {
	case line.HasCommandName("mtllib"):
		return s.processMaterialLibrary(line, handler)
	case line.HasCommandName("v"):
		return s.processVertex(line, handler)
	case line.HasCommandName("vt"):
		return s.processTexCoord(line, handler)
	case line.HasCommandName("vn"):
		return s.processNormal(line, handler)
	case line.HasCommandName("o"):
		return s.processObject(line, handler)
	case line.HasCommandName("usemtl"):
		return s.processMaterialReference(line, handler)
	case line.HasCommandName("f"):
		return s.processFace(line, handler)
	default:
		return nil
	}
}

func (s *scanner) processMaterialLibrary(line common.Line, handler common.EventHandler) error {
	for i := 0; i < line.ParamCount(); i++ {
		path := line.StringParam(i)
		event := MaterialLibraryEvent{
			FilePath: path,
		}
		err := handler(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *scanner) processVertex(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 3 {
		return errors.New("Insufficient vertex data.")
	}
	var err error
	event := VertexEvent{
		X: 0.0, Y: 0.0, Z: 0.0, W: 1.0,
	}
	event.X, err = line.FloatParam(0)
	if err != nil {
		return err
	}
	event.Y, err = line.FloatParam(1)
	if err != nil {
		return err
	}
	event.Z, err = line.FloatParam(2)
	if err != nil {
		return err
	}
	if line.ParamCount() >= 4 {
		event.W, err = line.FloatParam(3)
		if err != nil {
			return err
		}
	}
	return handler(event)
}

func (s *scanner) processTexCoord(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() == 0 {
		return errors.New("Insufficient texture coordinate data.")
	}

	var err error
	event := TexCoordEvent{
		U: 0.0, V: 0.0, W: 0.0,
	}
	event.U, err = line.FloatParam(0)
	if err != nil {
		return err
	}
	if line.ParamCount() >= 2 {
		event.V, err = line.FloatParam(1)
		if err != nil {
			return err
		}
	}
	if line.ParamCount() >= 3 {
		event.W, err = line.FloatParam(2)
		if err != nil {
			return err
		}
	}
	return handler(event)
}

func (s *scanner) processNormal(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 3 {
		return errors.New("Insufficient normal data.")
	}
	var err error
	event := NormalEvent{
		X: 0.0, Y: 0.0, Z: 0.0,
	}
	event.X, err = line.FloatParam(0)
	if err != nil {
		return err
	}
	event.Y, err = line.FloatParam(1)
	if err != nil {
		return err
	}
	event.Z, err = line.FloatParam(2)
	if err != nil {
		return err
	}
	return handler(event)
}

func (s *scanner) processObject(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() == 0 {
		return errors.New("No name specified for object.")
	}
	name := line.StringParam(0)
	event := ObjectEvent{
		ObjectName: name,
	}
	return handler(event)
}

func (s *scanner) processMaterialReference(line common.Line, handler common.EventHandler) error {
	event := MaterialReferenceEvent{}
	if line.ParamCount() > 0 {
		event.MaterialName = line.StringParam(0)
	}
	return handler(event)
}

func (s *scanner) processFace(line common.Line, handler common.EventHandler) error {
	err := handler(FaceStartEvent{})
	if err != nil {
		return err
	}

	for i := 0; i < line.ParamCount(); i++ {
		err := s.processReferenceSet(line.ReferenceSetParam(i), handler)
		if err != nil {
			return err
		}
	}

	return handler(FaceEndEvent{})
}

func (s *scanner) processReferenceSet(referenceSet common.ReferenceSet, handler common.EventHandler) error {
	err := handler(ReferenceSetStartEvent{})
	if err != nil {
		return err
	}

	if referenceSet.Count() == 0 {
		return errors.New("Reference set has no references.")
	}

	vertexIndex, err := referenceSet.IntReference(0)
	if err != nil {
		return err
	}
	handler(VertexReferenceEvent{
		VertexIndex: vertexIndex,
	})

	if (referenceSet.Count() > 1) && !referenceSet.IsBlank(1) {
		texCoordIndex, err := referenceSet.IntReference(1)
		if err != nil {
			return err
		}
		handler(TexCoordReferenceEvent{
			TexCoordIndex: texCoordIndex,
		})
	}

	if (referenceSet.Count() > 2) && !referenceSet.IsBlank(2) {
		normalIndex, err := referenceSet.IntReference(2)
		if err != nil {
			return err
		}
		handler(NormalReferenceEvent{
			NormalIndex: normalIndex,
		})
	}

	return handler(ReferenceSetEndEvent{})
}
