package otree_test

import (
	"fmt"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/otree"
)

func ExampleRewriteI() {

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
	rewriteHandler := RewriteOpI(func(path *PathNode[int], tree TreeNode[int, any, Pure]) (TreeNode[int, any, Pure], bool) {

		//Use index to prevent evaluation at the root of the tree
		if path == nil {
			return tree, false
		}

		//If all the children are literal values we may be able to evaluate the expression
		childrenAllValues, _ := MustGetFirst(All(OTree[int, any, Pure]().Children().Traverse(), NotEmpty(value)), tree)
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

	newTree := MustGet(RewriteI(TraverseTreeChildren[int, any, Pure](), rewriteHandler), tree)

	fmt.Println(newTree)

	data := lo.T2("alpha", tree)

	newTuple := MustModify(T2B[string, TreeNode[int, any, Pure]](), RewriteI(TraverseTreeChildren[int, any, Pure](), rewriteHandler), data)
	fmt.Println(newTuple.A, newTuple.B)

	//Output:
	//{+ : Col[0:{3 : Col[]} 1:{3 : Col[]}]}
	//alpha {+ : Col[0:{3 : Col[]} 1:{3 : Col[]}]}
}

func ExampleTopDownFilteredI() {

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

	pred := Compose3(
		ValueIIndex[*PathNode[int], TreeNode[int, string, Pure]](),
		EqPath(Path(0, 1), EqT2[int]()),
		Not(),
	)

	optic := TopDownFilteredI(
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

func ExampleBottomUpFilteredI() {

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

	pred := Compose3(
		ValueIIndex[*PathNode[int], TreeNode[int, string, Pure]](),
		EqPath(Path(0, 1), EqT2[int]()),
		Not(),
	)

	optic := BottomUpFilteredI(
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

func ExampleBreadthFirstFilteredI() {

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

	pred := Compose3(
		ValueIIndex[*PathNode[int], TreeNode[int, string, Pure]](),
		EqPath(Path(0, 1), EqT2[int]()),
		Not(),
	)

	optic := BreadthFirstFilteredI(
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
