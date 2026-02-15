package otree_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/otree"
)

func ExampleRewrite() {

	tree := Tree[any](
		"+",
		ValCol(
			Tree[any](
				"+",
				ValCol(
					Tree[any](1),
					Tree[any](2),
				),
			),
			Tree[any](3),
		),
	)

	opName := Compose(
		TreeValue[int, any, Pure](),
		DownCast[any, string](),
	)

	value := Compose(
		TreeValue[int, any, Pure](),
		DownCast[any, int](),
	)

	//Rewrite handler that evaluates the "+" function
	rewriteHandler := RewriteOp(func(tree TreeNode[int, any, Pure]) (TreeNode[int, any, Pure], bool) {
		//If all the children are literal values we may be able to evaluate the expression
		childrenAllValues, _ := MustGetFirst(
			All(
				OTree[int, any, Pure]().Children().Traverse(),
				NotEmpty(value),
			),
			tree,
		)

		if childrenAllValues {

			funcName, _ := MustGetFirst(opName, tree)

			switch funcName {
			case "+":
				sumParams, _ := MustGetFirst(Reduce(Compose(OTree[int, any, Pure]().Children().Traverse(), value), Sum[int]()), tree)
				return Tree[any](sumParams), true
			}
		}

		return tree, false
	})

	newTree := MustGet(Rewrite(TraverseTreeChildren[int, any, Pure](), rewriteHandler), tree)

	fmt.Println(newTree)

	data := lo.T2("alpha", tree)

	newTuple := MustModify(T2B[string, TreeNode[int, any, Pure]](), Rewrite(TraverseTreeChildren[int, any, Pure](), rewriteHandler), data)
	fmt.Println(newTuple.A, newTuple.B)

	//Output:
	//{6 : Col[]}
	//alpha {6 : Col[]}
}

func ExampleTopDownMatch() {

	data := Tree(
		":=",
		ValCol(
			Tree("x"),
			Tree("+",
				ValCol(
					Tree("1"),
					Tree("2"),
				),
			),
		),
	)

	//Match + nodes
	matchAdd := Compose(
		OTree[int, string, Pure]().Value(),
		Eq("+"),
	)

	//TopDownMatch will find all matching sub trees
	deepAdd := TopDownMatch(
		TraverseTreeChildren[int, string, Pure](),
		matchAdd,
	)

	var result []TreeNode[int, string, Pure] = MustGet(SliceOf(deepAdd, 10), data)
	fmt.Println(result)

	var modifyResult TreeNode[int, string, Pure]
	modifyResult, err := Set(
		Compose(
			deepAdd,
			TreeValue[int, string, Pure](),
		),
		"-",
		data,
	)
	fmt.Println(modifyResult, err)

	//Output:
	//[{+ : Col[0:{1 : Col[]} 1:{2 : Col[]}]}]
	//{:= : Col[0:{x : Col[]} 1:{- : Col[0:{1 : Col[]} 1:{2 : Col[]}]}]} <nil>
}

func ExampleResolvePath() {

	data := Tree(
		"alpha",
		ValCol(
			Tree("beta"),
			Tree("gamma",
				ValCol(Tree("delta")),
			),
		),
	)

	optic := ResolvePath(TraverseTreeChildren[int, string, Pure](), Path(1, 0))

	var result TreeNode[int, string, Pure]
	var ok bool
	result, ok = MustGetFirst(optic, data)
	fmt.Println(result, ok)

	modifyResult := MustSet(optic, Tree("omega"), data)
	fmt.Println(modifyResult)

	//Output:
	//{delta : Col[]} true
	//{alpha : Col[0:{beta : Col[]} 1:{gamma : Col[0:{omega : Col[]}]}]}
}

func TestRewriteTransform(t *testing.T) {

	tree := Tree[any](
		"+",
		ValCol(
			Tree[any](
				"*",
				ValCol(
					Tree[any](2),
					Tree[any](10),
				),
			),
			Tree[any](5),
		),
	)

	if r, err := Get(SliceOf(Compose(TopDown(TraverseTreeChildren[int, any, Pure]()), TreeValue[int, any, Pure]()), 3), tree); err != nil || !reflect.DeepEqual(r, []any{"+", "*", 2, 10, 5}) {
		t.Fatalf("IDDF 1 %v", r)
	}

	rewriteHandler := RewriteOp(func(tree TreeNode[int, any, Pure]) (TreeNode[int, any, Pure], bool) {

		switch op := MustGet(TreeValue[int, any, Pure](), tree).(type) {
		case string:

			children := MustGet(SliceOf(Compose3(TraverseTreeChildren[int, any, Pure](), TreeValue[int, any, Pure](), DownCast[any, int]()), 2), tree)
			if len(children) == 2 {
				leftNum := children[0]
				rightNum := children[1]

				switch op {
				case "+":
					return Tree[any](leftNum + rightNum), true
				case "-":
					return Tree[any](leftNum - rightNum), true
				case "*":
					return Tree[any](leftNum * rightNum), true
				case "/":
					return Tree[any](leftNum / rightNum), true
				}
			}
		}

		return tree, false
	})

	if r := MustGet(Compose(Rewrite(TraverseTreeChildren[int, any, Pure](), rewriteHandler), TreeValue[int, any, Pure]()), tree); r != 25 {
		t.Fatalf("Rewrite %v", r)
	}

	swap := Op(func(tree TreeNode[int, any, Pure]) TreeNode[int, any, Pure] {

		switch op := MustGet(TreeValue[int, any, Pure](), tree).(type) {
		case string:

			children := MustGet(
				ColOf(
					Indexing(
						Reversed(TraverseTreeChildren[int, any, Pure]()),
					),
				),
				tree,
			)

			return Tree[any](
				op,
				children,
			)
		}

		return tree
	})

	//{+ : SeqF[
	//	0:{* : SeqF[
	//		0:{10 : SeqF[]}
	//		1:{2 : SeqF[]}]}
	//	1:{5 : SeqF[]}]}

	swapExpected := Tree[any](
		"+",
		ValCol(
			Tree[any](5),
			Tree[any](
				"*",
				ValCol(
					Tree[any](10),
					Tree[any](2),
				),
			),
		),
	)

	if r := MustModify(
		TopDown(
			TraverseTreeChildren[int, any, Pure](),
		),
		swap,
		tree,
	); !MustGet(EqTreeI(r, EqT2[int](), EqT2[any]()), swapExpected) {
		t.Fatalf("Transform %v", r)
	}

}

func ExampleTopDownFiltered() {

	tree := Tree(
		"A",
		ValCol(
			Tree(
				"B",
				ValCol(
					Tree("C"),
					Tree("FilterMe"),
				),
			),
			Tree("D"),
		),
	)

	children := TraverseTreeChildren[int, string, Pure]()

	pred := Compose(
		OTree[int, string, Pure]().Value(),
		Ne("FilterMe"),
	)

	optic := TopDownFiltered(
		children,
		pred,
	)

	res := MustGet(
		SliceOf(
			Compose(
				optic,
				TreeValue[int, string, Pure](),
			),
			5,
		),
		tree,
	)

	fmt.Println(res)

	//Output:
	//[A B C D]
}

func ExampleBottomUpFiltered() {

	tree := Tree(
		"A",
		ValCol(
			Tree(
				"B",
				ValCol(
					Tree("C"),
					Tree("FilterMe"),
				),
			),
			Tree("D"),
		),
	)

	children := TraverseTreeChildren[int, string, Pure]()

	pred := Compose(
		OTree[int, string, Pure]().Value(),
		Ne("FilterMe"),
	)

	optic := BottomUpFiltered(
		children,
		pred,
	)

	res := MustGet(
		SliceOf(
			Compose(
				optic,
				TreeValue[int, string, Pure](),
			),
			5,
		),
		tree,
	)

	fmt.Println(res)

	//Output:
	//[C B D A]
}

func ExampleBreadthFirstFiltered() {

	tree := Tree(
		"A",
		ValCol(
			Tree(
				"B",
				ValCol(
					Tree("C"),
					Tree("FilterMe"),
				),
			),
			Tree("D"),
		),
	)

	children := TraverseTreeChildren[int, string, Pure]()

	pred := Compose(
		OTree[int, string, Pure]().Value(),
		Ne("FilterMe"),
	)

	optic := BreadthFirstFiltered(
		children,
		pred,
	)

	res := MustGet(
		SliceOf(
			Compose(
				optic,
				TreeValue[int, string, Pure](),
			),
			5,
		),
		tree,
	)

	fmt.Println(res)

	//Output:
	//[A B D C]
}

func ExampleTopDown() {

	tree := Tree(
		"A",
		ValCol(
			Tree(
				"B",
				ValCol(
					Tree("C"),
				),
			),
			Tree("D"),
		),
	)

	children := TraverseTreeChildren[int, string, Pure]()

	optic := TopDown(
		children,
	)

	res := MustGet(
		SliceOf(
			Compose(
				optic,
				TreeValue[int, string, Pure](),
			),
			5,
		),
		tree,
	)

	fmt.Println(res)

	//Output:
	//[A B C D]
}

func ExampleBottomUp() {

	tree := Tree(
		"A",
		ValCol(
			Tree(
				"B",
				ValCol(
					Tree("C"),
				),
			),
			Tree("D"),
		),
	)

	children := TraverseTreeChildren[int, string, Pure]()

	optic := BottomUp(
		children,
	)

	res := MustGet(
		SliceOf(
			Compose(
				optic,
				TreeValue[int, string, Pure](),
			),
			5,
		),
		tree,
	)

	fmt.Println(res)

	//Output:
	//[C B D A]
}

func ExampleBreadthFirst() {

	tree := Tree(
		"A",
		ValCol(
			Tree(
				"B",
				ValCol(
					Tree("C"),
				),
			),
			Tree("D"),
		),
	)

	children := TraverseTreeChildren[int, string, Pure]()

	optic := BreadthFirst(
		children,
	)

	res := MustGet(
		SliceOf(
			Compose(
				optic,
				TreeValue[int, string, Pure](),
			),
			5,
		),
		tree,
	)

	fmt.Println(res)

	//Output:
	//[A B D C]
}
