package optic_test

import (
	"context"
	"fmt"
	"strconv"

	. "github.com/spearson78/go-optic"
)

func ExampleOrderByOp() {
	data := []int{1, 10, 2, 9, 3, 8, 4, 7, 5, 6}

	optic := Ordered(
		TraverseSlice[int](),
		OrderByOp(
			func(left, right int) bool {
				return left < right
			},
		),
	)

	result := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(result)

	//Output:
	//[1 2 3 4 5 6 7 8 9 10]
}

func ExampleOrderByOpE() {
	data := []string{"1", "10", "2", "9", "3", "8", "4", "7", "5", "6"}

	optic := Ordered(
		TraverseSlice[string](),
		OrderByOpE(
			func(ctx context.Context, left, right string) (bool, error) {

				leftInt, err := strconv.ParseInt(left, 10, 32)
				if err != nil {
					return false, err
				}

				rightInt, err := strconv.ParseInt(right, 10, 32)
				if err != nil {
					return false, err
				}

				return leftInt < rightInt, nil
			},
		),
	)

	result, err := Get(SliceOf(optic, len(data)), data)
	fmt.Println(result, err)

	//Output:
	//[1 2 3 4 5 6 7 8 9 10] <nil>
}

func ExampleOrderByOpI() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	optic := OrderedI(
		TraverseSlice[int](),
		OrderByOpI(
			func(indexLeft int, left int, indexRight int, right int) bool {
				return -indexLeft < -indexRight
			},
		),
	)

	result := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(result)

	//Output:
	//[10 9 8 7 6 5 4 3 2 1]
}

func ExampleOrderByOpIE() {
	data := []string{"1", "10", "2", "9", "3", "8", "4", "7", "5", "6"}

	optic := OrderedI(
		TraverseSlice[string](),
		OrderByOpIE(
			func(ctx context.Context, indexLeft int, left string, indexRight int, right string) (bool, error) {

				leftInt, err := strconv.ParseInt(left, 10, 32)
				if err != nil {
					return false, err
				}

				rightInt, err := strconv.ParseInt(right, 10, 32)
				if err != nil {
					return false, err
				}

				return leftInt*int64(indexLeft) < rightInt*int64(indexRight), nil
			},
		),
	)

	result, err := Get(SliceOf(optic, len(data)), data)
	fmt.Println(result, err)

	//Output:
	//[1 2 10 3 4 9 5 8 7 6] <nil>
}

func ExampleDesc() {
	data := []int{1, 10, 2, 9, 3, 8, 4, 7, 5, 6}

	optic := Ordered(
		TraverseSlice[int](),
		Desc(
			OrderBy(
				Identity[int](),
			),
		),
	)

	result := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(result)

	//Output:
	//[10 9 8 7 6 5 4 3 2 1]
}

func ExampleOrderBy() {

	type Person struct {
		name string
		age  int
	}

	personAge := FieldLens(func(source *Person) *int { return &source.age })

	data := []Person{
		{
			name: "Alice",
			age:  46,
		},
		{
			name: "Bob",
			age:  45,
		},
		{
			name: "Carol",
			age:  44,
		},
	}

	optic := Ordered(
		TraverseSlice[Person](),
		OrderBy(
			personAge,
		),
	)

	result := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(result)

	//Output:
	//[{Carol 44} {Bob 45} {Alice 46}]
}

func ExampleOrderByI() {

	data := map[int]string{
		4: "alpha",
		3: "beta",
		2: "gamma",
		1: "delta",
	}

	optic := OrderedI(
		TraverseMap[int, string](),
		OrderByI(
			ValueIIndex[int, string](),
		),
	)

	result := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(result)

	//Output:
	//[delta gamma beta alpha]
}

func ExampleOrderByN() {

	type Person struct {
		name string
		age  int
	}

	personName := FieldLens(func(source *Person) *string { return &source.name })
	personAge := FieldLens(func(source *Person) *int { return &source.age })

	data := []Person{
		{
			name: "Bob",
			age:  46,
		},
		{
			name: "Alice",
			age:  46,
		},
		{
			name: "Carol",
			age:  44,
		},
	}

	optic := Ordered(
		TraverseSlice[Person](),
		OrderByN(
			OrderBy(personAge),
			OrderBy(personName),
		),
	)

	result := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(result)

	//Output:
	//[{Carol 44} {Alice 46} {Bob 46}]
}

func ExampleOrderBy2() {

	type Person struct {
		name string
		age  int
	}

	personName := FieldLens(func(source *Person) *string { return &source.name })
	personAge := FieldLens(func(source *Person) *int { return &source.age })

	data := []Person{
		{
			name: "Bob",
			age:  46,
		},
		{
			name: "Alice",
			age:  46,
		},
		{
			name: "Carol",
			age:  44,
		},
	}

	optic := Ordered(
		TraverseSlice[Person](),
		OrderBy2(
			OrderBy(personAge),
			OrderBy(personName),
		),
	)

	result := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(result)

	//Output:
	//[{Carol 44} {Alice 46} {Bob 46}]
}

func ExampleDescI() {

	data := map[int]string{
		4: "alpha",
		3: "beta",
		2: "gamma",
		1: "delta",
	}

	optic := OrderedI(
		TraverseMap[int, string](),
		DescI(
			OrderByI(
				ValueIIndex[int, string](),
			),
		),
	)

	result := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(result)

	//Output:
	//[alpha beta gamma delta]
}
