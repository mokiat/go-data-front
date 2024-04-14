package mtl_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMTL(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MTL Scanner Suite")
}
