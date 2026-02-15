package expr

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spearson78/go-optic/internal/util"
)

type ReducerExpression interface {
	ReducerS() reflect.Type
	ReducerE() reflect.Type
	ReducerR() reflect.Type

	//Provides a string representation of this Reducer for diagnosis/logging purposes.
	String() string
	Short() string
}

type ReducerTypeExpr struct {
	reducerS reflect.Type
	reducerA reflect.Type
	reducerR reflect.Type
}

func (o ReducerTypeExpr) ReducerS() reflect.Type {
	return o.reducerS
}

func (o ReducerTypeExpr) ReducerE() reflect.Type {
	return o.reducerA
}

func (o ReducerTypeExpr) ReducerR() reflect.Type {
	return o.reducerR
}

func NewReducerTypeExpr[S, A, R any]() ReducerTypeExpr {
	return ReducerTypeExpr{
		reducerS: reflect.TypeFor[S](),
		reducerA: reflect.TypeFor[A](),
		reducerR: reflect.TypeFor[R](),
	}
}

type ReducerN struct {
	ReducerTypeExpr
	Reducers []ReducerExpression
}

func (m ReducerN) Short() string {

	var sb strings.Builder

	sb.WriteString("ReducerN(")

	for i, v := range m.Reducers {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(v.Short())
	}

	sb.WriteString(")")

	return sb.String()
}

func (m ReducerN) String() string {

	var sb strings.Builder

	sb.WriteString("ReducerN[")
	sb.WriteString(util.FullTypeName(m.reducerS))
	sb.WriteString(",")
	sb.WriteString(util.FullTypeName(m.reducerA))
	sb.WriteString("](")

	for i, v := range m.Reducers {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(v.String())
	}

	sb.WriteString(")")

	return sb.String()
}

type AsReducer struct {
	ReducerTypeExpr
	Empty  OpticExpression
	Append OpticExpression
	End    OpticExpression
}

func (m AsReducer) Short() string {
	return fmt.Sprintf("AsReducer(%v,%v)", m.Empty, m.Append)
}

func (m AsReducer) String() string {
	return fmt.Sprintf("AsReducer[%v,%v](%v,%v)", util.FullTypeName(m.reducerS), util.FullTypeName(m.reducerA), m.Empty, m.Append)
}

type MeanReducer struct {
	ReducerTypeExpr
}

func (m MeanReducer) Short() string {
	return "MeanReducer"
}

func (m MeanReducer) String() string {
	return fmt.Sprintf("MeanReducer[%v,%v]()", util.FullTypeName(m.reducerS), util.FullTypeName(m.reducerA))
}

type MedianReducer struct {
	ReducerTypeExpr
}

func (m MedianReducer) Short() string {
	return "MedianReducer"
}

func (m MedianReducer) String() string {
	return fmt.Sprintf("MedianReducer[%v,%v]()", util.FullTypeName(m.reducerS), util.FullTypeName(m.reducerA))
}

type ModeReducer struct {
	ReducerTypeExpr
}

func (m ModeReducer) Short() string {
	return "ModeReducer"
}

func (m ModeReducer) String() string {
	return fmt.Sprintf("ModeReducer[%v,%v]()", util.FullTypeName(m.reducerS), util.FullTypeName(m.reducerA))
}

func CustomReducer(id string) func(ot ReducerTypeExpr) ReducerExpression {
	return func(ot ReducerTypeExpr) ReducerExpression {
		return CustomReducerExpr{
			ReducerTypeExpr: ot,
			Id:              id,
		}
	}
}

type CustomReducerExpr struct {
	ReducerTypeExpr
	Id string
}

func (m CustomReducerExpr) Short() string {
	return fmt.Sprintf("CustomReducer(%v)", m.Id)
}

func (m CustomReducerExpr) String() string {
	return fmt.Sprintf("CustomReducer(%v)", m.Id)
}

type StringBuilderReducerExpr struct {
	ReducerTypeExpr
}

func (m StringBuilderReducerExpr) Short() string {
	return "StringBuilderReducerExpr()"
}

func (m StringBuilderReducerExpr) String() string {
	return "StringBuilderReducerExpr()"
}

type AppendSliceReducerExpr struct {
	ReducerTypeExpr
}

func (m AppendSliceReducerExpr) Short() string {
	return "AppendSliceReducerExpr()"
}

func (m AppendSliceReducerExpr) String() string {
	return "AppendSliceReducerExpr()"
}

type ErrDuplicateKeyReducerExpr struct {
	ReducerTypeExpr
}

func (m ErrDuplicateKeyReducerExpr) Short() string {
	return "ErrDuplicateKeyReducerExpr()"
}

func (m ErrDuplicateKeyReducerExpr) String() string {
	return "ErrDuplicateKeyReducerExpr()"
}
