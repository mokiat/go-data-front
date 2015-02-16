package mtl

import (
	"io"

	"github.com/momchil-atanasov/go-data-front/common"
)

// MaterialEvent indicates that a material declaration (`newmtl`) has
// been scanned.
type MaterialEvent struct {

	// MaterialName holds the name of the material in the declaration
	MaterialName string
}

// NewScanner creates a new Scanner object that can scan through
// Wavefront MTL resources.
func NewScanner() common.Scanner {
	return &scanner{}
}

type scanner struct {
}

func (s *scanner) Scan(reader io.Reader, handler common.EventHandler) error {
	lineScanner := common.NewLineScanner(reader)

	var err error
	for lineScanner.Scan() {
		line := lineScanner.Line()
		switch {
		case line.IsBlank():
			break
		case line.IsComment():
			err = s.processComment(line, handler)
			break
		case line.IsCommand():
			err = s.processCommand(line, handler)
			break
		}
		if err != nil {
			return err
		}
	}

	if lineScanner.Err() != nil {
		return lineScanner.Err()
	}
	return nil
}

func (s *scanner) processComment(line common.Line, handler common.EventHandler) error {
	event := common.CommentEvent{
		Comment: line.Comment(),
	}
	return handler(event)
}

func (s *scanner) processCommand(line common.Line, handler common.EventHandler) error {
	switch {
	case line.HasCommandName("newmtl"):
		return s.processMaterial(line, handler)
	default:
		return nil
	}
}

func (s *scanner) processMaterial(line common.Line, handler common.EventHandler) error {
	// TODO: Test missing material name
	name := line.StringParam(0)
	event := MaterialEvent{
		MaterialName: name,
	}
	return handler(event)
}
