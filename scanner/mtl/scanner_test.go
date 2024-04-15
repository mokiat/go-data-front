package mtl_test

import (
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mokiat/go-data-front/common"
	"github.com/mokiat/go-data-front/internal/testutil"
	"github.com/mokiat/go-data-front/scanner/mtl"
)

var _ = Describe("Scanner", func() {
	var errStubbed = errors.New("stubbed to fail")

	var (
		testFile       string
		handler        common.EventHandler
		trackedHandler *testutil.EventHandlerTracker
		eventCounter   int
		errorHandler   common.EventHandler
		scanErr        error
	)

	itShouldNotHaveReturnedAnError := func() {
		GinkgoHelper()
		It("scanner should not have returned error", func() {
			Expect(scanErr).ToNot(HaveOccurred())
		})
	}

	itShouldHaveReturnedAnError := func() {
		GinkgoHelper()
		It("should have returned an error", func() {
			Expect(scanErr).To(HaveOccurred())
		})
	}

	itShouldHaveReturnedHandlerError := func() {
		GinkgoHelper()
		It("should have returned handler error", func() {
			Expect(scanErr).To(Equal(errStubbed))
		})
	}

	assertEvent := func(expected interface{}) {
		GinkgoHelper()
		Expect(len(trackedHandler.Events)).To(BeNumerically(">", eventCounter))
		Expect(trackedHandler.Events[eventCounter]).To(Equal(expected))
		eventCounter++
	}

	assertNoMoreEvents := func() {
		GinkgoHelper()
		Expect(trackedHandler.Events).To(HaveLen(eventCounter))
	}

	JustBeforeEach(func() {
		file, err := os.Open(filepath.Join("testdata", testFile))
		Expect(err).ToNot(HaveOccurred())
		defer file.Close()

		scanner := mtl.NewScanner()
		scanErr = scanner.Scan(file, handler)
	})

	BeforeEach(func() {
		trackedHandler = new(testutil.EventHandlerTracker)
		eventCounter = 0

		errorHandler = func(event common.Event) error {
			return errStubbed
		}

		handler = trackedHandler.Handle
	})

	Describe("basic MTL file", func() {
		BeforeEach(func() {
			testFile = "valid_basic.mtl"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned elements in order", func() {
			assertEvent(common.CommentEvent{
				Comment: "This is the beginning of this MTL file.",
			})
			assertEvent(mtl.MaterialEvent{
				MaterialName: "MyMaterial",
			})
			assertEvent(mtl.RGBAmbientColorEvent{
				R: 0.8,
				G: 0.5,
				B: 0.2,
			})
			assertEvent(mtl.RGBDiffuseColorEvent{
				R: 0.1,
				G: 0.4,
				B: 0.7,
			})
			assertEvent(mtl.RGBSpecularColorEvent{
				R: 0.3,
				G: 0.2,
				B: 1.0,
			})
			assertEvent(mtl.RGBTransmissionFilterEvent{
				R: 0.6,
				G: 0.7,
				B: 0.8,
			})
			assertEvent(mtl.RGBEmissiveColorEvent{
				R: 0.4,
				G: 0.3,
				B: 0.9,
			})
			assertEvent(mtl.DissolveEvent{
				Amount: 0.4,
			})
			assertEvent(mtl.SpecularExponentEvent{
				Amount: 330.0,
			})
			assertEvent(mtl.AmbientTextureEvent{
				TexturePath: "textures/ambient.bmp",
			})
			assertEvent(mtl.DiffuseTextureEvent{
				TexturePath: "textures/diffuse.bmp",
			})
			assertEvent(mtl.SpecularTextureEvent{
				TexturePath: "textures/specular.bmp",
			})
			assertEvent(mtl.SpecularExponentTextureEvent{
				TexturePath: "textures/specular_exponent.bmp",
			})
			assertEvent(mtl.DissolveTextureEvent{
				TexturePath: "textures/dissolve.bmp",
			})
			assertEvent(mtl.EmissiveTextureEvent{
				TexturePath: "textures/emissive.png",
			})
			assertEvent(mtl.BumpTextureEvent{
				TexturePath: "textures/bump.png",
			})
			assertNoMoreEvents()
		})
	})

	When("reading all kinds of ambient colors", func() {
		BeforeEach(func() {
			testFile = "valid_ambient_colors.mtl"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned all the colors", func() {
			assertEvent(mtl.RGBAmbientColorEvent{
				R: 0.3,
				G: 0.3,
				B: 0.3,
			})
			assertEvent(mtl.RGBAmbientColorEvent{
				R: 0.2,
				G: 0.2,
				B: 0.2,
			})
			assertEvent(mtl.RGBAmbientColorEvent{
				R: 0.5,
				G: 0.6,
				B: 0.7,
			})
			assertNoMoreEvents()
		})
	})

	When("reading all kinds of diffuse colors", func() {
		BeforeEach(func() {
			testFile = "valid_diffuse_colors.mtl"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned all the colors", func() {
			assertEvent(mtl.RGBDiffuseColorEvent{
				R: 0.3,
				G: 0.3,
				B: 0.3,
			})
			assertEvent(mtl.RGBDiffuseColorEvent{
				R: 0.2,
				G: 0.2,
				B: 0.2,
			})
			assertEvent(mtl.RGBDiffuseColorEvent{
				R: 0.5,
				G: 0.6,
				B: 0.7,
			})
			assertNoMoreEvents()
		})
	})

	When("reading all kinds of specular colors", func() {
		BeforeEach(func() {
			testFile = "valid_specular_colors.mtl"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned all the colors", func() {
			assertEvent(mtl.RGBSpecularColorEvent{
				R: 0.3,
				G: 0.3,
				B: 0.3,
			})
			assertEvent(mtl.RGBSpecularColorEvent{
				R: 0.2,
				G: 0.2,
				B: 0.2,
			})
			assertEvent(mtl.RGBSpecularColorEvent{
				R: 0.5,
				G: 0.6,
				B: 0.7,
			})
			assertNoMoreEvents()
		})
	})

	When("reading all kinds of transmission filters", func() {
		BeforeEach(func() {
			testFile = "valid_transmission_filters.mtl"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned all the filters", func() {
			assertEvent(mtl.RGBTransmissionFilterEvent{
				R: 0.5,
				G: 0.6,
				B: 0.7,
			})
			assertEvent(mtl.RGBTransmissionFilterEvent{
				R: 0.1,
				G: 0.9,
				B: 0.2,
			})
			assertNoMoreEvents()
		})
	})

	When("reading all kinds of emissive colors", func() {
		BeforeEach(func() {
			testFile = "valid_emissive_colors.mtl"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned all the colors", func() {
			assertEvent(mtl.RGBEmissiveColorEvent{
				R: 0.3,
				G: 0.3,
				B: 0.3,
			})
			assertEvent(mtl.RGBEmissiveColorEvent{
				R: 0.2,
				G: 0.2,
				B: 0.2,
			})
			assertEvent(mtl.RGBEmissiveColorEvent{
				R: 0.5,
				G: 0.6,
				B: 0.7,
			})
			assertNoMoreEvents()
		})
	})

	When("reading unsupported declarations", func() {
		BeforeEach(func() {
			testFile = "valid_unsupported_declarations.mtl"
		})

		itShouldNotHaveReturnedAnError()

		It("should not have scanned anything", func() {
			assertNoMoreEvents()
		})
	})

	When("reading material without name", func() {
		BeforeEach(func() {
			testFile = "error_missing_material_name.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading ambient color without enough values", func() {
		BeforeEach(func() {
			testFile = "error_missing_ambient_color_values.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading ambient color with invalid values", func() {
		BeforeEach(func() {
			testFile = "error_invalid_ambient_color_values.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading diffuse color without enough values", func() {
		BeforeEach(func() {
			testFile = "error_missing_diffuse_color_values.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading diffuse color with invalid values", func() {
		BeforeEach(func() {
			testFile = "error_invalid_diffuse_color_values.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading specular color without enough values", func() {
		BeforeEach(func() {
			testFile = "error_missing_specular_color_values.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading specular color with invalid values", func() {
		BeforeEach(func() {
			testFile = "error_invalid_specular_color_values.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading transmission filter without enough values", func() {
		BeforeEach(func() {
			testFile = "error_missing_transmission_filter_value.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading transmission filter with invalid values", func() {
		BeforeEach(func() {
			testFile = "error_invalid_transmission_filter_values.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading emissive color without enough values", func() {
		BeforeEach(func() {
			testFile = "error_missing_emissive_color_values.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading emissive color with invalid values", func() {
		BeforeEach(func() {
			testFile = "error_invalid_emissive_color_values.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading dissolve without value", func() {
		BeforeEach(func() {
			testFile = "error_missing_dissolve_value.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading dissolve with invalid value", func() {
		BeforeEach(func() {
			testFile = "error_invalid_dissolve_value.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading specular exponent without value", func() {
		BeforeEach(func() {
			testFile = "error_missing_specular_exponent_value.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading specular exponent with invalid value", func() {
		BeforeEach(func() {
			testFile = "error_invalid_specular_exponent_value.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading ambient texture without filename param", func() {
		BeforeEach(func() {
			testFile = "error_missing_ambient_texture_filename.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading diffuse texture without filename param", func() {
		BeforeEach(func() {
			testFile = "error_missing_diffuse_texture_filename.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading specular texture without filename param", func() {
		BeforeEach(func() {
			testFile = "error_missing_specular_texture_filename.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading specular exponent texture without filename param", func() {
		BeforeEach(func() {
			testFile = "error_missing_specular_exponent_texture_filename.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading dissolve texture without filename param", func() {
		BeforeEach(func() {
			testFile = "error_missing_dissolve_texture_filename.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading emissive texture without filename param", func() {
		BeforeEach(func() {
			testFile = "error_missing_emissive_texture_filename.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("reading bump texture without filename param", func() {
		BeforeEach(func() {
			testFile = "error_missing_bump_texture_filename.mtl"
		})

		itShouldHaveReturnedAnError()
	})

	When("handler returns an error", func() {
		BeforeEach(func() {
			handler = errorHandler
		})

		When("on ambient colors", func() {
			BeforeEach(func() {
				testFile = "valid_ambient_colors.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})

		When("on diffuse colors", func() {
			BeforeEach(func() {
				testFile = "valid_diffuse_colors.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})

		When("on specular colors", func() {
			BeforeEach(func() {
				testFile = "valid_specular_colors.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})

		When("on transmission filters", func() {
			BeforeEach(func() {
				testFile = "valid_transmission_filters.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})

		When("on emissive colors", func() {
			BeforeEach(func() {
				testFile = "valid_emissive_colors.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})

		When("on dissolves", func() {
			BeforeEach(func() {
				testFile = "valid_dissolves.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})

		When("on specular exponents", func() {
			BeforeEach(func() {
				testFile = "valid_specular_exponents.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})

		When("on ambient textures", func() {
			BeforeEach(func() {
				testFile = "valid_ambient_textures.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})

		When("on diffuse textures", func() {
			BeforeEach(func() {
				testFile = "valid_diffuse_textures.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})

		When("on specular textures", func() {
			BeforeEach(func() {
				testFile = "valid_specular_textures.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})

		When("on specular exponent textures", func() {
			BeforeEach(func() {
				testFile = "valid_specular_exponent_textures.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})

		When("on dissolve textures", func() {
			BeforeEach(func() {
				testFile = "valid_dissolve_textures.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})

		When("on emissive textures", func() {
			BeforeEach(func() {
				testFile = "valid_emissive_textures.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})

		When("on bump textures", func() {
			BeforeEach(func() {
				testFile = "valid_bump_textures.mtl"
			})

			itShouldHaveReturnedHandlerError()
		})
	})
})
