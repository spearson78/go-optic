package optic_test

import (
	"context"
	"fmt"
	"iter"
	"log"
	"strings"

	. "github.com/spearson78/go-optic"
)

func ExamplePull() {

	col := ValCol(10, 20, 30, 40, 50).AsIter()(context.Background())

	next, stop := PullIE(col)
	defer stop()

	for {
		index, focus, err, ok := next()
		if !ok {
			break
		}
		fmt.Println(index, focus, err, ok)
	}

	//Output:
	//0 10 <nil> true
	//1 20 <nil> true
	//2 30 <nil> true
	//3 40 <nil> true
	//4 50 <nil> true

}

func ExampleSeqOf() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//This optic focuses each string in the slice.
	optic := SeqOf(TraverseSlice[string]())

	for focus := range MustGet(optic, data) {
		fmt.Println(focus)
	}

	var modifyResult []string = MustModify(optic, Op(func(focus iter.Seq[string]) iter.Seq[string] {
		return func(yield func(string) bool) {
			focus(func(str string) bool {
				return yield(strings.ToUpper(str))
			})
		}
	}), data)
	fmt.Println(modifyResult)

	//Output:
	//alpha
	//beta
	//gamma
	//delta
	//[ALPHA BETA GAMMA DELTA]
}

func ExampleSeqOfP() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	//Polymorphic map traversal with a source type of map[int]string and result type of map[int]int
	optic := SeqOfP(TraverseMapP[int, string, int]())

	var viewResult iter.Seq[string]
	var err error

	viewResult, err = Get(optic, data)
	if err != nil {
		log.Fatal(err)
	}

	for v := range viewResult {
		fmt.Print(v, " ")
	}
	fmt.Println(".")

	//Note the return type is map[int]int not map[int]string
	var modifyResult map[int]int
	modifyResult, err = Modify(optic, Op(func(focus iter.Seq[string]) iter.Seq[int] {
		return func(yield func(int) bool) {
			focus(func(str string) bool {
				return yield(len(str))
			})
		}
	}), data)
	fmt.Println(modifyResult, err)

	//Output:
	//alpha beta gamma delta .
	//map[1:5 2:4 3:5 4:5] <nil>
}

func ExampleSeqIOf() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//This optic focuses each string in the slice.
	optic := SeqIOf(TraverseSlice[string]())

	for index, focus := range MustGet(optic, data) {
		fmt.Println(index, focus)
	}

	var modifyResult []string = MustModify(optic, Op(func(focus SeqI[int, string]) SeqI[int, string] {
		return func(yield func(int, string) bool) {
			focus(func(i int, str string) bool {
				return yield(i, strings.ToUpper(str))
			})
		}
	}), data)
	fmt.Println(modifyResult)

	//Output:
	//0 alpha
	//1 beta
	//2 gamma
	//3 delta
	//[ALPHA BETA GAMMA DELTA]
}

func ExampleSeqIOfP() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	//Polymorphic map traversal with a source type of map[int]string and result type of map[int]int
	optic := SeqIOfP(TraverseMapP[int, string, int]())

	var viewResult SeqI[int, string]
	var err error
	viewResult, err = Get(optic, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Result ")
	for i, v := range viewResult {
		fmt.Print(i, ":", v, " ")
	}
	fmt.Println(".")

	//Note the return type is map[int]int not map[int]string
	var modifyResult map[int]int
	modifyResult, err = Modify(optic, Op(func(focus SeqI[int, string]) SeqI[int, int] {
		return func(yield func(int, int) bool) {
			focus(func(index int, str string) bool {
				return yield(index, len(str))
			})
		}
	}), data)
	fmt.Println(modifyResult, err)

	//Output:
	//Result 1:alpha 2:beta 3:gamma 4:delta .
	//map[1:5 2:4 3:5 4:5] <nil>
}

func ExampleSeqEOf() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	optic := SeqEOf(TraverseMap[int, string]())

	viewResult, err := Get(optic, data)
	if err != nil {
		log.Fatal(err)
	}

	for val := range viewResult {
		v, err := val.Get()
		fmt.Printf("%v:%v ", v, err)
	}
	fmt.Println(".")

	modifyResult, err := Modify(optic, Op(func(focus SeqE[string]) SeqE[string] {
		return func(yield func(ValueE[string]) bool) {
			focus(func(val ValueE[string]) bool {
				return yield(ValE(strings.ToUpper(val.Value()), val.Error()))
			})
		}
	}), data)
	fmt.Println(modifyResult, err)

	//Output:
	//alpha:<nil> beta:<nil> gamma:<nil> delta:<nil> .
	//map[1:ALPHA 2:BETA 3:GAMMA 4:DELTA] <nil>
}

func ExampleSeqEOfP() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	//Polymorphic map traversal with a source type of map[int]string and result type of map[int]int
	optic := SeqEOfP(TraverseMapP[int, string, int]())

	var viewResult SeqE[string]
	var err error
	viewResult, err = Get(optic, data)
	if err != nil {
		log.Fatal(err)
	}
	for val := range viewResult {
		v, err := val.Get()
		fmt.Print(v, ":", err, " ")
	}
	fmt.Println(".")

	//Note the return type is map[int]int not map[int]string
	var modifyResult map[int]int
	modifyResult, err = Modify(optic, Op(func(focus SeqE[string]) SeqE[int] {
		return func(yield func(ValueE[int]) bool) {
			focus(func(val ValueE[string]) bool {
				return yield(ValE(len(val.Value()), val.Error()))
			})
		}
	}), data)
	fmt.Println(modifyResult, err)

	//Output:
	//alpha:<nil> beta:<nil> gamma:<nil> delta:<nil> .
	//map[1:5 2:4 3:5 4:5] <nil>
}

func ExampleSeqIEOf() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	optic := SeqIEOf(TraverseMap[int, string]())

	var viewResult SeqIE[int, string] = MustGet(optic, data)
	for valIE := range viewResult {
		i, v, e := valIE.Get()
		fmt.Println(i, v, e)
	}

	modifyResult := MustModify(optic, Op(func(focus SeqIE[int, string]) SeqIE[int, string] {
		return func(yield func(ValueIE[int, string]) bool) {
			focus(func(val ValueIE[int, string]) bool {
				index, str, err := val.Get()
				return yield(ValIE(index, strings.ToUpper(str), err))
			})
		}
	}), data)
	fmt.Println(modifyResult)

	//Output:
	//1 alpha <nil>
	//2 beta <nil>
	//3 gamma <nil>
	//4 delta <nil>
	//map[1:ALPHA 2:BETA 3:GAMMA 4:DELTA]
}

func ExampleSeqIEOfP() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	//Polymorphic map traversal with a source type of map[int]string and result type of map[int]int
	optic := SeqIEOfP(TraverseMapP[int, string, int]())

	var viewResult SeqIE[int, string]
	var err error
	viewResult, err = Get(optic, data)
	if err != nil {
		log.Fatal(err)
	}
	for valIE := range viewResult {
		i, v, e := valIE.Get()
		fmt.Println(i, v, e)
	}

	//Note the return type is map[int]int not map[int]string
	var modifyResult map[int]int
	modifyResult, err = Modify(optic, Op(func(focus SeqIE[int, string]) SeqIE[int, int] {
		return func(yield func(ValueIE[int, int]) bool) {
			focus(func(val ValueIE[int, string]) bool {
				index, str, err := val.Get()
				return yield(ValIE(index, len(str), err))
			})
		}
	}), data)
	fmt.Println(modifyResult, err)

	//Output:
	//1 alpha <nil>
	//2 beta <nil>
	//3 gamma <nil>
	//4 delta <nil>
	//map[1:5 2:4 3:5 4:5] <nil>
}
