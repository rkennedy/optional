package optional

import (
	"encoding/json"
	"errors"
)

// ErrEmpty indicates that an optional value was empty when its value was requested.
var ErrEmpty = errors.New("value not present")

type Value[T any] struct {
	value *T
}

func New[T any](v T) Value[T] {
	return Value[T]{
		value: &v,
	}
}

func (o Value[T]) Get() (T, error) {
	if !o.Present() {
		var zero T
		return zero, ErrEmpty
	}
	return *o.value, nil
}

func (o Value[T]) MustGet() T {
	if !o.Present() {
		panic(ErrEmpty)
	}
	return *o.value
}

func (o Value[T]) If(fn func(T)) {
	if o.Present() {
		fn(*o.value)
	}
}

func (o Value[T]) MarshalJSON() ([]byte, error) {
	if o.Present() {
		return json.Marshal(o.value)
	}
	return json.Marshal(nil)
}

func (o Value[T]) OrElse(v T) T {
	if o.Present() {
		return *o.value
	}
	return v
}

func (o Value[T]) Present() bool {
	return o.value != nil
}

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
