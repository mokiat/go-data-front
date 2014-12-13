package obj

import (
	"bufio"
	"io"
)

type ScannerHandler interface {
	OnComment(comment string) error
	OnMaterialLibrary(path string) error
	OnMaterialReference(name string) error
	OnVertexStart() error
	OnVertexX(x float32) error
	OnVertexY(y float32) error
	OnVertexZ(z float32) error
	OnVertexW(w float32) error
	OnVertexEnd() error
	OnTexCoordStart() error
	OnTexCoordU(u float32) error
	OnTexCoordV(v float32) error
	OnTexCoordW(w float32) error
	OnTexCoordEnd() error
	OnNormal(x, y, z float32) error
	OnObject(name string) error
	OnFaceStart() error
	OnCoordReferenceStart() error
	OnVertexIndex(index int) error
	OnTexCoordIndex(index int) error
	OnNormalIndex(index int) error
	OnCoordReferenceEnd() error
	OnFaceEnd() error
}

const commandComment = "v"
const commandTexCoord = "vt"
const commandNormal = "vn"
const commandObject = "o"
const commandFace = "f"
const commandMaterialRef = "usemtl"
const commandMaterialLib = "mtllib"

func Scan(reader io.Reader, handler ScannerHandler) error {
	line := NewScanLine()

	scanner := bufio.NewScanner(reader)
	for {
		line.Parse(scanner)
		// TODO: Handle Parse error
		switch {
		case line.IsBlank():
			break
		case line.IsComment():
			processComment(line, handler)
			break
		case line.IsCommand(commandComment):
			processVertex(line, handler)
		case line.IsCommand(commandTexCoord):
			processTexCoord(line, handler)
		case line.IsCommand(commandNormal):
			processNormal(line, handler)
		case line.IsCommand(commandObject):
			processObject(line, handler)
		}
		if line.IsAtEOF() {
			return nil
		}
	}
}

func processComment(line ScanLine, handler ScannerHandler) error {
	handler.OnComment(line.GetComment())
	return nil
}

func processVertex(line ScanLine, handler ScannerHandler) error {
	handler.OnVertexStart()
	x := line.FloatParam(0)
	handler.OnVertexX(x)
	y := line.FloatParam(1)
	handler.OnVertexY(y)
	z := line.FloatParam(2)
	handler.OnVertexZ(z)
	handler.OnVertexEnd()
	return nil
}

func processTexCoord(line ScanLine, handler ScannerHandler) error {
	u := line.FloatParam(0)
	v := line.FloatParam(1)
	handler.OnTexCoordStart()
	handler.OnTexCoordU(u)
	handler.OnTexCoordV(v)
	handler.OnTexCoordEnd()
	return nil
}

func processNormal(line ScanLine, handler ScannerHandler) error {
	x := line.FloatParam(0)
	y := line.FloatParam(1)
	z := line.FloatParam(2)
	handler.OnNormal(x, y, z)
	return nil
}

func processObject(line ScanLine, handler ScannerHandler) error {
	name := line.StringParam(0)
	handler.OnObject(name)
	return nil
}
