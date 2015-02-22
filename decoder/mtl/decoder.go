package mtl

import (
	"io"

	"github.com/momchil-atanasov/go-data-front/common"
	scanMTL "github.com/momchil-atanasov/go-data-front/scanner/mtl"
)

// DecodeLimits specifies restrictions on parsing.
//
// One will generally use this to limit the number of
// data that is parsed in order to prevent out of memory
// errors
type DecodeLimits struct {

	// MaxMaterialCount specifies the maximum number of
	// material declarations that can be parsed.
	MaxMaterialCount int64
}

// Decoder is an API that allows one to decode MTL
// Wavefront resources into an object model.
type Decoder interface {

	// Decode decodes the MTL Wavefron resource, specified
	// through the io.Reader, into a Library model.
	//
	// If decoding fails for some reason, an error is returned.
	Decode(io.Reader) (*Library, error)
}

// NewDecoder creates a new Decoder instance with the
// specified DecodeLimits.
func NewDecoder(limits DecodeLimits) Decoder {
	return &decoder{
		scanner: scanMTL.NewScanner(),
	}
}

type decoder struct {
	scanner common.Scanner
}

func (d *decoder) Decode(reader io.Reader) (*Library, error) {
	context := newDecodeContext()
	err := d.scanner.Scan(reader, context.HandleEvent)
	if err != nil {
		return nil, err
	}
	return context.Library, nil
}

func newDecodeContext() *decodeContext {
	return &decodeContext{
		Library:         new(Library),
		CurrentMaterial: nil,
	}
}

type decodeContext struct {
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
	case scanMTL.RGBTransmissionFilterEvent:
		return c.handleTransmissionFilter(actual)
	case scanMTL.SpecularExponentEvent:
		return c.handleSpecularExponent(actual)
	case scanMTL.DissolveEvent:
		return c.handleDissolve(actual)
	case scanMTL.AmbientTextureEvent:
		return c.handleAmbientTexture(actual)
	case scanMTL.DiffuseTextureEvent:
		return c.handleDiffuseTexture(actual)
	case scanMTL.SpecularTextureEvent:
		return c.handleSpecularTexture(actual)
	case scanMTL.SpecularExponentTextureEvent:
		return c.handleSpecularExponentTexture(actual)
	case scanMTL.DissolveTextureEvent:
		return c.handleDissolveTexture(actual)
	}
	return nil
}

func (c *decodeContext) handleMaterial(event scanMTL.MaterialEvent) error {
	// TODO: Check limits!!!
	c.CurrentMaterial = DefaultMaterial()
	c.CurrentMaterial.Name = event.MaterialName
	c.Library.Materials = append(c.Library.Materials, c.CurrentMaterial)
	return nil
}

func (c *decodeContext) handleAmbientColor(event scanMTL.RGBAmbientColorEvent) error {
	// TODO: Verify there is a current material
	c.CurrentMaterial.AmbientColor = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleDiffuseColor(event scanMTL.RGBDiffuseColorEvent) error {
	// TODO: Verify there is a current material
	c.CurrentMaterial.DiffuseColor = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleSpecularColor(event scanMTL.RGBSpecularColorEvent) error {
	// TODO: Verify there is a current material
	c.CurrentMaterial.SpecularColor = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleTransmissionFilter(event scanMTL.RGBTransmissionFilterEvent) error {
	// TODO: Verify there is a current material
	c.CurrentMaterial.TransmissionFilter = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleSpecularExponent(event scanMTL.SpecularExponentEvent) error {
	// TODO: Verify there is a current material
	c.CurrentMaterial.SpecularExponent = event.Amount
	return nil
}

func (c *decodeContext) handleDissolve(event scanMTL.DissolveEvent) error {
	// TODO: Verify there is a current material
	c.CurrentMaterial.Dissolve = event.Amount
	return nil
}

func (c *decodeContext) handleAmbientTexture(event scanMTL.AmbientTextureEvent) error {
	// TODO: Verify there is a current material
	c.CurrentMaterial.AmbientTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleDiffuseTexture(event scanMTL.DiffuseTextureEvent) error {
	// TODO: Verify there is a current material
	c.CurrentMaterial.DiffuseTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleSpecularTexture(event scanMTL.SpecularTextureEvent) error {
	// TODO: Verify there is a current material
	c.CurrentMaterial.SpecularTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleSpecularExponentTexture(event scanMTL.SpecularExponentTextureEvent) error {
	// TODO: Verify there is a current material
	c.CurrentMaterial.SpecularExponentTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleDissolveTexture(event scanMTL.DissolveTextureEvent) error {
	// TODO: Verify there is a current material
	c.CurrentMaterial.DissolveTexture = event.TexturePath
	return nil
}
