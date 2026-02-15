package osql

import (
	"context"
	"database/sql"
	"reflect"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

type tableHandler struct {
}

// Get implements optic.ExprHandler.
func (t *tableHandler) Get(ctx context.Context, expr expr.OpticExpression, source any) (index any, value any, found bool, err error) {

	ix, res, err := execGet(ctx, expr, source.(*sql.DB))

	return ix, res, err == nil, err
}

// Modify implements optic.ExprHandler.
func (t *tableHandler) Modify(ctx context.Context, o expr.OpticExpression, fmapExpr expr.OpticExpression, fmap func(index any, focus any, focusErr error) (any, error), source any) (any, bool, error) {

	res, err := execModify(ctx, o, fmapExpr, source.(*sql.DB))

	return res, err == nil, err
}

// ReverseGet implements optic.ExprHandler.
func (t *tableHandler) ReverseGet(ctx context.Context, expr expr.OpticExpression, focus any) (any, error) {
	return nil, UnsupportedOpticMethod
}

// Set implements optic.ExprHandler.
func (t *tableHandler) Set(ctx context.Context, o expr.OpticExpression, source any, val any) (any, error) {
	res, err := execModify(ctx, o, expr.Const{
		Value: val,
	}, source.(*sql.DB))

	return res, err
}

// TypeId implements optic.ExprHandler.
func (t *tableHandler) TypeId() string {
	return "github.com/spearson78/go-optic/osql/tablehandler"
}

type TableExpr struct {
	expr.OpticTypeExpr
	Name string
	I    reflect.Type
	A    reflect.Type
}

// Short implements expr.OpticExpression.
func (t TableExpr) Short() string {
	return "TableExpr"
}

// String implements expr.OpticExpression.
func (t TableExpr) String() string {
	return "TableExpr"
}

type JoinExpr struct {
	expr.OpticTypeExpr
	TableName  string
	ColumnName string
}

// Short implements expr.OpticExpression.
func (t JoinExpr) Short() string {
	return "JoinExpr"
}

// String implements expr.OpticExpression.
func (t JoinExpr) String() string {
	return "JoinExpr"
}

type JoinMExpr struct {
	expr.OpticTypeExpr
	TableName  string
	ColumnName string
}

// Short implements expr.OpticExpression.
func (t JoinMExpr) Short() string {
	return "JoinMExpr"
}

// String implements expr.OpticExpression.
func (t JoinMExpr) String() string {
	return "JoinMExpr"
}

func Table[I, A any](name string) Optic[Void, *sql.DB, *sql.DB, Collection[I, A, Err], Collection[I, A, Err], ReturnMany, ReadWrite, UniDir, Err] {

	return ExprOptic[Void, *sql.DB, *sql.DB, Collection[I, A, Err], Collection[I, A, Err], ReturnMany, ReadWrite, UniDir, Err](
		&tableHandler{},
		func(t expr.OpticTypeExpr) expr.OpticExpression {
			return TableExpr{
				OpticTypeExpr: t,
				Name:          name,
				I:             reflect.TypeFor[I](),
				A:             reflect.TypeFor[A](),
			}
		},
	)

}

func Join[S, A any](tableName string, columnName string) Optic[Void, S, S, A, A, ReturnMany, ReadWrite, UniDir, Err] {

	return ExprOptic[Void, S, S, A, A, ReturnMany, ReadWrite, UniDir, Err](
		&tableHandler{},
		func(t expr.OpticTypeExpr) expr.OpticExpression {
			return JoinExpr{
				OpticTypeExpr: t,
				TableName:     tableName,
				ColumnName:    columnName,
			}
		},
	)

}

func JoinM[S, A any](tableName string, columnName string) Optic[Void, S, S, A, A, ReturnMany, ReadWrite, UniDir, Err] {

	return ExprOptic[Void, S, S, A, A, ReturnMany, ReadWrite, UniDir, Err](
		&tableHandler{},
		func(t expr.OpticTypeExpr) expr.OpticExpression {
			return JoinMExpr{
				OpticTypeExpr: t,
				TableName:     tableName,
				ColumnName:    columnName,
			}
		},
	)

}
