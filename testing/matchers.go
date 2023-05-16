// Package testing provides Gomega matchers for use with [optional.Value] values.
package testing

import (
	"reflect"

	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"

	"github.com/rkennedy/optional"
)

// BeEmpty asserts that the tested value is an empty [optional.Value] with type T.
func BeEmpty[T any]() types.GomegaMatcher {
	return gcustom.MakeMatcher(func(opt optional.Value[T]) (bool, error) {
		return !opt.Present(), nil
	}).WithMessage("be empty")
}

// HaveValue asserts that the tested value is a non-empty [optional.Value] holding a value equal to the argument.
func HaveValue[T any](arg T) types.GomegaMatcher {
	return gcustom.
		MakeMatcher(func(opt optional.Value[T]) (bool, error) {
			val, err := opt.Get()
			if err != nil {
				return false, nil
			}
			return reflect.DeepEqual(val, arg), nil
		}).
		WithTemplate("Expected:\n{{.FormattedActual}}\n{{.To}} have a value equal to:\n{{format .Data 1}}").
		WithTemplateData(arg)
}
