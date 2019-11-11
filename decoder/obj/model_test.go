package obj_test

import (
	. "github.com/DanTulovsky/go-data-front/decoder/obj"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Model", func() {

	Describe("Reference", func() {
		var reference Reference

		BeforeEach(func() {
			reference = Reference{
				VertexIndex:   10,
				TexCoordIndex: 15,
				NormalIndex:   20,
			}
		})

		It("has tex coord index", func() {
			Ω(reference.HasTexCoord()).Should(BeTrue())
		})

		It("has normal index", func() {
			Ω(reference.HasNormal()).Should(BeTrue())
		})

		Context("when tex coord index is set to undefined", func() {
			BeforeEach(func() {
				reference.TexCoordIndex = UndefinedIndex
			})

			It("does not have a tex coord index", func() {
				Ω(reference.HasTexCoord()).Should(BeFalse())
			})
		})

		Context("when normal index is set to undefined", func() {
			BeforeEach(func() {
				reference.NormalIndex = UndefinedIndex
			})

			It("does not have a normal index", func() {
				Ω(reference.HasNormal()).Should(BeFalse())
			})
		})
	})

	Describe("Object", func() {
		var object *Object

		BeforeEach(func() {
			object = new(Object)
		})

		Context("when object has multiple meshes", func() {
			var firstMesh *Mesh
			var secondMesh *Mesh

			BeforeEach(func() {
				firstMesh = &Mesh{
					MaterialName: "First",
				}
				secondMesh = &Mesh{
					MaterialName: "Second",
				}
				object.Meshes = append(object.Meshes, firstMesh, secondMesh)
			})

			It("is possible to find mesh by material name", func() {
				mesh, found := object.FindMesh("First")
				Ω(found).Should(BeTrue())
				Ω(mesh).Should(Equal(firstMesh))

				mesh, found = object.FindMesh("Second")
				Ω(found).Should(BeTrue())
				Ω(mesh).Should(Equal(secondMesh))
			})

			It("will not find unexisting meshes", func() {
				_, found := object.FindMesh("Missing")
				Ω(found).Should(BeFalse())
			})
		})
	})

	Describe("Model", func() {
		var model *Model

		BeforeEach(func() {
			model = new(Model)
		})

		Context("when model has multiple vertices", func() {
			var firstVertex Vertex
			var secondVertex Vertex

			BeforeEach(func() {
				firstVertex = Vertex{
					X: 1.0, Y: 2.0, Z: 3.0,
				}
				secondVertex = Vertex{
					X: 4.0, Y: 5.0, Z: 6.0,
				}
				model.Vertices = append(model.Vertices, firstVertex, secondVertex)
			})

			It("is possible to get vertex from reference", func() {
				vertex := model.GetVertexFromReference(Reference{
					VertexIndex: 0,
				})
				Ω(vertex).Should(Equal(firstVertex))

				vertex = model.GetVertexFromReference(Reference{
					VertexIndex: 1,
				})
				Ω(vertex).Should(Equal(secondVertex))
			})
		})

		Context("when model has multiple texture coordinates", func() {
			var firstTexCoord TexCoord
			var secondTexCoord TexCoord

			BeforeEach(func() {
				firstTexCoord = TexCoord{
					U: 1.0, V: 2.0, W: 3.0,
				}
				secondTexCoord = TexCoord{
					U: 4.0, V: 5.0, W: 6.0,
				}
				model.TexCoords = append(model.TexCoords, firstTexCoord, secondTexCoord)
			})

			It("is possible to get vertex from reference", func() {
				texCoord := model.GetTexCoordFromReference(Reference{
					TexCoordIndex: 0,
				})
				Ω(texCoord).Should(Equal(firstTexCoord))

				texCoord = model.GetTexCoordFromReference(Reference{
					TexCoordIndex: 1,
				})
				Ω(texCoord).Should(Equal(secondTexCoord))
			})
		})

		Context("when model has multiple normals", func() {
			var firstNormal Normal
			var secondNormal Normal

			BeforeEach(func() {
				firstNormal = Normal{
					X: 1.0, Y: 2.0, Z: 3.0,
				}
				secondNormal = Normal{
					X: 4.0, Y: 5.0, Z: 6.0,
				}
				model.Normals = append(model.Normals, firstNormal, secondNormal)
			})

			It("is possible to get vertex from reference", func() {
				normal := model.GetNormalFromReference(Reference{
					NormalIndex: 0,
				})
				Ω(normal).Should(Equal(firstNormal))

				normal = model.GetNormalFromReference(Reference{
					NormalIndex: 1,
				})
				Ω(normal).Should(Equal(secondNormal))
			})
		})

		Context("when model has multiple objects", func() {
			var firstObject *Object
			var secondObject *Object

			BeforeEach(func() {
				firstObject = &Object{
					Name: "First",
				}
				secondObject = &Object{
					Name: "Second",
				}
				model.Objects = append(model.Objects, firstObject, secondObject)
			})

			It("is possible to find existing objects", func() {
				object, found := model.FindObject("First")
				Ω(found).Should(BeTrue())
				Ω(object).Should(Equal(firstObject))

				object, found = model.FindObject("Second")
				Ω(found).Should(BeTrue())
				Ω(object).Should(Equal(secondObject))
			})

			It("will not find unexisting objects", func() {
				_, found := model.FindObject("Missing")
				Ω(found).Should(BeFalse())
			})
		})
	})
})
