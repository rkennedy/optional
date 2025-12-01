package testing_test

import (
	"fmt"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"

	"sweetkennedy.net/optional"
	"sweetkennedy.net/optional/testing"
)

func exampleOutput(message string, _ ...int) {
	_, _ = fmt.Println(message)
}

func ExampleHaveValueEqualing() {
	g := NewGomega(exampleOutput)

	foo := optional.New("foo")
	g.Expect(foo).To(testing.HaveValueEqualing("foo")) // passes

	bar := optional.New("bar")
	g.Expect(bar).To(testing.HaveValueEqualing("baz")) // fails

	empty := optional.Value[string]{}
	g.Expect(empty).To(testing.HaveValueEqualing("")) // fails

	// Output:
	// Expected
	//     <string>: bar
	// to equal
	//     <string>: baz
	// Expected
	//     <optional.Value[string]>: {value: nil}
	// to hold a matching value
}

func ExampleHaveValueEqualing_negated() {
	g := NewGomega(exampleOutput)

	one := optional.New("one")
	g.Expect(one).NotTo(testing.HaveValueEqualing("one")) // fails

	two := optional.New("two")
	g.Expect(two).NotTo(testing.HaveValueEqualing("three")) // passes

	empty := optional.Value[string]{}
	g.Expect(empty).NotTo(testing.HaveValueEqualing("four")) // passes

	// Output:
	// Expected
	//     <optional.Value[string]>: {value: "one"}
	// not to equal
	//     <string>: one
	// nor to be empty
}

func ExampleHaveValueMatching() {
	g := NewGomega(exampleOutput)

	red := optional.New("red")
	g.Expect(red).To(testing.HaveValueMatching(HavePrefix("r"))) // passes

	one := optional.New(1)
	g.Expect(one).To(testing.HaveValueMatching(HavePrefix("r"))) // fails

	empty := optional.Value[string]{}
	g.Expect(empty).To(testing.HaveValueMatching(gstruct.Ignore())) // fails

	// Output:
	// HavePrefix matcher requires a string or stringer.  Got:
	//     <int>: 1
	// Expected
	//     <optional.Value[string]>: {value: nil}
	// to hold a matching value
}

func ExampleHaveValueMatching_negated() {
	g := NewGomega(exampleOutput)

	blue := optional.New("blue")
	g.Expect(blue).NotTo(testing.HaveValueMatching(HaveLen(1))) // passes

	green := optional.New("green")
	g.Expect(green).NotTo(testing.HaveValueMatching(HavePrefix("g"))) // fails

	empty := optional.Value[string]{}
	g.Expect(empty).NotTo(testing.HaveValueMatching(HaveLen(1))) // passes

	// Output:
	// Expected
	//     <optional.Value[string]>: {value: "green"}
	// not to have prefix
	//     <string>: g
	// nor to be empty
}
