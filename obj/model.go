package obj

type Normal struct {
	X float32
	Y float32
	Z float32
}

type Vertex struct {
	X float32
	Y float32
	Z float32
}

type TexCoord struct {
	U float32
	V float32
	W float32
}

const UndefinedIndex int = -1

type Reference struct {
	VertexIndex   int
	NormalIndex   int
	TexCoordIndex int
	TexCoordSize  int
}

type Face struct {
	References []Reference
}

type Mesh struct {
	Faces        []Face
	MaterialName string
}

type Object struct {
	Meshes []Mesh
	Name   string
}

type Model struct {
	Vertices          []Vertex
	Normals           []Normal
	TexCoords         []TexCoord
	Objects           []Object
	MaterialLibraries []string
}
