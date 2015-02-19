package mtl_test

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/momchil-atanasov/go-data-front/common"
	. "github.com/momchil-atanasov/go-data-front/common/common_test_help"
	. "github.com/momchil-atanasov/go-data-front/mtl"
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
			assertEvent(DissolveEvent{
				Amount: 0.4,
			})
			assertEvent(SpecularExponentEvent{
				Amount: 330.0,
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

	Context("when reading dissolve without value", func() {
		BeforeEach(func() {
			scanFile("error_missing_dissolve_value.mtl", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading dissolve with invalid value", func() {
		BeforeEach(func() {
			scanFile("error_invalid_dissolve_value.mtl", trackedHandler)
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

	// TODO: Handle handler errors

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
