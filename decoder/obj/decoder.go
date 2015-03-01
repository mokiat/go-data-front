package obj

import (
	"errors"
	"io"

	"github.com/momchil-atanasov/go-data-front/common"
	scanOBJ "github.com/momchil-atanasov/go-data-front/scanner/obj"
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

func NewDecoder(limits DecodeLimits) Decoder {
	return &decoder{
		limits:  &limits,
		scanner: scanOBJ.NewScanner(),
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
	return context.Model, nil
}

func newDecodeContext(limits *DecodeLimits) *decodeContext {
	return &decodeContext{
		Limits:        limits,
		Model:         new(Model),
		CurrentObject: nil,
		CurrentMesh:   nil,
		CurrentFace:   nil,
	}
}

type decodeContext struct {
	Limits           *DecodeLimits
	Model            *Model
	CurrentObject    *Object
	CurrentMesh      *Mesh
	CurrentFace      *Face
	CurrentReference *Reference
}

func (c *decodeContext) HandleEvent(event common.Event) error {
	switch actual := event.(type) {
	case scanOBJ.MaterialLibraryEvent:
		return c.handleMaterialLibrary(actual)
	case scanOBJ.VertexEvent:
		return c.handleVertex(actual)
	case scanOBJ.NormalEvent:
		return c.handleNormal(actual)
	case scanOBJ.TexCoordEvent:
		return c.handleTexCoord(actual)
	case scanOBJ.ObjectEvent:
		return c.handleObject(actual)
	case scanOBJ.MaterialReferenceEvent:
		return c.handleMaterialReference(actual)
	case scanOBJ.FaceStartEvent:
		return c.handleFaceStart()
	case scanOBJ.FaceEndEvent:
		return c.handleFaceEnd()
	case scanOBJ.ReferenceSetStartEvent:
		return c.handleReferencesStart()
	case scanOBJ.ReferenceSetEndEvent:
		return c.handleReferencesEnd()
	case scanOBJ.VertexReferenceEvent:
		return c.handleVertexReference(actual)
	case scanOBJ.TexCoordReferenceEvent:
		return c.handleTexCoordReference(actual)
	case scanOBJ.NormalReferenceEvent:
		return c.handleNormalReference(actual)
	}
	return nil
}

func (c *decodeContext) handleMaterialLibrary(event scanOBJ.MaterialLibraryEvent) error {
	if len(c.Model.MaterialLibraries) >= c.Limits.MaxMaterialLibraryCount {
		return errors.New("Maximum number of material libraries reached!")
	}
	c.Model.MaterialLibraries = append(c.Model.MaterialLibraries, event.FilePath)
	return nil
}

func (c *decodeContext) handleVertex(event scanOBJ.VertexEvent) error {
	if len(c.Model.Vertices) >= c.Limits.MaxVertexCount {
		return errors.New("Maximum number of vertices reached!")
	}
	c.Model.Vertices = append(c.Model.Vertices, Vertex{
		X: event.X,
		Y: event.Y,
		Z: event.Z,
		W: event.W,
	})
	return nil
}

func (c *decodeContext) handleTexCoord(event scanOBJ.TexCoordEvent) error {
	if len(c.Model.TexCoords) >= c.Limits.MaxTexCoordCount {
		return errors.New("Maximum number of texture coordinates reached!")
	}
	c.Model.TexCoords = append(c.Model.TexCoords, TexCoord{
		U: event.U,
		V: event.V,
		W: event.W,
	})
	return nil
}

func (c *decodeContext) handleNormal(event scanOBJ.NormalEvent) error {
	if len(c.Model.Normals) >= c.Limits.MaxNormalCount {
		return errors.New("Maximum number of normals reached!")
	}
	c.Model.Normals = append(c.Model.Normals, Normal{
		X: event.X,
		Y: event.Y,
		Z: event.Z,
	})
	return nil
}

func (c *decodeContext) handleObject(event scanOBJ.ObjectEvent) error {
	if len(c.Model.Objects) >= c.Limits.MaxObjectCount {
		return errors.New("Maximum number of objects reached!")
	}
	c.CurrentMesh = nil
	c.CurrentObject = new(Object)
	c.CurrentObject.Name = event.ObjectName
	c.Model.Objects = append(c.Model.Objects, c.CurrentObject)
	return nil
}

func (c *decodeContext) handleMaterialReference(event scanOBJ.MaterialReferenceEvent) error {
	c.assureCurrentObject()
	mesh, found := c.CurrentObject.FindMesh(event.MaterialName)
	if found {
		c.CurrentMesh = mesh
	} else {
		if len(c.CurrentObject.Meshes) >= c.Limits.MaxMaterialReferenceCount {
			return errors.New("Maximum number of material references reached!")
		}
		c.CurrentMesh = new(Mesh)
		c.CurrentMesh.MaterialName = event.MaterialName
		c.CurrentObject.Meshes = append(c.CurrentObject.Meshes, c.CurrentMesh)
	}
	return nil
}

func (c *decodeContext) handleFaceStart() error {
	c.assureCurrentMesh()
	if len(c.CurrentMesh.Faces) >= c.Limits.MaxFaceCount {
		return errors.New("Maximum number of faces reached!")
	}
	c.CurrentFace = new(Face)
	return nil
}

func (c *decodeContext) handleFaceEnd() error {
	if len(c.CurrentFace.References) < 3 {
		return errors.New("Face needs to have at least three vertices.")
	}
	c.CurrentMesh.Faces = append(c.CurrentMesh.Faces, c.CurrentFace)
	return nil
}

func (c *decodeContext) handleReferencesStart() error {
	if len(c.CurrentFace.References) >= c.Limits.MaxReferenceCount {
		return errors.New("Maximum number of vertex references reached!")
	}
	c.CurrentReference = &Reference{
		TexCoordIndex: UndefinedIndex,
		NormalIndex:   UndefinedIndex,
	}
	return nil
}

func (c *decodeContext) handleReferencesEnd() error {
	c.CurrentFace.References = append(c.CurrentFace.References, *c.CurrentReference)
	return nil
}

func (c *decodeContext) handleVertexReference(event scanOBJ.VertexReferenceEvent) error {
	if event.VertexIndex > 0 {
		c.CurrentReference.VertexIndex = event.VertexIndex - 1
	} else {
		c.CurrentReference.VertexIndex = int64(len(c.Model.Vertices)) + event.VertexIndex
	}
	return nil
}

func (c *decodeContext) handleTexCoordReference(event scanOBJ.TexCoordReferenceEvent) error {
	if event.TexCoordIndex > 0 {
		c.CurrentReference.TexCoordIndex = event.TexCoordIndex - 1
	} else {
		c.CurrentReference.TexCoordIndex = int64(len(c.Model.TexCoords)) + event.TexCoordIndex
	}
	return nil
}

func (c *decodeContext) handleNormalReference(event scanOBJ.NormalReferenceEvent) error {
	if event.NormalIndex > 0 {
		c.CurrentReference.NormalIndex = event.NormalIndex - 1
	} else {
		c.CurrentReference.NormalIndex = int64(len(c.Model.Normals)) + event.NormalIndex
	}
	return nil
}

func (c *decodeContext) assureCurrentObject() {
	if c.CurrentObject != nil {
		return
	}
	c.CurrentObject = &Object{
		Name: "Default",
	}
	c.Model.Objects = append(c.Model.Objects, c.CurrentObject)
}

func (c *decodeContext) assureCurrentMesh() {
	if c.CurrentMesh != nil {
		return
	}
	c.assureCurrentObject()
	c.CurrentMesh = new(Mesh)
	c.CurrentObject.Meshes = append(c.CurrentObject.Meshes, c.CurrentMesh)
}
