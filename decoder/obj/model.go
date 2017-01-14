package obj

// Model represents the data of a single Wavefront OBJ resource.
type Model struct {

	// Vertex holds a list of all the vertices
	Vertices []Vertex

	// Normals holds a list of all the normals
	Normals []Normal

	// TexCoords holds a list of the texture coordinates
	TexCoords []TexCoord

	// Objects holds a list of all the objects
	Objects []*Object

	// MaterialLibraries holds a list of filenames to MTL
	// resources that should be used together with the current
	// OBJ resource
	MaterialLibraries []string
}

// GetVertexFromReference is a helper method that allows one
// to get the Vertex directly from a Reference
func (m *Model) GetVertexFromReference(ref Reference) Vertex {
	return m.Vertices[ref.VertexIndex]
}

// GetTexCoordFromReference is a helper method that allows one
// to get the TexCoord directly from a Reference
func (m *Model) GetTexCoordFromReference(ref Reference) TexCoord {
	return m.TexCoords[ref.TexCoordIndex]
}

// GetNormalFromReference is a helper method that allows one
// to get the Normal directly from a Reference
func (m *Model) GetNormalFromReference(ref Reference) Normal {
	return m.Normals[ref.NormalIndex]
}

// FindObject is a helper method that allows one to search
// for an object in this model based on name
func (m *Model) FindObject(name string) (*Object, bool) {
	for _, object := range m.Objects {
		if object.Name == name {
			return object, true
		}
	}
	return nil, false
}

// Vertex is used to define the positional
// information for objects.
type Vertex struct {

	// X coordinate of this vertex.
	X float64

	// Y coordinate of this vertex.
	Y float64

	// Z coordinate of this vertex.
	Z float64

	// W coordinate of this vertex. (By default 1.0)
	W float64
}

// Normal is used to define the directional
// information for objects.
type Normal struct {

	// X coordinate of this normal.
	X float64

	// Y coordinate of this normal.
	Y float64

	// Z coordinate of this normal.
	Z float64
}

// TexCoord is used to define the texture
// mapping on an object's surface
type TexCoord struct {

	// U coordinate of this texture coordinate.
	U float64

	// V coordinate of this texture coordinate.
	V float64

	// W coordinate of this texture coordinate.
	W float64
}

// Object represents an object in the model
type Object struct {

	// Name holds the name of the object
	Name string

	// Meshes holds a list of meshes that make up
	// the object's shape
	Meshes []*Mesh
}

// FindMesh is a helper function that allows one to
// find a Mesh within an Object by searching by
// its material name
func (o *Object) FindMesh(materialName string) (*Mesh, bool) {
	for _, mesh := range o.Meshes {
		if mesh.MaterialName == materialName {
			return mesh, true
		}
	}
	return nil, false
}

// Mesh is a concept that cannot be directly mapped
// to OBJ specification. It is here to allow for the
// separation of mesh data based on material.
//
// Since a single object can have sections that use
// various materials, the Mesh is a mechanism through
// which such sections are separated.
type Mesh struct {

	// MaterialName holds the name of the material that
	// should be used for the rendering of this mesh
	MaterialName string

	// Faces holds all the faces that comprise this mesh
	Faces []*Face
}

// Face defines a single face that is part of a mesh
type Face struct {

	// References holds an array of Reference objects
	//
	// Each Reference holds information for a single
	// point in space. The list of all references compose
	// a polygon shape.
	References []Reference
}

// UndefinedIndex is used to mark an index as undefined.
const UndefinedIndex int64 = -1

// Reference is used to describe a single point from
// the set of points that define a shape
//
// Since a single point can have position, direction and
// texture information, all of that data is held in this
// structure.
type Reference struct {

	// VertexIndex holds the index into the array of vertices
	// for the positional data of this point.
	VertexIndex int64

	// TexCoordIndex holds the index into the array of texture
	// coordinates for the texture data of this point.
	//
	// If this value is equal to UndefinedIndex, then this point
	// does not have texture information.
	TexCoordIndex int64

	// NormalIndex holds the index into the array of normals
	// for the directional data of this point.
	//
	// If this value is equal to UndefinedIndex, then this point
	// does not have directional information.
	NormalIndex int64
}

// HasTexCoord is a helper method that helps determine
// whether the current Reference has texture coordinate
// information.
func (r Reference) HasTexCoord() bool {
	return r.TexCoordIndex != UndefinedIndex
}

// HasNormal is a helper method that helps determine
// whether the current Reference has normal information.
func (r Reference) HasNormal() bool {
	return r.NormalIndex != UndefinedIndex
}
