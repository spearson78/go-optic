package expr

import (
	"iter"
	"reflect"
	"testing"

	"github.com/samber/lo"
	"github.com/samber/mo"
)

type besideTestStruct struct {
	A lo.Tuple2[int, int]
	B lo.Tuple2[int, int]
}

func TestBeside(t *testing.T) {
	b := Beside{
		OpticTypeExpr: NewOpticTypeExpr[mo.Either[void, void], lo.Tuple2[int, int], lo.Tuple2[int, int], int, int, returnMany, readWrite, uniDir, pure](),
		Optic1: FieldLens{
			OpticTypeExpr: NewOpticTypeExpr[void, besideTestStruct, besideTestStruct, lo.Tuple2[int, int], lo.Tuple2[int, int], returnOne, readWrite, uniDir, pure](),
			Field:         reflect.TypeFor[besideTestStruct]().Field(0),
		},
		Optic2: FieldLens{
			OpticTypeExpr: NewOpticTypeExpr[void, besideTestStruct, besideTestStruct, lo.Tuple2[int, int], lo.Tuple2[int, int], returnOne, readWrite, uniDir, pure](),
			Field:         reflect.TypeFor[besideTestStruct]().Field(1),
		},
	}

	if e := b.String(); e != "Beside(FieldLens(expr.besideTestStruct.A),FieldLens(expr.besideTestStruct.B))" {
		t.Fatal("Beside.String()", e)
	}
}

func TestTraverse(t *testing.T) {
	e := Traverse{
		OpticTypeExpr: NewOpticTypeExpr[int, iter.Seq2[int, string], iter.Seq2[int, float64], int, float64, returnMany, readWrite, uniDir, pure](),
	}

	if s := e.String(); s != "Traverse(int,iter.Seq2[int,string],iter.Seq2[int,float64],int,float64)" {
		t.Fatal("Traverse(int,int,iter.Seq2[int,string],iter.Seq2[int,float64],int,float64)", e)
	}
}
