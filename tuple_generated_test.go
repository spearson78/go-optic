package optic

import (
	"fmt"

	"github.com/samber/lo"
)

func ExampleT2A() {

	t2 := lo.T2("one", 1)

	var viewResult string = MustGet(T2A[string, int](), t2)
	fmt.Println(viewResult)

	var setResult lo.Tuple2[string, int] = MustSet(T2A[string, int](), "two", t2)
	fmt.Println(setResult)

	//Output:one
	//{two 1}
}
func ExampleT2AP() {

	t2 := lo.T2("one", 1)

	var viewResult string = MustGet(T2AP[string, int, int](), t2)
	fmt.Println(viewResult)

	var setResult lo.Tuple2[int, int] = MustSet(T2AP[string, int, int](), 2, t2)
	fmt.Println(setResult)

	//Output:one
	//{2 1}
}
func ExampleT2B() {

	t2 := lo.T2("one", 1)

	var viewResult int = MustGet(T2B[string, int](), t2)
	fmt.Println(viewResult)

	var setResult lo.Tuple2[string, int] = MustSet(T2B[string, int](), 2, t2)
	fmt.Println(setResult)

	//Output:1
	//{one 2}
}
func ExampleT2BP() {

	t2 := lo.T2("one", 1)

	var viewResult int = MustGet(T2BP[string, int, string](), t2)
	fmt.Println(viewResult)

	var setResult lo.Tuple2[string, string] = MustSet(T2BP[string, int, string](), "two", t2)
	fmt.Println(setResult)

	//Output:1
	//{one two}
}
func ExampleT2Of() {

	type Person struct {
		Name string
		Age  int
	}

	personName := FieldLens(func(p *Person) *string { return &p.Name })
	personAge := FieldLens(func(p *Person) *int { return &p.Age })

	optic := T2Of(personName, personAge)

	var res lo.Tuple2[string, int] = MustGet(optic, Person{Name: "Max Mustermann", Age: 46})
	fmt.Println(res)

	var modifyRes Person = MustModify(Compose(optic, T2B[string, int]()), Add(1), Person{Name: "Max Mustermann", Age: 46})
	fmt.Println("Name:", modifyRes.Name)
	fmt.Println("Age:", modifyRes.Age)

	// Output:
	// {Max Mustermann 46}
	// Name: Max Mustermann
	// Age: 47
}
func ExampleDupT2() {

	data := []int{10, 20, 30}

	//This optic takes slice of single ints and converts it to slice of lo.Tuple2[int,int] where A and B have duplicated values.
	optic := Compose(TraverseSlice[int](), DupT2[int]())

	var listRes []lo.Tuple2[int, int] = MustGet(SliceOf(optic, 3), data)
	fmt.Println(listRes)

	compositeReducer := ReducerT2(Sum[int](), Product[int]())
	//MustReduce applies both Reducers at the same time and returns a lo.Tuple2 with the aggregated results.
	var res lo.Tuple2[int, int]
	var ok bool
	res, ok = MustGetFirst(Reduce(optic, compositeReducer), data)

	fmt.Println(res, ok)

	//Output:
	//[{10 10} {20 20} {30 30}]
	//{60 6000} true
}
func ExampleDupT3() {

	data := []int{10, 20, 30}

	//This optic takes slice of single ints and converts it to slice of lo.Tuple2[int,int] where A and B have duplicated values.
	optic := Compose(TraverseSlice[int](), DupT3[int]())

	var listRes []lo.Tuple3[int, int, int] = MustGet(SliceOf(optic, 3), data)
	fmt.Println(listRes)

	compositeReducer := ReducerT3(Sum[int](), Product[int](), Median[int]())
	//MustReduce applies both Reducers at the same time and returns a lo.Tuple2 with the aggregated results.
	var res lo.Tuple3[int, int, int]
	var ok bool
	res, ok = MustGetFirst(Reduce(optic, compositeReducer), data)

	fmt.Println(res, ok)

	//Output:
	//[{10 10 10} {20 20 20} {30 30 30}]
	//{60 6000 20} true
}
func ExampleTraverseT2() {

	data := lo.T2[int, int](1, 2)

	optic := TraverseT2[int]()

	var viewResult []int = MustGet(SliceOf(optic, 2), data)
	fmt.Println(viewResult)

	var overResult lo.Tuple2[int, int] = MustModify(optic, Add(10), data)
	fmt.Println(overResult)

	//Output:
	//[1 2]
	//{11 12}
}
func ExampleTraverseT2P() {

	data := lo.T2[string, string]("1", "2")

	optic := Compose(TraverseT2P[string, int](), ParseIntP[int](10, 32))

	var viewResult []int
	var err error
	viewResult, err = Get(SliceOfP(optic, 2), data)
	fmt.Println(viewResult, err)

	var overResult lo.Tuple2[int, int]
	overResult, err = Modify(optic, Add(10), data)
	fmt.Println(overResult, err)

	//Output:
	//[1 2] <nil>
	//{11 12} <nil>
}
