package obj_test

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/momchil-atanasov/go-data-front/common"
	. "github.com/momchil-atanasov/go-data-front/common/common_test_help"
	. "github.com/momchil-atanasov/go-data-front/scanner/obj"

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
		file, err := os.Open(fmt.Sprintf("obj_test_res/%s", filename))
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

	// itShouldHaveReturnedAnError := func() {
	// 	It("should have returned an error", func() {
	// 		Ω(scanErr).Should(HaveOccurred())
	// 	})
	// }
	//
	// itShouldHaveReturnedHandlerError := func() {
	// 	It("should have returned handler error", func() {
	// 		Ω(scanErr).Should(Equal(errorHandlerErr))
	// 	})
	// }

	assertEvent := func(expected interface{}) {
		Ω(len(handlerTracker.Events)).Should(BeNumerically(">", eventCounter))
		Ω(handlerTracker.Events[eventCounter]).Should(Equal(expected))
		eventCounter++
	}

	assertNoMoreEvents := func() {
		Ω(handlerTracker.Events).Should(HaveLen(eventCounter))
	}

	Describe("basic OBJ file", func() {
		BeforeEach(func() {
			scanFile("valid_basic.obj", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned elements in order", func() {
			assertEvent(common.CommentEvent{
				Comment: "This is the beginning of this OBJ file.",
			})
			assertEvent(MaterialLibraryEvent{
				FilePath: "valid_basic.mtl",
			})
			assertEvent(VertexEvent{
				X: -1.0, Y: 1.0, Z: -1.0, W: 1.0,
			})
			assertEvent(VertexEvent{
				X: -1.0, Y: -1.0, Z: 1.0, W: 1.0,
			})
			assertEvent(VertexEvent{
				X: 1.0, Y: -1.0, Z: -1.0, W: 1.0,
			})
			assertEvent(VertexEvent{
				X: 1.0, Y: 1.0, Z: 1.0, W: 1.0,
			})
			assertEvent(TexCoordEvent{
				U: 0.0, V: 0.0, W: 0.0,
			})
			assertEvent(TexCoordEvent{
				U: 1.0, V: 1.0, W: 0.0,
			})
			assertEvent(TexCoordEvent{
				U: 1.0, V: 0.0, W: 0.0,
			})
			assertEvent(TexCoordEvent{
				U: 0.0, V: 1.0, W: 0.0,
			})
			assertEvent(NormalEvent{
				X: 0.0, Y: 1.0, Z: 0.0,
			})
			assertEvent(NormalEvent{
				X: 1.0, Y: 0.0, Z: 0.0,
			})
			assertEvent(NormalEvent{
				X: 0.0, Y: 0.0, Z: 1.0,
			})
			assertEvent(ObjectEvent{
				ObjectName: "MyObject",
			})
			assertEvent(MaterialReferenceEvent{
				MaterialName: "BlueMaterial",
			})
			assertNoMoreEvents()
		})
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
