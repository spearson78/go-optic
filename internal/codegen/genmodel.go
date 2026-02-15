package codegen

//go:generate ../../makelens --root CG codegen genmodel.go genmodel_generated.go

type FileDef struct {
	Package    string
	Imports    []string
	DotImports []string
	Structs    []StructDef
	Vars       []VarDef
	Funcs      []FuncDef
}

type TypeExpression any

type TypeDef struct {
	Name       string
	TypeParams []TypeExpression
	Constraint string
}

type Star struct {
	Type TypeExpression
}

type DeRef struct {
	Value Expression
}

type SliceDef struct {
	Type TypeExpression
}

type MapDef struct {
	Key  TypeExpression
	Type TypeExpression
}

type FuncType struct {
	Params      []Param
	ReturnTypes []TypeExpression
}

type StructDef struct {
	Name       string
	TypeParams []TypeExpression
	Fields     []FieldDef
	Methods    []MethodDef
}

type FieldDef struct {
	Name string
	Type TypeExpression
}

type VarDef struct {
	Name  string
	Type  TypeExpression
	Value Expression
}

type FuncDef struct {
	Docs        []string
	Name        string
	TypeParams  []TypeExpression
	Params      []Param
	ReturnTypes []TypeExpression
	Body        []Statement
}

type MethodDef struct {
	Docs        []string
	Name        string
	Params      []Param
	ReturnTypes []TypeExpression
	Body        []Statement
}

type Statement any

type Expression any

type ReturnStmnt struct {
	Values []Expression
}

type CallExpr struct {
	Func       Expression
	TypeParams []TypeExpression
	Params     []Expression
}

type MethodCallExpr struct {
	Receiver Expression
	Name     string
	Params   []Expression
}

type GenFuncExpr struct {
	Params      []Param
	Body        []Statement
	ReturnTypes []TypeExpression
}

type Param struct {
	Name string
	Type TypeExpression
}

type BinaryExpr struct {
	Op    string
	Left  Expression
	Right Expression
}

type StructExpr struct {
	Type   TypeDef
	Fields []AssignField
}

type SliceExpr struct {
	Type   TypeDef
	Values []Expression
}

type AssignField struct {
	Name  string
	Value Expression
}

type AssignVar struct {
	Declare bool
	Vars    []string
	Value   Expression
}

type AddressExpr struct {
	Target Expression
}

type StringLiteral string

type IfStmnt struct {
	Condition Expression
	OnTrue    []Statement
	OnFalse   []Statement
}
