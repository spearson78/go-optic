package expr

import (
	"fmt"
	"reflect"

	"github.com/spearson78/go-optic/internal/util"
)

type UnaryExpr struct {
	OpticTypeExpr
	Op string
}

func (e UnaryExpr) Short() string {
	return fmt.Sprintf("Unary(%v)", e.Op)
}

func (e UnaryExpr) String() string {
	return fmt.Sprintf("Unary(%v,%v)", util.FullTypeName(e.opticS), e.Op)
}

type ParseInt struct {
	OpticTypeExpr
	Base    int
	BitSize int
}

func (e ParseInt) Short() string {
	return fmt.Sprintf("ParseInt(%v,%v)", e.Base, e.BitSize)
}

func (e ParseInt) String() string {
	return fmt.Sprintf("ParseInt(%v,%v,%v)", util.FullTypeName(e.opticA), e.Base, e.BitSize)
}

type FormatInt struct {
	OpticTypeExpr
	Base int
}

func (e FormatInt) Short() string {
	return fmt.Sprintf("FormatInt(%v)", e.Base)
}

func (e FormatInt) String() string {
	return fmt.Sprintf("FormatInt(%v,%v)", util.FullTypeName(e.opticA), e.Base)
}

type ParseFloat struct {
	OpticTypeExpr
	Fmt     byte
	Prec    int
	BitSize int
}

func (e ParseFloat) Short() string {
	return fmt.Sprintf("ParseFloat(%v,%v,%v)", e.Fmt, e.Prec, e.BitSize)
}

func (e ParseFloat) String() string {
	return fmt.Sprintf("ParseFloat(%v,%v,%v,%v)", util.FullTypeName(e.opticA), e.Fmt, e.Prec, e.BitSize)
}

type FormatFloat struct {
	OpticTypeExpr
	Fmt     byte
	Prec    int
	BitSize int
}

func (e FormatFloat) Short() string {
	return fmt.Sprintf("FormatFloat(%v,%v,%v)", e.Fmt, e.Prec, e.BitSize)
}

func (e FormatFloat) String() string {
	return fmt.Sprintf("FormatFloat(%v,%v,%v,%v)", util.FullTypeName(e.opticA), e.Fmt, e.Prec, e.BitSize)
}

type ParseFloatP struct {
	OpticTypeExpr
	BitSize int
}

func (e ParseFloatP) Short() string {
	return fmt.Sprintf("ParseFloatP(%v)", e.BitSize)
}

func (e ParseFloatP) String() string {
	return fmt.Sprintf("ParseFloatP(%v,%v)", util.FullTypeName(e.opticA), e.BitSize)
}

type BinaryExpr struct {
	OpticTypeExpr
	L  reflect.Type
	R  reflect.Type
	Op string
}

func (e BinaryExpr) Short() string {
	return e.Op
}

func (e BinaryExpr) String() string {
	return fmt.Sprintf("Binary(%v)", e.Op)
}

type OpT2BindExpr struct {
	OpticTypeExpr
	Op         OpticExpression
	S          reflect.Type
	RightValue any
}

func (e OpT2BindExpr) Short() string {
	return fmt.Sprintf("%v %v", e.Op.Short(), e.RightValue)
}

func (e OpT2BindExpr) String() string {
	return fmt.Sprintf("OpT2Bind(%v,%v,%v)", util.FullTypeName(e.S), e.Op, e.RightValue)
}

type IsoOpT2BindExpr struct {
	OpticTypeExpr
	Op    OpticExpression
	InvOp OpticExpression
	S     reflect.Type
	Right any
}

func (e IsoOpT2BindExpr) Short() string {
	return fmt.Sprintf("%v %v", e.Op.Short(), e.Right)
}

func (e IsoOpT2BindExpr) String() string {
	return fmt.Sprintf("IsoOpT2BindExpr(%v,%v,%v,%v)", util.FullTypeName(e.S), e.Op, e.InvOp, e.Right)
}
