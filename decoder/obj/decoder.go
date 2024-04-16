package obj

import (
	"fmt"
	"io"

	"github.com/mokiat/go-data-front/common"
	objscan "github.com/mokiat/go-data-front/scanner/obj"
)

// DecodeLimits specifies restrictions on parsing.
//
// One will generally use this to limit the number of
// data that is parsed in order to prevent out of memory
// errors
type DecodeLimits struct {

	// MaxVertexCount specifies the maximum number of vertices
	// that can be parsed before an error is thrown.
	MaxVertexCount int

	// MaxTexCoordCount specifies the maximum number of texture
	// coordinates that can be parsed before an error is thrown.
	MaxTexCoordCount int

	// MaxNormalCount specifies the maximum number of normals
	// that can be parsed before an error is thrown.
	MaxNormalCount int

	// MaxObjectCount specifies the maximum number of objects
	// that can be parsed before an error is thrown.
	MaxObjectCount int

	// MaxFaceCount specifies the maximum number of faces
	// that can be parsed per mesh before an error is thrown.
	MaxFaceCount int

	// MaxReferenceCount specifies the maximum number of vertex
	// references that a given face can have.
	MaxReferenceCount int

	// MaxMaterialLibraryCount specifies the maximum number of
	// material library references that can be parsed before
	// an error is thrown.
	MaxMaterialLibraryCount int

	// MaxMaterialReferenceCount specifies the maximum number of
	// material references that can be parsed per object before
	// an error is thrown.
	MaxMaterialReferenceCount int
}

// DefaultLimits returns some default DecodeLimits.
// Users can take the result and modify specific parameters.
func DefaultLimits() DecodeLimits {
	return DecodeLimits{
		MaxVertexCount:            65536,
		MaxTexCoordCount:          65536,
		MaxNormalCount:            65536,
		MaxObjectCount:            1024,
		MaxFaceCount:              65536,
		MaxReferenceCount:         16,
		MaxMaterialReferenceCount: 64,
		MaxMaterialLibraryCount:   32,
	}
}

// Decoder is an API that allows one to decode OBJ
// Wavefront resources into an object model.
type Decoder interface {

	// Decode decodes the OBJ Wavefront resource, specified
	// through the io.Reader, into a Library model.
	//
	// If decoding fails for some reason, an error is returned.
	Decode(io.Reader) (*Model, error)
}

// NewDecoder creates a new Decoder instance with the
// specified DecodeLimits.
func NewDecoder(limits DecodeLimits) Decoder {
	return &decoder{
		limits:  &limits,
		scanner: objscan.NewScanner(),
	}
}

type decoder struct {
	limits  *DecodeLimits
	scanner common.Scanner
}

func (d *decoder) Decode(reader io.Reader) (*Model, error) {
	context := newDecodeContext(d.limits)
	err := d.scanner.Scan(reader, context.HandleEvent)
	if err != nil {
		return nil, err
	}
	return context.Model(), nil
}

func newDecodeContext(limits *DecodeLimits) *decodeContext {
	return &decodeContext{
		limits:        limits,
		model:         new(Model),
		currentObject: nil,
		currentMesh:   nil,
		currentFace:   nil,
	}
}

type decodeContext struct {
	limits           *DecodeLimits
	model            *Model
	currentObject    *Object
	currentMesh      *Mesh
	currentFace      *Face
	currentReference *Reference
}

func (c *decodeContext) Model() *Model {
	return c.model
}

func (c *decodeContext) HandleEvent(event common.Event) error {
	switch actual := event.(type) {
	case objscan.MaterialLibraryEvent:
		return c.handleMaterialLibrary(actual)
	case objscan.VertexEvent:
		return c.handleVertex(actual)
	case objscan.NormalEvent:
		return c.handleNormal(actual)
	case objscan.TexCoordEvent:
		return c.handleTexCoord(actual)
	case objscan.ObjectEvent:
		return c.handleObject(actual)
	case objscan.MaterialReferenceEvent:
		return c.handleMaterialReference(actual)
	case objscan.FaceStartEvent:
		return c.handleFaceStart()
	case objscan.FaceEndEvent:
		return c.handleFaceEnd()
	case objscan.ReferenceSetStartEvent:
		return c.handleReferencesStart()
	case objscan.ReferenceSetEndEvent:
		return c.handleReferencesEnd()
	case objscan.VertexReferenceEvent:
		return c.handleVertexReference(actual)
	case objscan.TexCoordReferenceEvent:
		return c.handleTexCoordReference(actual)
	case objscan.NormalReferenceEvent:
		return c.handleNormalReference(actual)
	}
	return nil
}

func (c *decodeContext) handleMaterialLibrary(event objscan.MaterialLibraryEvent) error {
	if len(c.model.MaterialLibraries) >= c.limits.MaxMaterialLibraryCount {
		return fmt.Errorf("%w: maximum number of material libraries reached", common.ErrLimitsExceeded)
	}
	c.model.MaterialLibraries = append(c.model.MaterialLibraries, event.FilePath)
	return nil
}

func (c *decodeContext) handleVertex(event objscan.VertexEvent) error {
	if len(c.model.Vertices) >= c.limits.MaxVertexCount {
		return fmt.Errorf("%w: maximum number of vertices reached", common.ErrLimitsExceeded)
	}
	c.model.Vertices = append(c.model.Vertices, Vertex{
		X: event.X,
		Y: event.Y,
		Z: event.Z,
		W: event.W,
	})
	return nil
}

func (c *decodeContext) handleTexCoord(event objscan.TexCoordEvent) error {
	if len(c.model.TexCoords) >= c.limits.MaxTexCoordCount {
		return fmt.Errorf("%w: maximum number of texture coordinates reached", common.ErrLimitsExceeded)
	}
	c.model.TexCoords = append(c.model.TexCoords, TexCoord{
		U: event.U,
		V: event.V,
		W: event.W,
	})
	return nil
}

func (c *decodeContext) handleNormal(event objscan.NormalEvent) error {
	if len(c.model.Normals) >= c.limits.MaxNormalCount {
		return fmt.Errorf("%w: maximum number of normals reached", common.ErrLimitsExceeded)
	}
	c.model.Normals = append(c.model.Normals, Normal{
		X: event.X,
		Y: event.Y,
		Z: event.Z,
	})
	return nil
}

func (c *decodeContext) handleObject(event objscan.ObjectEvent) error {
	if len(c.model.Objects) >= c.limits.MaxObjectCount {
		return fmt.Errorf("%w: maximum number of objects reached", common.ErrLimitsExceeded)
	}
	c.currentMesh = nil
	c.currentObject = new(Object)
	c.currentObject.Name = event.ObjectName
	c.model.Objects = append(c.model.Objects, c.currentObject)
	return nil
}

func (c *decodeContext) handleMaterialReference(event objscan.MaterialReferenceEvent) error {
	c.assureCurrentObject()
	mesh, found := c.currentObject.FindMesh(event.MaterialName)
	if found {
		c.currentMesh = mesh
	} else {
		if len(c.currentObject.Meshes) >= c.limits.MaxMaterialReferenceCount {
			return fmt.Errorf("%w: maximum number of material references reached", common.ErrLimitsExceeded)
		}
		c.currentMesh = new(Mesh)
		c.currentMesh.MaterialName = event.MaterialName
		c.currentObject.Meshes = append(c.currentObject.Meshes, c.currentMesh)
	}
	return nil
}

func (c *decodeContext) handleFaceStart() error {
	c.assureCurrentMesh()
	if len(c.currentMesh.Faces) >= c.limits.MaxFaceCount {
		return fmt.Errorf("%w: maximum number of faces reached", common.ErrLimitsExceeded)
	}
	c.currentFace = new(Face)
	return nil
}

func (c *decodeContext) handleFaceEnd() error {
	if len(c.currentFace.References) < 3 {
		return fmt.Errorf("%w: face needs to have at least three vertices", common.ErrInvalid)
	}
	c.currentMesh.Faces = append(c.currentMesh.Faces, c.currentFace)
	return nil
}

func (c *decodeContext) handleReferencesStart() error {
	if len(c.currentFace.References) >= c.limits.MaxReferenceCount {
		return fmt.Errorf("%w: maximum number of vertex references reached", common.ErrLimitsExceeded)
	}
	c.currentReference = &Reference{
		TexCoordIndex: UndefinedIndex,
		NormalIndex:   UndefinedIndex,
	}
	return nil
}

func (c *decodeContext) handleReferencesEnd() error {
	c.currentFace.References = append(c.currentFace.References, *c.currentReference)
	return nil
}

func (c *decodeContext) handleVertexReference(event objscan.VertexReferenceEvent) error {
	if event.VertexIndex > 0 {
		c.currentReference.VertexIndex = event.VertexIndex - 1
	} else {
		c.currentReference.VertexIndex = int64(len(c.model.Vertices)) + event.VertexIndex
	}
	return nil
}

func (c *decodeContext) handleTexCoordReference(event objscan.TexCoordReferenceEvent) error {
	if event.TexCoordIndex > 0 {
		c.currentReference.TexCoordIndex = event.TexCoordIndex - 1
	} else {
		c.currentReference.TexCoordIndex = int64(len(c.model.TexCoords)) + event.TexCoordIndex
	}
	return nil
}

func (c *decodeContext) handleNormalReference(event objscan.NormalReferenceEvent) error {
	if event.NormalIndex > 0 {
		c.currentReference.NormalIndex = event.NormalIndex - 1
	} else {
		c.currentReference.NormalIndex = int64(len(c.model.Normals)) + event.NormalIndex
	}
	return nil
}

func (c *decodeContext) assureCurrentObject() {
	if c.currentObject != nil {
		return
	}
	c.currentObject = &Object{
		Name: "Default",
	}
	c.model.Objects = append(c.model.Objects, c.currentObject)
}

func (c *decodeContext) assureCurrentMesh() {
	if c.currentMesh != nil {
		return
	}
	c.assureCurrentObject()
	c.currentMesh = new(Mesh)
	c.currentObject.Meshes = append(c.currentObject.Meshes, c.currentMesh)
}
