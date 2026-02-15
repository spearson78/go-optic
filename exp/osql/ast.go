package osql

import (
	"reflect"
)

type Expression interface {
	Accept(v Visitor) error
}

type SelectStmnt struct {
	Distinct      bool
	PkColumns     []Expression
	IndexColumns  []Expression
	SelectColumns []Expression
	DeriveColumns []Expression
	ExceptColumns []ColumnIdentifier
	From          Expression
	Join          []JoinDef
	Where         []Expression
	Limit         Expression
	Offset        Expression
	OrderBy       []Expression
	GroupBy       []Expression
	Having        []Expression

	ContainsAggregate bool

	Scan RowScan
}

func (e *SelectStmnt) Accept(v Visitor) error {
	return v.VisitSelect(e)
}

type DeleteStmnt struct {
	From  Expression
	Where []Expression
}

func (e *DeleteStmnt) Accept(v Visitor) error {
	return v.VisitDelete(e)
}

type SetColumn struct {
	ColumnIdentifier
	Value Expression
}

type UpdateStmnt struct {
	SetColumns []SetColumn
	Table      Expression
	Where      []Expression
}

func (e *UpdateStmnt) Accept(v Visitor) error {
	return v.VisitUpdate(e)
}

type InsertStmnt struct {
	Table  Expression
	Values []SetColumn
}

func (e *InsertStmnt) Accept(v Visitor) error {
	return v.VisitInsert(e)
}

func GetTableName(e Expression) (*TableIdentifier, bool) {
	switch t := e.(type) {
	case *TableIdentifier:
		return t, true
	case *TableDef:
		return &TableIdentifier{
			Table: t.Name,
		}, true
	default:
		return &TableIdentifier{}, false
	}
}

type TableIdentifier struct {
	Database string
	Schema   string
	Table    string

	TableAlias string
	Type       reflect.Type
}

func (e *TableIdentifier) Accept(v Visitor) error {
	return v.VisitTableIdentifier(e)
}

type ColumnIdentifier struct {
	TableIdentifier
	Name string

	ReflectFieldNum int
}

func (e *ColumnIdentifier) Accept(v Visitor) error {
	return v.VisitColumnIdentifier(e)
}

type Alias struct {
	Column Expression
	Alias  string
}

func (e *Alias) Accept(v Visitor) error {
	return v.VisitAlias(e)
}

type ColumnExpression struct {
	Value Expression
	Alias string
}

func (e *ColumnExpression) Accept(v Visitor) error {
	return v.VisitColumnExpression(e)
}

type Keyword string

func (e Keyword) Accept(v Visitor) error {
	return v.VisitKeyword(e)
}

type TableDef struct {
	Name string
	Expr Expression
}

func (e *TableDef) Accept(v Visitor) error {
	return v.VisitTableDef(e)
}

type ColumnDef struct {
	Name string
	Expr Expression
}

func (e *ColumnDef) Accept(v Visitor) error {
	return v.VisitColumnDef(e)
}

type List []Expression

func (e List) Accept(v Visitor) error {
	return v.VisitList(e)
}

type Number float64

func (e Number) Accept(v Visitor) error {
	return v.VisitNumber(e)
}

type String string

func (e String) Accept(v Visitor) error {
	return v.VisitString(e)
}

type Date struct {
	DateString string
	Type       string
}

func (e Date) Accept(v Visitor) error {
	return v.VisitDate(e)
}

type BinaryExpression struct {
	Left  Expression
	Op    string
	Right Expression
}

func (e *BinaryExpression) Accept(v Visitor) error {
	return v.VisitBinaryExpression(e)
}

type DateDiff struct {
	Left  Expression
	Right Expression
}

func (e *DateDiff) Accept(v Visitor) error {
	return v.VisitDateDiff(e)
}

type DateAdd struct {
	Duration *Duration
	Date     Expression
}

func (e *DateAdd) Accept(v Visitor) error {
	return v.VisitDateAdd(e)
}

type Duration struct {
	Val  float64
	Type string
}

func (e Duration) Accept(v Visitor) error {
	return v.VisitDuration(e)
}

type Call struct {
	Function  Expression
	Params    []Expression
	Aggregate bool
}

func (e *Call) Accept(v Visitor) error {
	return v.VisitCall(e)
}

type JoinSide int

const (
	JOIN_INNER = iota
	JOIN_LEFT
	JOIN_RIGHT
)

type JoinDef struct {
	Side      JoinSide
	JoinTable Expression
	JoinExpr  Expression
}

func (e *JoinDef) Accept(v Visitor) error {
	return v.VisitJoin(e)
}

type Exists struct {
	JoinTable Expression
	JoinExpr  Expression
	SubSelect Expression
}

func (e *Exists) Accept(v Visitor) error {
	return v.VisitExists(e)
}

type Between struct {
	Value Expression
	From  Expression
	To    Expression
}

func (e *Between) Accept(v Visitor) error {
	return v.VisitBetween(e)
}

type CoalesceExpr struct {
	Value   Expression
	Default Expression
}

func (e *CoalesceExpr) Accept(v Visitor) error {
	return v.VisitCoalesce(e)
}

type NegateExpr struct {
	Expr Expression
}

func (e *NegateExpr) Accept(v Visitor) error {
	return v.VisitNegate(e)
}

type Positive struct {
	Expr Expression
}

func (e *Positive) Accept(v Visitor) error {
	return v.VisitPositive(e)
}

type SqlNot struct {
	Expr Expression
}

func (e *SqlNot) Accept(v Visitor) error {
	return v.VisitNot(e)
}

type ColumnList struct {
	Columns []Expression
}

func (e *ColumnList) Accept(v Visitor) error {
	return v.VisitColumnList(e)
}

type CaseCondition struct {
	Condition Expression
	Value     Expression
}

type CaseExpr struct {
	Conditions []*CaseCondition
}

func (e *CaseExpr) Accept(v Visitor) error {
	return v.VisitCase(e)
}

type Bool bool

const (
	SqlTrue  Bool = true
	SqlFalse Bool = false
)

func (e Bool) Accept(v Visitor) error {
	return v.VisitBool(e)
}

type Param string

func (e Param) Accept(v Visitor) error {
	return v.VisitParam(e)
}

type FuncDef struct {
	Params      []ColumnIdentifier
	NamedParams map[string]Expression
}

func (e *FuncDef) Accept(v Visitor) error {
	return v.VisitFuncDef(e)
}

type NamedValue struct {
	Name  string
	Value Expression
}

func (e NamedValue) Accept(v Visitor) error {
	return v.VisitNamedValue(e)
}

type SetOperation struct {
	Left      *SelectStmnt
	Operation string
	Right     *SelectStmnt
}

func (e *SetOperation) Accept(v Visitor) error {
	return v.VisitSetOperation(e)
}

type BoundaryType int

const (
	RowsBoundary BoundaryType = iota
	RangeBoundary
)

type OverExpr struct {
	Expr        Expression
	PartitionBy []Expression
	OrderBy     []Expression
	Between     Expression
}

func (e *OverExpr) Accept(v Visitor) error {
	return v.VisitOver(e)
}

type ValueExpr []Expression

func (e ValueExpr) Accept(v Visitor) error {
	return v.VisitValue(e)
}

type Values []ValueExpr

func (e *Values) Accept(v Visitor) error {
	return v.VisitValues(e)
}

type Null struct{}

func (e Null) Accept(v Visitor) error {
	return v.VisitNull()
}

type RawSql []Expression

func (e RawSql) Accept(v Visitor) error {
	return v.VisitRawSql(e)
}

type FrameBetween struct {
	BoundaryType BoundaryType
	From         Expression
	To           Expression
}

func (e *FrameBetween) Accept(v Visitor) error {
	return v.VisitFrameBetween(e)
}
