package obj_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestObj(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Obj Suite")
}
