package optic

import (
	"fmt"
	"math"

	"github.com/samber/lo"
)

func ExampleColTypeOp() {
	sliceCol := SliceColTypeP[string, string]()

	//ColOp converts a col op into a col op for our specific collection
	stringReversed := ColTypeOp(sliceCol, ReversedCol[int, string]())

	result := MustGet(stringReversed, []string{"alpha", "beta", "gamma", "delta"})

	fmt.Println(result)
	//Output:
	//[delta gamma beta alpha]
}

func ExampleDiffColTypeT2() {

	diff := Compose(
		AsModify(
			TraverseT2P[int, float64](),
			UpCast[int, float64](),
		),
		SubT2[float64](),
	)

	sliceDiff := DiffColTypeT2(
		SliceColType[int](),
		1,
		EPure(diff),
		DiffNone,
	)

	result := MustGet(
		SliceOf(WithIndex(sliceDiff), 5),
		lo.T2(
			[]int{10, 20, 30},
			[]int{10, 21},
		),
	)
	fmt.Println(result)

	//Output:
	//[DiffModify(Pos=1->1,Index=1->1,BeforeValue=21,Dist=1):{true 20} DiffAdd(Index=2):{true 30}]
}

func ExampleDiffColTypeT2I() {

	before := []string{"alpha", "beta", "gamma", "delta"}
	after := []string{"gamma", "Alpha", "epsilon", "alpha"}

	optic := DiffColTypeT2I[int](
		SliceColType[string](),
		2,
		DistanceI(func(indexA int, a string, indexB int, b string) float64 {
			lenA := len(a)
			lenB := len(b)

			dist := math.Abs(float64(lenA-lenB)) + math.Abs(float64(indexA-indexB))

			minLen := min(lenA, lenB)
			for i := 0; i < minLen; i++ {
				if a[i] != b[i] {
					dist = dist + 1
				}
			}

			return dist
		}),
		DiffNone,
	)

	for diffType, value := range MustGet(SeqIOf(optic), lo.T2(after, before)) {

		fmt.Println(diffType, "After Value", value)
	}

	//Output:
	//DiffModify(Pos=2->0,Index=2->0,BeforeValue=gamma,Dist=2) After Value {true gamma}
	//DiffRemove(Index=1,Value=beta) After Value {false }
	//DiffModify(Pos=0->1,Index=0->1,BeforeValue=alpha,Dist=2) After Value {true Alpha}
	//DiffAdd(Index=2) After Value {true epsilon}
	//DiffRemove(Index=3,Value=delta) After Value {false }
	//DiffAdd(Index=3) After Value {true alpha}
}
