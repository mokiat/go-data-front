package obj_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mokiat/go-data-front/decoder/obj"
)

var _ = Describe("Decoder", func() {
	var (
		testFile string
		limits   obj.DecodeLimits

		model     *obj.Model
		decodeErr error
	)

	itShouldHaveReturnedAnError := func() {
		GinkgoHelper()
		It("should have returned an error", func() {
			Expect(decodeErr).To(HaveOccurred())
		})
	}

	itShouldNotHaveReturnedAnError := func() {
		GinkgoHelper()
		It("should not have returned an error", func() {
			Expect(decodeErr).ToNot(HaveOccurred())
		})
	}

	JustBeforeEach(func() {
		file, err := os.Open(filepath.Join("testdata", testFile))
		Expect(err).ToNot(HaveOccurred())
		defer file.Close()

		decoder := obj.NewDecoder(limits)
		model, decodeErr = decoder.Decode(file)
	})

	BeforeEach(func() {
		limits = obj.DefaultLimits()
	})

	When("a basic file is decoded", func() {
		BeforeEach(func() {
			testFile = "valid_basic.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have returned a non-nil model", func() {
			Expect(model).ToNot(BeNil())
		})

		It("should have decoded material libraries", func() {
			Expect(model.MaterialLibraries).To(HaveLen(1))
			Expect(model.MaterialLibraries[0]).To(Equal("materials.mtl"))
		})

		It("should have decoded vertices", func() {
			Expect(model.Vertices).To(HaveLen(4))
			Expect(model.Vertices[0]).To(Equal(obj.Vertex{
				X: 0.1, Y: 0.2, Z: 0.3, W: 1.0,
			}))
			Expect(model.Vertices[1]).To(Equal(obj.Vertex{
				X: 0.4, Y: 0.5, Z: 0.6, W: 1.0,
			}))
			Expect(model.Vertices[2]).To(Equal(obj.Vertex{
				X: 0.7, Y: 0.8, Z: 0.9, W: 1.0,
			}))
			Expect(model.Vertices[3]).To(Equal(obj.Vertex{
				X: 1.0, Y: 0.9, Z: 0.8, W: 1.0,
			}))
		})

		It("should have decoded normals", func() {
			Expect(model.Normals).To(HaveLen(3))
			Expect(model.Normals[0]).To(Equal(obj.Normal{
				X: 0.0, Y: 1.0, Z: 0.0,
			}))
			Expect(model.Normals[1]).To(Equal(obj.Normal{
				X: 1.0, Y: 0.0, Z: 0.0,
			}))
			Expect(model.Normals[2]).To(Equal(obj.Normal{
				X: 0.0, Y: 0.0, Z: 1.0,
			}))
		})

		It("should have decoded texture coordinates", func() {
			Expect(model.TexCoords).To(HaveLen(4))
			Expect(model.TexCoords[0]).To(Equal(obj.TexCoord{
				U: 0.1, V: 0.2, W: 0.3,
			}))
			Expect(model.TexCoords[1]).To(Equal(obj.TexCoord{
				U: 0.4, V: 0.5, W: 0.6,
			}))
			Expect(model.TexCoords[2]).To(Equal(obj.TexCoord{
				U: 0.7, V: 0.8, W: 0.9,
			}))
			Expect(model.TexCoords[3]).To(Equal(obj.TexCoord{
				U: 1.0, V: 0.9, W: 0.8,
			}))
		})

		It("should have decoded objects", func() {
			Expect(model.Objects).To(HaveLen(1))
			Expect(model.Objects[0].Name).To(Equal("MyObject"))
		})

		It("should have decoded meshes", func() {
			object := model.Objects[0]
			Expect(object.Meshes).To(HaveLen(1))
			Expect(object.Meshes[0].MaterialName).To(Equal("BlueMaterial"))
		})

		It("should have decoded faces", func() {
			mesh := model.Objects[0].Meshes[0]
			Expect(mesh.Faces).To(HaveLen(2))
		})

		It("should have decoded references", func() {
			face := model.Objects[0].Meshes[0].Faces[0]
			Expect(face.References).To(HaveLen(3))

			Expect(face.References[0]).To(Equal(obj.Reference{
				VertexIndex:   0,
				TexCoordIndex: 3,
				NormalIndex:   0,
			}))
			Expect(face.References[0].HasTexCoord()).To(BeTrue())
			Expect(face.References[0].HasNormal()).To(BeTrue())

			Expect(face.References[1]).To(Equal(obj.Reference{
				VertexIndex:   1,
				TexCoordIndex: 0,
				NormalIndex:   0,
			}))
			Expect(face.References[1].HasTexCoord()).To(BeTrue())
			Expect(face.References[1].HasNormal()).To(BeTrue())

			Expect(face.References[2]).To(Equal(obj.Reference{
				VertexIndex:   2,
				TexCoordIndex: 2,
				NormalIndex:   1,
			}))
			Expect(face.References[2].HasTexCoord()).To(BeTrue())
			Expect(face.References[2].HasNormal()).To(BeTrue())
		})
	})

	When("a file with multiple objects is decoded", func() {
		BeforeEach(func() {
			testFile = "valid_objects.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all objects", func() {
			Expect(model.Objects).To(HaveLen(2))
			Expect(model.Objects[0].Name).To(Equal("First"))
			Expect(model.Objects[1].Name).To(Equal("Second"))
		})

		When("the number of objects is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxObjectCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	When("a file with all kinds of vertices is decoded", func() {
		BeforeEach(func() {
			testFile = "valid_vertices.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all vertices", func() {
			Expect(model.Vertices).To(HaveLen(2))
			Expect(model.Vertices[0]).To(Equal(obj.Vertex{
				X: 1.0, Y: 2.0, Z: 3.0, W: 1.0,
			}))
			Expect(model.Vertices[1]).To(Equal(obj.Vertex{
				X: 4.0, Y: 5.0, Z: 6.0, W: 7.0,
			}))
		})

		When("the number of vertices is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxVertexCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	When("a file with all kinds of texture coordinates is decoded", func() {
		BeforeEach(func() {
			testFile = "valid_tex_coords.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all texture coordinates", func() {
			Expect(model.TexCoords).To(HaveLen(3))
			Expect(model.TexCoords[0]).To(Equal(obj.TexCoord{
				U: 0.1, V: 0.0, W: 0.0,
			}))
			Expect(model.TexCoords[1]).To(Equal(obj.TexCoord{
				U: 0.4, V: 0.5, W: 0.0,
			}))
			Expect(model.TexCoords[2]).To(Equal(obj.TexCoord{
				U: 0.7, V: 0.8, W: 0.9,
			}))
		})

		When("the number of texture coordinates is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxTexCoordCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	When("a file with all kinds of normals is decoded", func() {
		BeforeEach(func() {
			testFile = "valid_normals.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all normals", func() {
			Expect(model.Normals).To(HaveLen(2))
			Expect(model.Normals[0]).To(Equal(obj.Normal{
				X: 1.0, Y: -2.0, Z: 3.0,
			}))
			Expect(model.Normals[1]).To(Equal(obj.Normal{
				X: -0.5, Y: 0.4, Z: -0.7,
			}))
		})

		When("the number of normals is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxNormalCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	When("a file with all kinds of material libraries is scanned", func() {
		BeforeEach(func() {
			testFile = "valid_material_libraries.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all material libraries", func() {
			Expect(model.MaterialLibraries).To(HaveLen(3))
			Expect(model.MaterialLibraries[0]).To(Equal("first.mtl"))
			Expect(model.MaterialLibraries[1]).To(Equal("second.mtl"))
			Expect(model.MaterialLibraries[2]).To(Equal("third.mtl"))
		})

		When("the number of material libraries is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxMaterialLibraryCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	When("a file with multiple material references is scanned", func() {
		BeforeEach(func() {
			testFile = "valid_material_references.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all material references", func() {
			object := model.Objects[0]
			Expect(object.Meshes).To(HaveLen(2))
			Expect(object.Meshes[0].MaterialName).To(Equal("Red"))
			Expect(object.Meshes[1].MaterialName).To(Equal("Blue"))
		})

		When("the number of material references is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxMaterialReferenceCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	When("a file with multiple faces is scanned", func() {
		BeforeEach(func() {
			testFile = "valid_faces.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all faces", func() {
			mesh := model.Objects[0].Meshes[0]
			Expect(mesh.Faces).To(HaveLen(2))
		})

		When("the number of faces is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxFaceCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	When("a file with multiple references is scanned", func() {
		BeforeEach(func() {
			testFile = "valid_references.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all faces", func() {
			face := model.Objects[0].Meshes[0].Faces[0]
			Expect(face.References).To(HaveLen(4))
		})

		When("the number of faces is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxReferenceCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	When("an object with duplicate materials is scanned", func() {
		BeforeEach(func() {
			limits.MaxMaterialReferenceCount = 2
			testFile = "valid_mesh_reuse.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have parsed only two meshes", func() {
			object := model.Objects[0]
			Expect(object.Meshes).To(HaveLen(2))
			Expect(object.Meshes[0].MaterialName).To(Equal("First"))
			Expect(object.Meshes[0].Faces).To(HaveLen(2))
			Expect(object.Meshes[1].MaterialName).To(Equal("Second"))
			Expect(object.Meshes[1].Faces).To(HaveLen(1))
		})
	})

	When("an object without a mesh is scanned", func() {
		BeforeEach(func() {
			testFile = "valid_no_mesh.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have created a default mesh", func() {
			meshes := model.Objects[0].Meshes
			Expect(meshes).To(HaveLen(1))
			Expect(meshes[0]).ToNot(BeNil())
			Expect(meshes[0].MaterialName).To(Equal(""))
		})
	})

	When("faces without an object are scanned", func() {
		BeforeEach(func() {
			testFile = "valid_no_object.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have created a default object", func() {
			Expect(model.Objects).To(HaveLen(1))
			Expect(model.Objects[0].Name).To(Equal("Default"))

			Expect(model.Objects[0].Meshes).To(HaveLen(1))
			Expect(model.Objects[0].Meshes[0].MaterialName).To(Equal("Red"))
		})
	})

	When("faces without an object and mesh are scanned", func() {
		BeforeEach(func() {
			testFile = "valid_no_object_no_mesh.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have created a default object and default mesh", func() {
			Expect(model.Objects).To(HaveLen(1))
			Expect(model.Objects[0].Name).To(Equal("Default"))

			Expect(model.Objects[0].Meshes).To(HaveLen(1))
			Expect(model.Objects[0].Meshes[0].MaterialName).To(Equal(""))
		})
	})

	When("negative indices are scanned", func() {
		BeforeEach(func() {
			testFile = "valid_negative_indices.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("have normalized the references", func() {
			face := model.Objects[0].Meshes[0].Faces[0]
			Expect(face.References).To(HaveLen(3))
			Expect(face.References[0]).To(Equal(obj.Reference{
				VertexIndex:   0,
				TexCoordIndex: 3,
				NormalIndex:   0,
			}))
			Expect(face.References[1]).To(Equal(obj.Reference{
				VertexIndex:   1,
				TexCoordIndex: 0,
				NormalIndex:   0,
			}))
			Expect(face.References[2]).To(Equal(obj.Reference{
				VertexIndex:   2,
				TexCoordIndex: 2,
				NormalIndex:   1,
			}))
		})
	})

	When("decoding face without enough references", func() {
		BeforeEach(func() {
			testFile = "error_missing_face_data.obj"
		})

		itShouldHaveReturnedAnError()
	})
})

var _ = Describe("DecodeLimits", func() {
	var limits obj.DecodeLimits

	Describe("DefaultLimits", func() {
		BeforeEach(func() {
			limits = obj.DefaultLimits()
		})

		It("vertex limit should be 65536", func() {
			Expect(limits.MaxVertexCount).To(Equal(65536))
		})

		It("normal limit should be 65536", func() {
			Expect(limits.MaxNormalCount).To(Equal(65536))
		})

		It("tex coord limit should be 65536", func() {
			Expect(limits.MaxTexCoordCount).To(Equal(65536))
		})

		It("object limit should be 1024", func() {
			Expect(limits.MaxObjectCount).To(Equal(1024))
		})

		It("face limit should be 65536", func() {
			Expect(limits.MaxFaceCount).To(Equal(65536))
		})

		It("reference limit should be 16", func() {
			Expect(limits.MaxReferenceCount).To(Equal(16))
		})

		It("material reference limit should be 64", func() {
			Expect(limits.MaxMaterialReferenceCount).To(Equal(64))
		})

		It("material library limit should be 32", func() {
			Expect(limits.MaxMaterialLibraryCount).To(Equal(32))
		})
	})
})
