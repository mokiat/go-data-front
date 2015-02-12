package mtl

import (
	"io"

	"github.com/momchil-atanasov/go-data-front/common"
)

// ScannerHandler interface acts as a callback that is used by the
// MTL Scanner implementation whenever an interesting MTL feature is scanned.
//
// Note: All methods allow for an error to be returned. Implementations can
// use this to stop any further scanning of the resource as well as to
// convey that the sequence of data is not as expected.
type ScannerHandler interface {

	// OnComment is called when a comment section has been parsed.
	OnComment(comment string) error

	// OnMaterial is called when a material declaration section has been parsed.
	OnMaterial(name string) error
}

//go:generate counterfeiter -o mtl_test_fake/fake_scanner_handler.go . ScannerHandler

type scanner struct {
	handler ScannerHandler
}

const commandMaterial = "newmtl"

// NewScanner creates a new Scanner object that can scan through MTL resources
// and trigger events whenever a feature is parsed.
func NewScanner(handler ScannerHandler) common.Scanner {
	return &scanner{
		handler: handler,
	}
}

func (s *scanner) Scan(reader io.Reader) error {
	lineScanner := common.NewLineScanner(reader)

	var err error
	for lineScanner.Scan() {
		line := lineScanner.Line()
		switch {
		case line.IsBlank():
			break
		case line.IsComment():
			err = s.processComment(line)
			break
		case line.IsCommand():
			err = s.processCommand(line)
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

func (s *scanner) processComment(line common.Line) error {
	return s.handler.OnComment(line.Comment())
}

func (s *scanner) processCommand(line common.Line) error {
	switch {
	case line.HasCommandName(commandMaterial):
		return s.processMaterial(line)
	default:
		return nil
	}
}

func (s *scanner) processMaterial(line common.Line) error {
	// TODO: Test missing material name
	name := line.StringParam(0)
	return s.handler.OnMaterial(name)
}
