// Package expr contains the expression types for all built in optics, combinators and functions.
package expr

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/spearson78/go-optic/internal/util"
	"golang.org/x/text/encoding"
)

type Identity struct {
	Index any
	OpticTypeExpr
}

func (e Identity) Short() string {
	return "Identity"
}

func (e Identity) String() string {
	return fmt.Sprintf("Identity(%v)", util.FullTypeName(e.opticS))
}

type Const struct {
	OpticTypeExpr
	Index any
	Value any
}

func (e Const) Short() string {
	return fmt.Sprintf("Const(%v)", e.Value)
}

func (e Const) String() string {
	return fmt.Sprintf("Const[%v](%v)", util.FullTypeName(e.opticS), e.Value)
}

type Value struct {
	OpticTypeExpr
	Value any
}

func (e Value) Short() string {
	return fmt.Sprintf("Value(%v)", e.Value)
}

func (e Value) String() string {
	return fmt.Sprintf("Value[%v](%v)", util.FullTypeName(e.opticS), e.Value)
}

type IgnoreWrite struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e IgnoreWrite) Short() string {
	return fmt.Sprintf("IgnoreWrite(%v)", e.Optic)
}

func (e IgnoreWrite) String() string {
	return fmt.Sprintf("IgnoreWrite[%v](%v)", util.FullTypeName(e.opticS), e.Optic)
}

type TupleElement struct {
	OpticTypeExpr
	Index int
}

func (e TupleElement) Short() string {
	return fmt.Sprintf("TupleElement(%v)", e.Index)
}

func (e TupleElement) String() string {
	return fmt.Sprintf("TupleElement(%v)", e.Index)
}

type TupleOf struct {
	OpticTypeExpr
	Elements []OpticExpression
}

func (e TupleOf) Short() string {
	var sb strings.Builder

	sb.WriteString("TupleOf(")

	for i, v := range e.Elements {
		if i != 0 {
			sb.WriteString(" , ")
		}
		sb.WriteString(v.Short())
	}

	sb.WriteString(")")

	return sb.String()
}

func (e TupleOf) String() string {
	var sb strings.Builder

	sb.WriteString("TupleOf(")

	for i, v := range e.Elements {
		if i != 0 {
			sb.WriteString(" , ")
		}
		sb.WriteString(v.String())
	}

	sb.WriteString(")")

	return sb.String()
}

type TupleDup struct {
	OpticTypeExpr
	N int
}

func (e TupleDup) Short() string {
	return fmt.Sprintf("TupleDup(%v)", e.N)
}

func (e TupleDup) String() string {
	return fmt.Sprintf("TupleDup(%v,%v)", util.FullTypeName(e.opticS), e.N)
}

type SwapTuple struct {
	OpticTypeExpr
	A reflect.Type
	B reflect.Type
}

func (e SwapTuple) Short() string {
	return "SwapTuple"
}

func (e SwapTuple) String() string {
	return fmt.Sprintf("SwapTuple(%v,%v)", util.FullTypeName(e.A), util.FullTypeName(e.B))
}

type SplitString struct {
	OpticTypeExpr
	SplitOn *regexp.Regexp
}

func (e SplitString) Short() string {
	return fmt.Sprintf("StringSplitRegexp(%v)", e.SplitOn)
}

func (e SplitString) String() string {
	return fmt.Sprintf("StringSplitRegexp(%v)", e.SplitOn)
}

type MatchString struct {
	OpticTypeExpr
	Match *regexp.Regexp
}

func (e MatchString) Short() string {
	return fmt.Sprintf("MatchRegexp(%v)", e.Match)
}

func (e MatchString) String() string {
	return fmt.Sprintf("MatchRegexp(%v)", e.Match)
}

type CaptureString struct {
	OpticTypeExpr
	MatchOn *regexp.Regexp
}

func (e CaptureString) Short() string {
	return fmt.Sprintf("CaptureString(%v)", e.MatchOn)
}

func (e CaptureString) String() string {
	return fmt.Sprintf("CaptureString(%v)", e.MatchOn)
}

type Reverse struct {
	OpticTypeExpr
	I reflect.Type
	A reflect.Type
	B reflect.Type
}

func (e Reverse) Short() string {
	return "Reverse"
}

func (e Reverse) String() string {
	return fmt.Sprintf("Reverse(%v,%v,%v)", util.FullTypeName(e.opticI), util.FullTypeName(e.opticA), util.FullTypeName(e.opticB))
}

type Ordered struct {
	OpticTypeExpr

	Optic   OpticExpression
	OrderBy OpticExpression
}

func (e Ordered) Short() string {
	return fmt.Sprintf("Ordered(%v,%v)", e.Optic.Short(), e.OrderBy.Short())
}

func (e Ordered) String() string {
	return fmt.Sprintf("Ordered(%v,%v)", e.Optic.String(), e.OrderBy.String())
}

type OrderBy struct {
	OpticTypeExpr

	Optic OpticExpression
}

func (e OrderBy) Short() string {
	return fmt.Sprintf("OrderBy(%v)", e.Optic.Short())
}

func (e OrderBy) String() string {
	return fmt.Sprintf("OrderBy(%v)", e.Optic.String())
}

type Desc struct {
	OpticTypeExpr

	Optic OpticExpression
}

func (e Desc) Short() string {
	return fmt.Sprintf("Desc(%v)", e.Optic.Short())
}

func (e Desc) String() string {
	return fmt.Sprintf("Desc(%v)", e.Optic.String())
}

type OrderByN struct {
	OpticTypeExpr

	Optics []OpticExpression
}

func (e OrderByN) Short() string {
	var sb strings.Builder

	sb.WriteString("OrderByN(")

	for i, v := range e.Optics {
		if i != 0 {
			sb.WriteString(" , ")
		}
		sb.WriteString(v.Short())
	}

	sb.WriteString(")")

	return sb.String()
}

func (e OrderByN) String() string {
	var sb strings.Builder

	sb.WriteString("OrderByN(")

	for i, v := range e.Optics {
		if i != 0 {
			sb.WriteString(" , ")
		}
		sb.WriteString(v.String())
	}

	sb.WriteString(")")

	return sb.String()
}

type None struct {
	OpticTypeExpr
}

func (e None) Short() string {
	return "None"
}

func (e None) String() string {
	return fmt.Sprintf("None(%v)", util.FullTypeName(e.opticA))
}

type Some struct {
	OpticTypeExpr
}

func (e Some) Short() string {
	return "Some"
}

func (e Some) String() string {
	return fmt.Sprintf("Some(%v)", util.FullTypeName(e.opticA))
}

type PtrOption struct {
	OpticTypeExpr
	A reflect.Type
	B reflect.Type
}

func (e PtrOption) Short() string {
	return "PtrOption"
}

func (e PtrOption) String() string {
	return fmt.Sprintf("PtrOption(%v)", util.FullTypeName(e.A))
}

type OptionOfFirst struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e OptionOfFirst) Short() string {
	return fmt.Sprintf("OptionOfFirst(%v)", e.Optic.Short())
}

func (e OptionOfFirst) String() string {
	return fmt.Sprintf("OptionOfFirst(%v)", e.Optic.String())
}

type Traverse struct {
	OpticTypeExpr
}

func (e Traverse) Short() string {
	return "Traverse"
}

func (e Traverse) String() string {
	return fmt.Sprintf("Traverse(%v,%v,%v,%v,%v)", util.FullTypeName(e.opticI), util.FullTypeName(e.opticS), util.FullTypeName(e.opticT), util.FullTypeName(e.opticA), util.FullTypeName(e.opticB))
}

type FieldLens struct {
	OpticTypeExpr
	Field    reflect.StructField
	FieldNum int
}

func (e FieldLens) Short() string {
	return fmt.Sprintf("%v.%v", e.opticS.Name(), e.Field.Name)
}

func (e FieldLens) String() string {
	return fmt.Sprintf("FieldLens(%v.%v)", util.FullTypeName(e.opticS), e.Field.Name)
}

type MethodGetter struct {
	OpticTypeExpr
	MethodName string
	Method     reflect.Type
}

func (e MethodGetter) Short() string {
	return fmt.Sprintf("MethodGetter(%v)", e.MethodName)
}

func (e MethodGetter) String() string {
	return fmt.Sprintf("MethodGetter(%v.%v)", util.FullTypeName(e.opticS), e.MethodName)
}

type MethodLens struct {
	OpticTypeExpr
	GetterName string
	Getter     reflect.Type
	SetterName string
	Setter     reflect.Type
}

func (e MethodLens) Short() string {
	return fmt.Sprintf("MethodLens(%v,%v)", e.GetterName, e.SetterName)
}

func (e MethodLens) String() string {
	return fmt.Sprintf("MethodLens(%v.%v,%v)", util.FullTypeName(e.opticS), e.GetterName, e.SetterName)
}

type Chosen struct {
	OpticTypeExpr
}

func (e Chosen) Short() string {
	return "Chosen"
}

func (e Chosen) String() string {
	return fmt.Sprintf("Chosen(%v)", util.FullTypeName(e.opticA))
}

type Left struct {
	OpticTypeExpr
	B reflect.Type
}

func (e Left) Short() string {
	return "Left"
}

func (e Left) String() string {
	return fmt.Sprintf("Left(%v,%v)", util.FullTypeName(e.opticA), util.FullTypeName(e.B))
}

type Right struct {
	OpticTypeExpr
	A reflect.Type
}

func (e Right) Short() string {
	return "Right"
}

func (e Right) String() string {
	return fmt.Sprintf("Right(%v,%v)", util.FullTypeName(e.A), util.FullTypeName(e.opticB))
}

type IxMapperType int

const (
	IxMapperCustom IxMapperType = iota
	IxMapperLeft
	IxMapperRight
	IxMapperBoth
)

func (i IxMapperType) String() string {
	switch i {
	case IxMapperCustom:
		return "custom"
	case IxMapperLeft:
		return "left"
	case IxMapperRight:
		return "right"
	case IxMapperBoth:
		return "both"
	default:
		return "unknown"
	}
}

type IxMap struct {
	OpticTypeExpr

	Type IxMapperType
}

func (e IxMap) Short() string {
	return "IxMap " + e.Type.String()
}

func (e IxMap) String() string {
	return fmt.Sprintf("IxMap(%v)", e.Type)
}

type Compose struct {
	OpticTypeExpr

	Left  OpticExpression
	Right OpticExpression
	IxMap OpticExpression
}

func (e Compose) Short() string {
	return fmt.Sprintf("%v | %v", e.Left.Short(), e.Right.Short())
}

func (e Compose) String() string {
	return fmt.Sprintf("%v | %v", e.Left.String(), e.Right.String())
}

type Length struct {
	OpticTypeExpr

	Optic OpticExpression
}

func (e Length) Short() string {
	return fmt.Sprintf("Length(%v)", e.Optic.Short())
}

func (e Length) String() string {
	return fmt.Sprintf("Length(%v)", e.Optic.String())
}

type Reversed struct {
	OpticTypeExpr

	Optic OpticExpression
}

func (e Reversed) Short() string {
	return "BacReversedkwards"
}

func (e Reversed) String() string {
	return fmt.Sprintf("Reversed(%v)", e.Optic.String())
}

type StringOf struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e StringOf) Short() string {
	return fmt.Sprintf("StringOf(%v)", e.Optic.Short())
}

func (e StringOf) String() string {
	return fmt.Sprintf("StringOf(%v)", e.Optic.String())
}

type CollectionOf struct {
	OpticTypeExpr
	I     reflect.Type
	A     reflect.Type
	B     reflect.Type
	Optic OpticExpression
}

func (e CollectionOf) Short() string {
	return fmt.Sprintf("SeqOf(%v)", e.Optic.Short())
}

func (e CollectionOf) String() string {
	return fmt.Sprintf("SeqOf(%v)", e.Optic.String())
}

type SeqOf struct {
	OpticTypeExpr
	I     reflect.Type
	A     reflect.Type
	B     reflect.Type
	Optic OpticExpression
}

func (e SeqOf) Short() string {
	return fmt.Sprintf("SeqOf(%v)", e.Optic.Short())
}

func (e SeqOf) String() string {
	return fmt.Sprintf("SeqOf(%v)", e.Optic.String())
}

type EqT2Of struct {
	OpticTypeExpr
	Optic OpticExpression
	Eq    OpticExpression
}

func (e EqT2Of) Short() string {
	return fmt.Sprintf("SeqEq(%v)", e.Eq.Short())
}

func (e EqT2Of) String() string {
	return fmt.Sprintf("SeqEq(%v)", e.Eq.String())
}

type SliceOf struct {
	OpticTypeExpr
	I     reflect.Type
	A     reflect.Type
	B     reflect.Type
	Optic OpticExpression
}

func (e SliceOf) Short() string {
	return fmt.Sprintf("SliceOf(%v)", e.Optic.Short())
}

func (e SliceOf) String() string {
	return fmt.Sprintf("SliceOf(%v)", e.Optic.String())
}

type MapOfReduced struct {
	OpticTypeExpr

	Optic   OpticExpression
	Reducer ReducerExpression
}

func (e MapOfReduced) Short() string {
	return fmt.Sprintf("MapOfLast(%v)", e.Optic.Short())
}

func (e MapOfReduced) String() string {
	return fmt.Sprintf("MapOfLast(%v)", e.Optic.String())
}

type Make struct {
	OpticTypeExpr
	Size []int
}

func (e Make) Short() string {
	return fmt.Sprintf("Make(%v)", e.Size)
}

func (e Make) String() string {
	return fmt.Sprintf("Make(%v)", e.Size)
}

type AtT2 struct {
	OpticTypeExpr
	V      reflect.Type
	Equals OpticExpression
}

func (e AtT2) Short() string {
	return fmt.Sprintf("AtT2()")
}

func (e AtT2) String() string {
	return fmt.Sprintf("AtT2()")
}

type Coalesce struct {
	OpticTypeExpr

	Optics []OpticExpression
}

func (e Coalesce) Short() string {
	var sb strings.Builder

	sb.WriteString("Coalesce(")

	for i, v := range e.Optics {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(v.Short())
	}

	sb.WriteString(")")

	return sb.String()
}

func (e Coalesce) String() string {
	var sb strings.Builder

	sb.WriteString("Case(")

	for i, v := range e.Optics {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(v.String())
	}

	sb.WriteString(")")

	return sb.String()
}

type Case struct {
	OpticTypeExpr
	Condition OpticExpression
	Default   OpticExpression
}

func (e Case) Short() string {
	return fmt.Sprintf("Case(%v -> %v)", e.Condition.Short(), e.Default.Short())
}

func (e Case) String() string {
	return fmt.Sprintf("Case(%v -> %v)", e.Condition.String(), e.Default.String())
}

type Switch struct {
	OpticTypeExpr

	Whens   []Case
	Default OpticExpression
}

func (e Switch) Short() string {
	var sb strings.Builder

	sb.WriteString("Switch(")

	for i, v := range e.Whens {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(v.Short())
	}

	sb.WriteString("Default(")
	sb.WriteString(e.Default.Short())

	sb.WriteString("))")

	return sb.String()
}

func (e Switch) String() string {
	var sb strings.Builder

	sb.WriteString("Switch(")

	for i, v := range e.Whens {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(v.String())
	}

	sb.WriteString("Default(")
	sb.WriteString(e.Default.String())

	sb.WriteString("))")

	return sb.String()
}

type Non struct {
	OpticTypeExpr
	Default any
	Equal   OpticExpression
}

func (e Non) Short() string {
	return fmt.Sprintf("Non(%v)", e.Default)
}

func (e Non) String() string {
	return fmt.Sprintf("Non(%v)", e.Default)
}

type ReverseGet struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e ReverseGet) Short() string {
	return fmt.Sprintf("ReverseGet(%v)", e.Optic.Short())
}

func (e ReverseGet) String() string {
	return fmt.Sprintf("ReverseGet(%v)", e.Optic.String())
}

type FilterMode byte

func (f FilterMode) String() string {
	switch f {
	case FilterStop:
		return "FilterStop"
	case FilterContinue:
		return "FilterContinue"
	case FilterYieldAll:
		return "FilterYieldAll"
	default:
		return "unknown"
	}
}

const (
	FilterStop FilterMode = iota
	FilterContinue
	FilterYieldAll
)

type Filtered struct {
	OpticTypeExpr

	PosPred     OpticExpression
	Pred        OpticExpression
	Optic       OpticExpression
	NoMatchMode FilterMode
	MatchMode   FilterMode
}

func (e Filtered) Short() string {
	return fmt.Sprintf("Filtered(%v,%v,%v,%v)", e.PosPred.Short(), e.Pred.Short(), e.MatchMode, e.NoMatchMode)
}

func (e Filtered) String() string {
	return fmt.Sprintf("Filtered(%v,%v,%v,%v,%v)", e.Optic.String(), e.PosPred.String(), e.Pred.String(), e.MatchMode, e.NoMatchMode)
}

type Index struct {
	OpticTypeExpr

	Index any
	Optic OpticExpression
}

func (e Index) Short() string {
	return fmt.Sprintf("Index(%v)", e.Index)
}

func (e Index) String() string {
	return fmt.Sprintf("Index(%v,%v)", e.Index, e.Optic.String())
}

type OpIxGet struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e OpIxGet) Short() string {
	return fmt.Sprintf("OpIxGet(%v)", e.Optic.Short())
}

func (e OpIxGet) String() string {
	return fmt.Sprintf("OpIxGet(%v,%v)", util.FullTypeName(e.opticA), e.Optic.String())
}

type Modify struct {
	OpticTypeExpr
	Optic OpticExpression
	Fmap  OpticExpression
}

func (e Modify) Short() string {
	return fmt.Sprintf("Modify(%v,%v)", e.Optic.Short(), e.Fmap.Short())
}

func (e Modify) String() string {
	return fmt.Sprintf("Modify(%v,%v,%v)", util.FullTypeName(e.opticA), e.Optic.String(), e.Fmap.String())
}

type Set struct {
	OpticTypeExpr
	Optic OpticExpression
	Val   any
}

func (e Set) Short() string {
	return fmt.Sprintf("Set(%v,%v)", e.Optic.Short(), e.Val)
}

func (e Set) String() string {
	return fmt.Sprintf("Set(%v,%v,%v)", util.FullTypeName(e.opticA), e.Optic.String(), e.Val)
}

type WithIndex struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e WithIndex) Short() string {
	return fmt.Sprintf("WithIndex(%v)", e.Optic.Short())
}

func (e WithIndex) String() string {
	return fmt.Sprintf("WithIndex(%v)", e.Optic.String())
}

type Indexing struct {
	OpticTypeExpr

	Optic OpticExpression
}

func (e Indexing) Short() string {
	return fmt.Sprintf("Indexing(%v)", e.Optic.Short())
}

func (e Indexing) String() string {
	return fmt.Sprintf("Indexing(%v)", e.Optic.String())
}

type ReIndexed struct {
	OpticTypeExpr

	IxMap   OpticExpression
	IxMatch OpticExpression
	Optic   OpticExpression
}

func (e ReIndexed) Short() string {
	return fmt.Sprintf("ReIndexed(%v,%v)", e.IxMap.Short(), e.Optic.Short())
}

func (e ReIndexed) String() string {
	return fmt.Sprintf("ReIndexed(%v,%v)", e.IxMap.String(), e.Optic.String())
}

type ReIndexedCol struct {
	OpticTypeExpr

	IxMap    OpticExpression
	IxMatchI OpticExpression
	IxMatchJ OpticExpression
}

func (e ReIndexedCol) Short() string {
	return fmt.Sprintf("ReIndexedCol(%v)", e.IxMap.Short())
}

func (e ReIndexedCol) String() string {
	return fmt.Sprintf("ReIndexedCol(%v)", e.IxMap.String())
}

type DiffCol struct {
	OpticTypeExpr

	Threshold float64
	Distance  OpticExpression
	IxMatch   OpticExpression
}

func (e DiffCol) Short() string {
	return fmt.Sprintf("DiffCol(%v,%v)", e.Threshold, e.Distance.Short())
}

func (e DiffCol) String() string {
	return fmt.Sprintf("DiffCol(%v,%v)", e.Threshold, e.Distance.String())
}

type SelfIndex struct {
	OpticTypeExpr

	IxMatch OpticExpression
	Optic   OpticExpression
}

func (e SelfIndex) Short() string {
	return fmt.Sprintf("SelfIndex(%v)", e.Optic.Short())
}

func (e SelfIndex) String() string {
	return fmt.Sprintf("SelfIndex(%v)", e.Optic.String())
}

type Beside struct {
	OpticTypeExpr
	Optic1 OpticExpression
	Optic2 OpticExpression
}

func (e Beside) Short() string {
	return fmt.Sprintf("Beside(%v,%v)", e.Optic1.Short(), e.Optic2.Short())
}

func (e Beside) String() string {
	return fmt.Sprintf("Beside(%v,%v)", e.Optic1.String(), e.Optic2.String())
}

type BesideEither struct {
	OpticTypeExpr
	OnLeft  OpticExpression
	OnRight OpticExpression
}

func (e BesideEither) Short() string {
	return fmt.Sprintf("BesideEither(%v,%v)", e.OnLeft.Short(), e.OnRight.Short())
}

func (e BesideEither) String() string {
	return fmt.Sprintf("BesideEither(%v,%v)", e.OnLeft.String(), e.OnRight.String())
}

type Ignore struct {
	OpticTypeExpr

	Optic     OpticExpression
	Predicate OpticExpression
}

func (e Ignore) Short() string {
	return fmt.Sprintf("Ignore(%v,%v)", e.Optic.Short(), e.Predicate.Short())
}

func (e Ignore) String() string {
	return fmt.Sprintf("Ignore(%v,%v)", e.Optic.String(), e.Predicate.String())
}

type Stop struct {
	OpticTypeExpr

	Optic     OpticExpression
	Predicate OpticExpression
}

func (e Stop) Short() string {
	return fmt.Sprintf("Stop(%v,%v)", e.Optic.Short(), e.Predicate.Short())
}

func (e Stop) String() string {
	return fmt.Sprintf("Stop(%v,%v)", e.Optic.String(), e.Predicate.String())
}

type Catch struct {
	OpticTypeExpr

	Optic  OpticExpression
	CatchA OpticExpression
	CatchT OpticExpression
}

func (e Catch) Short() string {
	return fmt.Sprintf("Catch(%v,%v,%v)", e.Optic.Short(), e.CatchA.Short(), e.CatchT.Short())
}

func (e Catch) String() string {
	return fmt.Sprintf("Catch(%v,%v,%v)", e.Optic.String(), e.CatchA.String(), e.CatchT.String())
}

type Cast struct {
	OpticTypeExpr
}

func (e Cast) Short() string {
	return "Cast"
}

func (e Cast) String() string {
	return fmt.Sprintf("Cast(%v,%v)", util.FullTypeName(e.opticS), util.FullTypeName(e.opticA))
}

type Last struct {
	OpticTypeExpr

	Optic OpticExpression
}

func (e Last) Short() string {
	return fmt.Sprintf("Last(%v)", e.Optic.Short())
}

func (e Last) String() string {
	return fmt.Sprintf("Last(%v)", e.Optic.String())
}

type Lookup struct {
	OpticTypeExpr
	Optic  OpticExpression
	Source any
}

func (e Lookup) Short() string {
	return fmt.Sprintf("Lookup(%v,%v)", e.Optic.Short(), e.Source)
}

func (e Lookup) String() string {
	return fmt.Sprintf("Lookup(%v,%v)", e.Optic.String(), e.Source)
}

type Matching struct {
	OpticTypeExpr
	Match OpticExpression
}

func (e Matching) Short() string {
	return fmt.Sprintf("Matching(%v)", e.Match.Short())
}

func (e Matching) String() string {
	return fmt.Sprintf("Matching(%v)", e.Match.String())
}

type Polymorphic struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e Polymorphic) Short() string {
	return fmt.Sprintf("Polymorphic(%v)", e.Optic.Short())
}

func (e Polymorphic) String() string {
	return fmt.Sprintf("Polymorphic(%v)", e.Optic.String())
}

type EditType byte

const (
	EditInsert EditType = 1 << iota
	EditDelete
	EditSubstitute
	EditTranspose

	EditLevenshtein = EditInsert | EditDelete | EditSubstitute
	EditOSA         = EditLevenshtein | EditTranspose
	EditLCS         = EditInsert | EditDelete

	EditAny = 0xFF
)

type EditDistance struct {
	OpticTypeExpr
	Optic        OpticExpression
	AllowedEdits EditType
}

func (e EditDistance) Short() string {
	return fmt.Sprintf("EditDistance(%v)", e.Optic.Short())
}

func (e EditDistance) String() string {
	return fmt.Sprintf("EditDistance(%v)", e.Optic.String())
}

type ForEach struct {
	OpticTypeExpr
	ForEach OpticExpression
	Op      OpticExpression
}

func (e ForEach) Short() string {
	return fmt.Sprintf("ForEach(%v,%v)", e.ForEach.Short(), e.Op.Short())
}

func (e ForEach) String() string {
	return fmt.Sprintf("ForEach(%v,%v)", e.ForEach.String(), e.Op.String())
}

type MapReduce struct {
	OpticTypeExpr

	Optic   OpticExpression
	Reducer ReducerExpression
	Mapper  OpticExpression
}

func (e MapReduce) Short() string {
	return fmt.Sprintf("MapReduce(%v,%v,%v)", e.Optic.Short(), e.Reducer, e.Mapper.Short())
}

func (e MapReduce) String() string {
	return fmt.Sprintf("MapReduce(%v,%v,%v)", e.Optic.String(), e.Reducer, e.Mapper.String())
}

type ToCol struct {
	OpticTypeExpr
	I reflect.Type
	A reflect.Type
	B reflect.Type
}

func (e ToCol) Short() string {
	return "ToCol"
}

func (e ToCol) String() string {
	return fmt.Sprintf("ToCol(%v,%v,%v,%v,%v)", util.FullTypeName(e.I), util.FullTypeName(e.opticS), util.FullTypeName(e.opticT), util.FullTypeName(e.A), util.FullTypeName(e.B))
}

type EncodeString struct {
	OpticTypeExpr
	Encoding encoding.Encoding
}

func (e EncodeString) Short() string {
	return fmt.Sprintf("EncodeString(%v)", e.Encoding)
}

func (e EncodeString) String() string {
	return fmt.Sprintf("EncodeString(%v)", e.Encoding)
}

type StringHasPrefix struct {
	OpticTypeExpr
	Prefix string
}

func (e StringHasPrefix) Short() string {
	return fmt.Sprintf("StringHasPrefix(%v)", e.Prefix)
}

func (e StringHasPrefix) String() string {
	return fmt.Sprintf("StringHasPrefix(%v)", e.Prefix)
}

type StringHasSuffix struct {
	OpticTypeExpr
	Suffix string
}

func (e StringHasSuffix) Short() string {
	return fmt.Sprintf("StringHasSuffix(%v)", e.Suffix)
}

func (e StringHasSuffix) String() string {
	return fmt.Sprintf("StringHasSuffix(%v)", e.Suffix)
}

type Error struct {
	OpticTypeExpr

	Err error
}

func (e Error) Short() string {
	return fmt.Sprintf("Error(%v)", e.Err)
}

func (e Error) String() string {
	return fmt.Sprintf("Error(%v)", e.Err)
}

type Throw struct {
	OpticTypeExpr
}

func (e Throw) Short() string {
	return "Throw()"
}

func (e Throw) String() string {
	return "Throw()"
}

type ErrorIs struct {
	OpticTypeExpr

	Target error
}

func (e ErrorIs) Short() string {
	return fmt.Sprintf("ErrorIs(%v)", e.Target)
}

func (e ErrorIs) String() string {
	return fmt.Sprintf("ErrorIs(%v)", e.Target)
}

type ErrorAs struct {
	OpticTypeExpr
}

func (e ErrorAs) Short() string {
	return fmt.Sprintf("ErrorAs(%v)", e.opticA)
}

func (e ErrorAs) String() string {
	return fmt.Sprintf("ErrorAs(%v)", e.opticA)
}

type WithPanic struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e WithPanic) Short() string {
	return fmt.Sprintf("WithPanic(%v)", e.Optic.Short())
}

func (e WithPanic) String() string {
	return fmt.Sprintf("WithPanic(%v)", e.Optic.String())
}

type SeqExprValue struct {
	Index    any
	Value    any
	ValuePtr any
	Error    error
}

type SeqExpr func(ctx context.Context, yield func(SeqExprValue) bool)

type AsSeqExpr interface {
	AsExpr() SeqExpr
}

type AppendCol struct {
	OpticTypeExpr
	I reflect.Type
	A reflect.Type
}

func (e AppendCol) Short() string {
	return "Appending()"
}

func (e AppendCol) String() string {
	return "Appending()"
}

type Any struct {
	OpticTypeExpr
	Optic OpticExpression
	Pred  OpticExpression
}

func (e Any) Short() string {
	return fmt.Sprintf("Any(%v,%v)", e.Optic.Short(), e.Pred.Short())
}

func (e Any) String() string {
	return fmt.Sprintf("Any(%v,%v)", e.Optic.String(), e.Pred.String())
}

type All struct {
	OpticTypeExpr
	Optic OpticExpression
	Pred  OpticExpression
}

func (e All) Short() string {
	return fmt.Sprintf("All(%v,%v)", e.Optic.Short(), e.Pred.Short())
}

func (e All) String() string {
	return fmt.Sprintf("All(%v,%v)", e.Optic.String(), e.Pred.String())
}

type SubCol struct {
	OpticTypeExpr
	Start  int
	Length int
}

func (e SubCol) Short() string {
	return fmt.Sprintf("SubCol(%v,%v)", e.Start, e.Length)
}

func (e SubCol) String() string {
	return fmt.Sprintf("SubCol(%v,%v)", e.Start, e.Length)
}

type Empty struct {
	OpticTypeExpr
}

func (e Empty) Short() string {
	return "Empty()"
}

func (e Empty) String() string {
	return "Empty()"
}

type EqDeep struct {
	OpticTypeExpr
}

func (e EqDeep) Short() string {
	return "EqDeep()"
}

func (e EqDeep) String() string {
	return "EqDeep()"
}

type IxMatch struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e IxMatch) Short() string {
	return fmt.Sprintf("IxMatch(%v)", e.Optic.Short())
}

func (e IxMatch) String() string {
	return fmt.Sprintf("IxMatch(%v)", e.Optic.String())
}

type Concat struct {
	OpticTypeExpr
	Optics []OpticExpression
}

func (e Concat) Short() string {
	var sb strings.Builder

	sb.WriteString("Concat(")

	for i, v := range e.Optics {
		if i != 0 {
			sb.WriteString(" , ")
		}
		sb.WriteString(v.Short())
	}

	sb.WriteString(")")

	return sb.String()
}

func (e Concat) String() string {
	var sb strings.Builder

	sb.WriteString("Concat(")

	for i, v := range e.Optics {
		if i != 0 {
			sb.WriteString(" , ")
		}
		sb.WriteString(v.String())
	}

	sb.WriteString(")")

	return sb.String()
}

type Grouped struct {
	OpticTypeExpr
	Optic   OpticExpression
	Reducer ReducerExpression
}

func (e Grouped) Short() string {
	return fmt.Sprintf("Grouped(%v,%v)", e.Optic.Short(), e.Reducer.Short())
}

func (e Grouped) String() string {
	return fmt.Sprintf("Grouped(%v,%v)", e.Optic.String(), e.Reducer.String())
}

type Prefixed struct {
	OpticTypeExpr
	Path     OpticExpression
	Val      any
	ValMatch OpticExpression
}

func (e Prefixed) Short() string {
	return fmt.Sprintf("Prefixed(%v,%v,%v)", e.Path.Short(), e.Val, e.ValMatch.Short())
}

func (e Prefixed) String() string {
	return fmt.Sprintf("Prefixed(%v,%v,%v)", e.Path.String(), e.Val, e.ValMatch.String())
}

type WithVar struct {
	OpticTypeExpr
	WithVar OpticExpression
	Name    string
	Value   OpticExpression
}

func (e WithVar) Short() string {
	return fmt.Sprintf("WithVar(%v,%v,%v)", e.WithVar.Short(), e.Name, e.Value.Short())
}

func (e WithVar) String() string {
	return fmt.Sprintf("WithVar(%v,%v,%v)", e.WithVar.String(), e.Name, e.Value.String())
}

type Var struct {
	OpticTypeExpr
	Name string
}

func (e Var) Short() string {
	return fmt.Sprintf("WithVar(%v)", e.Name)
}

func (e Var) String() string {
	return fmt.Sprintf("WithVar(%v)", e.Name)
}

type ColSourceErr struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e ColSourceErr) Short() string {
	return fmt.Sprintf("ColSourceErr(%v)", e.Optic)
}

func (e ColSourceErr) String() string {
	return fmt.Sprintf("ColSourceErr(%v)", e.Optic)
}

type ColFocusErr struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e ColFocusErr) Short() string {
	return fmt.Sprintf("ColFocusErr(%v)", e.Optic)
}

func (e ColFocusErr) String() string {
	return fmt.Sprintf("ColFocusErr(%v)", e.Optic)
}

type ColSourceFocusErr struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e ColSourceFocusErr) Short() string {
	return fmt.Sprintf("ColSourceFocusErr(%v)", e.Optic)
}

func (e ColSourceFocusErr) String() string {
	return fmt.Sprintf("ColSourceFocusErr(%v)", e.Optic)
}

type WithMetrics struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e WithMetrics) Short() string {
	return fmt.Sprintf("WithMetrics(%v)", e.Optic)
}

func (e WithMetrics) String() string {
	return fmt.Sprintf("WithMetrics(%v)", e.Optic)
}

type SliceChildren struct {
	OpticTypeExpr
	A reflect.Type
}

func (e SliceChildren) Short() string {
	return "SliceChildren"
}

func (e SliceChildren) String() string {
	return fmt.Sprintf("SliceChildren(%v)", util.FullTypeName(e.A))
}
