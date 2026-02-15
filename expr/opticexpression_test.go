package expr

import "testing"

type void struct{}

// Indicates this optic returns exactly one result. Used to prevent compilation when passing a [ReturnMany] optic to a [ReturnOne] action e.g. [View]
type returnOne void

// Indicates this optic returns exactly zero or more results.
type returnMany []void

// Indicates this optic is read only. Used to prevent compilation when passing a [ReadOnly] optic to a [ReadWrite] action e.g. [Over]
type readOnly []void

// Indicates this optic supports read and write actions.
type readWrite void

// Indicates this optic supports reverseget
type biDir void

// Indicates this optic does not support reverseget
type uniDir []void

type pure void

type impure []void

func TestOpticExpressionString(t *testing.T) {

	if e := NewOpticTypeExpr[int, []string, []string, string, string, returnOne, readOnly, uniDir, pure]().Signature(); e != "int,[]string,[]string,string,string,Getter,Pure" {
		t.Fatalf("1 NewOpticTypeExpr().String() : %v", e)
	}

	if e := NewOpticTypeExpr[int, []string, []string, string, string, returnOne, readOnly, biDir, pure]().Signature(); e != "int,[]string,[]string,string,string,ReturnOneReadOnlyBiDir,Pure" {
		t.Fatalf("2 NewOpticTypeExpr().String() : %v", e)
	}

	if e := NewOpticTypeExpr[int, []string, []string, string, string, returnOne, readWrite, uniDir, pure]().Signature(); e != "int,[]string,[]string,string,string,Lens,Pure" {
		t.Fatalf("3 NewOpticTypeExpr().String() : %v", e)
	}

	if e := NewOpticTypeExpr[int, []string, []string, string, string, returnOne, readWrite, biDir, pure]().Signature(); e != "int,[]string,[]string,string,string,Iso,Pure" {
		t.Fatalf("4 NewOpticTypeExpr().String() : %v", e)
	}

	if e := NewOpticTypeExpr[int, []string, []string, string, string, returnMany, readOnly, uniDir, pure]().Signature(); e != "int,[]string,[]string,string,string,Iteration,Pure" {
		t.Fatalf("5 NewOpticTypeExpr().String() : %v", e)
	}

	if e := NewOpticTypeExpr[int, []string, []string, string, string, returnMany, readOnly, biDir, pure]().Signature(); e != "int,[]string,[]string,string,string,ReturnManyReadOnlyBiDir,Pure" {
		t.Fatalf("6 NewOpticTypeExpr().String() : %v", e)
	}

	if e := NewOpticTypeExpr[int, []string, []string, string, string, returnMany, readWrite, uniDir, pure]().Signature(); e != "int,[]string,[]string,string,string,Traversal,Pure" {
		t.Fatalf("7 NewOpticTypeExpr().String() : %v", e)
	}

	if e := NewOpticTypeExpr[int, []string, []string, string, string, returnMany, readWrite, biDir, pure]().Signature(); e != "int,[]string,[]string,string,string,Prism,Pure" {
		t.Fatalf("8 NewOpticTypeExpr().String() : %v", e)
	}

}
