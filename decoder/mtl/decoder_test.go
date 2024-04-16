package mtl_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mokiat/go-data-front/decoder/mtl"
)

var _ = Describe("Decoder", func() {
	var (
		testFile string
		limits   mtl.DecodeLimits

		library   *mtl.Library
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

		decoder := mtl.NewDecoder(limits)
		library, decodeErr = decoder.Decode(file)
	})

	BeforeEach(func() {
		limits = mtl.DefaultLimits()
	})

	When("a basic file is decoded", func() {
		BeforeEach(func() {
			testFile = "valid_basic.mtl"
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded one single material", func() {
			Expect(library.Materials).To(HaveLen(1))
		})

		Context("material data", func() {
			var material *mtl.Material

			JustBeforeEach(func() {
				material = library.Materials[0]
			})

			It("should have decoded material name", func() {
				Expect(material.Name).To(Equal("TestMaterial"))
			})

			It("should have decoded ambient color", func() {
				Expect(material.AmbientColor).To(Equal(mtl.RGBColor{
					R: 1.0,
					G: 0.5,
					B: 0.1,
				}))
			})

			It("should have decoded diffuse color", func() {
				Expect(material.DiffuseColor).To(Equal(mtl.RGBColor{
					R: 0.5,
					G: 0.7,
					B: 0.3,
				}))
			})

			It("should have decoded specular color", func() {
				Expect(material.SpecularColor).To(Equal(mtl.RGBColor{
					R: 0.2,
					G: 0.4,
					B: 0.8,
				}))
			})

			It("should have decoded transmission filter", func() {
				Expect(material.TransmissionFilter).To(Equal(mtl.RGBColor{
					R: 0.3,
					G: 1.0,
					B: 0.4,
				}))
			})

			It("should have decoded emissive color", func() {
				Expect(material.EmissiveColor).To(Equal(mtl.RGBColor{
					R: 0.4,
					G: 0.3,
					B: 0.9,
				}))
			})

			It("should have decoded specular exponent", func() {
				Expect(material.SpecularExponent).To(Equal(650.0))
			})

			It("should have decoded dissolve", func() {
				Expect(material.Dissolve).To(Equal(0.7))
			})

			It("should have decoded illumination model", func() {
				Expect(material.Illumination).To(Equal(int64(2)))
			})

			It("should have decoded ambient texture", func() {
				Expect(material.AmbientTexture).To(Equal("ambient.png"))
			})

			It("should have decoded diffuse texture", func() {
				Expect(material.DiffuseTexture).To(Equal("diffuse.png"))
			})

			It("should have decoded specular texture", func() {
				Expect(material.SpecularTexture).To(Equal("specular.png"))
			})

			It("should have decoded specular exponent texture", func() {
				Expect(material.SpecularExponentTexture).To(Equal("specular_exponent.png"))
			})

			It("should have decoded dissolve texture", func() {
				Expect(material.DissolveTexture).To(Equal("dissolve.png"))
			})

			It("should have decoded emissive texture", func() {
				Expect(material.EmissiveTexture).To(Equal("emissive.png"))
			})

			It("should have decoded bump texture", func() {
				Expect(material.BumpTexture).To(Equal("bump.png"))
			})
		})
	})

	When("a file with multiple materials is decoded", func() {
		BeforeEach(func() {
			testFile = "valid_multiple_materials.mtl"
		})

		itShouldNotHaveReturnedAnError()

		It("should have decoded all materials", func() {
			Expect(library.Materials).To(HaveLen(2))
			Expect(library.Materials[0].Name).To(Equal("FirstMaterial"))
			Expect(library.Materials[1].Name).To(Equal("SecondMaterial"))
		})

		When("the number of materials is larger than the limit", func() {
			BeforeEach(func() {
				limits.MaxMaterialCount = 1
			})

			itShouldHaveReturnedAnError()
		})
	})

	When("decoding ambient color without material", func() {
		BeforeEach(func() {
			testFile = "error_ambient_color_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding diffuse color without material", func() {
		BeforeEach(func() {
			testFile = "error_diffuse_color_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding specular color without material", func() {
		BeforeEach(func() {
			testFile = "error_specular_color_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding transmission filter without material", func() {
		BeforeEach(func() {
			testFile = "error_transmission_filter_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding emissive color without material", func() {
		BeforeEach(func() {
			testFile = "error_emissive_color_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding specular exponent without material", func() {
		BeforeEach(func() {
			testFile = "error_specular_exponent_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding dissolve without material", func() {
		BeforeEach(func() {
			testFile = "error_dissolve_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding illumination without material", func() {
		BeforeEach(func() {
			testFile = "error_illumination_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding ambient texture without material", func() {
		BeforeEach(func() {
			testFile = "error_ambient_texture_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding diffuse texture without material", func() {
		BeforeEach(func() {
			testFile = "error_diffuse_texture_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding specular texture without material", func() {
		BeforeEach(func() {
			testFile = "error_specular_texture_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding specular exponent texture without material", func() {
		BeforeEach(func() {
			testFile = "error_specular_exponent_texture_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding dissolve texture without material", func() {
		BeforeEach(func() {
			testFile = "error_dissolve_texture_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding emissive texture without material", func() {
		BeforeEach(func() {
			testFile = "error_emissive_texture_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("decoding bump texture without material", func() {
		BeforeEach(func() {
			testFile = "error_bump_texture_no_material.mtl"
		})

		itShouldHaveReturnedAnError()
	})
})

var _ = Describe("DecodeLimits", func() {
	var limits mtl.DecodeLimits

	Describe("DefaultLimits", func() {
		BeforeEach(func() {
			limits = mtl.DefaultLimits()
		})

		Specify("default material limit should be 512", func() {
			Expect(limits.MaxMaterialCount).To(Equal(512))
		})
	})
})
