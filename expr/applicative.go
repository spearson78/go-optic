package expr

import (
	"fmt"
	"reflect"

	"github.com/spearson78/go-optic/internal/util"
)

type WithOption struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e WithOption) Short() string {
	return fmt.Sprintf("WithOption(%v)", e.Optic.Short())
}

func (e WithOption) String() string {
	return fmt.Sprintf("WithOption(%v)", e.Optic.String())
}

type WithFunc struct {
	OpticTypeExpr
	P     reflect.Type
	Optic OpticExpression
}

func (e WithFunc) Short() string {
	return fmt.Sprintf("WithFunc(%v)", e.Optic.Short())
}

func (e WithFunc) String() string {
	return fmt.Sprintf("WithFunc(%v,%v)", util.FullTypeName(e.P), e.Optic.String())
}

type WithComprehension struct {
	OpticTypeExpr
	Optic OpticExpression
}

func (e WithComprehension) Short() string {
	return fmt.Sprintf("WithComprehension(%v)", e.Optic.Short())
}

func (e WithComprehension) String() string {
	return fmt.Sprintf("WithComprehension(%v)", e.Optic.String())
}

type WithEither struct {
	OpticTypeExpr
	P     reflect.Type
	Optic OpticExpression
}

func (e WithEither) Short() string {
	return fmt.Sprintf("WithEither(%v)", e.Optic.Short())
}

func (e WithEither) String() string {
	return fmt.Sprintf("WithEither(%v,%v)", util.FullTypeName(e.P), e.Optic.String())
}

type WithValidation struct {
	OpticTypeExpr
	V     reflect.Type
	Optic OpticExpression
}

func (e WithValidation) Short() string {
	return fmt.Sprintf("WithValidation(%v)", e.Optic.Short())
}

func (e WithValidation) String() string {
	return fmt.Sprintf("WithValidation(%v,%v)", util.FullTypeName(e.V), e.Optic.String())
}
