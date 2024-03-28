package optional_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/rkennedy/optional"
	opt "github.com/rkennedy/optional/testing"
)

const (
	emptyCase = "empty"
	fullCase  = "full"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	t.Run(emptyCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		var o Value[string]
		g.Expect(o.Present()).To(BeFalse())
	})
	t.Run(fullCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		o := New("foo")
		g.Expect(o.Present()).To(BeTrue())
	})
}

func TestGet(t *testing.T) {
	t.Parallel()
	t.Run(emptyCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		var o Value[int]
		g.Expect(o.Get()).Error().To(MatchError(ErrEmpty))
	})
	t.Run(fullCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		o := New(9)
		g.Expect(o.Get()).To(Equal(9))
	})
}

func TestMustGet(t *testing.T) {
	t.Parallel()
	t.Run(emptyCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		var o Value[float64]
		g.Expect(func() { o.MustGet() }).To(PanicWith(ErrEmpty))
	})
	t.Run(fullCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		o := New(9.5)
		g.Expect(o.MustGet()).To(Equal(9.5))
	})
}

func TestIf(t *testing.T) {
	t.Parallel()
	t.Run(emptyCase, func(t *testing.T) {
		t.Parallel()

		var o Value[any]
		o.If(func(any) {
			t.Fail()
		})
	})
	t.Run(fullCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		o := New[any]('a')
		called := false
		o.If(func(a any) {
			g.Expect(a).To(Equal('a'))
			called = true
		})
		g.Expect(called).To(BeTrue())
	})
}

func TestMarshal(t *testing.T) {
	t.Parallel()
	t.Run(emptyCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		var o Value[bool]
		b, err := o.MarshalJSON()
		g.Expect(string(b), err).To(Equal("null"))
	})
	t.Run(fullCase, func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			Input   any
			Marshal string
		}{
			{4, `4`},
			{true, `true`},
			{nil, `null`},
			{"foo", `"foo"`},
		}
		for _, value := range cases {
			input, marshal := value.Input, value.Marshal
			t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
				t.Parallel()
				g := NewWithT(t)
				o := New(input)
				b, err := o.MarshalJSON()
				g.Expect(string(b), err).To(Equal(marshal))
			})
		}
	})
}

func TestOrElse(t *testing.T) {
	t.Parallel()
	t.Run(emptyCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		var o Value[rune]
		g.Expect(o.OrElse('r')).To(Equal('r'))
	})
	t.Run(fullCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		o := New('q')
		g.Expect(o.OrElse('r')).To(Equal('q'))
	})
}

func TestOrElseGet(t *testing.T) {
	t.Parallel()
	t.Run(emptyCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		var o Value[float64]
		var calls int
		fallback := func() float64 {
			calls++
			return 2.5
		}
		g.Expect(o.OrElseGet(fallback)).To(Equal(2.5))
		g.Expect(calls).To(Equal(1))
	})
	t.Run(fullCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		o := New(3.5)
		var calls int
		fallback := func() float64 {
			calls++
			return 2.5
		}
		g.Expect(o.OrElseGet(fallback)).To(Equal(3.5))
		g.Expect(calls).To(Equal(0))
	})
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	t.Run("null", func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		var o Value[int]
		g.Expect(o.UnmarshalJSON([]byte(`null`))).To(Succeed())
		g.Expect(o).To(opt.BeEmpty[int]())
	})
	t.Run("right type", func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		var o Value[int]
		g.Expect(o.UnmarshalJSON([]byte(`1`))).To(Succeed())
		g.Expect(o).To(opt.HaveValue(1))
	})
	t.Run("wrong type", func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		var o Value[int]
		g.Expect(o.UnmarshalJSON([]byte(`true`))).Error().To(HaveOccurred())
		g.Expect(o).To(opt.BeEmpty[int]())
	})
}

func TestGoStringFull(t *testing.T) {
	t.Parallel()
	t.Run("int", func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)
		o := New(42)
		g.Expect(fmt.Sprintf("%#v", o)).To(Equal("optional.Value[int]{42}"))
	})
	t.Run("string", func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)
		o := New("foo")
		g.Expect(fmt.Sprintf("%#v", o)).To(Equal(`optional.Value[string]{"foo"}`))
	})
}

func TestStringize(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		o := New(42)
		g.Expect(fmt.Sprintf("%v", o)).To(Equal(`42`))
	})
	t.Run("string", func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		o := New("foo")
		g.Expect(fmt.Sprintf("%v", o)).To(Equal(`foo`))
	})
	t.Run(emptyCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		var o Value[int]
		g.Expect(fmt.Sprintf("%v", o)).To(Equal(`None`))
	})
}

func TestTransform(t *testing.T) {
	t.Parallel()
	t.Run("same type", func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		o := New(4)
		g.Expect(Transform(o, func(i int) int { return i + 1 })).
			To(opt.HaveValue(5))
	})
	t.Run("different type", func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		o := New(4)
		g.Expect(Transform(o, func(i int) string { return fmt.Sprintf("%v", i) })).
			To(opt.HaveValue("4"))
	})
	t.Run(emptyCase, func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		var o Value[int]
		g.Expect(Transform(o, func(i int) int { return i + 1 })).
			To(opt.BeEmpty[int]())
	})
}

//revive:disable-next-line:cognitive-complexity
func TestTransformWithError(t *testing.T) {
	t.Parallel()
	err := errors.New("test error sentinel")

	t.Run("same type", func(t *testing.T) {
		t.Parallel()
		o := New(4)

		t.Run("without error", func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)
			g.Expect(TransformWithError(o, func(i int) (int, error) { return i + 1, nil })).
				To(opt.HaveValue(5))
		})
		t.Run("with error", func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)
			g.Expect(TransformWithError(o, func(i int) (int, error) { return i + 1, err })).
				Error().To(MatchError(err))
		})
	})
	t.Run("different type", func(t *testing.T) {
		t.Parallel()
		o := New(4)

		t.Run("without error", func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)
			g.Expect(TransformWithError(o, func(i int) (string, error) { return fmt.Sprintf("%v", i), nil })).
				To(opt.HaveValue("4"))
		})
		t.Run("with error", func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)
			g.Expect(TransformWithError(o, func(i int) (string, error) { return fmt.Sprintf("%v", i), err })).
				Error().To(MatchError(err))
		})
	})
	t.Run(emptyCase, func(t *testing.T) {
		t.Parallel()
		var o Value[int]

		t.Run("without error", func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)
			g.Expect(TransformWithError(o, func(i int) (int, error) { return i + 1, nil })).
				To(opt.BeEmpty[int]())
		})
		t.Run("with error", func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)
			g.Expect(TransformWithError(o, func(i int) (int, error) { return i + 1, err })).
				To(opt.BeEmpty[int]())
		})
	})
}
