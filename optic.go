// Package optic implements functional optics for querying and manipulating immutable data structures.
//
// Optics are formed by mixing the following 3 properties
//   - ReturnType, [ReturnOne] or [ReturnMany] , where ReturnMany means 0 or more values.
//   - ReadWriteType, [ReadWrite] or [ReadOnly]
//   - DirectionType, [UniDir] or [BiDir]
//
// The following optics are provided
//   - [Operator] single value getter, [ReturnOne] [ReadOnly] [UniDir]
//   - [Lens] single value getter and setter, [ReturnOne] [ReadOnly] [UniDir]
//   - [Iteration] multi value iteration. [ReturnMany] [ReadOnly] [UniDir]
//   - [Traversal] multi value iteration and modification [ReturnMany] [ReadWrite] [UniDir]
//   - [Traversal] multi value iteration and modification with optimized retrieval by index lookup.  [ReturnMany] [ReadWrite] [UniDir]
//   - [Iso] Bidirectional conversion between equivalent (isomorphic) types. [ReturnOne] [ReadWrite] [BiDir]
//   - [Prism] Bidirectional conversion between possibly incompatible types. [ReturnMany] [ReadWrite] [BiDir]
//
// Each optic is available in 4 varieties. The basic constructors listed above create optics that are
//   - Pure i.e. not able to return errors
//   - Non indexed i.e. no index values are provided in the results
//   - Non polymorphic i.e. the input types are identical to the output types)
//
// Additional constructors are suffixed with
//   - R for impure,non indexed and non polymorphic
//   - I for pure, indexed and non polymorphic
//   - P for pure, non indexed and polymorphic
//   - F for impure, indexed and polymorphic
//
// Optics typically have 8 type paramaters. I,S,T,A,B,RET,RW,DIR
//
// I is the index type. optics focusing on elements of a slice have an int index. optics focusing on elements of a map have the same index type as the map key. optics focusing on elements without a natural index type use the [Void] type
//
// S is the source type. This is the type of the value that will be passed as input to an optic.
//
// T is the result type. This is the type of the value returned by modification actions. For non polymorphic optics the return type is the same as the source type S
//
// A is the focus type. This is the type of value focused by the optic.
//
// B is the result focus type. This is the type of the value that will be embedded in the result type.  For non polymorphic optics the return type is the same as the focus type A
//
// RET is the return type. This indicates if the optic returns a single value like a lens [ReturnOne] or multiple values like a traversal [ReturnMany]
//
// RW is the read write type. This indicates if the optic can be written to [ReadWrite] or only read from [ReadOnly]
//
// DIR is the direction type. This indicates if the optic can be reversed [BiDir] or is one way [UniDir]
//
// Optics can be combined using the provided combinators like
//   - [Compose] which chains optics together.
//   - [Filtered] which filters the values focused by an optic.
//
// Optics can be used to both query and update a data structure using the provided actions (e.g. [ToListOf],[Modify]). Modification actions do not alter the original data structure instead a new modified copy is returned.
// Actions are available in 4 varieties
//   - standard: The action has a simple verb as a name and is able to return errors. e.g. [Modify]
//   - indexed: The action name is prefixed with an I and is able to return errors. e.g. [ModifyI]
//   - must: The action name is prefixed with Must and will panic on error. e.g. [MustModify]
//   - must indexed: The action name is prefixed with MustI and will panic on error. e.g. [MustModifyI]
//
// The Optic package also enables external packages to implement efficient combinators and actions by provided access to the low level optic function definitions via the [Optic] interface. The [Optic.AsExpr] method provides access to the internal expression tree representation of the optic. This can be used to transform optics to queries for external systems e.g. Sql.
package optic

import (
	"context"
	"errors"

	"github.com/spearson78/go-optic/expr"
)

var DisableOpticTypeOptimizations = false

func isType(ot expr.OpticType, t expr.OpticType) bool {
	return (ot & t) == t
}

// View function for [Iteration] type optics.
type IterFunc[I, S, A any] func(ctx context.Context, source S) SeqIE[I, A]

// Modification function for [Traversal] type optics.
type ModifyFunc[I, S, T, A, B any] func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error)

// View function for [Operation] type optics.
type GetterFunc[I, S, A any] func(ctx context.Context, source S) (I, A, error)

// Modification function for [Lens] type optics.
type SetterFunc[S, T, B any] func(ctx context.Context, focus B, source S) (T, error)

// Reverse conversion function for [Iso] and [Prism] type optics
type ReverseGetterFunc[T, B any] func(ctx context.Context, focus B) (T, error)

// View by index function for [Traversal] type optics
type IxGetterFunc[I, S, A any] func(ctx context.Context, index I, source S) SeqIE[I, A]

type IxMatchFunc[I any] func(indexA I, indexB I) bool

type OpGetFunc[S, A any] func(ctx context.Context, source S) (A, error)

// Length getter function for [Iteration] type optics
type LengthGetterFunc[S any] func(ctx context.Context, source S) (int, error)

func maybeUpgradeLengthGetter[S any](lengthGetter func(source S) int) LengthGetterFunc[S] {
	if lengthGetter != nil {
		return func(ctx context.Context, source S) (int, error) {
			return lengthGetter(source), nil
		}
	}

	return nil
}

type OpticRO[I, S, A, RET, RW, DIR, ERR any] interface {
	// Provides this optics [Iteration] style view function.
	// This methods behavior is defined for all optic types
	AsIter() IterFunc[I, S, A]

	// Provides this optics [Operation] style view function.
	// This methods behavior is only defined for [ReturnOne] optics.
	AsGetter() GetterFunc[I, S, A]

	// Provides this optics [Traversal] style view by index function.
	// This methods behavior is only defined for all optics.
	AsIxGetter() IxGetterFunc[I, S, A]

	AsIxMatch() IxMatchFunc[I]

	// Provides this optics [Iteration] style length getter function.
	// This methods behavior is defined for all optic types
	AsLengthGetter() LengthGetterFunc[S]

	AsOpGet() OpGetFunc[S, A]

	// Provides this optics expression tree representation.
	AsExpr() expr.OpticExpression

	//Provides the custom expression handler. Used by non go backends to execute expressions
	AsExprHandler() func(ctx context.Context) (ExprHandler, error)

	// Provides this optics internal type. This indicates which methods have efficient implementations.
	OpticType() expr.OpticType

	// Indicates whether this optic returns exactly 1 result [ReturnOne] or 0 or more [ReturnMany]
	// or was composed of several optics [CompositionTree]. Used to prevent compilation if a [ReturnMany]
	// optic is passed to an action that requires a [ReturnOne] result. e.g. [View].
	ReturnType() RET

	// Indicates whether this optic supports modification [ReadWrite] or can only be viewed [ReadOnly]
	// or was composed of several optics [CompositionTree]. Used to prevent compilation if a [ReadOnly]
	// optic is passed to a modification action e.g. [Over].
	ReadWriteType() RW

	// Indicates whether this optic supports ReverseGet [BiDir] or is unidirectional [UniDir]
	// or was composed of several optics [CompositionTree]. Used to prevent compilation if a [UniDir]
	// optic is passed to an action requiring a reverseget e.g. [ReView].
	DirType() DIR

	ErrType() ERR
}

// Optic interface implemented by all optic types. Provides access to the low level optic implementation functions. These functions are not intended for general user use but are provided to enable efficient combinators and actions to be implemented in external packages.
type Optic[I, S, T, A, B, RET, RW, DIR, ERR any] interface {
	// Provides this optics [Traversal] style modification function.
	// This methods behavior is only defined for [ReadWrite] optics.
	AsModify() ModifyFunc[I, S, T, A, B]

	// Provides this optics [Iteration] style view function.
	// This methods behavior is defined for all optic types
	AsIter() IterFunc[I, S, A]

	// Provides this optics [Operation] style view function.
	// This methods behavior is only defined for [ReturnOne] optics.
	AsGetter() GetterFunc[I, S, A]

	// Provides this optics [Lens] style modification function.
	// This methods behavior is only defined for [ReadWrite] optics.
	AsSetter() SetterFunc[S, T, B]

	// Provides this optics [Traversal] style view by index function.
	// This methods behavior is only defined for all optics.
	AsIxGetter() IxGetterFunc[I, S, A]

	AsIxMatch() IxMatchFunc[I]

	// Provides this optics [Iso] or [Prism] style modification function.
	// This methods behavior is only defined for all [BiDir]
	AsReverseGetter() ReverseGetterFunc[T, B]

	// Provides this optics [Iteration] style length getter function.
	// This methods behavior is defined for all optic types
	AsLengthGetter() LengthGetterFunc[S]

	AsOpGet() OpGetFunc[S, A]

	// Provides this optics expression tree representation.
	AsExpr() expr.OpticExpression

	//Provides the custom expression handler. Used by non go backends to execute expressions
	AsExprHandler() func(ctx context.Context) (ExprHandler, error)

	// Provides this optics internal type. This indicates which methods have efficient implementations.
	OpticType() expr.OpticType

	// Indicates whether this optic returns exactly 1 result [ReturnOne] or 0 or more [ReturnMany]
	// or was composed of several optics [CompositionTree]. Used to prevent compilation if a [ReturnMany]
	// optic is passed to an action that requires a [ReturnOne] result. e.g. [View].
	ReturnType() RET

	// Indicates whether this optic supports modification [ReadWrite] or can only be viewed [ReadOnly]
	// or was composed of several optics [CompositionTree]. Used to prevent compilation if a [ReadOnly]
	// optic is passed to a modification action e.g. [Over].
	ReadWriteType() RW

	// Indicates whether this optic supports ReverseGet [BiDir] or is unidirectional [UniDir]
	// or was composed of several optics [CompositionTree]. Used to prevent compilation if a [UniDir]
	// optic is passed to an action requiring a reverseget e.g. [ReView].
	DirType() DIR

	ErrType() ERR
}

// Indicates this optic returns exactly one result. Used to prevent compilation when passing a [ReturnMany] optic to a [ReturnOne] action e.g. [View]
type ReturnOne Void
type TReturnOne interface {
	comparable
}

// Indicates this optic returns exactly zero or more results.
type ReturnMany []Void

// Indicates this optic is read only. Used to prevent compilation when passing a [ReadOnly] optic to a [ReadWrite] action e.g. [Modify]
type ReadOnly []Void

// Indicates this optic supports read and write actions.
type ReadWrite Void
type TReadWrite interface {
	comparable
}

// Indicates this optic supports reverseget
type BiDir Void
type TBiDir interface {
	comparable
}

// Indicates this optic does not support reverseget
type UniDir []Void

type Pure Void
type TPure interface {
	comparable
}

type Err []Void

var UnsupportedOpticMethod = errors.New("unsupported optic method")

func unsupportedReverseGetter[B, T any](ctx context.Context, focus B) (T, error) {
	var t T
	return t, UnsupportedOpticMethod
}

// Return type that indicates this optic is composed of other optics with potentially different return types. If any part of the composition tree is [ReturnMany] this optic would be considered [ReturnMany] any will fail compilation if passed to a [ReturnOne] action.
//
// For some apis (e.g. [CoalesceN]) it may be necessary to condense the CompositionTree back to a [ReturnOne] or [ReturnMany] using the [Ret1] or [RetM] functions. When implementing custom combinators it is considered best practice to condense the return type.
type CompositionTree[L any, R any] struct {
	L L
	R R
}

type abortModifyErr struct {
	token *error
}

func handleAbortModify(token *error) {
	if r := recover(); r != nil {
		if abortModifyErr, ok := r.(*abortModifyErr); ok && abortModifyErr.token == token {
			//Ignore and let the error be returned
		} else {
			panic(r)
		}
	}
}

func abortModify(token *error) {
	panic(&abortModifyErr{
		token: token,
	})
}

func abortModifyError(err error, result *error) {
	if err != nil {
		*result = err
		abortModify(result)
	}
}
