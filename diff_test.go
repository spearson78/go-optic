package optic_test

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/samber/lo"

	. "github.com/spearson78/go-optic"
)

func ExampleDiffColT2I() {

	before := ValCol("alpha", "beta", "gamma", "delta")
	after := ValCol("gamma", "Alpha", "epsilon", "alpha")

	optic := DiffColT2I[int](
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
		EqT2[int](),
		DiffNone,
		true,
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

func ExampleDiffCol() {

	before := ValCol("alpha", "beta", "gamma", "delta")
	after := ValCol("Alpha", "gamma", "beta", "epsilon")

	optic := DiffCol[int](
		before,
		2,
		EPure(EditDistance(TraverseString(), EditLevenshtein, EqT2[rune](), 4)),
		EqT2[int](),
		DiffNone,
		true,
	)

	for i, v := range MustGet(SeqIOf(optic), after) {
		fmt.Println(i, "After Value", v)
	}

	//Output:
	//DiffModify(Pos=0->0,Index=0->0,BeforeValue=alpha,Dist=1) After Value {true Alpha}
	//DiffModify(Pos=2->1,Index=2->1,BeforeValue=gamma,Dist=0) After Value {true gamma}
	//DiffModify(Pos=1->2,Index=1->2,BeforeValue=beta,Dist=0) After Value {true beta}
	//DiffRemove(Index=3,Value=delta) After Value {false }
	//DiffAdd(Index=3) After Value {true epsilon}

}

func ExampleDiffColT2() {

	before := ValCol("alpha", "beta", "gamma", "delta")
	after := ValCol("Alpha", "gamma", "beta", "epsilon")

	optic := DiffColT2[int](2, EPure(EditDistance(TraverseString(), EditLevenshtein, EqT2[rune](), 4)), EqT2[int](), DiffNone, true)

	for i, v := range MustGet(SeqIOf(optic), lo.T2(after, before)) {
		fmt.Println(i, "After Value", v)
	}

	//Output:
	//DiffModify(Pos=0->0,Index=0->0,BeforeValue=alpha,Dist=1) After Value {true Alpha}
	//DiffModify(Pos=2->1,Index=2->1,BeforeValue=gamma,Dist=0) After Value {true gamma}
	//DiffModify(Pos=1->2,Index=1->2,BeforeValue=beta,Dist=0) After Value {true beta}
	//DiffRemove(Index=3,Value=delta) After Value {false }
	//DiffAdd(Index=3) After Value {true epsilon}

}

func ExampleDistancePercent() {

	before := ValCol("alpha", "beta", "gamma", "delta")
	after := ValCol("Alpha", "gamma", "beta", "epsilon")

	optic := DiffColT2[int](
		0.5,
		DistancePercent(
			EditDistance(TraverseString(), EditLevenshtein, EqT2[rune](), 4),
			Length(TraverseString()),
		),
		EqT2[int](),
		DiffNone,
		true,
	)

	for i, v := range MustGet(SeqIOf(optic), lo.T2(after, before)) {
		fmt.Println(i, "After Value", v)
	}

	//Output:
	//DiffModify(Pos=0->0,Index=0->0,BeforeValue=alpha,Dist=0.2) After Value {true Alpha}
	//DiffModify(Pos=2->1,Index=2->1,BeforeValue=gamma,Dist=0) After Value {true gamma}
	//DiffModify(Pos=1->2,Index=1->2,BeforeValue=beta,Dist=0) After Value {true beta}
	//DiffRemove(Index=3,Value=delta) After Value {false }
	//DiffAdd(Index=3) After Value {true epsilon}
}

func ExampleDistancePercentI() {

	before := ValCol("alpha", "beta", "gamma", "delta")
	after := ValCol("gamma", "Alpha", "epsilon", "alpha")

	optic := DiffColT2I[int](
		0.5,
		EPure(DistancePercentI(
			AddOp(
				Compose(
					DelveT2(ValueIValue[int, string]()),
					EditDistance(TraverseString(), EditLevenshtein, EqT2[rune](), 4),
				),
				Compose4(
					DelveT2(ValueIIndex[int, string]()),
					SubT2[int](),
					Abs[int](),
					IsoCast[int, float64](),
				),
			),
			Compose(Length(TraverseString()), Mul(2)),
		)),
		EqT2[int](),
		DiffNone,
		true,
	)

	for index, val := range MustGet(SeqIOf(optic), lo.T2(after, before)) {
		fmt.Println(index, "After Value", val)
	}

	//Output:
	//DiffModify(Pos=2->0,Index=2->0,BeforeValue=gamma,Dist=0.2) After Value {true gamma}
	//DiffRemove(Index=1,Value=beta) After Value {false }
	//DiffModify(Pos=0->1,Index=0->1,BeforeValue=alpha,Dist=0.2) After Value {true Alpha}
	//DiffAdd(Index=2) After Value {true epsilon}
	//DiffModify(Pos=3->3,Index=3->3,BeforeValue=delta,Dist=0.4) After Value {true alpha}

}

func ExampleDistance() {

	before := ValCol("alpha", "beta", "gamma", "delta")
	after := ValCol("Alpha", "gamma", "beta", "epsilon")

	optic := DiffColT2[int](
		2,
		Distance(func(a, b string) float64 {
			lenA := len(a)
			lenB := len(b)

			dist := math.Abs(float64(lenA - lenB))

			minLen := min(lenA, lenB)
			for i := 0; i < minLen; i++ {
				if a[i] != b[i] {
					dist = dist + 1
				}
			}

			return dist
		}),
		EqT2[int](),
		DiffNone,
		true,
	)

	for i, v := range MustGet(SeqIOf(optic), lo.T2(after, before)) {
		fmt.Println(i, "After Value", v)
	}

	//Output:
	//DiffModify(Pos=0->0,Index=0->0,BeforeValue=alpha,Dist=1) After Value {true Alpha}
	//DiffModify(Pos=2->1,Index=2->1,BeforeValue=gamma,Dist=0) After Value {true gamma}
	//DiffModify(Pos=1->2,Index=1->2,BeforeValue=beta,Dist=0) After Value {true beta}
	//DiffRemove(Index=3,Value=delta) After Value {false }
	//DiffAdd(Index=3) After Value {true epsilon}
}

func ExampleDistanceI() {

	before := ValCol("alpha", "beta", "gamma", "delta")
	after := ValCol("gamma", "Alpha", "epsilon", "alpha")

	optic := DiffColT2I[int](
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
		EqT2[int](),
		DiffNone,
		true,
	)

	for i, v := range MustGet(SeqIOf(optic), lo.T2(after, before)) {
		fmt.Println(i, "After Value", v)
	}

	//Output:
	//DiffModify(Pos=2->0,Index=2->0,BeforeValue=gamma,Dist=2) After Value {true gamma}
	//DiffRemove(Index=1,Value=beta) After Value {false }
	//DiffModify(Pos=0->1,Index=0->1,BeforeValue=alpha,Dist=2) After Value {true Alpha}
	//DiffAdd(Index=2) After Value {true epsilon}
	//DiffRemove(Index=3,Value=delta) After Value {false }
	//DiffAdd(Index=3) After Value {true alpha}
}

func ExampleDistanceE() {

	before := ColErr(ValCol(
		"alpha",
		"beta",
		"gamma",
		"delta",
	))
	after := ColErr(ValCol(
		"Alpha",
		"gamma",
		"beta",
		"epsilon",
	))

	optic := DiffColT2[int](
		2,
		DistanceE[int](func(ctx context.Context, a, b string) (float64, error) {
			lenA := len(a)
			lenB := len(b)

			dist := math.Abs(float64(lenA - lenB))

			minLen := min(lenA, lenB)
			for i := 0; i < minLen; i++ {
				if a[i] != b[i] {
					dist = dist + 1
				}
			}

			return dist, nil
		}),
		EqT2[int](),
		DiffNone,
		true,
	)

	seq, err := Get(SeqIEOf(optic), lo.T2(after, before))
	if err != nil {
		log.Fatal(err)
	}
	for val := range seq {
		i, v, err := val.Get()
		fmt.Println(i, "After Value", v, "Err:", err)
	}

	//Output:
	//DiffModify(Pos=0->0,Index=0->0,BeforeValue=alpha,Dist=1) After Value {true Alpha} Err: <nil>
	//DiffModify(Pos=2->1,Index=2->1,BeforeValue=gamma,Dist=0) After Value {true gamma} Err: <nil>
	//DiffModify(Pos=1->2,Index=1->2,BeforeValue=beta,Dist=0) After Value {true beta} Err: <nil>
	//DiffRemove(Index=3,Value=delta) After Value {false } Err: <nil>
	//DiffAdd(Index=3) After Value {true epsilon} Err: <nil>
}

func ExampleDistanceIE() {

	before := ColErr(ValCol("alpha", "beta", "gamma", "delta"))
	after := ColErr(ValCol("gamma", "Alpha", "epsilon", "alpha"))

	optic := DiffColT2I[int](
		2,
		DistanceIE(func(ctx context.Context, indexA int, a string, indexB int, b string) (float64, error) {
			lenA := len(a)
			lenB := len(b)

			dist := math.Abs(float64(lenA-lenB)) + math.Abs(float64(indexA-indexB))

			minLen := min(lenA, lenB)
			for i := 0; i < minLen; i++ {
				if a[i] != b[i] {
					dist = dist + 1
				}
			}

			return dist, nil
		}),
		EqT2[int](),
		DiffNone,
		true,
	)

	seq, err := Get(SeqIEOf(optic), lo.T2(after, before))
	if err != nil {
		log.Fatal(err)
	}
	for val := range seq {
		i, v, err := val.Get()
		fmt.Println(i, "After Value", v, "Err:", err)
	}

	//Output:
	//DiffModify(Pos=2->0,Index=2->0,BeforeValue=gamma,Dist=2) After Value {true gamma} Err: <nil>
	//DiffRemove(Index=1,Value=beta) After Value {false } Err: <nil>
	//DiffModify(Pos=0->1,Index=0->1,BeforeValue=alpha,Dist=2) After Value {true Alpha} Err: <nil>
	//DiffAdd(Index=2) After Value {true epsilon} Err: <nil>
	//DiffRemove(Index=3,Value=delta) After Value {false } Err: <nil>
	//DiffAdd(Index=3) After Value {true alpha} Err: <nil>
}
