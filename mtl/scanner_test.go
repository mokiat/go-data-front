package mtl_test

import (
	"errors"

	"github.com/momchil-atanasov/go-data-front/common/common_test_help"
	"github.com/momchil-atanasov/go-data-front/mtl/mtl_test_help"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Scanner", func() {
	var handlerFixture *mtl_test_help.ScannerFixture

	BeforeEach(func() {
		handlerFixture = mtl_test_help.NewScannerFixture()
	})

	Describe("basic MTL file", func() {
		var handlerStopErr error

		itScannerShouldHaveReturnedHandlerStopError := func() {
			It("scanner should have returned handler error", func() {
				handlerFixture.AssertScannerReturnedError(handlerStopErr)
			})
		}

		BeforeEach(func() {
			handlerStopErr = errors.New("Stop!")
		})

		JustBeforeEach(func() {
			handlerFixture.ScanFile("valid_basic.mtl")
		})

		Context("when handler behaves properly", func() {
			It("scanner should not have returned error", func() {
				handlerFixture.AssertScannerDidNotReturnError()
			})

			It("should have scanned the comments", func() {
				handlerFixture.AssertCommentCall("This is the beginning of this MTL file.")
				handlerFixture.AssertNoMoreCommentCalls()
			})

			It("should have scanned the material declarations", func() {
				handlerFixture.AssertMaterialCall("MyMaterial")
				handlerFixture.AssertNoMoreMaterialCalls()
			})
		})

		Context("when handler returns error on comment", func() {
			BeforeEach(func() {
				handlerFixture.Handler().OnCommentReturns(handlerStopErr)
			})

			itScannerShouldHaveReturnedHandlerStopError()
		})

		Context("when handler returns error on material", func() {
			BeforeEach(func() {
				handlerFixture.Handler().OnMaterialReturns(handlerStopErr)
			})

			itScannerShouldHaveReturnedHandlerStopError()
		})
	})

	Context("when reading fails", func() {
		var readerErr error

		BeforeEach(func() {
			readerErr = errors.New("Failed to read!")
			reader := common_test_help.NewFailingReader(readerErr)
			handlerFixture.Scan(reader)
		})

		It("scanner should have returned reader error", func() {
			handlerFixture.AssertScannerReturnedError(readerErr)
		})
	})
})
