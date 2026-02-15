package construction

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/spearson78/go-optic"
)

func TestConstructSliceOf(t *testing.T) {

	//BEGIN construct_sliceof
	data := []int{10, 20, 30, 40, 50}

	optic := Filtered(
		TraverseSlice[int](),
		Gte(30),
	)

	result := MustGet(
		SliceOf(
			optic,
			5, //Initial slice capacity
		),
		data,
	)

	fmt.Println(result)
	//END construct_sliceof

}

func TestConstructMapOf(t *testing.T) {

	//BEGIN construct_mapof
	data := []string{"alpha", "beta", "gamma"}

	result, err := Get(
		MapOf(
			TraverseSlice[string](),
			5, //Initial slice capacity
		),
		data,
	)

	fmt.Println(result, err)
	//END construct_mapof

	if !reflect.DeepEqual([]any{result, err}, []any{
		//BEGIN construct_mapof_result
		map[int]string{
			0: "alpha",
			1: "beta",
			2: "gamma",
		},
		//END construct_mapof_result
		nil,
	}) {
		t.Fatal(result, err)
	}

}

func TestConstructSliceOfModify(t *testing.T) {

	//BEGIN construct_sliceof_modify
	data := [][]string{
		[]string{"alpha", "beta"},
		[]string{"gamma", "delta"},
	}

	optic := SliceOf(
		Compose(
			TraverseSlice[[]string](),
			TraverseSlice[string](),
		),
		4, //Initial slice capacity
	)

	var result [][]string = MustModify(optic,
		Op(func(s []string) []string {

			fmt.Println(s)
			// []string{"alpha","beta","gamma","delta"}

			s[2] = "EDIT"
			return s
		}),
		data,
	)

	fmt.Println(result)
	//END construct_sliceof_modify

	if !reflect.DeepEqual([]any{result}, []any{
		//BEGIN construct_sliceof_modify_result
		[][]string{
			{"alpha", "beta"},
			{"EDIT", "delta"},
		},
		//END construct_sliceof_modify_result
	}) {
		t.Fatal(result)
	}

}
