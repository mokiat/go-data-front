package mtl_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mokiat/go-data-front/decoder/mtl"
)

var _ = Describe("Material", func() {
	var material *mtl.Material

	Describe("DefaultMaterial", func() {
		BeforeEach(func() {
			material = mtl.DefaultMaterial()
		})

		It("should have white ambient color", func() {
			Expect(material.AmbientColor).To(Equal(mtl.RGBColor{
				R: 1.0,
				G: 1.0,
				B: 1.0,
			}))
		})

		It("should have white diffuse color", func() {
			Expect(material.DiffuseColor).To(Equal(mtl.RGBColor{
				R: 1.0,
				G: 1.0,
				B: 1.0,
			}))
		})

		It("should have no specular color", func() {
			Expect(material.SpecularColor).To(Equal(mtl.RGBColor{
				R: 0.0,
				G: 0.0,
				B: 0.0,
			}))
		})

		It("should have no emissive color", func() {
			Expect(material.EmissiveColor).To(Equal(mtl.RGBColor{
				R: 0.0,
				G: 0.0,
				B: 0.0,
			}))
		})

		It("should have a factor of 1.0 dissolve", func() {
			Expect(material.Dissolve).To(Equal(1.0))
		})

		It("should have a white transmission filter", func() {
			Expect(material.TransmissionFilter).To(Equal(mtl.RGBColor{
				R: 1.0,
				G: 1.0,
				B: 1.0,
			}))
		})
	})
})

var _ = Describe("Library", func() {
	var library *mtl.Library

	BeforeEach(func() {
		library = new(mtl.Library)
	})

	When("a library has multiple materials", func() {
		var (
			blueMaterial *mtl.Material
			redMaterial  *mtl.Material
		)

		BeforeEach(func() {
			blueMaterial = &mtl.Material{
				Name: "Blue",
			}
			redMaterial = &mtl.Material{
				Name: "Red",
			}
			library.Materials = []*mtl.Material{blueMaterial, redMaterial}
		})

		It("is possible to find existing material by name", func() {
			material, found := library.FindMaterial("Blue")
			Expect(found).To(BeTrue())
			Expect(material).To(Equal(blueMaterial))

			material, found = library.FindMaterial("Red")
			Expect(found).To(BeTrue())
			Expect(material).To(Equal(redMaterial))
		})

		It("is will not find unexisting materials", func() {
			_, found := library.FindMaterial("Green")
			Expect(found).To(BeFalse())
		})
	})
})
