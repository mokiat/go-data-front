package mtl

// RGBColor represents a color represented
// by the three basic colors - Red, Green, and Blue
type RGBColor struct {

	// Specifies the amount of Red in this color. Usually in the
	// range 0.0 to 1.0.
	R float64

	// Specifies the amount of Green in this color. Usually in the
	// range 0.0 to 1.0.
	G float64

	// Specifies the amount of Blue in this color. Usually in the
	// range 0.0 to 1.0.
	B float64
}

// Material represents a material in a MTL
// wavefront resource.
//
// Materials define how objects should be rendered.
type Material struct {

	// Name holds the name of this material.
	Name string

	// AmbientColor holds the ambient color to be used
	// when rendering objects.
	AmbientColor RGBColor

	// DiffuseColor holds the diffuse color to be used
	// when rendering objects.
	DiffuseColor RGBColor

	// SpecularColor holds the specular color to be used
	// when rendering objects.
	SpecularColor RGBColor

	// EmissiveColor holds the emissive color to be used
	// when rendering objects.
	EmissiveColor RGBColor

	// TransmissionFilter holds the filter to be used on
	// colors when rendering objects.
	TransmissionFilter RGBColor

	// SpecularExponent defines the specular exponent for
	// this material.
	//
	// The specular exponent defines the sharpness of the
	// specular reflection. The value will generally range
	// between 0.0 (sharp) and 1000.0 (soft).
	SpecularExponent float64

	// Dissolve defines the dissolve for this material.
	//
	// Dissolve indicates how much an object should blend.
	// The value should range between 0.0 (fully transparent)
	// and 1.0 (opaque).
	Dissolve float64

	// AmbientTexture defines the location of the ambient
	// texture to be used when rendering objects.
	//
	// If this value is the empty string, then there is no
	// ambient texture provided.
	AmbientTexture string

	// DiffuseTexture defines the location of the diffuse
	// texture to be used when rendering objects.
	//
	// If this value is the empty string, then there is no
	// diffuse texture provided.
	DiffuseTexture string

	// SpecularTexture defines the location of the specular
	// texture to be used when rendering objects.
	//
	// If this value is the empty string, then there is no
	// specular texture provided.
	SpecularTexture string

	// EmissiveTexture defines the location of the emissive
	// texture to be used when rendering objects.
	//
	// If this value is the empty string, then there is no
	// emissive texture provided.
	EmissiveTexture string

	// SpecularExponentTexture defines the location of the specular
	// exponent texture to be used when rendering objects.
	//
	// If this value is the empty string, then there is no
	// specular exponent texture provided.
	SpecularExponentTexture string

	// DissolveTexture defines the location of the dissolve
	// texture to be used when rendering objects.
	//
	// If this value is the empty string, then there is no
	// dissolve texture provided.
	DissolveTexture string

	// BumpTexture defines the location of the bump
	// texture to be used when rendering objects.
	//
	// If this value is the empty string, then there is no
	// Bump texture provided.
	BumpTexture string

	// http://paulbourke.net/dataformats/mtl/
	// The "illum" statement specifies the illumination model to use in the
	// material.  Illumination models are mathematical equations that represent
	// various material lighting and shading effects.

	//  "illum_#"can be a number from 0 to 10.  The illumination models are
	// summarized below

	//  Illumination    Properties that are turned on in the
	//  model           Property Editor

	//  0		Color on and Ambient off
	//  1		Color on and Ambient on
	//  2		Highlight on
	//  3		Reflection on and Ray trace on
	//  4		Transparency: Glass on
	// 		 Reflection: Ray trace on
	//  5		Reflection: Fresnel on and Ray trace on
	//  6		Transparency: Refraction on
	// 		 Reflection: Fresnel off and Ray trace on
	//  7		Transparency: Refraction on
	// 		 Reflection: Fresnel on and Ray trace on
	//  8		Reflection on and Ray trace off
	//  9		Transparency: Glass on
	// 		 Reflection: Ray trace off
	//  10		Casts shadows onto invisible surfaces
	Illum int64
}

// DefaultMaterial returns a new Material which is
// initialized with some proper default values
func DefaultMaterial() *Material {
	return &Material{
		AmbientColor: RGBColor{
			R: 1.0,
			G: 1.0,
			B: 1.0,
		},
		DiffuseColor: RGBColor{
			R: 1.0,
			G: 1.0,
			B: 1.0,
		},
		Dissolve: 1.0,
		TransmissionFilter: RGBColor{
			R: 1.0,
			G: 1.0,
			B: 1.0,
		},
	}
}

// Library represents a material library.
//
// A material library can be though of as a
// single MTL resource which can contain information
// on multiple materials.
type Library struct {

	// Materials contains a list of all the materials that
	// were defined in the given library.
	Materials []*Material
}

// FindMaterial finds a material in the given Library
// with the specified name or returns false, otherwise.
func (l *Library) FindMaterial(name string) (*Material, bool) {
	for _, material := range l.Materials {
		if material.Name == name {
			return material, true
		}
	}
	return nil, false
}
