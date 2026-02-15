package optic_test

import (
	"reflect"
	"testing"
	"unicode"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

func FuzzFieldLens(f *testing.F) {

	fieldA := FieldLens(func(source *TestStruct) *string { return &source.A })

	f.Add("Hello World", "new value")
	f.Fuzz(func(t *testing.T, a, b string) {
		testStruct := TestStruct{
			A: a,
		}

		ValidateOpticTest(t, fieldA, testStruct, b)

	})

}

func TestFieldLens(t *testing.T) {
	testStruct := TestStruct{
		A: "Hello World",
	}

	fieldA := FieldLens(func(source *TestStruct) *string { return &source.A })
	ret, err := Modify(fieldA, Op(func(focus string) string {
		return "BOB"
	}), testStruct)

	if err != nil || ret.A != "BOB" {
		t.Fatal("TestFieldLens")
	}

}

func TestFieldLensExpr(t *testing.T) {

	personAge := FieldLens(func(source *Person) *int { return &source.Age })

	foundField, found := personAge.AsExpr().(expr.FieldLens)

	if !found {
		t.Fatal("field not found")
	}

	if foundField.OpticS().PkgPath() != "github.com/spearson78/go-optic_test" ||
		foundField.Field.Name != "Age" ||
		foundField.OpticS().Name() != "Person" ||
		foundField.Field.Tag != `example:"age"` {
		t.Fatal("field", foundField.Field.PkgPath, foundField.Field.Name, foundField.OpticS().Name(), foundField.Field.Tag)
	}
}

func TestFieldLens2(t *testing.T) {

	p := Person{
		Name: "steve",
		Age:  45,
	}
	nameLens := FieldLens[Person, string](func(s *Person) *string { return &s.Name })

	if r, err := Modify(nameLens, Op(func(s string) string {
		return s + "!"
	}), p); err != nil || r.Name != "steve!" {
		t.Fatalf("nameLens 1 : %v", r)
	}
}

func TestStrBothTraverse(t *testing.T) {

	strBoth := TraverseT2[string]()
	strTraverse := TraverseString()

	strBothTraverse := Compose(strBoth, strTraverse)

	if res, err := Modify(strBothTraverse, Op(func(r rune) rune {
		return unicode.ToUpper(r)
	}), lo.T2("hello", "world")); err != nil || res.A != "HELLO" || res.B != "WORLD" {
		t.Fatal("")
	}

}

func TestStrBothTraverseSlice(t *testing.T) {
	strBoth := TraverseT2[string]()
	strTraverse := TraverseString()

	strBothTraverse := WithComprehension(Compose(strBoth, strTraverse))

	if res, err := Modify(strBothTraverse, Op(func(r rune) []rune {
		return []rune{
			unicode.ToUpper(r),
			unicode.ToLower(r),
		}
	}), lo.T2("AA", "BB")); err != nil || !reflect.DeepEqual(res, []lo.Tuple2[string, string]{
		lo.T2("AA", "BB"),
		lo.T2("AA", "Bb"),
		lo.T2("AA", "bB"),
		lo.T2("AA", "bb"),
		lo.T2("Aa", "BB"),
		lo.T2("Aa", "Bb"),
		lo.T2("Aa", "bB"),
		lo.T2("Aa", "bb"),
		lo.T2("aA", "BB"),
		lo.T2("aA", "Bb"),
		lo.T2("aA", "bB"),
		lo.T2("aA", "bb"),
		lo.T2("aa", "BB"),
		lo.T2("aa", "Bb"),
		lo.T2("aa", "bB"),
		lo.T2("aa", "bb"),
	}) {
		t.Fatalf("TestStrBothTraverseSlice : %v", res)
	}
}

func TestPtrFieldLensWithPtrField(t *testing.T) {

	type Tree struct {
		Value string
		Left  *Tree
		Right *Tree
	}

	value := PtrFieldLens(func(t *Tree) *string { return &t.Value })
	left := PtrFieldLens(func(t *Tree) **Tree { return &t.Left })
	//right := PtrFieldLens(func(t *Tree) **Tree { return &t.Right })

	tree := &Tree{
		Value: "root",
		Left: &Tree{
			Value: "left",
		},
		Right: &Tree{
			Value: "right",
		},
	}

	newTree := MustSet(Compose(left, value), "edit-left", tree)

	if newTree.Left.Value != "edit-left" {
		t.Fatal("newTree.left 1")
	}

	if tree.Left.Value != "left" {
		t.Fatal("newTree.left 1")
	}

	newTree = MustSet(left, &Tree{Value: "edit-left"}, tree)

	if newTree.Left.Value != "edit-left" {
		t.Fatal("newTree.left 2")
	}

	if tree.Left.Value != "left" {
		t.Fatal("newTree.left 2")
	}

}
