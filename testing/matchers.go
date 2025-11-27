// Package testing provides Gomega matchers for use with [optional.Value] values.
package testing

import (
	"errors"
	"fmt"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"

	"sweetkennedy.net/optional"
)

// BeEmpty asserts that the tested value is an empty [optional.Value].
func BeEmpty() types.GomegaMatcher {
	return &emptyMatcher{}
}

type emptyMatcher struct{}

func (*emptyMatcher) Match(actual any) (success bool, err error) {
	opt, ok := actual.(optional.AnyGetter)
	if !ok {
		// actual isn't an optional.Value
		return false, fmt.Errorf("BeEmpty matcher expects an optional.Value. Got:\n%s",
			format.Object(actual, 1))
	}
	return !opt.Present(), nil
}

func (*emptyMatcher) FailureMessage(actual any) (message string) {
	return format.Message(actual, "not to have a value")
}

func (*emptyMatcher) NegatedFailureMessage(actual any) (message string) {
	return format.Message(actual, "to have a value")
}

// HaveValueMatching checks whether an [optional.Value] holds a value that matches the given matcher.
//
// Be careful when negating this matcher. Doing so will cause it to pass either when the value is empty or when the
// wrapped matcher fails. For example:
//
//	// !!! Assertion passes when len(v.Get()) != 3 _or_ when !v.Present
//	Expect(v).NotTo(HaveValueMatching(HaveLen(3)))
//
// If you want to check for an empty value, then use [BeEmpty]:
//
//	Expect(v).To(BeEmpty())
//
// If you want to check that a value is present that doesn't match, then negate the wrapped matcher:
//
//	Expect(v).To(HaveValueMatching(Not(HaveLen(3))))
func HaveValueMatching(matcher types.GomegaMatcher) types.GomegaMatcher {
	return &optionalValueMatcher{matcher}
}

type optionalValueMatcher struct {
	Matcher types.GomegaMatcher
}

func (m *optionalValueMatcher) Match(actual any) (success bool, err error) {
	opt, ok := actual.(optional.AnyGetter)
	if !ok {
		// actual isn't an optional.Value
		return false, fmt.Errorf("HaveValueMatching matcher expects an optional.Value. Got:\n%s",
			format.Object(actual, 1))
	}
	val, err := opt.GetAny()
	if err != nil {
		// actual is empty
		return false, nil
	}
	return m.Matcher.Match(val)
}

func (m *optionalValueMatcher) FailureMessage(actual any) (message string) {
	opt, ok := actual.(optional.AnyGetter)
	if !ok {
		panic("Match should have failed.")
	}
	val, err := opt.GetAny()
	if err != nil {
		if !errors.Is(err, optional.ErrEmpty) {
			// The only thing Get can ever return is ErrEmpty, so if we get here, we're in trouble.
			panic(err)
		}
		return format.Message(actual, "to hold a matching value")
	}
	return m.Matcher.FailureMessage(val)
}

func (m *optionalValueMatcher) NegatedFailureMessage(actual any) (message string) {
	return m.Matcher.NegatedFailureMessage(actual) + "\nnor to be empty"
}

// HaveValueEqualing checks whether an [optional.Value] holds a value equal to the given value. It's syntactic sugar for
// [HaveValueMatching]([gomega.Equal](arg)).
func HaveValueEqualing(arg any) types.GomegaMatcher {
	return HaveValueMatching(gomega.Equal(arg))
}
