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

	itDidNotReturnAnError := func() {
		It("scanner should not have returned error", func() {
			Ω(scanErr).ShouldNot(HaveOccurred())
		})
	}

	assertEvent := func(expected interface{}) {
		Ω(len(handlerTracker.Events)).Should(BeNumerically(">", eventCounter))
		Ω(handlerTracker.Events[eventCounter]).Should(BeEquivalentTo(expected))
		eventCounter++
	}

	assertNoMoreEvents := func() {
		Ω(handlerTracker.Events).Should(HaveLen(eventCounter))
	}

	Describe("basic MTL file", func() {
		JustBeforeEach(func() {
			scanFile("valid_basic.mtl", trackedHandler)
		})

		Context("when handler behaves properly", func() {
			itDidNotReturnAnError()

			It("should have scanned elements in order", func() {
				assertEvent(common.CommentEvent{
					Comment: "This is the beginning of this MTL file.",
				})
				assertEvent(MaterialEvent{
					MaterialName: "MyMaterial",
				})
				assertNoMoreEvents()
			})
		})
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
