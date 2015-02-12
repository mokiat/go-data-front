package mtl_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMtl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mtl Suite")
}
