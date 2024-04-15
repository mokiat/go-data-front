# Go Library for Wavefront file parsing

[![Go Reference](https://pkg.go.dev/badge/github.com/mokiat/go-data-front.svg)](https://pkg.go.dev/github.com/mokiat/go-data-front)
![Build Status](https://github.com/mokiat/go-data-front/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/mokiat/go-data-front)](https://goreportcard.com/report/github.com/mokiat/go-data-front)

A Go library for reading Wavefront 3D model resources (OBJ, MTL).

## User's Guide

The OBJ and MTL file formats are one of the most popular 3D model formats used at the moment. This is mainly due to their simplicity.

I am not sure what is the exact history of these file formats, but I believe they were first introduced by the Wavefront Technologies company for their 3D software. Now they are used in practically all 3D modeling software out there.

### OBJ

OBJ files are used to describe the coordinates, connections and shapes that make up a 3D model.

```
v -1.0 1.0 0.0
v -1.0 -1.0 0.0
v 1.0 -1.0 0.0
v 1.0 1.0 0.0

o Rectangle
f 1 2 3
f 1 3 4
```

The library provides two mechanisms for OBJ reading.

#### Scanning

The Scanner API allows you to scan through an OBJ file and receive events for each supported element that is scanned along the way. It is up to you to record the data and to correlate it afterward.

**Example**

```go
func main() {
	file, _ := os.Open("example.obj")
	defer file.Close()

	scanner := obj.NewScanner()
	scanner.Scan(file, handleEvent)
}

func handleEvent(event common.Event) error {
	switch actual := event.(type) {
	case common.CommentEvent:
		fmt.Printf("Comment: %s\n", actual.Comment)
		return nil
	case obj.ObjectEvent:
		fmt.Printf("Object with name: %s\n", actual.ObjectName)
		return nil
	}
	return nil
}
```

You can find the API documentation **[here](https://pkg.go.dev/github.com/mokiat/go-data-front/scanner/obj)**.


#### Decoder

The Decoder API uses the Scanner API and does the correlation for you. What you end up with is an object model that represents the parsed OBJ file.

**Example**

```go
func main() {
	file, _ := os.Open("example.obj")
	defer file.Close()

	decoder := obj.NewDecoder(obj.DefaultLimits())
	model, _ := decoder.Decode(file)

	fmt.Printf("Model has %d vertices.\n", len(model.Vertices))
	fmt.Printf("Model has %d texture coordinates.\n", len(model.TexCoords))
	fmt.Printf("Model has %d normals.\n", len(model.Normals))
	fmt.Printf("Model has %d objects.\n", len(model.Objects))
	fmt.Printf("First object has name: %s\n", model.Objects[0].Name)
}
```

You can find the API documentation **[here](https://pkg.go.dev/github.com/mokiat/go-data-front/decoder/obj)**.

### MTL

MTL files are optional and are present when a 3D model uses materials.

```
newmtl TestMaterial
Ka 1.0 0.5 0.1
Kd 0.5 0.7 0.3
Ks 0.2 0.4 0.8
Ns 650
d 0.7
map_Kd vehicle.png
```

The library provides two mechanisms for MTL reading.

#### Scanner

The Scanner API allows you to scan through an MTL file and receive events for each supported element that is scanned along the way. It is up to you to record the data and to correlate it afterward.

**Example**

```go
func main() {
	file, _ := os.Open("example.mtl")
	defer file.Close()

	scanner := mtl.NewScanner()
	scanner.Scan(file, handleEvent)
}

func handleEvent(event common.Event) error {
	switch actual := event.(type) {
	case common.CommentEvent:
		fmt.Printf("Comment: %s\n", actual.Comment)
		return nil
	case mtl.MaterialEvent:
		fmt.Printf("Material with name: %s\n", actual.MaterialName)
		return nil
	}
	return nil
}
```

You can find the API documentation **[here](https://pkg.go.dev/github.com/mokiat/go-data-front/scanner/mtl)**.

#### Decoder

The Decoder API uses the Scanner API and does the correlation for you. What you end up with is an object model that represents the parsed MTL file.

**Example**

```go
func main() {
	file, _ := os.Open("example.mtl")
	defer file.Close()

	decoder := mtl.NewDecoder(mtl.DefaultLimits())
	model, _ := decoder.Decode(file)

	fmt.Printf("Library has %d materials\n", len(model.Materials))
	fmt.Printf("First material has name: %s\n",	model.Materials[0].Name)
}
```

You can find the API documentation **[here](https://pkg.go.dev/github.com/mokiat/go-data-front/decoder/mtl)**.

## Developer's Guide

This library uses the **[Ginkgo](https://github.com/onsi/ginkgo)** tool for testing.

You can run all the tests with the following command from within the project path.

```bash
go run github.com/onsi/ginkgo/v2/ginkgo -r
```

## Notes

This library was implemented to support other projects of mine. As such, it supports only a subset of the OBJ and MTL specifications. The features that are provided should be sufficient for most use cases.

This repository is a rewrite of the **[https://github.com/mokiat/java-data-front](https://github.com/mokiat/java-data-front)** one, which is in Java.
