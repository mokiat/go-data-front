package mtl

import (
	"errors"
	"io"

	"github.com/momchil-atanasov/go-data-front/common"
)

// MaterialEvent indicates that a material declaration (`newmtl`) has
// been scanned.
type MaterialEvent struct {

	// MaterialName holds the name of the material in the declaration
	MaterialName string
}

// RGBColorEvent indicates that some type of RGB color declaration
// has been scanned. You will likely receive a subtype of this
// structure so you will likely need to do a type-switch.
type RGBColorEvent struct {

	// Specifies the amount of Red this color has. Usually this is in the
	// range 0.0 to 1.0.
	R float64

	// Specifies the amount of Green this color has. Usually this is in the
	// range 0.0 to 1.0.
	G float64

	// Specifies the amount of Blue this color has. Usually this is in the
	// range 0.0 to 1.0.
	B float64
}

// RGBAmbientColorEvent indicates that an ambient color declaration (`Ka`)
// has been scanned.
type RGBAmbientColorEvent RGBColorEvent

// RGBDiffuseColorEvent indicates that a diffuse color declaration (`Kd`)
// has been scanned.
type RGBDiffuseColorEvent RGBColorEvent

// RGBSpecularColorEvent indicates that a specular color declaration (`Ks`)
// has been scanned.
type RGBSpecularColorEvent RGBColorEvent

// RGBTransmissionFilterEvent indicates that a specular color declaration (`Tf`)
// has been scanned.
type RGBTransmissionFilterEvent RGBColorEvent

// DissolveEvent indicates that a dissolve declaration (`d`) has been
// scanned.
type DissolveEvent struct {

	// Amount indicates the amount of dissolve, where 1.0 indicates fully
	// opaque objects and 0.0 fully transparent.
	Amount float64
}

// SpecularExponentEvent indicates that a specular exponent declaration (`Ns`)
// has been scanned.
type SpecularExponentEvent struct {

	// Amount specifies the specular exponent amount, which defines the focus of
	// the specular highlight. The value ranges between 0.0 and 1000.0, where the
	// former results in a wide highlight and the latter in a tight one.
	Amount float64
}

// TextureEvent indicates that a texture declaration has been scanned.
// You will likely receive a subtype of this structure so you will likely
// need to type-switch on this.
type TextureEvent struct {

	// TexturePath specifies the location of the texture on the filesystem.
	TexturePath string
}

// AmbientTextureEvent indicates that an ambient texture declaration (`map_Ka`)
// has been scanned.
type AmbientTextureEvent TextureEvent

// DiffuseTextureEvent indicates that a diffuse texture declaration (`map_Kd`)
// has been scanned.
type DiffuseTextureEvent TextureEvent

// SpecularTextureEvent indicates that a specular texture declaration (`map_Ks`)
// has been scanned.
type SpecularTextureEvent TextureEvent

// SpecularExponentTextureEvent indicates that a specular exponent texture
// declaration (`map_Ns`) has been scanned.
type SpecularExponentTextureEvent TextureEvent

// DissolveTextureEvent indicates that a dissolve texture
// declaration (`map_d`) has been scanned.
type DissolveTextureEvent TextureEvent

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
	// XXX: NOTE: There is no need to have RColorEvents as the Spec states that if only
	// r is specified, g and b equal that. So you can just reuse RGB events.
	switch {
	case line.HasCommandName("newmtl"):
		return s.processMaterial(line, handler)
	case line.HasCommandName("Ka"):
		return s.processAmbientColor(line, handler)
	case line.HasCommandName("Kd"):
		return s.processDiffuseColor(line, handler)
	case line.HasCommandName("Ks"):
		return s.processSpecularColor(line, handler)
	case line.HasCommandName("Tf"):
		return s.processTransmissionFilter(line, handler)
	case line.HasCommandName("d"):
		return s.processDissolve(line, handler)
	case line.HasCommandName("Ns"):
		return s.processSpecularExponent(line, handler)
	case line.HasCommandName("map_Ka"):
		return s.processAmbientTexture(line, handler)
	case line.HasCommandName("map_Kd"):
		return s.processDiffuseTexture(line, handler)
	case line.HasCommandName("map_Ks"):
		return s.processSpecularTexture(line, handler)
	case line.HasCommandName("map_Ns"):
		return s.processSpecularExponentTexture(line, handler)
	case line.HasCommandName("map_d"):
		return s.processDissolveTexture(line, handler)
	default:
		return nil
	}
}

func (s *scanner) processMaterial(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 1 {
		return errors.New("Material declaration lacks name!")
	}
	name := line.StringParam(0)
	event := MaterialEvent{
		MaterialName: name,
	}
	return handler(event)
}

func (s *scanner) processAmbientColor(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 1 {
		return errors.New("Ambient color declaration lacks parameters!")
	}

	// TODO: Handle other scenarios
	event, err := s.getRGBColorEvent(line)
	if err != nil {
		return err
	}
	return handler(RGBAmbientColorEvent(event))
}

func (s *scanner) processDiffuseColor(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 1 {
		return errors.New("Diffuse color declaration lacks parameters!")
	}

	// TODO: Handle other scenarios
	event, err := s.getRGBColorEvent(line)
	if err != nil {
		return err
	}
	return handler(RGBDiffuseColorEvent(event))
}

func (s *scanner) processSpecularColor(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 1 {
		return errors.New("Specular color declaration lacks parameters!")
	}

	// TODO: Handle other scenarios
	event, err := s.getRGBColorEvent(line)
	if err != nil {
		return err
	}
	return handler(RGBSpecularColorEvent(event))
}

func (s *scanner) processTransmissionFilter(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 3 {
		return errors.New("Transmission filter declaration lacks parameters!")
	}

	event, err := s.getRGBColorEvent(line)
	if err != nil {
		return err
	}
	return handler(RGBTransmissionFilterEvent(event))
}

func (s *scanner) getRGBColorEvent(line common.Line) (RGBColorEvent, error) {
	var err error
	var event RGBColorEvent

	event.R, err = line.FloatParam(0)
	if err != nil {
		return RGBColorEvent{}, err
	}

	if line.ParamCount() >= 3 {
		// TODO: Test this!
		event.G, err = line.FloatParam(1)
		if err != nil {
			return RGBColorEvent{}, err
		}

		event.B, err = line.FloatParam(2)
		if err != nil {
			return RGBColorEvent{}, err
		}
	} else {
		// TODO: Test this scenario
		event.G = event.R
		event.B = event.R
	}

	return event, nil
}

func (s *scanner) processDissolve(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 1 {
		return errors.New("Dissolve declaration lacks value parameter!")
	}
	amount, err := line.FloatParam(0)
	if err != nil {
		return err
	}
	event := DissolveEvent{
		Amount: amount,
	}
	return handler(event)
}

func (s *scanner) processSpecularExponent(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 1 {
		return errors.New("Specular exponent declaration lacks value parameter!")
	}
	amount, err := line.FloatParam(0)
	if err != nil {
		return err
	}
	event := SpecularExponentEvent{
		Amount: amount,
	}
	return handler(event)
}

func (s *scanner) processAmbientTexture(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 1 {
		return errors.New("Ambient texture declaration lacks filename parameter!")
	}
	event, err := s.getTextureEvent(line)
	if err != nil {
		return err
	}
	return handler(AmbientTextureEvent(event))
}

func (s *scanner) processDiffuseTexture(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 1 {
		return errors.New("Diffuse texture declaration lacks filename parameter!")
	}
	event, err := s.getTextureEvent(line)
	if err != nil {
		return err
	}
	return handler(DiffuseTextureEvent(event))
}

func (s *scanner) processSpecularTexture(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 1 {
		return errors.New("Specular texture declaration lacks filename parameter!")
	}
	event, err := s.getTextureEvent(line)
	if err != nil {
		return err
	}
	return handler(SpecularTextureEvent(event))
}

func (s *scanner) processSpecularExponentTexture(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 1 {
		return errors.New("Specular exponent texture declaration lacks filename parameter!")
	}
	event, err := s.getTextureEvent(line)
	if err != nil {
		return err
	}
	return handler(SpecularExponentTextureEvent(event))
}

func (s *scanner) processDissolveTexture(line common.Line, handler common.EventHandler) error {
	if line.ParamCount() < 1 {
		return errors.New("Dissolve texture declaration lacks filename parameter!")
	}
	event, err := s.getTextureEvent(line)
	if err != nil {
		return err
	}
	return handler(DissolveTextureEvent(event))
}

func (s *scanner) getTextureEvent(line common.Line) (TextureEvent, error) {
	path := line.StringParam(0)
	event := TextureEvent{
		TexturePath: path,
	}
	return event, nil
}