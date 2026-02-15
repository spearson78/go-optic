package optic

import (
	"context"
	"iter"
	"reflect"

	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
)

// TraverseT2 returns a [Traversal] focusing on the elements of a [lo.Tuple2]
//
// See: [TraverseT2P] for a polymorphic version
func TraverseT2[A any]() Optic[int, lo.Tuple2[A, A], lo.Tuple2[A, A], A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseT2P[A, A]()
}

// TraverseT2P returns a polymorphic [Traversal] focusing on the elements of a [lo.Tuple2]
//
// See: [TraverseT2] for a non polymorphic version
func TraverseT2P[A any, B any]() Optic[int, lo.Tuple2[A, A], lo.Tuple2[B, B], A, B, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, int, lo.Tuple2[A, A], lo.Tuple2[B, B], A, B](func(ctx context.Context, source lo.Tuple2[A, A]) SeqIE[int, A] {
		return traverseTNIter(ctx, &source.A, &source.B)
	}, func(ctx context.Context, source lo.Tuple2[A, A]) (int, error) {
		return 2, nil
	}, func(ctx context.Context, fmap func(index int, focus A) (B, error), source lo.Tuple2[A, A]) (lo.Tuple2[B, B], error) {
		ret, err := traverseTNModify(ctx, fmap, &source.A, &source.B)
		return lo.T2(ret[0], ret[1]), err
	}, func(ctx context.Context, index int, source lo.Tuple2[A, A]) SeqIE[int, A] {
		return traverseTNIxGet(index, &source.A, &source.B)
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.Traverse{
			OpticTypeExpr: ot,
		}
	}))
}

// TraverseT3 returns a [Traversal] focusing on the elements of a [lo.Tuple3]
//
// See: [TraverseT3P] for a polymorphic version
func TraverseT3[A any]() Optic[int, lo.Tuple3[A, A, A], lo.Tuple3[A, A, A], A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseT3P[A, A]()
}

// TraverseT3P returns a polymorphic [Traversal] focusing on the elements of a [lo.Tuple3]
//
// See: [TraverseT3] for a non polymorphic version
func TraverseT3P[A any, B any]() Optic[int, lo.Tuple3[A, A, A], lo.Tuple3[B, B, B], A, B, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, int, lo.Tuple3[A, A, A], lo.Tuple3[B, B, B], A, B](func(ctx context.Context, source lo.Tuple3[A, A, A]) SeqIE[int, A] {
		return traverseTNIter(ctx, &source.A, &source.B, &source.C)
	}, func(ctx context.Context, source lo.Tuple3[A, A, A]) (int, error) {
		return 3, nil
	}, func(ctx context.Context, fmap func(index int, focus A) (B, error), source lo.Tuple3[A, A, A]) (lo.Tuple3[B, B, B], error) {
		ret, err := traverseTNModify(ctx, fmap, &source.A, &source.B, &source.C)
		return lo.T3(ret[0], ret[1], ret[2]), err
	}, func(ctx context.Context, index int, source lo.Tuple3[A, A, A]) SeqIE[int, A] {
		return traverseTNIxGet(index, &source.A, &source.B, &source.C)
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.Traverse{
			OpticTypeExpr: ot,
		}
	}))
}

// TraverseT4 returns a [Traversal] focusing on the elements of a [lo.Tuple4]
//
// See: [TraverseT4P] for a polymorphic version
func TraverseT4[A any]() Optic[int, lo.Tuple4[A, A, A, A], lo.Tuple4[A, A, A, A], A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseT4P[A, A]()
}

// TraverseT4P returns a polymorphic [Traversal] focusing on the elements of a [lo.Tuple4]
//
// See: [TraverseT4] for a non polymorphic version
func TraverseT4P[A any, B any]() Optic[int, lo.Tuple4[A, A, A, A], lo.Tuple4[B, B, B, B], A, B, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, int, lo.Tuple4[A, A, A, A], lo.Tuple4[B, B, B, B], A, B](func(ctx context.Context, source lo.Tuple4[A, A, A, A]) SeqIE[int, A] {
		return traverseTNIter(ctx, &source.A, &source.B, &source.C, &source.D)
	}, func(ctx context.Context, source lo.Tuple4[A, A, A, A]) (int, error) {
		return 4, nil
	}, func(ctx context.Context, fmap func(index int, focus A) (B, error), source lo.Tuple4[A, A, A, A]) (lo.Tuple4[B, B, B, B], error) {
		ret, err := traverseTNModify(ctx, fmap, &source.A, &source.B, &source.C, &source.D)
		return lo.T4(ret[0], ret[1], ret[2], ret[3]), err
	}, func(ctx context.Context, index int, source lo.Tuple4[A, A, A, A]) SeqIE[int, A] {
		return traverseTNIxGet(index, &source.A, &source.B, &source.C, &source.D)
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.Traverse{
			OpticTypeExpr: ot,
		}
	}))
}

// TraverseT5 returns a [Traversal] focusing on the elements of a [lo.Tuple5]
//
// See: [TraverseT5P] for a polymorphic version
func TraverseT5[A any]() Optic[int, lo.Tuple5[A, A, A, A, A], lo.Tuple5[A, A, A, A, A], A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseT5P[A, A]()
}

// TraverseT5P returns a polymorphic [Traversal] focusing on the elements of a [lo.Tuple5]
//
// See: [TraverseT5] for a non polymorphic version
func TraverseT5P[A any, B any]() Optic[int, lo.Tuple5[A, A, A, A, A], lo.Tuple5[B, B, B, B, B], A, B, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, int, lo.Tuple5[A, A, A, A, A], lo.Tuple5[B, B, B, B, B], A, B](func(ctx context.Context, source lo.Tuple5[A, A, A, A, A]) SeqIE[int, A] {
		return traverseTNIter(ctx, &source.A, &source.B, &source.C, &source.D, &source.E)
	}, func(ctx context.Context, source lo.Tuple5[A, A, A, A, A]) (int, error) {
		return 5, nil
	}, func(ctx context.Context, fmap func(index int, focus A) (B, error), source lo.Tuple5[A, A, A, A, A]) (lo.Tuple5[B, B, B, B, B], error) {
		ret, err := traverseTNModify(ctx, fmap, &source.A, &source.B, &source.C, &source.D, &source.E)
		return lo.T5(ret[0], ret[1], ret[2], ret[3], ret[4]), err
	}, func(ctx context.Context, index int, source lo.Tuple5[A, A, A, A, A]) SeqIE[int, A] {
		return traverseTNIxGet(index, &source.A, &source.B, &source.C, &source.D, &source.E)
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.Traverse{
			OpticTypeExpr: ot,
		}
	}))
}

// TraverseT6 returns a [Traversal] focusing on the elements of a [lo.Tuple6]
//
// See: [TraverseT6P] for a polymorphic version
func TraverseT6[A any]() Optic[int, lo.Tuple6[A, A, A, A, A, A], lo.Tuple6[A, A, A, A, A, A], A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseT6P[A, A]()
}

// TraverseT6P returns a polymorphic [Traversal] focusing on the elements of a [lo.Tuple6]
//
// See: [TraverseT6] for a non polymorphic version
func TraverseT6P[A any, B any]() Optic[int, lo.Tuple6[A, A, A, A, A, A], lo.Tuple6[B, B, B, B, B, B], A, B, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, int, lo.Tuple6[A, A, A, A, A, A], lo.Tuple6[B, B, B, B, B, B], A, B](func(ctx context.Context, source lo.Tuple6[A, A, A, A, A, A]) SeqIE[int, A] {
		return traverseTNIter(ctx, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F)
	}, func(ctx context.Context, source lo.Tuple6[A, A, A, A, A, A]) (int, error) {
		return 6, nil
	}, func(ctx context.Context, fmap func(index int, focus A) (B, error), source lo.Tuple6[A, A, A, A, A, A]) (lo.Tuple6[B, B, B, B, B, B], error) {
		ret, err := traverseTNModify(ctx, fmap, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F)
		return lo.T6(ret[0], ret[1], ret[2], ret[3], ret[4], ret[5]), err
	}, func(ctx context.Context, index int, source lo.Tuple6[A, A, A, A, A, A]) SeqIE[int, A] {
		return traverseTNIxGet(index, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F)
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.Traverse{
			OpticTypeExpr: ot,
		}
	}))
}

// TraverseT7 returns a [Traversal] focusing on the elements of a [lo.Tuple7]
//
// See: [TraverseT7P] for a polymorphic version
func TraverseT7[A any]() Optic[int, lo.Tuple7[A, A, A, A, A, A, A], lo.Tuple7[A, A, A, A, A, A, A], A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseT7P[A, A]()
}

// TraverseT7P returns a polymorphic [Traversal] focusing on the elements of a [lo.Tuple7]
//
// See: [TraverseT7] for a non polymorphic version
func TraverseT7P[A any, B any]() Optic[int, lo.Tuple7[A, A, A, A, A, A, A], lo.Tuple7[B, B, B, B, B, B, B], A, B, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, int, lo.Tuple7[A, A, A, A, A, A, A], lo.Tuple7[B, B, B, B, B, B, B], A, B](func(ctx context.Context, source lo.Tuple7[A, A, A, A, A, A, A]) SeqIE[int, A] {
		return traverseTNIter(ctx, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F, &source.G)
	}, func(ctx context.Context, source lo.Tuple7[A, A, A, A, A, A, A]) (int, error) {
		return 7, nil
	}, func(ctx context.Context, fmap func(index int, focus A) (B, error), source lo.Tuple7[A, A, A, A, A, A, A]) (lo.Tuple7[B, B, B, B, B, B, B], error) {
		ret, err := traverseTNModify(ctx, fmap, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F, &source.G)
		return lo.T7(ret[0], ret[1], ret[2], ret[3], ret[4], ret[5], ret[6]), err
	}, func(ctx context.Context, index int, source lo.Tuple7[A, A, A, A, A, A, A]) SeqIE[int, A] {
		return traverseTNIxGet(index, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F, &source.G)
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.Traverse{
			OpticTypeExpr: ot,
		}
	}))
}

// TraverseT8 returns a [Traversal] focusing on the elements of a [lo.Tuple8]
//
// See: [TraverseT8P] for a polymorphic version
func TraverseT8[A any]() Optic[int, lo.Tuple8[A, A, A, A, A, A, A, A], lo.Tuple8[A, A, A, A, A, A, A, A], A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseT8P[A, A]()
}

// TraverseT8P returns a polymorphic [Traversal] focusing on the elements of a [lo.Tuple8]
//
// See: [TraverseT8] for a non polymorphic version
func TraverseT8P[A any, B any]() Optic[int, lo.Tuple8[A, A, A, A, A, A, A, A], lo.Tuple8[B, B, B, B, B, B, B, B], A, B, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, int, lo.Tuple8[A, A, A, A, A, A, A, A], lo.Tuple8[B, B, B, B, B, B, B, B], A, B](func(ctx context.Context, source lo.Tuple8[A, A, A, A, A, A, A, A]) SeqIE[int, A] {
		return traverseTNIter(ctx, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F, &source.G, &source.H)
	}, func(ctx context.Context, source lo.Tuple8[A, A, A, A, A, A, A, A]) (int, error) {
		return 8, nil
	}, func(ctx context.Context, fmap func(index int, focus A) (B, error), source lo.Tuple8[A, A, A, A, A, A, A, A]) (lo.Tuple8[B, B, B, B, B, B, B, B], error) {
		ret, err := traverseTNModify(ctx, fmap, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F, &source.G, &source.H)
		return lo.T8(ret[0], ret[1], ret[2], ret[3], ret[4], ret[5], ret[6], ret[7]), err
	}, func(ctx context.Context, index int, source lo.Tuple8[A, A, A, A, A, A, A, A]) SeqIE[int, A] {
		return traverseTNIxGet(index, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F, &source.G, &source.H)
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.Traverse{
			OpticTypeExpr: ot,
		}
	}))
}

// TraverseT9 returns a [Traversal] focusing on the elements of a [lo.Tuple9]
//
// See: [TraverseT9P] for a polymorphic version
func TraverseT9[A any]() Optic[int, lo.Tuple9[A, A, A, A, A, A, A, A, A], lo.Tuple9[A, A, A, A, A, A, A, A, A], A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseT9P[A, A]()
}

// TraverseT9P returns a polymorphic [Traversal] focusing on the elements of a [lo.Tuple9]
//
// See: [TraverseT9] for a non polymorphic version
func TraverseT9P[A any, B any]() Optic[int, lo.Tuple9[A, A, A, A, A, A, A, A, A], lo.Tuple9[B, B, B, B, B, B, B, B, B], A, B, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, int, lo.Tuple9[A, A, A, A, A, A, A, A, A], lo.Tuple9[B, B, B, B, B, B, B, B, B], A, B](func(ctx context.Context, source lo.Tuple9[A, A, A, A, A, A, A, A, A]) SeqIE[int, A] {
		return traverseTNIter(ctx, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F, &source.G, &source.H, &source.I)
	}, func(ctx context.Context, source lo.Tuple9[A, A, A, A, A, A, A, A, A]) (int, error) {
		return 9, nil
	}, func(ctx context.Context, fmap func(index int, focus A) (B, error), source lo.Tuple9[A, A, A, A, A, A, A, A, A]) (lo.Tuple9[B, B, B, B, B, B, B, B, B], error) {
		ret, err := traverseTNModify(ctx, fmap, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F, &source.G, &source.H, &source.I)
		return lo.T9(ret[0], ret[1], ret[2], ret[3], ret[4], ret[5], ret[6], ret[7], ret[8]), err
	}, func(ctx context.Context, index int, source lo.Tuple9[A, A, A, A, A, A, A, A, A]) SeqIE[int, A] {
		return traverseTNIxGet(index, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F, &source.G, &source.H, &source.I)
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.Traverse{
			OpticTypeExpr: ot,
		}
	}))
}

// DupT2 returns a [Lens] focusing a [lo.Tuple2] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT2P] for a polymorphic version
func DupT2[A any]() Optic[Void, A, A, lo.Tuple2[A, A], lo.Tuple2[A, A], ReturnOne, ReadWrite, UniDir, Pure] {
	return DupT2P[A, A]()
}

// DupT2P returns a polymorphic [Lens] focusing a [lo.Tuple2] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT2] for a non polymorphic version
func DupT2P[A any, B any]() Optic[Void, A, B, lo.Tuple2[A, A], lo.Tuple2[B, B], ReturnOne, ReadWrite, UniDir, Pure] {
	return LensP[A, B, lo.Tuple2[A, A], lo.Tuple2[B, B]](func(source A) lo.Tuple2[A, A] {
		return lo.T2(source, source)
	}, func(focus lo.Tuple2[B, B], source A) B {
		return focus.A
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleDup{
			OpticTypeExpr: ot,
			N:             2,
		}
	}))
}

// DupT3 returns a [Lens] focusing a [lo.Tuple3] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT3P] for a polymorphic version
func DupT3[A any]() Optic[Void, A, A, lo.Tuple3[A, A, A], lo.Tuple3[A, A, A], ReturnOne, ReadWrite, UniDir, Pure] {
	return DupT3P[A, A]()
}

// DupT3P returns a polymorphic [Lens] focusing a [lo.Tuple3] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT3] for a non polymorphic version
func DupT3P[A any, B any]() Optic[Void, A, B, lo.Tuple3[A, A, A], lo.Tuple3[B, B, B], ReturnOne, ReadWrite, UniDir, Pure] {
	return LensP[A, B, lo.Tuple3[A, A, A], lo.Tuple3[B, B, B]](func(source A) lo.Tuple3[A, A, A] {
		return lo.T3(source, source, source)
	}, func(focus lo.Tuple3[B, B, B], source A) B {
		return focus.A
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleDup{
			OpticTypeExpr: ot,
			N:             3,
		}
	}))
}

// DupT4 returns a [Lens] focusing a [lo.Tuple4] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT4P] for a polymorphic version
func DupT4[A any]() Optic[Void, A, A, lo.Tuple4[A, A, A, A], lo.Tuple4[A, A, A, A], ReturnOne, ReadWrite, UniDir, Pure] {
	return DupT4P[A, A]()
}

// DupT4P returns a polymorphic [Lens] focusing a [lo.Tuple4] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT4] for a non polymorphic version
func DupT4P[A any, B any]() Optic[Void, A, B, lo.Tuple4[A, A, A, A], lo.Tuple4[B, B, B, B], ReturnOne, ReadWrite, UniDir, Pure] {
	return LensP[A, B, lo.Tuple4[A, A, A, A], lo.Tuple4[B, B, B, B]](func(source A) lo.Tuple4[A, A, A, A] {
		return lo.T4(source, source, source, source)
	}, func(focus lo.Tuple4[B, B, B, B], source A) B {
		return focus.A
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleDup{
			OpticTypeExpr: ot,
			N:             4,
		}
	}))
}

// DupT5 returns a [Lens] focusing a [lo.Tuple5] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT5P] for a polymorphic version
func DupT5[A any]() Optic[Void, A, A, lo.Tuple5[A, A, A, A, A], lo.Tuple5[A, A, A, A, A], ReturnOne, ReadWrite, UniDir, Pure] {
	return DupT5P[A, A]()
}

// DupT5P returns a polymorphic [Lens] focusing a [lo.Tuple5] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT5] for a non polymorphic version
func DupT5P[A any, B any]() Optic[Void, A, B, lo.Tuple5[A, A, A, A, A], lo.Tuple5[B, B, B, B, B], ReturnOne, ReadWrite, UniDir, Pure] {
	return LensP[A, B, lo.Tuple5[A, A, A, A, A], lo.Tuple5[B, B, B, B, B]](func(source A) lo.Tuple5[A, A, A, A, A] {
		return lo.T5(source, source, source, source, source)
	}, func(focus lo.Tuple5[B, B, B, B, B], source A) B {
		return focus.A
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleDup{
			OpticTypeExpr: ot,
			N:             5,
		}
	}))
}

// DupT6 returns a [Lens] focusing a [lo.Tuple6] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT6P] for a polymorphic version
func DupT6[A any]() Optic[Void, A, A, lo.Tuple6[A, A, A, A, A, A], lo.Tuple6[A, A, A, A, A, A], ReturnOne, ReadWrite, UniDir, Pure] {
	return DupT6P[A, A]()
}

// DupT6P returns a polymorphic [Lens] focusing a [lo.Tuple6] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT6] for a non polymorphic version
func DupT6P[A any, B any]() Optic[Void, A, B, lo.Tuple6[A, A, A, A, A, A], lo.Tuple6[B, B, B, B, B, B], ReturnOne, ReadWrite, UniDir, Pure] {
	return LensP[A, B, lo.Tuple6[A, A, A, A, A, A], lo.Tuple6[B, B, B, B, B, B]](func(source A) lo.Tuple6[A, A, A, A, A, A] {
		return lo.T6(source, source, source, source, source, source)
	}, func(focus lo.Tuple6[B, B, B, B, B, B], source A) B {
		return focus.A
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleDup{
			OpticTypeExpr: ot,
			N:             6,
		}
	}))
}

// DupT7 returns a [Lens] focusing a [lo.Tuple7] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT7P] for a polymorphic version
func DupT7[A any]() Optic[Void, A, A, lo.Tuple7[A, A, A, A, A, A, A], lo.Tuple7[A, A, A, A, A, A, A], ReturnOne, ReadWrite, UniDir, Pure] {
	return DupT7P[A, A]()
}

// DupT7P returns a polymorphic [Lens] focusing a [lo.Tuple7] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT7] for a non polymorphic version
func DupT7P[A any, B any]() Optic[Void, A, B, lo.Tuple7[A, A, A, A, A, A, A], lo.Tuple7[B, B, B, B, B, B, B], ReturnOne, ReadWrite, UniDir, Pure] {
	return LensP[A, B, lo.Tuple7[A, A, A, A, A, A, A], lo.Tuple7[B, B, B, B, B, B, B]](func(source A) lo.Tuple7[A, A, A, A, A, A, A] {
		return lo.T7(source, source, source, source, source, source, source)
	}, func(focus lo.Tuple7[B, B, B, B, B, B, B], source A) B {
		return focus.A
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleDup{
			OpticTypeExpr: ot,
			N:             7,
		}
	}))
}

// DupT8 returns a [Lens] focusing a [lo.Tuple8] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT8P] for a polymorphic version
func DupT8[A any]() Optic[Void, A, A, lo.Tuple8[A, A, A, A, A, A, A, A], lo.Tuple8[A, A, A, A, A, A, A, A], ReturnOne, ReadWrite, UniDir, Pure] {
	return DupT8P[A, A]()
}

// DupT8P returns a polymorphic [Lens] focusing a [lo.Tuple8] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT8] for a non polymorphic version
func DupT8P[A any, B any]() Optic[Void, A, B, lo.Tuple8[A, A, A, A, A, A, A, A], lo.Tuple8[B, B, B, B, B, B, B, B], ReturnOne, ReadWrite, UniDir, Pure] {
	return LensP[A, B, lo.Tuple8[A, A, A, A, A, A, A, A], lo.Tuple8[B, B, B, B, B, B, B, B]](func(source A) lo.Tuple8[A, A, A, A, A, A, A, A] {
		return lo.T8(source, source, source, source, source, source, source, source)
	}, func(focus lo.Tuple8[B, B, B, B, B, B, B, B], source A) B {
		return focus.A
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleDup{
			OpticTypeExpr: ot,
			N:             8,
		}
	}))
}

// DupT9 returns a [Lens] focusing a [lo.Tuple9] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT9P] for a polymorphic version
func DupT9[A any]() Optic[Void, A, A, lo.Tuple9[A, A, A, A, A, A, A, A, A], lo.Tuple9[A, A, A, A, A, A, A, A, A], ReturnOne, ReadWrite, UniDir, Pure] {
	return DupT9P[A, A]()
}

// DupT9P returns a polymorphic [Lens] focusing a [lo.Tuple9] with all elements set to the source value.
//
// Note: under modification the first element of the tuple is used. All other elements are ignored.
//
// See: [DupT9] for a non polymorphic version
func DupT9P[A any, B any]() Optic[Void, A, B, lo.Tuple9[A, A, A, A, A, A, A, A, A], lo.Tuple9[B, B, B, B, B, B, B, B, B], ReturnOne, ReadWrite, UniDir, Pure] {
	return LensP[A, B, lo.Tuple9[A, A, A, A, A, A, A, A, A], lo.Tuple9[B, B, B, B, B, B, B, B, B]](func(source A) lo.Tuple9[A, A, A, A, A, A, A, A, A] {
		return lo.T9(source, source, source, source, source, source, source, source, source)
	}, func(focus lo.Tuple9[B, B, B, B, B, B, B, B, B], source A) B {
		return focus.A
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleDup{
			OpticTypeExpr: ot,
			N:             9,
		}
	}))
}

// T2A returns a [Lens] focusing on element 0 of a [lo.Tuple2]
//
// See: [T2AP] for a polymorphic version
func T2A[A any, B any]() Optic[int, lo.Tuple2[A, B], lo.Tuple2[A, B], A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return T2AP[A, B, A]()
}

// T2AP returns a polymorphic [Lens] focusing on element 0 of a [lo.Tuple2]
//
// See: [T2A] for a non polymorphic version
func T2AP[A any, B any, A2 any]() Optic[int, lo.Tuple2[A, B], lo.Tuple2[A2, B], A, A2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple2[A, B], lo.Tuple2[A2, B], A, A2](func(ctx context.Context, source lo.Tuple2[A, B]) (int, A, error) {
		return 0, source.A, nil
	}, func(ctx context.Context, focus A2, source lo.Tuple2[A, B]) (lo.Tuple2[A2, B], error) {
		return lo.T2(focus, source.B), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         0,
		}
	}))
}

// T2B returns a [Lens] focusing on element 1 of a [lo.Tuple2]
//
// See: [T2BP] for a polymorphic version
func T2B[A any, B any]() Optic[int, lo.Tuple2[A, B], lo.Tuple2[A, B], B, B, ReturnOne, ReadWrite, UniDir, Pure] {
	return T2BP[A, B, B]()
}

// T2BP returns a polymorphic [Lens] focusing on element 1 of a [lo.Tuple2]
//
// See: [T2B] for a non polymorphic version
func T2BP[A any, B any, B2 any]() Optic[int, lo.Tuple2[A, B], lo.Tuple2[A, B2], B, B2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple2[A, B], lo.Tuple2[A, B2], B, B2](func(ctx context.Context, source lo.Tuple2[A, B]) (int, B, error) {
		return 1, source.B, nil
	}, func(ctx context.Context, focus B2, source lo.Tuple2[A, B]) (lo.Tuple2[A, B2], error) {
		return lo.T2(source.A, focus), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         1,
		}
	}))
}

// T3A returns a [Lens] focusing on element 0 of a [lo.Tuple3]
//
// See: [T3AP] for a polymorphic version
func T3A[A any, B any, C any]() Optic[int, lo.Tuple3[A, B, C], lo.Tuple3[A, B, C], A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return T3AP[A, B, C, A]()
}

// T3AP returns a polymorphic [Lens] focusing on element 0 of a [lo.Tuple3]
//
// See: [T3A] for a non polymorphic version
func T3AP[A any, B any, C any, A2 any]() Optic[int, lo.Tuple3[A, B, C], lo.Tuple3[A2, B, C], A, A2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple3[A, B, C], lo.Tuple3[A2, B, C], A, A2](func(ctx context.Context, source lo.Tuple3[A, B, C]) (int, A, error) {
		return 0, source.A, nil
	}, func(ctx context.Context, focus A2, source lo.Tuple3[A, B, C]) (lo.Tuple3[A2, B, C], error) {
		return lo.T3(focus, source.B, source.C), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         0,
		}
	}))
}

// T3B returns a [Lens] focusing on element 1 of a [lo.Tuple3]
//
// See: [T3BP] for a polymorphic version
func T3B[A any, B any, C any]() Optic[int, lo.Tuple3[A, B, C], lo.Tuple3[A, B, C], B, B, ReturnOne, ReadWrite, UniDir, Pure] {
	return T3BP[A, B, C, B]()
}

// T3BP returns a polymorphic [Lens] focusing on element 1 of a [lo.Tuple3]
//
// See: [T3B] for a non polymorphic version
func T3BP[A any, B any, C any, B2 any]() Optic[int, lo.Tuple3[A, B, C], lo.Tuple3[A, B2, C], B, B2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple3[A, B, C], lo.Tuple3[A, B2, C], B, B2](func(ctx context.Context, source lo.Tuple3[A, B, C]) (int, B, error) {
		return 1, source.B, nil
	}, func(ctx context.Context, focus B2, source lo.Tuple3[A, B, C]) (lo.Tuple3[A, B2, C], error) {
		return lo.T3(source.A, focus, source.C), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         1,
		}
	}))
}

// T3C returns a [Lens] focusing on element 2 of a [lo.Tuple3]
//
// See: [T3CP] for a polymorphic version
func T3C[A any, B any, C any]() Optic[int, lo.Tuple3[A, B, C], lo.Tuple3[A, B, C], C, C, ReturnOne, ReadWrite, UniDir, Pure] {
	return T3CP[A, B, C, C]()
}

// T3CP returns a polymorphic [Lens] focusing on element 2 of a [lo.Tuple3]
//
// See: [T3C] for a non polymorphic version
func T3CP[A any, B any, C any, C2 any]() Optic[int, lo.Tuple3[A, B, C], lo.Tuple3[A, B, C2], C, C2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple3[A, B, C], lo.Tuple3[A, B, C2], C, C2](func(ctx context.Context, source lo.Tuple3[A, B, C]) (int, C, error) {
		return 2, source.C, nil
	}, func(ctx context.Context, focus C2, source lo.Tuple3[A, B, C]) (lo.Tuple3[A, B, C2], error) {
		return lo.T3(source.A, source.B, focus), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         2,
		}
	}))
}

// T4A returns a [Lens] focusing on element 0 of a [lo.Tuple4]
//
// See: [T4AP] for a polymorphic version
func T4A[A any, B any, C any, D any]() Optic[int, lo.Tuple4[A, B, C, D], lo.Tuple4[A, B, C, D], A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return T4AP[A, B, C, D, A]()
}

// T4AP returns a polymorphic [Lens] focusing on element 0 of a [lo.Tuple4]
//
// See: [T4A] for a non polymorphic version
func T4AP[A any, B any, C any, D any, A2 any]() Optic[int, lo.Tuple4[A, B, C, D], lo.Tuple4[A2, B, C, D], A, A2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple4[A, B, C, D], lo.Tuple4[A2, B, C, D], A, A2](func(ctx context.Context, source lo.Tuple4[A, B, C, D]) (int, A, error) {
		return 0, source.A, nil
	}, func(ctx context.Context, focus A2, source lo.Tuple4[A, B, C, D]) (lo.Tuple4[A2, B, C, D], error) {
		return lo.T4(focus, source.B, source.C, source.D), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         0,
		}
	}))
}

// T4B returns a [Lens] focusing on element 1 of a [lo.Tuple4]
//
// See: [T4BP] for a polymorphic version
func T4B[A any, B any, C any, D any]() Optic[int, lo.Tuple4[A, B, C, D], lo.Tuple4[A, B, C, D], B, B, ReturnOne, ReadWrite, UniDir, Pure] {
	return T4BP[A, B, C, D, B]()
}

// T4BP returns a polymorphic [Lens] focusing on element 1 of a [lo.Tuple4]
//
// See: [T4B] for a non polymorphic version
func T4BP[A any, B any, C any, D any, B2 any]() Optic[int, lo.Tuple4[A, B, C, D], lo.Tuple4[A, B2, C, D], B, B2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple4[A, B, C, D], lo.Tuple4[A, B2, C, D], B, B2](func(ctx context.Context, source lo.Tuple4[A, B, C, D]) (int, B, error) {
		return 1, source.B, nil
	}, func(ctx context.Context, focus B2, source lo.Tuple4[A, B, C, D]) (lo.Tuple4[A, B2, C, D], error) {
		return lo.T4(source.A, focus, source.C, source.D), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         1,
		}
	}))
}

// T4C returns a [Lens] focusing on element 2 of a [lo.Tuple4]
//
// See: [T4CP] for a polymorphic version
func T4C[A any, B any, C any, D any]() Optic[int, lo.Tuple4[A, B, C, D], lo.Tuple4[A, B, C, D], C, C, ReturnOne, ReadWrite, UniDir, Pure] {
	return T4CP[A, B, C, D, C]()
}

// T4CP returns a polymorphic [Lens] focusing on element 2 of a [lo.Tuple4]
//
// See: [T4C] for a non polymorphic version
func T4CP[A any, B any, C any, D any, C2 any]() Optic[int, lo.Tuple4[A, B, C, D], lo.Tuple4[A, B, C2, D], C, C2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple4[A, B, C, D], lo.Tuple4[A, B, C2, D], C, C2](func(ctx context.Context, source lo.Tuple4[A, B, C, D]) (int, C, error) {
		return 2, source.C, nil
	}, func(ctx context.Context, focus C2, source lo.Tuple4[A, B, C, D]) (lo.Tuple4[A, B, C2, D], error) {
		return lo.T4(source.A, source.B, focus, source.D), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         2,
		}
	}))
}

// T4D returns a [Lens] focusing on element 3 of a [lo.Tuple4]
//
// See: [T4DP] for a polymorphic version
func T4D[A any, B any, C any, D any]() Optic[int, lo.Tuple4[A, B, C, D], lo.Tuple4[A, B, C, D], D, D, ReturnOne, ReadWrite, UniDir, Pure] {
	return T4DP[A, B, C, D, D]()
}

// T4DP returns a polymorphic [Lens] focusing on element 3 of a [lo.Tuple4]
//
// See: [T4D] for a non polymorphic version
func T4DP[A any, B any, C any, D any, D2 any]() Optic[int, lo.Tuple4[A, B, C, D], lo.Tuple4[A, B, C, D2], D, D2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple4[A, B, C, D], lo.Tuple4[A, B, C, D2], D, D2](func(ctx context.Context, source lo.Tuple4[A, B, C, D]) (int, D, error) {
		return 3, source.D, nil
	}, func(ctx context.Context, focus D2, source lo.Tuple4[A, B, C, D]) (lo.Tuple4[A, B, C, D2], error) {
		return lo.T4(source.A, source.B, source.C, focus), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         3,
		}
	}))
}

// T5A returns a [Lens] focusing on element 0 of a [lo.Tuple5]
//
// See: [T5AP] for a polymorphic version
func T5A[A any, B any, C any, D any, E any]() Optic[int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A, B, C, D, E], A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return T5AP[A, B, C, D, E, A]()
}

// T5AP returns a polymorphic [Lens] focusing on element 0 of a [lo.Tuple5]
//
// See: [T5A] for a non polymorphic version
func T5AP[A any, B any, C any, D any, E any, A2 any]() Optic[int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A2, B, C, D, E], A, A2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A2, B, C, D, E], A, A2](func(ctx context.Context, source lo.Tuple5[A, B, C, D, E]) (int, A, error) {
		return 0, source.A, nil
	}, func(ctx context.Context, focus A2, source lo.Tuple5[A, B, C, D, E]) (lo.Tuple5[A2, B, C, D, E], error) {
		return lo.T5(focus, source.B, source.C, source.D, source.E), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         0,
		}
	}))
}

// T5B returns a [Lens] focusing on element 1 of a [lo.Tuple5]
//
// See: [T5BP] for a polymorphic version
func T5B[A any, B any, C any, D any, E any]() Optic[int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A, B, C, D, E], B, B, ReturnOne, ReadWrite, UniDir, Pure] {
	return T5BP[A, B, C, D, E, B]()
}

// T5BP returns a polymorphic [Lens] focusing on element 1 of a [lo.Tuple5]
//
// See: [T5B] for a non polymorphic version
func T5BP[A any, B any, C any, D any, E any, B2 any]() Optic[int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A, B2, C, D, E], B, B2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A, B2, C, D, E], B, B2](func(ctx context.Context, source lo.Tuple5[A, B, C, D, E]) (int, B, error) {
		return 1, source.B, nil
	}, func(ctx context.Context, focus B2, source lo.Tuple5[A, B, C, D, E]) (lo.Tuple5[A, B2, C, D, E], error) {
		return lo.T5(source.A, focus, source.C, source.D, source.E), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         1,
		}
	}))
}

// T5C returns a [Lens] focusing on element 2 of a [lo.Tuple5]
//
// See: [T5CP] for a polymorphic version
func T5C[A any, B any, C any, D any, E any]() Optic[int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A, B, C, D, E], C, C, ReturnOne, ReadWrite, UniDir, Pure] {
	return T5CP[A, B, C, D, E, C]()
}

// T5CP returns a polymorphic [Lens] focusing on element 2 of a [lo.Tuple5]
//
// See: [T5C] for a non polymorphic version
func T5CP[A any, B any, C any, D any, E any, C2 any]() Optic[int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A, B, C2, D, E], C, C2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A, B, C2, D, E], C, C2](func(ctx context.Context, source lo.Tuple5[A, B, C, D, E]) (int, C, error) {
		return 2, source.C, nil
	}, func(ctx context.Context, focus C2, source lo.Tuple5[A, B, C, D, E]) (lo.Tuple5[A, B, C2, D, E], error) {
		return lo.T5(source.A, source.B, focus, source.D, source.E), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         2,
		}
	}))
}

// T5D returns a [Lens] focusing on element 3 of a [lo.Tuple5]
//
// See: [T5DP] for a polymorphic version
func T5D[A any, B any, C any, D any, E any]() Optic[int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A, B, C, D, E], D, D, ReturnOne, ReadWrite, UniDir, Pure] {
	return T5DP[A, B, C, D, E, D]()
}

// T5DP returns a polymorphic [Lens] focusing on element 3 of a [lo.Tuple5]
//
// See: [T5D] for a non polymorphic version
func T5DP[A any, B any, C any, D any, E any, D2 any]() Optic[int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A, B, C, D2, E], D, D2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A, B, C, D2, E], D, D2](func(ctx context.Context, source lo.Tuple5[A, B, C, D, E]) (int, D, error) {
		return 3, source.D, nil
	}, func(ctx context.Context, focus D2, source lo.Tuple5[A, B, C, D, E]) (lo.Tuple5[A, B, C, D2, E], error) {
		return lo.T5(source.A, source.B, source.C, focus, source.E), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         3,
		}
	}))
}

// T5E returns a [Lens] focusing on element 4 of a [lo.Tuple5]
//
// See: [T5EP] for a polymorphic version
func T5E[A any, B any, C any, D any, E any]() Optic[int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A, B, C, D, E], E, E, ReturnOne, ReadWrite, UniDir, Pure] {
	return T5EP[A, B, C, D, E, E]()
}

// T5EP returns a polymorphic [Lens] focusing on element 4 of a [lo.Tuple5]
//
// See: [T5E] for a non polymorphic version
func T5EP[A any, B any, C any, D any, E any, E2 any]() Optic[int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A, B, C, D, E2], E, E2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple5[A, B, C, D, E], lo.Tuple5[A, B, C, D, E2], E, E2](func(ctx context.Context, source lo.Tuple5[A, B, C, D, E]) (int, E, error) {
		return 4, source.E, nil
	}, func(ctx context.Context, focus E2, source lo.Tuple5[A, B, C, D, E]) (lo.Tuple5[A, B, C, D, E2], error) {
		return lo.T5(source.A, source.B, source.C, source.D, focus), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         4,
		}
	}))
}

// T6A returns a [Lens] focusing on element 0 of a [lo.Tuple6]
//
// See: [T6AP] for a polymorphic version
func T6A[A any, B any, C any, D any, E any, F any]() Optic[int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C, D, E, F], A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return T6AP[A, B, C, D, E, F, A]()
}

// T6AP returns a polymorphic [Lens] focusing on element 0 of a [lo.Tuple6]
//
// See: [T6A] for a non polymorphic version
func T6AP[A any, B any, C any, D any, E any, F any, A2 any]() Optic[int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A2, B, C, D, E, F], A, A2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A2, B, C, D, E, F], A, A2](func(ctx context.Context, source lo.Tuple6[A, B, C, D, E, F]) (int, A, error) {
		return 0, source.A, nil
	}, func(ctx context.Context, focus A2, source lo.Tuple6[A, B, C, D, E, F]) (lo.Tuple6[A2, B, C, D, E, F], error) {
		return lo.T6(focus, source.B, source.C, source.D, source.E, source.F), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         0,
		}
	}))
}

// T6B returns a [Lens] focusing on element 1 of a [lo.Tuple6]
//
// See: [T6BP] for a polymorphic version
func T6B[A any, B any, C any, D any, E any, F any]() Optic[int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C, D, E, F], B, B, ReturnOne, ReadWrite, UniDir, Pure] {
	return T6BP[A, B, C, D, E, F, B]()
}

// T6BP returns a polymorphic [Lens] focusing on element 1 of a [lo.Tuple6]
//
// See: [T6B] for a non polymorphic version
func T6BP[A any, B any, C any, D any, E any, F any, B2 any]() Optic[int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B2, C, D, E, F], B, B2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B2, C, D, E, F], B, B2](func(ctx context.Context, source lo.Tuple6[A, B, C, D, E, F]) (int, B, error) {
		return 1, source.B, nil
	}, func(ctx context.Context, focus B2, source lo.Tuple6[A, B, C, D, E, F]) (lo.Tuple6[A, B2, C, D, E, F], error) {
		return lo.T6(source.A, focus, source.C, source.D, source.E, source.F), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         1,
		}
	}))
}

// T6C returns a [Lens] focusing on element 2 of a [lo.Tuple6]
//
// See: [T6CP] for a polymorphic version
func T6C[A any, B any, C any, D any, E any, F any]() Optic[int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C, D, E, F], C, C, ReturnOne, ReadWrite, UniDir, Pure] {
	return T6CP[A, B, C, D, E, F, C]()
}

// T6CP returns a polymorphic [Lens] focusing on element 2 of a [lo.Tuple6]
//
// See: [T6C] for a non polymorphic version
func T6CP[A any, B any, C any, D any, E any, F any, C2 any]() Optic[int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C2, D, E, F], C, C2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C2, D, E, F], C, C2](func(ctx context.Context, source lo.Tuple6[A, B, C, D, E, F]) (int, C, error) {
		return 2, source.C, nil
	}, func(ctx context.Context, focus C2, source lo.Tuple6[A, B, C, D, E, F]) (lo.Tuple6[A, B, C2, D, E, F], error) {
		return lo.T6(source.A, source.B, focus, source.D, source.E, source.F), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         2,
		}
	}))
}

// T6D returns a [Lens] focusing on element 3 of a [lo.Tuple6]
//
// See: [T6DP] for a polymorphic version
func T6D[A any, B any, C any, D any, E any, F any]() Optic[int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C, D, E, F], D, D, ReturnOne, ReadWrite, UniDir, Pure] {
	return T6DP[A, B, C, D, E, F, D]()
}

// T6DP returns a polymorphic [Lens] focusing on element 3 of a [lo.Tuple6]
//
// See: [T6D] for a non polymorphic version
func T6DP[A any, B any, C any, D any, E any, F any, D2 any]() Optic[int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C, D2, E, F], D, D2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C, D2, E, F], D, D2](func(ctx context.Context, source lo.Tuple6[A, B, C, D, E, F]) (int, D, error) {
		return 3, source.D, nil
	}, func(ctx context.Context, focus D2, source lo.Tuple6[A, B, C, D, E, F]) (lo.Tuple6[A, B, C, D2, E, F], error) {
		return lo.T6(source.A, source.B, source.C, focus, source.E, source.F), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         3,
		}
	}))
}

// T6E returns a [Lens] focusing on element 4 of a [lo.Tuple6]
//
// See: [T6EP] for a polymorphic version
func T6E[A any, B any, C any, D any, E any, F any]() Optic[int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C, D, E, F], E, E, ReturnOne, ReadWrite, UniDir, Pure] {
	return T6EP[A, B, C, D, E, F, E]()
}

// T6EP returns a polymorphic [Lens] focusing on element 4 of a [lo.Tuple6]
//
// See: [T6E] for a non polymorphic version
func T6EP[A any, B any, C any, D any, E any, F any, E2 any]() Optic[int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C, D, E2, F], E, E2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C, D, E2, F], E, E2](func(ctx context.Context, source lo.Tuple6[A, B, C, D, E, F]) (int, E, error) {
		return 4, source.E, nil
	}, func(ctx context.Context, focus E2, source lo.Tuple6[A, B, C, D, E, F]) (lo.Tuple6[A, B, C, D, E2, F], error) {
		return lo.T6(source.A, source.B, source.C, source.D, focus, source.F), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         4,
		}
	}))
}

// T6F returns a [Lens] focusing on element 5 of a [lo.Tuple6]
//
// See: [T6FP] for a polymorphic version
func T6F[A any, B any, C any, D any, E any, F any]() Optic[int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C, D, E, F], F, F, ReturnOne, ReadWrite, UniDir, Pure] {
	return T6FP[A, B, C, D, E, F, F]()
}

// T6FP returns a polymorphic [Lens] focusing on element 5 of a [lo.Tuple6]
//
// See: [T6F] for a non polymorphic version
func T6FP[A any, B any, C any, D any, E any, F any, F2 any]() Optic[int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C, D, E, F2], F, F2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple6[A, B, C, D, E, F], lo.Tuple6[A, B, C, D, E, F2], F, F2](func(ctx context.Context, source lo.Tuple6[A, B, C, D, E, F]) (int, F, error) {
		return 5, source.F, nil
	}, func(ctx context.Context, focus F2, source lo.Tuple6[A, B, C, D, E, F]) (lo.Tuple6[A, B, C, D, E, F2], error) {
		return lo.T6(source.A, source.B, source.C, source.D, source.E, focus), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         5,
		}
	}))
}

// T7A returns a [Lens] focusing on element 0 of a [lo.Tuple7]
//
// See: [T7AP] for a polymorphic version
func T7A[A any, B any, C any, D any, E any, F any, G any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D, E, F, G], A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return T7AP[A, B, C, D, E, F, G, A]()
}

// T7AP returns a polymorphic [Lens] focusing on element 0 of a [lo.Tuple7]
//
// See: [T7A] for a non polymorphic version
func T7AP[A any, B any, C any, D any, E any, F any, G any, A2 any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A2, B, C, D, E, F, G], A, A2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A2, B, C, D, E, F, G], A, A2](func(ctx context.Context, source lo.Tuple7[A, B, C, D, E, F, G]) (int, A, error) {
		return 0, source.A, nil
	}, func(ctx context.Context, focus A2, source lo.Tuple7[A, B, C, D, E, F, G]) (lo.Tuple7[A2, B, C, D, E, F, G], error) {
		return lo.T7(focus, source.B, source.C, source.D, source.E, source.F, source.G), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         0,
		}
	}))
}

// T7B returns a [Lens] focusing on element 1 of a [lo.Tuple7]
//
// See: [T7BP] for a polymorphic version
func T7B[A any, B any, C any, D any, E any, F any, G any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D, E, F, G], B, B, ReturnOne, ReadWrite, UniDir, Pure] {
	return T7BP[A, B, C, D, E, F, G, B]()
}

// T7BP returns a polymorphic [Lens] focusing on element 1 of a [lo.Tuple7]
//
// See: [T7B] for a non polymorphic version
func T7BP[A any, B any, C any, D any, E any, F any, G any, B2 any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B2, C, D, E, F, G], B, B2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B2, C, D, E, F, G], B, B2](func(ctx context.Context, source lo.Tuple7[A, B, C, D, E, F, G]) (int, B, error) {
		return 1, source.B, nil
	}, func(ctx context.Context, focus B2, source lo.Tuple7[A, B, C, D, E, F, G]) (lo.Tuple7[A, B2, C, D, E, F, G], error) {
		return lo.T7(source.A, focus, source.C, source.D, source.E, source.F, source.G), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         1,
		}
	}))
}

// T7C returns a [Lens] focusing on element 2 of a [lo.Tuple7]
//
// See: [T7CP] for a polymorphic version
func T7C[A any, B any, C any, D any, E any, F any, G any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D, E, F, G], C, C, ReturnOne, ReadWrite, UniDir, Pure] {
	return T7CP[A, B, C, D, E, F, G, C]()
}

// T7CP returns a polymorphic [Lens] focusing on element 2 of a [lo.Tuple7]
//
// See: [T7C] for a non polymorphic version
func T7CP[A any, B any, C any, D any, E any, F any, G any, C2 any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C2, D, E, F, G], C, C2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C2, D, E, F, G], C, C2](func(ctx context.Context, source lo.Tuple7[A, B, C, D, E, F, G]) (int, C, error) {
		return 2, source.C, nil
	}, func(ctx context.Context, focus C2, source lo.Tuple7[A, B, C, D, E, F, G]) (lo.Tuple7[A, B, C2, D, E, F, G], error) {
		return lo.T7(source.A, source.B, focus, source.D, source.E, source.F, source.G), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         2,
		}
	}))
}

// T7D returns a [Lens] focusing on element 3 of a [lo.Tuple7]
//
// See: [T7DP] for a polymorphic version
func T7D[A any, B any, C any, D any, E any, F any, G any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D, E, F, G], D, D, ReturnOne, ReadWrite, UniDir, Pure] {
	return T7DP[A, B, C, D, E, F, G, D]()
}

// T7DP returns a polymorphic [Lens] focusing on element 3 of a [lo.Tuple7]
//
// See: [T7D] for a non polymorphic version
func T7DP[A any, B any, C any, D any, E any, F any, G any, D2 any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D2, E, F, G], D, D2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D2, E, F, G], D, D2](func(ctx context.Context, source lo.Tuple7[A, B, C, D, E, F, G]) (int, D, error) {
		return 3, source.D, nil
	}, func(ctx context.Context, focus D2, source lo.Tuple7[A, B, C, D, E, F, G]) (lo.Tuple7[A, B, C, D2, E, F, G], error) {
		return lo.T7(source.A, source.B, source.C, focus, source.E, source.F, source.G), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         3,
		}
	}))
}

// T7E returns a [Lens] focusing on element 4 of a [lo.Tuple7]
//
// See: [T7EP] for a polymorphic version
func T7E[A any, B any, C any, D any, E any, F any, G any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D, E, F, G], E, E, ReturnOne, ReadWrite, UniDir, Pure] {
	return T7EP[A, B, C, D, E, F, G, E]()
}

// T7EP returns a polymorphic [Lens] focusing on element 4 of a [lo.Tuple7]
//
// See: [T7E] for a non polymorphic version
func T7EP[A any, B any, C any, D any, E any, F any, G any, E2 any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D, E2, F, G], E, E2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D, E2, F, G], E, E2](func(ctx context.Context, source lo.Tuple7[A, B, C, D, E, F, G]) (int, E, error) {
		return 4, source.E, nil
	}, func(ctx context.Context, focus E2, source lo.Tuple7[A, B, C, D, E, F, G]) (lo.Tuple7[A, B, C, D, E2, F, G], error) {
		return lo.T7(source.A, source.B, source.C, source.D, focus, source.F, source.G), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         4,
		}
	}))
}

// T7F returns a [Lens] focusing on element 5 of a [lo.Tuple7]
//
// See: [T7FP] for a polymorphic version
func T7F[A any, B any, C any, D any, E any, F any, G any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D, E, F, G], F, F, ReturnOne, ReadWrite, UniDir, Pure] {
	return T7FP[A, B, C, D, E, F, G, F]()
}

// T7FP returns a polymorphic [Lens] focusing on element 5 of a [lo.Tuple7]
//
// See: [T7F] for a non polymorphic version
func T7FP[A any, B any, C any, D any, E any, F any, G any, F2 any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D, E, F2, G], F, F2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D, E, F2, G], F, F2](func(ctx context.Context, source lo.Tuple7[A, B, C, D, E, F, G]) (int, F, error) {
		return 5, source.F, nil
	}, func(ctx context.Context, focus F2, source lo.Tuple7[A, B, C, D, E, F, G]) (lo.Tuple7[A, B, C, D, E, F2, G], error) {
		return lo.T7(source.A, source.B, source.C, source.D, source.E, focus, source.G), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         5,
		}
	}))
}

// T7G returns a [Lens] focusing on element 6 of a [lo.Tuple7]
//
// See: [T7GP] for a polymorphic version
func T7G[A any, B any, C any, D any, E any, F any, G any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D, E, F, G], G, G, ReturnOne, ReadWrite, UniDir, Pure] {
	return T7GP[A, B, C, D, E, F, G, G]()
}

// T7GP returns a polymorphic [Lens] focusing on element 6 of a [lo.Tuple7]
//
// See: [T7G] for a non polymorphic version
func T7GP[A any, B any, C any, D any, E any, F any, G any, G2 any]() Optic[int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D, E, F, G2], G, G2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple7[A, B, C, D, E, F, G], lo.Tuple7[A, B, C, D, E, F, G2], G, G2](func(ctx context.Context, source lo.Tuple7[A, B, C, D, E, F, G]) (int, G, error) {
		return 6, source.G, nil
	}, func(ctx context.Context, focus G2, source lo.Tuple7[A, B, C, D, E, F, G]) (lo.Tuple7[A, B, C, D, E, F, G2], error) {
		return lo.T7(source.A, source.B, source.C, source.D, source.E, source.F, focus), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         6,
		}
	}))
}

// T8A returns a [Lens] focusing on element 0 of a [lo.Tuple8]
//
// See: [T8AP] for a polymorphic version
func T8A[A any, B any, C any, D any, E any, F any, G any, H any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F, G, H], A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return T8AP[A, B, C, D, E, F, G, H, A]()
}

// T8AP returns a polymorphic [Lens] focusing on element 0 of a [lo.Tuple8]
//
// See: [T8A] for a non polymorphic version
func T8AP[A any, B any, C any, D any, E any, F any, G any, H any, A2 any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A2, B, C, D, E, F, G, H], A, A2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A2, B, C, D, E, F, G, H], A, A2](func(ctx context.Context, source lo.Tuple8[A, B, C, D, E, F, G, H]) (int, A, error) {
		return 0, source.A, nil
	}, func(ctx context.Context, focus A2, source lo.Tuple8[A, B, C, D, E, F, G, H]) (lo.Tuple8[A2, B, C, D, E, F, G, H], error) {
		return lo.T8(focus, source.B, source.C, source.D, source.E, source.F, source.G, source.H), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         0,
		}
	}))
}

// T8B returns a [Lens] focusing on element 1 of a [lo.Tuple8]
//
// See: [T8BP] for a polymorphic version
func T8B[A any, B any, C any, D any, E any, F any, G any, H any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F, G, H], B, B, ReturnOne, ReadWrite, UniDir, Pure] {
	return T8BP[A, B, C, D, E, F, G, H, B]()
}

// T8BP returns a polymorphic [Lens] focusing on element 1 of a [lo.Tuple8]
//
// See: [T8B] for a non polymorphic version
func T8BP[A any, B any, C any, D any, E any, F any, G any, H any, B2 any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B2, C, D, E, F, G, H], B, B2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B2, C, D, E, F, G, H], B, B2](func(ctx context.Context, source lo.Tuple8[A, B, C, D, E, F, G, H]) (int, B, error) {
		return 1, source.B, nil
	}, func(ctx context.Context, focus B2, source lo.Tuple8[A, B, C, D, E, F, G, H]) (lo.Tuple8[A, B2, C, D, E, F, G, H], error) {
		return lo.T8(source.A, focus, source.C, source.D, source.E, source.F, source.G, source.H), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         1,
		}
	}))
}

// T8C returns a [Lens] focusing on element 2 of a [lo.Tuple8]
//
// See: [T8CP] for a polymorphic version
func T8C[A any, B any, C any, D any, E any, F any, G any, H any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F, G, H], C, C, ReturnOne, ReadWrite, UniDir, Pure] {
	return T8CP[A, B, C, D, E, F, G, H, C]()
}

// T8CP returns a polymorphic [Lens] focusing on element 2 of a [lo.Tuple8]
//
// See: [T8C] for a non polymorphic version
func T8CP[A any, B any, C any, D any, E any, F any, G any, H any, C2 any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C2, D, E, F, G, H], C, C2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C2, D, E, F, G, H], C, C2](func(ctx context.Context, source lo.Tuple8[A, B, C, D, E, F, G, H]) (int, C, error) {
		return 2, source.C, nil
	}, func(ctx context.Context, focus C2, source lo.Tuple8[A, B, C, D, E, F, G, H]) (lo.Tuple8[A, B, C2, D, E, F, G, H], error) {
		return lo.T8(source.A, source.B, focus, source.D, source.E, source.F, source.G, source.H), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         2,
		}
	}))
}

// T8D returns a [Lens] focusing on element 3 of a [lo.Tuple8]
//
// See: [T8DP] for a polymorphic version
func T8D[A any, B any, C any, D any, E any, F any, G any, H any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F, G, H], D, D, ReturnOne, ReadWrite, UniDir, Pure] {
	return T8DP[A, B, C, D, E, F, G, H, D]()
}

// T8DP returns a polymorphic [Lens] focusing on element 3 of a [lo.Tuple8]
//
// See: [T8D] for a non polymorphic version
func T8DP[A any, B any, C any, D any, E any, F any, G any, H any, D2 any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D2, E, F, G, H], D, D2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D2, E, F, G, H], D, D2](func(ctx context.Context, source lo.Tuple8[A, B, C, D, E, F, G, H]) (int, D, error) {
		return 3, source.D, nil
	}, func(ctx context.Context, focus D2, source lo.Tuple8[A, B, C, D, E, F, G, H]) (lo.Tuple8[A, B, C, D2, E, F, G, H], error) {
		return lo.T8(source.A, source.B, source.C, focus, source.E, source.F, source.G, source.H), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         3,
		}
	}))
}

// T8E returns a [Lens] focusing on element 4 of a [lo.Tuple8]
//
// See: [T8EP] for a polymorphic version
func T8E[A any, B any, C any, D any, E any, F any, G any, H any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F, G, H], E, E, ReturnOne, ReadWrite, UniDir, Pure] {
	return T8EP[A, B, C, D, E, F, G, H, E]()
}

// T8EP returns a polymorphic [Lens] focusing on element 4 of a [lo.Tuple8]
//
// See: [T8E] for a non polymorphic version
func T8EP[A any, B any, C any, D any, E any, F any, G any, H any, E2 any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E2, F, G, H], E, E2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E2, F, G, H], E, E2](func(ctx context.Context, source lo.Tuple8[A, B, C, D, E, F, G, H]) (int, E, error) {
		return 4, source.E, nil
	}, func(ctx context.Context, focus E2, source lo.Tuple8[A, B, C, D, E, F, G, H]) (lo.Tuple8[A, B, C, D, E2, F, G, H], error) {
		return lo.T8(source.A, source.B, source.C, source.D, focus, source.F, source.G, source.H), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         4,
		}
	}))
}

// T8F returns a [Lens] focusing on element 5 of a [lo.Tuple8]
//
// See: [T8FP] for a polymorphic version
func T8F[A any, B any, C any, D any, E any, F any, G any, H any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F, G, H], F, F, ReturnOne, ReadWrite, UniDir, Pure] {
	return T8FP[A, B, C, D, E, F, G, H, F]()
}

// T8FP returns a polymorphic [Lens] focusing on element 5 of a [lo.Tuple8]
//
// See: [T8F] for a non polymorphic version
func T8FP[A any, B any, C any, D any, E any, F any, G any, H any, F2 any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F2, G, H], F, F2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F2, G, H], F, F2](func(ctx context.Context, source lo.Tuple8[A, B, C, D, E, F, G, H]) (int, F, error) {
		return 5, source.F, nil
	}, func(ctx context.Context, focus F2, source lo.Tuple8[A, B, C, D, E, F, G, H]) (lo.Tuple8[A, B, C, D, E, F2, G, H], error) {
		return lo.T8(source.A, source.B, source.C, source.D, source.E, focus, source.G, source.H), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         5,
		}
	}))
}

// T8G returns a [Lens] focusing on element 6 of a [lo.Tuple8]
//
// See: [T8GP] for a polymorphic version
func T8G[A any, B any, C any, D any, E any, F any, G any, H any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F, G, H], G, G, ReturnOne, ReadWrite, UniDir, Pure] {
	return T8GP[A, B, C, D, E, F, G, H, G]()
}

// T8GP returns a polymorphic [Lens] focusing on element 6 of a [lo.Tuple8]
//
// See: [T8G] for a non polymorphic version
func T8GP[A any, B any, C any, D any, E any, F any, G any, H any, G2 any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F, G2, H], G, G2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F, G2, H], G, G2](func(ctx context.Context, source lo.Tuple8[A, B, C, D, E, F, G, H]) (int, G, error) {
		return 6, source.G, nil
	}, func(ctx context.Context, focus G2, source lo.Tuple8[A, B, C, D, E, F, G, H]) (lo.Tuple8[A, B, C, D, E, F, G2, H], error) {
		return lo.T8(source.A, source.B, source.C, source.D, source.E, source.F, focus, source.H), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         6,
		}
	}))
}

// T8H returns a [Lens] focusing on element 7 of a [lo.Tuple8]
//
// See: [T8HP] for a polymorphic version
func T8H[A any, B any, C any, D any, E any, F any, G any, H any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F, G, H], H, H, ReturnOne, ReadWrite, UniDir, Pure] {
	return T8HP[A, B, C, D, E, F, G, H, H]()
}

// T8HP returns a polymorphic [Lens] focusing on element 7 of a [lo.Tuple8]
//
// See: [T8H] for a non polymorphic version
func T8HP[A any, B any, C any, D any, E any, F any, G any, H any, H2 any]() Optic[int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F, G, H2], H, H2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple8[A, B, C, D, E, F, G, H], lo.Tuple8[A, B, C, D, E, F, G, H2], H, H2](func(ctx context.Context, source lo.Tuple8[A, B, C, D, E, F, G, H]) (int, H, error) {
		return 7, source.H, nil
	}, func(ctx context.Context, focus H2, source lo.Tuple8[A, B, C, D, E, F, G, H]) (lo.Tuple8[A, B, C, D, E, F, G, H2], error) {
		return lo.T8(source.A, source.B, source.C, source.D, source.E, source.F, source.G, focus), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         7,
		}
	}))
}

// T9A returns a [Lens] focusing on element 0 of a [lo.Tuple9]
//
// See: [T9AP] for a polymorphic version
func T9A[A any, B any, C any, D any, E any, F any, G any, H any, I any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G, H, I], A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return T9AP[A, B, C, D, E, F, G, H, I, A]()
}

// T9AP returns a polymorphic [Lens] focusing on element 0 of a [lo.Tuple9]
//
// See: [T9A] for a non polymorphic version
func T9AP[A any, B any, C any, D any, E any, F any, G any, H any, I any, A2 any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A2, B, C, D, E, F, G, H, I], A, A2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A2, B, C, D, E, F, G, H, I], A, A2](func(ctx context.Context, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (int, A, error) {
		return 0, source.A, nil
	}, func(ctx context.Context, focus A2, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (lo.Tuple9[A2, B, C, D, E, F, G, H, I], error) {
		return lo.T9(focus, source.B, source.C, source.D, source.E, source.F, source.G, source.H, source.I), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         0,
		}
	}))
}

// T9B returns a [Lens] focusing on element 1 of a [lo.Tuple9]
//
// See: [T9BP] for a polymorphic version
func T9B[A any, B any, C any, D any, E any, F any, G any, H any, I any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G, H, I], B, B, ReturnOne, ReadWrite, UniDir, Pure] {
	return T9BP[A, B, C, D, E, F, G, H, I, B]()
}

// T9BP returns a polymorphic [Lens] focusing on element 1 of a [lo.Tuple9]
//
// See: [T9B] for a non polymorphic version
func T9BP[A any, B any, C any, D any, E any, F any, G any, H any, I any, B2 any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B2, C, D, E, F, G, H, I], B, B2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B2, C, D, E, F, G, H, I], B, B2](func(ctx context.Context, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (int, B, error) {
		return 1, source.B, nil
	}, func(ctx context.Context, focus B2, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (lo.Tuple9[A, B2, C, D, E, F, G, H, I], error) {
		return lo.T9(source.A, focus, source.C, source.D, source.E, source.F, source.G, source.H, source.I), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         1,
		}
	}))
}

// T9C returns a [Lens] focusing on element 2 of a [lo.Tuple9]
//
// See: [T9CP] for a polymorphic version
func T9C[A any, B any, C any, D any, E any, F any, G any, H any, I any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G, H, I], C, C, ReturnOne, ReadWrite, UniDir, Pure] {
	return T9CP[A, B, C, D, E, F, G, H, I, C]()
}

// T9CP returns a polymorphic [Lens] focusing on element 2 of a [lo.Tuple9]
//
// See: [T9C] for a non polymorphic version
func T9CP[A any, B any, C any, D any, E any, F any, G any, H any, I any, C2 any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C2, D, E, F, G, H, I], C, C2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C2, D, E, F, G, H, I], C, C2](func(ctx context.Context, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (int, C, error) {
		return 2, source.C, nil
	}, func(ctx context.Context, focus C2, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (lo.Tuple9[A, B, C2, D, E, F, G, H, I], error) {
		return lo.T9(source.A, source.B, focus, source.D, source.E, source.F, source.G, source.H, source.I), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         2,
		}
	}))
}

// T9D returns a [Lens] focusing on element 3 of a [lo.Tuple9]
//
// See: [T9DP] for a polymorphic version
func T9D[A any, B any, C any, D any, E any, F any, G any, H any, I any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G, H, I], D, D, ReturnOne, ReadWrite, UniDir, Pure] {
	return T9DP[A, B, C, D, E, F, G, H, I, D]()
}

// T9DP returns a polymorphic [Lens] focusing on element 3 of a [lo.Tuple9]
//
// See: [T9D] for a non polymorphic version
func T9DP[A any, B any, C any, D any, E any, F any, G any, H any, I any, D2 any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D2, E, F, G, H, I], D, D2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D2, E, F, G, H, I], D, D2](func(ctx context.Context, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (int, D, error) {
		return 3, source.D, nil
	}, func(ctx context.Context, focus D2, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (lo.Tuple9[A, B, C, D2, E, F, G, H, I], error) {
		return lo.T9(source.A, source.B, source.C, focus, source.E, source.F, source.G, source.H, source.I), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         3,
		}
	}))
}

// T9E returns a [Lens] focusing on element 4 of a [lo.Tuple9]
//
// See: [T9EP] for a polymorphic version
func T9E[A any, B any, C any, D any, E any, F any, G any, H any, I any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G, H, I], E, E, ReturnOne, ReadWrite, UniDir, Pure] {
	return T9EP[A, B, C, D, E, F, G, H, I, E]()
}

// T9EP returns a polymorphic [Lens] focusing on element 4 of a [lo.Tuple9]
//
// See: [T9E] for a non polymorphic version
func T9EP[A any, B any, C any, D any, E any, F any, G any, H any, I any, E2 any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E2, F, G, H, I], E, E2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E2, F, G, H, I], E, E2](func(ctx context.Context, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (int, E, error) {
		return 4, source.E, nil
	}, func(ctx context.Context, focus E2, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (lo.Tuple9[A, B, C, D, E2, F, G, H, I], error) {
		return lo.T9(source.A, source.B, source.C, source.D, focus, source.F, source.G, source.H, source.I), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         4,
		}
	}))
}

// T9F returns a [Lens] focusing on element 5 of a [lo.Tuple9]
//
// See: [T9FP] for a polymorphic version
func T9F[A any, B any, C any, D any, E any, F any, G any, H any, I any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G, H, I], F, F, ReturnOne, ReadWrite, UniDir, Pure] {
	return T9FP[A, B, C, D, E, F, G, H, I, F]()
}

// T9FP returns a polymorphic [Lens] focusing on element 5 of a [lo.Tuple9]
//
// See: [T9F] for a non polymorphic version
func T9FP[A any, B any, C any, D any, E any, F any, G any, H any, I any, F2 any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F2, G, H, I], F, F2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F2, G, H, I], F, F2](func(ctx context.Context, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (int, F, error) {
		return 5, source.F, nil
	}, func(ctx context.Context, focus F2, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (lo.Tuple9[A, B, C, D, E, F2, G, H, I], error) {
		return lo.T9(source.A, source.B, source.C, source.D, source.E, focus, source.G, source.H, source.I), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         5,
		}
	}))
}

// T9G returns a [Lens] focusing on element 6 of a [lo.Tuple9]
//
// See: [T9GP] for a polymorphic version
func T9G[A any, B any, C any, D any, E any, F any, G any, H any, I any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G, H, I], G, G, ReturnOne, ReadWrite, UniDir, Pure] {
	return T9GP[A, B, C, D, E, F, G, H, I, G]()
}

// T9GP returns a polymorphic [Lens] focusing on element 6 of a [lo.Tuple9]
//
// See: [T9G] for a non polymorphic version
func T9GP[A any, B any, C any, D any, E any, F any, G any, H any, I any, G2 any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G2, H, I], G, G2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G2, H, I], G, G2](func(ctx context.Context, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (int, G, error) {
		return 6, source.G, nil
	}, func(ctx context.Context, focus G2, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (lo.Tuple9[A, B, C, D, E, F, G2, H, I], error) {
		return lo.T9(source.A, source.B, source.C, source.D, source.E, source.F, focus, source.H, source.I), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         6,
		}
	}))
}

// T9H returns a [Lens] focusing on element 7 of a [lo.Tuple9]
//
// See: [T9HP] for a polymorphic version
func T9H[A any, B any, C any, D any, E any, F any, G any, H any, I any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G, H, I], H, H, ReturnOne, ReadWrite, UniDir, Pure] {
	return T9HP[A, B, C, D, E, F, G, H, I, H]()
}

// T9HP returns a polymorphic [Lens] focusing on element 7 of a [lo.Tuple9]
//
// See: [T9H] for a non polymorphic version
func T9HP[A any, B any, C any, D any, E any, F any, G any, H any, I any, H2 any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G, H2, I], H, H2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G, H2, I], H, H2](func(ctx context.Context, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (int, H, error) {
		return 7, source.H, nil
	}, func(ctx context.Context, focus H2, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (lo.Tuple9[A, B, C, D, E, F, G, H2, I], error) {
		return lo.T9(source.A, source.B, source.C, source.D, source.E, source.F, source.G, focus, source.I), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         7,
		}
	}))
}

// T9I returns a [Lens] focusing on element 8 of a [lo.Tuple9]
//
// See: [T9IP] for a polymorphic version
func T9I[A any, B any, C any, D any, E any, F any, G any, H any, I any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G, H, I], I, I, ReturnOne, ReadWrite, UniDir, Pure] {
	return T9IP[A, B, C, D, E, F, G, H, I, I]()
}

// T9IP returns a polymorphic [Lens] focusing on element 8 of a [lo.Tuple9]
//
// See: [T9I] for a non polymorphic version
func T9IP[A any, B any, C any, D any, E any, F any, G any, H any, I any, I2 any]() Optic[int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G, H, I2], I, I2, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, int, lo.Tuple9[A, B, C, D, E, F, G, H, I], lo.Tuple9[A, B, C, D, E, F, G, H, I2], I, I2](func(ctx context.Context, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (int, I, error) {
		return 8, source.I, nil
	}, func(ctx context.Context, focus I2, source lo.Tuple9[A, B, C, D, E, F, G, H, I]) (lo.Tuple9[A, B, C, D, E, F, G, H, I2], error) {
		return lo.T9(source.A, source.B, source.C, source.D, source.E, source.F, source.G, source.H, focus), nil
	}, IxMatchComparable[int](), ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleElement{
			OpticTypeExpr: ot,
			Index:         8,
		}
	}))
}

// T2ToCol returns an [Iso] that converts a [lo.Tuple2] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T2PCol for a polymorphic version
func T2ToCol[A any]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], Collection[int, A, Pure], Collection[int, A, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return T2ToColP[A, A]()
}

// T2CToolP returns a polymorphic [Iso] that converts a [lo.Tuple2] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T2ToCol for a non polymorphic version
func T2ToColP[A any, B any]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[B, B], Collection[int, A, Pure], Collection[int, B, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, lo.Tuple2[A, A], lo.Tuple2[B, B], Collection[int, A, Pure], Collection[int, B, Pure]](func(ctx context.Context, source lo.Tuple2[A, A]) (Collection[int, A, Pure], error) {
		return tnColGetter(ctx, &source.A, &source.B), nil
	}, func(ctx context.Context, focus Collection[int, B, Pure]) (lo.Tuple2[B, B], error) {
		var ret lo.Tuple2[B, B]

		err := tnColReverseGet(ctx, focus, &ret.A, &ret.B)
		return ret, err
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.ToCol{
			OpticTypeExpr: ot,
			I:             reflect.TypeFor[int](),
			A:             reflect.TypeFor[A](),
			B:             reflect.TypeFor[B](),
		}
	}))
}

// T3ToCol returns an [Iso] that converts a [lo.Tuple3] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T3PCol for a polymorphic version
func T3ToCol[A any]() Optic[Void, lo.Tuple3[A, A, A], lo.Tuple3[A, A, A], Collection[int, A, Pure], Collection[int, A, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return T3ToColP[A, A]()
}

// T3CToolP returns a polymorphic [Iso] that converts a [lo.Tuple3] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T3ToCol for a non polymorphic version
func T3ToColP[A any, B any]() Optic[Void, lo.Tuple3[A, A, A], lo.Tuple3[B, B, B], Collection[int, A, Pure], Collection[int, B, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, lo.Tuple3[A, A, A], lo.Tuple3[B, B, B], Collection[int, A, Pure], Collection[int, B, Pure]](func(ctx context.Context, source lo.Tuple3[A, A, A]) (Collection[int, A, Pure], error) {
		return tnColGetter(ctx, &source.A, &source.B, &source.C), nil
	}, func(ctx context.Context, focus Collection[int, B, Pure]) (lo.Tuple3[B, B, B], error) {
		var ret lo.Tuple3[B, B, B]

		err := tnColReverseGet(ctx, focus, &ret.A, &ret.B, &ret.C)
		return ret, err
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.ToCol{
			OpticTypeExpr: ot,
			I:             reflect.TypeFor[int](),
			A:             reflect.TypeFor[A](),
			B:             reflect.TypeFor[B](),
		}
	}))
}

// T4ToCol returns an [Iso] that converts a [lo.Tuple4] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T4PCol for a polymorphic version
func T4ToCol[A any]() Optic[Void, lo.Tuple4[A, A, A, A], lo.Tuple4[A, A, A, A], Collection[int, A, Pure], Collection[int, A, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return T4ToColP[A, A]()
}

// T4CToolP returns a polymorphic [Iso] that converts a [lo.Tuple4] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T4ToCol for a non polymorphic version
func T4ToColP[A any, B any]() Optic[Void, lo.Tuple4[A, A, A, A], lo.Tuple4[B, B, B, B], Collection[int, A, Pure], Collection[int, B, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, lo.Tuple4[A, A, A, A], lo.Tuple4[B, B, B, B], Collection[int, A, Pure], Collection[int, B, Pure]](func(ctx context.Context, source lo.Tuple4[A, A, A, A]) (Collection[int, A, Pure], error) {
		return tnColGetter(ctx, &source.A, &source.B, &source.C, &source.D), nil
	}, func(ctx context.Context, focus Collection[int, B, Pure]) (lo.Tuple4[B, B, B, B], error) {
		var ret lo.Tuple4[B, B, B, B]

		err := tnColReverseGet(ctx, focus, &ret.A, &ret.B, &ret.C, &ret.D)
		return ret, err
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.ToCol{
			OpticTypeExpr: ot,
			I:             reflect.TypeFor[int](),
			A:             reflect.TypeFor[A](),
			B:             reflect.TypeFor[B](),
		}
	}))
}

// T5ToCol returns an [Iso] that converts a [lo.Tuple5] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T5PCol for a polymorphic version
func T5ToCol[A any]() Optic[Void, lo.Tuple5[A, A, A, A, A], lo.Tuple5[A, A, A, A, A], Collection[int, A, Pure], Collection[int, A, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return T5ToColP[A, A]()
}

// T5CToolP returns a polymorphic [Iso] that converts a [lo.Tuple5] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T5ToCol for a non polymorphic version
func T5ToColP[A any, B any]() Optic[Void, lo.Tuple5[A, A, A, A, A], lo.Tuple5[B, B, B, B, B], Collection[int, A, Pure], Collection[int, B, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, lo.Tuple5[A, A, A, A, A], lo.Tuple5[B, B, B, B, B], Collection[int, A, Pure], Collection[int, B, Pure]](func(ctx context.Context, source lo.Tuple5[A, A, A, A, A]) (Collection[int, A, Pure], error) {
		return tnColGetter(ctx, &source.A, &source.B, &source.C, &source.D, &source.E), nil
	}, func(ctx context.Context, focus Collection[int, B, Pure]) (lo.Tuple5[B, B, B, B, B], error) {
		var ret lo.Tuple5[B, B, B, B, B]

		err := tnColReverseGet(ctx, focus, &ret.A, &ret.B, &ret.C, &ret.D, &ret.E)
		return ret, err
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.ToCol{
			OpticTypeExpr: ot,
			I:             reflect.TypeFor[int](),
			A:             reflect.TypeFor[A](),
			B:             reflect.TypeFor[B](),
		}
	}))
}

// T6ToCol returns an [Iso] that converts a [lo.Tuple6] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T6PCol for a polymorphic version
func T6ToCol[A any]() Optic[Void, lo.Tuple6[A, A, A, A, A, A], lo.Tuple6[A, A, A, A, A, A], Collection[int, A, Pure], Collection[int, A, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return T6ToColP[A, A]()
}

// T6CToolP returns a polymorphic [Iso] that converts a [lo.Tuple6] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T6ToCol for a non polymorphic version
func T6ToColP[A any, B any]() Optic[Void, lo.Tuple6[A, A, A, A, A, A], lo.Tuple6[B, B, B, B, B, B], Collection[int, A, Pure], Collection[int, B, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, lo.Tuple6[A, A, A, A, A, A], lo.Tuple6[B, B, B, B, B, B], Collection[int, A, Pure], Collection[int, B, Pure]](func(ctx context.Context, source lo.Tuple6[A, A, A, A, A, A]) (Collection[int, A, Pure], error) {
		return tnColGetter(ctx, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F), nil
	}, func(ctx context.Context, focus Collection[int, B, Pure]) (lo.Tuple6[B, B, B, B, B, B], error) {
		var ret lo.Tuple6[B, B, B, B, B, B]

		err := tnColReverseGet(ctx, focus, &ret.A, &ret.B, &ret.C, &ret.D, &ret.E, &ret.F)
		return ret, err
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.ToCol{
			OpticTypeExpr: ot,
			I:             reflect.TypeFor[int](),
			A:             reflect.TypeFor[A](),
			B:             reflect.TypeFor[B](),
		}
	}))
}

// T7ToCol returns an [Iso] that converts a [lo.Tuple7] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T7PCol for a polymorphic version
func T7ToCol[A any]() Optic[Void, lo.Tuple7[A, A, A, A, A, A, A], lo.Tuple7[A, A, A, A, A, A, A], Collection[int, A, Pure], Collection[int, A, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return T7ToColP[A, A]()
}

// T7CToolP returns a polymorphic [Iso] that converts a [lo.Tuple7] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T7ToCol for a non polymorphic version
func T7ToColP[A any, B any]() Optic[Void, lo.Tuple7[A, A, A, A, A, A, A], lo.Tuple7[B, B, B, B, B, B, B], Collection[int, A, Pure], Collection[int, B, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, lo.Tuple7[A, A, A, A, A, A, A], lo.Tuple7[B, B, B, B, B, B, B], Collection[int, A, Pure], Collection[int, B, Pure]](func(ctx context.Context, source lo.Tuple7[A, A, A, A, A, A, A]) (Collection[int, A, Pure], error) {
		return tnColGetter(ctx, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F, &source.G), nil
	}, func(ctx context.Context, focus Collection[int, B, Pure]) (lo.Tuple7[B, B, B, B, B, B, B], error) {
		var ret lo.Tuple7[B, B, B, B, B, B, B]

		err := tnColReverseGet(ctx, focus, &ret.A, &ret.B, &ret.C, &ret.D, &ret.E, &ret.F, &ret.G)
		return ret, err
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.ToCol{
			OpticTypeExpr: ot,
			I:             reflect.TypeFor[int](),
			A:             reflect.TypeFor[A](),
			B:             reflect.TypeFor[B](),
		}
	}))
}

// T8ToCol returns an [Iso] that converts a [lo.Tuple8] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T8PCol for a polymorphic version
func T8ToCol[A any]() Optic[Void, lo.Tuple8[A, A, A, A, A, A, A, A], lo.Tuple8[A, A, A, A, A, A, A, A], Collection[int, A, Pure], Collection[int, A, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return T8ToColP[A, A]()
}

// T8CToolP returns a polymorphic [Iso] that converts a [lo.Tuple8] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T8ToCol for a non polymorphic version
func T8ToColP[A any, B any]() Optic[Void, lo.Tuple8[A, A, A, A, A, A, A, A], lo.Tuple8[B, B, B, B, B, B, B, B], Collection[int, A, Pure], Collection[int, B, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, lo.Tuple8[A, A, A, A, A, A, A, A], lo.Tuple8[B, B, B, B, B, B, B, B], Collection[int, A, Pure], Collection[int, B, Pure]](func(ctx context.Context, source lo.Tuple8[A, A, A, A, A, A, A, A]) (Collection[int, A, Pure], error) {
		return tnColGetter(ctx, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F, &source.G, &source.H), nil
	}, func(ctx context.Context, focus Collection[int, B, Pure]) (lo.Tuple8[B, B, B, B, B, B, B, B], error) {
		var ret lo.Tuple8[B, B, B, B, B, B, B, B]

		err := tnColReverseGet(ctx, focus, &ret.A, &ret.B, &ret.C, &ret.D, &ret.E, &ret.F, &ret.G, &ret.H)
		return ret, err
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.ToCol{
			OpticTypeExpr: ot,
			I:             reflect.TypeFor[int](),
			A:             reflect.TypeFor[A](),
			B:             reflect.TypeFor[B](),
		}
	}))
}

// T9ToCol returns an [Iso] that converts a [lo.Tuple9] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T9PCol for a polymorphic version
func T9ToCol[A any]() Optic[Void, lo.Tuple9[A, A, A, A, A, A, A, A, A], lo.Tuple9[A, A, A, A, A, A, A, A, A], Collection[int, A, Pure], Collection[int, A, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return T9ToColP[A, A]()
}

// T9CToolP returns a polymorphic [Iso] that converts a [lo.Tuple9] to a [Collection]
//
// Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.
//
// See: T9ToCol for a non polymorphic version
func T9ToColP[A any, B any]() Optic[Void, lo.Tuple9[A, A, A, A, A, A, A, A, A], lo.Tuple9[B, B, B, B, B, B, B, B, B], Collection[int, A, Pure], Collection[int, B, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, lo.Tuple9[A, A, A, A, A, A, A, A, A], lo.Tuple9[B, B, B, B, B, B, B, B, B], Collection[int, A, Pure], Collection[int, B, Pure]](func(ctx context.Context, source lo.Tuple9[A, A, A, A, A, A, A, A, A]) (Collection[int, A, Pure], error) {
		return tnColGetter(ctx, &source.A, &source.B, &source.C, &source.D, &source.E, &source.F, &source.G, &source.H, &source.I), nil
	}, func(ctx context.Context, focus Collection[int, B, Pure]) (lo.Tuple9[B, B, B, B, B, B, B, B, B], error) {
		var ret lo.Tuple9[B, B, B, B, B, B, B, B, B]

		err := tnColReverseGet(ctx, focus, &ret.A, &ret.B, &ret.C, &ret.D, &ret.E, &ret.F, &ret.G, &ret.H, &ret.I)
		return ret, err
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.ToCol{
			OpticTypeExpr: ot,
			I:             reflect.TypeFor[int](),
			A:             reflect.TypeFor[A](),
			B:             reflect.TypeFor[B](),
		}
	}))
}

// T2ColType returns a [CollectionType] wrapper for [lo.Tuple2].
//
// See: T2ColTypeP for a polymorphic version
func T2ColType[A any]() CollectionType[int, lo.Tuple2[A, A], lo.Tuple2[A, A], A, A, Pure] {
	return T2ColTypeP[A, A]()
}

// T2ColTypeP returns a polymorphic [CollectionType] wrapper for [lo.Tuple2].
//
// See: T2ColType for a non polymorphic version
func T2ColTypeP[A any, B any]() CollectionType[int, lo.Tuple2[A, A], lo.Tuple2[B, B], A, B, Pure] {
	return ColTypeP(T2ToColP[A, B](), AsReverseGet(T2ToColP[B, A]()), TraverseT2P[A, B]())
}

// T3ColType returns a [CollectionType] wrapper for [lo.Tuple3].
//
// See: T3ColTypeP for a polymorphic version
func T3ColType[A any]() CollectionType[int, lo.Tuple3[A, A, A], lo.Tuple3[A, A, A], A, A, Pure] {
	return T3ColTypeP[A, A]()
}

// T3ColTypeP returns a polymorphic [CollectionType] wrapper for [lo.Tuple3].
//
// See: T3ColType for a non polymorphic version
func T3ColTypeP[A any, B any]() CollectionType[int, lo.Tuple3[A, A, A], lo.Tuple3[B, B, B], A, B, Pure] {
	return ColTypeP(T3ToColP[A, B](), AsReverseGet(T3ToColP[B, A]()), TraverseT3P[A, B]())
}

// T4ColType returns a [CollectionType] wrapper for [lo.Tuple4].
//
// See: T4ColTypeP for a polymorphic version
func T4ColType[A any]() CollectionType[int, lo.Tuple4[A, A, A, A], lo.Tuple4[A, A, A, A], A, A, Pure] {
	return T4ColTypeP[A, A]()
}

// T4ColTypeP returns a polymorphic [CollectionType] wrapper for [lo.Tuple4].
//
// See: T4ColType for a non polymorphic version
func T4ColTypeP[A any, B any]() CollectionType[int, lo.Tuple4[A, A, A, A], lo.Tuple4[B, B, B, B], A, B, Pure] {
	return ColTypeP(T4ToColP[A, B](), AsReverseGet(T4ToColP[B, A]()), TraverseT4P[A, B]())
}

// T5ColType returns a [CollectionType] wrapper for [lo.Tuple5].
//
// See: T5ColTypeP for a polymorphic version
func T5ColType[A any]() CollectionType[int, lo.Tuple5[A, A, A, A, A], lo.Tuple5[A, A, A, A, A], A, A, Pure] {
	return T5ColTypeP[A, A]()
}

// T5ColTypeP returns a polymorphic [CollectionType] wrapper for [lo.Tuple5].
//
// See: T5ColType for a non polymorphic version
func T5ColTypeP[A any, B any]() CollectionType[int, lo.Tuple5[A, A, A, A, A], lo.Tuple5[B, B, B, B, B], A, B, Pure] {
	return ColTypeP(T5ToColP[A, B](), AsReverseGet(T5ToColP[B, A]()), TraverseT5P[A, B]())
}

// T6ColType returns a [CollectionType] wrapper for [lo.Tuple6].
//
// See: T6ColTypeP for a polymorphic version
func T6ColType[A any]() CollectionType[int, lo.Tuple6[A, A, A, A, A, A], lo.Tuple6[A, A, A, A, A, A], A, A, Pure] {
	return T6ColTypeP[A, A]()
}

// T6ColTypeP returns a polymorphic [CollectionType] wrapper for [lo.Tuple6].
//
// See: T6ColType for a non polymorphic version
func T6ColTypeP[A any, B any]() CollectionType[int, lo.Tuple6[A, A, A, A, A, A], lo.Tuple6[B, B, B, B, B, B], A, B, Pure] {
	return ColTypeP(T6ToColP[A, B](), AsReverseGet(T6ToColP[B, A]()), TraverseT6P[A, B]())
}

// T7ColType returns a [CollectionType] wrapper for [lo.Tuple7].
//
// See: T7ColTypeP for a polymorphic version
func T7ColType[A any]() CollectionType[int, lo.Tuple7[A, A, A, A, A, A, A], lo.Tuple7[A, A, A, A, A, A, A], A, A, Pure] {
	return T7ColTypeP[A, A]()
}

// T7ColTypeP returns a polymorphic [CollectionType] wrapper for [lo.Tuple7].
//
// See: T7ColType for a non polymorphic version
func T7ColTypeP[A any, B any]() CollectionType[int, lo.Tuple7[A, A, A, A, A, A, A], lo.Tuple7[B, B, B, B, B, B, B], A, B, Pure] {
	return ColTypeP(T7ToColP[A, B](), AsReverseGet(T7ToColP[B, A]()), TraverseT7P[A, B]())
}

// T8ColType returns a [CollectionType] wrapper for [lo.Tuple8].
//
// See: T8ColTypeP for a polymorphic version
func T8ColType[A any]() CollectionType[int, lo.Tuple8[A, A, A, A, A, A, A, A], lo.Tuple8[A, A, A, A, A, A, A, A], A, A, Pure] {
	return T8ColTypeP[A, A]()
}

// T8ColTypeP returns a polymorphic [CollectionType] wrapper for [lo.Tuple8].
//
// See: T8ColType for a non polymorphic version
func T8ColTypeP[A any, B any]() CollectionType[int, lo.Tuple8[A, A, A, A, A, A, A, A], lo.Tuple8[B, B, B, B, B, B, B, B], A, B, Pure] {
	return ColTypeP(T8ToColP[A, B](), AsReverseGet(T8ToColP[B, A]()), TraverseT8P[A, B]())
}

// T9ColType returns a [CollectionType] wrapper for [lo.Tuple9].
//
// See: T9ColTypeP for a polymorphic version
func T9ColType[A any]() CollectionType[int, lo.Tuple9[A, A, A, A, A, A, A, A, A], lo.Tuple9[A, A, A, A, A, A, A, A, A], A, A, Pure] {
	return T9ColTypeP[A, A]()
}

// T9ColTypeP returns a polymorphic [CollectionType] wrapper for [lo.Tuple9].
//
// See: T9ColType for a non polymorphic version
func T9ColTypeP[A any, B any]() CollectionType[int, lo.Tuple9[A, A, A, A, A, A, A, A, A], lo.Tuple9[B, B, B, B, B, B, B, B, B], A, B, Pure] {
	return ColTypeP(T9ToColP[A, B](), AsReverseGet(T9ToColP[B, A]()), TraverseT9P[A, B]())
}

// The T2Of Combinator constructs a [lo.Tuple2] whose elements are the focuses of the given [Optic]s
//
// Note: The number of focused tuples is limited by the optic that focuses the least elements.
func T2Of[I0 any, I1 any, S, T any, A0 any, A1 any, RET0 any, RW0 any, DIR0 any, ERR0 any, RET1 any, RW1 any, DIR1 any, ERR1 any](o0 Optic[I0, S, S, A0, A0, RET0, RW0, DIR0, ERR0], o1 Optic[I1, S, T, A1, A1, RET1, RW1, DIR1, ERR1]) Optic[lo.Tuple2[I0, I1], S, T, lo.Tuple2[A0, A1], lo.Tuple2[A0, A1], CompositionTree[RET0, RET1], CompositionTree[RW0, RW1], UniDir, CompositionTree[ERR0, ERR1]] {
	// TODO: I have edited this function I need to update maketuple to generate the slightly polymorphic version here.
	return CombiTraversal[CompositionTree[RET0, RET1], CompositionTree[RW0, RW1], CompositionTree[ERR0, ERR1], lo.Tuple2[I0, I1], S, T, lo.Tuple2[A0, A1], lo.Tuple2[A0, A1]](func(ctx context.Context, source S) SeqIE[lo.Tuple2[I0, I1], lo.Tuple2[A0, A1]] {
		return func(yield func(focusHello ValueIE[lo.Tuple2[I0, I1], lo.Tuple2[A0, A1]]) bool) {
			next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
			defer stop1()
			o0.AsIter()(ctx, source)(func(val ValueIE[I0, A0]) bool {
				i0, a0, err := val.Get()
				if err != nil {
					var i lo.Tuple2[I0, I1]

					var a lo.Tuple2[A0, A1]

					return yield(ValIE(i, a, err))
				}

				v1, ok := next1()
				if !ok {
					return false
				}

				i1, a1, err := v1.Get()
				if err != nil {
					var i lo.Tuple2[I0, I1]

					var a lo.Tuple2[A0, A1]

					return yield(ValIE(i, a, err))
				}

				return yield(ValIE(lo.T2(i0, i1), lo.T2(a0, a1), nil))
			})
		}
	}, func(ctx context.Context, source S) (int, error) {
		l0, err := o0.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l1, err := o1.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		return min(l0, l1), nil
	}, func(ctx context.Context, fmap func(index lo.Tuple2[I0, I1], focus lo.Tuple2[A0, A1]) (lo.Tuple2[A0, A1], error), source S) (T, error) {
		next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
		defer stop1()
		var mapping1 []A1

		source, err := o0.AsModify()(ctx, func(index I0, focus A0) (A0, error) {
			v1, ok := next1()
			if !ok {
				return focus, nil
			}

			i1, a1, err := v1.Get()
			if err != nil {
				var a A0

				return a, err
			}

			mapped, err := fmap(lo.T2(index, i1), lo.T2(focus, a1))
			if err != nil {
				var a A0

				return a, err
			}

			mapping1 = append(mapping1, mapped.B)
			return mapped.A, err
		}, source)
		if err != nil {
			var t T

			return t, err
		}

		i := 0
		retT, err := o1.AsModify()(ctx, func(index I1, focus A1) (A1, error) {
			if i >= len(mapping1) {
				return focus, nil
			}

			return mapping1[i], nil
		}, source)
		if err != nil {
			var t T

			return t, err
		}

		return retT, nil
	}, nil, func(a lo.Tuple2[I0, I1], b lo.Tuple2[I0, I1]) bool {
		if o0.AsIxMatch()(a.A, b.A) != true {
			return false
		}

		if o1.AsIxMatch()(a.B, b.B) != true {
			return false
		}

		return true
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleOf{
			OpticTypeExpr: ot,
			Elements:      []expr.OpticExpression{o0.AsExpr(), o1.AsExpr()},
		}
	}, o0, o1))
}

// The T3Of Combinator constructs a [lo.Tuple3] whose elements are the focuses of the given [Optic]s
//
// Note: The number of focused tuples is limited by the optic that focuses the least elements.
func T3Of[I0 any, I1 any, I2 any, S any, A0 any, A1 any, A2 any, RET0 any, RW0 any, DIR0 any, ERR0 any, RET1 any, RW1 any, DIR1 any, ERR1 any, RET2 any, RW2 any, DIR2 any, ERR2 any](o0 Optic[I0, S, S, A0, A0, RET0, RW0, DIR0, ERR0], o1 Optic[I1, S, S, A1, A1, RET1, RW1, DIR1, ERR1], o2 Optic[I2, S, S, A2, A2, RET2, RW2, DIR2, ERR2]) Optic[lo.Tuple3[I0, I1, I2], S, S, lo.Tuple3[A0, A1, A2], lo.Tuple3[A0, A1, A2], CompositionTree[CompositionTree[RET0, RET1], RET2], CompositionTree[CompositionTree[RW0, RW1], RW2], UniDir, CompositionTree[CompositionTree[ERR0, ERR1], ERR2]] {
	return CombiTraversal[CompositionTree[CompositionTree[RET0, RET1], RET2], CompositionTree[CompositionTree[RW0, RW1], RW2], CompositionTree[CompositionTree[ERR0, ERR1], ERR2], lo.Tuple3[I0, I1, I2], S, S, lo.Tuple3[A0, A1, A2], lo.Tuple3[A0, A1, A2]](func(ctx context.Context, source S) SeqIE[lo.Tuple3[I0, I1, I2], lo.Tuple3[A0, A1, A2]] {
		return func(yield func(focusHello ValueIE[lo.Tuple3[I0, I1, I2], lo.Tuple3[A0, A1, A2]]) bool) {
			next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
			defer stop1()
			next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
			defer stop2()
			o0.AsIter()(ctx, source)(func(val ValueIE[I0, A0]) bool {
				i0, a0, err := val.Get()
				if err != nil {
					var i lo.Tuple3[I0, I1, I2]

					var a lo.Tuple3[A0, A1, A2]

					return yield(ValIE(i, a, err))
				}

				v1, ok := next1()
				if !ok {
					return false
				}

				i1, a1, err := v1.Get()
				if err != nil {
					var i lo.Tuple3[I0, I1, I2]

					var a lo.Tuple3[A0, A1, A2]

					return yield(ValIE(i, a, err))
				}

				v2, ok := next2()
				if !ok {
					return false
				}

				i2, a2, err := v2.Get()
				if err != nil {
					var i lo.Tuple3[I0, I1, I2]

					var a lo.Tuple3[A0, A1, A2]

					return yield(ValIE(i, a, err))
				}

				return yield(ValIE(lo.T3(i0, i1, i2), lo.T3(a0, a1, a2), nil))
			})
		}
	}, func(ctx context.Context, source S) (int, error) {
		l0, err := o0.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l1, err := o1.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l2, err := o2.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		return min(l0, l1, l2), nil
	}, func(ctx context.Context, fmap func(index lo.Tuple3[I0, I1, I2], focus lo.Tuple3[A0, A1, A2]) (lo.Tuple3[A0, A1, A2], error), source S) (S, error) {
		next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
		defer stop1()
		var mapping1 []A1

		next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
		defer stop2()
		var mapping2 []A2

		ret, err := o0.AsModify()(ctx, func(index I0, focus A0) (A0, error) {
			v1, ok := next1()
			if !ok {
				return focus, nil
			}

			i1, a1, err := v1.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v2, ok := next2()
			if !ok {
				return focus, nil
			}

			i2, a2, err := v2.Get()
			if err != nil {
				var a A0

				return a, err
			}

			mapped, err := fmap(lo.T3(index, i1, i2), lo.T3(focus, a1, a2))
			if err != nil {
				var a A0

				return a, err
			}

			mapping1 = append(mapping1, mapped.B)
			mapping2 = append(mapping2, mapped.C)
			return mapped.A, err
		}, source)
		if err != nil {
			var s S

			return s, err
		}

		i := 0
		ret, err = o1.AsModify()(ctx, func(index I1, focus A1) (A1, error) {
			if i >= len(mapping1) {
				return focus, nil
			}

			return mapping1[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o2.AsModify()(ctx, func(index I2, focus A2) (A2, error) {
			if i >= len(mapping2) {
				return focus, nil
			}

			return mapping2[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		return ret, nil
	}, nil, func(a lo.Tuple3[I0, I1, I2], b lo.Tuple3[I0, I1, I2]) bool {
		if o0.AsIxMatch()(a.A, b.A) != true {
			return false
		}

		if o1.AsIxMatch()(a.B, b.B) != true {
			return false
		}

		if o2.AsIxMatch()(a.C, b.C) != true {
			return false
		}

		return true
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleOf{
			OpticTypeExpr: ot,
			Elements:      []expr.OpticExpression{o0.AsExpr(), o1.AsExpr(), o2.AsExpr()},
		}
	}, o0, o1, o2))
}

// The T4Of Combinator constructs a [lo.Tuple4] whose elements are the focuses of the given [Optic]s
//
// Note: The number of focused tuples is limited by the optic that focuses the least elements.
func T4Of[I0 any, I1 any, I2 any, I3 any, S any, A0 any, A1 any, A2 any, A3 any, RET0 any, RW0 any, DIR0 any, ERR0 any, RET1 any, RW1 any, DIR1 any, ERR1 any, RET2 any, RW2 any, DIR2 any, ERR2 any, RET3 any, RW3 any, DIR3 any, ERR3 any](o0 Optic[I0, S, S, A0, A0, RET0, RW0, DIR0, ERR0], o1 Optic[I1, S, S, A1, A1, RET1, RW1, DIR1, ERR1], o2 Optic[I2, S, S, A2, A2, RET2, RW2, DIR2, ERR2], o3 Optic[I3, S, S, A3, A3, RET3, RW3, DIR3, ERR3]) Optic[lo.Tuple4[I0, I1, I2, I3], S, S, lo.Tuple4[A0, A1, A2, A3], lo.Tuple4[A0, A1, A2, A3], CompositionTree[CompositionTree[RET0, RET1], CompositionTree[RET2, RET3]], CompositionTree[CompositionTree[RW0, RW1], CompositionTree[RW2, RW3]], UniDir, CompositionTree[CompositionTree[ERR0, ERR1], CompositionTree[ERR2, ERR3]]] {
	return CombiTraversal[CompositionTree[CompositionTree[RET0, RET1], CompositionTree[RET2, RET3]], CompositionTree[CompositionTree[RW0, RW1], CompositionTree[RW2, RW3]], CompositionTree[CompositionTree[ERR0, ERR1], CompositionTree[ERR2, ERR3]], lo.Tuple4[I0, I1, I2, I3], S, S, lo.Tuple4[A0, A1, A2, A3], lo.Tuple4[A0, A1, A2, A3]](func(ctx context.Context, source S) SeqIE[lo.Tuple4[I0, I1, I2, I3], lo.Tuple4[A0, A1, A2, A3]] {
		return func(yield func(focusHello ValueIE[lo.Tuple4[I0, I1, I2, I3], lo.Tuple4[A0, A1, A2, A3]]) bool) {
			next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
			defer stop1()
			next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
			defer stop2()
			next3, stop3 := iter.Pull(iter.Seq[ValueIE[I3, A3]](o3.AsIter()(ctx, source)))
			defer stop3()
			o0.AsIter()(ctx, source)(func(val ValueIE[I0, A0]) bool {
				i0, a0, err := val.Get()
				if err != nil {
					var i lo.Tuple4[I0, I1, I2, I3]

					var a lo.Tuple4[A0, A1, A2, A3]

					return yield(ValIE(i, a, err))
				}

				v1, ok := next1()
				if !ok {
					return false
				}

				i1, a1, err := v1.Get()
				if err != nil {
					var i lo.Tuple4[I0, I1, I2, I3]

					var a lo.Tuple4[A0, A1, A2, A3]

					return yield(ValIE(i, a, err))
				}

				v2, ok := next2()
				if !ok {
					return false
				}

				i2, a2, err := v2.Get()
				if err != nil {
					var i lo.Tuple4[I0, I1, I2, I3]

					var a lo.Tuple4[A0, A1, A2, A3]

					return yield(ValIE(i, a, err))
				}

				v3, ok := next3()
				if !ok {
					return false
				}

				i3, a3, err := v3.Get()
				if err != nil {
					var i lo.Tuple4[I0, I1, I2, I3]

					var a lo.Tuple4[A0, A1, A2, A3]

					return yield(ValIE(i, a, err))
				}

				return yield(ValIE(lo.T4(i0, i1, i2, i3), lo.T4(a0, a1, a2, a3), nil))
			})
		}
	}, func(ctx context.Context, source S) (int, error) {
		l0, err := o0.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l1, err := o1.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l2, err := o2.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l3, err := o3.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		return min(l0, l1, l2, l3), nil
	}, func(ctx context.Context, fmap func(index lo.Tuple4[I0, I1, I2, I3], focus lo.Tuple4[A0, A1, A2, A3]) (lo.Tuple4[A0, A1, A2, A3], error), source S) (S, error) {
		next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
		defer stop1()
		var mapping1 []A1

		next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
		defer stop2()
		var mapping2 []A2

		next3, stop3 := iter.Pull(iter.Seq[ValueIE[I3, A3]](o3.AsIter()(ctx, source)))
		defer stop3()
		var mapping3 []A3

		ret, err := o0.AsModify()(ctx, func(index I0, focus A0) (A0, error) {
			v1, ok := next1()
			if !ok {
				return focus, nil
			}

			i1, a1, err := v1.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v2, ok := next2()
			if !ok {
				return focus, nil
			}

			i2, a2, err := v2.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v3, ok := next3()
			if !ok {
				return focus, nil
			}

			i3, a3, err := v3.Get()
			if err != nil {
				var a A0

				return a, err
			}

			mapped, err := fmap(lo.T4(index, i1, i2, i3), lo.T4(focus, a1, a2, a3))
			if err != nil {
				var a A0

				return a, err
			}

			mapping1 = append(mapping1, mapped.B)
			mapping2 = append(mapping2, mapped.C)
			mapping3 = append(mapping3, mapped.D)
			return mapped.A, err
		}, source)
		if err != nil {
			var s S

			return s, err
		}

		i := 0
		ret, err = o1.AsModify()(ctx, func(index I1, focus A1) (A1, error) {
			if i >= len(mapping1) {
				return focus, nil
			}

			return mapping1[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o2.AsModify()(ctx, func(index I2, focus A2) (A2, error) {
			if i >= len(mapping2) {
				return focus, nil
			}

			return mapping2[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o3.AsModify()(ctx, func(index I3, focus A3) (A3, error) {
			if i >= len(mapping3) {
				return focus, nil
			}

			return mapping3[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		return ret, nil
	}, nil, func(a lo.Tuple4[I0, I1, I2, I3], b lo.Tuple4[I0, I1, I2, I3]) bool {
		if o0.AsIxMatch()(a.A, b.A) != true {
			return false
		}

		if o1.AsIxMatch()(a.B, b.B) != true {
			return false
		}

		if o2.AsIxMatch()(a.C, b.C) != true {
			return false
		}

		if o3.AsIxMatch()(a.D, b.D) != true {
			return false
		}

		return true
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleOf{
			OpticTypeExpr: ot,
			Elements:      []expr.OpticExpression{o0.AsExpr(), o1.AsExpr(), o2.AsExpr(), o3.AsExpr()},
		}
	}, o0, o1, o2, o3))
}

// The T5Of Combinator constructs a [lo.Tuple5] whose elements are the focuses of the given [Optic]s
//
// Note: The number of focused tuples is limited by the optic that focuses the least elements.
func T5Of[I0 any, I1 any, I2 any, I3 any, I4 any, S any, A0 any, A1 any, A2 any, A3 any, A4 any, RET0 any, RW0 any, DIR0 any, ERR0 any, RET1 any, RW1 any, DIR1 any, ERR1 any, RET2 any, RW2 any, DIR2 any, ERR2 any, RET3 any, RW3 any, DIR3 any, ERR3 any, RET4 any, RW4 any, DIR4 any, ERR4 any](o0 Optic[I0, S, S, A0, A0, RET0, RW0, DIR0, ERR0], o1 Optic[I1, S, S, A1, A1, RET1, RW1, DIR1, ERR1], o2 Optic[I2, S, S, A2, A2, RET2, RW2, DIR2, ERR2], o3 Optic[I3, S, S, A3, A3, RET3, RW3, DIR3, ERR3], o4 Optic[I4, S, S, A4, A4, RET4, RW4, DIR4, ERR4]) Optic[lo.Tuple5[I0, I1, I2, I3, I4], S, S, lo.Tuple5[A0, A1, A2, A3, A4], lo.Tuple5[A0, A1, A2, A3, A4], CompositionTree[CompositionTree[CompositionTree[RET0, RET1], RET2], CompositionTree[RET3, RET4]], CompositionTree[CompositionTree[CompositionTree[RW0, RW1], RW2], CompositionTree[RW3, RW4]], UniDir, CompositionTree[CompositionTree[CompositionTree[ERR0, ERR1], ERR2], CompositionTree[ERR3, ERR4]]] {
	return CombiTraversal[CompositionTree[CompositionTree[CompositionTree[RET0, RET1], RET2], CompositionTree[RET3, RET4]], CompositionTree[CompositionTree[CompositionTree[RW0, RW1], RW2], CompositionTree[RW3, RW4]], CompositionTree[CompositionTree[CompositionTree[ERR0, ERR1], ERR2], CompositionTree[ERR3, ERR4]], lo.Tuple5[I0, I1, I2, I3, I4], S, S, lo.Tuple5[A0, A1, A2, A3, A4], lo.Tuple5[A0, A1, A2, A3, A4]](func(ctx context.Context, source S) SeqIE[lo.Tuple5[I0, I1, I2, I3, I4], lo.Tuple5[A0, A1, A2, A3, A4]] {
		return func(yield func(focusHello ValueIE[lo.Tuple5[I0, I1, I2, I3, I4], lo.Tuple5[A0, A1, A2, A3, A4]]) bool) {
			next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
			defer stop1()
			next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
			defer stop2()
			next3, stop3 := iter.Pull(iter.Seq[ValueIE[I3, A3]](o3.AsIter()(ctx, source)))
			defer stop3()
			next4, stop4 := iter.Pull(iter.Seq[ValueIE[I4, A4]](o4.AsIter()(ctx, source)))
			defer stop4()
			o0.AsIter()(ctx, source)(func(val ValueIE[I0, A0]) bool {
				i0, a0, err := val.Get()
				if err != nil {
					var i lo.Tuple5[I0, I1, I2, I3, I4]

					var a lo.Tuple5[A0, A1, A2, A3, A4]

					return yield(ValIE(i, a, err))
				}

				v1, ok := next1()
				if !ok {
					return false
				}

				i1, a1, err := v1.Get()
				if err != nil {
					var i lo.Tuple5[I0, I1, I2, I3, I4]

					var a lo.Tuple5[A0, A1, A2, A3, A4]

					return yield(ValIE(i, a, err))
				}

				v2, ok := next2()
				if !ok {
					return false
				}

				i2, a2, err := v2.Get()
				if err != nil {
					var i lo.Tuple5[I0, I1, I2, I3, I4]

					var a lo.Tuple5[A0, A1, A2, A3, A4]

					return yield(ValIE(i, a, err))
				}

				v3, ok := next3()
				if !ok {
					return false
				}

				i3, a3, err := v3.Get()
				if err != nil {
					var i lo.Tuple5[I0, I1, I2, I3, I4]

					var a lo.Tuple5[A0, A1, A2, A3, A4]

					return yield(ValIE(i, a, err))
				}

				v4, ok := next4()
				if !ok {
					return false
				}

				i4, a4, err := v4.Get()
				if err != nil {
					var i lo.Tuple5[I0, I1, I2, I3, I4]

					var a lo.Tuple5[A0, A1, A2, A3, A4]

					return yield(ValIE(i, a, err))
				}

				return yield(ValIE(lo.T5(i0, i1, i2, i3, i4), lo.T5(a0, a1, a2, a3, a4), nil))
			})
		}
	}, func(ctx context.Context, source S) (int, error) {
		l0, err := o0.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l1, err := o1.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l2, err := o2.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l3, err := o3.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l4, err := o4.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		return min(l0, l1, l2, l3, l4), nil
	}, func(ctx context.Context, fmap func(index lo.Tuple5[I0, I1, I2, I3, I4], focus lo.Tuple5[A0, A1, A2, A3, A4]) (lo.Tuple5[A0, A1, A2, A3, A4], error), source S) (S, error) {
		next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
		defer stop1()
		var mapping1 []A1

		next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
		defer stop2()
		var mapping2 []A2

		next3, stop3 := iter.Pull(iter.Seq[ValueIE[I3, A3]](o3.AsIter()(ctx, source)))
		defer stop3()
		var mapping3 []A3

		next4, stop4 := iter.Pull(iter.Seq[ValueIE[I4, A4]](o4.AsIter()(ctx, source)))
		defer stop4()
		var mapping4 []A4

		ret, err := o0.AsModify()(ctx, func(index I0, focus A0) (A0, error) {
			v1, ok := next1()
			if !ok {
				return focus, nil
			}

			i1, a1, err := v1.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v2, ok := next2()
			if !ok {
				return focus, nil
			}

			i2, a2, err := v2.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v3, ok := next3()
			if !ok {
				return focus, nil
			}

			i3, a3, err := v3.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v4, ok := next4()
			if !ok {
				return focus, nil
			}

			i4, a4, err := v4.Get()
			if err != nil {
				var a A0

				return a, err
			}

			mapped, err := fmap(lo.T5(index, i1, i2, i3, i4), lo.T5(focus, a1, a2, a3, a4))
			if err != nil {
				var a A0

				return a, err
			}

			mapping1 = append(mapping1, mapped.B)
			mapping2 = append(mapping2, mapped.C)
			mapping3 = append(mapping3, mapped.D)
			mapping4 = append(mapping4, mapped.E)
			return mapped.A, err
		}, source)
		if err != nil {
			var s S

			return s, err
		}

		i := 0
		ret, err = o1.AsModify()(ctx, func(index I1, focus A1) (A1, error) {
			if i >= len(mapping1) {
				return focus, nil
			}

			return mapping1[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o2.AsModify()(ctx, func(index I2, focus A2) (A2, error) {
			if i >= len(mapping2) {
				return focus, nil
			}

			return mapping2[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o3.AsModify()(ctx, func(index I3, focus A3) (A3, error) {
			if i >= len(mapping3) {
				return focus, nil
			}

			return mapping3[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o4.AsModify()(ctx, func(index I4, focus A4) (A4, error) {
			if i >= len(mapping4) {
				return focus, nil
			}

			return mapping4[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		return ret, nil
	}, nil, func(a lo.Tuple5[I0, I1, I2, I3, I4], b lo.Tuple5[I0, I1, I2, I3, I4]) bool {
		if o0.AsIxMatch()(a.A, b.A) != true {
			return false
		}

		if o1.AsIxMatch()(a.B, b.B) != true {
			return false
		}

		if o2.AsIxMatch()(a.C, b.C) != true {
			return false
		}

		if o3.AsIxMatch()(a.D, b.D) != true {
			return false
		}

		if o4.AsIxMatch()(a.E, b.E) != true {
			return false
		}

		return true
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleOf{
			OpticTypeExpr: ot,
			Elements:      []expr.OpticExpression{o0.AsExpr(), o1.AsExpr(), o2.AsExpr(), o3.AsExpr(), o4.AsExpr()},
		}
	}, o0, o1, o2, o3, o4))
}

// The T6Of Combinator constructs a [lo.Tuple6] whose elements are the focuses of the given [Optic]s
//
// Note: The number of focused tuples is limited by the optic that focuses the least elements.
func T6Of[I0 any, I1 any, I2 any, I3 any, I4 any, I5 any, S any, A0 any, A1 any, A2 any, A3 any, A4 any, A5 any, RET0 any, RW0 any, DIR0 any, ERR0 any, RET1 any, RW1 any, DIR1 any, ERR1 any, RET2 any, RW2 any, DIR2 any, ERR2 any, RET3 any, RW3 any, DIR3 any, ERR3 any, RET4 any, RW4 any, DIR4 any, ERR4 any, RET5 any, RW5 any, DIR5 any, ERR5 any](o0 Optic[I0, S, S, A0, A0, RET0, RW0, DIR0, ERR0], o1 Optic[I1, S, S, A1, A1, RET1, RW1, DIR1, ERR1], o2 Optic[I2, S, S, A2, A2, RET2, RW2, DIR2, ERR2], o3 Optic[I3, S, S, A3, A3, RET3, RW3, DIR3, ERR3], o4 Optic[I4, S, S, A4, A4, RET4, RW4, DIR4, ERR4], o5 Optic[I5, S, S, A5, A5, RET5, RW5, DIR5, ERR5]) Optic[lo.Tuple6[I0, I1, I2, I3, I4, I5], S, S, lo.Tuple6[A0, A1, A2, A3, A4, A5], lo.Tuple6[A0, A1, A2, A3, A4, A5], CompositionTree[CompositionTree[CompositionTree[RET0, RET1], RET2], CompositionTree[CompositionTree[RET3, RET4], RET5]], CompositionTree[CompositionTree[CompositionTree[RW0, RW1], RW2], CompositionTree[CompositionTree[RW3, RW4], RW5]], UniDir, CompositionTree[CompositionTree[CompositionTree[ERR0, ERR1], ERR2], CompositionTree[CompositionTree[ERR3, ERR4], ERR5]]] {
	return CombiTraversal[CompositionTree[CompositionTree[CompositionTree[RET0, RET1], RET2], CompositionTree[CompositionTree[RET3, RET4], RET5]], CompositionTree[CompositionTree[CompositionTree[RW0, RW1], RW2], CompositionTree[CompositionTree[RW3, RW4], RW5]], CompositionTree[CompositionTree[CompositionTree[ERR0, ERR1], ERR2], CompositionTree[CompositionTree[ERR3, ERR4], ERR5]], lo.Tuple6[I0, I1, I2, I3, I4, I5], S, S, lo.Tuple6[A0, A1, A2, A3, A4, A5], lo.Tuple6[A0, A1, A2, A3, A4, A5]](func(ctx context.Context, source S) SeqIE[lo.Tuple6[I0, I1, I2, I3, I4, I5], lo.Tuple6[A0, A1, A2, A3, A4, A5]] {
		return func(yield func(focusHello ValueIE[lo.Tuple6[I0, I1, I2, I3, I4, I5], lo.Tuple6[A0, A1, A2, A3, A4, A5]]) bool) {
			next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
			defer stop1()
			next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
			defer stop2()
			next3, stop3 := iter.Pull(iter.Seq[ValueIE[I3, A3]](o3.AsIter()(ctx, source)))
			defer stop3()
			next4, stop4 := iter.Pull(iter.Seq[ValueIE[I4, A4]](o4.AsIter()(ctx, source)))
			defer stop4()
			next5, stop5 := iter.Pull(iter.Seq[ValueIE[I5, A5]](o5.AsIter()(ctx, source)))
			defer stop5()
			o0.AsIter()(ctx, source)(func(val ValueIE[I0, A0]) bool {
				i0, a0, err := val.Get()
				if err != nil {
					var i lo.Tuple6[I0, I1, I2, I3, I4, I5]

					var a lo.Tuple6[A0, A1, A2, A3, A4, A5]

					return yield(ValIE(i, a, err))
				}

				v1, ok := next1()
				if !ok {
					return false
				}

				i1, a1, err := v1.Get()
				if err != nil {
					var i lo.Tuple6[I0, I1, I2, I3, I4, I5]

					var a lo.Tuple6[A0, A1, A2, A3, A4, A5]

					return yield(ValIE(i, a, err))
				}

				v2, ok := next2()
				if !ok {
					return false
				}

				i2, a2, err := v2.Get()
				if err != nil {
					var i lo.Tuple6[I0, I1, I2, I3, I4, I5]

					var a lo.Tuple6[A0, A1, A2, A3, A4, A5]

					return yield(ValIE(i, a, err))
				}

				v3, ok := next3()
				if !ok {
					return false
				}

				i3, a3, err := v3.Get()
				if err != nil {
					var i lo.Tuple6[I0, I1, I2, I3, I4, I5]

					var a lo.Tuple6[A0, A1, A2, A3, A4, A5]

					return yield(ValIE(i, a, err))
				}

				v4, ok := next4()
				if !ok {
					return false
				}

				i4, a4, err := v4.Get()
				if err != nil {
					var i lo.Tuple6[I0, I1, I2, I3, I4, I5]

					var a lo.Tuple6[A0, A1, A2, A3, A4, A5]

					return yield(ValIE(i, a, err))
				}

				v5, ok := next5()
				if !ok {
					return false
				}

				i5, a5, err := v5.Get()
				if err != nil {
					var i lo.Tuple6[I0, I1, I2, I3, I4, I5]

					var a lo.Tuple6[A0, A1, A2, A3, A4, A5]

					return yield(ValIE(i, a, err))
				}

				return yield(ValIE(lo.T6(i0, i1, i2, i3, i4, i5), lo.T6(a0, a1, a2, a3, a4, a5), nil))
			})
		}
	}, func(ctx context.Context, source S) (int, error) {
		l0, err := o0.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l1, err := o1.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l2, err := o2.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l3, err := o3.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l4, err := o4.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l5, err := o5.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		return min(l0, l1, l2, l3, l4, l5), nil
	}, func(ctx context.Context, fmap func(index lo.Tuple6[I0, I1, I2, I3, I4, I5], focus lo.Tuple6[A0, A1, A2, A3, A4, A5]) (lo.Tuple6[A0, A1, A2, A3, A4, A5], error), source S) (S, error) {
		next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
		defer stop1()
		var mapping1 []A1

		next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
		defer stop2()
		var mapping2 []A2

		next3, stop3 := iter.Pull(iter.Seq[ValueIE[I3, A3]](o3.AsIter()(ctx, source)))
		defer stop3()
		var mapping3 []A3

		next4, stop4 := iter.Pull(iter.Seq[ValueIE[I4, A4]](o4.AsIter()(ctx, source)))
		defer stop4()
		var mapping4 []A4

		next5, stop5 := iter.Pull(iter.Seq[ValueIE[I5, A5]](o5.AsIter()(ctx, source)))
		defer stop5()
		var mapping5 []A5

		ret, err := o0.AsModify()(ctx, func(index I0, focus A0) (A0, error) {
			v1, ok := next1()
			if !ok {
				return focus, nil
			}

			i1, a1, err := v1.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v2, ok := next2()
			if !ok {
				return focus, nil
			}

			i2, a2, err := v2.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v3, ok := next3()
			if !ok {
				return focus, nil
			}

			i3, a3, err := v3.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v4, ok := next4()
			if !ok {
				return focus, nil
			}

			i4, a4, err := v4.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v5, ok := next5()
			if !ok {
				return focus, nil
			}

			i5, a5, err := v5.Get()
			if err != nil {
				var a A0

				return a, err
			}

			mapped, err := fmap(lo.T6(index, i1, i2, i3, i4, i5), lo.T6(focus, a1, a2, a3, a4, a5))
			if err != nil {
				var a A0

				return a, err
			}

			mapping1 = append(mapping1, mapped.B)
			mapping2 = append(mapping2, mapped.C)
			mapping3 = append(mapping3, mapped.D)
			mapping4 = append(mapping4, mapped.E)
			mapping5 = append(mapping5, mapped.F)
			return mapped.A, err
		}, source)
		if err != nil {
			var s S

			return s, err
		}

		i := 0
		ret, err = o1.AsModify()(ctx, func(index I1, focus A1) (A1, error) {
			if i >= len(mapping1) {
				return focus, nil
			}

			return mapping1[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o2.AsModify()(ctx, func(index I2, focus A2) (A2, error) {
			if i >= len(mapping2) {
				return focus, nil
			}

			return mapping2[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o3.AsModify()(ctx, func(index I3, focus A3) (A3, error) {
			if i >= len(mapping3) {
				return focus, nil
			}

			return mapping3[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o4.AsModify()(ctx, func(index I4, focus A4) (A4, error) {
			if i >= len(mapping4) {
				return focus, nil
			}

			return mapping4[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o5.AsModify()(ctx, func(index I5, focus A5) (A5, error) {
			if i >= len(mapping5) {
				return focus, nil
			}

			return mapping5[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		return ret, nil
	}, nil, func(a lo.Tuple6[I0, I1, I2, I3, I4, I5], b lo.Tuple6[I0, I1, I2, I3, I4, I5]) bool {
		if o0.AsIxMatch()(a.A, b.A) != true {
			return false
		}

		if o1.AsIxMatch()(a.B, b.B) != true {
			return false
		}

		if o2.AsIxMatch()(a.C, b.C) != true {
			return false
		}

		if o3.AsIxMatch()(a.D, b.D) != true {
			return false
		}

		if o4.AsIxMatch()(a.E, b.E) != true {
			return false
		}

		if o5.AsIxMatch()(a.F, b.F) != true {
			return false
		}

		return true
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleOf{
			OpticTypeExpr: ot,
			Elements:      []expr.OpticExpression{o0.AsExpr(), o1.AsExpr(), o2.AsExpr(), o3.AsExpr(), o4.AsExpr(), o5.AsExpr()},
		}
	}, o0, o1, o2, o3, o4, o5))
}

// The T7Of Combinator constructs a [lo.Tuple7] whose elements are the focuses of the given [Optic]s
//
// Note: The number of focused tuples is limited by the optic that focuses the least elements.
func T7Of[I0 any, I1 any, I2 any, I3 any, I4 any, I5 any, I6 any, S any, A0 any, A1 any, A2 any, A3 any, A4 any, A5 any, A6 any, RET0 any, RW0 any, DIR0 any, ERR0 any, RET1 any, RW1 any, DIR1 any, ERR1 any, RET2 any, RW2 any, DIR2 any, ERR2 any, RET3 any, RW3 any, DIR3 any, ERR3 any, RET4 any, RW4 any, DIR4 any, ERR4 any, RET5 any, RW5 any, DIR5 any, ERR5 any, RET6 any, RW6 any, DIR6 any, ERR6 any](o0 Optic[I0, S, S, A0, A0, RET0, RW0, DIR0, ERR0], o1 Optic[I1, S, S, A1, A1, RET1, RW1, DIR1, ERR1], o2 Optic[I2, S, S, A2, A2, RET2, RW2, DIR2, ERR2], o3 Optic[I3, S, S, A3, A3, RET3, RW3, DIR3, ERR3], o4 Optic[I4, S, S, A4, A4, RET4, RW4, DIR4, ERR4], o5 Optic[I5, S, S, A5, A5, RET5, RW5, DIR5, ERR5], o6 Optic[I6, S, S, A6, A6, RET6, RW6, DIR6, ERR6]) Optic[lo.Tuple7[I0, I1, I2, I3, I4, I5, I6], S, S, lo.Tuple7[A0, A1, A2, A3, A4, A5, A6], lo.Tuple7[A0, A1, A2, A3, A4, A5, A6], CompositionTree[CompositionTree[CompositionTree[RET0, RET1], CompositionTree[RET2, RET3]], CompositionTree[CompositionTree[RET4, RET5], RET6]], CompositionTree[CompositionTree[CompositionTree[RW0, RW1], CompositionTree[RW2, RW3]], CompositionTree[CompositionTree[RW4, RW5], RW6]], UniDir, CompositionTree[CompositionTree[CompositionTree[ERR0, ERR1], CompositionTree[ERR2, ERR3]], CompositionTree[CompositionTree[ERR4, ERR5], ERR6]]] {
	return CombiTraversal[CompositionTree[CompositionTree[CompositionTree[RET0, RET1], CompositionTree[RET2, RET3]], CompositionTree[CompositionTree[RET4, RET5], RET6]], CompositionTree[CompositionTree[CompositionTree[RW0, RW1], CompositionTree[RW2, RW3]], CompositionTree[CompositionTree[RW4, RW5], RW6]], CompositionTree[CompositionTree[CompositionTree[ERR0, ERR1], CompositionTree[ERR2, ERR3]], CompositionTree[CompositionTree[ERR4, ERR5], ERR6]], lo.Tuple7[I0, I1, I2, I3, I4, I5, I6], S, S, lo.Tuple7[A0, A1, A2, A3, A4, A5, A6], lo.Tuple7[A0, A1, A2, A3, A4, A5, A6]](func(ctx context.Context, source S) SeqIE[lo.Tuple7[I0, I1, I2, I3, I4, I5, I6], lo.Tuple7[A0, A1, A2, A3, A4, A5, A6]] {
		return func(yield func(focusHello ValueIE[lo.Tuple7[I0, I1, I2, I3, I4, I5, I6], lo.Tuple7[A0, A1, A2, A3, A4, A5, A6]]) bool) {
			next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
			defer stop1()
			next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
			defer stop2()
			next3, stop3 := iter.Pull(iter.Seq[ValueIE[I3, A3]](o3.AsIter()(ctx, source)))
			defer stop3()
			next4, stop4 := iter.Pull(iter.Seq[ValueIE[I4, A4]](o4.AsIter()(ctx, source)))
			defer stop4()
			next5, stop5 := iter.Pull(iter.Seq[ValueIE[I5, A5]](o5.AsIter()(ctx, source)))
			defer stop5()
			next6, stop6 := iter.Pull(iter.Seq[ValueIE[I6, A6]](o6.AsIter()(ctx, source)))
			defer stop6()
			o0.AsIter()(ctx, source)(func(val ValueIE[I0, A0]) bool {
				i0, a0, err := val.Get()
				if err != nil {
					var i lo.Tuple7[I0, I1, I2, I3, I4, I5, I6]

					var a lo.Tuple7[A0, A1, A2, A3, A4, A5, A6]

					return yield(ValIE(i, a, err))
				}

				v1, ok := next1()
				if !ok {
					return false
				}

				i1, a1, err := v1.Get()
				if err != nil {
					var i lo.Tuple7[I0, I1, I2, I3, I4, I5, I6]

					var a lo.Tuple7[A0, A1, A2, A3, A4, A5, A6]

					return yield(ValIE(i, a, err))
				}

				v2, ok := next2()
				if !ok {
					return false
				}

				i2, a2, err := v2.Get()
				if err != nil {
					var i lo.Tuple7[I0, I1, I2, I3, I4, I5, I6]

					var a lo.Tuple7[A0, A1, A2, A3, A4, A5, A6]

					return yield(ValIE(i, a, err))
				}

				v3, ok := next3()
				if !ok {
					return false
				}

				i3, a3, err := v3.Get()
				if err != nil {
					var i lo.Tuple7[I0, I1, I2, I3, I4, I5, I6]

					var a lo.Tuple7[A0, A1, A2, A3, A4, A5, A6]

					return yield(ValIE(i, a, err))
				}

				v4, ok := next4()
				if !ok {
					return false
				}

				i4, a4, err := v4.Get()
				if err != nil {
					var i lo.Tuple7[I0, I1, I2, I3, I4, I5, I6]

					var a lo.Tuple7[A0, A1, A2, A3, A4, A5, A6]

					return yield(ValIE(i, a, err))
				}

				v5, ok := next5()
				if !ok {
					return false
				}

				i5, a5, err := v5.Get()
				if err != nil {
					var i lo.Tuple7[I0, I1, I2, I3, I4, I5, I6]

					var a lo.Tuple7[A0, A1, A2, A3, A4, A5, A6]

					return yield(ValIE(i, a, err))
				}

				v6, ok := next6()
				if !ok {
					return false
				}

				i6, a6, err := v6.Get()
				if err != nil {
					var i lo.Tuple7[I0, I1, I2, I3, I4, I5, I6]

					var a lo.Tuple7[A0, A1, A2, A3, A4, A5, A6]

					return yield(ValIE(i, a, err))
				}

				return yield(ValIE(lo.T7(i0, i1, i2, i3, i4, i5, i6), lo.T7(a0, a1, a2, a3, a4, a5, a6), nil))
			})
		}
	}, func(ctx context.Context, source S) (int, error) {
		l0, err := o0.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l1, err := o1.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l2, err := o2.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l3, err := o3.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l4, err := o4.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l5, err := o5.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l6, err := o6.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		return min(l0, l1, l2, l3, l4, l5, l6), nil
	}, func(ctx context.Context, fmap func(index lo.Tuple7[I0, I1, I2, I3, I4, I5, I6], focus lo.Tuple7[A0, A1, A2, A3, A4, A5, A6]) (lo.Tuple7[A0, A1, A2, A3, A4, A5, A6], error), source S) (S, error) {
		next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
		defer stop1()
		var mapping1 []A1

		next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
		defer stop2()
		var mapping2 []A2

		next3, stop3 := iter.Pull(iter.Seq[ValueIE[I3, A3]](o3.AsIter()(ctx, source)))
		defer stop3()
		var mapping3 []A3

		next4, stop4 := iter.Pull(iter.Seq[ValueIE[I4, A4]](o4.AsIter()(ctx, source)))
		defer stop4()
		var mapping4 []A4

		next5, stop5 := iter.Pull(iter.Seq[ValueIE[I5, A5]](o5.AsIter()(ctx, source)))
		defer stop5()
		var mapping5 []A5

		next6, stop6 := iter.Pull(iter.Seq[ValueIE[I6, A6]](o6.AsIter()(ctx, source)))
		defer stop6()
		var mapping6 []A6

		ret, err := o0.AsModify()(ctx, func(index I0, focus A0) (A0, error) {
			v1, ok := next1()
			if !ok {
				return focus, nil
			}

			i1, a1, err := v1.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v2, ok := next2()
			if !ok {
				return focus, nil
			}

			i2, a2, err := v2.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v3, ok := next3()
			if !ok {
				return focus, nil
			}

			i3, a3, err := v3.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v4, ok := next4()
			if !ok {
				return focus, nil
			}

			i4, a4, err := v4.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v5, ok := next5()
			if !ok {
				return focus, nil
			}

			i5, a5, err := v5.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v6, ok := next6()
			if !ok {
				return focus, nil
			}

			i6, a6, err := v6.Get()
			if err != nil {
				var a A0

				return a, err
			}

			mapped, err := fmap(lo.T7(index, i1, i2, i3, i4, i5, i6), lo.T7(focus, a1, a2, a3, a4, a5, a6))
			if err != nil {
				var a A0

				return a, err
			}

			mapping1 = append(mapping1, mapped.B)
			mapping2 = append(mapping2, mapped.C)
			mapping3 = append(mapping3, mapped.D)
			mapping4 = append(mapping4, mapped.E)
			mapping5 = append(mapping5, mapped.F)
			mapping6 = append(mapping6, mapped.G)
			return mapped.A, err
		}, source)
		if err != nil {
			var s S

			return s, err
		}

		i := 0
		ret, err = o1.AsModify()(ctx, func(index I1, focus A1) (A1, error) {
			if i >= len(mapping1) {
				return focus, nil
			}

			return mapping1[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o2.AsModify()(ctx, func(index I2, focus A2) (A2, error) {
			if i >= len(mapping2) {
				return focus, nil
			}

			return mapping2[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o3.AsModify()(ctx, func(index I3, focus A3) (A3, error) {
			if i >= len(mapping3) {
				return focus, nil
			}

			return mapping3[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o4.AsModify()(ctx, func(index I4, focus A4) (A4, error) {
			if i >= len(mapping4) {
				return focus, nil
			}

			return mapping4[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o5.AsModify()(ctx, func(index I5, focus A5) (A5, error) {
			if i >= len(mapping5) {
				return focus, nil
			}

			return mapping5[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o6.AsModify()(ctx, func(index I6, focus A6) (A6, error) {
			if i >= len(mapping6) {
				return focus, nil
			}

			return mapping6[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		return ret, nil
	}, nil, func(a lo.Tuple7[I0, I1, I2, I3, I4, I5, I6], b lo.Tuple7[I0, I1, I2, I3, I4, I5, I6]) bool {
		if o0.AsIxMatch()(a.A, b.A) != true {
			return false
		}

		if o1.AsIxMatch()(a.B, b.B) != true {
			return false
		}

		if o2.AsIxMatch()(a.C, b.C) != true {
			return false
		}

		if o3.AsIxMatch()(a.D, b.D) != true {
			return false
		}

		if o4.AsIxMatch()(a.E, b.E) != true {
			return false
		}

		if o5.AsIxMatch()(a.F, b.F) != true {
			return false
		}

		if o6.AsIxMatch()(a.G, b.G) != true {
			return false
		}

		return true
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleOf{
			OpticTypeExpr: ot,
			Elements:      []expr.OpticExpression{o0.AsExpr(), o1.AsExpr(), o2.AsExpr(), o3.AsExpr(), o4.AsExpr(), o5.AsExpr(), o6.AsExpr()},
		}
	}, o0, o1, o2, o3, o4, o5, o6))
}

// The T8Of Combinator constructs a [lo.Tuple8] whose elements are the focuses of the given [Optic]s
//
// Note: The number of focused tuples is limited by the optic that focuses the least elements.
func T8Of[I0 any, I1 any, I2 any, I3 any, I4 any, I5 any, I6 any, I7 any, S any, A0 any, A1 any, A2 any, A3 any, A4 any, A5 any, A6 any, A7 any, RET0 any, RW0 any, DIR0 any, ERR0 any, RET1 any, RW1 any, DIR1 any, ERR1 any, RET2 any, RW2 any, DIR2 any, ERR2 any, RET3 any, RW3 any, DIR3 any, ERR3 any, RET4 any, RW4 any, DIR4 any, ERR4 any, RET5 any, RW5 any, DIR5 any, ERR5 any, RET6 any, RW6 any, DIR6 any, ERR6 any, RET7 any, RW7 any, DIR7 any, ERR7 any](o0 Optic[I0, S, S, A0, A0, RET0, RW0, DIR0, ERR0], o1 Optic[I1, S, S, A1, A1, RET1, RW1, DIR1, ERR1], o2 Optic[I2, S, S, A2, A2, RET2, RW2, DIR2, ERR2], o3 Optic[I3, S, S, A3, A3, RET3, RW3, DIR3, ERR3], o4 Optic[I4, S, S, A4, A4, RET4, RW4, DIR4, ERR4], o5 Optic[I5, S, S, A5, A5, RET5, RW5, DIR5, ERR5], o6 Optic[I6, S, S, A6, A6, RET6, RW6, DIR6, ERR6], o7 Optic[I7, S, S, A7, A7, RET7, RW7, DIR7, ERR7]) Optic[lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7], S, S, lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7], lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7], CompositionTree[CompositionTree[CompositionTree[RET0, RET1], CompositionTree[RET2, RET3]], CompositionTree[CompositionTree[RET4, RET5], CompositionTree[RET6, RET7]]], CompositionTree[CompositionTree[CompositionTree[RW0, RW1], CompositionTree[RW2, RW3]], CompositionTree[CompositionTree[RW4, RW5], CompositionTree[RW6, RW7]]], UniDir, CompositionTree[CompositionTree[CompositionTree[ERR0, ERR1], CompositionTree[ERR2, ERR3]], CompositionTree[CompositionTree[ERR4, ERR5], CompositionTree[ERR6, ERR7]]]] {
	return CombiTraversal[CompositionTree[CompositionTree[CompositionTree[RET0, RET1], CompositionTree[RET2, RET3]], CompositionTree[CompositionTree[RET4, RET5], CompositionTree[RET6, RET7]]], CompositionTree[CompositionTree[CompositionTree[RW0, RW1], CompositionTree[RW2, RW3]], CompositionTree[CompositionTree[RW4, RW5], CompositionTree[RW6, RW7]]], CompositionTree[CompositionTree[CompositionTree[ERR0, ERR1], CompositionTree[ERR2, ERR3]], CompositionTree[CompositionTree[ERR4, ERR5], CompositionTree[ERR6, ERR7]]], lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7], S, S, lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7], lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7]](func(ctx context.Context, source S) SeqIE[lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7], lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7]] {
		return func(yield func(focusHello ValueIE[lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7], lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7]]) bool) {
			next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
			defer stop1()
			next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
			defer stop2()
			next3, stop3 := iter.Pull(iter.Seq[ValueIE[I3, A3]](o3.AsIter()(ctx, source)))
			defer stop3()
			next4, stop4 := iter.Pull(iter.Seq[ValueIE[I4, A4]](o4.AsIter()(ctx, source)))
			defer stop4()
			next5, stop5 := iter.Pull(iter.Seq[ValueIE[I5, A5]](o5.AsIter()(ctx, source)))
			defer stop5()
			next6, stop6 := iter.Pull(iter.Seq[ValueIE[I6, A6]](o6.AsIter()(ctx, source)))
			defer stop6()
			next7, stop7 := iter.Pull(iter.Seq[ValueIE[I7, A7]](o7.AsIter()(ctx, source)))
			defer stop7()
			o0.AsIter()(ctx, source)(func(val ValueIE[I0, A0]) bool {
				i0, a0, err := val.Get()
				if err != nil {
					var i lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7]

					var a lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7]

					return yield(ValIE(i, a, err))
				}

				v1, ok := next1()
				if !ok {
					return false
				}

				i1, a1, err := v1.Get()
				if err != nil {
					var i lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7]

					var a lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7]

					return yield(ValIE(i, a, err))
				}

				v2, ok := next2()
				if !ok {
					return false
				}

				i2, a2, err := v2.Get()
				if err != nil {
					var i lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7]

					var a lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7]

					return yield(ValIE(i, a, err))
				}

				v3, ok := next3()
				if !ok {
					return false
				}

				i3, a3, err := v3.Get()
				if err != nil {
					var i lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7]

					var a lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7]

					return yield(ValIE(i, a, err))
				}

				v4, ok := next4()
				if !ok {
					return false
				}

				i4, a4, err := v4.Get()
				if err != nil {
					var i lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7]

					var a lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7]

					return yield(ValIE(i, a, err))
				}

				v5, ok := next5()
				if !ok {
					return false
				}

				i5, a5, err := v5.Get()
				if err != nil {
					var i lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7]

					var a lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7]

					return yield(ValIE(i, a, err))
				}

				v6, ok := next6()
				if !ok {
					return false
				}

				i6, a6, err := v6.Get()
				if err != nil {
					var i lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7]

					var a lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7]

					return yield(ValIE(i, a, err))
				}

				v7, ok := next7()
				if !ok {
					return false
				}

				i7, a7, err := v7.Get()
				if err != nil {
					var i lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7]

					var a lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7]

					return yield(ValIE(i, a, err))
				}

				return yield(ValIE(lo.T8(i0, i1, i2, i3, i4, i5, i6, i7), lo.T8(a0, a1, a2, a3, a4, a5, a6, a7), nil))
			})
		}
	}, func(ctx context.Context, source S) (int, error) {
		l0, err := o0.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l1, err := o1.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l2, err := o2.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l3, err := o3.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l4, err := o4.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l5, err := o5.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l6, err := o6.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l7, err := o7.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		return min(l0, l1, l2, l3, l4, l5, l6, l7), nil
	}, func(ctx context.Context, fmap func(index lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7], focus lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7]) (lo.Tuple8[A0, A1, A2, A3, A4, A5, A6, A7], error), source S) (S, error) {
		next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
		defer stop1()
		var mapping1 []A1

		next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
		defer stop2()
		var mapping2 []A2

		next3, stop3 := iter.Pull(iter.Seq[ValueIE[I3, A3]](o3.AsIter()(ctx, source)))
		defer stop3()
		var mapping3 []A3

		next4, stop4 := iter.Pull(iter.Seq[ValueIE[I4, A4]](o4.AsIter()(ctx, source)))
		defer stop4()
		var mapping4 []A4

		next5, stop5 := iter.Pull(iter.Seq[ValueIE[I5, A5]](o5.AsIter()(ctx, source)))
		defer stop5()
		var mapping5 []A5

		next6, stop6 := iter.Pull(iter.Seq[ValueIE[I6, A6]](o6.AsIter()(ctx, source)))
		defer stop6()
		var mapping6 []A6

		next7, stop7 := iter.Pull(iter.Seq[ValueIE[I7, A7]](o7.AsIter()(ctx, source)))
		defer stop7()
		var mapping7 []A7

		ret, err := o0.AsModify()(ctx, func(index I0, focus A0) (A0, error) {
			v1, ok := next1()
			if !ok {
				return focus, nil
			}

			i1, a1, err := v1.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v2, ok := next2()
			if !ok {
				return focus, nil
			}

			i2, a2, err := v2.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v3, ok := next3()
			if !ok {
				return focus, nil
			}

			i3, a3, err := v3.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v4, ok := next4()
			if !ok {
				return focus, nil
			}

			i4, a4, err := v4.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v5, ok := next5()
			if !ok {
				return focus, nil
			}

			i5, a5, err := v5.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v6, ok := next6()
			if !ok {
				return focus, nil
			}

			i6, a6, err := v6.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v7, ok := next7()
			if !ok {
				return focus, nil
			}

			i7, a7, err := v7.Get()
			if err != nil {
				var a A0

				return a, err
			}

			mapped, err := fmap(lo.T8(index, i1, i2, i3, i4, i5, i6, i7), lo.T8(focus, a1, a2, a3, a4, a5, a6, a7))
			if err != nil {
				var a A0

				return a, err
			}

			mapping1 = append(mapping1, mapped.B)
			mapping2 = append(mapping2, mapped.C)
			mapping3 = append(mapping3, mapped.D)
			mapping4 = append(mapping4, mapped.E)
			mapping5 = append(mapping5, mapped.F)
			mapping6 = append(mapping6, mapped.G)
			mapping7 = append(mapping7, mapped.H)
			return mapped.A, err
		}, source)
		if err != nil {
			var s S

			return s, err
		}

		i := 0
		ret, err = o1.AsModify()(ctx, func(index I1, focus A1) (A1, error) {
			if i >= len(mapping1) {
				return focus, nil
			}

			return mapping1[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o2.AsModify()(ctx, func(index I2, focus A2) (A2, error) {
			if i >= len(mapping2) {
				return focus, nil
			}

			return mapping2[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o3.AsModify()(ctx, func(index I3, focus A3) (A3, error) {
			if i >= len(mapping3) {
				return focus, nil
			}

			return mapping3[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o4.AsModify()(ctx, func(index I4, focus A4) (A4, error) {
			if i >= len(mapping4) {
				return focus, nil
			}

			return mapping4[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o5.AsModify()(ctx, func(index I5, focus A5) (A5, error) {
			if i >= len(mapping5) {
				return focus, nil
			}

			return mapping5[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o6.AsModify()(ctx, func(index I6, focus A6) (A6, error) {
			if i >= len(mapping6) {
				return focus, nil
			}

			return mapping6[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o7.AsModify()(ctx, func(index I7, focus A7) (A7, error) {
			if i >= len(mapping7) {
				return focus, nil
			}

			return mapping7[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		return ret, nil
	}, nil, func(a lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7], b lo.Tuple8[I0, I1, I2, I3, I4, I5, I6, I7]) bool {
		if o0.AsIxMatch()(a.A, b.A) != true {
			return false
		}

		if o1.AsIxMatch()(a.B, b.B) != true {
			return false
		}

		if o2.AsIxMatch()(a.C, b.C) != true {
			return false
		}

		if o3.AsIxMatch()(a.D, b.D) != true {
			return false
		}

		if o4.AsIxMatch()(a.E, b.E) != true {
			return false
		}

		if o5.AsIxMatch()(a.F, b.F) != true {
			return false
		}

		if o6.AsIxMatch()(a.G, b.G) != true {
			return false
		}

		if o7.AsIxMatch()(a.H, b.H) != true {
			return false
		}

		return true
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleOf{
			OpticTypeExpr: ot,
			Elements:      []expr.OpticExpression{o0.AsExpr(), o1.AsExpr(), o2.AsExpr(), o3.AsExpr(), o4.AsExpr(), o5.AsExpr(), o6.AsExpr(), o7.AsExpr()},
		}
	}, o0, o1, o2, o3, o4, o5, o6, o7))
}

// The T9Of Combinator constructs a [lo.Tuple9] whose elements are the focuses of the given [Optic]s
//
// Note: The number of focused tuples is limited by the optic that focuses the least elements.
func T9Of[I0 any, I1 any, I2 any, I3 any, I4 any, I5 any, I6 any, I7 any, I8 any, S any, A0 any, A1 any, A2 any, A3 any, A4 any, A5 any, A6 any, A7 any, A8 any, RET0 any, RW0 any, DIR0 any, ERR0 any, RET1 any, RW1 any, DIR1 any, ERR1 any, RET2 any, RW2 any, DIR2 any, ERR2 any, RET3 any, RW3 any, DIR3 any, ERR3 any, RET4 any, RW4 any, DIR4 any, ERR4 any, RET5 any, RW5 any, DIR5 any, ERR5 any, RET6 any, RW6 any, DIR6 any, ERR6 any, RET7 any, RW7 any, DIR7 any, ERR7 any, RET8 any, RW8 any, DIR8 any, ERR8 any](o0 Optic[I0, S, S, A0, A0, RET0, RW0, DIR0, ERR0], o1 Optic[I1, S, S, A1, A1, RET1, RW1, DIR1, ERR1], o2 Optic[I2, S, S, A2, A2, RET2, RW2, DIR2, ERR2], o3 Optic[I3, S, S, A3, A3, RET3, RW3, DIR3, ERR3], o4 Optic[I4, S, S, A4, A4, RET4, RW4, DIR4, ERR4], o5 Optic[I5, S, S, A5, A5, RET5, RW5, DIR5, ERR5], o6 Optic[I6, S, S, A6, A6, RET6, RW6, DIR6, ERR6], o7 Optic[I7, S, S, A7, A7, RET7, RW7, DIR7, ERR7], o8 Optic[I8, S, S, A8, A8, RET8, RW8, DIR8, ERR8]) Optic[lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8], S, S, lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8], lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8], CompositionTree[CompositionTree[CompositionTree[CompositionTree[RET0, RET1], RET2], CompositionTree[RET3, RET4]], CompositionTree[CompositionTree[RET5, RET6], CompositionTree[RET7, RET8]]], CompositionTree[CompositionTree[CompositionTree[CompositionTree[RW0, RW1], RW2], CompositionTree[RW3, RW4]], CompositionTree[CompositionTree[RW5, RW6], CompositionTree[RW7, RW8]]], UniDir, CompositionTree[CompositionTree[CompositionTree[CompositionTree[ERR0, ERR1], ERR2], CompositionTree[ERR3, ERR4]], CompositionTree[CompositionTree[ERR5, ERR6], CompositionTree[ERR7, ERR8]]]] {
	return CombiTraversal[CompositionTree[CompositionTree[CompositionTree[CompositionTree[RET0, RET1], RET2], CompositionTree[RET3, RET4]], CompositionTree[CompositionTree[RET5, RET6], CompositionTree[RET7, RET8]]], CompositionTree[CompositionTree[CompositionTree[CompositionTree[RW0, RW1], RW2], CompositionTree[RW3, RW4]], CompositionTree[CompositionTree[RW5, RW6], CompositionTree[RW7, RW8]]], CompositionTree[CompositionTree[CompositionTree[CompositionTree[ERR0, ERR1], ERR2], CompositionTree[ERR3, ERR4]], CompositionTree[CompositionTree[ERR5, ERR6], CompositionTree[ERR7, ERR8]]], lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8], S, S, lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8], lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8]](func(ctx context.Context, source S) SeqIE[lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8], lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8]] {
		return func(yield func(focusHello ValueIE[lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8], lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8]]) bool) {
			next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
			defer stop1()
			next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
			defer stop2()
			next3, stop3 := iter.Pull(iter.Seq[ValueIE[I3, A3]](o3.AsIter()(ctx, source)))
			defer stop3()
			next4, stop4 := iter.Pull(iter.Seq[ValueIE[I4, A4]](o4.AsIter()(ctx, source)))
			defer stop4()
			next5, stop5 := iter.Pull(iter.Seq[ValueIE[I5, A5]](o5.AsIter()(ctx, source)))
			defer stop5()
			next6, stop6 := iter.Pull(iter.Seq[ValueIE[I6, A6]](o6.AsIter()(ctx, source)))
			defer stop6()
			next7, stop7 := iter.Pull(iter.Seq[ValueIE[I7, A7]](o7.AsIter()(ctx, source)))
			defer stop7()
			next8, stop8 := iter.Pull(iter.Seq[ValueIE[I8, A8]](o8.AsIter()(ctx, source)))
			defer stop8()
			o0.AsIter()(ctx, source)(func(val ValueIE[I0, A0]) bool {
				i0, a0, err := val.Get()
				if err != nil {
					var i lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8]

					var a lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8]

					return yield(ValIE(i, a, err))
				}

				v1, ok := next1()
				if !ok {
					return false
				}

				i1, a1, err := v1.Get()
				if err != nil {
					var i lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8]

					var a lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8]

					return yield(ValIE(i, a, err))
				}

				v2, ok := next2()
				if !ok {
					return false
				}

				i2, a2, err := v2.Get()
				if err != nil {
					var i lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8]

					var a lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8]

					return yield(ValIE(i, a, err))
				}

				v3, ok := next3()
				if !ok {
					return false
				}

				i3, a3, err := v3.Get()
				if err != nil {
					var i lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8]

					var a lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8]

					return yield(ValIE(i, a, err))
				}

				v4, ok := next4()
				if !ok {
					return false
				}

				i4, a4, err := v4.Get()
				if err != nil {
					var i lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8]

					var a lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8]

					return yield(ValIE(i, a, err))
				}

				v5, ok := next5()
				if !ok {
					return false
				}

				i5, a5, err := v5.Get()
				if err != nil {
					var i lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8]

					var a lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8]

					return yield(ValIE(i, a, err))
				}

				v6, ok := next6()
				if !ok {
					return false
				}

				i6, a6, err := v6.Get()
				if err != nil {
					var i lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8]

					var a lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8]

					return yield(ValIE(i, a, err))
				}

				v7, ok := next7()
				if !ok {
					return false
				}

				i7, a7, err := v7.Get()
				if err != nil {
					var i lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8]

					var a lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8]

					return yield(ValIE(i, a, err))
				}

				v8, ok := next8()
				if !ok {
					return false
				}

				i8, a8, err := v8.Get()
				if err != nil {
					var i lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8]

					var a lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8]

					return yield(ValIE(i, a, err))
				}

				return yield(ValIE(lo.T9(i0, i1, i2, i3, i4, i5, i6, i7, i8), lo.T9(a0, a1, a2, a3, a4, a5, a6, a7, a8), nil))
			})
		}
	}, func(ctx context.Context, source S) (int, error) {
		l0, err := o0.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l1, err := o1.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l2, err := o2.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l3, err := o3.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l4, err := o4.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l5, err := o5.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l6, err := o6.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l7, err := o7.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		l8, err := o8.AsLengthGetter()(ctx, source)
		if err != nil {
			return 0, err
		}

		return min(l0, l1, l2, l3, l4, l5, l6, l7, l8), nil
	}, func(ctx context.Context, fmap func(index lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8], focus lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8]) (lo.Tuple9[A0, A1, A2, A3, A4, A5, A6, A7, A8], error), source S) (S, error) {
		next1, stop1 := iter.Pull(iter.Seq[ValueIE[I1, A1]](o1.AsIter()(ctx, source)))
		defer stop1()
		var mapping1 []A1

		next2, stop2 := iter.Pull(iter.Seq[ValueIE[I2, A2]](o2.AsIter()(ctx, source)))
		defer stop2()
		var mapping2 []A2

		next3, stop3 := iter.Pull(iter.Seq[ValueIE[I3, A3]](o3.AsIter()(ctx, source)))
		defer stop3()
		var mapping3 []A3

		next4, stop4 := iter.Pull(iter.Seq[ValueIE[I4, A4]](o4.AsIter()(ctx, source)))
		defer stop4()
		var mapping4 []A4

		next5, stop5 := iter.Pull(iter.Seq[ValueIE[I5, A5]](o5.AsIter()(ctx, source)))
		defer stop5()
		var mapping5 []A5

		next6, stop6 := iter.Pull(iter.Seq[ValueIE[I6, A6]](o6.AsIter()(ctx, source)))
		defer stop6()
		var mapping6 []A6

		next7, stop7 := iter.Pull(iter.Seq[ValueIE[I7, A7]](o7.AsIter()(ctx, source)))
		defer stop7()
		var mapping7 []A7

		next8, stop8 := iter.Pull(iter.Seq[ValueIE[I8, A8]](o8.AsIter()(ctx, source)))
		defer stop8()
		var mapping8 []A8

		ret, err := o0.AsModify()(ctx, func(index I0, focus A0) (A0, error) {
			v1, ok := next1()
			if !ok {
				return focus, nil
			}

			i1, a1, err := v1.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v2, ok := next2()
			if !ok {
				return focus, nil
			}

			i2, a2, err := v2.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v3, ok := next3()
			if !ok {
				return focus, nil
			}

			i3, a3, err := v3.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v4, ok := next4()
			if !ok {
				return focus, nil
			}

			i4, a4, err := v4.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v5, ok := next5()
			if !ok {
				return focus, nil
			}

			i5, a5, err := v5.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v6, ok := next6()
			if !ok {
				return focus, nil
			}

			i6, a6, err := v6.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v7, ok := next7()
			if !ok {
				return focus, nil
			}

			i7, a7, err := v7.Get()
			if err != nil {
				var a A0

				return a, err
			}

			v8, ok := next8()
			if !ok {
				return focus, nil
			}

			i8, a8, err := v8.Get()
			if err != nil {
				var a A0

				return a, err
			}

			mapped, err := fmap(lo.T9(index, i1, i2, i3, i4, i5, i6, i7, i8), lo.T9(focus, a1, a2, a3, a4, a5, a6, a7, a8))
			if err != nil {
				var a A0

				return a, err
			}

			mapping1 = append(mapping1, mapped.B)
			mapping2 = append(mapping2, mapped.C)
			mapping3 = append(mapping3, mapped.D)
			mapping4 = append(mapping4, mapped.E)
			mapping5 = append(mapping5, mapped.F)
			mapping6 = append(mapping6, mapped.G)
			mapping7 = append(mapping7, mapped.H)
			mapping8 = append(mapping8, mapped.I)
			return mapped.A, err
		}, source)
		if err != nil {
			var s S

			return s, err
		}

		i := 0
		ret, err = o1.AsModify()(ctx, func(index I1, focus A1) (A1, error) {
			if i >= len(mapping1) {
				return focus, nil
			}

			return mapping1[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o2.AsModify()(ctx, func(index I2, focus A2) (A2, error) {
			if i >= len(mapping2) {
				return focus, nil
			}

			return mapping2[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o3.AsModify()(ctx, func(index I3, focus A3) (A3, error) {
			if i >= len(mapping3) {
				return focus, nil
			}

			return mapping3[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o4.AsModify()(ctx, func(index I4, focus A4) (A4, error) {
			if i >= len(mapping4) {
				return focus, nil
			}

			return mapping4[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o5.AsModify()(ctx, func(index I5, focus A5) (A5, error) {
			if i >= len(mapping5) {
				return focus, nil
			}

			return mapping5[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o6.AsModify()(ctx, func(index I6, focus A6) (A6, error) {
			if i >= len(mapping6) {
				return focus, nil
			}

			return mapping6[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o7.AsModify()(ctx, func(index I7, focus A7) (A7, error) {
			if i >= len(mapping7) {
				return focus, nil
			}

			return mapping7[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		ret, err = o8.AsModify()(ctx, func(index I8, focus A8) (A8, error) {
			if i >= len(mapping8) {
				return focus, nil
			}

			return mapping8[i], nil
		}, ret)
		if err != nil {
			var s S

			return s, err
		}

		return ret, nil
	}, nil, func(a lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8], b lo.Tuple9[I0, I1, I2, I3, I4, I5, I6, I7, I8]) bool {
		if o0.AsIxMatch()(a.A, b.A) != true {
			return false
		}

		if o1.AsIxMatch()(a.B, b.B) != true {
			return false
		}

		if o2.AsIxMatch()(a.C, b.C) != true {
			return false
		}

		if o3.AsIxMatch()(a.D, b.D) != true {
			return false
		}

		if o4.AsIxMatch()(a.E, b.E) != true {
			return false
		}

		if o5.AsIxMatch()(a.F, b.F) != true {
			return false
		}

		if o6.AsIxMatch()(a.G, b.G) != true {
			return false
		}

		if o7.AsIxMatch()(a.H, b.H) != true {
			return false
		}

		if o8.AsIxMatch()(a.I, b.I) != true {
			return false
		}

		return true
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.TupleOf{
			OpticTypeExpr: ot,
			Elements:      []expr.OpticExpression{o0.AsExpr(), o1.AsExpr(), o2.AsExpr(), o3.AsExpr(), o4.AsExpr(), o5.AsExpr(), o6.AsExpr(), o7.AsExpr(), o8.AsExpr()},
		}
	}, o0, o1, o2, o3, o4, o5, o6, o7, o8))
}

// The DelveT2 combinator focuses each element of a [lo.Tuple2] using the given [Optic]
func DelveT2[I any, S any, A any, RET any, RW any, DIR any, ERR any](o Optic[I, S, S, A, A, RET, RW, DIR, ERR]) Optic[lo.Tuple2[I, I], lo.Tuple2[S, S], lo.Tuple2[S, S], lo.Tuple2[A, A], lo.Tuple2[A, A], RET, RW, UniDir, ERR] {
	return UnsafeReconstrain[RET, RW, UniDir, ERR](T2Of(Compose(T2A[S, S](), o), Compose(T2B[S, S](), o)))
}

// The DelveT3 combinator focuses each element of a [lo.Tuple3] using the given [Optic]
func DelveT3[I any, S any, A any, RET any, RW any, DIR any, ERR any](o Optic[I, S, S, A, A, RET, RW, DIR, ERR]) Optic[lo.Tuple3[I, I, I], lo.Tuple3[S, S, S], lo.Tuple3[S, S, S], lo.Tuple3[A, A, A], lo.Tuple3[A, A, A], RET, RW, UniDir, ERR] {
	return UnsafeReconstrain[RET, RW, UniDir, ERR](T3Of(Compose(T3A[S, S, S](), o), Compose(T3B[S, S, S](), o), Compose(T3C[S, S, S](), o)))
}

// The DelveT4 combinator focuses each element of a [lo.Tuple4] using the given [Optic]
func DelveT4[I any, S any, A any, RET any, RW any, DIR any, ERR any](o Optic[I, S, S, A, A, RET, RW, DIR, ERR]) Optic[lo.Tuple4[I, I, I, I], lo.Tuple4[S, S, S, S], lo.Tuple4[S, S, S, S], lo.Tuple4[A, A, A, A], lo.Tuple4[A, A, A, A], RET, RW, UniDir, ERR] {
	return UnsafeReconstrain[RET, RW, UniDir, ERR](T4Of(Compose(T4A[S, S, S, S](), o), Compose(T4B[S, S, S, S](), o), Compose(T4C[S, S, S, S](), o), Compose(T4D[S, S, S, S](), o)))
}

// The DelveT5 combinator focuses each element of a [lo.Tuple5] using the given [Optic]
func DelveT5[I any, S any, A any, RET any, RW any, DIR any, ERR any](o Optic[I, S, S, A, A, RET, RW, DIR, ERR]) Optic[lo.Tuple5[I, I, I, I, I], lo.Tuple5[S, S, S, S, S], lo.Tuple5[S, S, S, S, S], lo.Tuple5[A, A, A, A, A], lo.Tuple5[A, A, A, A, A], RET, RW, UniDir, ERR] {
	return UnsafeReconstrain[RET, RW, UniDir, ERR](T5Of(Compose(T5A[S, S, S, S, S](), o), Compose(T5B[S, S, S, S, S](), o), Compose(T5C[S, S, S, S, S](), o), Compose(T5D[S, S, S, S, S](), o), Compose(T5E[S, S, S, S, S](), o)))
}

// The DelveT6 combinator focuses each element of a [lo.Tuple6] using the given [Optic]
func DelveT6[I any, S any, A any, RET any, RW any, DIR any, ERR any](o Optic[I, S, S, A, A, RET, RW, DIR, ERR]) Optic[lo.Tuple6[I, I, I, I, I, I], lo.Tuple6[S, S, S, S, S, S], lo.Tuple6[S, S, S, S, S, S], lo.Tuple6[A, A, A, A, A, A], lo.Tuple6[A, A, A, A, A, A], RET, RW, UniDir, ERR] {
	return UnsafeReconstrain[RET, RW, UniDir, ERR](T6Of(Compose(T6A[S, S, S, S, S, S](), o), Compose(T6B[S, S, S, S, S, S](), o), Compose(T6C[S, S, S, S, S, S](), o), Compose(T6D[S, S, S, S, S, S](), o), Compose(T6E[S, S, S, S, S, S](), o), Compose(T6F[S, S, S, S, S, S](), o)))
}

// The DelveT7 combinator focuses each element of a [lo.Tuple7] using the given [Optic]
func DelveT7[I any, S any, A any, RET any, RW any, DIR any, ERR any](o Optic[I, S, S, A, A, RET, RW, DIR, ERR]) Optic[lo.Tuple7[I, I, I, I, I, I, I], lo.Tuple7[S, S, S, S, S, S, S], lo.Tuple7[S, S, S, S, S, S, S], lo.Tuple7[A, A, A, A, A, A, A], lo.Tuple7[A, A, A, A, A, A, A], RET, RW, UniDir, ERR] {
	return UnsafeReconstrain[RET, RW, UniDir, ERR](T7Of(Compose(T7A[S, S, S, S, S, S, S](), o), Compose(T7B[S, S, S, S, S, S, S](), o), Compose(T7C[S, S, S, S, S, S, S](), o), Compose(T7D[S, S, S, S, S, S, S](), o), Compose(T7E[S, S, S, S, S, S, S](), o), Compose(T7F[S, S, S, S, S, S, S](), o), Compose(T7G[S, S, S, S, S, S, S](), o)))
}

// The DelveT8 combinator focuses each element of a [lo.Tuple8] using the given [Optic]
func DelveT8[I any, S any, A any, RET any, RW any, DIR any, ERR any](o Optic[I, S, S, A, A, RET, RW, DIR, ERR]) Optic[lo.Tuple8[I, I, I, I, I, I, I, I], lo.Tuple8[S, S, S, S, S, S, S, S], lo.Tuple8[S, S, S, S, S, S, S, S], lo.Tuple8[A, A, A, A, A, A, A, A], lo.Tuple8[A, A, A, A, A, A, A, A], RET, RW, UniDir, ERR] {
	return UnsafeReconstrain[RET, RW, UniDir, ERR](T8Of(Compose(T8A[S, S, S, S, S, S, S, S](), o), Compose(T8B[S, S, S, S, S, S, S, S](), o), Compose(T8C[S, S, S, S, S, S, S, S](), o), Compose(T8D[S, S, S, S, S, S, S, S](), o), Compose(T8E[S, S, S, S, S, S, S, S](), o), Compose(T8F[S, S, S, S, S, S, S, S](), o), Compose(T8G[S, S, S, S, S, S, S, S](), o), Compose(T8H[S, S, S, S, S, S, S, S](), o)))
}

// The DelveT9 combinator focuses each element of a [lo.Tuple9] using the given [Optic]
func DelveT9[I any, S any, A any, RET any, RW any, DIR any, ERR any](o Optic[I, S, S, A, A, RET, RW, DIR, ERR]) Optic[lo.Tuple9[I, I, I, I, I, I, I, I, I], lo.Tuple9[S, S, S, S, S, S, S, S, S], lo.Tuple9[S, S, S, S, S, S, S, S, S], lo.Tuple9[A, A, A, A, A, A, A, A, A], lo.Tuple9[A, A, A, A, A, A, A, A, A], RET, RW, UniDir, ERR] {
	return UnsafeReconstrain[RET, RW, UniDir, ERR](T9Of(Compose(T9A[S, S, S, S, S, S, S, S, S](), o), Compose(T9B[S, S, S, S, S, S, S, S, S](), o), Compose(T9C[S, S, S, S, S, S, S, S, S](), o), Compose(T9D[S, S, S, S, S, S, S, S, S](), o), Compose(T9E[S, S, S, S, S, S, S, S, S](), o), Compose(T9F[S, S, S, S, S, S, S, S, S](), o), Compose(T9G[S, S, S, S, S, S, S, S, S](), o), Compose(T9H[S, S, S, S, S, S, S, S, S](), o), Compose(T9I[S, S, S, S, S, S, S, S, S](), o)))
}
