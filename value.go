package optic

import "fmt"

// Constructor for [ValueI]
func ValI[I, T any](i I, v T) ValueI[I, T] {
	return ValueI[I, T]{
		index: i,
		value: v,
	}
}

// Constructor for [ValueE]
func ValE[T any](v T, err error) ValueE[T] {
	return ValueE[T]{
		value: v,
		err:   err,
	}
}

// Constructor for [ValueIE]
func ValIE[I, T any](i I, v T, err error) ValueIE[I, T] {
	return ValueIE[I, T]{
		index: i,
		value: v,
		err:   err,
	}
}

// Tuple of index and value. used by index optic operations.
type ValueI[I, T any] struct {
	index I
	value T
}

// Get returns the index and value as a tuple
func (iv ValueI[I, T]) Get() (I, T) {
	return iv.index, iv.value
}

func (iv ValueI[I, T]) Index() I {
	return iv.index
}

func (iv ValueI[I, T]) Value() T {
	return iv.value
}

// Lens for the [ValueI] Value field
func ValueIValue[I, T any]() Optic[Void, ValueI[I, T], ValueI[I, T], T, T, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens[ValueI[I, T], T](func(source *ValueI[I, T]) *T { return &source.value })
}

// Lens for the [ValueI] Index field
func ValueIIndex[I, T any]() Optic[Void, ValueI[I, T], ValueI[I, T], I, I, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens[ValueI[I, T], I](func(source *ValueI[I, T]) *I { return &source.index })
}

func (iv ValueI[I, T]) String() string {
	return fmt.Sprintf("%v:%v", iv.index, iv.value)
}

// Tuple of index and value. used by index optic operations.
type ValueE[T any] struct {
	value T
	err   error
}

// Get returns the index and value as a tuple
func (v ValueE[T]) Get() (T, error) {
	return v.value, v.err
}

func (iv ValueE[T]) Value() T {
	return iv.value
}

func (iv ValueE[T]) Error() error {
	return iv.err
}

// Lens for the [ValueR] Value field
func ValueEValue[T any]() Optic[Void, ValueE[T], ValueE[T], T, T, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens[ValueE[T], T](func(source *ValueE[T]) *T { return &source.value })
}

// Lens for the [ValueE] error field
func ValueEError[T any]() Optic[Void, ValueE[T], ValueE[T], error, error, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens[ValueE[T], error](func(source *ValueE[T]) *error { return &source.err })
}

func (iv ValueE[T]) String() string {
	return fmt.Sprintf("%v:%v", iv.value, iv.err)
}

// Tuple of index and value. used by [Range].
type ValueIE[I, T any] struct {
	index I
	value T
	err   error
}

// Get returns the index and value as a tuple
func (v ValueIE[I, T]) Get() (I, T, error) {
	return v.index, v.value, v.err
}

func (iv ValueIE[I, T]) Index() I {
	return iv.index
}

func (iv ValueIE[I, T]) Value() T {
	return iv.value
}

func (iv ValueIE[I, T]) Error() error {
	return iv.err
}

// Lens for the [ValueR] Value field
func ValueIEValue[I, T any]() Optic[Void, ValueIE[I, T], ValueIE[I, T], T, T, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens(func(source *ValueIE[I, T]) *T { return &source.value })
}

// Lens for the [ValueR] Value field
func ValueIEIndex[I, T any]() Optic[Void, ValueIE[I, T], ValueIE[I, T], I, I, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens(func(source *ValueIE[I, T]) *I { return &source.index })
}

// Lens for the [ValueE] error field
func ValueIEError[I, T any]() Optic[Void, ValueIE[I, T], ValueIE[I, T], error, error, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens[ValueIE[I, T], error](func(source *ValueIE[I, T]) *error { return &source.err })
}

func (iv ValueIE[I, T]) String() string {
	return fmt.Sprintf("%v:%v:%v", iv.index, iv.value, iv.err)
}
