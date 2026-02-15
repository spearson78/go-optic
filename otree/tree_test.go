package otree_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/samber/mo"
	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/otree"
)

func TestTreeFilter(t *testing.T) {

	data := Tree(
		"a",
		ValCol(
			Tree(
				"b",
				ValCol(
					Tree("b.1"),
					Tree("b.2"),
				),
			),
			Tree("c"),
		),
	)

	//Filtered determines whether to focus on the tree node value. It does not make structural changes to the tree.
	optic := Filtered(
		Compose(TopDown(TraverseTreeChildren[int, string, Pure]()), TreeValue[int, string, Pure]()),
		Ne("b"),
	)

	res, err := Get(SliceOf(optic, 10), data)
	if err != nil || !reflect.DeepEqual(res, []string{"a", "b.1", "b.2", "c"}) {
		t.Fatal("get filter tree seq", res, err)
	}

	newTree, err := Modify(optic, Op(strings.ToUpper), data)

	expectedTree := Tree(
		"A",
		ValCol(
			Tree(
				"b",
				ValCol(
					Tree("B.1"),
					Tree("B.2"),
				),
			),
			Tree("C"),
		),
	)

	if !MustGet(EqTreeI(expectedTree, EqT2[int](), EqT2[string]()), newTree) {
		t.Fatal("modify filter branch", newTree)
	}
}

func TestTreeBranchFilter(t *testing.T) {

	data := Tree(
		"a",
		ValCol(
			Tree(
				"b",
				ValCol(
					Tree("b.1"),
					Tree("b.2"),
				),
			),
			Tree("c"),
		),
	)

	//BranchFiltered performas a top down traversal and stops at branches that are filtered
	optic := Compose(
		TopDownFiltered(
			TraverseTreeChildren[int, string, Pure](),
			Compose(
				TreeValue[int, string, Pure](),
				Ne("b"),
			),
		),
		TreeValue[int, string, Pure](),
	)

	res, err := Get(SliceOf(optic, 10), data)
	if err != nil || !reflect.DeepEqual(res, []string{"a", "c"}) {
		t.Fatal("get filter tree seq", res, err)
	}

	newTree, err := Modify(optic, Op(strings.ToUpper), data)

	expectedTree := Tree(
		"A",
		ValCol(
			Tree(
				"b",
				ValCol(
					Tree("b.1"),
					Tree("b.2"),
				),
			),
			Tree("C"),
		),
	)

	if !MustGet(EqTreeI(expectedTree, EqT2[int](), EqT2[string]()), newTree) {
		t.Fatal("modify filter branch", newTree)
	}
}

func TestTreeFilteredCol(t *testing.T) {

	data := Tree(
		"a",
		ValCol(
			Tree(
				"b",
				ValCol(
					Tree("b.1"),
					Tree("b.2"),
				),
			),
			Tree("c"),
		),
	)

	expectedTree := TreeI(
		"a",
		IxMatchComparable[int](),
		ValColI(IxMatchComparable[int](), ValI(1, Tree("c"))),
	)

	if res, err := Modify(
		ComposeLeft(
			BottomUp(
				TraverseTreeChildren[int, string, Pure](),
			),
			TreeChildren[int, string, Pure](),
		),
		FilteredCol[int](
			EPure(Compose(
				TreeValue[int, string, Pure](),
				Ne("b"),
			)),
		),
		data,
	); err != nil || !MustGet(EqTreeI(expectedTree, EqT2[int](), EqT2[string]()), res) {
		t.Fatal("TestTreeFilteredCol op", res, err)
	}

}

func ExampleTree() {

	data := Tree[string](
		"Root",
		ValCol(
			Tree("Leaf 1"),
			Tree("Leaf 2"),
		),
	)

	fmt.Println(data)

	//Output:
	//{Root : Col[0:{Leaf 1 : Col[]} 1:{Leaf 2 : Col[]}]}
}

func ExampleTreeValue() {

	data := Tree[string](
		"Root",
		ValCol(
			Tree("Leaf 1"),
			Tree("Leaf 2"),
		),
	)

	var result string = MustGet(TreeValue[int, string, Pure](), data)
	fmt.Println(result)

	modifyResult := MustSet(TreeValue[int, string, Pure](), "New Root", data)

	fmt.Println(modifyResult)

	//Output:
	//Root
	//{New Root : Col[0:{Leaf 1 : Col[]} 1:{Leaf 2 : Col[]}]}

}

func TestDiffTree(t *testing.T) {

	before := Tree(
		"alpha", //[]
		ValCol(
			Tree(
				"beta", //[0]
				ValCol(
					Tree("gamma"), //[0 0]
					Tree(
						"delta", //[0 1]
						ValCol(
							Tree("epsilon"), //[0 1 0]
							Tree("zeta"),    //[0 1 1]
							Tree("eta"),     //[0 1 2]
						),
					),
				),
			),
			Tree("theta"), //[1]
		),
	)

	after := Tree(
		"alpha", //0
		ValCol(
			Tree(
				"beta", //[0]
				ValCol(
					Tree(
						"Delta", //[0 0]
						ValCol(
							Tree("Zeta"), //[0 0 0]
							Tree("eta"),  //[0 0 1]
							Tree("Iota"), //[0 0 2]
						),
					),
					Tree("gamma"), //[0 1]
					Tree("Kappa"), //[0 2]
				),
			),
		),
	)

	dist := DistancePercent(
		EditDistance(Compose(TreeValue[int, string, Pure](), TraverseString()), EditLevenshtein, EqT2[rune](), 10),
		Compose(TreeValue[int, string, Pure](), Length(TraverseString())),
	)

	diffType := FieldLens(func(source *Diff[*PathNode[int], TreeNode[int, string, Pure]]) *DiffType { return &source.Type })

	a := Compose(Non[TreeNode[int, string, Pure]](TreeNode[int, string, Pure]{}, EqDeepT2[TreeNode[int, string, Pure]]()), TreeValue[int, string, Pure]())

	c := ComposeLeft(DiffTreeT2(TreeChildren[int, string, Pure](), 0.5, dist, EqT2[int](), DiffNone, true), a)

	diff := ReIndexed(c, diffType, EqT2[DiffType]())

	if res := MustGet(SliceOf(WithIndex(diff), 10), lo.T2(mo.Some(after), mo.Some(before))); !MustGet(EqDeepT2[[]ValueI[DiffType, string]](), lo.T2(res, []ValueI[DiffType, string]{
		ValI(
			DiffModify,
			"Delta",
		),
		ValI(
			DiffRemove,
			"",
		),
		ValI(
			DiffModify,
			"Zeta",
		),
		ValI(
			DiffModify,
			"eta",
		),
		ValI(
			DiffAdd,
			"Iota",
		),
		ValI(
			DiffModify,
			"gamma",
		),
		ValI(
			DiffAdd,
			"Kappa",
		),
		ValI(
			DiffRemove,
			"",
		),
	})) {
		t.Fatal("iterate", res)
	}

	fmapSeq := 0
	expected := Tree(
		"alpha",
		ValCol(
			Tree(
				"beta",
				ValCol(
					Tree(
						"Delta:0",
						ValCol(
							Tree("Zeta:1"),
							Tree("eta:2"),
							Tree("Iota:3"),
						),
					),
					Tree("gamma:4"),
					Tree("Kappa:5"),
				),
			),
		),
	)

	if res := MustModify(diff, Op(func(focus string) string {
		if focus == "" {
			return ""
		}

		ret := fmt.Sprintf("%v:%v", focus, fmapSeq)
		fmapSeq++
		return ret
	}), lo.T2(mo.Some(after), mo.Some(before))); !MustGet(EqTreeI(res.A.MustGet(), EqT2[int](), EqT2[string]()), expected) {
		t.Fatal("modify", res.A.MustGet(), "expected", expected)
	}

}

func TestReIndexedTree(t *testing.T) {

	data := TreeE(
		"alpha", //[]
		ValColE(
			ValE(TreeE(
				"beta",
				ValColE( //[0]
					ValE(TreeE("gamma"), nil), //[0 0]
					ValE(TreeE(
						"delta",
						ValColE( //[0 1]
							ValE(TreeE("epsilon"), nil), //[0 1 0]
							ValE(TreeE("zeta"), nil),    //[0 1 1]
							ValE(TreeE("eta"), nil),     //[0 1 2]
						),
					), nil),
				),
			), nil),
			ValE(TreeE("theta"), nil), //[1]
		),
	)

	ixmatch := IxMatchComparable[string]()

	expected := TreeIE(
		"alpha", //[]
		ixmatch,
		ValColIE[string, TreeNode[string, string, Err], Err](
			ixmatch,
			ValIE("0", TreeIE(
				"beta", //[0]
				ixmatch,
				ValColIE[string, TreeNode[string, string, Err], Err](
					ixmatch,
					ValIE("0", TreeIE[string, string, Err]("gamma", ixmatch), nil), //[0 0]
					ValIE("1", TreeIE(
						"delta", //[0 1]
						ixmatch,
						ValColIE[string, TreeNode[string, string, Err], Err](
							ixmatch,
							ValIE("0", TreeIE[string, string, Err]("epsilon", ixmatch), nil), //[0 1 0]
							ValIE("1", TreeIE[string, string, Err]("zeta", ixmatch), nil),    //[0 1 1]
							ValIE("2", TreeIE[string, string, Err]("eta", ixmatch), nil),     //[0 1 2]
						),
					), nil),
				),
			), nil),
			ValIE("1", TreeIE[string, string, Err]("theta", ixmatch), nil), //[1]
		),
	)

	optic := ReIndexedTree[string](AsReverseGet(ParseInt[int](10, 32)))

	res, err := Get(optic, mo.Some(data))
	if err != nil {
		t.Fatal(res.MustGet(), "!=", expected, err)
	}

	treeEq, err := Get(EqTreeI(res.MustGet(), EqT2[string](), EqT2[string]()), expected)
	if err != nil || !treeEq {
		t.Fatal(err)
	}

}

type BinaryTree struct {
	Left  *BinaryTree
	Right *BinaryTree
	Value string
}

func (t *BinaryTree) String() string {
	return fmt.Sprintf("(%v:%v:%v)", t.Value, t.Left, t.Right)
}

func TestLanguagesShouldHaveTreePrimitive(t *testing.T) {

	left := PtrFieldLens(func(source *BinaryTree) **BinaryTree { return &source.Left })
	right := PtrFieldLens(func(source *BinaryTree) **BinaryTree { return &source.Right })
	value := PtrFieldLens(func(source *BinaryTree) *string { return &source.Value })

	children := Filtered(
		Concat(
			left,
			right,
		),
		Ne[*BinaryTree](nil),
	)

	data := &BinaryTree{
		Value: "root",
		Left: &BinaryTree{
			Value: "root.left",
			Left: &BinaryTree{
				Value: "root.left.left",
			},
		},
		Right: &BinaryTree{
			Value: "root.right",
			Left: &BinaryTree{
				Value: "root.right.left",
			},
		},
	}

	var res []*BinaryTree

	optic := TopDown(children)

	for v := range MustGet(SeqOf(optic), data) {
		res = append(res, v)
	}

	if !reflect.DeepEqual(res, []*BinaryTree{data, data.Left, data.Left.Left, data.Right, data.Right.Left}) {
		t.Fatal("1", res)
	}

	res = nil
	for v := range MustGet(SeqOf(optic), data) {
		res = append(res, v)
		if v.Value == "root.left" {
			break
		}
	}

	if !reflect.DeepEqual(res, []*BinaryTree{data, data.Left}) {
		t.Fatal("2", res)
	}

	pruned := TopDown(If(
		Compose(value, Eq("root.left")),
		Nothing[*BinaryTree, *BinaryTree](),
		EPure(children),
	))

	res = nil
	for v := range MustGet(SeqOf(pruned), data) {
		res = append(res, v)
	}

	if !reflect.DeepEqual(res, []*BinaryTree{data, data.Left, data.Right, data.Right.Left}) {
		t.Fatal("3", res)
	}

	iterChildren := Concat3(
		AppendString(StringCol("a")),
		AppendString(StringCol("b")),
		AppendString(StringCol("c")),
	)

	var strRes []string
	for v := range MustGet(SeqOf(BreadthFirst(iterChildren)), "") {
		if len(v) > 3 {
			break
		}
		strRes = append(strRes, v)
	}

	if !reflect.DeepEqual(strRes, []string{"", "a", "b", "c", "aa", "ab", "ac", "ba", "bb", "bc", "ca", "cb", "cc", "aaa", "aab", "aac", "aba", "abb", "abc", "aca", "acb", "acc", "baa", "bab", "bac", "bba", "bbb", "bbc", "bca", "bcb", "bcc", "caa", "cab", "cac", "cba", "cbb", "cbc", "cca", "ccb", "ccc"}) {
		t.Fatal("4", strRes)
	}

}

func TestTreeEq(t *testing.T) {
	a := Tree(
		"alpha", //[]
		ValCol(
			Tree(
				"beta", //[0]
				ValCol(
					Tree("gamma"), //[0 0]
					Tree(
						"delta", //[0 1]
						ValCol(
							Tree("epsilon"), //[0 1 0]
							Tree("zeta"),    //[0 1 1]
							Tree("eta"),     //[0 1 2]
						),
					),
				),
			),
			Tree("theta"), //[1]
		),
	)

	b := Tree(
		"alpha", //[]
		ValCol(
			Tree(
				"beta", //[0]
				ValCol(
					Tree("gamma"), //[0 0]
					Tree(
						"delta", //[0 1]
						ValCol(
							Tree("epsilon"), //[0 1 0]
							Tree("zeta"),    //[0 1 1]
							Tree("eta"),     //[0 1 2]
						),
					),
				),
			),
			Tree("theta"), //[1]
			Tree("extra"), //[2]
		),
	)

	if res := MustGet(EqTreeI(a, EqT2[int](), EqT2[string]()), a); res != true {
		t.Fatal(res)
	}

	if res := MustGet(EqTreeI(b, EqT2[int](), EqT2[string]()), a); res != false {
		t.Fatal(res)
	}

	if res := MustGet(EqTreeI(a, EqT2[int](), EqT2[string]()), b); res != false {
		t.Fatal(res)
	}

	if res := MustGet(EqTreeI(b, EqT2[int](), EqT2[string]()), b); res != true {
		t.Fatal(res)
	}

}

func TestDiffTreeRoot(t *testing.T) {

	before := Tree(
		"alpha", //[]
		ValCol(
			Tree(
				"beta", //[0]
				ValCol(
					Tree("gamma"), //[0 0]
					Tree(
						"delta", //[0 1]
						ValCol(
							Tree("epsilon"), //[0 1 0]
							Tree("zeta"),    //[0 1 1]
							Tree("eta"),     //[0 1 2]
						),
					),
				),
			),
			Tree("theta"), //[1]
		),
	)

	after := Tree(
		"Alpha", //0
		ValCol(
			Tree(
				"beta", //[0]
				ValCol(
					Tree("gamma"), //[0 0]
					Tree(
						"delta", //[0 1]
						ValCol(
							Tree("epsilon"), //[0 1 0]
							Tree("zeta"),    //[0 1 1]
							Tree("eta"),     //[0 1 2]
						),
					),
				),
			),
			Tree("theta"), //[1]
		),
	)

	dist := DistancePercent(
		EditDistance(Compose(TreeValue[int, string, Pure](), TraverseString()), EditLevenshtein, EqT2[rune](), 10),
		Compose(TreeValue[int, string, Pure](), Length(TraverseString())),
	)

	diffType := FieldLens(func(source *Diff[*PathNode[int], TreeNode[int, string, Pure]]) *DiffType { return &source.Type })

	a := Compose(Non[TreeNode[int, string, Pure]](TreeNode[int, string, Pure]{}, EqDeepT2[TreeNode[int, string, Pure]]()), TreeValue[int, string, Pure]())

	c := ComposeLeft(DiffTreeT2(TreeChildren[int, string, Pure](), 0.5, dist, EqT2[int](), DiffNone, true), a)

	diff := ReIndexed(c, diffType, EqT2[DiffType]())

	if res := MustGet(SliceOf(WithIndex(diff), 10), lo.T2(mo.Some(after), mo.Some(before))); !MustGet(EqDeepT2[[]ValueI[DiffType, string]](), lo.T2(res, []ValueI[DiffType, string]{
		ValI(
			DiffModify,
			"Alpha",
		),
	})) {
		t.Fatal("iterate", res)
	}

	fmapSeq := 0
	expected := Tree(
		"Alpha:0",
		ValCol(
			Tree(
				"beta", //[0]
				ValCol(
					Tree("gamma"), //[0 0]
					Tree(
						"delta", //[0 1]
						ValCol(
							Tree("epsilon"), //[0 1 0]
							Tree("zeta"),    //[0 1 1]
							Tree("eta"),     //[0 1 2]
						),
					),
				),
			),
			Tree("theta"), //[1]
		),
	)

	if res := MustModify(diff, Op(func(focus string) string {
		if focus == "" {
			return ""
		}

		ret := fmt.Sprintf("%v:%v", focus, fmapSeq)
		fmapSeq++
		return ret
	}), lo.T2(mo.Some(after), mo.Some(before))); !MustGet(EqTreeI(res.A.MustGet(), EqT2[int](), EqT2[string]()), expected) {
		t.Fatal("modify", res.A, "expected", expected)
	}
}

func TestDiffTreeRootDelete(t *testing.T) {

	tree := Tree(
		"alpha", //[]
		ValCol(
			Tree(
				"beta", //[0]
				ValCol(
					Tree("gamma"), //[0 0]
					Tree(
						"delta", //[0 1]
						ValCol(
							Tree("epsilon"), //[0 1 0]
							Tree("zeta"),    //[0 1 1]
							Tree("eta"),     //[0 1 2]
						),
					),
				),
			),
			Tree("theta"), //[1]
		),
	)

	dist := DistancePercent(
		EditDistance(Compose(TreeValue[int, string, Pure](), TraverseString()), EditLevenshtein, EqT2[rune](), 10),
		Compose(TreeValue[int, string, Pure](), Length(TraverseString())),
	)

	diffType := FieldLens(func(source *Diff[*PathNode[int], TreeNode[int, string, Pure]]) *DiffType { return &source.Type })

	diffTree := DiffTreeT2(TreeChildren[int, string, Pure](), 0.5, dist, EqT2[int](), DiffNone, true)

	optic := ReIndexed(diffTree, diffType, EqT2[DiffType]())

	res := MustGet(SliceOf(WithIndex(optic), 10), lo.T2(mo.Some(tree), mo.None[TreeNode[int, string, Pure]]()))

	diff, optNode := res[0].Get()
	if diff != DiffAdd {
		t.Fatal(diff, optNode)
	}

	res = MustGet(SliceOf(WithIndex(optic), 10), lo.T2(mo.None[TreeNode[int, string, Pure]](), mo.Some(tree)))

	diff, optNode = res[0].Get()
	if diff != DiffRemove {
		t.Fatal(diff, optNode)
	}

}

func TestTraverseTreeNodePIxGet(t *testing.T) {

	tree := Tree(
		"alpha", //[]
		ValCol(
			Tree(
				"beta", //[0]
				ValCol(
					Tree("gamma"), //[0 0]
					Tree(
						"delta", //[0 1]
						ValCol(
							Tree("epsilon"), //[0 1 0]
							Tree("zeta"),    //[0 1 1]
							Tree("eta"),     //[0 1 2]
						),
					),
				),
			),
			Tree("theta"), //[1]
		),
	)

	res := MustGet(
		SliceOf(
			Index(
				TraverseTreeChildrenP[int, string, string, Pure](EqT2[int]()),
				Path(0, 1),
			),
			1,
		),
		tree,
	)

	if !reflect.DeepEqual(
		res,
		[]string{"delta"},
	) {
		t.Fatal(res)
	}

}
