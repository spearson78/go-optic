package optic

import (
	"context"

	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
)

// Identity returns an [Iso] that focuses on the source.
func Identity[S any]() Optic[Void, S, S, S, S, ReturnOne, ReadWrite, BiDir, Pure] {
	return IdentityP[S, S]()
}

// IdentityP returns a polymorphic [Iso] that focuses on the source.
func IdentityP[S, T any]() Optic[Void, S, T, S, T, ReturnOne, ReadWrite, BiDir, Pure] {
	return IdentityIP[Void, S, T](Void{}, IxMatchVoid())
}

// IdentityIP returns an index aware polymorphic [Iso] that focuses on the source and given index.
func IdentityI[I, S, T any](i I, ixMatch func(a, b I) bool) Optic[I, S, S, S, S, ReturnOne, ReadWrite, BiDir, Pure] {
	return IdentityIP[I, S, S](i, ixMatch)
}

// IdentityIP returns an index aware polymorphic [Iso] that focuses on the source and given index.
func IdentityIP[I, S, T any](i I, ixMatch func(a, b I) bool) Optic[I, S, T, S, T, ReturnOne, ReadWrite, BiDir, Pure] {

	omni := UnsafeOmni[I, S, T, S, T, ReturnOne, ReadWrite, BiDir, Pure](
		func(ctx context.Context, source S) (I, S, error) {
			return i, source, nil
		},
		func(ctx context.Context, focus T, source S) (T, error) {
			return focus, nil
		},
		func(ctx context.Context, source S) SeqIE[I, S] {
			return func(yield func(ValueIE[I, S]) bool) {
				yield(ValIE(i, source, nil))
			}
		},
		func(ctx context.Context, source S) (int, error) {
			return 1, nil
		},
		func(ctx context.Context, fmap func(index I, focus S) (T, error), source S) (T, error) {
			return fmap(i, source)
		},
		func(ctx context.Context, index I, source S) SeqIE[I, S] {
			return func(yield func(ValueIE[I, S]) bool) {
				yield(ValIE(i, source, nil))
			}
		},
		ixMatch,
		func(ctx context.Context, focus T) (T, error) {
			return focus, nil
		},
		nil,
		func() expr.OpticExpression {
			return expr.Identity{
				Index:         i,
				OpticTypeExpr: expr.NewOpticTypeExpr[Void, S, T, S, T, ReturnOne, ReadWrite, BiDir, Pure](),
			}
		},
	).(omniOptic[I, S, T, S, T, ReturnOne, ReadWrite, BiDir, Pure])

	omni.opticType |= expr.OpticTypeIdentity

	return omni
}

// Value returns an optic that focuses the given value that can be changed under modification.
//
// See:
//   - [Const] for a polymorphic version that doesn't allow updates under modification
func Value[S any](v S) Optic[Void, S, S, S, S, ReturnOne, ReadWrite, BiDir, Pure] {
	return Iso[S, S](
		func(source S) S {
			return v
		},
		func(focus S) S {
			return focus
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Value{
				OpticTypeExpr: ot,
				Value:         v,
			}
		}),
	)
}

// The IgnoreWrite modifes the given optic to be ReadWrite at the cost of ignoring any updates the original source is returned instead.
// This can be useful for making [Const] ReadWrite.
func IgnoreWrite[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, S, A, B, RET, ReadWrite, UniDir, ERR] {
	return Omni[I, S, S, A, B, RET, ReadWrite, UniDir, ERR](
		o.AsGetter(),
		func(ctx context.Context, focus B, source S) (S, error) {
			return source, nil
		},
		o.AsIter(),
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (S, error) {
			return source, nil
		},
		o.AsIxGetter(),
		o.AsIxMatch(),
		unsupportedReverseGetter[B, S],
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.IgnoreWrite{
				OpticTypeExpr: ot,
				Optic:         o.AsExpr(),
			}
		}, o),
	)
}

func Const[S any, A any](v A) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return ConstP[S, S, A, A](v)
}

func ConstP[S, T, A, B any](v A) Optic[Void, S, T, A, B, ReturnOne, ReadOnly, UniDir, Pure] {
	//WARNING: Do not call ConstI or ConstIP with a true predicate as true is implemented using Const. This will cause a stack overflow.
	return CombiGetter[Pure, Void, S, T, A, B](
		func(ctx context.Context, source S) (Void, A, error) {
			return Void{}, v, nil
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Const{
					OpticTypeExpr: ot,
					Index:         Void{},
					Value:         v,
				}
			},
		),
	)
}

// ConstI returns an optic that focuses a constant index and value.
//
// See:
//   - [Const] for a non index aware version
func ConstI[S, A any, I any, ERR TPure](i I, v A, ixMatch Predicate[lo.Tuple2[I, I], ERR]) Optic[I, S, S, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return ConstIP[S, S, A, A, I](i, v, ixMatch)
}

func ConstIP[S, T, A, B any, I any, ERR TPure](i I, v A, ixMatch Predicate[lo.Tuple2[I, I], ERR]) Optic[I, S, T, A, B, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, I, S, T, A, B](
		func(ctx context.Context, source S) (I, A, error) {
			return i, v, nil
		},
		func(indexA, indexB I) bool {
			return Must(PredGet(context.Background(), ixMatch, lo.T2(indexA, indexB)))
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Const{
					OpticTypeExpr: ot,
					Index:         i,
					Value:         v,
				}
			},
			ixMatch,
		),
	)
}

// Nothing returns a [Traversal] that focuses no elements
func Nothing[S any, A any]() Optic[Void, S, S, A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return NothingI[S, A, Void](Const[lo.Tuple2[Void, Void]](true))
}

// NothingI returns an index aware [Traversal] that focuses no elements
func NothingI[S any, A, I any, ERR TPure](ixMatch Predicate[lo.Tuple2[I, I], ERR]) Optic[I, S, S, A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, I, S, S, A, A](
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {}
		},
		func(ctx context.Context, source S) (int, error) {
			return 0, nil
		},
		func(ctx context.Context, fmap func(index I, focus A) (A, error), source S) (S, error) {
			return source, nil
		},
		nil,
		func(indexA, indexB I) bool {
			return Must(PredGet(context.Background(), ixMatch, lo.T2(indexA, indexB)))
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Empty{
					OpticTypeExpr: ot,
				}
			},
			ixMatch,
		),
	)
}
