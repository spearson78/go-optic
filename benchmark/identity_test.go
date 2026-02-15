package benchmark

import (
	"context"
	"testing"
	"unsafe"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

func identityComposeLeft[I, J, S, T, A, B, C, D, RETI, RETJ, RWI, RWJ, DIRI, DIRJ, ERRI, ERRJ any](o1 Optic[I, S, T, A, B, RETI, RWI, DIRI, ERRI], o2 Optic[J, A, B, C, D, RETJ, RWJ, DIRJ, ERRJ]) Optic[I, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]] {
	return UnsafeOmni[I, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]](
		func(ctx context.Context, source S) (I, C, error) {
			panic("Not Implemented")
		},
		func(ctx context.Context, focus D, source S) (T, error) {
			panic("not implemented")
		},
		func(ctx context.Context, source S) SeqIE[I, C] {
			panic("not implemented")
		},
		func(ctx context.Context, source S) (int, error) {
			panic("not implemented")
		},
		func(ctx context.Context, fmap func(index I, focus C) (D, error), source S) (T, error) {
			return o1.AsModify()(ctx, func(index I, focus A) (B, error) {
				ret, err := fmap(index, *(*C)(unsafe.Pointer(&focus)))
				return *(*B)(unsafe.Pointer(&ret)), err
			}, source)
		},
		func(ctx context.Context, index I, source S) SeqIE[I, C] {
			panic("not implemented")
		},
		func(indexA, indexB I) bool {
			panic("not implemented")
		},
		func(ctx context.Context, focus D) (T, error) {
			panic("not implemented")
		},
		func(ctx context.Context) (ExprHandler, error) {
			return nil, nil
		},
		func() expr.OpticExpression {
			return expr.Custom("custom")(expr.NewOpticTypeExpr[I, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]]())
		},
	)

}

func identityCompose[I, J, S, T, A, B, C, D, RETI, RETJ any, RWI, RWJ any, DIRI, DIRJ any, ERRI, ERRJ any](left Optic[I, S, T, A, B, RETI, RWI, DIRI, ERRI], right Optic[J, A, B, C, D, RETJ, RWJ, DIRJ, ERRJ]) Optic[J, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]] {
	return UnsafeOmni[J, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]](
		func(ctx context.Context, source S) (J, C, error) {
			panic("not implemented")
		},
		func(ctx context.Context, focus D, source S) (T, error) {
			panic("not implemented")
		},
		func(ctx context.Context, source S) SeqIE[J, C] {
			panic("not implemented")
		},
		func(ctx context.Context, source S) (int, error) {
			panic("not implemented")
		},
		func(ctx context.Context, fmap func(index J, focus C) (D, error), source S) (T, error) {
			ret, err := right.AsModify()(ctx, fmap, *(*A)(unsafe.Pointer((&source))))
			return *(*T)(unsafe.Pointer(&ret)), err
		},
		func(ctx context.Context, index J, source S) SeqIE[J, C] {
			panic("not implemented")
		},
		func(indexA, indexB J) bool {
			panic("not implemented")
		},
		func(ctx context.Context, focus D) (T, error) {
			panic("not implemented")
		},
		func(ctx context.Context) (ExprHandler, error) {
			return nil, nil
		},
		func() expr.OpticExpression {
			return expr.Custom("custom")(expr.NewOpticTypeExpr[I, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]]())
		},
	)
}

func BenchmarkIdentity(b *testing.B) {

	slice := MustGet(ColToSlice[int](), RangeCol[int](1, 1000))
	add1 := Add(1)
	identity := Identity[int]()
	traverse := TraverseSlice[int]()
	traverseIdentity := ComposeLeft(traverse, Identity[int]())
	optimizedTraverseIdentity := identityComposeLeft(traverse, Identity[int]())
	identityTraverseIdentity := Compose(Identity[[]int](), traverseIdentity)

	optimizedIdentityTraverseIdentity := identityCompose(Identity[[]int](), optimizedTraverseIdentity)

	b.Run("traverse", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			MustModify(traverse, add1, slice)
		}
	})

	b.Run("traverse op identity", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			MustModify(traverse, identity, slice)
		}
	})

	b.Run("traverse identity", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			MustModify(traverseIdentity, add1, slice)
		}
	})

	b.Run("opti traverse identity", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			MustModify(optimizedTraverseIdentity, add1, slice)
		}
	})

	b.Run("identity traverse identity", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			MustModify(identityTraverseIdentity, add1, slice)
		}
	})

	b.Run("opti identity traverse identity", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			MustModify(optimizedIdentityTraverseIdentity, add1, slice)
		}
	})

}
