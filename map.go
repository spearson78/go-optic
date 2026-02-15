package optic

import (
	"cmp"
	"context"
	"errors"
	"iter"
	"reflect"
	"sort"
	"unsafe"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

// TraverseMap returns a [Traversal] that focuses the elements of a map. The key of the map element is used as the index.
//
// Note: TraverseMap always iterates in key order.
//
// See: [TraverseMapP] for a polymorphic version
func TraverseMap[K comparable, A any]() Optic[K, map[K]A, map[K]A, A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseMapP[K, A, A]()
}

func unsortedMapSeqF[I comparable, A any](source map[I]A) func(yield func(val ValueIE[I, A]) bool) {
	return func(yield func(val ValueIE[I, A]) bool) {
		for k, v := range source {
			if !yield(ValIE(k, v, nil)) {
				break
			}
		}
	}

}

// TraverseMapP returns a polymorphic [Traversal] that focuses the elements of a map. The key of the map element is used as the index.
//
// Note: TraverseMapP always iterates in key order.
//
// See: [TraverseMap] for a non polymorphic version
func TraverseMapP[K comparable, A any, B any]() Optic[K, map[K]A, map[K]B, A, B, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, K, map[K]A, map[K]B, A, B](
		func(ctx context.Context, source map[K]A) SeqIE[K, A] {
			return func(yield func(ValueIE[K, A]) bool) {

				unsortedSeq := unsortedMapSeqF(source)
				sortedCol, _, _, err := heapSort(ctx, unsortedSeq, func(ctx context.Context, a, b ValueI[K, A]) (bool, error) {
					return mapKeyCompare(a.index, b.index) < 0, nil
				})

				err = JoinCtxErr(ctx, err)
				if err != nil {
					var k K
					var a A
					yield(ValIE(k, a, err))
					return
				}

				sortedCol(yield)
			}
		},
		func(ctx context.Context, source map[K]A) (int, error) {
			return len(source), nil
		},
		func(ctx context.Context, fmap func(index K, focus A) (B, error), source map[K]A) (map[K]B, error) {
			var sorted []ValueI[K, A]
			for i, a := range source {
				sorted = append(sorted, ValI(i, a))
			}

			//Map iteration is sometimes done in reverse to prevent dependencies on random map ordering
			//We need to work with a dependable order to ensure that operations like PartsOf work consistently
			sort.Slice(sorted, func(i, j int) bool {
				return mapKeyCompare(sorted[i].index, sorted[j].index) < 0
			})

			ret := make(map[K]B, len(source))

			for _, a := range sorted {
				res, err := fmap(a.index, a.value)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					return nil, err
				}
				ret[a.index] = res
			}
			return ret, ctx.Err()
		},
		func(ctx context.Context, index K, source map[K]A) SeqIE[K, A] {
			return func(yield func(ValueIE[K, A]) bool) {
				r, ok := source[index]
				if ok {
					yield(ValIE(index, r, nil))
				}
			}

		},
		IxMatchComparable[K](),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Traverse{
				OpticTypeExpr: ot,
			}
		}),
	)
}

// AtMap returns a [Lens] that focuses the index slot in a map.
// When viewing the none option means no element was present in the map with the given index.
// When modifying the none option will remove the element from the map.
//
// See:
//   - [AtMapT2] for a version that can use a dynamic index.
//   - [Element] for an optic that focuses by index position.
//   - [IndexMap] for a version without the optional focus
func AtMap[K comparable, V any](index K) Optic[K, map[K]V, map[K]V, mo.Option[V], mo.Option[V], ReturnOne, ReadWrite, UniDir, Pure] {

	return Ret1(Rw(Ud((EPure(Compose(
		T2Of(
			Identity[map[K]V](),
			IgnoreWrite(Const[map[K]V](index)), //I lose RW here.
		),
		AtMapT2[K, V](),
	))))))

}

// AtMapT2 returns a [Lens] that focuses the element at tuple.B in a map tuple.A.
// When viewing the none option means no element was present in the map with the given index.
// When modifying the none option will remove the element from the map.
//
// See:
//   - [AtMap] for a version that uses a fixed index..
//   - [Element] for an optic that focuses by index position.
//   - [IndexMap] for a version without the optional focus
func AtMapT2[K comparable, V any]() Optic[K, lo.Tuple2[map[K]V, K], lo.Tuple2[map[K]V, K], mo.Option[V], mo.Option[V], ReturnOne, ReadWrite, UniDir, Pure] {
	return LensI[K, lo.Tuple2[map[K]V, K], mo.Option[V]](
		func(source lo.Tuple2[map[K]V, K]) (K, mo.Option[V]) {
			c, ok := source.A[source.B]
			return source.B, mo.TupleToOption(c, ok)
		},
		func(focus mo.Option[V], source lo.Tuple2[map[K]V, K]) lo.Tuple2[map[K]V, K] {
			clone := make(map[K]V, len(source.A))
			for k, v := range source.A {
				clone[k] = v
			}
			if nv, ok := focus.Get(); ok {
				clone[source.B] = nv
			} else {
				delete(clone, source.B)
			}

			return lo.T2(clone, source.B)
		},
		IxMatchComparable[K](),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.AtT2{
				OpticTypeExpr: ot,
				V:             reflect.TypeFor[V](),
			}
		}),
	)
}

// The MapOf combinator focuses on a map of all the elements in the given optic.
// If multiple focuses have the same index a [ErrDuplicateKey] is returned
//
// Under modification this map can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified map does not contain a key that is present in the source the result will use values from the original source.
// If the modified map contains additional keys they will be ignored.
//
// See:
//   - [MapOfP] for a polymorphic version.
func MapOf[I comparable, S, T, A, RETI, RW, DIR, ERR any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], size int) Optic[Void, S, T, map[I]A, map[I]A, ReturnOne, RW, UniDir, Err] {
	return EErr(MapOfReduced(
		o,
		duplicateKeyErrReducer[A](),
		size,
	))
}

func duplicateKeyErrReducer[A any]() ReductionP[mo.Option[A], A, A, Err] {
	return ReducerEP[mo.Option[A], A, A](
		func(ctx context.Context) (mo.Option[A], error) {
			return mo.None[A](), nil
		},
		func(ctx context.Context, state mo.Option[A], appendVal A) (mo.Option[A], error) {
			if state.IsPresent() {
				return mo.None[A](), ErrDuplicateKey
			}
			return mo.Some(appendVal), nil
		},
		func(ctx context.Context, state mo.Option[A]) (A, error) {
			return state.MustGet(), nil
		},
		ReducerExprDef(
			func(t expr.ReducerTypeExpr) expr.ReducerExpression {
				return expr.ErrDuplicateKeyReducerExpr{
					ReducerTypeExpr: t,
				}
			},
		),
	)

}

func toMapReduced[I comparable, S, T, A, B, RETI, RW, DIR, ERR, SR, RERR any](ctx context.Context, o Optic[I, S, T, A, B, RETI, RW, DIR, ERR], r ReductionP[SR, A, A, RERR], source S, size int) (Void, map[I]A, error) {
	reduction := make(map[I]SR)
	var err error

	o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
		index, focus, focusErr := val.Get()
		err = JoinCtxErr(ctx, focusErr)
		if err != nil {
			return false
		}

		rs, ok := reduction[index]
		if ok {
			rs, err = r.Reduce(ctx, rs, focus)
		} else {
			var empty SR
			empty, err = r.Empty(ctx)
			if err != nil {
				return false
			}

			rs, err = r.Reduce(ctx, empty, focus)
		}
		err = JoinCtxErr(ctx, focusErr)
		if err != nil {
			return false
		}

		reduction[index] = rs
		return true
	})

	if err != nil {
		return Void{}, nil, err
	}

	ret := make(map[I]A, size)

	for k, v := range reduction {
		reduced, err := r.End(ctx, v)
		if err != nil {
			return Void{}, nil, err
		}
		ret[k] = reduced
	}

	return Void{}, ret, err
}

// The MapOfReduced combinator focuses on a map of all the elements in the given optic.
// If multiple focuses have the same index the last focused value is used.
//
// Under modification this map can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified map does not contain a key that is present in the source the result will use values from the original source.
// If the modified map contains additional keys they will be ignored.
//
// See:
//   - [MapOfReducedP] for a polymorphic version.
//   - [MapOf] for a version that uses the first focus.
func MapOfReduced[I comparable, S, T, A, RETI, RW, DIR, ERR, SR, RERR any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], r ReductionP[SR, A, A, RERR], size int) Optic[Void, S, T, map[I]A, map[I]A, ReturnOne, RW, UniDir, CompositionTree[ERR, RERR]] {

	return CombiLens[RW, CompositionTree[ERR, RERR], Void, S, T, map[I]A, map[I]A](
		func(ctx context.Context, source S) (Void, map[I]A, error) {
			return toMapReduced(ctx, o, r, source, size)
		},
		func(ctx context.Context, newFocus map[I]A, source S) (T, error) {
			var innerErr error
			ret, err := o.AsModify()(ctx, func(index I, oldFocus A) (A, error) {
				newVal, found := newFocus[index]
				if !found {
					return oldFocus, nil
				}
				return newVal, nil
			}, source)

			return ret, JoinCtxErr(ctx, errors.Join(innerErr, err))
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.MapOfReduced{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					Reducer:       r.AsExpr(),
				}
			},
			o,
			r,
		),
	)
}

// The MapOfReducedP combinator focuses on a map of all the elements in the given optic.
// If multiple focuses have the same index the last focused value is used.
//
// Under modification this map can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified map does not contain a key that is present in the source then a [ErrUnsafeMissingElement] will be returned
// If the modified map contains additional keys they will be ignored.
//
// See:
//   - [MapOfReduced] for a non polymorphic version.
//   - [MapOfP] for a version that uses the first focus.
func MapOfReducedP[I comparable, S, T, A, B, RETI, RW, DIR, ERR, SR, RERR any](o Optic[I, S, T, A, B, RETI, RW, DIR, ERR], r ReductionP[SR, A, A, RERR], size int) Optic[Void, S, T, map[I]A, map[I]B, ReturnOne, RW, UniDir, Err] {

	return CombiLens[RW, Err, Void, S, T, map[I]A, map[I]B](
		func(ctx context.Context, source S) (Void, map[I]A, error) {
			return toMapReduced(ctx, o, r, source, size)
		},
		func(ctx context.Context, newFocus map[I]B, source S) (T, error) {
			var innerErr error
			ret, err := o.AsModify()(ctx, func(index I, oldFocus A) (B, error) {
				newVal, found := newFocus[index]
				if !found {
					var b B
					return b, ErrUnsafeMissingElement
				}
				return newVal, nil
			}, source)

			return ret, JoinCtxErr(ctx, errors.Join(innerErr, err))
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.MapOfReduced{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					Reducer:       r.AsExpr(),
				}
			},
			o,
			r,
		),
	)
}

// The MapOfP combinator focuses on a polymorphic map of all the elements in the given optic.
// If multiple focuses have the same index a [ErrDuplicateKey] is returned
//
// Under modification this map can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified map does not contain a key that is present in the source then [ErrUnsafeMissingElement] will be returned
// If the modified map contains additional keys they will be ignored.
//
// See:
//   - [MapOf] for a safe non polymorphic version.
func MapOfP[I comparable, S, T, A, B, RETI, RW, DIR, ERR any](o Optic[I, S, T, A, B, RETI, RW, DIR, ERR], size int) Optic[Void, S, T, map[I]A, map[I]B, ReturnOne, RW, UniDir, Err] {
	return EErr(MapOfReducedP(
		o,
		duplicateKeyErrReducer[A](),
		size,
	))
}

// MapToCol returns an [Iso] that converts a map to an [Collection]
//
// Note: In the case of duplicate indices ReverseGet will use the last occurrence.
func MapToCol[K comparable, A any]() Optic[Void, map[K]A, map[K]A, Collection[K, A, Pure], Collection[K, A, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return MapToColP[K, A, A]()
}

// MapToColP returns a polymorphic [Iso] that converts a map to an [Collection]
//
// Note: In the case of duplicate indices ReverseGet will use the last occurrence.
func MapToColP[I comparable, A, B any]() Optic[Void, map[I]A, map[I]B, Collection[I, A, Pure], Collection[I, B, Pure], ReturnOne, ReadWrite, BiDir, Pure] {

	return CombiIso[ReadWrite, BiDir, Pure, map[I]A, map[I]B, Collection[I, A, Pure], Collection[I, B, Pure]](
		func(ctx context.Context, source map[I]A) (Collection[I, A, Pure], error) {
			return ColI(
				func(yield func(index I, focus A) bool) {
					unsortedSeq := unsortedMapSeqF(source)
					sortedCol, _, _, err := heapSort(ctx, unsortedSeq, func(ctx context.Context, a, b ValueI[I, A]) (bool, error) {
						return mapKeyCompare(a.index, b.index) < 0, nil
					})

					if err != nil {
						panic(err) //less fnc above is pure
					}

					sortedCol(func(val ValueIE[I, A]) bool {
						index, focus, err := val.Get()
						if err != nil {
							panic(err) //less fnc above is pure
						}
						return yield(index, focus)
					})
				},
				func(index I) iter.Seq2[I, A] {
					return func(yield func(index I, focus A) bool) {
						if v, ok := source[index]; ok {
							yield(index, v)
						}
					}
				},
				IxMatchComparable[I](),
				func() int {
					return len(source)
				},
			), nil
		},
		func(ctx context.Context, focus Collection[I, B, Pure]) (map[I]B, error) {
			var retErr error
			ret := make(map[I]B)
			focus.AsIter()(ctx)(func(val ValueIE[I, B]) bool {
				index, value, focusErr := val.Get()
				if focusErr != nil {
					retErr = focusErr
					return false
				}
				ret[index] = value
				return true
			})
			return ret, retErr
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ToCol{
				OpticTypeExpr: ot,
				I:             reflect.TypeFor[I](),
				A:             reflect.TypeFor[A](),
				B:             reflect.TypeFor[B](),
			}
		}),
	)
}

//go:generate ./makecolops  -nocolof -unordered map_ops.go "optic" "Map" "K comparable, S any" "K comparable, S, T any" "" "MapColTypeP[K, S, S]()" "MapColTypeP[K, S, T]()" "K" "map[K]S" "map[K]T" "S" "T"

func MapColType[K comparable, S any]() CollectionType[K, map[K]S, map[K]S, S, S, Pure] {
	return MapColTypeP[K, S, S]()
}

// MapColTypeP returns a [CollectionType] wrapper for maps.
func MapColTypeP[K comparable, S, T any]() CollectionType[K, map[K]S, map[K]T, S, T, Pure] {
	return ColTypeP(
		MapToColP[K, S, T](),
		AsReverseGet(MapToColP[K, T, S]()),
		TraverseMapP[K, S, T](),
	)
}

//Map key comparison extracted from fmt.Print

func unsafeCompare[K comparable, T cmp.Ordered](aVal, bVal K) int {
	a := *(*T)(unsafe.Pointer(&aVal))
	b := *(*T)(unsafe.Pointer(&bVal))
	return cmp.Compare(a, b)
}

func mapKeyCompare[K comparable](aVal, bVal K) int {
	aType := reflect.TypeFor[K]()

	switch aType.Kind() {
	case reflect.Int:
		return unsafeCompare[K, int](aVal, bVal)
	case reflect.Int8:
		return unsafeCompare[K, int8](aVal, bVal)
	case reflect.Int16:
		return unsafeCompare[K, int16](aVal, bVal)
	case reflect.Int32:
		return unsafeCompare[K, int32](aVal, bVal)
	case reflect.Int64:
		return unsafeCompare[K, int64](aVal, bVal)
	case reflect.Uint:
		return unsafeCompare[K, uint](aVal, bVal)
	case reflect.Uint8:
		return unsafeCompare[K, uint8](aVal, bVal)
	case reflect.Uint16:
		return unsafeCompare[K, uint16](aVal, bVal)
	case reflect.Uint32:
		return unsafeCompare[K, uint32](aVal, bVal)
	case reflect.Uint64:
		return unsafeCompare[K, uint64](aVal, bVal)
	case reflect.Uintptr:
		return unsafeCompare[K, uintptr](aVal, bVal)
	case reflect.String:
		return unsafeCompare[K, string](aVal, bVal)
	case reflect.Float32:
		return unsafeCompare[K, float32](aVal, bVal)
	case reflect.Float64:
		return unsafeCompare[K, float64](aVal, bVal)
	case reflect.Complex64:

		a := *(*complex64)(unsafe.Pointer(&aVal))
		b := *(*complex64)(unsafe.Pointer(&bVal))
		if c := cmp.Compare(real(a), real(b)); c != 0 {
			return c
		}
		return cmp.Compare(imag(a), imag(b))

	case reflect.Complex128:

		a := *(*complex128)(unsafe.Pointer(&aVal))
		b := *(*complex128)(unsafe.Pointer(&bVal))
		if c := cmp.Compare(real(a), real(b)); c != 0 {
			return c
		}
		return cmp.Compare(imag(a), imag(b))

	case reflect.Bool:

		a := *(*bool)(unsafe.Pointer(&aVal))
		b := *(*bool)(unsafe.Pointer(&bVal))

		switch {
		case a == b:
			return 0
		case a:
			return 1
		default:
			return -1
		}

	case reflect.Pointer, reflect.UnsafePointer:
		a := reflect.ValueOf(aVal)
		b := reflect.ValueOf(bVal)
		return cmp.Compare(a.Pointer(), b.Pointer())

	case reflect.Chan:
		a := reflect.ValueOf(aVal)
		b := reflect.ValueOf(bVal)
		if c, ok := nilMapKeyCompare(a, b); ok {
			return c
		}
		return cmp.Compare(a.Pointer(), b.Pointer())

	case reflect.Struct:
		a := reflect.ValueOf(aVal)
		b := reflect.ValueOf(bVal)

		for i := 0; i < a.NumField(); i++ {
			if c := reflectMapKeyCompare(a.Field(i), b.Field(i)); c != 0 {
				return c
			}
		}
		return 0

	case reflect.Array:
		a := reflect.ValueOf(aVal)
		b := reflect.ValueOf(bVal)

		for i := 0; i < a.Len(); i++ {
			if c := reflectMapKeyCompare(a.Index(i), b.Index(i)); c != 0 {
				return c
			}
		}
		return 0

	case reflect.Interface:
		a := reflect.ValueOf(aVal).Convert(aType)
		b := reflect.ValueOf(bVal).Convert(aType)

		if c, ok := nilMapKeyCompare(a, b); ok {
			return c
		}
		c := reflectMapKeyCompare(reflect.ValueOf(a.Elem().Type()), reflect.ValueOf(b.Elem().Type()))
		if c != 0 {
			return c
		}
		return reflectMapKeyCompare(a.Elem(), b.Elem())

	default:
		// Certain types cannot appear as keys (maps, funcs, slices), but be explicit.
		panic("bad type in compare: " + aType.String())
	}
}

func isNilMapKeyable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return true
	default:
		return false
	}
}

// nilMapKeyCompare checks whether either value is nil. If not, the boolean is false.
// If either value is nil, the boolean is true and the integer is the comparison
// value. The comparison is defined to be 0 if both are nil, otherwise the one
// nil value compares low. Both arguments must represent a chan, func,
// interface, map, pointer, or slice.
func nilMapKeyCompare(aVal, bVal reflect.Value) (int, bool) {

	if aVal.IsNil() {
		if isNilMapKeyable(bVal) && bVal.IsNil() {
			return 0, true
		}
		return -1, true
	}
	if bVal.IsNil() {
		return 1, true
	}
	return 0, false
}

func reflectMapKeyCompare(aVal, bVal reflect.Value) int {
	aType, bType := aVal.Type(), bVal.Type()
	if aType != bType {
		return -1 // No good answer possible, but don't return 0: they're not equal.
	}
	switch aVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cmp.Compare(aVal.Int(), bVal.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return cmp.Compare(aVal.Uint(), bVal.Uint())
	case reflect.String:
		return cmp.Compare(aVal.String(), bVal.String())
	case reflect.Float32, reflect.Float64:
		return cmp.Compare(aVal.Float(), bVal.Float())
	case reflect.Complex64, reflect.Complex128:
		a, b := aVal.Complex(), bVal.Complex()
		if c := cmp.Compare(real(a), real(b)); c != 0 {
			return c
		}
		return cmp.Compare(imag(a), imag(b))
	case reflect.Bool:
		a, b := aVal.Bool(), bVal.Bool()
		switch {
		case a == b:
			return 0
		case a:
			return 1
		default:
			return -1
		}
	case reflect.Pointer, reflect.UnsafePointer:
		return cmp.Compare(aVal.Pointer(), bVal.Pointer())
	case reflect.Chan:
		if c, ok := nilMapKeyCompare(aVal, bVal); ok {
			return c
		}
		return cmp.Compare(aVal.Pointer(), bVal.Pointer())
	case reflect.Struct:
		for i := 0; i < aVal.NumField(); i++ {
			if c := reflectMapKeyCompare(aVal.Field(i), bVal.Field(i)); c != 0 {
				return c
			}
		}
		return 0
	case reflect.Array:
		for i := 0; i < aVal.Len(); i++ {
			if c := reflectMapKeyCompare(aVal.Index(i), bVal.Index(i)); c != 0 {
				return c
			}
		}
		return 0
	case reflect.Interface:
		if c, ok := nilMapKeyCompare(aVal, bVal); ok {
			return c
		}
		c := reflectMapKeyCompare(reflect.ValueOf(aVal.Elem().Type()), reflect.ValueOf(bVal.Elem().Type()))
		if c != 0 {
			return c
		}
		return reflectMapKeyCompare(aVal.Elem(), bVal.Elem())
	default:
		// Certain types cannot appear as keys (maps, funcs, slices), but be explicit.
		panic("bad type in compare: " + aType.String())
	}
}

func MakeMap[S any, K comparable, V any](size int) Optic[Void, S, S, map[K]V, map[K]V, ReturnOne, ReadOnly, UniDir, Pure] {
	return Getter[S, map[K]V](
		func(source S) map[K]V {
			return make(map[K]V, size)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Make{
				OpticTypeExpr: ot,
				Size:          []int{size},
			}
		}),
	)
}
