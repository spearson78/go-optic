package main

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func TestColOp1(t *testing.T) {

	//BEGIN collection_op1
	result := MustModify(
		SliceToCol[int](),
		ReversedCol[int, int](),
		[]int{1, 2, 3},
	)

	fmt.Println(result)
	//END collection_op1

	if !reflect.DeepEqual([]any{result}, []any{
		//BEGIN result_collection_op1
		[]int{3, 2, 1},
		//END result_collection_op1
	}) {
		t.Fatal(result)
	}
}

func TestMakeLensColOp(t *testing.T) {

	//BEGIN makelenscolop
	result := MustModify(
		O.BlogPost().Ratings(), //Ratings focuses a Collection[int,Rating] not a []Rating
		ReversedCol[int, Rating](),
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)
	fmt.Println(result)
	//END makelenscolop

}

func TestColOp2(t *testing.T) {

	//BEGIN collection_op2
	result := MustModify(
		SliceToCol[int](),
		FilteredCol[int](
			Ne(2),
		),
		[]int{1, 2, 3},
	)

	fmt.Println(result)
	//END collection_op2

	if !reflect.DeepEqual([]any{result}, []any{
		//BEGIN result_collection_op2
		[]int{1, 3},
		//END result_collection_op2
	}) {
		t.Fatal(result)
	}
}

func TestColOpReconstrain1(t *testing.T) {

	//BEGIN colreconstrain1
	data := []string{"3", "beta", "2"}

	//The Sort may generate an error in the collection during sorting.
	//Collection[int, string, Err]
	sort := OrderedCol[int](
		OrderBy(
			ParseInt[int](10, 0),
		),
	)
	//END colreconstrain1

	//BEGIN colreconstrain2
	res, err := Modify(
		ColFocusErr( //Collection[int,string,Err]
			SliceToCol[string](), //Collection[int,string,Pure]
		),
		sort,
		data,
	)
	fmt.Println(res, err)
	//END colreconstrain2

	if err == nil {
		t.Fatal(res, err)
	}

}

func TestColOpReconstrain2(t *testing.T) {
	//BEGIN colreconstrain3
	//Collection[int,string,Pure]
	data := ValCol("3", "2", "1")

	//Optic[Void, Collection[int, string, Err], Collection[int, string, Err], Collection[int, string, Err], Collection[int, string, Err], ReturnOne, ReadWrite, UniDir, Err]
	sort := OrderedCol[int](
		OrderBy(
			ParseInt[int](10, 0),
		),
	)

	res, err := Modify(
		Identity[Collection[int, string, Err]](),
		sort,
		//Convert to Collection[int,string,Err]
		ColErr(data),
	)

	fmt.Println(res, err)
	//END colreconstrain3

	//Output:
	//Col[1:3 2:3 0:2] <nil>
}
