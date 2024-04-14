package obj_test

import (
	"fmt"
	"os"

	. "github.com/mokiat/go-data-front/decoder/obj"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decoder", func() {
	var limits DecodeLimits
	var decoder Decoder
	var model *Model
	var decodeErr error

	decodeFile := func(filename string) {
		file, err := os.Open(fmt.Sprintf("testdata/%s", filename))
		if err != nil {
			panic(err)
		}
		defer file.Close()
		decoder = NewDecoder(limits)
		model, decodeErr = decoder.Decode(file)
	}

	itShouldHaveReturnedAnError := func() {
		It("should have returned an error", func() {
			Ω(decodeErr).Should(HaveOccurred())
		})
	}

	itShouldNotHaveReturnedAnError := func() {
		It("should not have returned an error", func() {
			Ω(decodeErr).ShouldNot(HaveOccurred())
		})
	}

	BeforeEach(func() {
		limits = DefaultLimits()
	})

	Context("when a basic file is decoded", func() {
		BeforeEach(func() {
			decodeFile("valid_basic.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have returned a non-nil model", func() {
			Ω(model).ShouldNot(BeNil())
		})

		It("should have decoded material libraries", func() {
			Ω(model.MaterialLibraries).Should(HaveLen(1))
			Ω(model.MaterialLibraries[0]).Should(Equal("materials.mtl"))
		})

		It("should have decoded vertices", func() {
			Ω(model.Vertices).Should(HaveLen(4))
			Ω(model.Vertices[0]).Should(Equal(Vertex{
				X: 0.1, Y: 0.2, Z: 0.3, W: 1.0,
			}))
			Ω(model.Vertices[1]).Should(Equal(Vertex{
				X: 0.4, Y: 0.5, Z: 0.6, W: 1.0,
			}))
			Ω(model.Vertices[2]).Should(Equal(Vertex{
				X: 0.7, Y: 0.8, Z: 0.9, W: 1.0,
			}))
			Ω(model.Vertices[3]).Should(Equal(Vertex{
				X: 1.0, Y: 0.9, Z: 0.8, W: 1.0,
			}))
		})

		It("should have decoded normals", func() {
			Ω(model.Normals).Should(HaveLen(3))
			Ω(model.Normals[0]).Should(Equal(Normal{
				X: 0.0, Y: 1.0, Z: 0.0,
			}))
			Ω(model.Normals[1]).Should(Equal(Normal{
				X: 1.0, Y: 0.0, Z: 0.0,
			}))
			Ω(model.Normals[2]).Should(Equal(Normal{
				X: 0.0, Y: 0.0, Z: 1.0,
			}))
		})

		It("should have decoded texture coordinates", func() {
			Ω(model.TexCoords).Should(HaveLen(4))
			Ω(model.TexCoords[0]).Should(Equal(TexCoord{
				U: 0.1, V: 0.2, W: 0.3,
			}))
			Ω(model.TexCoords[1]).Should(Equal(TexCoord{
				U: 0.4, V: 0.5, W: 0.6,
			}))
			Ω(model.TexCoords[2]).Should(Equal(TexCoord{
				U: 0.7, V: 0.8, W: 0.9,
			}))
			Ω(model.TexCoords[3]).Should(Equal(TexCoord{
				U: 1.0, V: 0.9, W: 0.8,
			}))
		})

		It("should have decoded objects", func() {
			Ω(model.Objects).Should(HaveLen(1))
			Ω(model.Objects[0].Name).Should(Equal("MyObject"))
		})

		It("should have decoded meshes", func() {
			object := model.Objects[0]
			Ω(object.Meshes).Should(HaveLen(1))
			Ω(object.Meshes[0].MaterialName).Should(Equal("BlueMaterial"))
		})

		It("should have decoded faces", func() {
			mesh := model.Objects[0].Meshes[0]
			Ω(mesh.Faces).Should(HaveLen(2))
		})

		It("should have decoded references", func() {
			face := model.Objects[0].Meshes[0].Faces[0]
			Ω(face.References).Should(HaveLen(3))

			Ω(face.References[0]).Should(Equal(Reference{
				VertexIndex:   0,
				TexCoordIndex: 3,
				NormalIndex:   0,
			}))
			Ω(face.References[0].HasTexCoord()).Should(BeTrue())
			Ω(face.References[0].HasNormal()).Should(BeTrue())

			Ω(face.References[1]).Should(Equal(Reference{
				VertexIndex:   1,
				TexCoordIndex: 0,
				NormalIndex:   0,
			}))
			Ω(face.References[1].HasTexCoord()).Should(BeTrue())
			Ω(face.References[1].HasNormal()).Should(BeTrue())

			Ω(face.References[2]).Should(Equal(Reference{
				VertexIndex:   2,
				TexCoordIndex: 2,
				NormalIndex:   1,
			}))
			Ω(face.References[2].HasTexCoord()).Should(BeTrue())
			Ω(face.References[2].HasNormal()).Should(BeTrue())
		})
	})

	Context("when a file with multiple objects is decoded", func() {
		JustBeforeEach(func() {
			decodeFile("valid_objects.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all objects", func() {
			Ω(model.Objects).Should(HaveLen(2))
			Ω(model.Objects[0].Name).Should(Equal("First"))
			Ω(model.Objects[1].Name).Should(Equal("Second"))
		})

		Context("when the number of objects is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxObjectCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	Context("when a file with all kinds of vertices is decoded", func() {
		JustBeforeEach(func() {
			decodeFile("valid_vertices.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all vertices", func() {
			Ω(model.Vertices).Should(HaveLen(2))
			Ω(model.Vertices[0]).Should(Equal(Vertex{
				X: 1.0, Y: 2.0, Z: 3.0, W: 1.0,
			}))
			Ω(model.Vertices[1]).Should(Equal(Vertex{
				X: 4.0, Y: 5.0, Z: 6.0, W: 7.0,
			}))
		})

		Context("when the number of vertices is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxVertexCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	Context("when a file with all kinds of texture coordinates is decoded", func() {
		JustBeforeEach(func() {
			decodeFile("valid_tex_coords.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all texture coordinates", func() {
			Ω(model.TexCoords).Should(HaveLen(3))
			Ω(model.TexCoords[0]).Should(Equal(TexCoord{
				U: 0.1, V: 0.0, W: 0.0,
			}))
			Ω(model.TexCoords[1]).Should(Equal(TexCoord{
				U: 0.4, V: 0.5, W: 0.0,
			}))
			Ω(model.TexCoords[2]).Should(Equal(TexCoord{
				U: 0.7, V: 0.8, W: 0.9,
			}))
		})

		Context("when the number of texture coordinates is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxTexCoordCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	Context("when a file with all kinds of normals is decoded", func() {
		JustBeforeEach(func() {
			decodeFile("valid_normals.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all normals", func() {
			Ω(model.Normals).Should(HaveLen(2))
			Ω(model.Normals[0]).Should(Equal(Normal{
				X: 1.0, Y: -2.0, Z: 3.0,
			}))
			Ω(model.Normals[1]).Should(Equal(Normal{
				X: -0.5, Y: 0.4, Z: -0.7,
			}))
		})

		Context("when the number of normals is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxNormalCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	Context("when a file with all kinds of material libraries is scanned", func() {
		JustBeforeEach(func() {
			decodeFile("valid_material_libraries.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all material libraries", func() {
			Ω(model.MaterialLibraries).Should(HaveLen(3))
			Ω(model.MaterialLibraries[0]).Should(Equal("first.mtl"))
			Ω(model.MaterialLibraries[1]).Should(Equal("second.mtl"))
			Ω(model.MaterialLibraries[2]).Should(Equal("third.mtl"))
		})

		Context("when the number of material libraries is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxMaterialLibraryCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	Context("when a file with multiple material references is scanned", func() {
		JustBeforeEach(func() {
			decodeFile("valid_material_references.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all material references", func() {
			object := model.Objects[0]
			Ω(object.Meshes).Should(HaveLen(2))
			Ω(object.Meshes[0].MaterialName).Should(Equal("Red"))
			Ω(object.Meshes[1].MaterialName).Should(Equal("Blue"))
		})

		Context("when the number of material references is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxMaterialReferenceCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	Context("when a file with multiple faces is scanned", func() {
		JustBeforeEach(func() {
			decodeFile("valid_faces.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all faces", func() {
			mesh := model.Objects[0].Meshes[0]
			Ω(mesh.Faces).Should(HaveLen(2))
		})

		Context("when the number of faces is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxFaceCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	Context("when a file with multiple references is scanned", func() {
		JustBeforeEach(func() {
			decodeFile("valid_references.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all faces", func() {
			face := model.Objects[0].Meshes[0].Faces[0]
			Ω(face.References).Should(HaveLen(4))
		})

		Context("when the number of faces is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxReferenceCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	Context("when an object with duplicate materials is scanned", func() {
		JustBeforeEach(func() {
			limits.MaxMaterialReferenceCount = 2
			decodeFile("valid_mesh_reuse.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have parsed only two meshes", func() {
			object := model.Objects[0]
			Ω(object.Meshes).Should(HaveLen(2))
			Ω(object.Meshes[0].MaterialName).Should(Equal("First"))
			Ω(object.Meshes[0].Faces).Should(HaveLen(2))
			Ω(object.Meshes[1].MaterialName).Should(Equal("Second"))
			Ω(object.Meshes[1].Faces).Should(HaveLen(1))
		})
	})

	Context("when an object without a mesh is scanned", func() {
		JustBeforeEach(func() {
			decodeFile("valid_no_mesh.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have created a default mesh", func() {
			meshes := model.Objects[0].Meshes
			Ω(meshes).Should(HaveLen(1))
			Ω(meshes[0]).ShouldNot(BeNil())
			Ω(meshes[0].MaterialName).Should(Equal(""))
		})
	})

	Context("when faces without an object are scanned", func() {
		JustBeforeEach(func() {
			decodeFile("valid_no_object.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have created a default object", func() {
			Ω(model.Objects).Should(HaveLen(1))
			Ω(model.Objects[0].Name).Should(Equal("Default"))

			Ω(model.Objects[0].Meshes).Should(HaveLen(1))
			Ω(model.Objects[0].Meshes[0].MaterialName).Should(Equal("Red"))
		})
	})

	Context("when faces without an object and mesh are scanned", func() {
		JustBeforeEach(func() {
			decodeFile("valid_no_object_no_mesh.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have created a default object and default mesh", func() {
			Ω(model.Objects).Should(HaveLen(1))
			Ω(model.Objects[0].Name).Should(Equal("Default"))

			Ω(model.Objects[0].Meshes).Should(HaveLen(1))
			Ω(model.Objects[0].Meshes[0].MaterialName).Should(Equal(""))
		})
	})

	Context("when negative indices are scanned", func() {
		JustBeforeEach(func() {
			decodeFile("valid_negative_indices.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("have normalized the references", func() {
			face := model.Objects[0].Meshes[0].Faces[0]
			Ω(face.References).Should(HaveLen(3))
			Ω(face.References[0]).Should(Equal(Reference{
				VertexIndex:   0,
				TexCoordIndex: 3,
				NormalIndex:   0,
			}))
			Ω(face.References[1]).Should(Equal(Reference{
				VertexIndex:   1,
				TexCoordIndex: 0,
				NormalIndex:   0,
			}))
			Ω(face.References[2]).Should(Equal(Reference{
				VertexIndex:   2,
				TexCoordIndex: 2,
				NormalIndex:   1,
			}))
		})
	})

	Context("when decoding face without enough references", func() {
		BeforeEach(func() {
			decodeFile("error_missing_face_data.obj")
		})

		itShouldHaveReturnedAnError()
	})

	Describe("Default DecodeLimits", func() {
		var limits DecodeLimits

		BeforeEach(func() {
			limits = DefaultLimits()
		})

		It("vertex limit should be 65536", func() {
			Ω(limits.MaxVertexCount).Should(Equal(65536))
		})

		It("normal limit should be 65536", func() {
			Ω(limits.MaxNormalCount).Should(Equal(65536))
		})

		It("tex coord limit should be 65536", func() {
			Ω(limits.MaxTexCoordCount).Should(Equal(65536))
		})

		It("object limit should be 1024", func() {
			Ω(limits.MaxObjectCount).Should(Equal(1024))
		})

		It("face limit should be 65536", func() {
			Ω(limits.MaxFaceCount).Should(Equal(65536))
		})

		It("reference limit should be 16", func() {
			Ω(limits.MaxReferenceCount).Should(Equal(16))
		})

		It("material reference limit should be 64", func() {
			Ω(limits.MaxMaterialReferenceCount).Should(Equal(64))
		})

		It("material library limit should be 32", func() {
			Ω(limits.MaxMaterialLibraryCount).Should(Equal(32))
		})
	})
})
