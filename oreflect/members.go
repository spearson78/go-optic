package oreflect

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
	"github.com/spearson78/go-optic/otree"
)

type Member struct {
	Kind        reflect.Kind
	Type        reflect.Type
	StructField reflect.StructField
	SliceIndex  int
	MapKey      any
}

func (m Member) String() string {
	switch m.Kind {
	case reflect.Struct:
		return fmt.Sprintf(".%v", m.StructField.Name)
	case reflect.Slice:
		return fmt.Sprintf("[%v]", m.SliceIndex)
	case reflect.Map:
		return fmt.Sprintf("[%v]", m.MapKey)
	default:
		return ".??"
	}
}

// TraverseMembers returns a [Traversal] that iterates recursively through the source struct focusing on instances of the focus type.
// TraverseMembers uses reflection to traverse through
// - fields of a struct
// - elements of a slice
// - values in a map
//
// If the traversal is recursive then only the first level of matches will be focused.
// To fully iterate a recursive structure use TraverseMembers as with [TraverseTopDown] or [TraverseBottomUp]
func TraverseMembers[S, A any]() Optic[*otree.PathNode[Member], S, S, A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseFilteredMembers[S, A](True[reflect.Type]())
}

func TraverseFilteredMembers[S, A, ERR any](reflectInto Predicate[reflect.Type, ERR]) Optic[*otree.PathNode[Member], S, S, A, A, ReturnMany, ReadWrite, UniDir, ERR] {

	return CombiTraversal[ReturnMany, ReadWrite, ERR, *otree.PathNode[Member], S, S, A, A](
		func(ctx context.Context, source S) SeqIE[*otree.PathNode[Member], A] {
			return func(yield func(ValueIE[*otree.PathNode[Member], A]) bool) {

				targetType := reflect.TypeFor[A]()
				original := reflect.ValueOf(source)

				visitedPtrs := make(map[unsafe.Pointer]Void)

				membersIterRecursive(ctx, nil, original, targetType, visitedPtrs, reflectInto, func(val ValueIE[*otree.PathNode[Member], reflect.Value]) bool {
					index, focus, focusErr := val.Get()
					if focusErr != nil {
						var a A
						return yield(ValIE(index, a, focusErr))
					}
					a := focus.Interface().(A)
					return yield(ValIE(index, a, ctx.Err()))
				})

			}
		},
		nil,
		func(ctx context.Context, fmap func(index *otree.PathNode[Member], focus A) (A, error), source S) (S, error) {
			targetType := reflect.TypeFor[A]()

			original := reflect.ValueOf(source)
			copy := reflect.New(original.Type()).Elem()
			copy.Set(original)

			clonedPtrs := make(map[unsafe.Pointer]reflect.Value)

			err := membersModifyRecursive(ctx, nil, copy, targetType, clonedPtrs, reflectInto, func(index *otree.PathNode[Member], focus reflect.Value) (reflect.Value, error) {
				a := focus.Interface().(A)
				ret, err := fmap(index, a)
				err = JoinCtxErr(ctx, err)
				return reflect.ValueOf(ret), err
			})

			return copy.Interface().(S), JoinCtxErr(ctx, err)
		},
		nil,
		PredToIxMatch(otree.EqPathT2(EqDeepT2[Member]())),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.TraverseMembers{
					OpticTypeExpr: ot,
				}
			},
			reflectInto,
		),
	)
}

// modified version of https://gist.github.com/hvoecking/10772475
func membersModifyRecursive[ERR any](ctx context.Context, path *otree.PathNode[Member], copy reflect.Value, targetType reflect.Type, clonedPtrs map[unsafe.Pointer]reflect.Value, into Predicate[reflect.Type, ERR], fmap func(index *otree.PathNode[Member], focus reflect.Value) (reflect.Value, error)) error {

	allowed, err := PredGet(ctx, into, copy.Type())
	err = JoinCtxErr(ctx, err)
	if err != nil {
		return err
	}

	if !allowed {
		return ctx.Err()
	}

	if path != nil && copy.Type().ConvertibleTo(targetType) {
		origAsTarget := copy.Convert(targetType)
		val, err := fmap(path, origAsTarget)
		err = JoinCtxErr(ctx, err)
		if err == nil {
			copy.Set(val)
		}
		return err
	} else {
		switch copy.Kind() {
		case reflect.Ptr:
			if copy.IsNil() {
				return ctx.Err()
			}

			if copy.Elem().Type().PkgPath() == "reflect" {
				return ctx.Err()
			}

			ptrVal := copy.UnsafePointer()

			cloned, ok := clonedPtrs[ptrVal]
			if ok {
				copy.Set(cloned)
			}

			originalValue := copy.Elem()
			cloned = reflect.New(originalValue.Type())
			clonedPtrs[ptrVal] = cloned

			copy.Set(cloned)
			copy.Elem().Set(originalValue)
			err := membersModifyRecursive(ctx, path, copy.Elem(), targetType, clonedPtrs, into, fmap)
			return JoinCtxErr(ctx, err)
		case reflect.Interface:
			originalValue := copy.Elem()
			if !originalValue.IsValid() {
				return ctx.Err()
			}

			copyValue := reflect.New(originalValue.Type()).Elem()
			copyValue.Set(originalValue)

			err := membersModifyRecursive(ctx, path, copyValue, targetType, clonedPtrs, into, fmap)

			copy.Set(copyValue)
			return JoinCtxErr(ctx, err)
		case reflect.Struct:

			if copy.Type().PkgPath() == "reflect" {
				return ctx.Err()
			}

			for i := 0; i < copy.NumField(); i += 1 {

				childPath := path.Append(Member{
					Kind:        reflect.Struct,
					Type:        copy.Type(),
					StructField: copy.Type().Field(i),
				})

				copyField := copy.Field(i)
				copyUnexpField := reflect.NewAt(copyField.Type(), unsafe.Pointer(copyField.UnsafeAddr())).Elem()

				err := membersModifyRecursive(ctx, childPath, copyUnexpField, targetType, clonedPtrs, into, fmap)

				err = JoinCtxErr(ctx, err)
				if err != nil {
					return err
				}
			}

			return ctx.Err()

		case reflect.Slice:
			copyValue := reflect.MakeSlice(copy.Type(), copy.Len(), copy.Cap())
			reflect.Copy(copyValue, copy)
			copy.Set(copyValue)

			for i := 0; i < copy.Len(); i += 1 {

				childPath := path.Append(Member{
					Kind:       reflect.Slice,
					Type:       copy.Type().Elem(),
					SliceIndex: i,
				})

				err := membersModifyRecursive(ctx, childPath, copy.Index(i), targetType, clonedPtrs, into, fmap)

				err = JoinCtxErr(ctx, err)
				if err != nil {
					return err
				}
			}

			return ctx.Err()
		case reflect.Map:
			copiedMapValue := reflect.MakeMap(copy.Type())

			for _, key := range copy.MapKeys() {
				originalValue := copy.MapIndex(key)

				childPath := path.Append(Member{
					Kind:   reflect.Map,
					Type:   copy.Type().Elem(),
					MapKey: key.Interface(),
				})

				copyValue := reflect.New(originalValue.Type()).Elem()
				copyValue.Set(originalValue)
				err := membersModifyRecursive(ctx, childPath, copyValue, targetType, clonedPtrs, into, fmap)

				err = JoinCtxErr(ctx, err)
				if err != nil {
					return err
				}
				copiedMapValue.SetMapIndex(key, copyValue)
			}

			copy.Set(copiedMapValue)

			return ctx.Err()
		default:
			return ctx.Err()
		}
	}
}

func membersIterRecursive[ERR any](ctx context.Context, path *otree.PathNode[Member], original reflect.Value, targetType reflect.Type, visitedPtrs map[unsafe.Pointer]Void, into Predicate[reflect.Type, ERR], yield func(ValueIE[*otree.PathNode[Member], reflect.Value]) bool) bool {

	allowed, err := PredGet(ctx, into, original.Type())
	err = JoinCtxErr(ctx, err)
	if err != nil {
		return yield(ValIE(path, reflect.ValueOf(nil), err))
	}

	if !allowed {
		return true
	}

	if path != nil && original.Type().ConvertibleTo(targetType) {
		origAsTarget := original.Convert(targetType)
		return yield(ValIE(path, origAsTarget, ctx.Err()))
	} else {
		switch original.Kind() {
		case reflect.Ptr:
			if original.IsNil() {
				return true
			}

			if original.Elem().Type().PkgPath() == "reflect" {
				return true
			}

			ptrVal := original.UnsafePointer()

			_, ok := visitedPtrs[ptrVal]
			if ok {
				return true
			}

			return membersIterRecursive(ctx, path, original.Elem(), targetType, visitedPtrs, into, yield)

		case reflect.Interface:
			originalValue := original.Elem()
			if !originalValue.IsValid() {
				return true
			}

			return membersIterRecursive(ctx, path, originalValue, targetType, visitedPtrs, into, yield)
		case reflect.Struct:

			cloned := false
			var unexpClone reflect.Value

			if original.Type().PkgPath() == "reflect" {
				return true
			} else {
				cont := true
				for i := 0; i < original.NumField(); i += 1 {

					childPath := path.Append(Member{
						Kind:        reflect.Struct,
						Type:        original.Type(),
						StructField: original.Type().Field(i),
					})

					if !original.Type().Field(i).IsExported() {
						if !cloned {
							unexpClone = reflect.New(original.Type()).Elem()
							unexpClone.Set(original)
						}

						cloneField := unexpClone.Field(1)
						cloneField = reflect.NewAt(cloneField.Type(), unsafe.Pointer(cloneField.UnsafeAddr())).Elem()

						cont = membersIterRecursive(ctx, childPath, cloneField, targetType, visitedPtrs, into, yield)
						if !cont {
							break
						}
					} else {
						cont = membersIterRecursive(ctx, childPath, original.Field(i), targetType, visitedPtrs, into, yield)
						if !cont {
							break
						}
					}
				}
				return cont
			}
		case reflect.Slice:

			cont := true

			for i := 0; i < original.Len(); i += 1 {

				childPath := path.Append(Member{
					Kind:       reflect.Slice,
					Type:       original.Type().Elem(),
					SliceIndex: i,
				})

				cont = membersIterRecursive(ctx, childPath, original.Index(i), targetType, visitedPtrs, into, yield)
				if !cont {
					break
				}
			}

			return cont
		case reflect.Map:

			cont := true

			for _, key := range original.MapKeys() {
				originalValue := original.MapIndex(key)

				childPath := path.Append(Member{
					Kind:   reflect.Map,
					Type:   original.Type().Elem(),
					MapKey: key.Interface(),
				})

				cont = membersIterRecursive(ctx, childPath, originalValue, targetType, visitedPtrs, into, yield)
				if !cont {
					break
				}
			}

			return cont
		default:
			return true
		}
	}
}
