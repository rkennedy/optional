package optional_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/rkennedy/optional"
	opt "github.com/rkennedy/optional/testing"
)

var _ = Describe("Optional", func() {
	It("can be created empty", func() {
		var o Value[string]
		Expect(o).To(opt.BeEmpty())
	})

	It("can be created full", func() {
		o := New("foo")
		Expect(o).NotTo(opt.BeEmpty())
	})

	Context("when empty", func() {
		It("Get returns error", func() {
			var o Value[int]
			Expect(o.Get()).Error().To(MatchError(ErrEmpty))
		})

		It("MustGet panics", func() {
			var o Value[float64]
			Expect(func() { o.MustGet() }).To(PanicWith(ErrEmpty))
		})

		It("If isn't called", func() {
			var o Value[any]
			o.If(func(any) {
				Fail("If was called on an empty Optional")
			})
		})

		It("Marshall emits null", func() {
			var o Value[bool]
			b, err := o.MarshalJSON()
			Expect(string(b), err).To(Equal("null"))
		})

		It("OrElse returns the alternative value", func() {
			var o Value[rune]
			Expect(o.OrElse('r')).To(Equal('r'))
		})

		It("OrElseGet returns the value from the function", func() {
			var o Value[float64]
			var calls int
			fallback := func() float64 {
				calls++
				return 2.5
			}
			Expect(o.OrElseGet(fallback)).To(Equal(2.5))
			Expect(calls).To(Equal(1))
		})

		Context("formatted with Sprintf", func() {
			It("formats with %#v", func() {
				var o Value[int]
				Expect(fmt.Sprintf("%#v", o)).To(Equal(`optional.Value[int]{}`))
			})

			It("formats with %v", func() {
				var o Value[int]
				Expect(fmt.Sprintf("%v", o)).To(Equal(`None`))
			})
		})

		It("transforms to empty", func() {
			var o Value[int]
			Expect(Transform(o, func(i int) int { return i + 1 })).
				To(opt.BeEmpty())
		})
	})

	Context("when full", func() {
		It("Get returns the item", func() {
			o := New(9)
			Expect(o.Get()).To(Equal(9))
		})

		It("MustGet returns the item", func() {
			o := New(9.5)
			Expect(o.MustGet()).To(Equal(9.5))
		})

		It("If is called with the item", func() {
			o := New[any]('a')
			called := false
			o.If(func(a any) {
				Expect(a).To(Equal('a'))
				called = true
			})
			Expect(called).To(BeTrue())
		})

		DescribeTable("Marshal emits JSON values",
			func(value any, expected string) {
				o := New(value)
				b, err := o.MarshalJSON()
				Expect(string(b), err).To(Equal(expected))
			},
			Entry("for int", 4, `4`),
			Entry("for bool", true, `true`),
			Entry("for nil", nil, `null`),
			Entry("for string", "foo", `"foo"`),
		)

		It("OrElse returns the item", func() {
			o := New('q')
			Expect(o.OrElse('r')).To(Equal('q'))
		})

		It("OrElseGet doesn't use the function", func() {
			o := New(3.5)
			var calls int
			fallback := func() float64 {
				calls++
				return 2.5
			}
			Expect(o.OrElseGet(fallback)).To(Equal(3.5))
			Expect(calls).To(Equal(0))
		})

		Context("stringizes with Sprintf %#v", func() {
			It("produces an int", func() {
				o := New(42)
				Expect(fmt.Sprintf("%#v", o)).To(Equal(`optional.Value[int]{42}`))
			})

			It("produces a string", func() {
				o := New("foo")
				Expect(fmt.Sprintf("%#v", o)).To(Equal(`optional.Value[string]{"foo"}`))
			})
		})

		Context("stringizes with Sprintf %v", func() {
			It("produces an int", func() {
				o := New(42)
				Expect(fmt.Sprintf("%v", o)).To(Equal(`42`))
			})

			It("produces a string", func() {
				o := New("foo")
				Expect(fmt.Sprintf("%v", o)).To(Equal(`foo`))
			})
		})

		Context("transforms", func() {
			It("to the same type", func() {
				o := New(4)
				Expect(Transform(o, func(i int) int { return i + 1 })).
					To(opt.HaveValueEqualing(5))
			})

			It("to a different type", func() {
				o := New(4)
				Expect(Transform(o, func(i int) string { return fmt.Sprintf("%v", i) })).
					To(opt.HaveValueEqualing("4"))
			})
		})
	})

	Context("when unmarshalled from JSON", func() {
		It("is empty on null", func() {
			var o Value[int]
			Expect(o.UnmarshalJSON([]byte(`null`))).To(Succeed())
			Expect(o).To(opt.BeEmpty())
		})

		It("is full when the types match", func() {
			var o Value[int]
			Expect(o.UnmarshalJSON([]byte(`1`))).To(Succeed())
			Expect(o).To(opt.HaveValueEqualing(1))
		})

		It("is empty when types mismatch", func() {
			var o Value[int]
			Expect(o.UnmarshalJSON([]byte(`true`))).Error().To(HaveOccurred())
			Expect(o).To(opt.BeEmpty())
		})
	})

	Context("TransformWithError", func() {
		err := errors.New("test error sentinel")

		Context("to the same type", func() {
			o := New(4)

			It("succeeds", func() {
				Expect(TransformWithError(o, func(i int) (int, error) { return i + 1, nil })).
					To(opt.HaveValueEqualing(5))
			})

			It("reports an error", func() {
				Expect(TransformWithError(o, func(i int) (int, error) { return i + 1, err })).
					Error().To(MatchError(err))
			})
		})

		Context("to a different type", func() {
			o := New(4)

			It("succeeds", func() {
				Expect(TransformWithError(o, func(i int) (string, error) { return fmt.Sprintf("%v", i), nil })).
					To(opt.HaveValueEqualing("4"))
			})

			It("reports an error", func() {
				Expect(TransformWithError(o, func(i int) (string, error) { return fmt.Sprintf("%v", i), err })).
					Error().To(MatchError(err))
			})
		})

		Context("from an empty value", func() {
			var o Value[int]

			It("returns an empty value", func() {
				Expect(TransformWithError(o, func(i int) (int, error) { return i + 1, nil })).
					To(opt.BeEmpty())
			})

			It("returns an empty value when there's an error", func() {
				Expect(TransformWithError(o, func(i int) (int, error) { return i + 1, err })).
					To(opt.BeEmpty())
			})
		})
	})
})
