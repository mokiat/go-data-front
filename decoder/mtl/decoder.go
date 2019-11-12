package mtl

import (
	"errors"
	"io"

	scanMTL "github.com/mokiat/go-data-front/scanner/mtl"
	"github.com/mokiat/go-data-front/common"
)

// DecodeLimits specifies restrictions on parsing.
//
// One will generally use this to limit the number of
// data that is parsed in order to prevent out of memory
// errors
type DecodeLimits struct {

	// MaxMaterialCount specifies the maximum number of
	// material declarations that can be parsed.
	MaxMaterialCount int
}

// DefaultLimits returns some default DecodeLimits.
// Users can take the result and modify specific parameters.
func DefaultLimits() DecodeLimits {
	return DecodeLimits{
		MaxMaterialCount: 512,
	}
}

// Decoder is an API that allows one to decode MTL
// Wavefront resources into an object model.
type Decoder interface {

	// Decode decodes the MTL Wavefront resource, specified
	// through the io.Reader, into a Library model.
	//
	// If decoding fails for some reason, an error is returned.
	Decode(io.Reader) (*Library, error)
}

// NewDecoder creates a new Decoder instance with the
// specified DecodeLimits.
func NewDecoder(limits DecodeLimits) Decoder {
	return &decoder{
		limits:  &limits,
		scanner: scanMTL.NewScanner(),
	}
}

type decoder struct {
	limits  *DecodeLimits
	scanner common.Scanner
}

func (d *decoder) Decode(reader io.Reader) (*Library, error) {
	context := newDecodeContext(d.limits)
	err := d.scanner.Scan(reader, context.HandleEvent)
	if err != nil {
		return nil, err
	}
	return context.Library, nil
}

func newDecodeContext(limits *DecodeLimits) *decodeContext {
	return &decodeContext{
		Limits:          limits,
		Library:         new(Library),
		CurrentMaterial: nil,
	}
}

type decodeContext struct {
	Limits          *DecodeLimits
	Library         *Library
	CurrentMaterial *Material
}

func (c *decodeContext) HandleEvent(event common.Event) error {
	switch actual := event.(type) {
	case scanMTL.MaterialEvent:
		return c.handleMaterial(actual)
	case scanMTL.RGBAmbientColorEvent:
		return c.handleAmbientColor(actual)
	case scanMTL.RGBDiffuseColorEvent:
		return c.handleDiffuseColor(actual)
	case scanMTL.RGBSpecularColorEvent:
		return c.handleSpecularColor(actual)
	case scanMTL.RGBEmissiveColorEvent:
		return c.handleEmissiveColor(actual)
	case scanMTL.RGBTransmissionFilterEvent:
		return c.handleTransmissionFilter(actual)
	case scanMTL.SpecularExponentEvent:
		return c.handleSpecularExponent(actual)
	case scanMTL.IllumEvent:
		return c.handleIllum(actual)
	case scanMTL.DissolveEvent:
		return c.handleDissolve(actual)
	case scanMTL.AmbientTextureEvent:
		return c.handleAmbientTexture(actual)
	case scanMTL.DiffuseTextureEvent:
		return c.handleDiffuseTexture(actual)
	case scanMTL.SpecularTextureEvent:
		return c.handleSpecularTexture(actual)
	case scanMTL.EmissiveTextureEvent:
		return c.handleEmissiveTexture(actual)
	case scanMTL.SpecularExponentTextureEvent:
		return c.handleSpecularExponentTexture(actual)
	case scanMTL.DissolveTextureEvent:
		return c.handleDissolveTexture(actual)
	case scanMTL.BumpTextureEvent:
		return c.handleBumpTexture(actual)
	}
	return nil
}

func (c *decodeContext) handleMaterial(event scanMTL.MaterialEvent) error {
	if len(c.Library.Materials) >= c.Limits.MaxMaterialCount {
		return errors.New("Max number of allowed materials reached!")
	}
	c.CurrentMaterial = DefaultMaterial()
	c.CurrentMaterial.Name = event.MaterialName
	c.Library.Materials = append(c.Library.Materials, c.CurrentMaterial)
	return nil
}

func (c *decodeContext) handleAmbientColor(event scanMTL.RGBAmbientColorEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.AmbientColor = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleDiffuseColor(event scanMTL.RGBDiffuseColorEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.DiffuseColor = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleSpecularColor(event scanMTL.RGBSpecularColorEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.SpecularColor = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleEmissiveColor(event scanMTL.RGBEmissiveColorEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.EmissiveColor = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleTransmissionFilter(event scanMTL.RGBTransmissionFilterEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.TransmissionFilter = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleSpecularExponent(event scanMTL.SpecularExponentEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.SpecularExponent = event.Amount
	return nil
}

func (c *decodeContext) handleIllum(event scanMTL.IllumEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.Illum = event.Amount
	return nil
}

func (c *decodeContext) handleDissolve(event scanMTL.DissolveEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.Dissolve = event.Amount
	return nil
}

func (c *decodeContext) handleAmbientTexture(event scanMTL.AmbientTextureEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.AmbientTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleDiffuseTexture(event scanMTL.DiffuseTextureEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.DiffuseTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleSpecularTexture(event scanMTL.SpecularTextureEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.SpecularTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleEmissiveTexture(event scanMTL.EmissiveTextureEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.EmissiveTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleSpecularExponentTexture(event scanMTL.SpecularExponentTextureEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.SpecularExponentTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleDissolveTexture(event scanMTL.DissolveTextureEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.DissolveTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleBumpTexture(event scanMTL.BumpTextureEvent) error {
	if c.CurrentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.CurrentMaterial.BumpTexture = event.TexturePath
	return nil
}

func (c *decodeContext) newMissingMaterialError() error {
	return errors.New("Declaration requires a material to be defined first.")
}
