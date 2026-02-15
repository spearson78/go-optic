package optic_test

import (
	"fmt"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
)

func ExampleSwappedT2() {

	t2 := lo.T2("one", 1)

	var viewResult lo.Tuple2[int, string] = MustGet(SwappedT2[string, int](), t2)
	fmt.Println(viewResult)

	//Note we set a swapped value and the lens converts it back to the original order
	var setResult lo.Tuple2[string, int] = MustSet(SwappedT2[string, int](), lo.T2(2, "two"), t2)
	fmt.Println(setResult)

	//Output:{1 one}
	//{two 2}
}

func ExampleBesideT2() {

	data1 := []int{10, 20, 30}
	data2 := map[string]int{
		"alpha": 40,
		"beta":  50,
		"gamma": 60,
		"delta": 70,
	}
	data := lo.T2(data1, data2)

	optic1 := TraverseSlice[int]()
	optic2 := TraverseMap[string, int]()

	beside := BesideT2(optic1, optic2)

	var listRes []int = MustGet(SliceOf(beside, 4), data)
	fmt.Println(listRes)

	var overResult lo.Tuple2[[]int, map[string]int] = MustModify(beside, Mul(2), data)
	fmt.Println(overResult)

	//Output:
	//[10 20 30 40 50 70 60]
	//{[20 40 60] map[alpha:80 beta:100 delta:140 gamma:120]}

}
