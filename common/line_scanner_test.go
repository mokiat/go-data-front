package common_test

import (
	"bytes"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mokiat/go-data-front/common"
)

var _ = Describe("LineScanner", func() {
	var lineScanner common.LineScanner

	openForScanning := func(fileName string) common.LineScanner {
		GinkgoHelper()
		data, err := os.ReadFile(filepath.Join("testdata", fileName))
		Expect(err).ToNot(HaveOccurred())
		return common.NewLineScanner(bytes.NewReader(data))
	}

	readNextLine := func() common.Line {
		GinkgoHelper()
		Expect(lineScanner.Scan()).To(BeTrue())
		Expect(lineScanner.Err()).ToNot(HaveOccurred())
		return lineScanner.Line()
	}

	assertNoMoreLines := func() {
		GinkgoHelper()
		Expect(lineScanner.Scan()).To(BeFalse())
		Expect(lineScanner.Err()).ToNot(HaveOccurred())
	}

	assertIsBlank := func(line common.Line) {
		GinkgoHelper()
		Expect(line.IsBlank()).To(BeTrue())
	}

	assertIsNotBlank := func(line common.Line) {
		GinkgoHelper()
		Expect(line.IsBlank()).To(BeFalse())
	}

	assertIsComment := func(line common.Line, value string) {
		GinkgoHelper()
		Expect(line.IsComment()).To(BeTrue())
		Expect(line.Comment()).To(Equal(value))
	}

	assertIsNotComment := func(line common.Line) {
		GinkgoHelper()
		Expect(line.IsComment()).To(BeFalse())
	}

	assertIsCommand := func(line common.Line, name string) {
		GinkgoHelper()
		Expect(line.IsCommand()).To(BeTrue())
		Expect(line.HasCommandName(name)).To(BeTrue())
		Expect(line.CommandName()).To(Equal(name))
	}

	assertCommandParams := func(line common.Line, params ...string) {
		GinkgoHelper()
		Expect(line.ParamCount()).To(Equal(len(params)))
		for index, param := range params {
			Expect(line.StringParam(index)).To(Equal(param))
		}
	}

	BeforeEach(func() {
		lineScanner = nil
	})

	Describe("scanning blank lines", func() {
		var (
			comment common.Line
			blank   common.Line
			command common.Line
		)

		BeforeEach(func() {
			lineScanner = openForScanning("line_scanner_blank_lines.txt")
			comment = readNextLine()
			blank = readNextLine()
			command = readNextLine()
			assertNoMoreLines()
		})

		It("can scan blank lines", func() {
			assertIsBlank(blank)
		})

		It("ignores all non-blanks", func() {
			assertIsNotBlank(comment)
			assertIsNotBlank(command)
		})
	})

	Describe("scanning comments", func() {
		var (
			firstComment           common.Line
			command                common.Line
			tightComment           common.Line
			point                  common.Line
			firstMultiLineComment  common.Line
			secondMultiLineComment common.Line
			normal                 common.Line
			emptyComment           common.Line
			blank                  common.Line
			object                 common.Line
			hashComment            common.Line
			complicated            common.Line
			lastComment            common.Line
		)

		BeforeEach(func() {
			lineScanner = openForScanning("line_scanner_comments.txt")
			firstComment = readNextLine()
			command = readNextLine()
			tightComment = readNextLine()
			point = readNextLine()
			firstMultiLineComment = readNextLine()
			secondMultiLineComment = readNextLine()
			normal = readNextLine()
			emptyComment = readNextLine()
			blank = readNextLine()
			object = readNextLine()
			hashComment = readNextLine()
			complicated = readNextLine()
			lastComment = readNextLine()
			assertNoMoreLines()
		})

		It("can scan all the comment lines", func() {
			assertIsComment(firstComment, "Comment at file start")
			assertIsComment(tightComment, "Comment that is right next to special char")
			assertIsComment(firstMultiLineComment, "This comment uses")
			assertIsComment(secondMultiLineComment, "two lines")
			assertIsComment(emptyComment, "")
			assertIsComment(hashComment, "Previous comment was empty. This one contains the # character twice.")
			assertIsComment(lastComment, "Comment at file end")
		})

		It("ignores all non-comments", func() {
			assertIsNotComment(command)
			assertIsNotComment(point)
			assertIsNotComment(normal)
			assertIsNotComment(blank)
			assertIsNotComment(object)
			assertIsNotComment(complicated)
		})
	})

	Describe("scanning command names", func() {
		BeforeEach(func() {
			lineScanner = openForScanning("line_scanner_command_names.txt")
		})

		It("can scan all the command names", func() {
			firstCommand := readNextLine()
			assertIsCommand(firstCommand, "firstCommand")

			secondCommand := readNextLine()
			assertIsCommand(secondCommand, "secondCommand")

			thirdCommand := readNextLine()
			assertIsCommand(thirdCommand, "thirdCommand")

			assertNoMoreLines()
		})
	})

	Describe("scanning parameters", func() {
		var (
			noParams         common.Line
			blank            common.Line
			stringParams     common.Line
			intParams        common.Line
			floatParams      common.Line
			referencesParams common.Line
		)

		BeforeEach(func() {
			lineScanner = openForScanning("line_scanner_parameters.txt")
			noParams = readNextLine()
			blank = readNextLine()
			stringParams = readNextLine()
			intParams = readNextLine()
			floatParams = readNextLine()
			referencesParams = readNextLine()
			assertNoMoreLines()
		})

		It("can scan commands without parameters", func() {
			Expect(noParams.ParamCount()).To(BeZero())
			Expect(blank.ParamCount()).To(BeZero())
		})

		It("can scan proper number of string parameters", func() {
			Expect(stringParams.ParamCount()).To(Equal(4))
		})

		It("can scan proper number of int parameters", func() {
			Expect(intParams.ParamCount()).To(Equal(2))
		})

		It("can scan proper number of float parameters", func() {
			Expect(floatParams.ParamCount()).To(Equal(3))
		})

		It("can scan string parameters when strings", func() {
			Expect(stringParams.StringParam(0)).To(Equal("hello"))
			Expect(stringParams.StringParam(1)).To(Equal("complex/param"))
			Expect(stringParams.StringParam(2)).To(Equal("\"quoted\""))
			Expect(stringParams.StringParam(3)).To(Equal("?123?"))
		})

		It("can scan string parameters when ints", func() {
			Expect(intParams.StringParam(0)).To(Equal("3"))
			Expect(intParams.StringParam(1)).To(Equal("-50"))
		})

		It("can scan string parameters when floats", func() {
			Expect(floatParams.StringParam(0)).To(Equal("1.0"))
			Expect(floatParams.StringParam(1)).To(Equal("-13.33"))
			Expect(floatParams.StringParam(2)).To(Equal(".3"))
		})

		It("can scan int parameters when ints", func() {
			value, err := intParams.IntParam(0)
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(Equal(int64(3)))

			value, err = intParams.IntParam(1)
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(Equal(int64(-50)))
		})

		It("cannot scan int parameters when strings", func() {
			_, err := stringParams.IntParam(0)
			Expect(err).To(HaveOccurred())
		})

		It("cannot scan int parameters when floats", func() {
			_, err := floatParams.IntParam(0)
			Expect(err).To(HaveOccurred())
		})

		It("can scan float parameters when floats", func() {
			value, err := floatParams.FloatParam(0)
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(BeNumerically("~", 1.0))

			value, err = floatParams.FloatParam(1)
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(BeNumerically("~", -13.33))

			value, err = floatParams.FloatParam(2)
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(BeNumerically("~", 0.3))
		})

		It("cannot scan float parameters when strings", func() {
			_, err := stringParams.FloatParam(0)
			Expect(err).To(HaveOccurred())
		})

		It("can scan float parameters when ints", func() {
			value, err := intParams.FloatParam(0)
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(BeNumerically("~", 3))

			value, err = intParams.FloatParam(1)
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(BeNumerically("~", -50))
		})

		It("can scan reference sets when strings", func() {
			referenceSet := stringParams.ReferenceSetParam(0)
			Expect(referenceSet.Count()).To(Equal(1))
			Expect(referenceSet.StringReference(0)).To(Equal("hello"))
		})

		It("can scan reference sets when ints", func() {
			referenceSet := intParams.ReferenceSetParam(0)
			Expect(referenceSet.Count()).To(Equal(1))
			Expect(referenceSet.StringReference(0)).To(Equal("3"))
		})

		It("can scan reference sets when floats", func() {
			referenceSet := floatParams.ReferenceSetParam(0)
			Expect(referenceSet.Count()).To(Equal(1))
			Expect(referenceSet.StringReference(0)).To(Equal("1.0"))
		})

		Describe("ReferenceSet", func() {
			var (
				singleReferenceSet  common.ReferenceSet
				doubleReferenceSet  common.ReferenceSet
				trippleReferenceSet common.ReferenceSet
				skipReferenceSet    common.ReferenceSet
			)

			BeforeEach(func() {
				singleReferenceSet = referencesParams.ReferenceSetParam(0)
				doubleReferenceSet = referencesParams.ReferenceSetParam(1)
				trippleReferenceSet = referencesParams.ReferenceSetParam(2)
				skipReferenceSet = referencesParams.ReferenceSetParam(3)
			})

			It("can parse single references", func() {
				Expect(singleReferenceSet.Count()).To(Equal(1))
				Expect(singleReferenceSet.StringReference(0)).To(Equal("1"))
			})

			It("can parse double references", func() {
				Expect(doubleReferenceSet.Count()).To(Equal(2))
				Expect(doubleReferenceSet.StringReference(0)).To(Equal("2"))
				Expect(doubleReferenceSet.StringReference(1)).To(Equal("3"))
			})

			It("can parse trippe references", func() {
				Expect(trippleReferenceSet.Count()).To(Equal(3))
				Expect(trippleReferenceSet.StringReference(0)).To(Equal("4"))
				Expect(trippleReferenceSet.StringReference(1)).To(Equal("abc"))
				Expect(trippleReferenceSet.StringReference(2)).To(Equal("6"))
			})

			It("can parse skip references", func() {
				Expect(skipReferenceSet.Count()).To(Equal(3))
				Expect(skipReferenceSet.StringReference(0)).To(Equal("7.3"))
				Expect(skipReferenceSet.IsBlank(0)).To(BeFalse())
				Expect(skipReferenceSet.StringReference(1)).To(Equal(""))
				Expect(skipReferenceSet.IsBlank(1)).To(BeTrue())
				Expect(skipReferenceSet.StringReference(2)).To(Equal("8"))
				Expect(skipReferenceSet.IsBlank(2)).To(BeFalse())
			})

			It("can return int reference when int", func() {
				reference, err := singleReferenceSet.IntReference(0)
				Expect(err).ToNot(HaveOccurred())
				Expect(reference).To(Equal(int64(1)))
			})

			It("cannot return int reference when string", func() {
				_, err := trippleReferenceSet.IntReference(1)
				Expect(err).To(HaveOccurred())
			})

			It("cannot return int reference when float", func() {
				_, err := skipReferenceSet.IntReference(0)
				Expect(err).To(HaveOccurred())
			})

			It("can return float reference when float", func() {
				reference, err := skipReferenceSet.FloatReference(0)
				Expect(err).ToNot(HaveOccurred())
				Expect(reference).To(Equal(float64(7.3)))
			})

			It("can return float reference when int", func() {
				reference, err := singleReferenceSet.FloatReference(0)
				Expect(err).ToNot(HaveOccurred())
				Expect(reference).To(Equal(float64(1.0)))
			})

			It("cannot return float reference when string", func() {
				_, err := trippleReferenceSet.FloatReference(1)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("scanning logical lines", func() {
		var (
			first  common.Line
			second common.Line
			third  common.Line
		)

		BeforeEach(func() {
			lineScanner = openForScanning("line_scanner_logical_line.txt")
			first = readNextLine()
			second = readNextLine()
			third = readNextLine()
			assertNoMoreLines()
		})

		It("can scan logical lines", func() {
			assertIsCommand(first, "firstLogical")
			assertCommandParams(first, "has", "tight", "separation")

			assertIsCommand(second, "secondLogical")
			assertCommandParams(second, "has", "blank", "space", "separation")

			assertIsCommand(third, "thirdLogical")
			assertCommandParams(third, "ends", "with", "separation")
		})
	})
})
