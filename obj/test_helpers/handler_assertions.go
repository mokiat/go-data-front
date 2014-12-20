package test_helpers

import (
	objfakes "github.com/momchil-atanasov/go-data-front/obj/fakes"
	. "github.com/onsi/gomega"
)

type HandlerFixture struct {
	handler                    *objfakes.FakeScannerHandler
	commentCounter             int
	vertexStartCounter         int
	vertexEndCounter           int
	vertexXCounter             int
	vertexYCounter             int
	vertexZCounter             int
	texCoordStartCounter       int
	texCoordEndCounter         int
	texCoordUCounter           int
	texCoordVCounter           int
	texCoordWCounter           int
	normalCounter              int
	objectCounter              int
	faceStartCounter           int
	faceEndCounter             int
	coordReferenceStartCounter int
	coordReferenceEndCounter   int
	vertexIndexCounter         int
	texCoordIndexCounter       int
	normalIndexCounter         int
	materialLibraryCounter     int
	materialReferenceCounter   int
}

func NewHandlerFixture() *HandlerFixture {
	return &HandlerFixture{
		handler: new(objfakes.FakeScannerHandler),
	}
}

func (f *HandlerFixture) Handler() *objfakes.FakeScannerHandler {
	return f.handler
}

func (f *HandlerFixture) AssertCommentCall(expectedComment string) {
	Ω(f.handler.OnCommentCallCount()).Should(BeNumerically(">", f.commentCounter))
	Ω(f.handler.OnCommentArgsForCall(f.commentCounter)).Should(Equal(expectedComment))
	f.commentCounter++
}

func (f *HandlerFixture) AssertNoMoreCommentCalls() {
	Ω(f.handler.OnCommentCallCount()).Should(Equal(f.commentCounter))
}

func (f *HandlerFixture) AssertVertexStart() {
	Ω(f.handler.OnVertexStartCallCount()).Should(BeNumerically(">", f.vertexStartCounter))
	f.vertexStartCounter++
}

func (f *HandlerFixture) AssertVertexEnd() {
	Ω(f.handler.OnVertexEndCallCount()).Should(BeNumerically(">", f.vertexEndCounter))
	f.vertexEndCounter++
}

func (f *HandlerFixture) AssertVertexX(expectedX float32) {
	Ω(f.handler.OnVertexXCallCount()).Should(BeNumerically(">", f.vertexXCounter))
	Ω(f.handler.OnVertexXArgsForCall(f.vertexXCounter)).Should(BeNumerically("~", expectedX, 0.0001))
	f.vertexXCounter++
}

func (f *HandlerFixture) AssertVertexY(expectedY float32) {
	Ω(f.handler.OnVertexYCallCount()).Should(BeNumerically(">", f.vertexYCounter))
	Ω(f.handler.OnVertexYArgsForCall(f.vertexYCounter)).Should(BeNumerically("~", expectedY, 0.0001))
	f.vertexYCounter++
}

func (f *HandlerFixture) AssertVertexZ(expectedZ float32) {
	Ω(f.handler.OnVertexZCallCount()).Should(BeNumerically(">", f.vertexZCounter))
	Ω(f.handler.OnVertexZArgsForCall(f.vertexZCounter)).Should(BeNumerically("~", expectedZ, 0.0001))
	f.vertexZCounter++
}

func (f *HandlerFixture) AssertVertexXYZ(expectedX, expectedY, expectedZ float32) {
	f.AssertVertexStart()
	f.AssertVertexX(expectedX)
	f.AssertVertexY(expectedY)
	f.AssertVertexZ(expectedZ)
	f.AssertVertexEnd()
}

func (f *HandlerFixture) AssertTexCoordStart() {
	Ω(f.handler.OnTexCoordStartCallCount()).Should(BeNumerically(">", f.texCoordStartCounter))
	f.texCoordStartCounter++
}

func (f *HandlerFixture) AssertTexCoordU(expectedU float32) {
	Ω(f.handler.OnTexCoordUCallCount()).Should(BeNumerically(">", f.texCoordUCounter))
	Ω(f.handler.OnTexCoordUArgsForCall(f.texCoordUCounter)).Should(BeNumerically("~", expectedU, 0.0001))
	f.texCoordUCounter++
}

func (f *HandlerFixture) AssertTexCoordV(expectedV float32) {
	Ω(f.handler.OnTexCoordVCallCount()).Should(BeNumerically(">", f.texCoordVCounter))
	Ω(f.handler.OnTexCoordVArgsForCall(f.texCoordVCounter)).Should(BeNumerically("~", expectedV, 0.0001))
	f.texCoordVCounter++
}

func (f *HandlerFixture) AssertTexCoordEnd() {
	Ω(f.handler.OnTexCoordEndCallCount()).Should(BeNumerically(">", f.texCoordEndCounter))
	f.texCoordEndCounter++
}

func (f *HandlerFixture) AssertTexCoordUV(expectedU, expectedV float32) {
	f.AssertTexCoordStart()
	f.AssertTexCoordU(expectedU)
	f.AssertTexCoordV(expectedV)
	f.AssertTexCoordEnd()
}

func (f *HandlerFixture) AssertNormal(expectedX, expectedY, expectedZ float32) {
	Ω(f.handler.OnNormalCallCount()).Should(BeNumerically(">", f.normalCounter))
	argX, argY, argZ := f.handler.OnNormalArgsForCall(f.normalCounter)
	Ω(argX).Should(BeNumerically("~", expectedX, 0.0001))
	Ω(argY).Should(BeNumerically("~", expectedY, 0.0001))
	Ω(argZ).Should(BeNumerically("~", expectedZ, 0.0001))
	f.normalCounter++
}

func (f *HandlerFixture) AssertObject(expectedName string) {
	Ω(f.handler.OnObjectCallCount()).Should(BeNumerically(">", f.objectCounter))
	argName := f.handler.OnObjectArgsForCall(f.objectCounter)
	Ω(argName).Should(Equal(expectedName))
	f.objectCounter++
}

func (f *HandlerFixture) AssertFaceStart() {
	Ω(f.handler.OnFaceStartCallCount()).Should(BeNumerically(">", f.faceStartCounter))
	f.faceStartCounter++
}

func (f *HandlerFixture) AssertFaceEnd() {
	Ω(f.handler.OnFaceEndCallCount()).Should(BeNumerically(">", f.faceEndCounter))
	f.faceEndCounter++
}

func (f *HandlerFixture) AssertFaceCallCount(expectedCount int) {
	Ω(f.handler.OnFaceStartCallCount()).Should(Equal(expectedCount))
	Ω(f.handler.OnFaceEndCallCount()).Should(Equal(expectedCount))
}

func (f *HandlerFixture) AssertCoordReferenceStart() {
	Ω(f.handler.OnCoordReferenceStartCallCount()).Should(BeNumerically(">", f.coordReferenceStartCounter))
	f.coordReferenceStartCounter++
}

func (f *HandlerFixture) AssertCoordReferenceEnd() {
	Ω(f.handler.OnCoordReferenceEndCallCount()).Should(BeNumerically(">", f.coordReferenceEndCounter))
	f.coordReferenceEndCounter++
}

func (f *HandlerFixture) AssertCoordReferenceCallCount(expectedCount int) {
	Ω(f.handler.OnCoordReferenceStartCallCount()).Should(Equal(expectedCount))
	Ω(f.handler.OnCoordReferenceEndCallCount()).Should(Equal(expectedCount))
}

func (f *HandlerFixture) AssertVertexIndex(expectedIndex int) {
	Ω(f.handler.OnVertexIndexCallCount()).Should(BeNumerically(">", f.vertexIndexCounter))
	argIndex := f.handler.OnVertexIndexArgsForCall(f.vertexIndexCounter)
	Ω(argIndex).Should(Equal(expectedIndex))
	f.vertexIndexCounter++
}

func (f *HandlerFixture) AssertTexCoordIndex(expectedIndex int) {
	Ω(f.handler.OnTexCoordIndexCallCount()).Should(BeNumerically(">", f.texCoordIndexCounter))
	argIndex := f.handler.OnTexCoordIndexArgsForCall(f.texCoordIndexCounter)
	Ω(argIndex).Should(Equal(expectedIndex))
	f.texCoordIndexCounter++
}

func (f *HandlerFixture) AssertNormalIndex(expectedIndex int) {
	Ω(f.handler.OnNormalIndexCallCount()).Should(BeNumerically(">", f.normalIndexCounter))
	argIndex := f.handler.OnNormalIndexArgsForCall(f.normalIndexCounter)
	Ω(argIndex).Should(Equal(expectedIndex))
	f.normalIndexCounter++
}

func (f *HandlerFixture) AssertMaterialLibrary(expectedPath string) {
	Ω(f.handler.OnMaterialLibraryCallCount()).Should(BeNumerically(">", f.materialLibraryCounter))
	pathArg := f.handler.OnMaterialLibraryArgsForCall(f.materialLibraryCounter)
	Ω(pathArg).Should(Equal(expectedPath))
	f.materialLibraryCounter++
}

func (f *HandlerFixture) AssertMaterialReference(expectedName string) {
	Ω(f.handler.OnMaterialReferenceCallCount()).Should(BeNumerically(">", f.materialReferenceCounter))
	nameArg := f.handler.OnMaterialReferenceArgsForCall(f.materialReferenceCounter)
	Ω(nameArg).Should(Equal(expectedName))
	f.materialReferenceCounter++
}
