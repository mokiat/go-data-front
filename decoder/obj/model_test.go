package obj_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mokiat/go-data-front/decoder/obj"
)

var _ = Describe("Reference", func() {
	var reference obj.Reference

	BeforeEach(func() {
		reference = obj.Reference{
			VertexIndex:   10,
			TexCoordIndex: 15,
			NormalIndex:   20,
		}
	})

	It("has tex coord index", func() {
		Expect(reference.HasTexCoord()).To(BeTrue())
	})

	It("has normal index", func() {
		Expect(reference.HasNormal()).To(BeTrue())
	})

	When("the tex coord index is set to undefined", func() {
		BeforeEach(func() {
			reference.TexCoordIndex = obj.UndefinedIndex
		})

		It("does not have a tex coord index", func() {
			Expect(reference.HasTexCoord()).To(BeFalse())
		})
	})

	When("the normal index is set to undefined", func() {
		BeforeEach(func() {
			reference.NormalIndex = obj.UndefinedIndex
		})

		It("does not have a normal index", func() {
			Expect(reference.HasNormal()).To(BeFalse())
		})
	})
})

var _ = Describe("Object", func() {
	var object *obj.Object

	BeforeEach(func() {
		object = new(obj.Object)
	})

	When("the object has multiple meshes", func() {
		var (
			firstMesh  *obj.Mesh
			secondMesh *obj.Mesh
		)

		BeforeEach(func() {
			firstMesh = &obj.Mesh{
				MaterialName: "First",
			}
			secondMesh = &obj.Mesh{
				MaterialName: "Second",
			}
			object.Meshes = append(object.Meshes, firstMesh, secondMesh)
		})

		It("is possible to find mesh by material name", func() {
			mesh, found := object.FindMesh("First")
			Expect(found).To(BeTrue())
			Expect(mesh).To(Equal(firstMesh))

			mesh, found = object.FindMesh("Second")
			Expect(found).To(BeTrue())
			Expect(mesh).To(Equal(secondMesh))
		})

		It("will not find unexisting meshes", func() {
			_, found := object.FindMesh("Missing")
			Expect(found).To(BeFalse())
		})
	})
})

var _ = Describe("Model", func() {
	var model *obj.Model

	BeforeEach(func() {
		model = new(obj.Model)
	})

	When("the model has multiple vertices", func() {
		var (
			firstVertex  obj.Vertex
			secondVertex obj.Vertex
		)

		BeforeEach(func() {
			firstVertex = obj.Vertex{
				X: 1.0, Y: 2.0, Z: 3.0,
			}
			secondVertex = obj.Vertex{
				X: 4.0, Y: 5.0, Z: 6.0,
			}
			model.Vertices = append(model.Vertices, firstVertex, secondVertex)
		})

		It("is possible to get vertex from reference", func() {
			vertex := model.GetVertexFromReference(obj.Reference{
				VertexIndex: 0,
			})
			Expect(vertex).To(Equal(firstVertex))

			vertex = model.GetVertexFromReference(obj.Reference{
				VertexIndex: 1,
			})
			Expect(vertex).To(Equal(secondVertex))
		})
	})

	When("the model has multiple texture coordinates", func() {
		var (
			firstTexCoord  obj.TexCoord
			secondTexCoord obj.TexCoord
		)

		BeforeEach(func() {
			firstTexCoord = obj.TexCoord{
				U: 1.0, V: 2.0, W: 3.0,
			}
			secondTexCoord = obj.TexCoord{
				U: 4.0, V: 5.0, W: 6.0,
			}
			model.TexCoords = append(model.TexCoords, firstTexCoord, secondTexCoord)
		})

		It("is possible to get vertex from reference", func() {
			texCoord := model.GetTexCoordFromReference(obj.Reference{
				TexCoordIndex: 0,
			})
			Expect(texCoord).To(Equal(firstTexCoord))

			texCoord = model.GetTexCoordFromReference(obj.Reference{
				TexCoordIndex: 1,
			})
			Expect(texCoord).To(Equal(secondTexCoord))
		})
	})

	When("the model has multiple normals", func() {
		var (
			firstNormal  obj.Normal
			secondNormal obj.Normal
		)

		BeforeEach(func() {
			firstNormal = obj.Normal{
				X: 1.0, Y: 2.0, Z: 3.0,
			}
			secondNormal = obj.Normal{
				X: 4.0, Y: 5.0, Z: 6.0,
			}
			model.Normals = append(model.Normals, firstNormal, secondNormal)
		})

		It("is possible to get vertex from reference", func() {
			normal := model.GetNormalFromReference(obj.Reference{
				NormalIndex: 0,
			})
			Expect(normal).To(Equal(firstNormal))

			normal = model.GetNormalFromReference(obj.Reference{
				NormalIndex: 1,
			})
			Expect(normal).To(Equal(secondNormal))
		})
	})

	When("the model has multiple objects", func() {
		var (
			firstObject  *obj.Object
			secondObject *obj.Object
		)

		BeforeEach(func() {
			firstObject = &obj.Object{
				Name: "First",
			}
			secondObject = &obj.Object{
				Name: "Second",
			}
			model.Objects = append(model.Objects, firstObject, secondObject)
		})

		It("is possible to find existing objects", func() {
			object, found := model.FindObject("First")
			Expect(found).To(BeTrue())
			Expect(object).To(Equal(firstObject))

			object, found = model.FindObject("Second")
			Expect(found).To(BeTrue())
			Expect(object).To(Equal(secondObject))
		})

		It("will not find unexisting objects", func() {
			_, found := model.FindObject("Missing")
			Expect(found).To(BeFalse())
		})
	})
})
