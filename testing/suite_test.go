package testing_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestingSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Optional testing")
}
