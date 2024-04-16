package mtl

import (
	"fmt"
	"io"

	"github.com/mokiat/go-data-front/common"
	mtlscan "github.com/mokiat/go-data-front/scanner/mtl"
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
		scanner: mtlscan.NewScanner(),
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
	return context.Library(), nil
}

func newDecodeContext(limits *DecodeLimits) *decodeContext {
	return &decodeContext{
		limits:          limits,
		library:         new(Library),
		currentMaterial: nil,
	}
}

type decodeContext struct {
	limits          *DecodeLimits
	library         *Library
	currentMaterial *Material
}

func (c *decodeContext) Library() *Library {
	return c.library
}

func (c *decodeContext) HandleEvent(event common.Event) error {
	switch actual := event.(type) {
	case mtlscan.MaterialEvent:
		return c.handleMaterial(actual)
	case mtlscan.RGBAmbientColorEvent:
		return c.handleAmbientColor(actual)
	case mtlscan.RGBDiffuseColorEvent:
		return c.handleDiffuseColor(actual)
	case mtlscan.RGBSpecularColorEvent:
		return c.handleSpecularColor(actual)
	case mtlscan.RGBEmissiveColorEvent:
		return c.handleEmissiveColor(actual)
	case mtlscan.RGBTransmissionFilterEvent:
		return c.handleTransmissionFilter(actual)
	case mtlscan.SpecularExponentEvent:
		return c.handleSpecularExponent(actual)
	case mtlscan.DissolveEvent:
		return c.handleDissolve(actual)
	case mtlscan.IlluminationEvent:
		return c.handleIllumination(actual)
	case mtlscan.AmbientTextureEvent:
		return c.handleAmbientTexture(actual)
	case mtlscan.DiffuseTextureEvent:
		return c.handleDiffuseTexture(actual)
	case mtlscan.SpecularTextureEvent:
		return c.handleSpecularTexture(actual)
	case mtlscan.EmissiveTextureEvent:
		return c.handleEmissiveTexture(actual)
	case mtlscan.SpecularExponentTextureEvent:
		return c.handleSpecularExponentTexture(actual)
	case mtlscan.DissolveTextureEvent:
		return c.handleDissolveTexture(actual)
	case mtlscan.BumpTextureEvent:
		return c.handleBumpTexture(actual)
	}
	return nil
}

func (c *decodeContext) handleMaterial(event mtlscan.MaterialEvent) error {
	if len(c.library.Materials) >= c.limits.MaxMaterialCount {
		return fmt.Errorf("%w: max number of materials reached", common.ErrLimitsExceeded)
	}
	c.currentMaterial = DefaultMaterial()
	c.currentMaterial.Name = event.MaterialName
	c.library.Materials = append(c.library.Materials, c.currentMaterial)
	return nil
}

func (c *decodeContext) handleAmbientColor(event mtlscan.RGBAmbientColorEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.AmbientColor = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleDiffuseColor(event mtlscan.RGBDiffuseColorEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.DiffuseColor = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleSpecularColor(event mtlscan.RGBSpecularColorEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.SpecularColor = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleEmissiveColor(event mtlscan.RGBEmissiveColorEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.EmissiveColor = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleTransmissionFilter(event mtlscan.RGBTransmissionFilterEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.TransmissionFilter = RGBColor{
		R: event.R,
		G: event.G,
		B: event.B,
	}
	return nil
}

func (c *decodeContext) handleSpecularExponent(event mtlscan.SpecularExponentEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.SpecularExponent = event.Amount
	return nil
}

func (c *decodeContext) handleDissolve(event mtlscan.DissolveEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.Dissolve = event.Amount
	return nil
}

func (c *decodeContext) handleIllumination(event mtlscan.IlluminationEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.Illumination = event.Model
	return nil
}

func (c *decodeContext) handleAmbientTexture(event mtlscan.AmbientTextureEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.AmbientTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleDiffuseTexture(event mtlscan.DiffuseTextureEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.DiffuseTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleSpecularTexture(event mtlscan.SpecularTextureEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.SpecularTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleEmissiveTexture(event mtlscan.EmissiveTextureEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.EmissiveTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleSpecularExponentTexture(event mtlscan.SpecularExponentTextureEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.SpecularExponentTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleDissolveTexture(event mtlscan.DissolveTextureEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.DissolveTexture = event.TexturePath
	return nil
}

func (c *decodeContext) handleBumpTexture(event mtlscan.BumpTextureEvent) error {
	if c.currentMaterial == nil {
		return c.newMissingMaterialError()
	}
	c.currentMaterial.BumpTexture = event.TexturePath
	return nil
}

func (c *decodeContext) newMissingMaterialError() error {
	return fmt.Errorf("%w: material declaration outside of material block", common.ErrInvalid)
}
