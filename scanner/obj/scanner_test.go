package obj_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mokiat/go-data-front/common"
	"github.com/mokiat/go-data-front/internal/testutil"
	"github.com/mokiat/go-data-front/scanner/obj"
)

var _ = Describe("Scanner", func() {
	var (
		testFile       string
		handler        common.EventHandler
		trackedHandler *testutil.EventHandlerTracker
		eventCounter   int
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

	assertEvent := func(expected interface{}) {
		GinkgoHelper()
		Expect(len(trackedHandler.Events)).To(BeNumerically(">", eventCounter))
		Expect(trackedHandler.Events[eventCounter]).To(Equal(expected))
		eventCounter++
	}

	assertAnyEvent := func() {
		GinkgoHelper()
		Expect(len(trackedHandler.Events)).To(BeNumerically(">", eventCounter))
		eventCounter++
	}

	assertAnyEvents := func(count int) {
		GinkgoHelper()
		Expect(len(trackedHandler.Events)).To(BeNumerically(">", eventCounter+count))
		eventCounter += count
	}

	assertNoMoreEvents := func() {
		GinkgoHelper()
		Expect(trackedHandler.Events).To(HaveLen(eventCounter))
	}

	JustBeforeEach(func() {
		file, err := os.Open(filepath.Join("testdata", testFile))
		Expect(err).ToNot(HaveOccurred())
		defer file.Close()

		scanner := obj.NewScanner()
		scanErr = scanner.Scan(file, handler)
	})

	BeforeEach(func() {
		trackedHandler = new(testutil.EventHandlerTracker)
		eventCounter = 0

		handler = trackedHandler.Handle
	})

	Describe("basic OBJ file", func() {
		BeforeEach(func() {
			testFile = "valid_basic.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned elements in order", func() {
			assertEvent(common.CommentEvent{
				Comment: "This is the beginning of this OBJ file.",
			})
			assertEvent(obj.MaterialLibraryEvent{
				FilePath: "valid_basic.mtl",
			})
			assertEvent(obj.VertexEvent{
				X: -1.0, Y: 1.0, Z: -1.0, W: 1.0,
			})
			assertEvent(obj.VertexEvent{
				X: -1.0, Y: -1.0, Z: 1.0, W: 1.0,
			})
			assertEvent(obj.VertexEvent{
				X: 1.0, Y: -1.0, Z: -1.0, W: 1.0,
			})
			assertEvent(obj.VertexEvent{
				X: 1.0, Y: 1.0, Z: 1.0, W: 1.0,
			})
			assertEvent(obj.TexCoordEvent{
				U: 0.0, V: 0.0, W: 0.0,
			})
			assertEvent(obj.TexCoordEvent{
				U: 1.0, V: 1.0, W: 0.0,
			})
			assertEvent(obj.TexCoordEvent{
				U: 1.0, V: 0.0, W: 0.0,
			})
			assertEvent(obj.TexCoordEvent{
				U: 0.0, V: 1.0, W: 0.0,
			})
			assertEvent(obj.NormalEvent{
				X: 0.0, Y: 1.0, Z: 0.0,
			})
			assertEvent(obj.NormalEvent{
				X: 1.0, Y: 0.0, Z: 0.0,
			})
			assertEvent(obj.NormalEvent{
				X: 0.0, Y: 0.0, Z: 1.0,
			})
			assertEvent(obj.ObjectEvent{
				ObjectName: "MyObject",
			})
			assertEvent(obj.MaterialReferenceEvent{
				MaterialName: "BlueMaterial",
			})
			assertEvent(obj.FaceStartEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 1,
			})
			assertEvent(obj.TexCoordReferenceEvent{
				TexCoordIndex: 4,
			})
			assertEvent(obj.NormalReferenceEvent{
				NormalIndex: 1,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 2,
			})
			assertEvent(obj.TexCoordReferenceEvent{
				TexCoordIndex: 1,
			})
			assertEvent(obj.NormalReferenceEvent{
				NormalIndex: 1,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 3,
			})
			assertEvent(obj.TexCoordReferenceEvent{
				TexCoordIndex: 3,
			})
			assertEvent(obj.NormalReferenceEvent{
				NormalIndex: 2,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.FaceEndEvent{})
			assertNoMoreEvents()
		})
	})

	When("a file with all kinds of comments is scanned", func() {
		BeforeEach(func() {
			testFile = "valid_comments.obj"
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

	When("a file with all kinds of vertices is scanned", func() {
		BeforeEach(func() {
			testFile = "valid_vertices.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the vertices", func() {
			assertEvent(obj.VertexEvent{
				X: 1.0, Y: 1.0, Z: -1.0, W: 1.0,
			})
			assertEvent(obj.VertexEvent{
				X: -1.0, Y: -1.0, Z: 1.0, W: 0.5,
			})
			assertNoMoreEvents()
		})
	})
	When("a file with all kinds of texture coordinates is scanned", func() {
		BeforeEach(func() {
			testFile = "valid_texcoords.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the texture coordinates", func() {
			assertEvent(obj.TexCoordEvent{
				U: 1.6, V: 0.0, W: 0.0,
			})
			assertEvent(obj.TexCoordEvent{
				U: 0.0, V: -0.5, W: 0.0,
			})
			assertEvent(obj.TexCoordEvent{
				U: -0.2, V: 1.4, W: 3.0,
			})
			assertNoMoreEvents()
		})
	})

	When("a file with all kinds of normals is scanned", func() {
		BeforeEach(func() {
			testFile = "valid_normals.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the normals", func() {
			assertEvent(obj.NormalEvent{
				X: 1.0, Y: 0.0, Z: 0.0,
			})
			assertEvent(obj.NormalEvent{
				X: -1.0, Y: -1.0, Z: 1.0,
			})
			assertNoMoreEvents()
		})
	})

	When("a file with all kinds of objects is scanned", func() {
		BeforeEach(func() {
			testFile = "valid_objects.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the objects", func() {
			assertEvent(obj.ObjectEvent{
				ObjectName: "FirstObject",
			})
			assertEvent(obj.ObjectEvent{
				ObjectName: "LastObject",
			})
			assertNoMoreEvents()
		})
	})
	When("a file with all kinds of coord references is scanned", func() {
		BeforeEach(func() {
			testFile = "valid_faces.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned them correctly", func() {
			// First face
			assertEvent(obj.FaceStartEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 1,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 2,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 3,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.FaceEndEvent{})

			// Second face
			assertEvent(obj.FaceStartEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 2,
			})
			assertEvent(obj.TexCoordReferenceEvent{
				TexCoordIndex: 1,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 3,
			})
			assertEvent(obj.TexCoordReferenceEvent{
				TexCoordIndex: 2,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 4,
			})
			assertEvent(obj.TexCoordReferenceEvent{
				TexCoordIndex: 3,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.FaceEndEvent{})

			// Third face
			assertEvent(obj.FaceStartEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 4,
			})
			assertEvent(obj.NormalReferenceEvent{
				NormalIndex: 5,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 6,
			})
			assertEvent(obj.NormalReferenceEvent{
				NormalIndex: 7,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 8,
			})
			assertEvent(obj.NormalReferenceEvent{
				NormalIndex: 9,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.FaceEndEvent{})

			// Fourth face
			assertEvent(obj.FaceStartEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 1,
			})
			assertEvent(obj.TexCoordReferenceEvent{
				TexCoordIndex: 2,
			})
			assertEvent(obj.NormalReferenceEvent{
				NormalIndex: 3,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 2,
			})
			assertEvent(obj.TexCoordReferenceEvent{
				TexCoordIndex: 3,
			})
			assertEvent(obj.NormalReferenceEvent{
				NormalIndex: 4,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.ReferenceSetStartEvent{})
			assertEvent(obj.VertexReferenceEvent{
				VertexIndex: 3,
			})
			assertEvent(obj.TexCoordReferenceEvent{
				TexCoordIndex: 4,
			})
			assertEvent(obj.NormalReferenceEvent{
				NormalIndex: 5,
			})
			assertEvent(obj.ReferenceSetEndEvent{})
			assertEvent(obj.FaceEndEvent{})

			assertNoMoreEvents()
		})
	})

	When("a file with all kinds of material libraries is scanned", func() {
		BeforeEach(func() {
			testFile = "valid_material_libraries.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the material libraries", func() {
			assertEvent(obj.MaterialLibraryEvent{
				FilePath: "valid.mtl",
			})
			assertEvent(obj.MaterialLibraryEvent{
				FilePath: "extension01.mtl",
			})
			assertEvent(obj.MaterialLibraryEvent{
				FilePath: "extension02.mtl",
			})
			assertNoMoreEvents()
		})
	})

	When("a file with all kinds of material references is scanned", func() {
		BeforeEach(func() {
			testFile = "valid_material_references.obj"
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the material references", func() {
			assertEvent(obj.MaterialReferenceEvent{
				MaterialName: "",
			})
			assertEvent(obj.MaterialReferenceEvent{
				MaterialName: "MyMaterial",
			})
			assertNoMoreEvents()
		})
	})

	When("a file with insufficient vertex data is scanned", func() {
		BeforeEach(func() {
			testFile = "error_insufficient_vertex_data.obj"
		})

		itShouldHaveReturnedAnError()
	})

	When("a file with insufficient texture coordinate data is scanned", func() {
		BeforeEach(func() {
			testFile = "error_insufficient_texcoord_data.obj"
		})

		itShouldHaveReturnedAnError()
	})

	When("a file with insufficient normal data is scanned", func() {
		BeforeEach(func() {
			testFile = "error_insufficient_normal_data.obj"
		})

		itShouldHaveReturnedAnError()
	})

	When("a file with an unnamed object is scanned", func() {
		BeforeEach(func() {
			testFile = "error_empty_object_name.obj"
		})

		itShouldHaveReturnedAnError()
	})

	When("a file with corrupt vertex is scanned", func() {
		BeforeEach(func() {
			testFile = "error_corrupt_vertex.obj"
		})

		itShouldHaveReturnedAnError()
	})

	When("a file with corrupt texture coordinate is scanned", func() {
		BeforeEach(func() {
			testFile = "error_corrupt_texcoord.obj"
		})

		itShouldHaveReturnedAnError()
	})

	When("a file with corrupt normal is scanned", func() {
		BeforeEach(func() {
			testFile = "error_corrupt_normal.obj"
		})

		itShouldHaveReturnedAnError()
	})

	When("a file with corrupt vertex reference is scanned", func() {
		BeforeEach(func() {
			testFile = "error_corrupt_vertex_reference.obj"
		})

		itShouldHaveReturnedAnError()
	})

	When("a file with corrupt texture coordinate reference is scanned", func() {
		BeforeEach(func() {
			testFile = "error_corrupt_texcoord_reference.obj"
		})

		itShouldHaveReturnedAnError()
	})

	When("a file with corrupt normal reference is scanned", func() {
		BeforeEach(func() {
			testFile = "error_corrupt_normal_reference.obj"
		})

		itShouldHaveReturnedAnError()
	})
})
