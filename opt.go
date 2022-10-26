// Package optional provides an _option type_, which can either be empty or
// hold a value. In that respect, it's very similar to an ordinary pointer
// type, except it has methods that make its possible emptiness more explicit.
package optional

import (
	"encoding/json"
	"errors"
)

// ErrEmpty indicates that an optional value was empty when its value was requested.
var ErrEmpty = errors.New("value not present")

// Value is a type that may or may not hold a value. Its interface is modeled
// on Java's java.util.Optional type and C++'s std::optional.
type Value[T any] struct {
	value *T
}

// New creates a new Value holding the given value.
func New[T any](v T) Value[T] {
	return Value[T]{
		value: &v,
	}
}

// Get returns the current value, if there is one. If the Value is empty, then
// ErrEmpty is returned, and the value result is unspecified.
func (o Value[T]) Get() (result T, err error) {
	if !o.Present() {
		return result, ErrEmpty
	}
	return *o.value, nil
}

// MustGet returns the current value, if there is one. If the Value is empty,
// then MustGet panics with ErrEmpty.
func (o Value[T]) MustGet() T {
	if !o.Present() {
		panic(ErrEmpty)
	}
	return *o.value
}

// If calls the given function if the Value holds a value. If Value is empty,
// then If is a no-op.
func (o Value[T]) If(fn func(T)) {
	if o.Present() {
		fn(*o.value)
	}
}

// MarshalJSON converts the value to a JSON value. If the Value is empty, then
// the JSON result is null.
func (o Value[T]) MarshalJSON() ([]byte, error) {
	if o.Present() {
		return json.Marshal(o.value)
	}
	return json.Marshal(nil)
}

// OrElse returns the stored value, if there is one. If the Value is empty,
// then OrElse returns the given argument.
func (o Value[T]) OrElse(v T) T {
	if o.Present() {
		return *o.value
	}
	return v
}

// Present returns true if there is a value stored, false if the Value is empty.
func (o Value[T]) Present() bool {
	return o.value != nil
}

// UnmarshalJSON converts the given JSON value to an optional Value[T]. If the
// JSON value is null, then the result is empty. Otherwise, the JSON is
// unmarshaled in the same way values of type T are unmarshaled.
func (o *Value[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.value = nil
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = &value
	return nil
}
