package obj_test

import (
	"os"

	. "github.com/momchil-atanasov/go-data-front/obj"
	"github.com/momchil-atanasov/go-data-front/obj/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scanner", func() {

	var handler *fakes.FakeScannerHandler
	var scanError error
	var commentCounter int
	var vertexStartCounter int
	var vertexXCounter int
	var vertexYCounter int
	var vertexZCounter int
	var vertexEndCounter int
	var texCoordStartCounter int
	var texCoordEndCounter int
	var texCoordUCounter int
	var texCoordVCounter int
	var texCoordWCounter int
	var normalCounter int
	var objectCounter int
	var faceStartCounter int
	var faceEndCounter int
	var materialLibraryCounter int
	var materialReferenceCounter int

	BeforeEach(func() {
		handler = new(fakes.FakeScannerHandler)
		scanError = nil

		commentCounter = 0
		vertexStartCounter = 0
		vertexEndCounter = 0
		vertexXCounter = 0
		vertexYCounter = 0
		vertexZCounter = 0
		texCoordStartCounter = 0
		texCoordEndCounter = 0
		texCoordUCounter = 0
		texCoordVCounter = 0
		texCoordWCounter = 0
		normalCounter = 0
		objectCounter = 0
		faceStartCounter = 0
		faceEndCounter = 0
		materialLibraryCounter = 0
		materialReferenceCounter = 0
	})

	scanFile := func(filename string) {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scanError = Scan(file, handler)
	}

	assertComment := func(expectedComment string) {
		Ω(handler.OnCommentCallCount()).Should(BeNumerically(">", commentCounter))
		Ω(handler.OnCommentArgsForCall(commentCounter)).Should(Equal(expectedComment))
		commentCounter++
	}

	assertVertexStart := func() {
		Ω(handler.OnVertexStartCallCount()).Should(BeNumerically(">", vertexStartCounter))
		vertexStartCounter++
	}

	assertVertexEnd := func() {
		Ω(handler.OnVertexEndCallCount()).Should(BeNumerically(">", vertexEndCounter))
		vertexEndCounter++
	}

	assertVertexX := func(expectedX float32) {
		Ω(handler.OnVertexXCallCount()).Should(BeNumerically(">", vertexXCounter))
		Ω(handler.OnVertexXArgsForCall(vertexXCounter)).Should(BeNumerically("~", expectedX, 0.0001))
		vertexXCounter++
	}

	assertVertexY := func(expectedY float32) {
		Ω(handler.OnVertexYCallCount()).Should(BeNumerically(">", vertexYCounter))
		Ω(handler.OnVertexYArgsForCall(vertexYCounter)).Should(BeNumerically("~", expectedY, 0.0001))
		vertexYCounter++
	}

	assertVertexZ := func(expectedZ float32) {
		Ω(handler.OnVertexZCallCount()).Should(BeNumerically(">", vertexZCounter))
		Ω(handler.OnVertexZArgsForCall(vertexZCounter)).Should(BeNumerically("~", expectedZ, 0.0001))
		vertexZCounter++
	}

	assertVertexXYZ := func(expectedX, expectedY, expectedZ float32) {
		assertVertexStart()
		assertVertexX(expectedX)
		assertVertexY(expectedY)
		assertVertexZ(expectedZ)
		assertVertexEnd()
	}

	assertTexCoordStart := func() {
		Ω(handler.OnTexCoordStartCallCount()).Should(BeNumerically(">", texCoordStartCounter))
		texCoordStartCounter++
	}

	assertTexCoordU := func(expectedU float32) {
		Ω(handler.OnTexCoordUCallCount()).Should(BeNumerically(">", texCoordUCounter))
		Ω(handler.OnTexCoordUArgsForCall(texCoordUCounter)).Should(BeNumerically("~", expectedU, 0.0001))
		texCoordUCounter++
	}

	assertTexCoordV := func(expectedV float32) {
		Ω(handler.OnTexCoordVCallCount()).Should(BeNumerically(">", texCoordVCounter))
		Ω(handler.OnTexCoordVArgsForCall(texCoordVCounter)).Should(BeNumerically("~", expectedV, 0.0001))
		texCoordVCounter++
	}

	assertTexCoordEnd := func() {
		Ω(handler.OnTexCoordEndCallCount()).Should(BeNumerically(">", texCoordEndCounter))
		texCoordEndCounter++
	}

	assertTexCoordUV := func(expectedU, expectedV float32) {
		assertTexCoordStart()
		assertTexCoordU(expectedU)
		assertTexCoordV(expectedV)
		assertTexCoordEnd()
	}

	assertNormal := func(expectedX, expectedY, expectedZ float32) {
		Ω(handler.OnNormalCallCount()).Should(BeNumerically(">", normalCounter))
		argX, argY, argZ := handler.OnNormalArgsForCall(normalCounter)
		Ω(argX).Should(BeNumerically("~", expectedX, 0.0001))
		Ω(argY).Should(BeNumerically("~", expectedY, 0.0001))
		Ω(argZ).Should(BeNumerically("~", expectedZ, 0.0001))
		normalCounter++
	}

	assertObject := func(expectedName string) {
		Ω(handler.OnObjectCallCount()).Should(BeNumerically(">", objectCounter))
		argName := handler.OnObjectArgsForCall(objectCounter)
		Ω(argName).Should(Equal(expectedName))
		objectCounter++
	}

	assertFaceStart := func() {
		Ω(handler.OnFaceStartCallCount()).Should(BeNumerically(">", faceStartCounter))
		faceStartCounter++
	}

	assertFaceEnd := func() {
		Ω(handler.OnFaceEndCallCount()).Should(BeNumerically(">", faceEndCounter))
		faceEndCounter++
	}

	assertMaterialLibrary := func(expectedPath string) {
		Ω(handler.OnMaterialLibraryCallCount()).Should(BeNumerically(">", materialLibraryCounter))
		pathArg := handler.OnMaterialLibraryArgsForCall(materialLibraryCounter)
		Ω(pathArg).Should(Equal(expectedPath))
		materialLibraryCounter++
	}

	assertMaterialReference := func(expectedName string) {
		Ω(handler.OnMaterialReferenceCallCount()).Should(BeNumerically(">", materialReferenceCounter))
		nameArg := handler.OnMaterialReferenceArgsForCall(materialReferenceCounter)
		Ω(nameArg).Should(Equal(expectedName))
		materialReferenceCounter++
	}

	Context("when a basic OBJ file is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_basic.obj")
		})

		It("should have scanned the comments", func() {
			assertComment("This is the beginning of this OBJ file.")
		})

		It("should have scanned the vertices", func() {
			assertVertexXYZ(-1.0, 1.0, -1.0)
			assertVertexXYZ(-1.0, -1.0, 1.0)
			assertVertexXYZ(1.0, -1.0, -1.0)
			assertVertexXYZ(1.0, 1.0, 1.0)
		})

		It("should have scanned the texture coordinates", func() {
			assertTexCoordUV(0.0, 0.0)
			assertTexCoordUV(1.0, 1.0)
			assertTexCoordUV(1.0, 0.0)
			assertTexCoordUV(0.0, 1.0)
		})

		It("should have scanned the normals", func() {
			assertNormal(0.0, 1.0, 0.0)
			assertNormal(1.0, 0.0, 0.0)
			assertNormal(0.0, 0.0, 1.0)
		})

		It("should have scanned the objects", func() {
			assertObject("MyObject")
		})

		It("should have scanned the faces", func() {
			assertFaceStart()
			assertFaceEnd()
		})

		// TODO: Test for data references

		It("should have scanned material libraries", func() {
			assertMaterialLibrary("valid_basic.mtl")
		})

		It("should have scanned material references", func() {
			assertMaterialReference("BlueMaterial")
		})
	})

	Context("when a file with all kinds of comments is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_comments.obj")
		})

		It("should have scanned the comments", func() {
			assertComment("Comment at file start")
			assertComment("Comment that is right next to special char")
			assertComment("This comment uses")
			assertComment("two lines")
			assertComment("")
			assertComment("Previous comment was empty. This one contain the # character twice.")
			assertComment("Comment at file end")
		})

	})

})
