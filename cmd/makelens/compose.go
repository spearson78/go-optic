package main

import . "github.com/spearson78/go-optic/internal/codegen"

func RwL(param Expression) Expression {
	return CallExpr{
		Func:   "optic.RwL",
		Params: []Expression{param},
	}
}

func RetM(param Expression) Expression {
	return CallExpr{
		Func:   "optic.RetM",
		Params: []Expression{param},
	}
}

func EErrL(param Expression) Expression {
	return CallExpr{
		Func:   "optic.EErrL",
		Params: []Expression{param},
	}
}

func DirL(param Expression) Expression {
	return CallExpr{
		Func:   "optic.DirL",
		Params: []Expression{param},
	}
}

func Compose(params ...Expression) Expression {
	return CallExpr{
		Func:   "optic.Compose",
		Params: params,
	}
}

func ComposeLeft(params ...Expression) Expression {
	return CallExpr{
		Func:   "optic.ComposeLeft",
		Params: params,
	}
}

func RetL(o Expression) Expression {
	return CallExpr{
		Func: "optic.RetL",
		Params: []Expression{
			o,
		},
	}
}

func Ud(o Expression) Expression {
	return CallExpr{
		Func: "optic.Ud",
		Params: []Expression{
			o,
		},
	}
}
