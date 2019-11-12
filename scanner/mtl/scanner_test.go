package mtl_test

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mokiat/go-data-front/common"
	. "github.com/mokiat/go-data-front/common/common_test_help"
	. "github.com/mokiat/go-data-front/scanner/mtl"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scanner", func() {
	var handlerTracker *EventHandlerTracker
	var trackedHandler common.EventHandler
	var errorHandlerErr error
	var errorHandler common.EventHandler
	var scanErr error
	var scanner common.Scanner
	var eventCounter int

	BeforeEach(func() {
		handlerTracker = new(EventHandlerTracker)
		trackedHandler = handlerTracker.Handle
		eventCounter = 0

		errorHandlerErr = errors.New("Handler returned error!")
		errorHandler = func(event common.Event) error {
			return errorHandlerErr
		}

		scanErr = nil
		scanner = NewScanner()
	})

	scan := func(reader io.Reader, handler common.EventHandler) {
		scanErr = scanner.Scan(reader, handler)
	}

	scanFile := func(filename string, handler common.EventHandler) {
		file, err := os.Open(fmt.Sprintf("mtl_test_res/%s", filename))
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scan(file, handler)
	}

	itShouldNotHaveReturnedAnError := func() {
		It("scanner should not have returned error", func() {
			Ω(scanErr).ShouldNot(HaveOccurred())
		})
	}

	itShouldHaveReturnedAnError := func() {
		It("should have returned an error", func() {
			Ω(scanErr).Should(HaveOccurred())
		})
	}

	itShouldHaveReturnedHandlerError := func() {
		It("should have returned handler error", func() {
			Ω(scanErr).Should(Equal(errorHandlerErr))
		})
	}

	assertEvent := func(expected interface{}) {
		Ω(len(handlerTracker.Events)).Should(BeNumerically(">", eventCounter))
		Ω(handlerTracker.Events[eventCounter]).Should(Equal(expected))
		eventCounter++
	}

	assertNoMoreEvents := func() {
		Ω(handlerTracker.Events).Should(HaveLen(eventCounter))
	}

	Describe("basic MTL file", func() {
		BeforeEach(func() {
			scanFile("valid_basic.mtl", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned elements in order", func() {
			assertEvent(common.CommentEvent{
				Comment: "This is the beginning of this MTL file.",
			})
			assertEvent(MaterialEvent{
				MaterialName: "MyMaterial",
			})
			assertEvent(RGBAmbientColorEvent{
				R: 0.8,
				G: 0.5,
				B: 0.2,
			})
			assertEvent(RGBDiffuseColorEvent{
				R: 0.1,
				G: 0.4,
				B: 0.7,
			})
			assertEvent(RGBSpecularColorEvent{
				R: 0.3,
				G: 0.2,
				B: 1.0,
			})
			assertEvent(RGBTransmissionFilterEvent{
				R: 0.6,
				G: 0.7,
				B: 0.8,
			})
			assertEvent(RGBEmissiveColorEvent{
				R: 0.4,
				G: 0.3,
				B: 0.9,
			})
			assertEvent(DissolveEvent{
				Amount: 0.4,
			})
			assertEvent(SpecularExponentEvent{
				Amount: 330.0,
			})
			assertEvent(IllumEvent{
				Amount: 2,
			})
			assertEvent(AmbientTextureEvent{
				TexturePath: "textures/ambient.bmp",
			})
			assertEvent(DiffuseTextureEvent{
				TexturePath: "textures/diffuse.bmp",
			})
			assertEvent(SpecularTextureEvent{
				TexturePath: "textures/specular.bmp",
			})
			assertEvent(SpecularExponentTextureEvent{
				TexturePath: "textures/specular_exponent.bmp",
			})
			assertEvent(DissolveTextureEvent{
				TexturePath: "textures/dissolve.bmp",
			})
			assertEvent(EmissiveTextureEvent{
				TexturePath: "textures/emissive.png",
			})
			assertEvent(BumpTextureEvent{
				TexturePath: "textures/bump.png",
			})
			assertNoMoreEvents()
		})
	})

	Context("when reading all kinds of ambient colors", func() {
		BeforeEach(func() {
			scanFile("valid_ambient_colors.mtl", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned all the colors", func() {
			assertEvent(RGBAmbientColorEvent{
				R: 0.3,
				G: 0.3,
				B: 0.3,
			})
			assertEvent(RGBAmbientColorEvent{
				R: 0.2,
				G: 0.2,
				B: 0.2,
			})
			assertEvent(RGBAmbientColorEvent{
				R: 0.5,
				G: 0.6,
				B: 0.7,
			})
			assertNoMoreEvents()
		})
	})

	Context("when reading all kinds of diffuse colors", func() {
		BeforeEach(func() {
			scanFile("valid_diffuse_colors.mtl", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned all the colors", func() {
			assertEvent(RGBDiffuseColorEvent{
				R: 0.3,
				G: 0.3,
				B: 0.3,
			})
			assertEvent(RGBDiffuseColorEvent{
				R: 0.2,
				G: 0.2,
				B: 0.2,
			})
			assertEvent(RGBDiffuseColorEvent{
				R: 0.5,
				G: 0.6,
				B: 0.7,
			})
			assertNoMoreEvents()
		})
	})

	Context("when reading all kinds of specular colors", func() {
		BeforeEach(func() {
			scanFile("valid_specular_colors.mtl", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned all the colors", func() {
			assertEvent(RGBSpecularColorEvent{
				R: 0.3,
				G: 0.3,
				B: 0.3,
			})
			assertEvent(RGBSpecularColorEvent{
				R: 0.2,
				G: 0.2,
				B: 0.2,
			})
			assertEvent(RGBSpecularColorEvent{
				R: 0.5,
				G: 0.6,
				B: 0.7,
			})
			assertNoMoreEvents()
		})
	})

	Context("when reading all kinds of transmission filters", func() {
		BeforeEach(func() {
			scanFile("valid_transmission_filters.mtl", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned all the filters", func() {
			assertEvent(RGBTransmissionFilterEvent{
				R: 0.5,
				G: 0.6,
				B: 0.7,
			})
			assertEvent(RGBTransmissionFilterEvent{
				R: 0.1,
				G: 0.9,
				B: 0.2,
			})
			assertNoMoreEvents()
		})
	})

	Context("when reading all kinds of emissive colors", func() {
		BeforeEach(func() {
			scanFile("valid_emissive_colors.mtl", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned all the colors", func() {
			assertEvent(RGBEmissiveColorEvent{
				R: 0.3,
				G: 0.3,
				B: 0.3,
			})
			assertEvent(RGBEmissiveColorEvent{
				R: 0.2,
				G: 0.2,
				B: 0.2,
			})
			assertEvent(RGBEmissiveColorEvent{
				R: 0.5,
				G: 0.6,
				B: 0.7,
			})
			assertNoMoreEvents()
		})
	})

	Context("when reading unsupported declarations", func() {
		BeforeEach(func() {
			scanFile("valid_unsupported_declarations.mtl", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should not have scanned anything", func() {
			assertNoMoreEvents()
		})
	})

	Context("when reading material without name", func() {
		BeforeEach(func() {
			scanFile("error_missing_material_name.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading ambient color without enough values", func() {
		BeforeEach(func() {
			scanFile("error_missing_ambient_color_values.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading ambient color with invalid values", func() {
		BeforeEach(func() {
			scanFile("error_invalid_ambient_color_values.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading diffuse color without enough values", func() {
		BeforeEach(func() {
			scanFile("error_missing_diffuse_color_values.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading diffuse color with invalid values", func() {
		BeforeEach(func() {
			scanFile("error_invalid_diffuse_color_values.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading specular color without enough values", func() {
		BeforeEach(func() {
			scanFile("error_missing_specular_color_values.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading specular color with invalid values", func() {
		BeforeEach(func() {
			scanFile("error_invalid_specular_color_values.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading transmission filter without enough values", func() {
		BeforeEach(func() {
			scanFile("error_missing_transmission_filter_value.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading transmission filter with invalid values", func() {
		BeforeEach(func() {
			scanFile("error_invalid_transmission_filter_values.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading emissive color without enough values", func() {
		BeforeEach(func() {
			scanFile("error_missing_emissive_color_values.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading emissive color with invalid values", func() {
		BeforeEach(func() {
			scanFile("error_invalid_emissive_color_values.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading dissolve without value", func() {
		BeforeEach(func() {
			scanFile("error_missing_dissolve_value.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading illum without value", func() {
		BeforeEach(func() {
			scanFile("error_missing_illum_value.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})
	Context("when reading dissolve with invalid value", func() {
		BeforeEach(func() {
			scanFile("error_invalid_dissolve_value.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading illum with invalid value", func() {
		BeforeEach(func() {
			scanFile("error_invalid_illum_value.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})
	Context("when reading specular exponent without value", func() {
		BeforeEach(func() {
			scanFile("error_missing_specular_exponent_value.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading specular exponent with invalid value", func() {
		BeforeEach(func() {
			scanFile("error_invalid_specular_exponent_value.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading ambient texture without filename param", func() {
		BeforeEach(func() {
			scanFile("error_missing_ambient_texture_filename.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading diffuse texture without filename param", func() {
		BeforeEach(func() {
			scanFile("error_missing_diffuse_texture_filename.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading specular texture without filename param", func() {
		BeforeEach(func() {
			scanFile("error_missing_specular_texture_filename.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading specular exponent texture without filename param", func() {
		BeforeEach(func() {
			scanFile("error_missing_specular_exponent_texture_filename.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading dissolve texture without filename param", func() {
		BeforeEach(func() {
			scanFile("error_missing_dissolve_texture_filename.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading emissive texture without filename param", func() {
		BeforeEach(func() {
			scanFile("error_missing_emissive_texture_filename.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading bump texture without filename param", func() {
		BeforeEach(func() {
			scanFile("error_missing_bump_texture_filename.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when handler returns error on ambient colors", func() {
		BeforeEach(func() {
			scanFile("valid_ambient_colors.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when handler returns error on diffuse colors", func() {
		BeforeEach(func() {
			scanFile("valid_diffuse_colors.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when handler returns error on specular colors", func() {
		BeforeEach(func() {
			scanFile("valid_specular_colors.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when handler returns error on transmission filters", func() {
		BeforeEach(func() {
			scanFile("valid_transmission_filters.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when handler returns error on emissive colors", func() {
		BeforeEach(func() {
			scanFile("valid_emissive_colors.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when handler returns error on dissolves", func() {
		BeforeEach(func() {
			scanFile("valid_dissolves.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when handler returns error on illums", func() {
		BeforeEach(func() {
			scanFile("valid_illums.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})
	Context("when handler returns error on specular exponents", func() {
		BeforeEach(func() {
			scanFile("valid_specular_exponents.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when handler returns error on ambient textures", func() {
		BeforeEach(func() {
			scanFile("valid_ambient_textures.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when handler returns error on diffuse textures", func() {
		BeforeEach(func() {
			scanFile("valid_diffuse_textures.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when handler returns error on specular textures", func() {
		BeforeEach(func() {
			scanFile("valid_specular_textures.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when handler returns error on specular exponent textures", func() {
		BeforeEach(func() {
			scanFile("valid_specular_exponent_textures.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when handler returns error on dissolve textures", func() {
		BeforeEach(func() {
			scanFile("valid_dissolve_textures.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when handler returns error on emissive textures", func() {
		BeforeEach(func() {
			scanFile("valid_emissive_textures.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when handler returns error on bump textures", func() {
		BeforeEach(func() {
			scanFile("valid_bump_textures.mtl", errorHandler)
		})

		itShouldHaveReturnedHandlerError()
	})

	Context("when reading fails", func() {
		var readerErr error

		BeforeEach(func() {
			readerErr = errors.New("Failed to read!")
			reader := NewFailingReader(readerErr)
			scan(reader, trackedHandler)
		})

		It("scanner should have returned reader error", func() {
			Ω(scanErr).Should(Equal(readerErr))
		})
	})
})
