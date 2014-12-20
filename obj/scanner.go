package obj

import (
	"bufio"
	"io"
)

const commandComment = "v"
const commandTexCoord = "vt"
const commandNormal = "vn"
const commandObject = "o"
const commandFace = "f"
const commandMaterialRef = "usemtl"
const commandMaterialLib = "mtllib"

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

type Scanner interface {
	Scan(io.Reader) error
}

type scanner struct {
	handler ScannerHandler
}

func NewScanner(handler ScannerHandler) Scanner {
	return &scanner{
		handler: handler,
	}
}

func (s *scanner) Scan(reader io.Reader) error {
	line := NewScanLine()

	scanner := bufio.NewScanner(reader)
	for {
		line.Parse(scanner)
		// TODO: Handle Parse error
		switch {
		case line.IsBlank():
			break
		case line.IsComment():
			s.processComment(line)
			break
		case line.IsCommand(commandComment):
			s.processVertex(line)
			break
		case line.IsCommand(commandTexCoord):
			s.processTexCoord(line)
			break
		case line.IsCommand(commandNormal):
			s.processNormal(line)
			break
		case line.IsCommand(commandObject):
			s.processObject(line)
			break
		case line.IsCommand(commandFace):
			s.processFace(line)
			break
		case line.IsCommand(commandMaterialLib):
			s.processMaterialLibrary(line)
			break
		case line.IsCommand(commandMaterialRef):
			s.processMaterialReference(line)
			break
		}
		if line.IsAtEOF() {
			return nil
		}
	}
}

func (s *scanner) processComment(line ScanLine) error {
	s.handler.OnComment(line.GetComment())
	return nil
}

func (s *scanner) processVertex(line ScanLine) error {
	s.handler.OnVertexStart()
	x := line.FloatParam(0)
	s.handler.OnVertexX(x)
	y := line.FloatParam(1)
	s.handler.OnVertexY(y)
	z := line.FloatParam(2)
	s.handler.OnVertexZ(z)
	s.handler.OnVertexEnd()
	return nil
}

func (s *scanner) processTexCoord(line ScanLine) error {
	u := line.FloatParam(0)
	v := line.FloatParam(1)
	s.handler.OnTexCoordStart()
	s.handler.OnTexCoordU(u)
	s.handler.OnTexCoordV(v)
	s.handler.OnTexCoordEnd()
	return nil
}

func (s *scanner) processNormal(line ScanLine) error {
	x := line.FloatParam(0)
	y := line.FloatParam(1)
	z := line.FloatParam(2)
	s.handler.OnNormal(x, y, z)
	return nil
}

func (s *scanner) processObject(line ScanLine) error {
	name := line.StringParam(0)
	s.handler.OnObject(name)
	return nil
}

func (s *scanner) processFace(line ScanLine) error {
	s.handler.OnFaceStart()
	for i := 0; i < line.ParamCount(); i++ {
		s.handler.OnCoordReferenceStart()
		coordReference := line.CoordReferenceParam(i)
		s.handler.OnVertexIndex(coordReference.VertexIndex)
		if coordReference.HasTexCoordIndex {
			s.handler.OnTexCoordIndex(coordReference.TexCoordIndex)
		}
		if coordReference.HasNormalIndex {
			s.handler.OnNormalIndex(coordReference.NormalIndex)
		}
		s.handler.OnCoordReferenceEnd()
	}
	s.handler.OnFaceEnd()
	return nil
}

func (s *scanner) processMaterialLibrary(line ScanLine) error {
	for i := 0; i < line.ParamCount(); i++ {
		path := line.StringParam(i)
		s.handler.OnMaterialLibrary(path)
	}
	return nil
}

func (s *scanner) processMaterialReference(line ScanLine) error {
	if line.ParamCount() > 0 {
		name := line.StringParam(0)
		s.handler.OnMaterialReference(name)
	} else {
		s.handler.OnMaterialReference("")
	}
	return nil
}
