package osql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type SqlContext struct {
	This      Expression
	That      Expression
	Params    map[string]Expression
	Resolving map[string]struct{}
}

type Sqlite3 struct {
	s *SqlContext
	b strings.Builder
}

// VisitDateAdd implements Visitor.
func (visitor *Sqlite3) VisitDateAdd(e *DateAdd) error {

	visitor.b.WriteString("(")
	e.Date.Accept(visitor)

	visitor.b.WriteString(" + ")

	dtype := ""

	switch e.Duration.Type {
	case "days":
		dtype = "DAY"
	case "hours":
		dtype = "HOUR"
	case "minutes":
		dtype = "MINUTE"
	case "seconds":
		dtype = "SECOND"
	case "months":
		dtype = "MONTH"
	case "years":
		dtype = "YEAR"
	default:
		return fmt.Errorf("unknown duration type %v", e.Duration.Type)
	}

	visitor.b.WriteString("INTERVAL ")
	visitor.b.WriteString(strconv.FormatFloat(e.Duration.Val, 'f', -1, 64))
	visitor.b.WriteString(" ")
	visitor.b.WriteString(dtype)
	visitor.b.WriteString(")")

	return nil
}

// VisitDateDiff implements Visitor.
func (visitor *Sqlite3) VisitDateDiff(e *DateDiff) error {
	visitor.b.WriteString("(")

	err := e.Left.Accept(visitor)
	if err != nil {
		return err
	}

	visitor.b.WriteString(" - ")

	err = e.Right.Accept(visitor)
	if err != nil {
		return err
	}

	visitor.b.WriteString(")")

	return nil
}

func ToSqlite3(expr Expression) (string, error) {

	v := &Sqlite3{
		s: &SqlContext{
			Params:    make(map[string]Expression),
			Resolving: make(map[string]struct{}),
		},
	}
	err := expr.Accept(v)
	return v.b.String(), err

}

func (visitor *Sqlite3) EscapeIdentifier(name string) {
	visitor.b.WriteString("\"")
	visitor.b.WriteString(name)
	visitor.b.WriteString("\"")
}

// VisitBetween implements Visitor.
func (visitor *Sqlite3) VisitBetween(e *Between) error {

	err := e.Value.Accept(visitor)
	if err != nil {
		return err
	}

	visitor.b.WriteString(" BETWEEN ")

	err = e.From.Accept(visitor)
	if err != nil {
		return err
	}

	visitor.b.WriteString(" AND ")

	err = e.To.Accept(visitor)
	if err != nil {
		return err
	}

	return nil
}

func (visitor *Sqlite3) VisitList(e List) error {

	visitor.b.WriteString("(")

	for i, v := range e {
		if i != 0 {
			visitor.b.WriteString(",")
		}
		err := v.Accept(visitor)
		if err != nil {
			return err
		}
	}

	visitor.b.WriteString(")")

	return nil
}

// VisitBinaryExpression implements Visitor.
func (visitor *Sqlite3) VisitBinaryExpression(e *BinaryExpression) error {

	_, rNull := e.Right.(Null)
	_, lNull := e.Left.(Null)

	if rNull {
		if lNull {
			switch e.Op {
			case "=":
				SqlTrue.Accept(visitor)
				return nil
			case "!=":
				SqlFalse.Accept(visitor)
				return nil
			default:
				//normal case fall through
			}
		} else {
			switch e.Op {
			case "=":
				visitor.b.WriteString("(")
				err := e.Left.Accept(visitor)
				visitor.b.WriteString(" IS NULL)")
				return err
			case "!=":
				visitor.b.WriteString("(")
				err := e.Left.Accept(visitor)
				visitor.b.WriteString(" IS NOT NULL)")
				return err
			default:
				//normal case fall through
			}
		}
	} else {
		if lNull {
			switch e.Op {
			case "=":
				visitor.b.WriteString("(")
				err := e.Right.Accept(visitor)
				visitor.b.WriteString(" IS NULL)")
				return err
			case "!=":
				visitor.b.WriteString("(")
				err := e.Right.Accept(visitor)
				visitor.b.WriteString(" IS NOT NULL)")
				return err
			default:
				//normal case fall through
			}
		} else {
			//normal case fall through
		}
	}

	visitor.b.WriteString("(")

	_, isLSelect := e.Left.(*SelectStmnt)
	if isLSelect {
		visitor.b.WriteString("(")
	}

	err := e.Left.Accept(visitor)
	if err != nil {
		return err
	}

	if isLSelect {
		visitor.b.WriteString(")")
	}

	visitor.b.WriteString(" ")
	visitor.b.WriteString(e.Op)
	visitor.b.WriteString(" ")

	_, isRSelect := e.Right.(*SelectStmnt)
	if isRSelect {
		visitor.b.WriteString("(")
	}

	err = e.Right.Accept(visitor)
	if err != nil {
		return err
	}

	if isRSelect {
		visitor.b.WriteString(")")
	}

	visitor.b.WriteString(")")

	return nil
}

func (visitor *Sqlite3) VisitExists(e *Exists) error {
	visitor.b.WriteString("EXISTS (")

	visitor.s.That = visitor.s.This

	err := e.SubSelect.Accept(visitor)
	if err != nil {
		return err
	}

	visitor.b.WriteString(")")

	return nil
}

// VisitBool implements Visitor.
func (visitor *Sqlite3) VisitBool(e Bool) error {
	if e {
		visitor.b.WriteString("TRUE")
	} else {
		visitor.b.WriteString("FALSE")
	}

	return nil
}

// VisitCall implements Visitor.
func (visitor *Sqlite3) VisitCall(e *Call) error {

	sqlFunc := ""
	distinct := false

	params := e.Params

	switch f := e.Function.(type) {
	case *ColumnIdentifier:
		switch f.Name {
		case "count":
			sqlFunc = "COUNT"
		case "avg":
			sqlFunc = "AVG"
		case "average":
			sqlFunc = "AVG"
		case "sum":
			sqlFunc = "SUM"
		case "min":
			sqlFunc = "MIN"
		case "max":
			sqlFunc = "MAX"
		case "stddev":
			sqlFunc = "STDDEV"
		case "count_distinct":
			sqlFunc = "COUNT"
			distinct = true
		case "round":
			sqlFunc = "ROUND"
			if len(e.Params) != 2 {
				return fmt.Errorf("round requirs 2 params")
			}
			params = []Expression{params[1], params[0]}
		case "rank":
			sqlFunc = "RANK"
			params = nil
		default:
			return fmt.Errorf("unknown function %v", f.Name)
		}
	case *FuncDef:
		newCtx := &SqlContext{
			Params:    make(map[string]Expression),
			Resolving: make(map[string]struct{}),
		}

		for k, v := range visitor.s.Params {
			newCtx.Params[k] = v
		}

		pos := 0
		for _, p := range params {
			if namedParam, isNamed := p.(*NamedValue); isNamed {
				newCtx.Params[namedParam.Name] = namedParam.Value
			} else {
				newCtx.Params[f.Params[pos].Name] = p
				pos++
			}
		}

		//Insert default values for missing named parameters
		for k, v := range f.NamedParams {
			if _, found := newCtx.Params[k]; !found {
				newCtx.Params[k] = v
			}
		}

		err := f.Accept(visitor)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("function not supported %T", f)
	}

	visitor.b.WriteString(sqlFunc)
	visitor.b.WriteString("(")

	if distinct {
		visitor.b.WriteString("DISTINCT ")
	}

	for i, p := range params {
		if i != 0 {
			visitor.b.WriteString(",")
		}

		switch pe := p.(type) {
		case *ColumnIdentifier:
			if pe.Name == "this" {
				visitor.b.WriteString("*")
			} else {
				err := p.Accept(visitor)
				if err != nil {
					return err
				}
			}
		default:
			err := p.Accept(visitor)
			if err != nil {
				return err
			}
		}
	}

	visitor.b.WriteString(")")

	return nil
}

// VisitCase implements Visitor.
func (visitor *Sqlite3) VisitCase(e *CaseExpr) error {

	if len(e.Conditions) == 0 {
		Null{}.Accept(visitor)
		return nil
	}

	visitor.b.WriteString("CASE ")

	hasDefault := false

	for i, v := range e.Conditions {
		if boolCond, isBool := v.Condition.(Bool); isBool {
			if boolCond {
				hasDefault = true

				if i != 0 {
					visitor.b.WriteString("ELSE ")
				}

				err := v.Value.Accept(visitor)
				if err != nil {
					return err
				}

				break // true => ends the case staements
			} else {
				//Ignore false conditions
			}
		} else {

			visitor.b.WriteString("WHEN ")

			err := v.Condition.Accept(visitor)
			if err != nil {
				return err
			}

			visitor.b.WriteString(" THEN ")

			err = v.Value.Accept(visitor)
			if err != nil {
				return err
			}

			visitor.b.WriteString(" ")
		}
	}

	if !hasDefault {
		visitor.b.WriteString("ELSE NULL")
	}

	visitor.b.WriteString(" END")

	return nil
}

// VisitCoalesce implements Visitor.
func (visitor *Sqlite3) VisitCoalesce(e *CoalesceExpr) error {

	visitor.b.WriteString("COALESCE(")

	err := e.Value.Accept(visitor)
	if err != nil {
		return err
	}

	visitor.b.WriteString(",")

	err = e.Default.Accept(visitor)
	if err != nil {
		return err
	}

	visitor.b.WriteString(")")

	return nil
}

// VisitColumnDef implements Visitor.
func (visitor *Sqlite3) VisitColumnDef(e *ColumnDef) error {

	if e.Name != "" {
		visitor.s.Resolving[e.Name] = struct{}{}
	}
	err := e.Expr.Accept(visitor)
	if err != nil {
		return err
	}
	if e.Name != "" {
		delete(visitor.s.Resolving, e.Name)

		visitor.b.WriteString(" AS ")
		visitor.EscapeIdentifier(e.Name)
	}

	return nil

}

// VisitColumnIdentifier implements Visitor.
func (visitor *Sqlite3) VisitColumnIdentifier(e *ColumnIdentifier) error {

	if e.TableAlias != "" {
		visitor.EscapeIdentifier(e.TableAlias)
		visitor.b.WriteString(".")
	} else if e.Table == "" {
		_, resolving := visitor.s.Resolving[e.Name]
		expr, isExpr := visitor.s.Params[e.Name]
		if !resolving && isExpr {
			visitor.s.Resolving[e.Name] = struct{}{}
			err := expr.Accept(visitor)
			delete(visitor.s.Resolving, e.Name)
			return err
		}
	} else {
		switch e.Table {
		case "this":
			if visitor.s.This == nil {
				return fmt.Errorf("no this defined")
			}
			err := visitor.s.This.Accept(visitor)
			if err != nil {
				return err
			}
		case "that":
			if visitor.s.That == nil {
				return fmt.Errorf("no that defined")
			}

			err := visitor.s.That.Accept(visitor)
			if err != nil {
				return err
			}
		default:
			visitor.EscapeIdentifier(e.Table)
		}

		visitor.b.WriteString(".")
	}

	if e.Name == "*" {
		visitor.b.WriteString("*")
	} else {
		visitor.EscapeIdentifier(e.Name)
	}

	return nil
}

func (visitor *Sqlite3) VisitAlias(e *Alias) error {

	err := e.Column.Accept(visitor)
	if err != nil {
		return err
	}

	if e.Alias != "" {
		visitor.b.WriteString(" AS ")
		visitor.EscapeIdentifier(e.Alias)
	}

	return nil

}

func (visitor *Sqlite3) VisitColumnExpression(e *ColumnExpression) error {

	visitor.b.WriteString("(")
	err := e.Value.Accept(visitor)
	if err != nil {
		return err
	}
	visitor.b.WriteString(")")

	if e.Alias != "" {
		visitor.b.WriteString(" AS ")
		visitor.EscapeIdentifier(e.Alias)
	}

	return nil
}

// VisitColumnList implements Visitor.
func (visitor *Sqlite3) VisitColumnList(e *ColumnList) error {

	for i, v := range e.Columns {
		if i != 0 {
			visitor.b.WriteString(",")
		}

		err := v.Accept(visitor)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitDate implements Visitor.
func (visitor *Sqlite3) VisitDate(e Date) error {
	visitor.b.WriteString(e.Type + " '" + e.DateString + "'")
	return nil
}

// VisitDuration implements Visitor.
func (visitor *Sqlite3) VisitDuration(e Duration) error {

	dtype := ""

	switch e.Type {
	case "days":
		dtype = "DAY"
	case "hours":
		dtype = "HOUR"
	case "minutes":
		dtype = "MINUTE"
	case "seconds":
		dtype = "SECOND"
	case "months":
		dtype = "MONTH"
	case "years":
		dtype = "YEAR"
	default:
		return fmt.Errorf("unknown duration type %v", e.Type)
	}

	visitor.b.WriteString("INTERVAL " + strconv.FormatFloat(e.Val, 'f', -1, 64) + " " + dtype)
	return nil

}

// VisitFrameBetween implements Visitor.
func (visitor *Sqlite3) VisitFrameBetween(e *FrameBetween) error {

	visitor.b.WriteString(" ")

	boundary := "ROWS"
	if e.BoundaryType == RangeBoundary {
		boundary = "RANGE"
	}
	visitor.b.WriteString(boundary)

	visitor.b.WriteString(" BETWEEN ")

	if e.From == nil {
		visitor.b.WriteString("UNBOUNDED PRECEDING")
	} else {

		fromNum, ok := e.From.(Number)
		if !ok {
			return fmt.Errorf("expected number for window from : %T", e.From)
		}

		if fromNum == 0 {
			visitor.b.WriteString("CURRENT ROW")
		} else if fromNum > 0 {
			err := fromNum.Accept(visitor)
			if err != nil {
				return err
			}
			visitor.b.WriteString(" FOLLOWING")
		} else {
			err := Number(-fromNum).Accept(visitor)
			if err != nil {
				return err
			}
			visitor.b.WriteString(" PRECEDING")
		}
	}

	visitor.b.WriteString(" AND ")

	if e.To == nil {
		visitor.b.WriteString("UNBOUNDED FOLLOWING")
	} else {

		toNum, ok := e.To.(Number)
		if !ok {
			return fmt.Errorf("expected number for window from : %T", e.To)
		}

		if toNum == 0 {
			visitor.b.WriteString("CURRENT ROW")
		} else if toNum > 0 {
			err := toNum.Accept(visitor)
			if err != nil {
				return err
			}
			visitor.b.WriteString(" FOLLOWING")
		} else {
			err := Number(-toNum).Accept(visitor)
			if err != nil {
				return err
			}
			visitor.b.WriteString(" PRECEDING")
		}
	}

	return nil
}

// VisitFuncDef implements Visitor.
func (*Sqlite3) VisitFuncDef(expr *FuncDef) error {
	return fmt.Errorf("FuncDef to SQL not supported")
}

// VisitJoin implements Visitor.
func (visitor *Sqlite3) VisitJoin(e *JoinDef) error {

	joinType := "INNER"

	switch e.Side {
	case JOIN_INNER:
		joinType = "INNER"
	case JOIN_LEFT:
		joinType = "LEFT"
	case JOIN_RIGHT:
		joinType = "RIGHT"
	default:
		return fmt.Errorf("unknown join type %v", e.Side)
	}

	visitor.b.WriteString(joinType)
	visitor.b.WriteString(" JOIN ")

	err := e.JoinTable.Accept(visitor)
	if err != nil {
		return err
	}

	visitor.b.WriteString(" ON ")

	switch t := e.JoinTable.(type) {
	case *TableDef:
		visitor.s.That = &TableIdentifier{Table: t.Name}
	default:
		visitor.s.That = e.JoinTable
	}

	err = e.JoinExpr.Accept(visitor)
	if err != nil {
		return err
	}

	return nil
}

// VisizKeyword implements Visitor.
func (visitor *Sqlite3) VisitKeyword(expr Keyword) error {
	visitor.b.WriteString(string(expr))
	return nil
}

// VisitNamedValue implements Visitor.
func (visitor *Sqlite3) VisitNamedValue(e NamedValue) error {
	value, ok := visitor.s.Params[e.Name]
	if !ok {
		return e.Value.Accept(visitor)
	} else {
		return value.Accept(visitor)
	}
}

// VisitNegate implements Visitor.
func (visitor *Sqlite3) VisitNegate(e *NegateExpr) error {
	visitor.b.WriteString("-")
	return e.Expr.Accept(visitor)
}

// VisitNot implements Visitor.
func (visitor *Sqlite3) VisitNot(e *SqlNot) error {
	visitor.b.WriteString("NOT ")
	return e.Expr.Accept(visitor)
}

// VisitNull implements Visitor.
func (visitor *Sqlite3) VisitNull() error {
	visitor.b.WriteString("NULL")
	return nil
}

// VisitNumber implements Visitor.
func (visitor *Sqlite3) VisitNumber(e Number) error {
	visitor.b.WriteString(strconv.FormatFloat(float64(e), 'f', -1, 64))
	return nil
}

// VisitOver implements Visitor.
func (visitor *Sqlite3) VisitOver(e *OverExpr) error {

	err := e.Expr.Accept(visitor)
	if err != nil {
		return err
	}

	visitor.b.WriteString(" OVER (")
	if len(e.PartitionBy) > 0 {
		visitor.b.WriteString(" PARTITION BY ")
		for i, v := range e.PartitionBy {
			if i != 0 {
				visitor.b.WriteString(",")
			}

			err := v.Accept(visitor)
			if err != nil {
				return err
			}
		}
	}

	if len(e.OrderBy) > 0 {
		visitor.b.WriteString(" ORDER BY ")
		for i, v := range e.OrderBy {
			if i != 0 {
				visitor.b.WriteString(",")
			}

			err := v.Accept(visitor)
			if err != nil {
				return err
			}
		}
	}

	if e.Between != nil {
		err := e.Between.Accept(visitor)
		if err != nil {
			return err
		}
	}

	visitor.b.WriteString(")")

	return nil
}

// VisitParam implements Visitor.
func (visitor *Sqlite3) VisitParam(e Param) error {
	visitor.b.WriteString("$")
	visitor.b.WriteString(string(e))
	return nil
}

// VisitPositive implements Visitor.
func (visitor *Sqlite3) VisitPositive(e *Positive) error {
	visitor.b.WriteString("+")
	return e.Expr.Accept(visitor)
}

// VisitRawSql implements Visitor.
func (visitor *Sqlite3) VisitRawSql(e RawSql) error {

	for _, v := range e {

		if str, ok := v.(String); ok {
			visitor.b.WriteString(string(str))
		} else {
			err := v.Accept(visitor)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// VisitSelect implements Visitor.
func (visitor *Sqlite3) VisitSelect(e *SelectStmnt) error {

	newCtx := SqlContext{
		This:      visitor.s.This,
		That:      visitor.s.That,
		Params:    make(map[string]Expression),
		Resolving: visitor.s.Resolving,
	}

	for k, v := range visitor.s.Params {
		newCtx.Params[k] = v
	}

	oldCtx := visitor.s
	visitor.s = &newCtx
	defer func() { visitor.s = oldCtx }()

	visitor.b.WriteString("SELECT ")

	if e.Distinct {
		visitor.b.WriteString("DISTINCT ")
	}

	colCount := 0

	newCtx.This = e.From

	anonTablePresent := false
	var remainingExcepts []ColumnIdentifier

	if len(e.SelectColumns)+len(e.IndexColumns)+len(e.PkColumns) == 0 {

		exceptAll := make(map[string]struct{})

		for _, v := range e.ExceptColumns {
			if v.Table != "" && v.Name == "*" {
				exceptAll[v.Table] = struct{}{}
			} else {
				remainingExcepts = append(remainingExcepts, v)
			}
		}

		if len(exceptAll) == 0 {
			colCount++
			visitor.b.WriteString("*")
		} else {

			var tableNames []string

			if name, ok := GetTableName(e.From); ok {
				tableNames = append(tableNames, name.Table)
			} else {
				anonTablePresent = true
			}

			for _, v := range e.Join {
				if name, ok := GetTableName(v.JoinTable); ok {
					tableNames = append(tableNames, name.Table)
				} else {
					anonTablePresent = true
					break
				}
			}

			if anonTablePresent {
				colCount++
				visitor.b.WriteString("*")
				remainingExcepts = e.ExceptColumns
			} else {
				for _, v := range tableNames {
					if _, except := exceptAll[v]; !except {
						visitor.EscapeIdentifier(v)
						visitor.b.WriteString(".*")
					}
				}
			}
		}

	} else {

		remainingExcepts = e.ExceptColumns

		for _, column := range e.PkColumns {

			if colCount != 0 {
				visitor.b.WriteString(",")
			}
			colCount++

			err := column.Accept(visitor)
			if err != nil {
				return err
			}
		}

		for _, column := range e.IndexColumns {
			if colCount != 0 {
				visitor.b.WriteString(",")
			}
			colCount++

			err := column.Accept(visitor)
			if err != nil {
				return err
			}
		}

		for _, column := range e.SelectColumns {

			if colCount != 0 {
				visitor.b.WriteString(",")
			}

			colCount++
			err := column.Accept(visitor)
			if err != nil {
				return err
			}
		}
	}

	for _, column := range e.DeriveColumns {

		if colCount != 0 {
			visitor.b.WriteString(",")
		}

		err := column.Accept(visitor)
		if err != nil {
			return err
		}
	}

	if len(remainingExcepts) > 0 {

		visitor.b.WriteString(" EXCEPT (")
		for i, column := range remainingExcepts {

			if i != 0 {
				visitor.b.WriteString(",")
			}

			err := column.Accept(visitor)
			if err != nil {
				return err
			}
		}
		visitor.b.WriteString(")")
	}

	visitor.b.WriteString(" FROM ")
	switch e.From.(type) {
	case *TableIdentifier:
		err := e.From.Accept(visitor)
		if err != nil {
			return err
		}
	case *TableDef:
		err := e.From.Accept(visitor)
		if err != nil {
			return err
		}
	default:
		visitor.b.WriteString("(")
		if e.From == nil {
			return errors.New("missing from")
		}
		err := e.From.Accept(visitor)
		if err != nil {
			return err
		}
		visitor.b.WriteString(")")
	}

	for _, join := range e.Join {
		visitor.b.WriteString(" ")
		err := join.Accept(visitor)
		visitor.s.That = nil
		if err != nil {
			return err
		}
	}

	if len(e.Where) > 0 {
		visitor.b.WriteString(" WHERE (")
		for i, v := range e.Where {
			if i != 0 {
				visitor.b.WriteString(") AND (")
			}

			err := v.Accept(visitor)
			if err != nil {
				return err
			}
		}
		visitor.b.WriteString(")")
	}

	if len(e.GroupBy) > 0 {
		visitor.b.WriteString(" GROUP BY ")
		for i, column := range e.GroupBy {

			colCount++
			if i != 0 {
				visitor.b.WriteString(",")
			}

			if def, isDef := column.(*ColumnDef); isDef {
				err := def.Expr.Accept(visitor)
				if err != nil {
					return err
				}
			} else if def, isDef := column.(*TableDef); isDef {
				err := def.Expr.Accept(visitor)
				if err != nil {
					return err
				}
			} else {
				err := column.Accept(visitor)
				if err != nil {
					return err
				}
			}
		}
	}

	if len(e.Having) > 0 {
		visitor.b.WriteString(" HAVING (")
		for i, v := range e.Having {
			if i != 0 {
				visitor.b.WriteString(") AND (")
			}

			err := v.Accept(visitor)
			if err != nil {
				return err
			}
		}
		visitor.b.WriteString(")")
	}

	if len(e.OrderBy) > 0 {
		visitor.b.WriteString(" ORDER BY ")
		for i, v := range e.OrderBy {
			if i != 0 {
				visitor.b.WriteString(",")
			}

			if colDef, isColDef := v.(*ColumnDef); isColDef {
				v = colDef.Expr
			}

			descExpr, isDesc := v.(*NegateExpr)
			if isDesc {
				v = descExpr.Expr
			}

			posExpr, isPos := v.(*Positive)
			if isPos {
				v = posExpr.Expr
			}

			err := v.Accept(visitor)
			if err != nil {
				return err
			}

			if isDesc {
				visitor.b.WriteString(" DESC")
			}

		}
	}

	if e.Limit != nil {

		visitor.b.WriteString(" LIMIT ")
		err := e.Limit.Accept(visitor)
		if err != nil {
			return err
		}
	}

	if e.Offset != nil {

		visitor.b.WriteString(" OFFSET ")
		err := e.Offset.Accept(visitor)
		if err != nil {
			return err
		}
	}

	return nil
}

func (visitor *Sqlite3) VisitDelete(e *DeleteStmnt) error {

	newCtx := SqlContext{
		This:      visitor.s.This,
		That:      visitor.s.That,
		Params:    make(map[string]Expression),
		Resolving: visitor.s.Resolving,
	}

	for k, v := range visitor.s.Params {
		newCtx.Params[k] = v
	}

	oldCtx := visitor.s
	visitor.s = &newCtx
	defer func() { visitor.s = oldCtx }()

	visitor.b.WriteString("DELETE FROM ")

	newCtx.This = e.From
	switch e.From.(type) {
	case *TableIdentifier:
		err := e.From.Accept(visitor)
		if err != nil {
			return err
		}
	case *TableDef:
		err := e.From.Accept(visitor)
		if err != nil {
			return err
		}
	default:
		visitor.b.WriteString("(")
		if e.From == nil {
			return errors.New("missing from")
		}
		err := e.From.Accept(visitor)
		if err != nil {
			return err
		}
		visitor.b.WriteString(")")
	}

	if len(e.Where) > 0 {
		visitor.b.WriteString(" WHERE (")
		for i, v := range e.Where {
			if i != 0 {
				visitor.b.WriteString(") AND (")
			}

			err := v.Accept(visitor)
			if err != nil {
				return err
			}
		}
		visitor.b.WriteString(")")
	}

	return nil
}

func (visitor *Sqlite3) VisitUpdate(e *UpdateStmnt) error {

	newCtx := SqlContext{
		This:      visitor.s.This,
		That:      visitor.s.That,
		Params:    make(map[string]Expression),
		Resolving: visitor.s.Resolving,
	}

	for k, v := range visitor.s.Params {
		newCtx.Params[k] = v
	}

	oldCtx := visitor.s
	visitor.s = &newCtx
	defer func() { visitor.s = oldCtx }()

	visitor.b.WriteString("UPDATE ")

	switch e.Table.(type) {
	case *TableIdentifier:
		err := e.Table.Accept(visitor)
		if err != nil {
			return err
		}
	case *TableDef:
		err := e.Table.Accept(visitor)
		if err != nil {
			return err
		}
	default:
		visitor.b.WriteString("(")
		if e.Table == nil {
			return errors.New("missing from")
		}
		err := e.Table.Accept(visitor)
		if err != nil {
			return err
		}
		visitor.b.WriteString(")")
	}

	visitor.b.WriteString(" SET ")

	colCount := 0

	newCtx.This = e.Table

	for i, column := range e.SetColumns {

		colCount++
		if i != 0 {
			visitor.b.WriteString(",")
		}

		visitor.b.WriteString(column.ColumnIdentifier.Name)
		visitor.b.WriteString(" = ")
		err := column.Value.Accept(visitor)
		if err != nil {
			return err
		}
	}

	if len(e.Where) > 0 {
		visitor.b.WriteString(" WHERE (")
		for i, v := range e.Where {
			if i != 0 {
				visitor.b.WriteString(") AND (")
			}

			err := v.Accept(visitor)
			if err != nil {
				return err
			}
		}
		visitor.b.WriteString(")")
	}

	return nil
}

func (visitor *Sqlite3) VisitInsert(e *InsertStmnt) error {

	newCtx := SqlContext{
		This:      visitor.s.This,
		That:      visitor.s.That,
		Params:    make(map[string]Expression),
		Resolving: visitor.s.Resolving,
	}

	for k, v := range visitor.s.Params {
		newCtx.Params[k] = v
	}

	oldCtx := visitor.s
	visitor.s = &newCtx
	defer func() { visitor.s = oldCtx }()

	visitor.b.WriteString("INSERT INTO ")

	switch e.Table.(type) {
	case *TableIdentifier:
		err := e.Table.Accept(visitor)
		if err != nil {
			return err
		}
	case *TableDef:
		err := e.Table.Accept(visitor)
		if err != nil {
			return err
		}
	default:
		visitor.b.WriteString("(")
		if e.Table == nil {
			return errors.New("missing from")
		}
		err := e.Table.Accept(visitor)
		if err != nil {
			return err
		}
		visitor.b.WriteString(")")
	}

	visitor.b.WriteString("(")

	colCount := 0

	newCtx.This = e.Table

	for i, column := range e.Values {

		colCount++
		if i != 0 {
			visitor.b.WriteString(",")
		}

		visitor.b.WriteString(column.ColumnIdentifier.Name)
	}

	visitor.b.WriteString(") VALUES(")

	for i, column := range e.Values {

		colCount++
		if i != 0 {
			visitor.b.WriteString(",")
		}

		column.Value.Accept(visitor)
	}

	visitor.b.WriteString(")")

	return nil
}

// VisitSetOperation implements Visitor.
func (visitor *Sqlite3) VisitSetOperation(e *SetOperation) error {

	err := e.Left.Accept(visitor)
	if err != nil {
		return err
	}

	visitor.b.WriteString(" ")
	visitor.b.WriteString(e.Operation)
	visitor.b.WriteString(" ")

	err = e.Right.Accept(visitor)
	if err != nil {
		return err
	}

	return nil
}

// VisitString implements Visitor.
func (visitor *Sqlite3) VisitString(e String) error {
	visitor.b.WriteString("'" + strings.ReplaceAll(string(e), "'", "''") + "'")
	return nil
}

// VisitTableDef implements Visitor.
func (visitor *Sqlite3) VisitTableDef(e *TableDef) error {

	switch e.Expr.(type) {
	case *TableIdentifier:
		err := e.Expr.Accept(visitor)
		if err != nil {
			return err
		}
	default:
		visitor.b.WriteString("(")
		err := e.Expr.Accept(visitor)
		if err != nil {
			return err
		}
		visitor.b.WriteString(")")
	}

	if e.Name != "" {
		visitor.b.WriteString(" AS ")
		visitor.EscapeIdentifier(e.Name)
	}

	return nil
}

// VisitTableIdentifier implements Visitor.
func (visitor *Sqlite3) VisitTableIdentifier(e *TableIdentifier) error {

	if e.Database == "" {
		if e.Schema == "" {
			expr, isExpr := visitor.s.Params[e.Table]
			if isExpr {
				return expr.Accept(visitor)
			}

			switch e.Table {
			case "this":
				if visitor.s.This == nil {
					return fmt.Errorf("no this defined")
				}
				err := visitor.s.This.Accept(visitor)
				if err != nil {
					return err
				}
			case "that":
				if visitor.s.That == nil {
					return fmt.Errorf("no that defined")
				}
				err := visitor.s.That.Accept(visitor)
				if err != nil {
					return err
				}
			default:
				visitor.EscapeIdentifier(e.Table)
			}
		} else {
			visitor.EscapeIdentifier(e.Schema)
			visitor.b.WriteString(".")
			visitor.EscapeIdentifier(e.Table)
		}
	} else {
		visitor.EscapeIdentifier(e.Database)
		visitor.b.WriteString(".")
		visitor.EscapeIdentifier(e.Schema)
		visitor.b.WriteString(".")
		visitor.EscapeIdentifier(e.Table)
	}

	if e.TableAlias != "" {
		visitor.b.WriteString(" AS ")
		visitor.EscapeIdentifier(e.TableAlias)
	}

	return nil
}

// VisitValue implements Visitor.
func (visitor *Sqlite3) VisitValue(e ValueExpr) error {

	visitor.b.WriteString("(")

	for i, v := range e {
		if i != 0 {
			visitor.b.WriteString(",")
		}

		err := v.Accept(visitor)
		if err != nil {
			return err
		}
	}

	visitor.b.WriteString(")")

	return nil
}

// VisitValues implements Visitor.
func (visitor *Sqlite3) VisitValues(e *Values) error {

	visitor.b.WriteString("(values ")

	for i, v := range *e {
		if i != 0 {
			visitor.b.WriteString(",")
		}

		err := v.Accept(visitor)
		if err != nil {
			return err
		}
	}

	visitor.b.WriteString(")")

	return nil

}
