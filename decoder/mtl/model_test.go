package mtl_test

import (
	. "github.com/momchil-atanasov/go-data-front/decoder/mtl"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Model", func() {
	Describe("DefaultMaterial", func() {
		var material Material

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
})
