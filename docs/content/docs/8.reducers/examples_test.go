package reducers

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/spearson78/go-optic"
)

func TestReducersSum(t *testing.T) {
	//BEGIN reducers_sum
	data := []int{1, 2, 3, 4}

	result, ok := MustGetFirst(
		Reduce(
			TraverseSlice[int](),
			Sum[int](),
		),
		data,
	)
	fmt.Println(result, ok)
	//END reducers_sum

	if !reflect.DeepEqual([]any{result, ok}, []any{
		//BEGIN reducers_sum_result
		10, true,
		//END reducers_sum_result
	}) {
		t.Fatal(result, ok)
	}
}

func TestReducersAddT2(t *testing.T) {

	//BEGIN reducers_playground_addt2
	data := []int{1, 2, 3, 4}

	//BEGIN reducers_addt2
	sum := AsReducer(
		0,            //Initial value
		AddT2[int](), //Reduction operation
	)
	//END reducers_addt2

	result, ok := MustGetFirst(
		Reduce(
			TraverseSlice[int](),
			sum,
		),
		data,
	)
	fmt.Println(result, ok)
	//END reducers_playground_addt2

	if !reflect.DeepEqual([]any{result, ok}, []any{
		//BEGIN reducers_addt2_result
		10, true,
		//END reducers_addt2_result
	}) {
		t.Fatal(result, ok)
	}
}
