package testing_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"

	"sweetkennedy.net/optional"
	opt "sweetkennedy.net/optional/testing"
)

var _ = Describe("Optional matchers", func() {
	Context("BeEmpty", func() {
		DescribeTable("matches values",
			func(o any, match types.GomegaMatcher) {
				Expect(o).To(match)
			},
			Entry("empty int", optional.Value[int]{}, opt.BeEmpty()),
			Entry("empty string", optional.Value[string]{}, opt.BeEmpty()),
			Entry("int", optional.New(1), Not(opt.BeEmpty())),
			Entry("string", optional.New("s"), Not(opt.BeEmpty())),
		)

		It("produces a message when not empty", func() {
			vi := optional.New(1)
			m := opt.BeEmpty()
			Expect(m.Match(vi)).To(BeFalse())
			Expect(m.FailureMessage(vi)).To(Equal(
				"Expected\n    <optional.Value[int] | len:1, cap:1>: [1]\nto be empty"))
		})

		It("produces a message when empty", func() {
			vi := optional.Value[int]{}
			m := Not(opt.BeEmpty())
			Expect(m.Match(vi)).To(BeFalse())
			Expect(m.FailureMessage(vi)).To(Equal(
				"Expected\n    <optional.Value[int] | len:0, cap:0>: []\nnot to be empty"))
		})
	})

	Context("HaveValue", func() {
		DescribeTable("matches values",
			func(o any, match types.GomegaMatcher) {
				Expect(o).To(match)
			},
			Entry("int", optional.New(1), opt.HaveValueEqualing(1)),
			Entry("string", optional.New("s"), opt.HaveValueEqualing("s")),
			Entry("string matcher", optional.New("s"), opt.HaveValueMatching(HaveLen(1))),
			Entry("except for empty", optional.Value[int]{}, Not(opt.HaveValueEqualing(1))),
		)

		It("has same type expectation has its matcher", func() {
			o := optional.New(1)

			m := opt.HaveValueMatching(HaveLen(1))
			Expect(m.Match(o)).Error().To(MatchError(MatchRegexp(`^HaveLen matcher expects a`)))
		})

		It("mismatches on non-optional type", func() {
			vi := int(1)

			m := opt.HaveValueMatching(gstruct.Ignore())
			Expect(m.Match(vi)).Error().To(MatchError(MatchRegexp(`expects an optional.Value`)))
		})

		It("produces a message on empty mismatch", func() {
			o := optional.Value[int]{}

			m := opt.HaveValueEqualing(1)
			Expect(m.Match(o)).To(BeFalse())
			Expect(m.FailureMessage(o)).To(HavePrefix(
				"Expected\n    <optional.Value[int] | len:0, cap:0>: []\nto hold a matching value"))
		})

		It("produces matcher's message on mismatch", func() {
			o := optional.New("s")

			submatch := HaveLen(2)
			submatchMessage := submatch.FailureMessage(o.MustGet())

			m := opt.HaveValueMatching(submatch)
			Expect(m.Match(o)).To(BeFalse())

			Expect(m.FailureMessage(o)).To(Equal(submatchMessage))
		})
	})
})
