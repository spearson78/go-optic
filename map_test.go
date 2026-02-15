package optic_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	"github.com/samber/mo"

	. "github.com/spearson78/go-optic"
)

func FuzzTraverseMap(f *testing.F) {
	f.Add(uint8(0), int64(0), int(10))
	f.Add(uint8(1), int64(1), int(20))
	f.Add(uint8(2), int64(2), int(30))

	f.Fuzz(func(t *testing.T, len uint8, seed int64, newVal int) {

		r := rand.New(rand.NewSource(seed))
		data := make(map[int]int, len)
		for i := byte(0); i < len; i++ {
			data[r.Int()] = r.Int()
		}

		ValidateOpticTestPred(t, TraverseMap[int, int](), data, newVal, EqDeepT2[map[int]int]())

	})
}

func ExampleTraverseMap() {

	mapData := map[string]int{
		"a": 5,
		"b": 10,
		"c": 15,
	}

	//Technically multiple elements of an Traversal may have the same index so Index returns multiple values. FirstOf selects the first one.
	elemOf, found := MustGetFirst(Index(TraverseMap[string, int](), "b"), mapData)
	fmt.Println(elemOf, found)

	result := MustModify(TraverseMap[string, int](), Mul(10), mapData)
	fmt.Println(result)

	//Output:10 true
	//map[a:50 b:100 c:150]
}

func ExampleTraverseMapP() {

	mapData := map[string]string{
		"a": "5",
		"b": "10",
		"c": "15",
	}

	//Not the result is a map[string]int not map[string]string
	var overResult map[string]int = MustModify(TraverseMapP[string, string, int](), Op(func(focus string) int {
		i, _ := strconv.ParseInt(focus, 10, 32)
		return int(i)
	}), mapData)
	fmt.Println(overResult)

	//Output:map[a:5 b:10 c:15]
}

func TestMapOfReducedHistogram(t *testing.T) {

	data := []string{"I", "know", "what", "I", "like", "and", "I", "like", "what", "I", "know"}

	//Word mapper converts the list int a list of maps with a single entry with the string as key and 1 as the value.
	wordMapper := MapOfReduced(
		ComposeLeft(
			SelfIndex( //Move the word into the index
				TraverseSlice[string](),
				EqT2[string](),
			),
			Const[string](1), //Set the value of each word to 1
		),
		Sum[int](), //Sum the counts of each word.
		5,
	)

	var wordMap map[string]int = MustGet(wordMapper, data)

	if !reflect.DeepEqual(map[string]int{
		"I":    4,
		"and":  1,
		"know": 2,
		"like": 2,
		"what": 2,
	}, wordMap) {
		t.Fatal(wordMap)
	}
}

func ExampleMapToCol() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	var getRes map[int]string = MustGet(
		Compose3(
			MapToCol[int, string](),
			FilteredCol[int](Ne("beta")),
			AsReverseGet(MapToCol[int, string]()),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes map[int]string = MustModify(
		MapToCol[int, string](),
		FilteredCol[int](Ne("beta")),
		data,
	)
	fmt.Println(modifyRes)

	// Output:
	//map[1:alpha 3:gamma 4:delta]
	//map[1:alpha 3:gamma 4:delta]
}

func TestMapMetrics(t *testing.T) {

	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
		"delta": 4,
	}

	var metrics Metrics
	optic := WithMetrics(TraverseMap[string, int](), &metrics)

	if index, viewResult, found := MustGetFirstI(optic, data); index != "alpha" || viewResult != 1 || found != true || metrics.Access != 1 || metrics.Focused != 1 || metrics.Custom["comparisons"] > 6 {
		t.Fatalf("IPreViewContext metrics : %v %v %v %v", index, viewResult, found, metrics)
	}

}

func ExampleMapOf() {

	data := []string{
		"alpha",
		"beta",
		"gamma",
		"delta",
	}

	optic := MapOf(TraverseSlice[string](), len(data))

	var viewResult map[int]string
	viewResult, err := Get(optic, data)
	fmt.Println(viewResult, err)

	var modifyResult []string
	modifyResult, err = Modify(optic, Op(func(focus map[int]string) map[int]string {
		focus[1] = "beta-test"
		return focus
	}), data)
	fmt.Println(modifyResult, err)

	//Output:
	//map[0:alpha 1:beta 2:gamma 3:delta] <nil>
	//[alpha beta-test gamma delta] <nil>
}

func ExampleMapOfReduced() {

	data := ValColI(
		IxMatchComparable[string](),
		ValI("alpha", 1),
		ValI("beta", 2),
		ValI("alpha", 3),
	)

	optic := MapOfReduced(TraverseCol[string, int](), FirstReducer[int](), 3)

	var viewResult map[string]int
	viewResult, err := Get(optic, data)
	fmt.Println(viewResult, err)

	var modifyResult Collection[string, int, Pure]
	modifyResult, err = Modify(optic, Op(func(focus map[string]int) map[string]int {
		focus["beta"] = 10
		return focus
	}), data)
	fmt.Println(modifyResult, err)

	//Output:
	//map[alpha:1 beta:2] <nil>
	//Col[alpha:1 beta:10 alpha:1] <nil>
}

func ExampleMapOfP() {

	data := []string{
		"1",
		"2",
		"3",
		"4",
	}

	optic := MapOfP(TraverseSliceP[string, int](), len(data))

	var viewResult map[int]string
	viewResult, err := Get(optic, data)
	fmt.Println(viewResult, err)

	var modifyResult []int
	modifyResult, err = Modify(optic, Op(func(focus map[int]string) map[int]int {
		ret := make(map[int]int)

		for k, v := range focus {
			intVal, _ := strconv.ParseInt(v, 10, 32)
			ret[k] = int(intVal*10) + k
		}

		return ret
	}), data)
	fmt.Println(modifyResult, err)

	//Output:
	//map[0:1 1:2 2:3 3:4] <nil>
	//[10 21 32 43] <nil>
}
func ExampleMapOfReducedP() {

	data := ValColI(
		IxMatchComparable[string](),
		ValI("alpha", 1),
		ValI("beta", 2),
		ValI("alpha", 3),
	)

	optic := MapOfReducedP(TraverseColP[string, int, string](), FirstReducer[int](), 3)

	var viewResult map[string]int
	viewResult, err := Get(optic, data)
	fmt.Println(viewResult, err)

	var modifyResult Collection[string, string, Pure]
	modifyResult, err = Modify(optic, Op(func(focus map[string]int) map[string]string {
		ret := make(map[string]string)

		for k, v := range focus {
			strVal := strconv.Itoa(v)
			ret[k] = strVal
		}

		return ret
	}), data)
	fmt.Println(modifyResult, err)

	//Output:
	//map[alpha:1 beta:2] <nil>
	//Col[alpha:1 beta:2 alpha:1] <nil>
}

func ExampleAtMap() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	optic := AtMap[int, string](5)

	var viewResult mo.Option[string] = MustGet(optic, data)
	fmt.Println(viewResult)

	var modifyResult map[int]string = MustSet(optic, mo.Some("epsilon"), data)
	fmt.Println(modifyResult)

	var deleteResult map[int]string = MustSet(AtMap[int, string](2), mo.None[string](), data)
	fmt.Println(deleteResult)

	//At can be composed with Non to provide a default value.

	defaultValueOptic := Compose(AtMap[int, string](6), Non("default", EqT2[string]()))

	var defaultViewResult string = MustGet(defaultValueOptic, data)
	fmt.Println(defaultViewResult)

	//Output:
	//{false }
	//map[1:alpha 2:beta 3:gamma 4:delta 5:epsilon]
	//map[1:alpha 3:gamma 4:delta]
	//default
}
