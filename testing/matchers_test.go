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
		DescribeTable("matches empty values",
			func(o any) {
				Expect(o).To(opt.BeEmpty())
			},
			Entry("empty int", optional.Value[int]{}),
			Entry("empty string", optional.Value[string]{}),
		)

		It("produces a message when not empty", func() {
			vi := optional.New(1)
			m := opt.BeEmpty()
			Expect(m.Match(vi)).To(BeFalse())
			Expect(m.FailureMessage(vi)).To(Equal(
				"Expected\n    <optional.Value[int]>: {value: 1}\nnot to have a value"))
		})

		When("negated", func() {
			DescribeTable("matches non-empty values",
				func(o any) {
					Expect(o).NotTo(opt.BeEmpty())
				},
				Entry("int", optional.New(1)),
				Entry("string", optional.New("s")),
			)

			It("produces a message when empty", func() {
				vi := optional.Value[int]{}
				m := Not(opt.BeEmpty())
				Expect(m.Match(vi)).To(BeFalse())
				Expect(m.FailureMessage(vi)).To(Equal(
					"Expected\n    <optional.Value[int]>: {value: nil}\nto have a value"))
			})
		})
	})

	Context("HaveValueMatching", func() {
		DescribeTable("matches values",
			func(o any, match types.GomegaMatcher) {
				Expect(o).To(opt.HaveValueMatching(match))
			},
			Entry("int", optional.New(1), Equal(1)),
			Entry("string", optional.New("s"), Equal("s")),
			Entry("string matcher", optional.New("s"), HaveLen(1)),
		)

		When("negated", func() {
			DescribeTable("matches values",
				func(o any, match types.GomegaMatcher) {
					Expect(o).NotTo(opt.HaveValueMatching(match))
				},
				Entry("empty int", optional.Value[int]{}, Equal(1)),
				Entry("non-matching value", optional.New("foo"), Equal("food")),
			)
		})

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

			m := opt.HaveValueMatching(gstruct.Ignore())
			Expect(m.Match(o)).To(BeFalse())
			Expect(m.FailureMessage(o)).To(HavePrefix(
				"Expected\n    <optional.Value[int]>: {value: nil}\nto hold a matching value"))
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
