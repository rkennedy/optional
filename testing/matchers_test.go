package testing_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/rkennedy/optional"
	opt "github.com/rkennedy/optional/testing"
)

func TestEmpty(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	vi := optional.Value[int]{}
	vs := optional.Value[string]{}

	g.Expect(vi).To(opt.BeEmpty())
	g.Expect(vs).To(opt.BeEmpty())
}

func TestNotEmpty(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	vi := optional.New(1)
	vs := optional.New("s")

	g.Expect(vi).NotTo(opt.BeEmpty())
	g.Expect(vs).NotTo(opt.BeEmpty())
}

func TestEmptyFailureMessage(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	vi := optional.New(1)
	m := opt.BeEmpty()
	g.Expect(m.Match(vi)).To(BeFalse())
	g.Expect(m.FailureMessage(vi)).To(Equal("Expected\n    <optional.Value[int] | len:1, cap:1>: [1]\nto be empty"))
}

func TestEmptyUnexpectedFailureMessage(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	vi := optional.Value[int]{}
	m := Not(opt.BeEmpty())
	g.Expect(m.Match(vi)).To(BeFalse())
	g.Expect(m.FailureMessage(vi)).To(Equal("Expected\n    <optional.Value[int] | len:0, cap:0>: []\nnot to be empty"))
}

func TestHaveValue(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	vi := optional.New(1)
	vs := optional.New("s")

	g.Expect(vi).To(opt.HaveValueEqualing(1))
	g.Expect(vs).To(opt.HaveValueEqualing("s"))
	g.Expect(vs).To(opt.HaveValueMatching[string](HaveLen(1)))
}

func TestNotHaveValueEqualing(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	vi := optional.Value[int]{}

	g.Expect(vi).NotTo(opt.HaveValueEqualing(1))
}

func TestWrongType(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	vi := optional.New(1)

	m := opt.HaveValueMatching[string](gstruct.Ignore())
	g.Expect(m.Match(vi)).Error().To(MatchError(MatchRegexp(`^Transform function expects.*but we have`)))
}

func TestEmptyHaveValue(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	vi := optional.Value[int]{}

	m := opt.HaveValueEqualing(1)
	g.Expect(m.Match(vi)).To(BeFalse())
	g.Expect(m.FailureMessage(vi)).To(Equal("Expected\n    <optional.Value[int] | len:0, cap:0>: []\nnot to be empty"))
}

func TestHaveValueFailureMessage(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	vs := optional.New("s")
	m := opt.HaveValueMatching[string](HaveLen(2))
	g.Expect(m.Match(vs)).To(BeFalse())
	g.Expect(m.FailureMessage(vs)).To(Equal("Expected\n    <string>: s\nto have length 2"))
}
