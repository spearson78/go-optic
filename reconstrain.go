package optic

// !!DANGER!!
// This function blindly applies the new constraints without regard to the current constraints.
// The resulting optic may return errors when an unimplemented optic function is called.
func UnsafeReconstrain[RET, RW, DIR, ERR, I, S, T, A, B, OLDRET, OLDRW, OLDDIR, OLDERR any](o Optic[I, S, T, A, B, OLDRET, OLDRW, OLDDIR, OLDERR]) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {

	//I considered checking the constraints and panicking here.
	//There are some complex composed combinators that use the type safe reconstraint.
	//If I replaced this with a panicky Reconstrain then a change in any of the optics would cause a panic.
	//THis panic would even depend on the optic passed into the combinator.
	return UnsafeOmni[I, S, T, A, B, RET, RW, DIR, ERR](
		o.AsGetter(),
		o.AsSetter(),

		o.AsIter(),
		o.AsLengthGetter(),
		o.AsModify(),

		o.AsIxGetter(),
		o.AsIxMatch(),

		o.AsReverseGetter(),

		o.AsExprHandler(),
		o.AsExpr,
	)
}
