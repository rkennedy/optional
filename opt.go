// Package optional provides an “option type,” which can either be empty or
// hold a value. In that respect, it's very similar to an ordinary pointer
// type, except it has methods that make its possible emptiness more explicit.
package optional

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrEmpty indicates that an optional value was empty when its value was requested.
var ErrEmpty = errors.New("value not present")

// Value is a type that may or may not hold a value. Its interface is modeled
// on Java's [java.util.Optional] type and C++'s [std::optional].
//
// Value implements the following interfaces:
//   - [fmt.GoStringer]
//   - [fmt.Stringer]
//   - [json.Marshaler]
//   - [json.Unmarshaler]
//
// [java.util.Optional]: https://docs.oracle.com/javase/8/docs/api/java/util/Optional.html
// [std::optional]: https://en.cppreference.com/w/cpp/utility/optional
type Value[T any] []T

var _ fmt.GoStringer = &Value[any]{}
var _ fmt.Stringer = &Value[any]{}
var _ json.Marshaler = &Value[any]{}
var _ json.Unmarshaler = &Value[any]{}

// New creates a new Value holding the given value.
func New[T any](v T) Value[T] {
	return []T{v}
}

// Get returns the current value, if there is one. If the Value is empty, then
// ErrEmpty is returned, and the value result is unspecified.
func (o Value[T]) Get() (result T, err error) {
	if !o.Present() {
		return result, ErrEmpty
	}
	return o[0], nil
}

// MustGet returns the current value, if there is one. If the Value is empty,
// then MustGet panics with ErrEmpty.
func (o Value[T]) MustGet() T {
	if !o.Present() {
		panic(ErrEmpty)
	}
	return o[0]
}

// If calls the given function if the Value holds a value. If Value is empty,
// then If is a no-op.
func (o Value[T]) If(fn func(T)) {
	if o.Present() {
		fn(o[0])
	}
}

// MarshalJSON converts the value to a JSON value. If the Value is empty, then
// the JSON result is null.
func (o Value[T]) MarshalJSON() ([]byte, error) {
	if o.Present() {
		return json.Marshal(o[0])
	}
	return json.Marshal(nil)
}

// OrElse returns the stored value, if there is one. If the Value is empty,
// then OrElse returns the given argument.
func (o Value[T]) OrElse(v T) T {
	if o.Present() {
		return o[0]
	}
	return v
}

// OrElseGet returns the stored value, if there is one. If the Value is empty,
// then OrElseGet calls the given function and returns the result. Use this
// instead of OrElse when calculation of the fallback value is relatively
// expensive. It will only be calculated when needed.
func (o Value[T]) OrElseGet(calculateFallback func() T) T {
	if o.Present() {
		return o[0]
	}
	return calculateFallback()
}

// Present returns true if there is a value stored, false if the Value is empty.
func (o Value[_]) Present() bool {
	return len(o) != 0
}

// UnmarshalJSON converts the given JSON value to an optional Value[T]. If the
// JSON value is null, then the result is empty. Otherwise, the JSON is
// unmarshaled in the same way values of type T are unmarshaled.
func (o *Value[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*o = []T{}
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*o = []T{value}
	return nil
}

// GoString fornats the Value as Go code, providing an implementation for the
// %#v format string.
func (o Value[T]) GoString() string {
	if o.Present() {
		return fmt.Sprintf("%T{%#v}", o, o[0])
	}
	return fmt.Sprintf("%T{}", o)
}

// String returns the string representation of the stored value, if present.
// Otherwise, it returns None.
func (o Value[T]) String() string {
	if o.Present() {
		return fmt.Sprintf("%v", o[0])
	}
	return "None"
}

// Transform applies the given function to the optional value if the input
// value is non-empty, and returns a new optional of the corresponding return
// type holding the returned value. Returns an empty value if the input is
// empty.
func Transform[T, U any](in Value[T], fn func(T) U) Value[U] {
	if in.Present() {
		return New(fn(in[0]))
	}
	return Value[U]{}
}

// TransformWithError applies the given function to the optional value if the
// input value is non-empty, and returns a new optional of the corresponding
// return type holding the returned value. Returns an empty value if the input
// is empty. If the transform function returns an error, then an empty value
// and that error are returned.
func TransformWithError[T, U any](in Value[T], fn func(T) (U, error)) (result Value[U], err error) {
	in.If(func(val T) {
		var newVal U
		newVal, err = fn(val)
		if err == nil {
			result = []U{newVal}
		}
	})
	return result, err
}
