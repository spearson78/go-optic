package optic_test

import (
	"fmt"

	. "github.com/spearson78/go-optic"
)

func ExampleAsModify() {

	data := "1"

	var result int
	result, err := Get(
		AsModify(
			ParseIntP[int](10, 32),
			Add(1),
		),
		data,
	)

	fmt.Println(result, err)

	//Output:
	//2 <nil>
}

func ExampleAsSet() {

	data := "1"

	var result int
	result, err := Get(
		AsSet(
			ParseIntP[int](10, 32),
			Const[string](2),
		),
		data,
	)

	fmt.Println(result, err)

	//Output:
	//2 <nil>
}
