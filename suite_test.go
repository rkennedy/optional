package optional_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOptionalSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Optional")
}
