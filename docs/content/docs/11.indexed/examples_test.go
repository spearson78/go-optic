package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
)

func TestIndexTraverseMap(t *testing.T) {
	//BEGIN traversemap

	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
		"delta": 4,
	}
	result := MustGet(
		SeqIOf(
			TraverseMap[string, int](),
		),
		data,
	)
	for index, value := range result {
		fmt.Println(index, value)
	}
	//END traversemap

	var sb strings.Builder
	sb.WriteString("\n")
	for index, value := range result {
		fmt.Fprintln(&sb, index, value)
	}

	resStr := sb.String()
	if resStr != `
alpha 1
beta 2
delta 4
gamma 3
` {
		t.Fatal(sb.String())
	}

}

func TestIndexLostIndex(t *testing.T) {
	//BEGIN lostindex

	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
		"delta": 4,
	}
	result, err := Get(
		SeqIOf(
			Compose(
				TraverseMap[string, int](),
				Mul(10),
			),
		),
		data,
	)

	if err != nil {
		fmt.Println(err)
	}

	for index, value := range result {
		fmt.Println(index, value)
	}
	//END lostindex

	var sb strings.Builder
	sb.WriteString("\n")
	for index, value := range result {
		fmt.Fprintln(&sb, index, value)
	}

	resStr := sb.String()
	if resStr != `
{} 10
{} 20
{} 40
{} 30
` {
		t.Fatal(sb.String())
	}
}

func TestIndexComposeLeft(t *testing.T) {

	//BEGIN play_composeleft
	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
		"delta": 4,
	}

	//BEGIN composeleft
	result, err := Get(
		SeqIOf(
			ComposeLeft(
				TraverseMap[string, int](),
				Mul(10),
			),
		),
		data,
	)

	if err != nil {
		fmt.Println(err)
	}

	for index, value := range result {
		fmt.Println(index, value)
	}
	//END composeleft
	//END play_composeleft

	var sb strings.Builder
	sb.WriteString("\n")
	for index, value := range result {
		fmt.Fprintln(&sb, index, value)
	}

	resStr := sb.String()
	if resStr != `
alpha 10
beta 20
delta 40
gamma 30
` {
		t.Fatal(sb.String())
	}
}

func TestIndexComposeLeftIxGet(t *testing.T) {

	//BEGIN play_ixget_composeleft
	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
		"delta": 4,
	}

	//BEGIN ixget_composeleft
	result, ok, err := GetFirst(
		Index(
			ComposeLeft(
				TraverseMap[string, int](),
				Mul(10),
			),
			"beta",
		),
		data,
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result, ok, err)

	//END ixget_composeleft
	//END play_ixget_composeleftixget

	var sb strings.Builder

	if !reflect.DeepEqual([]any{result, ok, err}, []any{
		//BEGIN res_ixget_composeleft
		20, true, nil,
		//END res_ixget_composeleft
	}) {
		t.Fatal(sb.String())
	}

}

func TestIndexComposeBoth(t *testing.T) {
	//BEGIN composeboth

	data := []map[string]int{
		map[string]int{
			"alpha": 1,
			"beta":  2,
		},
		map[string]int{
			"gamma": 3,
			"delta": 4,
		},
	}
	result, err := Get(
		SeqIOf(
			ComposeBoth(
				TraverseSlice[map[string]int](),
				TraverseMap[string, int](),
			),
		),
		data,
	)

	if err != nil {
		fmt.Println(err)
	}

	for index, value := range result {
		fmt.Println(index, value)
	}
	//END composeboth

	var sb strings.Builder
	sb.WriteString("\n")
	for index, value := range result {
		fmt.Fprintln(&sb, index, value)
	}

	resStr := sb.String()
	if resStr != `
{0 alpha} 1
{0 beta} 2
{1 delta} 4
{1 gamma} 3
` {
		t.Fatal(sb.String())
	}

}

func TestIndexComposeBothIxGet(t *testing.T) {
	//BEGIN play_ixget_composeboth

	data := []map[string]int{
		map[string]int{
			"alpha": 1,
			"beta":  2,
		},
		map[string]int{
			"gamma": 3,
			"delta": 4,
		},
	}

	//BEGIN ixget_composeboth
	result, ok, err := GetFirst(
		Index(
			ComposeBoth(
				TraverseSlice[map[string]int](),
				TraverseMap[string, int](),
			),
			lo.T2(1, "gamma"),
		),
		data,
	)

	fmt.Println(result, ok, err)
	//END ixget_composeboth
	//END play_ixget_composeboth

	if !reflect.DeepEqual([]any{result, ok, err}, []any{
		//BEGIN res_ixget_composeboth
		3, true, nil,
		//END res_ixget_composeboth
	}) {
		t.Fatal(result, ok, err)
	}

}

func TestComposeI(t *testing.T) {

	//BEGIN composei
	data := []map[string]int{
		map[string]int{
			"alpha": 1,
			"beta":  2,
		},
		map[string]int{
			"gamma": 3,
			"delta": 4,
		},
	}

	ixMap := IxMapIso[int, string, lo.Tuple2[int, string]](
		func(left int, right string) lo.Tuple2[int, string] {
			return lo.T2(left, right)
		},
		func(t1, t2 lo.Tuple2[int, string]) bool {
			return t1 == t2
		},
		func(mapped lo.Tuple2[int, string]) (int, bool, string, bool) {
			return mapped.A, true, mapped.B, true
		},
		ExprCustom("IxMapBoth"),
	)

	result, err := Get(
		SeqIOf(
			ComposeI(
				ixMap,
				TraverseSlice[map[string]int](),
				TraverseMap[string, int](),
			),
		),
		data,
	)

	if err != nil {
		fmt.Println(err)
	}

	for index, value := range result {
		fmt.Println(index, value)
	}
	//END composei

	var sb strings.Builder
	sb.WriteString("\n")
	for index, value := range result {
		fmt.Fprintln(&sb, index, value)
	}

	resStr := sb.String()
	if resStr != `
{0 alpha} 1
{0 beta} 2
{1 delta} 4
{1 gamma} 3
` {
		t.Fatal(sb.String())
	}

}

func TestReIndexed1(t *testing.T) {

	//BEGIN reindexed1
	data := []string{
		"alpha",
		"beta",
		"gamma",
		"delta",
	}

	optic := ReIndexed(
		TraverseSlice[string](),
		Compose(
			FormatInt[int](10),
			Op(func(focus string) string {
				return "KEY:" + focus
			}),
		),
		EqT2[string](), //IxMatch for the new index
	)

	res, err := Get(MapOf(optic, 10), data)
	fmt.Println(res, err)
	//END reindexed1
}

func TestReIndexed2(t *testing.T) {

	//BEGIN reindexed2
	data := []lo.Tuple2[string, int]{
		lo.T2("Max Mustermann", 42),
		lo.T2("Erika Mustermann", 37),
	}

	optic := ReIndexed(
		//Self index will set the index to the lo.Tuple2[string,int]
		SelfIndex(
			TraverseSlice[lo.Tuple2[string, int]](),
			EqT2[lo.Tuple2[string, int]](), //IxMatch for the new index
		),
		//Reindexed will then set the index to element A of the tuple.
		T2A[string, int](),
		EqT2[string](), //IxMatch for the new index
	)

	//We can then build a map the element A of the tuple as key and the fully tuple as the element.
	res, err := Get(MapOf(optic, 10), data)
	fmt.Println(res, err)
	//END reindexed2
}
