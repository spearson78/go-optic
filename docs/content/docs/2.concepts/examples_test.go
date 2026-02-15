package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	. "github.com/spearson78/go-optic"
)

func TestConceptsIntro(t *testing.T) {
	//BEGIN intro
	result := MustSet(
		Index(
			TraverseSlice[int](),
			3,
		),
		4,
		[]int{10, 20, 30, 40, 50},
	)
	fmt.Println(result)
	//END intro

	if !reflect.DeepEqual(result, []int{10, 20, 30, 4, 50}) {
		t.Fatal(result)
	}
}

func TestConceptsMustGetFirst(t *testing.T) {
	//BEGIN mustgetfirst
	result, ok := MustGetFirst(
		Index(
			TraverseSlice[int](),
			3,
		),
		[]int{10, 20, 30, 40, 50},
	)

	fmt.Println(result, ok)
	//END mustgetfirst

	if result != 40 || !ok {
		t.Fatal(result, ok)
	}
}

func TestConceptsMustSet(t *testing.T) {
	//BEGIN mustset
	result := MustSet(
		Index(
			TraverseMap[string, int](),
			"alpha",
		),
		1,
		map[string]int{
			"alpha": 10,
			"beta":  20,
		},
	)
	fmt.Println(result)
	//END mustset

	if !reflect.DeepEqual(
		result,
		//BEGIN mustset_result
		map[string]int{
			"alpha": 1,
			"beta":  20,
		},
		//END mustset_result
	) {
		t.Fatal(result)
	}
}

func TestConceptsMustModify(t *testing.T) {
	//BEGIN mustmodify
	result := MustModify(
		Filtered(
			TraverseSlice[int](),
			Lt(10),
		),
		Mul(2),
		[]int{1, 2, 30, 4, 5},
	)
	fmt.Println(result)
	//END mustmodify

	if !reflect.DeepEqual(
		result,
		//BEGIN mustmodify_result
		[]int{2, 4, 30, 8, 10},
		//END mustmodify_result
	) {

		t.Fatal(result)
	}

}

func TestConceptsActionsMustModify(t *testing.T) {
	//BEGIN actions_mustmodify
	result := MustModify(
		TraverseMap[string, int](),
		Mul(2),
		map[string]int{
			"alpha": 1,
			"beta":  2,
			"gamma": 3,
		},
	)
	fmt.Println(result)
	//END actions_mustmodify

	if !reflect.DeepEqual(
		result,
		//BEGIN actions_mustmodify_result
		map[string]int{
			"alpha": 2,
			"beta":  4,
			"gamma": 6,
		},
		//END actions_mustmodify_result
	) {

		t.Fatal(result)
	}
}

func TestConceptsActionsMustModifyI(t *testing.T) {
	//BEGIN actions_mustmodifyi
	result := MustModifyI(
		TraverseSlice[string](),
		OpI(func(index int, focus string) string {
			if index%2 == 0 {
				return strings.ToUpper(focus)
			} else {
				return focus
			}
		}),
		[]string{"alpha", "beta", "gamma", "delta"},
	)
	fmt.Println(result)
	//END actions_mustmodifyi

	if !reflect.DeepEqual(
		result,
		//BEGIN actions_mustmodifyi_result
		[]string{"ALPHA", "beta", "GAMMA", "delta"},
		//END actions_mustmodifyi_result
	) {

		t.Fatal(result)
	}
}

func TestConceptsActionsMustGetFirst(t *testing.T) {
	//BEGIN actions_mustgetfirst
	result, found := MustGetFirst(
		TraverseMap[string, int](),
		map[string]int{
			"alpha": 1,
			"beta":  2,
		},
	)
	fmt.Println(result, found)
	//END actions_mustgetfirst

	if !reflect.DeepEqual([]any{result, found}, []any{
		//BEGIN actions_mustgetfirst_result
		1, true,
		//END actions_mustgetfirst_result
	},
	) {

		t.Fatal(result)
	}
}

func TestConceptsActionsMustGetFirstI(t *testing.T) {
	//BEGIN actions_mustgetfirsti
	index, result, found := MustGetFirstI(
		TraverseMap[string, int](),
		map[string]int{
			"alpha": 1,
			"beta":  2,
		},
	)
	fmt.Println(index, result, found)
	//END actions_mustgetfirsti

	if !reflect.DeepEqual([]any{index, result, found}, []any{
		//BEGIN actions_mustgetfirsti_result
		"alpha", 1, true,
		//END actions_mustgetfirsti_result
	},
	) {

		t.Fatal(result)
	}
}

func TestConceptsActionsErrorAware(t *testing.T) {
	//BEGIN actions_erroraware
	result, err := Modify(
		ParseInt[int](10, 0),
		Mul(2),
		"1",
	)
	fmt.Println(result, err)
	//END actions_erroraware

	if !reflect.DeepEqual([]any{result, err}, []any{
		//BEGIN actions_erroraware_result
		"2", nil,
		//END actions_erroraware_result
	},
	) {
		t.Fatal(result)
	}
}

func TestConceptsActionsContextAware(t *testing.T) {
	//BEGIN actions_contextaware
	result, err := ModifyContext(
		context.Background(),
		ParseInt[int](10, 0),
		Mul(2),
		"1",
	)
	fmt.Println(result, err)
	//END actions_contextaware
	if !reflect.DeepEqual([]any{result, err}, []any{
		"2", nil,
	},
	) {
		t.Fatal(result)
	}
}

func TestConceptsOpticsIntroGetter(t *testing.T) {
	//BEGIN optics_intro_getter
	result := MustGet(Add(5), 10)
	fmt.Println(result)
	//END optics_intro_getter

	if !reflect.DeepEqual([]any{result}, []any{
		15,
	},
	) {
		t.Fatal(result)
	}
}

func TestConceptsOpticsIntroReverseGetter(t *testing.T) {
	//BEGIN optics_intro_reversegetter
	result := MustReverseGet(Add(5), 10)
	fmt.Println(result)
	//END optics_intro_reversegetter

	if !reflect.DeepEqual([]any{result}, []any{
		5,
	},
	) {
		t.Fatal(result)
	}
}

func TestConceptsOpticsGetters(t *testing.T) {
	//BEGIN optics_getter
	result := MustGet(
		Eq("alpha"),
		"beta",
	)
	fmt.Println(result)
	//END optics_getter

	if result != false {
		t.Fatal(result)
	}
}

func TestConceptsOpticsGettersPredicate(t *testing.T) {
	//BEGIN optics_getter_predicate
	result := MustSet(
		Filtered(
			TraverseSlice[int](), //Optic
			Gt(10),               //Predicate
		),
		10, //New value
		[]int{1, 2, 30, 4, 50},
	)
	fmt.Println(result)
	//END optics_getter_predicate

	if !reflect.DeepEqual([]any{result}, []any{
		//BEGIN optics_getter_predicate_result
		[]int{1, 2, 10, 4, 10},
		//END optics_getter_predicate_result
	},
	) {
		t.Fatal(result)
	}
}

func TestConceptsOpticsGettersOp(t *testing.T) {
	//BEGIN optics_getter_op
	result := MustGet(
		Op(
			func(source []string) int {
				return len(source)
			},
		),
		[]string{"alpha", "beta"},
	)
	fmt.Println(result)
	//END optics_getter_op

	if result != 2 {
		t.Fatal(result)
	}
}

func TestConceptsOpticsGettersOpExisting(t *testing.T) {
	//BEGIN optics_getter_op_existing
	result := MustModify(
		TraverseSlice[string](), //Optic
		Op(strings.ToUpper),     //Modify Operation
		[]string{"alpha", "beta"},
	)
	fmt.Println(result)
	//END optics_getter_op_existing

	if !reflect.DeepEqual([]any{result}, []any{
		//BEGIN optics_getter_op_existing_result
		[]string{"ALPHA", "BETA"},
		//END optics_getter_op_existing_result
	},
	) {
		t.Fatal(result)
	}
}

func TestConceptsOpticsLensesFieldLens(t *testing.T) {
	//BEGIN optics_lenses_fieldlens
	type ExampleStruct struct {
		name    string
		address string
	}

	nameField := FieldLens(
		func(source *ExampleStruct) *string {
			return &source.name
		},
	)

	data := ExampleStruct{
		name:    "Max Mustermann",
		address: "Musterstadt",
	}

	result := MustSet(nameField, "Erika Mustermann", data)
	fmt.Println(result)
	//END optics_lenses_fieldlens

	if !reflect.DeepEqual(result,
		//BEGIN optics_lenses_fieldlens_result
		ExampleStruct{
			name:    "Erika Mustermann",
			address: "Musterstadt",
		},
		//END optics_lenses_fieldlens_result
	) {
		t.Fatal(result)
	}
}

func TestConceptsOpticsIsoCelsius(t *testing.T) {
	//BEGIN optics_isos_celsius
	celsiusToFahrenheit := Compose(
		Mul(1.8),
		Add(32.0),
	)

	fahrenHeit := MustGet(celsiusToFahrenheit, 32.0)
	fmt.Println(fahrenHeit)

	celsius := MustReverseGet(celsiusToFahrenheit, 89.6)
	fmt.Println(celsius)
	//END optics_isos_celsius

	if !reflect.DeepEqual([]any{fahrenHeit, celsius}, []any{
		89.6,
		31.999999999999996,
	},
	) {
		t.Fatal(fahrenHeit, celsius)
	}
}

func TestConceptsOpticsIsoFrom(t *testing.T) {
	//BEGIN optics_isos_from
	celsiusToFahrenheit := Compose(Mul(1.8), Add(32.0))

	celsius := MustGet(
		AsReverseGet(
			celsiusToFahrenheit,
		),
		89.6, //Fahrenheit
	)
	fmt.Println(celsius)
	//END optics_isos_from

	if !reflect.DeepEqual([]any{celsius}, []any{
		31.999999999999996,
	},
	) {
		t.Fatal(celsius)
	}
}

func TestConceptsOpticsIterationsReduce(t *testing.T) {
	//BEGIN optics_iterations_reduce
	sum, found, err := GetFirst(
		MapReduce(
			TraverseSlice[string](), //Optic
			ParseInt[int](10, 32),   //Mapper
			Sum[int](),              //Reducer
		),
		[]string{"1", "2", "3", "4"},
	)
	fmt.Println(sum, found, err)
	//END optics_iterations_reduce

	if !reflect.DeepEqual([]any{sum, found, err}, []any{
		10, true, nil,
	},
	) {
		t.Fatal(sum, found, err)
	}
}

func TestConceptsOpticsTraversalSeqOf(t *testing.T) {
	//BEGIN optics_traversal_seqof
	for v := range MustGet(
		SeqOf(
			TraverseSlice[int](),
		),
		[]int{1, 2, 4, 5},
	) {
		fmt.Println(v)
	}
	//END optics_traversal_seqof
}

func TestConceptsOpticsTraversalModify(t *testing.T) {
	//BEGIN optics_traversal_modify
	result := MustModify(
		TraverseSlice[int](),
		Mul(2),
		[]int{1, 2, 3, 4},
	)
	fmt.Println(result)
	//END optics_traversal_modify

	if !reflect.DeepEqual([]any{result}, []any{
		//BEGIN optics_traversal_modify_result
		[]int{2, 4, 6, 8},
		//END optics_traversal_modify_result
	},
	) {
		t.Fatal(result)
	}
}

func TestConceptsOpticsTraversalIndex(t *testing.T) {
	//BEGIN optics_traversal_index
	result, found := MustGetFirst(
		Index(
			TraverseMap[string, int](), //Optic
			"beta",                     //Index to access
		),
		map[string]int{
			"alpha": 1,
			"beta":  2,
			"gamma": 3,
			"delta": 4,
		},
	)
	fmt.Println(result, found)
	//END optics_traversal_index

	if !reflect.DeepEqual([]any{result, found}, []any{
		//BEGIN optics_traversal_index_result
		2, true,
		//END optics_traversal_index_result
	},
	) {
		t.Fatal(result)
	}
}

func TestConceptsOpticsPrismsGet(t *testing.T) {
	//BEGIN optics_prism_get
	result := MustGet(
		SliceOf(
			Compose(
				TraverseSlice[any](),
				DownCast[any, int](),
			),
			3, //Initial size of slice
		),
		[]any{1, "two", 3},
	)
	fmt.Println(result)
	//END optics_prism_get

	if !reflect.DeepEqual([]any{result}, []any{
		//BEGIN optics_prism_get_result
		[]int{1, 3},
		//END optics_prism_get_result
	},
	) {
		t.Fatal(result)
	}
}

func TestConceptsOpticsPrismModify(t *testing.T) {
	//BEGIN optics_prism_modify
	result := MustModify(
		Compose(
			TraverseSlice[any](),
			DownCast[any, int](),
		),
		Mul(2),
		[]any{1, "two", 3},
	)
	fmt.Println(result)
	//END optics_prism_modify

	if !reflect.DeepEqual([]any{result}, []any{
		//BEGIN optics_prism_modify_result
		[]any{2, "two", 6},
		//END optics_prism_modify_result
	},
	) {
		t.Fatal(result)
	}
}

//BEGIN optics_compose_types

type Comment struct {
	title   string
	content string
}
type BlogPost struct {
	content  string
	comments []Comment
}

//END optics_compose_types

func TestConceptsOpticsComposeFieldLens(t *testing.T) {

	//BEGIN optics_compose_fieldlens
	blogComments := FieldLens(func(source *BlogPost) *[]Comment {
		return &source.comments
	})

	commentTitle := FieldLens(func(source *Comment) *string {
		return &source.title
	})
	//END optics_compose_fieldlens

	//BEGIN optics_compose_fieldlens_composed
	blogCommentTitles := Compose3(
		blogComments,
		TraverseSlice[Comment](),
		commentTitle,
	)
	//END optics_compose_fieldlens_composed

	//BEGIN optics_compose_get
	for commentTitle := range MustGet(
		SeqOf(
			blogCommentTitles,
		),
		BlogPost{ /*...*/ },
	) {
		fmt.Println(commentTitle)
	}
	//END optics_compose_get

	//BEGIN optics_compose_modify
	updatedBlogPost := MustModify(
		blogCommentTitles,
		Op(strings.ToUpper),
		BlogPost{ /*...*/ },
	)
	//END optics_compose_modify

	fmt.Println(updatedBlogPost)
}

func TestConceptsOpticsComposeFieldLensGet(t *testing.T) {

	//BEGIN optics_compose_fieldlens_playground
	for commentTitle := range MustGet(
		SeqOf(
			Compose3(
				FieldLens(func(source *BlogPost) *[]Comment { return &source.comments }),
				TraverseSlice[Comment](),
				FieldLens(func(source *Comment) *string { return &source.title }),
			),
		),
		BlogPost{
			comments: []Comment{
				Comment{title: "First Comment"},
				Comment{title: "Second Comment"},
			},
		},
	) {
		fmt.Println(commentTitle)
	}
	//END optics_compose_fieldlens_playground
}

func TestConceptsOpticsComposeFieldLensModify(t *testing.T) {

	//BEGIN optics_compose_fieldlens_modify_playground
	updatedBlogPost := MustModify(
		Compose3(
			FieldLens(func(source *BlogPost) *[]Comment { return &source.comments }),
			TraverseSlice[Comment](),
			FieldLens(func(source *Comment) *string { return &source.title }),
		),
		Op(strings.ToUpper),
		BlogPost{
			comments: []Comment{
				Comment{title: "First Comment"},
				Comment{title: "Second Comment"},
			},
		},
	)

	fmt.Println(updatedBlogPost)
	//END optics_compose_fieldlens_modify_playground
}

func TestConceptsOpticsComposeRight(t *testing.T) {

	//BEGIN optics_compose_right
	index, result, ok := MustGetFirstI(
		Compose3(
			FieldLens(func(source *BlogPost) *[]Comment { return &source.comments }),
			TraverseSlice[Comment](),
			FieldLens(func(source *Comment) *string { return &source.title }),
		),
		BlogPost{
			comments: []Comment{
				Comment{title: "First Comment"},
				Comment{title: "Second Comment"},
			},
		},
	)
	fmt.Println(index, result, ok)
	//END optics_compose_right

	if !reflect.DeepEqual([]any{index, result, ok}, []any{
		//BEGIN optics_compose_right_result
		Void{}, "First Comment", true,
		//END optics_compose_right_result
	},
	) {
		t.Fatal(index, result, ok)
	}
}

func TestConceptsOpticsComposeLeft(t *testing.T) {

	//BEGIN optics_compose_left
	index, result, ok := MustGetFirstI(
		Compose(
			FieldLens(func(source *BlogPost) *[]Comment { return &source.comments }),
			ComposeLeft(
				TraverseSlice[Comment](),
				FieldLens(func(source *Comment) *string { return &source.title }),
			),
		),
		BlogPost{
			comments: []Comment{
				Comment{title: "First Comment"},
				Comment{title: "Second Comment"},
			},
		},
	)
	fmt.Println(index, result, ok)
	//END optics_compose_left

	if !reflect.DeepEqual([]any{index, result, ok}, []any{
		//BEGIN optics_compose_left_result
		0, "First Comment", true,
		//END optics_compose_left_result
	},
	) {
		t.Fatal(index, result, ok)
	}
}

func TestConceptsOpticsCombinators(t *testing.T) {
	//BEGIN optics_combinators
	result := MustModify(
		Filtered( //Combinator
			TraverseSlice[int](),
			AndOp( //Combinator
				Gt(10),
				Lt(40),
			),
		),
		Compose( //Combinator
			Add(10),
			Mul(2),
		),
		[]int{10, 20, 30, 40, 50},
	)
	fmt.Println(result)
	//END optics_combinators

	if !reflect.DeepEqual([]any{result}, []any{
		//BEGIN optics_combinators_result
		[]int{10, 60, 80, 40, 50},
		//END optics_combinators_result
	},
	) {
		t.Fatal(result)
	}
}

func TestConceptsIdentityRule(t *testing.T) {
	//BEGIN identity_rule
	result := MustModify(
		Filtered(
			TraverseSlice[int](),
			Lt(10),
		),
		Mul(2),
		[]int{1, 2, 30, 4, 5},
	)
	fmt.Println(result)
	//END identity_rule

	if !reflect.DeepEqual(
		result,
		//BEGIN identity_rule_result
		[]int{2, 4, 30, 8, 10},
		//END identity_rule_result
	) {

		t.Fatal(result)
	}

}

func TestConceptsIdentityRuleIdentity(t *testing.T) {
	//BEGIN identity_rule_identity
	result := MustModify(
		Filtered(
			TraverseSlice[int](),
			Lt(10),
		),
		Identity[int](),
		[]int{1, 2, 30, 4, 5},
	)
	fmt.Println(result)
	//END identity_rule_identity

	if !reflect.DeepEqual(
		result,
		//BEGIN identity_rule_identity_result
		[]int{1, 2, 30, 4, 5},
		//END identity_rule_identity_result
	) {

		t.Fatal(result)
	}
}

func TestConceptsIdentityRuleReorder(t *testing.T) {
	//BEGIN identity_rule_identity_reorder
	result := MustGet(
		SliceOf(
			Ordered(
				TraverseSlice[int](),
				OrderBy[int, int](Identity[int]()),
			),
			10, //initial size of slice
		),
		[]int{5, 4, 30, 2, 1},
	)
	fmt.Println(result)
	//END identity_rule_identity_reorder

	if !reflect.DeepEqual(
		result,
		//BEGIN identity_rule_identity_reorder_result
		[]int{1, 2, 4, 5, 30},
		//END identity_rule_identity_reorder_result
	) {

		t.Fatal(result)
	}
}

func TestConceptsIdentityRuleReorderModify(t *testing.T) {
	//BEGIN identity_rule_identity_reorder_modify
	result := MustModify(
		Ordered(
			TraverseSlice[int](),
			OrderBy[int, int](Identity[int]()),
		),
		Mul(2),
		[]int{5, 4, 30, 2, 1},
	)
	fmt.Println(result)
	//END identity_rule_identity_reorder_modify

	if !reflect.DeepEqual(
		result,
		//BEGIN identity_rule_identity_reorder_modify_result
		[]int{10, 8, 60, 4, 2},
		//END identity_rule_identity_reorder_modify_result
	) {

		t.Fatal(result)
	}
}

func TestConceptsIdentityRuleReorderModifyTaking(t *testing.T) {
	//BEGIN identity_rule_identity_reorder_modify_taking
	result := MustModify(
		Taking(
			Ordered(
				TraverseSlice[int](),
				OrderBy[int, int](Identity[int]()),
			),
			3, //Take 3
		),
		Mul(2),
		[]int{5, 4, 30, 2, 1},
	)
	fmt.Println(result)
	//END identity_rule_identity_reorder_modify_taking

	if !reflect.DeepEqual(
		result,
		//BEGIN identity_rule_identity_reorder_modify_taking_result
		[]int{5, 8, 30, 4, 2},
		//END identity_rule_identity_reorder_modify_taking_result
	) {

		t.Fatal(result)
	}
}

func TestConceptsIdentityRuleContext(t *testing.T) {

	//BEGIN identity_rule_context
	celsiusToFahrenheit := Compose(
		Mul(1.8),
		Add(32.0),
	)

	celsiusResult := MustModify(
		celsiusToFahrenheit,
		Add(1.0),
		32,
	)
	fmt.Println(celsiusResult)
	//END identity_rule_context

	if !reflect.DeepEqual(
		celsiusResult,
		//BEGIN identity_rule_context_result
		32.55555555555555,
		//END identity_rule_context_result
	) {

		t.Fatal(celsiusResult)
	}
}

//BEGIN identity_rule_virtual

type Data struct {
	celsiusValue float64
}

func CelsiusValue() Optic[Void, Data, Data, float64, float64, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens(func(source *Data) *float64 { return &source.celsiusValue })
}

func FahrenheitValue() Optic[Void, Data, Data, float64, float64, ReturnOne, ReadWrite, UniDir, Pure] {
	return Ret1(Rw(Ud(EPure(Compose(
		CelsiusValue(),
		Compose(
			Mul(1.8),
			Add(32.0),
		),
	)))))
}

//END identity_rule_virtual
