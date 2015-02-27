package common_test

import (
	"fmt"
	"io"
	"os"

	. "github.com/momchil-atanasov/go-data-front/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LineScanner", func() {
	var reader io.ReadCloser
	var lineScanner LineScanner

	openFile := func(filename string) io.ReadCloser {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		return file
	}

	openForScanning := func(fileName string) {
		reader = openFile(fmt.Sprintf("common_test_res/%s", fileName))
		lineScanner = NewLineScanner(reader)
	}

	readNextLine := func() Line {
		Ω(lineScanner.Scan()).Should(BeTrue())
		Ω(lineScanner.Err()).ShouldNot(HaveOccurred())
		return lineScanner.Line()
	}

	assertNoMoreLines := func() {
		Ω(lineScanner.Scan()).Should(BeFalse())
		Ω(lineScanner.Err()).ShouldNot(HaveOccurred())
	}

	assertIsBlank := func(line Line) {
		Ω(line.IsBlank()).Should(BeTrue())
	}

	assertIsNotBlank := func(line Line) {
		Ω(line.IsBlank()).Should(BeFalse())
	}

	assertIsComment := func(line Line, value string) {
		Ω(line.IsComment()).Should(BeTrue())
		Ω(line.Comment()).Should(Equal(value))
	}

	assertIsNotComment := func(line Line) {
		Ω(line.IsComment()).Should(BeFalse())
	}

	assertIsCommand := func(line Line, name string) {
		Ω(line.IsCommand()).Should(BeTrue())
		Ω(line.HasCommandName(name)).Should(BeTrue())
		Ω(line.CommandName()).Should(Equal(name))
	}

	assertCommandParams := func(line Line, params ...string) {
		Ω(line.ParamCount()).Should(Equal(len(params)))
		for index, param := range params {
			Ω(line.StringParam(index)).Should(Equal(param))
		}
	}

	BeforeEach(func() {
		reader = nil
	})

	AfterEach(func() {
		if reader != nil {
			reader.Close()
		}
	})

	Describe("scanning blank lines", func() {
		var comment Line
		var blank Line
		var command Line

		BeforeEach(func() {
			openForScanning("line_scanner_blank_lines.txt")
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
		var firstComment Line
		var command Line
		var tightComment Line
		var point Line
		var firstMultiLineComment Line
		var secondMultiLineComment Line
		var normal Line
		var emptyComment Line
		var blank Line
		var object Line
		var hashComment Line
		var complicated Line
		var lastComment Line

		BeforeEach(func() {
			openForScanning("line_scanner_comments.txt")
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
			openForScanning("line_scanner_command_names.txt")
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
		var blank Line
		var noParams Line
		var stringParams Line
		var intParams Line
		var floatParams Line
		var referencesParams Line

		BeforeEach(func() {
			openForScanning("line_scanner_parameters.txt")
			noParams = readNextLine()
			blank = readNextLine()
			stringParams = readNextLine()
			intParams = readNextLine()
			floatParams = readNextLine()
			referencesParams = readNextLine()
			assertNoMoreLines()
		})

		It("can scan commands without parameters", func() {
			Ω(noParams.ParamCount()).Should(Equal(0))
		})

		It("blank lines have zero parameters", func() {
			Ω(blank.ParamCount()).Should(Equal(0))
		})

		It("can scan proper number of string parameters", func() {
			Ω(stringParams.ParamCount()).Should(Equal(4))
		})

		It("can scan proper number of int parameters", func() {
			Ω(intParams.ParamCount()).Should(Equal(2))
		})

		It("can scan proper number of float parameters", func() {
			Ω(floatParams.ParamCount()).Should(Equal(3))
		})

		It("can scan string parameters when strings", func() {
			Ω(stringParams.StringParam(0)).Should(Equal("hello"))
			Ω(stringParams.StringParam(1)).Should(Equal("complex/param"))
			Ω(stringParams.StringParam(2)).Should(Equal("\"quoted\""))
			Ω(stringParams.StringParam(3)).Should(Equal("?123?"))
		})

		It("can scan string parameters when ints", func() {
			Ω(intParams.StringParam(0)).Should(Equal("3"))
			Ω(intParams.StringParam(1)).Should(Equal("-50"))
		})

		It("can scan string parameters when floats", func() {
			Ω(floatParams.StringParam(0)).Should(Equal("1.0"))
			Ω(floatParams.StringParam(1)).Should(Equal("-13.33"))
			Ω(floatParams.StringParam(2)).Should(Equal(".3"))
		})

		It("can scan int parameters when ints", func() {
			value, err := intParams.IntParam(0)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(value).Should(Equal(int64(3)))

			value, err = intParams.IntParam(1)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(value).Should(Equal(int64(-50)))
		})

		It("cannot scan int parameters when strings", func() {
			_, err := stringParams.IntParam(0)
			Ω(err).Should(HaveOccurred())
		})

		It("cannot scan int parameters when floats", func() {
			_, err := floatParams.IntParam(0)
			Ω(err).Should(HaveOccurred())
		})

		It("can scan float parameters when floats", func() {
			value, err := floatParams.FloatParam(0)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(value).Should(BeNumerically("~", 1.0))

			value, err = floatParams.FloatParam(1)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(value).Should(BeNumerically("~", -13.33))

			value, err = floatParams.FloatParam(2)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(value).Should(BeNumerically("~", 0.3))
		})

		It("cannot scan float parameters when strings", func() {
			_, err := stringParams.FloatParam(0)
			Ω(err).Should(HaveOccurred())
		})

		It("can scan float paramters when ints", func() {
			value, err := intParams.FloatParam(0)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(value).Should(BeNumerically("~", 3))

			value, err = intParams.FloatParam(1)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(value).Should(BeNumerically("~", -50))
		})

		It("can scan reference sets when strings", func() {
			referenceSet, err := stringParams.ReferenceSetParam(0)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(referenceSet.Count()).Should(Equal(1))
			Ω(referenceSet.StringReference(0)).Should(Equal("hello"))
		})

		It("can scan reference sets when ints", func() {
			referenceSet, err := intParams.ReferenceSetParam(0)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(referenceSet.Count()).Should(Equal(1))
			Ω(referenceSet.StringReference(0)).Should(Equal("3"))
		})

		It("can scan reference sets when floats", func() {
			referenceSet, err := floatParams.ReferenceSetParam(0)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(referenceSet.Count()).Should(Equal(1))
			Ω(referenceSet.StringReference(0)).Should(Equal("1.0"))
		})

		Describe("ReferenceSet", func() {
			var singleReferenceSet ReferenceSet
			var doubleReferenceSet ReferenceSet
			var trippleReferenceSet ReferenceSet
			var skipReferenceSet ReferenceSet

			BeforeEach(func() {
				var err error
				singleReferenceSet, err = referencesParams.ReferenceSetParam(0)
				Ω(err).ShouldNot(HaveOccurred())
				doubleReferenceSet, err = referencesParams.ReferenceSetParam(1)
				Ω(err).ShouldNot(HaveOccurred())
				trippleReferenceSet, err = referencesParams.ReferenceSetParam(2)
				Ω(err).ShouldNot(HaveOccurred())
				skipReferenceSet, err = referencesParams.ReferenceSetParam(3)
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("can parse single references", func() {
				Ω(singleReferenceSet.Count()).Should(Equal(1))
				Ω(singleReferenceSet.StringReference(0)).Should(Equal("1"))
			})

			It("can parse double references", func() {
				Ω(doubleReferenceSet.Count()).Should(Equal(2))
				Ω(doubleReferenceSet.StringReference(0)).Should(Equal("2"))
				Ω(doubleReferenceSet.StringReference(1)).Should(Equal("3"))
			})

			It("can parse trippe references", func() {
				Ω(trippleReferenceSet.Count()).Should(Equal(3))
				Ω(trippleReferenceSet.StringReference(0)).Should(Equal("4"))
				Ω(trippleReferenceSet.StringReference(1)).Should(Equal("abc"))
				Ω(trippleReferenceSet.StringReference(2)).Should(Equal("6"))
			})

			It("can parse skip references", func() {
				Ω(skipReferenceSet.Count()).Should(Equal(3))
				Ω(skipReferenceSet.StringReference(0)).Should(Equal("7.3"))
				Ω(skipReferenceSet.IsBlank(0)).Should(BeFalse())
				Ω(skipReferenceSet.StringReference(1)).Should(Equal(""))
				Ω(skipReferenceSet.IsBlank(1)).Should(BeTrue())
				Ω(skipReferenceSet.StringReference(2)).Should(Equal("8"))
				Ω(skipReferenceSet.IsBlank(2)).Should(BeFalse())
			})

			It("can return int reference when int", func() {
				reference, err := singleReferenceSet.IntReference(0)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(reference).Should(Equal(int64(1)))
			})

			It("cannot return int reference when string", func() {
				_, err := trippleReferenceSet.IntReference(1)
				Ω(err).Should(HaveOccurred())
			})

			It("cannot return int reference when float", func() {
				_, err := skipReferenceSet.IntReference(0)
				Ω(err).Should(HaveOccurred())
			})

			It("can return float reference when float", func() {
				reference, err := skipReferenceSet.FloatReference(0)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(reference).Should(Equal(float64(7.3)))
			})

			It("can return float reference when int", func() {
				reference, err := singleReferenceSet.FloatReference(0)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(reference).Should(Equal(float64(1.0)))
			})

			It("cannot return float reference when string", func() {
				_, err := trippleReferenceSet.FloatReference(1)
				Ω(err).Should(HaveOccurred())
			})
		})
	})

	Describe("scanning logical lines", func() {
		var first Line
		var second Line
		var third Line

		BeforeEach(func() {
			openForScanning("line_scanner_logical_line.txt")
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
