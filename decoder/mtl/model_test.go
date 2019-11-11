package mtl_test

import (
	. "github.com/DanTulovsky/go-data-front/decoder/mtl"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Model", func() {
	Describe("DefaultMaterial", func() {
		var material *Material

		BeforeEach(func() {
			material = DefaultMaterial()
		})

		It("should have white ambient color", func() {
			Ω(material.AmbientColor).Should(Equal(RGBColor{
				R: 1.0,
				G: 1.0,
				B: 1.0,
			}))
		})

		It("should have white diffuse color", func() {
			Ω(material.DiffuseColor).Should(Equal(RGBColor{
				R: 1.0,
				G: 1.0,
				B: 1.0,
			}))
		})

		It("should have no specular color", func() {
			Ω(material.SpecularColor).Should(Equal(RGBColor{
				R: 0.0,
				G: 0.0,
				B: 0.0,
			}))
		})

		It("should have no emissive color", func() {
			Ω(material.EmissiveColor).Should(Equal(RGBColor{
				R: 0.0,
				G: 0.0,
				B: 0.0,
			}))
		})

		It("should have a factor of 1.0 dissolve", func() {
			Ω(material.Dissolve).Should(Equal(1.0))
		})

		It("should have a white transmission filter", func() {
			Ω(material.TransmissionFilter).Should(Equal(RGBColor{
				R: 1.0,
				G: 1.0,
				B: 1.0,
			}))
		})
	})

	Describe("Library", func() {
		var library *Library

		BeforeEach(func() {
			library = new(Library)
		})

		Context("when library has multiple materials", func() {
			var blueMaterial *Material
			var redMaterial *Material

			BeforeEach(func() {
				blueMaterial = &Material{
					Name: "Blue",
				}
				redMaterial = &Material{
					Name: "Red",
				}
				library.Materials = []*Material{blueMaterial, redMaterial}
			})

			It("is possible to find existing material by name", func() {
				material, found := library.FindMaterial("Blue")
				Ω(found).Should(BeTrue())
				Ω(material).Should(Equal(blueMaterial))

				material, found = library.FindMaterial("Red")
				Ω(found).Should(BeTrue())
				Ω(material).Should(Equal(redMaterial))
			})

			It("is will not find unexisting materials", func() {
				_, found := library.FindMaterial("Green")
				Ω(found).Should(BeFalse())
			})
		})
	})
})
