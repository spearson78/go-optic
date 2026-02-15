package osql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

func getReflectColumn(tbl *TableIdentifier, t reflect.StructField) (*ColumnIdentifier, error) {
	tag := t.Tag.Get("db")

	if tag == "" {
		return nil, fmt.Errorf("missing db column name %v", t.Name)
	}

	return &ColumnIdentifier{
		TableIdentifier: *tbl,
		Name:            tag,
	}, nil
}

type sqlContext struct {
	ActiveTable *TableIdentifier
}

func getColumn(tbl *TableIdentifier, o expr.OpticExpression) (*ColumnIdentifier, error) {
	switch t := o.(type) {
	case expr.Compose:
		//This is a side effect of makelens
		if _, isId := t.Left.(expr.Identity); !isId {
			return nil, fmt.Errorf("unknown column type %T", t.Left)
		}

		return getColumn(tbl, t.Right)
	case expr.FieldLens:
		return getReflectColumn(
			tbl,
			t.Field,
		)
	default:
		return nil, fmt.Errorf("unknown column type %T", t)
	}
}

func getReflectPrimaryKey(tbl *TableIdentifier) (*ColumnIdentifier, reflect.StructField, error) {

	for i := 0; i < tbl.Type.NumField(); i++ {
		f := tbl.Type.Field(i)
		tag := f.Tag.Get("osql")
		if tag == "PK" {
			col, err := getReflectColumn(
				tbl,
				f,
			)

			return col, f, err
		}
	}

	return nil, reflect.StructField{}, fmt.Errorf("missing db pk column %v", tbl.Type.Name())

}

func isFilterMode(e expr.Filtered) bool {
	return e.NoMatchMode == expr.FilterContinue && e.MatchMode == expr.FilterContinue
}

func toFloat(v any) (float64, error) {
	switch t := v.(type) {
	case int:
		return float64(t), nil
	case int32:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case float32:
		return float64(t), nil
	case float64:
		return t, nil
	case string:
		return strconv.ParseFloat(t, 64)
	default:
		return 0, fmt.Errorf("cannot coerce to number : %T", t)
	}
}

func toInt(v any) (int, error) {
	switch t := v.(type) {
	case int:
		return int(t), nil
	case int32:
		return int(t), nil
	case int64:
		return int(t), nil
	case float32:
		return int(t), nil
	case float64:
		return int(t), nil
	case string:
		i, err := strconv.ParseInt(t, 10, 0)
		return int(i), err
	default:
		return 0, fmt.Errorf("cannot coerce to number : %T", t)
	}
}

func buildExpression(sqlContext *sqlContext, left Expression, o expr.OpticExpression) (Expression, error) {
	switch t := o.(type) {
	case expr.OpT2BindExpr:

		op := t.Op.(expr.BinaryExpr).Op

		val, err := toLiteral(t.RightValue)
		if err != nil {
			return nil, err
		}

		return &BinaryExpression{
			Left:  left,
			Op:    op,
			Right: val,
		}, nil

	case expr.Identity:
		return nil, nil
	case expr.FieldLens:
		if left != nil {
			return nil, fmt.Errorf("left is non nil on field access : %T", left)
		}

		if strings.HasPrefix(t.OpticS().Name(), "ValueI[") && t.Field.Name == "value" {
			//Ignore the value field lens
			return nil, nil
			//TODO: Handle the ValueI index field lens.
		}

		return getColumn(sqlContext.ActiveTable, t)
	case expr.Compose:

		whereExp, err := buildExpression(sqlContext, left, t.Left)
		if err != nil {
			return nil, err
		}

		return buildExpression(sqlContext, whereExp, t.Right)
	case expr.IsoOpT2BindExpr:
		switch op := t.Op.(type) {
		case expr.BinaryExpr:

			val, err := toLiteral(t.Right)
			if err != nil {
				return nil, err
			}

			opExpr := &BinaryExpression{
				Left:  left,
				Op:    op.Op,
				Right: val,
			}

			return opExpr, nil

		default:
			return nil, fmt.Errorf("buildExpression: unknown op expression type %T", op)
		}
	default:
		return nil, fmt.Errorf("buildWhere: unknown expression type %T", t)
	}
}

func buildSelect(sqlContext *sqlContext, stmnt *SelectStmnt, o expr.OpticExpression) error {

	switch t := o.(type) {
	case expr.Identity:
		return nil
	case TableExpr:

		tbl := &TableIdentifier{
			Table: t.Name,
			Type:  t.A,
		}

		sqlContext.ActiveTable = tbl

		pkCol, pkField, err := getReflectPrimaryKey(tbl)
		if err != nil {
			return err
		}

		alias := &Alias{
			Column: pkCol,
			Alias:  "OSQL_IX",
		}

		stmnt.From = tbl
		stmnt.IndexColumns = []Expression{
			alias,
		}
		stmnt.SelectColumns = []Expression{
			&ColumnIdentifier{
				TableIdentifier: *tbl,
				Name:            "*",
			},
		}

		stmnt.Scan = RowScan{
			IxScan: func() (reflect.Value, []any, error) {
				ret := reflect.New(pkField.Type)
				return ret, []any{ret.Interface()}, nil
			},
			FocusScan: func() (reflect.Value, []any, error) {
				return ScanStruct(t.A)
			},
		}

		return nil
	case expr.Traverse:

		return nil
	case JoinMExpr:

		tbl := &TableIdentifier{
			Table: t.TableName,
			Type:  t.OpticA(),
		}

		prevTable := sqlContext.ActiveTable
		sqlContext.ActiveTable = tbl

		pkCol, pkField, err := getReflectPrimaryKey(tbl)
		if err != nil {
			return err
		}
		alias := &Alias{
			Column: pkCol,
			Alias:  "OSQL_IX",
		}

		stmnt.IndexColumns = []Expression{
			alias,
		}
		stmnt.SelectColumns = []Expression{
			&ColumnIdentifier{
				TableIdentifier: *tbl,
				Name:            "*",
			},
		}

		prevPkCol, _, err := getReflectPrimaryKey(prevTable)
		if err != nil {
			return err
		}

		stmnt.Join = append(stmnt.Join, JoinDef{
			Side:      JOIN_INNER,
			JoinTable: tbl,
			JoinExpr: &BinaryExpression{
				Left: prevPkCol,
				Op:   "=",
				Right: &ColumnIdentifier{
					TableIdentifier: *tbl,
					Name:            t.ColumnName,
				},
			},
		})

		stmnt.Scan = RowScan{
			IxScan: func() (reflect.Value, []any, error) {
				ret := reflect.New(pkField.Type)
				return ret, []any{ret.Interface()}, nil
			},
			FocusScan: func() (reflect.Value, []any, error) {
				return ScanStruct(t.OpticA())
			},
		}

		return nil

	case JoinExpr:

		tbl := &TableIdentifier{
			Table: t.TableName,
			Type:  t.OpticA(),
		}

		prevTable := sqlContext.ActiveTable

		sqlContext.ActiveTable = tbl

		pkCol, pkField, err := getReflectPrimaryKey(tbl)
		if err != nil {
			return err
		}
		alias := &Alias{
			Column: pkCol,
			Alias:  "OSQL_IX",
		}

		stmnt.IndexColumns = []Expression{
			alias,
		}
		stmnt.SelectColumns = []Expression{
			&ColumnIdentifier{
				TableIdentifier: *tbl,
				Name:            "*",
			},
		}

		stmnt.Join = append(stmnt.Join, JoinDef{
			Side:      JOIN_INNER,
			JoinTable: tbl,
			JoinExpr: &BinaryExpression{
				Left: &ColumnIdentifier{
					TableIdentifier: *prevTable,
					Name:            t.ColumnName,
				},
				Op:    "=",
				Right: pkCol,
			},
		})

		stmnt.Scan = RowScan{
			IxScan: func() (reflect.Value, []any, error) {
				ret := reflect.New(pkField.Type)
				return ret, []any{ret.Interface()}, nil
			},
			FocusScan: func() (reflect.Value, []any, error) {
				return ScanStruct(t.OpticA())
			},
		}

		return nil
	case expr.Filtered:
		if !isFilterMode(t) {
			return errors.New("unsupported filter mode")
		}

		if err := buildSelect(sqlContext, stmnt, t.Optic); err != nil {
			return err
		}

		expr, err := buildExpression(sqlContext, nil, t.Pred)
		if err != nil {
			return err
		}

		stmnt.Where = append(stmnt.Where, expr)

		return nil
	case expr.Compose:

		switch ixt := t.IxMap.(type) {
		case expr.IxMap:
			switch ixt.Type {
			case expr.IxMapperLeft:
				if err := buildSelect(sqlContext, stmnt, t.Left); err != nil {
					return err
				}
				leftIxCols := stmnt.IndexColumns
				leftIxScan := stmnt.Scan.IxScan

				err := buildSelect(sqlContext, stmnt, t.Right)
				if err != nil {
					return err
				}

				stmnt.IndexColumns = leftIxCols
				stmnt.Scan.IxScan = leftIxScan

				return nil
			case expr.IxMapperRight:
				if err := buildSelect(sqlContext, stmnt, t.Left); err != nil {
					return err
				}

				return buildSelect(sqlContext, stmnt, t.Right)
			default:
				return fmt.Errorf("unknown ixmap type %v", ixt.Type)
			}
		default:
			return fmt.Errorf("unknown ixmap type %T", ixt)
		}
	case expr.FieldLens:

		if t.OpticA().Kind() == reflect.Slice {
			return nil
		}

		col, err := getColumn(sqlContext.ActiveTable, o)
		if err != nil {
			return err
		}

		stmnt.IndexColumns = nil
		stmnt.SelectColumns = []Expression{col}

		stmnt.Scan = RowScan{
			IxScan: func() (reflect.Value, []any, error) {
				return reflect.ValueOf(&optic.Void{}), nil, nil
			},
			FocusScan: func() (reflect.Value, []any, error) {
				value := reflect.New(t.OpticA())
				return value, []any{value.Interface()}, nil
			},
		}

		return nil

	case expr.TupleOf:

		stmnt.IndexColumns = make([]Expression, 0, len(t.Elements))
		stmnt.SelectColumns = make([]Expression, 0, len(t.Elements))

		for _, element := range t.Elements {

			elementSelec := SelectStmnt{
				From: stmnt.From,
				Join: append([]JoinDef(nil), stmnt.Join...),
			}

			buildSelect(sqlContext, &elementSelec, element)

			stmnt.IndexColumns = append(stmnt.IndexColumns, elementSelec.IndexColumns...)
			stmnt.SelectColumns = append(stmnt.SelectColumns, elementSelec.SelectColumns...)
		}

		stmnt.Scan = RowScan{
			IxScan: func() (reflect.Value, []any, error) {
				tuplePtr := reflect.New(o.OpticI())

				var args []any
				for i := range tuplePtr.Elem().NumField() {
					var err error
					args, err = ScanInto(tuplePtr.Elem().Field(i).Addr(), args)
					if err != nil {
						return reflect.Value{}, nil, err
					}
				}

				return tuplePtr, args, nil
			},
			FocusScan: func() (reflect.Value, []any, error) {
				tuplePtr := reflect.New(o.OpticA())

				var args []any
				for i := range tuplePtr.Elem().NumField() {
					var err error
					args, err = ScanInto(tuplePtr.Elem().Field(i).Addr(), args)
					if err != nil {
						return reflect.Value{}, nil, err
					}
				}

				return tuplePtr, args, nil
			},
		}

		return nil
	case expr.SelfIndex:

		err := buildSelect(sqlContext, stmnt, t.Optic)
		if err != nil {
			return err
		}

		stmnt.IndexColumns = append([]Expression(nil), stmnt.SelectColumns...)
		stmnt.Scan.IxScan = stmnt.Scan.FocusScan

		return nil
	case expr.ReIndexed:
		err := buildSelect(sqlContext, stmnt, t.Optic)
		if err != nil {
			return err
		}

		if len(stmnt.IndexColumns) != 1 {
			return errors.New("ReIndexed: only primitive indexes can be mapped")
		}

		newIndex, err := buildExpression(sqlContext, stmnt.IndexColumns[0].(*Alias).Column, t.IxMap)
		if err != nil {
			return err
		}

		stmnt.IndexColumns[0] = &Alias{
			Column: newIndex,
			Alias:  "OSQL_IX",
		}

		stmnt.Scan.IxScan = func() (reflect.Value, []any, error) {
			v := reflect.New(t.IxMap.OpticA())
			args, err := ScanInto(v, nil)
			return v, args, err
		}

		return nil
	case expr.Index:
		err := buildSelect(sqlContext, stmnt, t.Optic)
		if err != nil {
			return err
		}

		if len(stmnt.IndexColumns) != 1 {
			return errors.New("Index: only primitive indexes are supported")
		}

		val, err := toLiteral(t.Index)
		if err != nil {
			return err
		}

		opExpr := &BinaryExpression{
			Left:  &ColumnIdentifier{Name: "OSQL_IX"},
			Op:    "=",
			Right: val,
		}

		stmnt.Where = append(stmnt.Where, opExpr)

		return nil

	default:
		return fmt.Errorf("buildSelect: unknown expression type %T", t)
	}

}

type RowScan struct {
	IxScan    func() (reflect.Value, []any, error)
	FocusScan func() (reflect.Value, []any, error)
}

func ScanStruct(structType reflect.Type) (reflect.Value, []any, error) {
	structPtr := reflect.New(structType)
	args, err := ScanInto(structPtr, nil)

	return structPtr, args, err
}

func ScanInto(structPtr reflect.Value, args []any) ([]any, error) {
	if structPtr.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("target must be a pointer, got %s", structPtr.Kind())
	}

	switch structPtr.Elem().Kind() {
	case reflect.Struct:
		for i := 0; i < structPtr.Elem().NumField(); i++ {

			fieldValue := structPtr.Elem().Field(i)

			writeableField := reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr()))

			args = append(args, writeableField.Interface())
		}

		return args, nil
	case reflect.String, reflect.Int:
		return append(args, structPtr.Interface()), nil
	default:
		return nil, fmt.Errorf("target must be a struct, got %s", structPtr.Elem().Kind())
	}
}

func setField(structVal reflect.Value, i int, setVal reflect.Value) {
	field := structVal.Field(i)
	writeableField := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	writeableField.Set(setVal)
}

func execGet(ctx context.Context, o expr.OpticExpression, db *sql.DB) (any, any, error) {

	var sqlContext sqlContext

	switch t := o.(type) {
	case expr.SeqOf:
		var stmnt SelectStmnt
		err := buildSelect(&sqlContext, &stmnt, t.Optic)
		if err != nil {
			return nil, nil, err
		}

		sqlBuilder := Sqlite3{
			s: &SqlContext{},
		}

		sqlBuilder.VisitSelect(&stmnt)

		sqlStr := sqlBuilder.b.String()

		name := t.OpticA().Name()
		switch {
		case strings.HasPrefix(name, "SeqE["):

			iterFnc := reflect.MakeFunc(t.OpticA(), func(args []reflect.Value) (results []reflect.Value) {

				yieldFnc := args[0]
				valType := yieldFnc.Type().In(0)

				row, err := db.QueryContext(ctx, sqlStr)
				if err != nil {
					val := reflect.New(valType).Elem()
					setField(val, 0, reflect.Zero(t.A))
					setField(val, 1, reflect.ValueOf(err))

					yieldFnc.Call([]reflect.Value{val})
					return
				}
				defer row.Close()

				for row.Next() {
					var args []any
					_, ixArgs, err := stmnt.Scan.IxScan()
					if err != nil {
						val := reflect.New(valType).Elem()
						setField(val, 0, reflect.Zero(t.A))
						setField(val, 1, reflect.ValueOf(err))

						yieldFnc.Call([]reflect.Value{val})
						return
					}
					args = append(args, ixArgs...)

					focusVal, focusArgs, err := stmnt.Scan.FocusScan()
					if err != nil {
						val := reflect.New(valType).Elem()
						setField(val, 0, reflect.Zero(t.A))
						setField(val, 1, reflect.ValueOf(err))

						yieldFnc.Call([]reflect.Value{val})
						return
					}
					args = append(args, focusArgs...)

					err = row.Scan(args...)
					if err != nil {
						val := reflect.New(valType).Elem()
						setField(val, 0, reflect.Zero(t.A))
						setField(val, 1, reflect.ValueOf(err))

						yieldFnc.Call([]reflect.Value{val})
						return
					}

					val := reflect.New(valType).Elem()

					setField(val, 0, focusVal.Elem())

					ret := yieldFnc.Call([]reflect.Value{val})
					if !ret[0].Bool() {
						break
					}

				}

				return nil
			})

			return optic.Void{}, iterFnc.Interface(), nil
		case strings.HasPrefix(name, "SeqIE["):

			iterFnc := reflect.MakeFunc(t.OpticA(), func(args []reflect.Value) (results []reflect.Value) {

				yieldFnc := args[0]
				valType := yieldFnc.Type().In(0)

				row, err := db.QueryContext(ctx, sqlStr)
				if err != nil {
					val := reflect.New(valType).Elem()
					setField(val, 0, reflect.Zero(t.I))
					setField(val, 1, reflect.Zero(t.A))
					setField(val, 2, reflect.ValueOf(err))

					yieldFnc.Call([]reflect.Value{val})
					return
				}
				defer row.Close()

				for row.Next() {
					var args []any
					ixVal, ixArgs, err := stmnt.Scan.IxScan()
					if err != nil {
						val := reflect.New(valType).Elem()
						setField(val, 0, reflect.Zero(t.I))
						setField(val, 1, reflect.Zero(t.A))
						setField(val, 2, reflect.ValueOf(err))

						yieldFnc.Call([]reflect.Value{val})
						return
					}
					args = append(args, ixArgs...)

					focusVal, focusArgs, err := stmnt.Scan.FocusScan()
					if err != nil {
						val := reflect.New(valType).Elem()
						setField(val, 0, reflect.Zero(t.I))
						setField(val, 1, reflect.Zero(t.A))
						setField(val, 2, reflect.ValueOf(err))

						yieldFnc.Call([]reflect.Value{val})
						return
					}
					args = append(args, focusArgs...)

					err = row.Scan(args...)
					if err != nil {
						val := reflect.New(valType).Elem()
						setField(val, 0, reflect.Zero(t.I))
						setField(val, 1, reflect.Zero(t.A))
						setField(val, 2, reflect.ValueOf(err))

						yieldFnc.Call([]reflect.Value{val})
						return
					}

					val := reflect.New(valType).Elem()

					setField(val, 0, ixVal.Elem())
					setField(val, 1, focusVal.Elem())
					setField(val, 2, reflect.Zero(reflect.TypeFor[error]()))

					ret := yieldFnc.Call([]reflect.Value{val})
					if !ret[0].Bool() {
						break
					}

				}

				return nil
			})

			return optic.Void{}, iterFnc.Interface(), nil

		default:
			return nil, nil, fmt.Errorf("unknown iter func type %v", t.OpticA().Name())
		}

	case expr.SliceOf:
		slice := reflect.New(reflect.SliceOf(t.A)).Elem()

		var stmnt SelectStmnt
		err := buildSelect(&sqlContext, &stmnt, t.Optic)
		if err != nil {
			return nil, slice.Interface(), err
		}

		sqlBuilder := Sqlite3{
			s: &SqlContext{},
		}

		sqlBuilder.VisitSelect(&stmnt)

		sql := sqlBuilder.b.String()

		row, err := db.QueryContext(ctx, sql)
		if err != nil {
			return nil, slice.Interface(), err
		}
		defer row.Close()

		for row.Next() {
			var args []any
			_, ixArgs, err := stmnt.Scan.IxScan()
			if err != nil {
				return nil, slice.Interface(), err
			}
			args = append(args, ixArgs...)

			focusVal, focusArgs, err := stmnt.Scan.FocusScan()
			if err != nil {
				return nil, slice.Interface(), err
			}
			args = append(args, focusArgs...)

			err = row.Scan(args...)
			if err != nil {
				return nil, slice.Interface(), err
			}

			slice = reflect.Append(slice, focusVal.Elem())

		}

		return optic.Void{}, slice.Interface(), err
	case expr.MapOfReduced:
		retMap := reflect.MakeMap(reflect.MapOf(t.Optic.OpticI(), t.Optic.OpticA()))

		var stmnt SelectStmnt
		err := buildSelect(&sqlContext, &stmnt, t.Optic)
		if err != nil {
			return nil, retMap.Interface(), err
		}

		sqlBuilder := Sqlite3{
			s: &SqlContext{},
		}

		sqlBuilder.VisitSelect(&stmnt)

		sql := sqlBuilder.b.String()

		row, err := db.QueryContext(ctx, sql)
		if err != nil {
			return nil, retMap.Interface(), err
		}
		defer row.Close()

		for row.Next() {
			var args []any
			ixVal, ixArgs, err := stmnt.Scan.IxScan()
			if err != nil {
				return nil, retMap.Interface(), err
			}
			args = append(args, ixArgs...)

			focusVal, focusArgs, err := stmnt.Scan.FocusScan()
			if err != nil {
				return nil, retMap.Interface(), err
			}
			args = append(args, focusArgs...)

			err = row.Scan(args...)
			if err != nil {
				return nil, retMap.Interface(), err
			}

			retMap.SetMapIndex(ixVal.Elem(), focusVal.Elem())
		}

		return optic.Void{}, retMap.Interface(), err

	default:
		var stmnt SelectStmnt
		err := buildSelect(&sqlContext, &stmnt, o)
		if err != nil {
			return nil, nil, err
		}

		if stmnt.Limit == nil {
			stmnt.Limit = Number(1)
		}

		sqlBuilder := Sqlite3{
			s: &SqlContext{},
		}

		sqlBuilder.VisitSelect(&stmnt)

		sqlStr := sqlBuilder.b.String()

		var args []any
		ixVal, ixArgs, err := stmnt.Scan.IxScan()
		if err != nil {
			return nil, nil, err
		}
		args = append(args, ixArgs...)

		focusVal, focusArgs, err := stmnt.Scan.FocusScan()
		if err != nil {
			return nil, nil, err
		}
		args = append(args, focusArgs...)

		row := db.QueryRowContext(ctx, sqlStr)
		err = row.Scan(args...)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, optic.ErrEmptyGet
		}
		if err != nil {
			return nil, nil, err
		}

		return ixVal.Elem().Interface(), focusVal.Elem().Interface(), nil
	}

}

func buildDelete(sqlContext *sqlContext, stmnt *DeleteStmnt, o expr.OpticExpression) error {
	switch t := o.(type) {
	case TableExpr:

		tbl := &TableIdentifier{
			Table: t.Name,
			Type:  t.A,
		}

		stmnt.From = tbl

		sqlContext.ActiveTable = tbl

		return nil
	default:
		return fmt.Errorf("buildDelete: unknown expression type %T", t)
	}
}

func buildDeleteFmap(sqlContext *sqlContext, stmnt *DeleteStmnt, fmapExpr expr.OpticExpression) error {
	switch t := fmapExpr.(type) {
	case expr.Filtered:
		if !isFilterMode(t) {
			return errors.New("unsupported filter mode")
		}

		expr, err := buildExpression(sqlContext, nil, t.Pred)
		if err != nil {
			return err
		}

		stmnt.Where = append(stmnt.Where, &SqlNot{Expr: expr})

		return nil
	default:
		return fmt.Errorf("buildDeleteFmap: unknown expression type %T", t)
	}
}

func buildUpdate(sqlContext *sqlContext, stmnt *UpdateStmnt, o expr.OpticExpression) error {
	switch t := o.(type) {
	case expr.Compose:
		err := buildUpdate(sqlContext, stmnt, t.Left)
		if err != nil {
			return err
		}
		return buildUpdate(sqlContext, stmnt, t.Right)
	case expr.Filtered:
		if !isFilterMode(t) {
			return errors.New("unsupported filter mode")
		}

		err := buildUpdate(sqlContext, stmnt, t.Optic)
		if err != nil {
			return err
		}

		expr, err := buildExpression(sqlContext, nil, t.Pred)
		if err != nil {
			return err
		}

		stmnt.Where = append(stmnt.Where, expr)

		return nil
	case expr.Traverse:
		return nil
	case expr.Identity:
		return nil
	case TableExpr:

		tbl := &TableIdentifier{
			Table: t.Name,
			Type:  t.A,
		}

		sqlContext.ActiveTable = tbl
		stmnt.Table = tbl
		stmnt.SetColumns = nil

		for i := 0; i < t.A.NumField(); i++ {
			field := t.A.Field(i)

			col, err := getReflectColumn(tbl, field)
			if err != nil {
				return err
			}

			stmnt.SetColumns = append(stmnt.SetColumns, SetColumn{
				ColumnIdentifier: *col,
			})
		}

		return nil
	case expr.FieldLens:

		col, err := getColumn(sqlContext.ActiveTable, t)
		if err != nil {
			return err
		}

		stmnt.SetColumns = []SetColumn{
			{
				ColumnIdentifier: *col,
			},
		}

		return nil
	case expr.Index:
		err := buildUpdate(sqlContext, stmnt, t.Optic)
		if err != nil {
			return err
		}

		val, err := toLiteral(t.Index)
		if err != nil {
			return err
		}

		col, _, err := getReflectPrimaryKey(sqlContext.ActiveTable)
		if err != nil {
			return err
		}

		opExpr := &BinaryExpression{
			Left:  col,
			Op:    "=",
			Right: val,
		}

		stmnt.Where = append(stmnt.Where, opExpr)

		return nil
	default:
		return fmt.Errorf("buildUpdate: unknown expression type %T", t)
	}
}

func toLiteral(t any) (Expression, error) {
	return toLiteralReflect(reflect.ValueOf(t), t)
}

func toLiteralReflect(val reflect.Value, t any) (Expression, error) {

	kind := val.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		val, err := toFloat(t)
		if err != nil {
			return nil, err
		}
		return Number(val), nil
	case reflect.Float32, reflect.Float64:
		val, err := toFloat(t)
		if err != nil {
			return nil, err
		}
		return Number(val), nil

	case reflect.Bool:
		if t.(bool) {
			return Number(1), nil
		} else {
			return Number(0), nil
		}
	case reflect.String:
		return String(t.(string)), nil
	default:
		return nil, fmt.Errorf("unknown literal kind %v", kind)
	}
}

func buildUpdateFmap(stmnt *UpdateStmnt, fmapExpr expr.OpticExpression) error {
	switch t := fmapExpr.(type) {
	case expr.IsoOpT2BindExpr:
		switch op := t.Op.(type) {
		case expr.BinaryExpr:
			if len(stmnt.SetColumns) == 0 {
				return errors.New("missing column focus in update")
			}

			val, err := toLiteral(t.Right)
			if err != nil {
				return err
			}

			opExpr := &BinaryExpression{
				Left:  &stmnt.SetColumns[0].ColumnIdentifier,
				Op:    op.Op,
				Right: val,
			}

			stmnt.SetColumns[0].Value = opExpr

			return nil

		default:
			return fmt.Errorf("buildUpdateFmap: unknown op expression type %T", op)
		}
	case expr.Const:

		val := reflect.ValueOf(t.Value)
		kind := val.Kind()
		switch kind {
		case reflect.Struct:
			for i := 0; i < val.NumField(); i++ {
				fieldVal := val.Field(i)

				litVal, err := toLiteralReflect(fieldVal, fieldVal.Interface())
				if err != nil {
					return err
				}

				stmnt.SetColumns[i].Value = litVal
			}

			return nil

		default:
			val, err := toLiteral(t.Value)
			if err != nil {
				return err
			}
			stmnt.SetColumns[0].Value = val
			return nil
		}
	default:
		return fmt.Errorf("buildUpdateFmap: unknown expression type %T", t)
	}
}

func buildInsert(sqlContext *sqlContext, stmnt *InsertStmnt, o expr.OpticExpression) error {
	switch t := o.(type) {
	case TableExpr:
		tbl := &TableIdentifier{
			Table: t.Name,
			Type:  t.A,
		}

		sqlContext.ActiveTable = tbl
		stmnt.Table = tbl

		for i := 0; i < t.A.NumField(); i++ {
			field := t.A.Field(i)

			if field.Tag.Get("osql") != "PK" {

				col, err := getReflectColumn(tbl, field)
				if err != nil {
					return err
				}

				stmnt.Values = append(stmnt.Values, SetColumn{
					ColumnIdentifier: *col,
					Value:            Param("P" + strconv.Itoa(i)),
				})
			}
		}

		return nil
	default:
		return fmt.Errorf("buildInsert: unknown expression type %T", t)
	}
}

func execModify(ctx context.Context, o expr.OpticExpression, fmapExpr expr.OpticExpression, db *sql.DB) (any, error) {

	sqlContext := &sqlContext{}

	switch tfmap := fmapExpr.(type) {
	case expr.CollectionOf:
		switch t := tfmap.Optic.(type) {
		case expr.Filtered:
			var stmnt DeleteStmnt
			err := buildDelete(sqlContext, &stmnt, o)
			if err != nil {
				return nil, err
			}

			err = buildDeleteFmap(sqlContext, &stmnt, tfmap.Optic)
			if err != nil {
				return nil, err
			}

			sqlBuilder := Sqlite3{
				s: &SqlContext{},
			}

			sqlBuilder.VisitDelete(&stmnt)

			sql := sqlBuilder.b.String()

			_, err = db.ExecContext(ctx, sql)

			return db, err
		default:
			return nil, fmt.Errorf("unknown collection operation %T", t)
		}
	case expr.Compose:
		switch compLeft := tfmap.Left.(type) {
		case expr.FieldLens:

			if strings.HasPrefix(compLeft.OpticS().Name(), "ValueI[") && compLeft.Field.Name == "value" {
				//skip the ValueI.value field lens
				return execModify(ctx, o, tfmap.Right, db)
			}

			return nil, errors.New("not implemented.")

		case expr.TupleOf:
			switch tfmap.Right.(type) {
			case expr.AppendCol:

				if len(compLeft.Elements) != 2 {
					return nil, fmt.Errorf("unknown appendcol left expression %T", compLeft)
				}

				constLeft, ok := compLeft.Elements[1].(expr.Const)
				if !ok {
					return nil, fmt.Errorf("unknown appendcol tuple[1] expression %T", compLeft.Elements[1])
				}

				asSeqExpr, ok := constLeft.Value.(expr.AsSeqExpr)
				if !ok {
					return nil, fmt.Errorf("unknown appendcol const expression %T", compLeft.Elements[1])
				}

				var stmnt InsertStmnt
				err := buildInsert(sqlContext, &stmnt, o)
				if err != nil {
					return nil, err
				}

				if len(stmnt.Values) == 0 {
					return nil, errors.New("execModify: insert missing values")
				}

				sqlBuilder := Sqlite3{
					s: &SqlContext{},
				}

				sqlBuilder.VisitInsert(&stmnt)

				sqlStr := sqlBuilder.b.String()

				tx, err := db.BeginTx(ctx, nil)
				if err != nil {
					return nil, err
				}

				var retErr error

				asSeqExpr.AsExpr()(ctx, func(seqExprVal expr.SeqExprValue) bool {
					if seqExprVal.Error != nil {
						retErr = seqExprVal.Error
					}

					val := reflect.ValueOf(seqExprVal.ValuePtr).Elem()

					var args []any
					for i := 0; i < val.NumField(); i++ {
						field := val.Field(i)

						if val.Type().Field(i).Tag.Get("osql") != "PK" {
							writeableField := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr()))
							args = append(args, sql.Named("P"+strconv.Itoa(i), writeableField.Interface()))
						}
					}

					_, err = tx.ExecContext(ctx, sqlStr, args...)
					if err != nil {
						retErr = err
						return false
					}

					return true
				})

				if retErr != nil {
					tx.Rollback()
					return nil, err
				}

				err = tx.Commit()

				return db, err

			default:
				return nil, fmt.Errorf("unknown compose right expression %T", tfmap)
			}

		default:
			return nil, fmt.Errorf("unknown compose left expression %T", tfmap)
		}

	default:

		var stmnt UpdateStmnt

		err := buildUpdate(sqlContext, &stmnt, o)

		if err != nil {
			return nil, err
		}

		err = buildUpdateFmap(&stmnt, fmapExpr)
		if err != nil {
			return nil, err
		}

		sqlBuilder := Sqlite3{
			s: &SqlContext{},
		}

		sqlBuilder.VisitUpdate(&stmnt)

		sql := sqlBuilder.b.String()

		_, err = db.Exec(sql)

		return db, err

	}

}
