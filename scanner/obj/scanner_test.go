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

	assertAnyEvent := func() {
		Ω(len(handlerTracker.Events)).Should(BeNumerically(">", eventCounter))
		eventCounter++
	}

	assertAnyEvents := func(count int) {
		Ω(len(handlerTracker.Events)).Should(BeNumerically(">", eventCounter+count))
		eventCounter += count
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
			assertEvent(FaceStartEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 1,
			})
			assertEvent(TexCoordReferenceEvent{
				TexCoordIndex: 4,
			})
			assertEvent(NormalReferenceEvent{
				NormalIndex: 1,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 2,
			})
			assertEvent(TexCoordReferenceEvent{
				TexCoordIndex: 1,
			})
			assertEvent(NormalReferenceEvent{
				NormalIndex: 1,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 3,
			})
			assertEvent(TexCoordReferenceEvent{
				TexCoordIndex: 3,
			})
			assertEvent(NormalReferenceEvent{
				NormalIndex: 2,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(FaceEndEvent{})
			assertNoMoreEvents()
		})
	})

	Context("when a file with all kinds of comments is scanned", func() {
		BeforeEach(func() {
			scanFile("valid_comments.obj", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the comments", func() {
			assertEvent(common.CommentEvent{
				Comment: "Comment at file start",
			})
			assertAnyEvent()
			assertEvent(common.CommentEvent{
				Comment: "Comment that is right next to special char",
			})
			assertAnyEvent()
			assertEvent(common.CommentEvent{
				Comment: "This comment uses",
			})
			assertEvent(common.CommentEvent{
				Comment: "two lines",
			})
			assertAnyEvent()
			assertEvent(common.CommentEvent{
				Comment: "",
			})
			assertAnyEvent()
			assertEvent(common.CommentEvent{
				Comment: "Previous comment was empty. This one contain the # character twice.",
			})
			assertAnyEvents(11)
			assertEvent(common.CommentEvent{
				Comment: "Comment at file end",
			})
			assertNoMoreEvents()
		})
	})

	Context("when a file with all kinds of vertices is scanned", func() {
		BeforeEach(func() {
			scanFile("valid_vertices.obj", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the vertices", func() {
			assertEvent(VertexEvent{
				X: 1.0, Y: 1.0, Z: -1.0, W: 1.0,
			})
			assertEvent(VertexEvent{
				X: -1.0, Y: -1.0, Z: 1.0, W: 0.5,
			})
			assertNoMoreEvents()
		})
	})
	Context("when a file with all kinds of texture coordinates is scanned", func() {
		BeforeEach(func() {
			scanFile("valid_texcoords.obj", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the texture coordinates", func() {
			assertEvent(TexCoordEvent{
				U: 1.6, V: 0.0, W: 0.0,
			})
			assertEvent(TexCoordEvent{
				U: 0.0, V: -0.5, W: 0.0,
			})
			assertEvent(TexCoordEvent{
				U: -0.2, V: 1.4, W: 3.0,
			})
			assertNoMoreEvents()
		})
	})

	Context("when a file with all kinds of normals is scanned", func() {
		BeforeEach(func() {
			scanFile("valid_normals.obj", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the normals", func() {
			assertEvent(NormalEvent{
				X: 1.0, Y: 0.0, Z: 0.0,
			})
			assertEvent(NormalEvent{
				X: -1.0, Y: -1.0, Z: 1.0,
			})
			assertNoMoreEvents()
		})
	})

	Context("when a file with all kinds of objects is scanned", func() {
		BeforeEach(func() {
			scanFile("valid_objects.obj", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the objects", func() {
			assertEvent(ObjectEvent{
				ObjectName: "FirstObject",
			})
			assertEvent(ObjectEvent{
				ObjectName: "LastObject",
			})
			assertNoMoreEvents()
		})
	})
	Context("when a file with all kinds of coord references is scanned", func() {
		BeforeEach(func() {
			scanFile("valid_faces.obj", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned them correctly", func() {
			// First face
			assertEvent(FaceStartEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 1,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 2,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 3,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(FaceEndEvent{})

			// Second face
			assertEvent(FaceStartEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 2,
			})
			assertEvent(TexCoordReferenceEvent{
				TexCoordIndex: 1,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 3,
			})
			assertEvent(TexCoordReferenceEvent{
				TexCoordIndex: 2,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 4,
			})
			assertEvent(TexCoordReferenceEvent{
				TexCoordIndex: 3,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(FaceEndEvent{})

			// Third face
			assertEvent(FaceStartEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 4,
			})
			assertEvent(NormalReferenceEvent{
				NormalIndex: 5,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 6,
			})
			assertEvent(NormalReferenceEvent{
				NormalIndex: 7,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 8,
			})
			assertEvent(NormalReferenceEvent{
				NormalIndex: 9,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(FaceEndEvent{})

			// Fourth face
			assertEvent(FaceStartEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 1,
			})
			assertEvent(TexCoordReferenceEvent{
				TexCoordIndex: 2,
			})
			assertEvent(NormalReferenceEvent{
				NormalIndex: 3,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 2,
			})
			assertEvent(TexCoordReferenceEvent{
				TexCoordIndex: 3,
			})
			assertEvent(NormalReferenceEvent{
				NormalIndex: 4,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(ReferenceSetStartEvent{})
			assertEvent(VertexReferenceEvent{
				VertexIndex: 3,
			})
			assertEvent(TexCoordReferenceEvent{
				TexCoordIndex: 4,
			})
			assertEvent(NormalReferenceEvent{
				NormalIndex: 5,
			})
			assertEvent(ReferenceSetEndEvent{})
			assertEvent(FaceEndEvent{})

			assertNoMoreEvents()
		})
	})

	Context("when a file with all kinds of material libraries is scanned", func() {
		BeforeEach(func() {
			scanFile("valid_material_libraries.obj", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the material libraries", func() {
			assertEvent(MaterialLibraryEvent{
				FilePath: "valid.mtl",
			})
			assertEvent(MaterialLibraryEvent{
				FilePath: "extension01.mtl",
			})
			assertEvent(MaterialLibraryEvent{
				FilePath: "extension02.mtl",
			})
			assertNoMoreEvents()
		})
	})

	Context("when a file with all kinds of material references is scanned", func() {
		BeforeEach(func() {
			scanFile("valid_material_references.obj", trackedHandler)
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the material references", func() {
			assertEvent(MaterialReferenceEvent{
				MaterialName: "",
			})
			assertEvent(MaterialReferenceEvent{
				MaterialName: "MyMaterial",
			})
			assertNoMoreEvents()
		})
	})

	Context("when a file with insufficient vertex data is scanned", func() {
		BeforeEach(func() {
			scanFile("error_insufficient_vertex_data.obj", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with insufficient texture coordinate data is scanned", func() {
		BeforeEach(func() {
			scanFile("error_insufficient_texcoord_data.obj", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with insufficient normal data is scanned", func() {
		BeforeEach(func() {
			scanFile("error_insufficient_normal_data.obj", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with an unnamed object is scanned", func() {
		BeforeEach(func() {
			scanFile("error_empty_object_name.obj", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with corrupt vertex is scanned", func() {
		BeforeEach(func() {
			scanFile("error_corrupt_vertex.obj", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with corrupt texture coordinate is scanned", func() {
		BeforeEach(func() {
			scanFile("error_corrupt_texcoord.obj", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with corrupt normal is scanned", func() {
		BeforeEach(func() {
			scanFile("error_corrupt_normal.obj", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with corrupt vertex reference is scanned", func() {
		BeforeEach(func() {
			scanFile("error_corrupt_vertex_reference.obj", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with corrupt texture coordinate reference is scanned", func() {
		BeforeEach(func() {
			scanFile("error_corrupt_texcoord_reference.obj", trackedHandler)
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with corrupt normal reference is scanned", func() {
		BeforeEach(func() {
			scanFile("error_corrupt_normal_reference.obj", trackedHandler)
		})

		itShouldHaveReturnedAnError()
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
