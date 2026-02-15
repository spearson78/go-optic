package optic

import (
	"context"
	"errors"
	"reflect"
	"unsafe"

	"github.com/spearson78/go-optic/expr"
)

// FieldLens returns a [Lens] focusing on a field within a struct
//
// A FieldLens is constructed from a single function
// - fref: should return the address of the field within the given struct.
//
// Note: FieldLenses are normally created by the makelens tool which also provides a convenient composition builder.
//
// See:
// - [PtrFieldLens] for a lens that focuses on a pointer of the source.
// - [PtrFieldLensE] for a lens that focuses on a pointer of the source where the source is not expected to be nil.
func FieldLens[S, A any](fref func(source *S) *A) Optic[Void, S, S, A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return fieldLensP(
		fref,
		func(ctx context.Context, source S) (Void, A, error) {
			ptrA := fref(&source)
			return Void{}, *ptrA, ctx.Err()
		},
		func(ctx context.Context, focus A, source S) (S, error) {
			ptrA := fref(&source)
			*ptrA = focus
			return source, nil
		},
	)
}

func FieldLensP[S, T, A, B any](get func(source *S) *A, set func(focus B, source S) T) Optic[Void, S, T, A, B, ReturnOne, ReadWrite, UniDir, Pure] {
	return fieldLensP(
		get,
		func(ctx context.Context, source S) (Void, A, error) {
			ptrA := get(&source)
			return Void{}, *ptrA, ctx.Err()
		},
		func(ctx context.Context, focus B, source S) (T, error) {
			b := set(focus, source)
			return b, ctx.Err()

		},
	)
}

func fieldLensP[S, T, A, B any](fref func(source *S) *A, get func(ctx context.Context, source S) (Void, A, error), set func(ctx context.Context, focus B, source S) (T, error)) Optic[Void, S, T, A, B, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, Void, S, T, A, B](
		get,
		set,
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			source := reflect.New(ot.OpticS())
			return getFieldLensExpr(source, ot, fref)
		}),
	)
}

// PtrFieldLens returns a [Traversal] focusing on a field within a pointer to a struct. If the source pointer is nil then nothing is focused.
//
// A PtrFieldLens is constructed from a single function
// - fref: should return the address of the field within the given struct.
//
// Note: PtrFieldLenses are normally created by the makelens tool which also provides a convenient composition builder.
//
// See:
// - [FieldLens] for a lens that focuses on a value of the source.
// - [PtrFieldLensE] for a lens that focuses on a pointer of the source where the source is not expected to be nil.
func PtrFieldLens[S, A any](fref func(source *S) *A) Optic[Void, *S, *S, A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, Void, *S, *S, A, A](
		func(ctx context.Context, source *S) SeqIE[Void, A] {
			return func(yield func(ValueIE[Void, A]) bool) {
				if source != nil {
					yield(ValIE(Void{}, *fref(source), ctx.Err()))
				}
			}
		},
		func(ctx context.Context, source *S) (int, error) {
			if source != nil {
				return 1, nil
			} else {
				return 0, nil
			}
		},
		func(ctx context.Context, fmap func(index Void, focus A) (A, error), source *S) (*S, error) {
			if source == nil {
				return nil, nil
			}

			ret := *source
			ptr := fref(&ret)
			val, err := fmap(Void{}, *ptr)
			err = JoinCtxErr(ctx, err)
			if err != nil {
				return source, err
			}

			*ptr = val

			return &ret, nil
		},
		nil,
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			source := reflect.New(ot.OpticS().Elem())
			return getFieldLensExpr(source, ot, fref)
		}),
	)
}

var ErrNilPointer = errors.New("nil pointer passed to MustPtrFieldLens")

// PtrFieldLensE returns a [Lens] focusing on a field within a pointer to a struct. If the source pointer is nil then [ErrNilPointer] is returned.
//
// A PtrFieldLensE is constructed from a single function
// - fref: should return the address of the field within the given struct.
//
// See:
// - [FieldLens] for a lens that focuses on a value of the source.
// - [PtrFieldLens] for a lens that focuses on a pointer of the source but ignores nil sources.
func PtrFieldLensE[S, A any](fref func(source *S) *A) Optic[Void, *S, *S, A, A, ReturnOne, ReadWrite, UniDir, Err] {
	return LensE[*S, A](
		func(ctx context.Context, source *S) (A, error) {
			if source == nil {
				var a A
				return a, ErrNilPointer
			}

			ptrA := fref(source)

			return *ptrA, ctx.Err()
		},
		func(ctx context.Context, focus A, source *S) (*S, error) {
			if source == nil {
				return source, ErrNilPointer
			}

			ret := *source
			ptrA := fref(&ret)
			*ptrA = focus
			return &ret, ctx.Err()
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			source := reflect.New(ot.OpticS().Elem())
			return getFieldLensExpr(source, ot, fref)
		}),
	)
}

func getFieldLensExpr[S, A any](sourceVal reflect.Value, ot expr.OpticTypeExpr, fref func(source *S) *A) expr.FieldLens {

	source := sourceVal.Interface().(*S)
	field := fref(source)

	fieldOffset := uintptr(unsafe.Pointer(field)) - uintptr(unsafe.Pointer(source))

	rType := reflect.ValueOf(source).Elem().Type()
	numFields := rType.NumField()

	for fieldNum := 0; fieldNum < numFields; fieldNum++ {
		field := rType.Field(fieldNum)

		if field.Name == "" || field.Type.Size() == 0 {
			continue
		}

		if fieldOffset == field.Offset {
			return expr.FieldLens{
				OpticTypeExpr: ot,
				Field:         field,
				FieldNum:      fieldNum,
			}
		}
	}

	panic("getFieldLensExpr: field not found")
}
