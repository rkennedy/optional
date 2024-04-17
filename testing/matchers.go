// Package testing provides Gomega matchers for use with [optional.Value] values.
package testing

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"

	"github.com/rkennedy/optional"
)

// EmptyMatcher is a Gomega matcher that checks whether an [optional.Value] is empty.
type EmptyMatcher struct {
}

var _ types.GomegaMatcher = &EmptyMatcher{}

// Match implements [types.GomegaMatcher]'s Match function.
func (*EmptyMatcher) Match(actual any) (bool, error) {
	opt, ok := actual.(optional.ValueHaver)
	if !ok {
		return false, fmt.Errorf("to have type %s", reflect.TypeFor[optional.ValueHaver]().String())
	}
	return !opt.Present(), nil
}

// FailureMessage implements [types.GomegaMatcher]'s FailureMessage function.
func (*EmptyMatcher) FailureMessage(actual any) string {
	_, ok := actual.(optional.ValueHaver)
	if !ok {
		return format.Message(actual, "to be an optional.Value", fmt.Sprintf("got type %T", actual))
	}
	return format.Message(actual, "not to hold a value")
}

// NegatedFailureMessage implements [types.GomegaMatcher]'s NegatedFailureMessage function.
func (*EmptyMatcher) NegatedFailureMessage(actual any) string {
	_, ok := actual.(optional.ValueHaver)
	if !ok {
		return format.Message(actual, "to be an optional.Value", fmt.Sprintf("got type %T", actual))
	}
	return format.Message(actual, "to hold a value")
}

// BeEmpty asserts that the tested value is an empty [optional.Value] with type T.
func BeEmpty() types.GomegaMatcher {
	return &EmptyMatcher{}
}

func get[T any](arg optional.Value[T]) (T, error) {
	return arg.Get()
}

// HaveValueMatching checks whether an [optional.Value] holds a value that matches the given matcher. Be careful when
// negating this matcher. You probably don't want to negate this matcher; doing so will cause it to pass either when the
// value is empty or when the wrapped matcher fails.
//
//	Expect(v).NotTo(HaveValueMatching[string](HaveLen(3))) // !!! Fails when !v.Present _or_ when len(v.Get()) != 3
//
// If you want to check for an empty value, then use [BeEmpty]:
//
//	Expect(v).To(BeEmpty())
//
// If you want to check that a value is present that doesn't match, then negate the wrapped matcher:
//
//	Expect(v).To(HaveValueMatching[string](Not(HaveLen(3))))
func HaveValueMatching[T any](matcher types.GomegaMatcher) types.GomegaMatcher {
	return &matchers.AndMatcher{
		Matchers: []types.GomegaMatcher{
			&matchers.NotMatcher{
				Matcher: &EmptyMatcher{},
			},
			matchers.NewWithTransformMatcher(get[T], matcher),
		},
	}
}

// HaveValueEqualing checks wheter an [optional.Value] holds a value equal to the given value. Be careful when negating
// this matcher.
//
//	Expect(v).NotTo(HaveValueEqualing(3)) // !!! Fails when !v.Present _or_ when v.Get() != 3
//
// More likely, you want to check that the [optional.Value] contains a value, and that the contained value is not equal
// to something. Use [HaveValueMatching] for that and negate the wrapped matcher. For example:
//
//	Expect(v).To(HaveValueMatching[int](Not(Equal(3))))
func HaveValueEqualing[T any](arg T) types.GomegaMatcher {
	return &matchers.AndMatcher{
		Matchers: []types.GomegaMatcher{
			&matchers.NotMatcher{
				Matcher: &EmptyMatcher{},
			},
			matchers.NewWithTransformMatcher(get[T], &matchers.EqualMatcher{Expected: arg}),
		},
	}
}
