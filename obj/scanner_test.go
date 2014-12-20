package obj_test

import (
	"os"

	. "github.com/momchil-atanasov/go-data-front/obj"
	"github.com/momchil-atanasov/go-data-front/obj/test_helpers"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Scanner", func() {

	var handlerFixture *test_helpers.HandlerFixture
	var scanner Scanner

	var scanError error

	BeforeEach(func() {
		handlerFixture = test_helpers.NewHandlerFixture()
		scanner = NewScanner(handlerFixture.Handler())
		scanError = nil
	})

	scanFile := func(filename string) {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scanError = scanner.Scan(file)
	}

	Context("when a basic OBJ file is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_basic.obj")
		})

		It("should have scanned the comments", func() {
			handlerFixture.AssertCommentCall("This is the beginning of this OBJ file.")
			handlerFixture.AssertNoMoreCommentCalls()
		})

		It("should have scanned the vertices", func() {
			handlerFixture.AssertVertexXYZ(-1.0, 1.0, -1.0)
			handlerFixture.AssertVertexXYZ(-1.0, -1.0, 1.0)
			handlerFixture.AssertVertexXYZ(1.0, -1.0, -1.0)
			handlerFixture.AssertVertexXYZ(1.0, 1.0, 1.0)
		})

		It("should have scanned the texture coordinates", func() {
			handlerFixture.AssertTexCoordUV(0.0, 0.0)
			handlerFixture.AssertTexCoordUV(1.0, 1.0)
			handlerFixture.AssertTexCoordUV(1.0, 0.0)
			handlerFixture.AssertTexCoordUV(0.0, 1.0)
		})

		It("should have scanned the normals", func() {
			handlerFixture.AssertNormal(0.0, 1.0, 0.0)
			handlerFixture.AssertNormal(1.0, 0.0, 0.0)
			handlerFixture.AssertNormal(0.0, 0.0, 1.0)
		})

		It("should have scanned the objects", func() {
			handlerFixture.AssertObject("MyObject")
		})

		It("should have scanned the faces", func() {
			handlerFixture.AssertFaceStart()
			handlerFixture.AssertFaceEnd()
		})

		It("should have scanned data references", func() {
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(1)
			handlerFixture.AssertTexCoordIndex(4)
			handlerFixture.AssertNormalIndex(1)
			handlerFixture.AssertCoordReferenceEnd()

			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(2)
			handlerFixture.AssertTexCoordIndex(1)
			handlerFixture.AssertNormalIndex(1)
			handlerFixture.AssertCoordReferenceEnd()

			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(3)
			handlerFixture.AssertTexCoordIndex(3)
			handlerFixture.AssertNormalIndex(2)
			handlerFixture.AssertCoordReferenceEnd()
		})

		It("should have scanned material libraries", func() {
			handlerFixture.AssertMaterialLibrary("valid_basic.mtl")
		})

		It("should have scanned material references", func() {
			handlerFixture.AssertMaterialReference("BlueMaterial")
		})
	})

	Context("when a file with all kinds of comments is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_comments.obj")
		})

		It("should have scanned the comments", func() {
			handlerFixture.AssertCommentCall("Comment at file start")
			handlerFixture.AssertCommentCall("Comment that is right next to special char")
			handlerFixture.AssertCommentCall("This comment uses")
			handlerFixture.AssertCommentCall("two lines")
			handlerFixture.AssertCommentCall("")
			handlerFixture.AssertCommentCall("Previous comment was empty. This one contain the # character twice.")
			handlerFixture.AssertCommentCall("Comment at file end")
			handlerFixture.AssertNoMoreCommentCalls()
		})
	})

	Context("when a file with all kinds of coord references is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_coord_references.obj")
		})

		It("should have scanned them correctly", func() {
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(1)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(2)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(3)
			handlerFixture.AssertCoordReferenceEnd()

			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(2)
			handlerFixture.AssertTexCoordIndex(1)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(3)
			handlerFixture.AssertTexCoordIndex(2)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(4)
			handlerFixture.AssertTexCoordIndex(3)
			handlerFixture.AssertCoordReferenceEnd()

			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(4)
			handlerFixture.AssertNormalIndex(5)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(6)
			handlerFixture.AssertNormalIndex(7)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(8)
			handlerFixture.AssertNormalIndex(9)
			handlerFixture.AssertCoordReferenceEnd()

			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(1)
			handlerFixture.AssertTexCoordIndex(2)
			handlerFixture.AssertNormalIndex(3)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(2)
			handlerFixture.AssertTexCoordIndex(3)
			handlerFixture.AssertNormalIndex(4)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(3)
			handlerFixture.AssertTexCoordIndex(4)
			handlerFixture.AssertNormalIndex(5)
			handlerFixture.AssertCoordReferenceEnd()
		})
	})

})
