package obj

import (
	"bufio"
	"errors"
	"io"
)

const commandComment = "v"
const commandTexCoord = "vt"
const commandNormal = "vn"
const commandObject = "o"
const commandFace = "f"
const commandMaterialRef = "usemtl"
const commandMaterialLib = "mtllib"

type scanner struct {
	handler ScannerHandler
}

func (s *scanner) Scan(reader io.Reader) error {
	line := NewScanLine()

	scanner := bufio.NewScanner(reader)
	for {
		var err error
		err = line.Parse(scanner)
		if err != nil {
			return err
		}
		switch {
		case line.IsBlank():
			break
		case line.IsComment():
			err = s.processComment(line)
			break
		case line.IsCommand(commandComment):
			err = s.processVertex(line)
			break
		case line.IsCommand(commandTexCoord):
			err = s.processTexCoord(line)
			break
		case line.IsCommand(commandNormal):
			err = s.processNormal(line)
			break
		case line.IsCommand(commandObject):
			err = s.processObject(line)
			break
		case line.IsCommand(commandFace):
			err = s.processFace(line)
			break
		case line.IsCommand(commandMaterialLib):
			err = s.processMaterialLibrary(line)
			break
		case line.IsCommand(commandMaterialRef):
			err = s.processMaterialReference(line)
			break
		}
		if err != nil {
			return err
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
	if line.ParamCount() < 3 {
		return errors.New("Insufficient vertex data.")
	}
	s.handler.OnVertexStart()

	x, err := line.FloatParam(0)
	if err != nil {
		return err
	}
	s.handler.OnVertexX(x)

	y, err := line.FloatParam(1)
	if err != nil {
		return err
	}
	s.handler.OnVertexY(y)

	z, err := line.FloatParam(2)
	if err != nil {
		return err
	}
	s.handler.OnVertexZ(z)

	if line.ParamCount() >= 4 {
		w, err := line.FloatParam(3)
		if err != nil {
			return err
		}
		s.handler.OnVertexW(w)
	}

	s.handler.OnVertexEnd()
	return nil
}

func (s *scanner) processTexCoord(line ScanLine) error {
	if line.ParamCount() == 0 {
		return errors.New("Insufficient texture coordinate data.")
	}
	s.handler.OnTexCoordStart()

	u, err := line.FloatParam(0)
	if err != nil {
		return err
	}
	s.handler.OnTexCoordU(u)

	if line.ParamCount() >= 2 {
		v, err := line.FloatParam(1)
		if err != nil {
			return err
		}
		s.handler.OnTexCoordV(v)
	}

	if line.ParamCount() >= 3 {
		w, err := line.FloatParam(2)
		if err != nil {
			return err
		}
		s.handler.OnTexCoordW(w)
	}

	s.handler.OnTexCoordEnd()
	return nil
}

func (s *scanner) processNormal(line ScanLine) error {
	if line.ParamCount() < 3 {
		return errors.New("Insufficient normal data.")
	}

	x, err := line.FloatParam(0)
	if err != nil {
		return err
	}

	y, err := line.FloatParam(1)
	if err != nil {
		return err
	}

	z, err := line.FloatParam(2)
	if err != nil {
		return err
	}

	s.handler.OnNormal(x, y, z)
	return nil
}

func (s *scanner) processObject(line ScanLine) error {
	if line.ParamCount() == 0 {
		return errors.New("No name specified for object.")
	}
	name := line.StringParam(0)
	s.handler.OnObject(name)
	return nil
}

func (s *scanner) processFace(line ScanLine) error {
	s.handler.OnFaceStart()
	for i := 0; i < line.ParamCount(); i++ {
		s.handler.OnCoordReferenceStart()

		coordReference, err := line.CoordReferenceParam(i)
		if err != nil {
			return err
		}

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
