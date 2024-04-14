package obj_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOBJ(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OBJ Decoder Suite")
}
