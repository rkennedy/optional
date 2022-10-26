package optional_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/rkennedy/optional"
)

func TestCreateEmpty(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	var o Value[string]
	g.Expect(o.Present()).To(BeFalse())
}

func TestCreateFull(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	o := New("foo")
	g.Expect(o.Present()).To(BeTrue())
}

func TestGetEmpty(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	var o Value[int]
	g.Expect(o.Get()).Error().To(MatchError(ErrEmpty))
}

func TestGetFull(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	o := New(9)
	g.Expect(o.Get()).To(Equal(9))
}

func TestMustGetEmpty(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	var o Value[float64]
	g.Expect(func() { o.MustGet() }).To(PanicWith(ErrEmpty))
}

func TestMustGetFull(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	o := New(9.5)
	g.Expect(o.MustGet()).To(Equal(9.5))
}

func TestIfEmpty(t *testing.T) {
	t.Parallel()

	var o Value[any]
	o.If(func(a any) {
		t.Fail()
	})
}

func TestIfFull(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	o := New[any]('a')
	o.If(func(a any) {
		g.Expect(a).To(Equal('a'))
	})
}

func TestMarshalEmpty(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	var o Value[bool]
	b, err := o.MarshalJSON()
	g.Expect(string(b), err).To(Equal("null"))
}

func TestMarshalFull(t *testing.T) {
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
}

func TestOrElseEmpty(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	var o Value[rune]
	g.Expect(o.OrElse('r')).To(Equal('r'))
}

func TestOrElseFull(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	o := New('q')
	g.Expect(o.OrElse('r')).To(Equal('q'))
}

func TestUnmarshalNull(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	var o Value[int]
	g.Expect(o.UnmarshalJSON([]byte(`null`))).To(Succeed())
	g.Expect(o.Present()).To(BeFalse())
}

type NonerrorError struct {
	Value any
}

func (e NonerrorError) Error() string {
	return fmt.Sprintf("%#v is not an error", e.Value)
}

func TestUnmarshalInt(t *testing.T) {
	t.Parallel()

	t.Run("right type", func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		var o Value[int]
		g.Expect(o.UnmarshalJSON([]byte(`1`))).To(Succeed())
		g.Expect(o.Get()).To(Equal(1))
	})
	t.Run("wrong type", func(t *testing.T) {
		t.Parallel()
		g := NewWithT(t)

		var o Value[int]
		g.Expect(o.UnmarshalJSON([]byte(`true`))).Error().To(HaveOccurred())
		g.Expect(o.Present()).To(BeFalse())
	})
}
