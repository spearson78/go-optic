package osql

type Visitor interface {
	VisitSelect(expr *SelectStmnt) error
	VisitUpdate(expr *UpdateStmnt) error
	VisitInsert(expr *InsertStmnt) error
	VisitDelete(expr *DeleteStmnt) error
	VisitTableIdentifier(expr *TableIdentifier) error
	VisitColumnIdentifier(expr *ColumnIdentifier) error
	VisitAlias(expr *Alias) error
	VisitColumnExpression(expr *ColumnExpression) error
	VisitKeyword(expr Keyword) error
	VisitTableDef(expr *TableDef) error
	VisitColumnDef(expr *ColumnDef) error
	VisitList(expr List) error
	VisitNumber(expr Number) error
	VisitString(expr String) error
	VisitDate(expr Date) error
	VisitBinaryExpression(expr *BinaryExpression) error
	VisitDuration(expr Duration) error
	VisitCall(expr *Call) error
	VisitJoin(expr *JoinDef) error
	VisitExists(expr *Exists) error
	VisitBetween(expr *Between) error
	VisitCoalesce(expr *CoalesceExpr) error
	VisitNegate(expr *NegateExpr) error
	VisitPositive(expr *Positive) error
	VisitNot(expr *SqlNot) error
	VisitColumnList(expr *ColumnList) error
	VisitCase(expr *CaseExpr) error
	VisitBool(expr Bool) error
	VisitParam(expr Param) error
	VisitFuncDef(expr *FuncDef) error
	VisitNamedValue(expr NamedValue) error
	VisitSetOperation(expr *SetOperation) error
	VisitOver(expr *OverExpr) error
	VisitValue(expr ValueExpr) error
	VisitValues(expr *Values) error
	VisitNull() error
	VisitRawSql(expr RawSql) error
	VisitFrameBetween(expr *FrameBetween) error
	VisitDateDiff(expr *DateDiff) error
	VisitDateAdd(expr *DateAdd) error
}
