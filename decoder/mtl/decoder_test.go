package mtl_test

import (
	"fmt"
	"os"

	. "github.com/mokiat/go-data-front/decoder/mtl"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decoder", func() {
	var decoder Decoder
	var limits DecodeLimits
	var library *Library
	var decodeErr error

	decodeFile := func(filename string) {
		file, err := os.Open(fmt.Sprintf("mtl_test_res/%s", filename))
		if err != nil {
			panic(err)
		}
		defer file.Close()
		decoder = NewDecoder(limits)
		library, decodeErr = decoder.Decode(file)
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
			decodeFile("valid_basic.mtl")
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded one single material", func() {
			Ω(library.Materials).Should(HaveLen(1))
		})

		Describe("material data", func() {
			var material *Material

			BeforeEach(func() {
				material = library.Materials[0]
			})

			It("should have decoded material name", func() {
				Ω(material.Name).Should(Equal("TestMaterial"))
			})

			It("should have decoded ambient color", func() {
				Ω(material.AmbientColor).Should(Equal(RGBColor{
					R: 1.0,
					G: 0.5,
					B: 0.1,
				}))
			})

			It("should have decoded diffuse color", func() {
				Ω(material.DiffuseColor).Should(Equal(RGBColor{
					R: 0.5,
					G: 0.7,
					B: 0.3,
				}))
			})

			It("should have decoded specular color", func() {
				Ω(material.SpecularColor).Should(Equal(RGBColor{
					R: 0.2,
					G: 0.4,
					B: 0.8,
				}))
			})

			It("should have decoded transmission filter", func() {
				Ω(material.TransmissionFilter).Should(Equal(RGBColor{
					R: 0.3,
					G: 1.0,
					B: 0.4,
				}))
			})

			It("should have decoded emissive color", func() {
				Ω(material.EmissiveColor).Should(Equal(RGBColor{
					R: 0.4,
					G: 0.3,
					B: 0.9,
				}))
			})

			It("should have decoded specular exponent", func() {
				Ω(material.SpecularExponent).Should(Equal(650.0))
			})

			It("should have decoded dissolve", func() {
				Ω(material.Dissolve).Should(Equal(0.7))
			})

			It("should have decoded illum", func() {
				Ω(material.Illum).Should(Equal(int64(2)))
			})

			It("should have decoded ambient texture", func() {
				Ω(material.AmbientTexture).Should(Equal("ambient.png"))
			})

			It("should have decoded diffuse texture", func() {
				Ω(material.DiffuseTexture).Should(Equal("diffuse.png"))
			})

			It("should have decoded specular texture", func() {
				Ω(material.SpecularTexture).Should(Equal("specular.png"))
			})

			It("should have decoded specular exponent texture", func() {
				Ω(material.SpecularExponentTexture).Should(Equal("specular_exponent.png"))
			})

			It("should have decoded dissolve texture", func() {
				Ω(material.DissolveTexture).Should(Equal("dissolve.png"))
			})

			It("should have decoded emissive texture", func() {
				Ω(material.EmissiveTexture).Should(Equal("emissive.png"))
			})

			It("should have decoded bump texture", func() {
				Ω(material.BumpTexture).Should(Equal("bump.png"))
			})
		})
	})

	Context("when a file with multiple materials is decoded", func() {
		JustBeforeEach(func() {
			decodeFile("valid_multiple_materials.mtl")
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all materials", func() {
			Ω(library.Materials).Should(HaveLen(2))
			Ω(library.Materials[0].Name).Should(Equal("FirstMaterial"))
			Ω(library.Materials[1].Name).Should(Equal("SecondMaterial"))
		})

		Context("when the number of materials is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxMaterialCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	Context("when decoding ambient color without material", func() {
		BeforeEach(func() {
			decodeFile("error_ambient_color_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding diffuse color without material", func() {
		BeforeEach(func() {
			decodeFile("error_diffuse_color_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding specular color without material", func() {
		BeforeEach(func() {
			decodeFile("error_specular_color_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding transmission filter without material", func() {
		BeforeEach(func() {
			decodeFile("error_transmission_filter_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding emissive color without material", func() {
		BeforeEach(func() {
			decodeFile("error_emissive_color_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding specular exponent without material", func() {
		BeforeEach(func() {
			decodeFile("error_specular_exponent_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding dissolve without material", func() {
		BeforeEach(func() {
			decodeFile("error_dissolve_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding illum without material", func() {
		BeforeEach(func() {
			decodeFile("error_illum_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding ambient texture without material", func() {
		BeforeEach(func() {
			decodeFile("error_ambient_texture_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding diffuse texture without material", func() {
		BeforeEach(func() {
			decodeFile("error_diffuse_texture_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding specular texture without material", func() {
		BeforeEach(func() {
			decodeFile("error_specular_texture_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding specular exponent texture without material", func() {
		BeforeEach(func() {
			decodeFile("error_specular_exponent_texture_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding dissolve texture without material", func() {
		BeforeEach(func() {
			decodeFile("error_dissolve_texture_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding emissive texture without material", func() {
		BeforeEach(func() {
			decodeFile("error_emissive_texture_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when decoding bump texture without material", func() {
		BeforeEach(func() {
			decodeFile("error_bump_texture_no_material.mtl")
		})

		itShouldHaveReturnedAnError()
	})

	Describe("Default DecodeLimits", func() {
		var limits DecodeLimits

		BeforeEach(func() {
			limits = DefaultLimits()
		})

		It("material limit should be 512", func() {
			Ω(limits.MaxMaterialCount).Should(Equal(512))
		})
	})
})
