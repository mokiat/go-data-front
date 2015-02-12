package mtl_test_help

import (
	"fmt"
	"io"
	"os"

	"github.com/momchil-atanasov/go-data-front/common"
	"github.com/momchil-atanasov/go-data-front/mtl"
	"github.com/momchil-atanasov/go-data-front/mtl/mtl_test_fake"
	. "github.com/onsi/gomega"
)

type ScannerFixture struct {
	scanner         common.Scanner
	handler         *mtl_test_fake.FakeScannerHandler
	scanErr         error
	commentCounter  int
	materialCounter int
}

func NewScannerFixture() *ScannerFixture {
	handler := new(mtl_test_fake.FakeScannerHandler)
	return &ScannerFixture{
		scanner: mtl.NewScanner(handler),
		handler: handler,
	}
}

func (f *ScannerFixture) Scan(reader io.Reader) {
	f.scanErr = f.scanner.Scan(reader)
}

func (f *ScannerFixture) ScanFile(filename string) {
	file, err := os.Open(fmt.Sprintf("mtl_test_res/%s", filename))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	f.Scan(file)
}

func (f *ScannerFixture) AssertScannerReturnedError(expectedErr error) {
	Ω(f.scanErr).Should(HaveOccurred())
	Ω(f.scanErr).Should(Equal(expectedErr))
}

func (f *ScannerFixture) AssertScannerDidNotReturnError() {
	Ω(f.scanErr).ShouldNot(HaveOccurred())
}

func (f *ScannerFixture) Handler() *mtl_test_fake.FakeScannerHandler {
	return f.handler
}

func (f *ScannerFixture) AssertCommentCall(expectedComment string) {
	Ω(f.handler.OnCommentCallCount()).Should(BeNumerically(">", f.commentCounter))
	Ω(f.handler.OnCommentArgsForCall(f.commentCounter)).Should(Equal(expectedComment))
	f.commentCounter++
}

func (f *ScannerFixture) AssertNoMoreCommentCalls() {
	Ω(f.handler.OnCommentCallCount()).Should(Equal(f.commentCounter))
}

func (f *ScannerFixture) AssertMaterialCall(expectedName string) {
	Ω(f.handler.OnMaterialCallCount()).Should(BeNumerically(">", f.materialCounter))
	Ω(f.handler.OnMaterialArgsForCall(f.materialCounter)).Should(Equal(expectedName))
	f.materialCounter++
}

func (f *ScannerFixture) AssertNoMoreMaterialCalls() {
	Ω(f.handler.OnMaterialCallCount()).Should(Equal(f.materialCounter))
}
