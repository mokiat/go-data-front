package common

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Line represents a single logical line in a Wavefront file.
// You should get a hold of such a structure through the use of the
// LineScanner API.
type Line struct {
	// contains filtered or unexported fields
	line     string
	segments []string
}

// IsBlank returns whether the current logical line is blank
func (l Line) IsBlank() bool {
	return l.line == ""
}

// IsComment returns whether the current logical line represents a comment
func (l Line) IsComment() bool {
	return strings.HasPrefix(l.line, "#")
}

// Comment returns the comment held by this logical line. One should first
// use IsComment to assure that this logical line is indeed a comment.
func (l Line) Comment() string {
	return strings.TrimSpace(strings.TrimPrefix(l.line, "#"))
}

// IsCommand returns whether the current logical line represents a command
// (e.g. usemtl, mtllib, vn, etc.)
func (l Line) IsCommand() bool {
	return len(l.segments) > 0
}

// HasCommandName checks whether the command held by this logical line has
// the specified name. One should first use IsCommand to assure that this
// logical line is indeed a command.
func (l Line) HasCommandName(name string) bool {
	return l.CommandName() == name
}

// CommandName returns the name of the command. One should first use the
// IsCommand method to assure that this logical line is indeed a command.
func (l Line) CommandName() string {
	return l.segments[0]
}

// ParamCount returns the number of parameters provided with the current
// command. Parameters are indexed from `0` upto the number (excluding)
// returned by this function.
func (l Line) ParamCount() int {
	count := len(l.segments) - 1
	if count < 0 {
		count = 0
	}
	return count
}

// Returns the parameter, converted to a string, at the specified index
func (l Line) StringParam(index int) string {
	return l.segments[index+1]
}

// Returns the parameter, converted to an integer, at the specified index
func (l Line) IntParam(index int) (int64, error) {
	value, err := strconv.ParseInt(l.StringParam(index), 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// Returns the parameter, converted to a float, at the specified index
func (l Line) FloatParam(index int) (float64, error) {
	value, err := strconv.ParseFloat(l.StringParam(index), 64)
	if err != nil {
		return 0.0, err
	}
	return value, nil
}

// LineScanner is an API that allows one to scan Wavefront files
// one logical line at a time.
type LineScanner interface {
	// Scan scans the next logical line and returns whether that was successful.
	// A scan is not successful when the end of the file is reached or when an
	// error has occurred during scanning.
	Scan() bool

	// Err is used to get the scanning error, should one have actually occurred.
	// If `nil` is returned, then scanning was successful and one can proceed
	// and access the line through the Line method.
	Err() error

	// Line returns the last scanned logical line.
	Line() Line
}

type lineScanner struct {
	scanner      *bufio.Scanner
	segmentRegex *regexp.Regexp
	lineBuffer   bytes.Buffer
	scanLine     Line
	scanErr      error
}

// NewLineScanner creates a new LineScanner instance that uses the
// specified io.Reader to read a Wavefront resource.
func NewLineScanner(reader io.Reader) LineScanner {
	regex, err := regexp.Compile("[\\s]+")
	if err != nil {
		panic(err)
	}
	return &lineScanner{
		scanner:      bufio.NewScanner(reader),
		segmentRegex: regex,
		lineBuffer:   bytes.Buffer{},
	}
}

func (s *lineScanner) Scan() bool {
	s.lineBuffer.Reset()

	scanIterations := 0
	for s.scanner.Scan() {
		if s.scanner.Err() != nil {
			s.scanErr = s.scanner.Err()
			return false
		}

		scanIterations++
		line := s.scanner.Text()
		if strings.HasSuffix(line, "\\") {
			s.lineBuffer.WriteString(strings.TrimSuffix(line, "\\"))
		} else {
			s.lineBuffer.WriteString(line)
			break
		}
	}

	s.scanLine = s.createLine(s.lineBuffer.String())
	return scanIterations > 0
}

func (s *lineScanner) createLine(logicalLine string) Line {
	return Line{
		line:     strings.TrimSpace(logicalLine),
		segments: s.segmentRegex.Split(logicalLine, -1),
	}
}

func (s *lineScanner) Err() error {
	return s.scanErr
}

func (s *lineScanner) Line() Line {
	return s.scanLine
}
