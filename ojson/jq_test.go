package ojson_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/samber/mo"
	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/ojson"
	"github.com/spearson78/go-optic/otree"
	"golang.org/x/text/encoding/unicode"
)

func ExampleJqMergeT2() {

	res, err := Get(
		Compose3(
			ParseString[any](),
			T2Of(
				OptionOfFirst(Nth(0)),
				OptionOfFirst(Nth(1)),
			),
			JqMergeT2(),
		),
		`[{"k": {"a": 1, "b": 2}} , {"k": {"a": 0,"c": 3}}]`,
	)
	fmt.Println(res, err)

	//Output:
	//{true map[k:map[a:0 b:2 c:3]]} <nil>
}

func ExampleJqSlice() {

	sliceRes, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqSlice(2, 2),
		),
		`["a","b","c","d","e"]`,
	)
	fmt.Println(sliceRes, ok, err)

	stringRes, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqSlice(2, 2),
		),
		`"abcdefghi"`,
	)
	fmt.Println(stringRes, ok, err)

	//Output:
	//[c d] true <nil>
	//cd true <nil>
}

func ExampleJqObjectOf() {

	data := `["alpha",2]`

	optic := Compose(
		ParseString[any](),
		JqObjectOf(
			Concat(
				ReIndexed(Nth(0), Const[int]("author"), EqT2[string]()),
				ReIndexed(Nth(1), Const[int]("name"), EqT2[string]()),
			),
			2,
		),
	)

	res, err := Get(optic, data)
	fmt.Println(res, err)

	//Output:
	//map[author:alpha name:2] <nil>

}

func ExampleJqObjectsOf() {

	res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				JqObjectsOf(
					EErr(Concat(
						ReIndexed(
							ColFocusErr(ColOf(
								ReIndexed(
									KeyE("name"),
									UpCast[string, any](),
									EqT2[any](),
								),
							)),
							Const[Void]("name"),
							EqT2[string](),
						),
						ReIndexed(
							ColFocusErr(ColOf(
								Key("slice").TraverseE(),
							)),
							Const[Void]("value"),
							EqT2[string](),
						),
					)),
				),
			),
			2,
		),
		`{"name":"alpha","slice":["value 1", "value 2"]}`,
	)

	fmt.Println(res, err)

	//Output:
	//[map[name:alpha value:value 1] map[name:alpha value:value 2]] <nil>

}

func ExampleJqPick() {

	res, err := Modify(
		ParseString[any](),
		JqPick(
			otree.Path[any]("a"),
			otree.Path[any]("b", "c"),
			otree.Path[any]("x"),
		),
		`{"a": 1, "b": {"c": 2, "d": 3}, "e": 4}`,
	)

	fmt.Println(res, err)

	//Output:
	//{"a":1,"b":{"c":2},"x":null} <nil>
}

func ExampleJqOrder() {

	res, err := Get(
		SliceOf(
			Ordered(
				Compose(
					ParseString[any](),
					Traverse(),
				),
				JqOrder(),
			),
			10,
		),
		`[8,3,null,6]`,
	)
	fmt.Println(res, err)

	//Output:
	//[<nil> 3 6 8] <nil>
}

func ExampleJqOrderBy() {

	res, err := Get(
		SliceOf(
			Ordered(
				Compose(
					ParseString[any](),
					Traverse(),
				),
				JqOrderBy(Key("foo")),
			),
			10,
		),
		`[{"foo":4, "bar":10}, {"foo":3, "bar":10}, {"foo":2, "bar":1}]`,
	)
	fmt.Println(res, err)

	//Output:
	//[map[bar:1 foo:2] map[bar:10 foo:3] map[bar:10 foo:4]] <nil>
}

func TestJQTutorial(t *testing.T) {
	//These tests are derived from https://jqlang.github.io/jq/tutorial/

	if res, ok, err := GetFirst(Compose(ParseString[any](), Nth(0)), testJson); !ok || err != nil || res.(map[string]any)["sha"] != "588ff1874c8c394253c231733047a550efe78260" {
		t.Fatal(res, ok, err)
	}

	buildObj := JqObjectOf(
		Concat(
			ReIndexed(Key("commit").Key("author").Key("name"), Const[string]("author"), EqT2[string]()),
			ReIndexed(Key("commit").Key("committer").Key("name"), Const[string]("name"), EqT2[string]()),
		),
		2,
	)

	if res, ok, err := GetFirst(Compose3(ParseString[any](), Nth(0), buildObj), testJson); !ok || err != nil || !reflect.DeepEqual(res, map[string]any{
		"author": "dependabot[bot]",
		"name":   "GitHub",
	}) {
		t.Fatalf("%v : %T , %v, %v", res, res, ok, err)
	}

	if res, err := Get(SliceOf(Compose3(ParseString[any](), Traverse(), buildObj), 5), testJson); err != nil || !reflect.DeepEqual(res, []any{
		map[string]any{
			"author": "dependabot[bot]",
			"name":   "GitHub",
		},
		map[string]any{
			"author": "lectrical",
			"name":   "GitHub",
		},
		map[string]any{
			"author": "Emanuele Torre",
			"name":   "Emanuele Torre",
		},
		map[string]any{
			"author": "Emanuele Torre",
			"name":   "Emanuele Torre",
		},
		map[string]any{
			"author": "itchyny",
			"name":   "GitHub",
		},
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	buildObj2 := JqObjectOf(
		Concat3(
			ReIndexed(Key("commit").Key("author").Key("name"), Const[string]("author"), EqT2[string]()),
			ReIndexed(Key("commit").Key("committer").Key("name"), Const[string]("name"), EqT2[string]()),
			ReIndexed(
				Compose(SliceOf(Key("parents").Traverse().Key("html_url"), 10), UpCast[[]any, any]()),
				Const[Void]("parents"),
				EqT2[string](),
			),
		),
		2,
	)

	if res, err := Get(SliceOf(Compose3(ParseString[any](), Traverse(), buildObj2), 5), testJson); err != nil || !reflect.DeepEqual(res, []any{
		map[string]any{
			"author": "dependabot[bot]",
			"name":   "GitHub",
			"parents": []any{
				"https://github.com/jqlang/jq/commit/bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195",
			},
		},
		map[string]any{
			"author": "lectrical",
			"name":   "GitHub",
			"parents": []any{
				"https://github.com/jqlang/jq/commit/8bcdc9304ace5f2cc9bf662ab8998d75537e05f0",
			},
		},
		map[string]any{
			"author": "Emanuele Torre",
			"name":   "Emanuele Torre",
			"parents": []any{
				"https://github.com/jqlang/jq/commit/0b82b3841b05faefe0ac18379bdb361b7e4e3464",
			},
		},
		map[string]any{
			"author": "Emanuele Torre",
			"name":   "Emanuele Torre",
			"parents": []any{
				"https://github.com/jqlang/jq/commit/96e8d893c10ed2f7656ccb8cfa39a9a291663a7e",
			},
		},
		map[string]any{
			"author": "itchyny",
			"name":   "GitHub",
			"parents": []any{
				"https://github.com/jqlang/jq/commit/8619f8a8ac746887cf43abd8e2116abba253cdcc",
			},
		},
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

}

func ExampleJqContains() {

	res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqContains([]any{"baz", "bar"}),
		),
		`["foobar", "foobaz", "blarp"]`,
	)
	fmt.Println(res, ok, err)

	//Output:
	//true true <nil>
}

func ExampleJqInside() {

	res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqInside([]any{"foobar", "foobaz", "blarp"}),
		),
		`["baz", "bar"]`,
	)
	fmt.Println(res, ok, err)

	//Output
	//true true <nil>
}

func ExampleJqIndices() {
	res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				JqIndices(", "),
			),
			3,
		),
		`"a,b, cd, efg, hijk"`,
	)

	fmt.Println(res, err)

	//Output:
	//[3 7 12] <nil>
}

func ExampleJqLabel() {

	res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				JqLabel(
					Compose3(
						Traverse(),
						Float(),
						If(
							Eq(-1.0),
							JqBreak[float64, float64]("out"),
							Identity[float64](),
						),
					),
					"out",
				),
			),
			2,
		),
		`[1,2,-1,3]`,
	)

	fmt.Println(res, err)

	//Output:
	//[1 2] <nil>
}

func ExampleJqBreak() {

	res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				JqLabel(
					Compose3(
						Traverse(),
						Float(),
						If(
							Eq(-1.0),
							JqBreak[float64, float64]("out"),
							Identity[float64](),
						),
					),
					"out",
				),
			),
			2,
		),
		`[1,2,-1,3]`,
	)

	fmt.Println(res, err)

	//Output:
	//[1 2] <nil>
}

func TestJQManual(t *testing.T) {
	//These tests are derived from https://jqlang.github.io/jq/manual/

	//Identity
	if res, ok, err := GetFirst(ParseString[any](), `"Hello, World!"`); !ok || err != nil || !reflect.DeepEqual(res, "Hello, World!") {
		t.Fatalf("%v : %T , %v, %v", res, res, ok, err)
	}

	if res, ok, err := GetFirst(ParseString[any](), "0.12345678901234567890123456789"); !ok || err != nil || !reflect.DeepEqual(res, 0.12345678901234567890123456789) {
		t.Fatalf("%v : %T , %v, %v", res, res, ok, err)
	}

	if res, err := Get(SliceOf(
		Compose(
			ParseString[any](),
			Concat(
				Identity[any](),
				Compose(
					AsReverseGet(ParseString[any]()),
					IsoCastE[string, any](),
				),
			),
		), 2),
		"12345678909876543212345"); err != nil || !reflect.DeepEqual(res, []any{1.2345678909876543e+22, "1.2345678909876543e+22"}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, ok, err := GetFirst(
		Compose(
			ParseString[float64](),
			Lt(0.12345678901234567890123456788),
		),
		"0.12345678901234567890123456789",
	); !ok || err != nil || !reflect.DeepEqual(res, false) {
		t.Fatalf("%v : %T , %v, %v", res, res, ok, err)
	}

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				Traverse(),
				SliceOf(
					Concat(
						Identity[any](),
						Coalesce(
							Compose3(DownCast[any, int](), Eq(1), UpCast[bool, any]()),
							Compose3(DownCast[any, float64](), Eq(1.0), UpCast[bool, any]()),
						),
					),
					2,
				),
			),
			4,
		),
		"[1, 1.000, 1.0, 100e-2]",
	); err != nil || !reflect.DeepEqual(res, [][]any{[]any{float64(1.0), true}, []any{float64(1.0), true}, []any{float64(1.0), true}, []any{float64(1.0), true}}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				WithVar(
					Compose(
						Concat(
							Var[any, int]("$big"),
							Compose(Var[any, int]("$big"), Add(1)),
						),
						Gt(100),
					),
					"$big",
					IsoCastE[any, int](),
				),
			),
			2,
		),
		"100",
	); err != nil || !reflect.DeepEqual(res, []bool{false, true}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Object Identifier-Index
	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			Key("foo"),
		),
		`{"foo": 42, "bar": "less interesting data"}`,
	); !ok || err != nil || !reflect.DeepEqual(res, float64(42.0)) {
		t.Fatalf("%v : %T , %v, %v", res, res, ok, err)
	}

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			Key("foo"),
		),
		`{"notfoo": true, "alsonotfoo": false}`,
	); !ok || err != nil || res != nil {
		t.Fatalf("%v : %T , %v, %v", res, res, ok, err)
	}

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			KeyE("foo"),
		),
		`{"notfoo": true, "alsonotfoo": false}`,
	); !ok || err != nil || res != nil {
		t.Fatalf("%v : %T , %v, %v", res, res, ok, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				Key("foo"),
			),
			0,
		),
		`[1,2]`,
	); err != nil || len(res) != 0 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Array Index
	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			Nth(0),
		),
		`[{"name":"JSON", "good":true}, {"name":"XML", "good":false}]`,
	); !ok || err != nil || !reflect.DeepEqual(res, map[string]any{
		"name": "JSON",
		"good": true,
	}) {
		t.Fatalf("%v : %T , %v, %v", res, res, ok, err)
	}

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			Nth(2),
		),
		`[{"name":"JSON", "good":true}, {"name":"XML", "good":false}]`,
	); !ok || err != nil || res != nil {
		t.Fatalf("%v : %T , %v, %v", res, res, ok, err)
	}

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			Nth(-2),
		),
		`[1,2,3]`,
	); !ok || err != nil || res != 2.0 {
		t.Fatalf("%v : %T , %v, %v", res, res, ok, err)
	}

	//Array/String Slice
	//ojson Slice/SubSeq takes a start and length parameter this was implemented for the https://github.com/tidwall/sjson tests

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqSlice(2, 2),
		),
		`["a","b","c","d","e"]`,
	); !ok || err != nil || !reflect.DeepEqual(res, []any{"c", "d"}) {
		t.Fatalf("%v : %T , %v, %v", res, res, ok, err)
	}

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqSlice(2, 2),
		),
		`"abcdefghi"`,
	); !ok || err != nil || !reflect.DeepEqual(res, "cd") {
		t.Fatalf("%v : %T , %v, %v", res, res, ok, err)
	}

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqSlice(-2, 2), //By definition the resulting subslice must be 2 in length.
		),
		`["a","b","c","d","e"]`,
	); !ok || err != nil || !reflect.DeepEqual(res, []any{"d", "e"}) {
		t.Fatalf("%v : %T , %v, %v", res, res, ok, err)
	}

	//Array Object Value iterator
	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				TraverseE(),
			),
			2,
		),
		`[{"name":"JSON", "good":true}, {"name":"XML", "good":false}]`,
	); err != nil || !reflect.DeepEqual(res, []any{
		map[string]any{"name": "JSON", "good": true},
		map[string]any{"name": "XML", "good": false},
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				TraverseE(),
			),
			2,
		),
		`[]`,
	); err != nil || !reflect.DeepEqual(res, []any{}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				Key("foo").TraverseE(),
			),
			2,
		),
		`{"foo":[1,2,3]}`,
	); err != nil || !reflect.DeepEqual(res, []any{1.0, 2.0, 3.0}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				TraverseE(),
			),
			2,
		),
		`{"a": 1, "b": 1}`,
	); err != nil || !reflect.DeepEqual(res, []any{1.0, 1.0}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				Traverse(),
			),
			2,
		),
		`1.0`,
	); err != nil || !reflect.DeepEqual(res, []any{}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				TraverseE(),
			),
			2,
		),
		`1.0`,
	); !errors.Is(err, ErrTraverse) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Comma

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				Concat(
					KeyE("foo"),
					KeyE("bar"),
				),
			),
			2,
		),
		`{"foo": 42, "bar": "something else", "baz": true}`,
	); err != nil || !reflect.DeepEqual(res, []any{42.0, "something else"}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				Concat(
					ReIndexed(KeyE("user"), UpCast[string, any](), EqT2[any]()),
					KeyE("projects").Traverse(),
				),
			),
			2,
		),
		`{"user":"stedolan", "projects": ["jq", "wikiflow"]}`,
	); err != nil || !reflect.DeepEqual(res, []any{"stedolan", "jq", "wikiflow"}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				Concat(
					NthE(4),
					NthE(2),
				),
			),
			2,
		),
		`["a","b","c","d","e"]`,
	); err != nil || !reflect.DeepEqual(res, []any{"e", "c"}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	// Pipe

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				TraverseE().KeyE("name"),
			),
			2,
		),
		`[{"name":"JSON", "good":true}, {"name":"XML", "good":false}]`,
	); err != nil || !reflect.DeepEqual(res, []any{"JSON", "XML"}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Parenthesis

	//jq '(. + 2) * 5'
	if res, err := Get(
		Compose4(
			ParseString[any](),
			IntE(),
			Add[int64](2),
			Mul[int64](5),
		),
		`1`,
	); err != nil || !reflect.DeepEqual(res, int64(15)) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//jq '. + (2 * 5)'
	if res, err := Get(
		Compose3(
			ParseString[any](),
			IntE(),
			Add[int64](2*5),
		),
		`1`,
	); err != nil || !reflect.DeepEqual(res, int64(11)) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Array Construction
	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				TraverseE().Int(),
				Mul[int64](2),
			),
			3,
		),
		`[1, 2, 3]`,
	); err != nil || !reflect.DeepEqual(res, []int64{2, 4, 6}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Object Construction

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				JqObjectsOf(
					EErr(Concat(
						ColFocusErr(ReIndexed(
							ColOf(
								ReIndexed(
									KeyE("user"),
									UpCast[string, any](),
									EqT2[any](),
								),
							),
							Const[Void]("user"),
							EqT2[string](),
						)),
						ColFocusErr(ReIndexed(
							ColOf(
								Key("titles").TraverseE(),
							),
							Const[Void]("title"),
							EqT2[string](),
						)),
					)),
				),
			),
			2,
		),
		`{"user":"stedolan","titles":["JQ Primer", "More JQ"]}`,
	); err != nil || !reflect.DeepEqual(res, []any{
		map[string]any{
			"user":  "stedolan",
			"title": "JQ Primer",
		},
		map[string]any{
			"user":  "stedolan",
			"title": "More JQ",
		},
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		Compose(
			ParseString[any](),
			WithVar(
				JqObjectOf(
					ReIndexed(
						Key("titles").Optic,
						FirstOrDefault(Ignore(Var[string, string]("user"), Const[error](true)), "user"),
						EqT2[string](),
					),
					1,
				),
				"user",
				FirstOrError(KeyE("user").String(), errors.New("user key not found")),
			),
		),
		`{"user":"stedolan","titles":["JQ Primer", "More JQ"]}`,
	); err != nil || !reflect.DeepEqual(res, map[string]any{
		"stedolan": []any{
			"JQ Primer",
			"More JQ",
		},
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//The above optic is read write!
	if res, err := Set(
		Compose(
			ParseString[any](),
			WithVar(
				JqObjectOf(
					ReIndexed(
						Key("titles").Optic,
						FirstOrDefault(Ignore(Var[string, string]("user"), Const[error](true)), "user"),
						EqT2[string](),
					),
					1,
				),
				"user",
				FirstOrError(KeyE("user").String(), errors.New("user key not found")),
			),
		),
		any(map[string]any{
			"stedolan": []any{
				"edit1",
				"More JQ",
			},
		}),
		`{"user":"stedolan","titles":["JQ Primer", "More JQ"]}`,
	); err != nil || res != `{"titles":["edit1","More JQ"],"user":"stedolan"}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Object Construction
	assignmentJson := `{"posts": [{"title": "First post", "author": "anon"},
    {"title": "A well-written article", "author": "person1"}],
    "realnames": {"anon": "Anonymous Coward","person1": "Person McPherson"}}`

	realNames, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			KeyE("realnames").Object(),
		),
		assignmentJson,
	)
	if !ok || err != nil {
		log.Fatal(realNames, ok, err)
	}

	if res, err := Get(SliceOf(
		Compose(
			ParseString[any](),
			WithVar(
				Compose(
					Key("posts").Traverse(),
					JqObjectOf(
						Concat(
							Key("title"),
							ComposeLeft(
								Key("author"),
								Compose3(
									DownCast[any, string](),
									Lookup(TraverseMap[string, any](), realNames),
									Some[any](),
								),
							),
						),
						2,
					),
				),
				"$names",
				FirstOrDefault(Key("realnames"), any(map[string]any{})),
			),
		),
		10,
	),
		assignmentJson,
	); err != nil || !reflect.DeepEqual(res, []any{
		map[string]any{
			"author": "Anonymous Coward",
			"title":  "First post",
		},
		map[string]any{
			"author": "Person McPherson",
			"title":  "A well-written article",
		},
	}) {
		t.Fatal(res, err)
	}

	// ... (Recursive descent.)

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				Compose(
					otree.TopDown(Traverse()),
					Key("a"),
				),
			),
			1,
		),
		`[[{"a":1}]]`,
	); err != nil || !reflect.DeepEqual(res, []any{
		float64(1.0),
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Set(
		Compose(
			ParseString[any](),
			Compose(
				otree.TopDown(Traverse()),
				Key("a"),
			),
		),
		2,
		`[[{"a":1}]]`,
	); err != nil || res != `[[{"a":2}]]` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	// Addition

	if res, ok, err := GetFirst(
		Compose4(
			ParseString[any](),
			Key("a"),
			DownCast[any, float64](),
			Add(1.0),
		),
		`{"a": 7}`,
	); !ok || err != nil || res != 8.0 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Set(
		Compose4(
			ParseString[any](),
			Key("a"),
			DownCast[any, float64](),
			Add(1.0),
		),
		10,
		`{"a": 7}`,
	); err != nil || res != `{"a":9}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				Concat(
					Key("a").Traverse(),
					Key("b").Traverse(),
				),
			),
			4,
		),
		`{"a": [1,2], "b": [3,4]}`,
	); err != nil || !reflect.DeepEqual(res, []any{1.0, 2.0, 3.0, 4.0}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, ok, err := GetFirst(
		Compose4(
			ParseString[any](),
			FirstOrDefault(Key("a"), 0.0),
			DownCast[any, float64](),
			Add(1.0),
		),
		`{}`,
	); !ok || err != nil || res != 1.0 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		MapOfReduced(
			Compose(
				ParseString[any](),
				Concat4(
					ConstI[any]("a", 1.0, EqT2[string]()),
					ConstI[any]("b", 2.0, EqT2[string]()),
					ConstI[any]("c", 3.0, EqT2[string]()),
					ConstI[any]("a", 4.0, EqT2[string]()),
				),
			),
			FirstReducer[float64](),
			3,
		),
		`null`,
	); err != nil || !reflect.DeepEqual(res, map[string]float64{
		"a": 1.0,
		"b": 2.0,
		"c": 3.0,
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Subtraction

	if res, ok, err := GetFirst(
		Compose4(
			ParseString[any](),
			Key("a"),
			Float(),
			SubFrom(4.0),
		),
		`{"a": 3}`,
	); !ok || err != nil || res != 1.0 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Set(
		Compose4(
			ParseString[any](),
			Key("a"),
			Float(),
			SubFrom(4.0),
		),
		10,
		`{"a": 3}`,
	); err != nil || res != `{"a":-6}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Set(
		Compose4(
			ParseString[any](),
			Key("a"),
			DownCast[any, float64](),
			Add(1.0),
		),
		10,
		`{"a": 7}`,
	); err != nil || res != `{"a":9}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Filtered(
				Compose(
					ParseString[any](),
					Traverse(),
				),
				NotOp(In[any]("xml", "yaml")),
			),
			1,
		),
		`["xml", "yaml", "json"]`,
	); err != nil || !reflect.DeepEqual(res, []any{
		"json",
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Set(
		Filtered(
			Compose(
				ParseString[any](),
				Traverse(),
			),
			NotOp(In[any]("xml", "yaml")),
		),
		"modify",
		`["xml", "yaml", "json"]`,
	); err != nil || res != `["xml","yaml","modify"]` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Division

	if res, err := Get(
		Compose3(
			ParseString[float64](),
			DivQuotient(10.0),
			Mul(3.0),
		),
		`5`,
	); err != nil || res != 6.0 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Set(
		Compose3(
			ParseString[float64](),
			DivQuotient(10.0),
			Mul(3.0),
		),
		10.0,
		`5`,
	); err != nil || res != `3` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose4(
				ParseString[any](),
				Traverse(),
				Float(),
				DivQuotient(1.0),
			),
			3,
		),
		`[1,0,-1]`,
	); err != nil || !reflect.DeepEqual(res, []float64{1.0, math.Inf(1), -1.0}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Modify(
		Compose4(
			ParseString[any](),
			Traverse(),
			Float(),
			DivQuotient(1.0),
		),
		Add(2.0),
		`[1,0,-1]`,
	); err != nil || res != `[0.3333333333333333,0,1]` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[string](),
				SplitString(regexp.MustCompile(", ")),
			),
			10,
		),
		`"a, b,c,d, e"`,
	); err != nil || !reflect.DeepEqual(res, []string{"a", "b,c,d", "e"}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Set(
		Compose(
			ParseString[string](),
			SplitString(regexp.MustCompile(", ")),
		),
		"modify",
		`"a, b,c,d, e"`,
	); err != nil || res != `"modify, modify, modify"` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Multiplication

	if res, err := Get(
		Compose3(
			ParseString[any](),
			T2Of(
				OptionOfFirst(Nth(0)),
				OptionOfFirst(Nth(1)),
			),
			JqMergeT2(),
		),
		`[{"k": {"a": 1, "b": 2}} , {"k": {"a": 0,"c": 3}}]`,
	); err != nil || !reflect.DeepEqual(res, mo.Some(any(map[string]any{
		"k": map[string]any{
			"a": 0.0,
			"b": 2.0,
			"c": 3.0,
		},
	}))) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//abs

	if res, err := Modify(
		Compose3(
			ParseString[any](),
			Traverse(),
			Float(),
		),
		Abs[float64](),
		`[-10, -1.1, -1e-1]`,
	); err != nil || res != `[10,1.1,0.1]` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//length
	jqLength := FirstOrDefault(
		Coalesce4(
			Compose(DownCast[any, map[string]any](), Length(TraverseMap[string, any]())),
			Compose(DownCast[any, []any](), Length(TraverseSlice[any]())),
			Compose(DownCast[any, string](), Length(TraverseString())),
			Compose(DownCast[any, int](), Abs[int]()),
		),
		0,
	)

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				Traverse(),
				jqLength,
			),
			10,
		),
		`[[1,2], "string", {"a":2}, null, -5]`,
	); err != nil || !reflect.DeepEqual(res, []int{2, 6, 1, 0, 5}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//ut8bytelength

	utf8bytelength := Length(Compose(EncodeString(unicode.UTF8), TraverseSlice[byte]()))

	if res, err := Get(
		Compose(
			ParseString[string](),
			utf8bytelength,
		),
		`"\u03bc"`,
	); err != nil || res != 2 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//keys , keys_unsorted

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				WithIndex(Traverse()),
				ValueIIndex[any, any](),
			),
			3,
		),
		`{"abc": 1, "abcd": 2, "Foo": 3}`,
	); err != nil || !reflect.DeepEqual(res, []any{"Foo", "abc", "abcd"}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				WithIndex(Traverse()),
				ValueIIndex[any, any](),
			),
			3,
		),
		`[42,3,35]`,
	); err != nil || !reflect.DeepEqual(res, []any{0, 1, 2}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Has

	if res, err := Modify(
		Compose(
			ParseString[any](),
			Traverse(),
		),
		Compose(
			NotEmpty(
				Index(
					Compose(
						DownCast[any, map[string]any](),
						TraverseMap[string, any](),
					),
					"foo",
				),
			),
			UpCast[bool, any](),
		),
		`[{"foo": 42}, {}]`,
	); err != nil || res != `[true,false]` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Modify(
		Compose(
			ParseString[any](),
			Traverse(),
		),
		Compose(
			NotEmpty(
				Index(
					Compose(
						DownCast[any, []any](),
						TraverseSlice[any](),
					),
					2,
				),
			),
			UpCast[bool, any](),
		),
		`[[0,1], ["a","b","c"]]`,
	); err != nil || res != `[false,true]` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//In
	if res, err := Get(
		SliceOf(
			Compose4(
				ParseString[any](),
				Traverse(),
				Lookup(Traverse(), any(map[string]any{
					"foo": 42,
				})),
				Present[any](),
			),
			3,
		),
		`["foo", "bar"]`,
	); err != nil || !reflect.DeepEqual(res, []bool{true, false}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose5(
				ParseString[any](),
				Traverse(),
				DownCast[any, int](),
				Lookup(TraverseSlice[any](), []any{0, 1}),
				Present[any](),
			),
			3,
		),
		`[2, 0]`,
	); err != nil || !reflect.DeepEqual(res, []bool{false, true}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Map

	if res, err := Get(
		SliceOf(
			Compose4(
				ParseString[any](),
				Traverse(),
				Float(),
				Add(1.0),
			),
			3,
		),
		`[1,2,3]`,
	); err != nil || !reflect.DeepEqual(res, []float64{2, 3, 4}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Modify(
		Compose3(
			ParseString[any](),
			Traverse(),
			Float(),
		),
		Add(1.0),
		`{"a": 1, "b": 2, "c": 3}`,
	); err != nil || res != `{"a":2,"b":3,"c":4}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				Traverse(),
				Concat(
					Identity[any](),
					Identity[any](),
				),
			),
			4,
		),
		`[1,2]`,
	); err != nil || !reflect.DeepEqual(res, []any{1.0, 1.0, 2.0, 2.0}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Modify(
		Compose(
			ParseString[any](),
			ChildrenCol(),
		),
		FilteredCol[any](EErr(NotOp(In[any](nil, false)))),
		`{"a": null, "b": true, "c": false}`,
	); err != nil || res != `{"b":true}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Pick
	if res, err := Modify(
		ParseString[any](),
		JqPick(
			otree.Path[any]("a"),
			otree.Path[any]("b", "c"),
			otree.Path[any]("x"),
		),
		`{"a": 1, "b": {"c": 2, "d": 3}, "e": 4}`,
	); err != nil || res != `{"a":1,"b":{"c":2},"x":null}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Modify(
		ParseString[any](),
		JqPick(
			otree.Path[any](2),
			otree.Path[any](0),
			otree.Path[any](0),
		),
		`[1,2,3,4]`,
	); err != nil || res != `[1,null,3]` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Path

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				WithIndex(otree.TopDown(Traverse())),
				ValueIIndex[*otree.PathNode[any], any](),
			),
			3,
		),
		`{"a":[{"b":1}]}`,
	); err != nil || !MustGet(EqDeepT2[any](), lo.T2[any, any](res, []*otree.PathNode[any]{
		nil,
		otree.Path[any]("a"),
		otree.Path[any]("a", 0),
		otree.Path[any]("a", 0, "b"),
	})) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Del
	if res, err := Modify(
		ParseString[map[string]any](),
		FilteredMapI[string, any](NeI[any]("foo")),
		`{"foo": 42, "bar": 9001, "baz": 42}`,
	); err != nil || res != `{"bar":9001,"baz":42}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Modify(
		ParseString[[]any](),
		FilteredSliceI[any](OpOnIx[any](NotOp(In(1, 2)))),
		`["foo", "bar", "baz"]`,
	); err != nil || res != `["foo"]` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Getpath

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				Concat(
					Index(otree.TopDown(Traverse()), otree.Path[any]("a", "b")),
					Index(otree.TopDown(Traverse()), otree.Path[any]("a", "c")),
				),
			),
			3,
		),
		`{"a":{"b":0, "c":1}}`,
	); err != nil || !MustGet(EqDeepT2[any](), lo.T2[any, any](res, []any{
		0.0,
		1.0,
	})) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Setpath

	if res, err := Set(
		Compose3(
			ParseString[any](),
			Key("a"),
			Key("b"),
		),
		1.0,
		`null`,
	); err != nil || res != `{"a":{"b":1}}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Set(
		Compose3(
			ParseString[any](),
			Key("a"),
			Key("b"),
		),
		1.0,
		`{"a":{"b":0}}`,
	); err != nil || res != `{"a":{"b":1}}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Set(
		Compose3(
			ParseString[any](),
			Nth(0),
			Key("a"),
		),
		1.0,
		`null`,
	); err != nil || res != `[{"a":1}]` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Delpaths

	if res, err := Modify(
		Compose(
			ParseString[any](),
			otree.WithChildPath(
				EErr(ComposeLeft(
					otree.BottomUp(
						Traverse(),
					),
					ChildrenCol(),
				)),
			),
		),
		FilteredColI(
			EErr(OpOnIx[any](NotOp(EqDeep(otree.Path[any]("a", "b"))))),
			PredToIxMatch(otree.EqPathT2(EqT2[any]())),
		),
		`{"a":{"b":1},"x":{"y":2}}`,
	); err != nil || res != `{"a":{},"x":{"y":2}}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//to_entries , from entries, with_entries

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				WithIndex(Traverse()),
			),
			1,
		),
		`{"a": 1, "b": 2}`,
	); err != nil || !MustGet(EqDeepT2[any](), lo.T2[any, any](res, []ValueI[any, any]{
		ValI[any, any]("a", 1.0),
		ValI[any, any]("b", 2.0),
	})) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		MapOf(
			Compose(
				ParseString[any](),
				ReIndexed(
					ComposeLeft(
						SelfIndex(Traverse(), EqT2[any]()),
						Key("value"),
					),
					FirstOrDefault(
						Compose(
							Key("key"),
							DownCast[any, string](),
						),
						"",
					),
					EqT2[string](),
				),
			),
			1,
		),
		`[{"key":"a", "value":1}, {"key":"b", "value":2}]`,
	); err != nil || !reflect.DeepEqual(res, map[string]any{
		"a": 1.0,
		"b": 2.0,
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		MapOf(
			Compose(
				ParseString[any](),
				ReIndexed(
					Traverse(),
					FirstOrDefault(
						Compose(
							DownCast[any, string](),
							PrependString(StringCol("KEY_")),
						),
						"",
					),
					EqT2[string](),
				),
			),
			1,
		),
		`{"a": 1, "b": 2}`,
	); err != nil || !reflect.DeepEqual(res, map[string]any{
		"KEY_a": 1.0,
		"KEY_b": 2.0,
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//select

	if res, err := Modify(
		Compose(
			ParseString[any](),
			ChildrenCol(),
		),
		FilteredCol[any](EErr(Compose(
			Float(),
			Gte(2.0),
		))),
		`[1,5,3,0,7]`,
	); err != nil || res != `[5,3,7]` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				Filtered(
					Traverse(),
					Compose(
						Key("id"),
						Eq[any]("second"),
					),
				),
			),
			1,
		),
		`[{"id": "first", "val": 1}, {"id": "second", "val": 2}]`,
	); err != nil || !reflect.DeepEqual(res, []any{
		map[string]any{
			"id":  "second",
			"val": 2.0,
		},
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//error
	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				Error[any, any](errors.ErrUnsupported),
			),
			1,
		),
		`null`,
	); !errors.Is(err, errors.ErrUnsupported) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//paths

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				WithIndex(Dropping(otree.TopDown(Traverse()), 1)),
				ValueIIndex[*otree.PathNode[any], any](),
			),
			3,
		),
		`[1,[[],{"a":2}]]`,
	); err != nil || !MustGet(EqDeepT2[any](), lo.T2[any, any](res, []*otree.PathNode[any]{
		otree.Path[any](0),
		otree.Path[any](1),
		otree.Path[any](1, 0),
		otree.Path[any](1, 1),
		otree.Path[any](1, 1, "a"),
	})) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				WithIndex(
					ComposeLeft(
						Dropping(otree.TopDown(Traverse()), 1),
						DownCast[any, float64](),
					),
				),
				ValueIIndex[*otree.PathNode[any], float64](),
			),
			3,
		),
		`[1,[[],{"a":2}]]`,
	); err != nil || !MustGet(EqDeepT2[any](), lo.T2[any, any](res, []*otree.PathNode[any]{
		otree.Path[any](0),
		otree.Path[any](1, 1, "a"),
	})) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//add

	if res, ok, err := GetFirst(
		Reduce(
			Compose(
				ParseString[any](),
				Traverse().String(),
			),
			StringBuilderReducer(),
		),
		`["a","b","c"]`,
	); !ok || err != nil || res != "abc" {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, ok, err := GetFirst(
		Reduce(
			Compose(
				ParseString[any](),
				Traverse().Float(),
			),
			Sum[float64](),
		),
		`[1, 2, 3]`,
	); !ok || err != nil || res != 6 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, ok, err := GetFirst(
		Reduce(
			Compose(
				ParseString[any](),
				Traverse().Float(),
			),
			Sum[float64](),
		),
		`[]`,
	); ok {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//any

	if res, err := Get(
		Compose(
			ParseString[any](),
			Any(
				Traverse().Bool(),
				Identity[bool](),
			),
		),
		`[true, false]`,
	); err != nil || res != true {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		Compose(
			ParseString[any](),
			Any(
				Traverse().Bool(),
				Identity[bool](),
			),
		),
		`[false, false]`,
	); err != nil || res != false {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		Compose(
			ParseString[any](),
			Any(
				Traverse().Bool(),
				Identity[bool](),
			),
		),
		`[]`,
	); err != nil || res != false {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//sort

	if res, err := Get(
		SliceOf(
			Ordered(
				Compose(
					ParseString[any](),
					Traverse(),
				),
				JqOrder(),
			),
			10,
		),
		`[8,3,null,6]`,
	); err != nil || !reflect.DeepEqual(res, []any{nil, 3.0, 6.0, 8.0}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Ordered(
				Compose(
					ParseString[any](),
					Traverse(),
				),
				JqOrderBy(Key("foo")),
			),
			10,
		),
		`[{"foo":4, "bar":10}, {"foo":3, "bar":10}, {"foo":2, "bar":1}]`,
	); err != nil || !reflect.DeepEqual(res, []any{
		map[string]any{
			"foo": 2.0,
			"bar": 1.0,
		},
		map[string]any{
			"foo": 3.0,
			"bar": 10.0,
		},
		map[string]any{
			"foo": 4.0,
			"bar": 10.0,
		},
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Ordered(
				Compose(
					ParseString[any](),
					Traverse(),
				),
				OrderBy2(
					JqOrderBy(Key("foo")),
					JqOrderBy(Key("bar")),
				),
			),
			10,
		),
		`[{"foo":4, "bar":10}, {"foo":3, "bar":20}, {"foo":2, "bar":1}, {"foo":3, "bar":10}]`,
	); err != nil || !reflect.DeepEqual(res, []any{
		map[string]any{
			"foo": 2.0,
			"bar": 1.0,
		},
		map[string]any{
			"foo": 3.0,
			"bar": 10.0,
		},
		map[string]any{
			"foo": 3.0,
			"bar": 20.0,
		},
		map[string]any{
			"foo": 4.0,
			"bar": 10.0,
		},
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Group By

	if res, err := Get(
		SliceOf(
			Grouped(
				Compose(
					ParseString[any](),
					ReIndexed(
						SelfIndex(Traverse(), EqDeepT2[any]()),
						FirstOrDefault(Key("foo"), nil),
						EqDeepT2[any](),
					),
				),
				AppendSliceReducer[any](10),
			),
			2,
		),
		`[{"foo":1, "bar":10}, {"foo":3, "bar":100}, {"foo":1, "bar":1}]`,
	); err != nil || !reflect.DeepEqual(res, [][]any{
		{
			map[string]any{
				"foo": 1.0,
				"bar": 10.0,
			},
			map[string]any{
				"foo": 1.0,
				"bar": 1.0,
			},
		},
		{
			map[string]any{
				"foo": 3.0,
				"bar": 100.0,
			},
		},
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//min,max,min_by,max_by

	if res, ok, err := GetFirst(
		Reduce(
			Compose3(
				ParseString[any](),
				Traverse(),
				Float(),
			),
			MinReducer[float64](),
		),

		`[5,4,2,7]`,
	); !ok || err != nil || res != 2.0 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, ok, err := GetFirst(
		MaxOf(
			Compose(
				ParseString[any](),
				Traverse(),
			),
			FirstOrDefault(Key("foo").Float(), 0.0),
		),
		`[{"foo":1, "bar":14}, {"foo":2, "bar":3}]`,
	); !ok || err != nil || !reflect.DeepEqual(res, map[string]any{
		"foo": 2.0,
		"bar": 3.0,
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Unique

	if res, err := Get(
		SliceOf(
			Grouped(
				Ordered(
					Compose(
						ParseString[any](),
						SelfIndex(Traverse(), EqDeepT2[any]()),
					),
					JqOrder(),
				),
				FirstReducer[any](),
			),
			2,
		),
		`[1,2,5,3,5,3,1,3]`,
	); err != nil || !reflect.DeepEqual(res, []any{1.0, 2.0, 3.0, 5.0}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Grouped(
				Ordered(
					Compose(
						ParseString[any](),
						ReIndexed(
							SelfIndex(Traverse(), EqDeepT2[any]()),
							FirstOrDefault(Key("foo"), nil),
							EqDeepT2[any](),
						),
					),
					JqOrder(),
				),
				FirstReducer[any](),
			),
			2,
		),
		`[{"foo": 1, "bar": 2}, {"foo": 1, "bar": 3}, {"foo": 4, "bar": 5}]`,
	); err != nil || !reflect.DeepEqual(res, []any{

		map[string]any{
			"foo": 1.0,
			"bar": 2.0,
		},
		map[string]any{
			"foo": 4.0,
			"bar": 5.0,
		},
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Grouped(
				Ordered(
					Compose(
						ParseString[any](),
						ReIndexed(
							SelfIndex(Traverse(), EqDeepT2[any]()),
							Length(
								Compose(
									DownCast[any, string](),
									TraverseString(),
								),
							),
							EqT2[int](),
						),
					),
					OrderBy(
						Length(
							Compose(
								DownCast[any, string](),
								TraverseString(),
							),
						),
					),
				),
				FirstReducer[any](),
			),
			2,
		),
		`["chunky", "bacon", "kitten", "cicada", "asparagus"]`,
	); err != nil || !reflect.DeepEqual(res, []any{"bacon", "chunky", "asparagus"}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Reverse

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				Reversed(Traverse()),
			),
			4,
		),
		`[1,2,3,4]`,
	); err != nil || !reflect.DeepEqual(res, []any{4.0, 3.0, 2.0, 1.0}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//contains

	if res, err := Get(
		Compose(
			ParseString[any](),
			JqContains("bar"),
		),
		`"foobar"`,
	); err != nil || res != true {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqContains([]any{"baz", "bar"}),
		),
		`["foobar", "foobaz", "blarp"]`,
	); !ok || err != nil || res != true {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqContains([]any{"bazzzzz", "bar"}),
		),
		`["foobar", "foobaz", "blarp"]`,
	); !ok || err != nil || res != false {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//jq 'contains({foo: 12, bar: [{barp: 12}]})'
	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqContains(map[string]any{
				"foo": 12.0,
				"bar": []any{
					map[string]any{
						"barp": 12.0,
					},
				},
			}),
		),
		`{"foo": 12, "bar":[1,2,{"barp":12, "blip":13}]}`,
	); !ok || err != nil || res != true {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//jq 'contains({foo: 12, bar: [{barp: 15}]})'
	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqContains(map[string]any{
				"foo": 12.0,
				"bar": []any{
					map[string]any{
						"barp": 15.0,
					},
				},
			}),
		),
		`{"foo": 12, "bar":[1,2,{"barp":12, "blip":13}]}`,
	); !ok || err != nil || res != false {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//indices
	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				JqIndices(", "),
			),
			3,
		),
		`"a,b, cd, efg, hijk"`,
	); err != nil || !reflect.DeepEqual(res, []any{3, 7, 12}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				JqIndices(1.0),
			),
			3,
		),
		`[0,1,2,1,3,1,4]`,
	); err != nil || !reflect.DeepEqual(res, []any{1, 3, 5}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//jq 'indices([1,2])'
	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				JqIndices([]any{1.0, 2.0}),
			),
			3,
		),
		`[0,1,2,3,1,4,2,5,1,2,6,7]`,
	); err != nil || !reflect.DeepEqual(res, []any{1, 8}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//index rindex
	if res, _, err := GetFirst(
		Compose(
			ParseString[any](),
			Taking(JqIndices(", "), 1),
		),
		`"a,b, cd, efg, hijk"`,
	); err != nil || res != 3 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, _, err := GetFirst(
		Compose(
			ParseString[any](),
			Taking(JqIndices(1.0), 1),
		),
		`[0,1,2,1,3,1,4]`,
	); err != nil || res != 1 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, _, err := GetFirst(
		Compose(
			ParseString[any](),
			Taking(JqIndices([]any{1.0, 2.0}), 1),
		),
		`[0,1,2,3,1,4,2,5,1,2,6,7]`,
	); err != nil || res != 1 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, _, err := GetFirst(
		Compose(
			ParseString[any](),
			Last(JqIndices(", ")),
		),
		`"a,b, cd, efg, hijk"`,
	); err != nil || res != 12 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, _, err := GetFirst(
		Compose(
			ParseString[any](),
			Last(JqIndices(1.0)),
		),
		`[0,1,2,1,3,1,4]`,
	); err != nil || res != 5 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, _, err := GetFirst(
		Compose(
			ParseString[any](),
			Last(JqIndices([]any{1.0, 2.0})),
		),
		`[0,1,2,3,1,4,2,5,1,2,6,7]`,
	); err != nil || res != 8 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Inside
	if res, err := Get(
		Compose(
			ParseString[any](),
			JqInside("foobar"),
		),
		`"bar"`,
	); err != nil || res != true {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqInside([]any{"foobar", "foobaz", "blarp"}),
		),
		`["baz", "bar"]`,
	); !ok || err != nil || res != true {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqInside([]any{"foobar", "foobaz", "blarp"}),
		),
		`["bazzzzz", "bar"]`,
	); !ok || err != nil || res != false {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqInside(map[string]any{
				"foo": 12.0,
				"bar": []any{
					1.0,
					2.0,
					map[string]any{
						"barp": 12.0,
						"blip": 13.0,
					},
				},
			}),
		),
		`{"foo": 12, "bar": [{"barp": 12}]}`,
	); !ok || err != nil || res != true {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			JqInside(map[string]any{
				"foo": 12.0,
				"bar": []any{
					1.0,
					2.0,
					map[string]any{
						"barp": 12.0,
						"blip": 13.0,
					},
				},
			}),
		),
		`{"foo": 12, "bar": [{"barp": 15}]}`,
	); !ok || err != nil || res != false {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//startswith

	if res, err := Get(
		SliceOf(
			Compose4(
				ParseString[any](),
				Traverse(),
				String(),
				StringHasPrefix("foo"),
			),
			5,
		),
		`["fo", "foo", "barfoo", "foobar", "barfoob"]`,
	); err != nil || !reflect.DeepEqual(res, []bool{false, true, false, true, false}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose4(
				ParseString[any](),
				Traverse(),
				String(),
				StringHasSuffix("foo"),
			),
			5,
		),
		`["foobar", "barfoo"]`,
	); err != nil || !reflect.DeepEqual(res, []bool{false, true}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//While
	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[float64](),
				otree.TopDownFiltered(
					Mul(2.0),
					Lt(100.0),
				),
			),
			10,
		),
		`1`,
	); err != nil || !reflect.DeepEqual(res, []float64{
		1.0,
		2.0,
		4.0,
		8.0,
		16.0,
		32.0,
		64.0,
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Recurse(F) recurse,recurse(f;condition)

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				otree.TopDown(Key("foo").Traverse()),
			),
			1,
		),
		`{"foo":[{"foo": []}, {"foo":[{"foo":[]}]}]}`,
	); err != nil || !reflect.DeepEqual(res, []any{
		map[string]any{
			"foo": []any{
				map[string]any{
					"foo": []any{},
				},
				map[string]any{
					"foo": []any{
						map[string]any{
							"foo": []any{},
						},
					},
				},
			},
		},
		map[string]any{
			"foo": []any{},
		},
		map[string]any{
			"foo": []any{
				map[string]any{
					"foo": []any{},
				},
			},
		},
		map[string]any{
			"foo": []any{},
		},
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Set(
		Compose3(
			ParseString[any](),
			otree.TopDown(Key("foo").Traverse()),
			Key("modified"),
		),
		"value",
		`{"foo":[{"foo": []}, {"foo":[{"foo":[]}]}]}`,
	); err != nil || res != `{"foo":[{"foo":[],"modified":"value"},{"foo":[{"foo":[],"modified":"value"}],"modified":"value"}],"modified":"value"}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				otree.TopDown(Traverse()),
			),
			1,
		),
		`{"a":0,"b":[1]}`,
	); err != nil || !reflect.DeepEqual(res, []any{
		//[map[a:0 b:[1]] 0 [1] 1]
		map[string]any{
			"a": float64(0),
			"b": []any{float64(1)},
		},
		float64(0),
		[]any{float64(1)},
		float64(1),
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Set(
		Dropping(
			Compose(
				ParseString[any](),
				otree.TopDown(Traverse()),
			),
			1,
		),
		"value",
		`{"a":0,"b":[1]}`,
	); err != nil || res != `{"a":"value","b":"value"}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[float64](),
				otree.TopDownFiltered(MulOp(Identity[float64](), Identity[float64]()), Lt(float64(20))),
			),
			1,
		),
		`2`,
	); err != nil || !reflect.DeepEqual(res, []float64{
		float64(2),
		float64(4),
		float64(16),
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//catch

	if res, ok, err := GetFirst(
		Compose(
			ParseString[any](),
			Catch(
				KeyE("a").Optic,
				CombiGetter[Pure, string, error, error, any, any](
					func(ctx context.Context, err error) (string, any, error) {
						return "a", ". is not an object", nil
					},
					IxMatchComparable[string](),
					ExprCustom("catch"),
				),
			),
		),
		`true`,
	); !ok || err != nil || res != ". is not an object" {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				Traverse(),
				Ignore(
					KeyE("a").Optic,
					Const[error](true),
				),
			),
			2,
		),
		`[{}, true, {"a":1}]`,
	); err != nil || !reflect.DeepEqual(res, []any{
		nil,
		1.0,
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//label break
	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				JqLabel(
					Compose3(
						Traverse(),
						Float(),
						If(
							Eq(-1.0),
							JqBreak[float64, float64]("out"),
							Identity[float64](),
						),
					),
					"out",
				),
			),
			2,
		),
		`[1,2,3]`,
	); err != nil || !reflect.DeepEqual(res, []float64{
		1.0,
		2.0,
		3.0,
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose(
				ParseString[any](),
				JqLabel(
					Compose3(
						Traverse(),
						Float(),
						If(
							Eq(-1.0),
							JqBreak[float64, float64]("out"),
							Identity[float64](),
						),
					),
					"out",
				),
			),
			2,
		),
		`[1,2,-1,3]`,
	); err != nil || !reflect.DeepEqual(res, []float64{
		1.0,
		2.0,
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Regexps

	if res, err := Get(
		NotEmpty(
			Compose(
				ParseString[string](),
				MatchString(regexp.MustCompile("foo"), -1),
			),
		),
		`"foo"`,
	); err != nil || res != true {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				Traverse().String(),
				NotEmpty(MatchString(regexp.MustCompile("(?i)abc"), -1)),
			),
			2,
		),
		`["xabcd", "ABC"]`,
	); err != nil || !reflect.DeepEqual(res, []bool{true, true}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				String(),
				WithIndex(MatchString(regexp.MustCompile("(abc)+"), -1)),
			),
			2,
		),
		`"abc abc"`,
	); err != nil || !reflect.DeepEqual(res, []ValueI[MatchIndex, string]{ValI(MatchIndex{Offsets: []int{0, 3, 0, 3}, Captures: []string{"abc"}}, "abc"), ValI(MatchIndex{Offsets: []int{4, 7, 4, 7}, Captures: []string{"abc"}}, "abc")}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, ok, err := GetFirst(
		Compose3(
			ParseString[any](),
			String(),
			WithIndex(MatchString(regexp.MustCompile("foo"), -1)),
		),
		`"foo bar foo"`,
	); !ok || err != nil || !reflect.DeepEqual(res, ValI(MatchIndex{Offsets: []int{0, 3}, Captures: []string{}}, "foo")) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				String(),
				WithIndex(MatchString(regexp.MustCompile("(?i)foo"), -1)),
			),
			2,
		),
		`"foo bar FOO"`,
	); err != nil || !reflect.DeepEqual(res, []ValueI[MatchIndex, string]{ValI(MatchIndex{Offsets: []int{0, 3}, Captures: []string{}}, "foo"), ValI(MatchIndex{Offsets: []int{8, 11}, Captures: []string{}}, "FOO")}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				String(),
				WithIndex(MatchString(regexp.MustCompile("(?i)foo (bar)? foo"), -1)),
			),
			2,
		),
		`"foo bar foo foo  foo"`,
	); err != nil || !reflect.DeepEqual(res, []ValueI[MatchIndex, string]{ValI(MatchIndex{Offsets: []int{0, 11, 4, 7}, Captures: []string{"bar"}}, "foo bar foo"), ValI(MatchIndex{Offsets: []int{12, 20, -1, -1}, Captures: []string{""}}, "foo  foo")}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		Length(
			Compose3(
				ParseString[any](),
				String(),
				MatchString(regexp.MustCompile("."), -1),
			),
		),
		`"abc"`,
	); err != nil || res != 3 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//scan

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				String(),
				WithIndex(CaptureMapString(regexp.MustCompile("([a-z]+)-([0-9]+)"), -1)),
			),
			2,
		),
		`"xyzzy-14"`,
	); err != nil || !reflect.DeepEqual(res, []ValueI[lo.Tuple2[[]int, string], map[string]string]{
		ValI(
			lo.T2([]int{0, 8, 0, 5, 6, 8}, "xyzzy-14"),
			map[string]string{
				"0": "xyzzy",
				"1": "14",
			},
		),
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//scan is writeanle in optics.
	if res, err := Set(
		Index(
			Compose4(
				ParseString[any](),
				String(),
				CaptureMapString(regexp.MustCompile("([a-z]+)-([0-9]+)"), -1),
				TraverseMap[string, string](),
			),
			"1",
		),
		"15",
		`"xyzzy-14"`,
	); err != nil || res != "\"xyzzy-15\"" {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//scan
	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				String(),
				MatchString(regexp.MustCompile("c"), -1),
			),
			2,
		),
		`"abcdefabc"`,
	); err != nil || !reflect.DeepEqual(res, []string{"c", "c"}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				String(),
				CaptureString(regexp.MustCompile("(a+)(b+)"), -1),
			),
			2,
		),
		`"abaabbaaabbb"`,
	); err != nil || !reflect.DeepEqual(res, [][]string{
		{"a", "b"},
		{"aa", "bb"},
		{"aaa", "bbb"},
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Split

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				String(),
				SplitString(regexp.MustCompile(", *")),
			),
			2,
		),
		`"ab,cd, ef, gh"`,
	); err != nil || !reflect.DeepEqual(res, []string{
		"ab",
		"cd",
		"ef",
		"gh",
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//sub

	if res, err := ModifyI(
		Compose(
			ParseString[string](),
			MatchString(
				regexp.MustCompile("[^a-z]*([a-z]+)"),
				1,
			),
		),
		FirstOrDefault(
			Reduce(
				Concat(
					Const[ValueI[MatchIndex, string]]("Z"),
					FirstOrDefault(
						Index(
							Compose3(
								ValueIIndex[MatchIndex, string](),
								MatchIndexCaptures(),
								TraverseSlice[string](),
							),
							0,
						),
						"",
					),
				),
				StringBuilderReducer(),
			),
			"",
		),
		`"123abc456def"`,
	); err != nil || res != `"Zabc456def"` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Modify(
		WithComprehension(
			Compose(
				ParseString[string](),
				CaptureMapString(regexp.MustCompile(`(?P<a>.)`), 1),
			),
		),
		SliceOf(
			Concat(
				AsModify(Index(TraverseMap[string, string](), "a"), Op(strings.ToUpper)),
				AsModify(Index(TraverseMap[string, string](), "a"), Op(strings.ToLower)),
			),
			2,
		),
		`"aB"`,
	); err != nil || !reflect.DeepEqual(res, []string{`"AB"`, `"aB"`}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//gsub

	if res, err := ModifyI(
		Compose(
			ParseString[string](),
			MatchString(
				regexp.MustCompile("(.)[^a]*"),
				-1,
			),
		),
		FirstOrDefault(
			Reduce(
				Concat3(
					Const[ValueI[MatchIndex, string]]("+"),
					FirstOrDefault(
						Index(
							Compose3(
								ValueIIndex[MatchIndex, string](),
								MatchIndexCaptures(),
								TraverseSlice[string](),
							),
							0,
						),
						"",
					),
					Const[ValueI[MatchIndex, string]]("-"),
				),
				StringBuilderReducer(),
			),
			"",
		),
		`"Abcabc"`,
	); err != nil || res != `"+A-+a-"` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Modify(
		WithComprehension(
			Compose(
				ParseString[string](),
				MatchString(regexp.MustCompile(`p`), -1),
			),
		),
		SliceOf(
			Concat(
				Const[string]("a"),
				Const[string]("b"),
			),
			2,
		),
		`"p"`,
	); err != nil || !reflect.DeepEqual(res, []string{`"a"`, `"b"`}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Advanced Features

	//Variables

	if res, err := Get(
		Compose(
			ParseString[any](),
			WithVar(
				Compose(
					T2Of(
						FirstOrDefault(Key("foo").Float(), 0),
						FirstOrDefault(Var[any, float64]("x"), 0),
					),
					AddT2[float64](),
				),
				"x",
				FirstOrDefault(Key("bar").Float(), 0.0),
			),
		),
		`{"foo":10, "bar":200}`,
	); err != nil || res != 210 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		Compose(
			ParseString[any](),
			WithVar(
				SliceOf(
					Concat(
						Compose3(
							FloatE(),
							Mul(2.0),
							WithVar(
								Var[float64, any]("i"),
								"i",
								Identity[float64](),
							),
						),
						Var[any, any]("i"),
					),
					2,
				),
				"i",
				Identity[any](),
			),
		),
		`5`,
	); err != nil || !reflect.DeepEqual(res, []any{10.0, 5.0}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		Compose(
			ParseString[any](),
			WithVar(
				WithVar(
					WithVar(
						Compose4(
							T2Of(Var[any, float64]("a"), Var[any, float64]("b")),
							AddT2[float64](),
							T2Of(Identity[float64](), Var[float64, float64]("c")),
							AddT2[float64](),
						),
						"a",
						FirstOrDefault(Nth(0), any(0.0)),
					),
					"b",
					FirstOrDefault(Nth(1), any(0.0)),
				),
				"c",
				FirstOrDefault(Nth(2).Key("c"), any(0.0)),
			),
		),
		`[2, 3, {"c": 4, "d": 5}]`,
	); err != nil || res != 9 {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose3(
				ParseString[any](),
				Traverse(),

				WithVar(
					WithVar(
						JqObjectOf(
							Concat(
								ReIndexed(Var[any, any]("a"), Const[Void]("a"), EqT2[string]()),
								ReIndexed(Var[any, any]("b"), Const[Void]("b"), EqT2[string]()),
							),
							2,
						),
						"a",
						FirstOrDefault(Nth(0), any(nil)),
					),
					"b",
					FirstOrDefault(Nth(1), any(nil)),
				),
			),
			3,
		),
		`[[0], [0, 1], [2, 1, 0]]`,
	); err != nil || !reflect.DeepEqual(res, []any{
		map[string]any{
			"a": 0.0,
			"b": nil,
		},
		map[string]any{
			"a": 0.0,
			"b": 1.0,
		},
		map[string]any{
			"a": 2.0,
			"b": 1.0,
		},
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Plain assignment
	if res, err := Get(
		Compose3(
			ParseString[any](),
			AsSet(
				Key("a"),
				FirstOrDefault(Key("b"), nil),
			),
			AsReverseGet(ParseString[any]()),
		),
		`{"a": {"b": 10}, "b": 20}`,
	); err != nil || res != `{"a":20,"b":20}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		Compose3(
			ParseString[any](),
			AsSet(
				Key("a"),
				FirstOrDefault(Key("a").Key("b"), nil),
			),
			AsReverseGet(ParseString[any]()),
		),
		`{"a": {"b": 10}, "b": 20}`,
	); err != nil || res != `{"a":10,"b":20}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		Compose3(
			ParseString[any](),
			AsSet(
				Key("a"),
				FirstOrDefault(Key("a").Key("b"), nil),
			),
			AsReverseGet(ParseString[any]()),
		),
		`{"a": {"b": 10}, "b": 20}`,
	); err != nil || res != `{"a":10,"b":20}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Get(
		SliceOf(
			Compose4(
				TraverseCol[int, int](),
				UpCast[int, any](),
				JqObjectOf(
					Concat(
						ReIndexed(Identity[any](), Const[Void]("a"), EqT2[string]()),
						ReIndexed(Identity[any](), Const[Void]("b"), EqT2[string]()),
					),
					2,
				),
				AsReverseGet(ParseString[any]()),
			),
			3,
		),
		RangeCol(0, 2),
	); err != nil || !reflect.DeepEqual(res, []string{
		`{"a":0,"b":0}`,
		`{"a":1,"b":1}`,
		`{"a":2,"b":2}`,
	}) {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	//Complex assignment

	complexJson := `{
  "posts": [
    {
      "title": "First post",
      "author": "anon",
	  "comments": ["first comment","second comment"]
    },
    {
      "title": "A well-written article",
      "author": "stedolan",
	  "comments": ["third comment","fourth comment"]
    }
  ]
}`

	if res, err := Set(
		Compose(
			ParseString[any](),
			Key("posts").Nth(0).Key("title"),
		),
		"JQ manual",
		complexJson,
	); err != nil || res != `{"posts":[{"author":"anon","comments":["first comment","second comment"],"title":"JQ manual"},{"author":"stedolan","comments":["third comment","fourth comment"],"title":"A well-written article"}]}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Modify(
		Compose(
			ParseString[any](),
			Key("posts").Traverse().Key("comments").Array(),
		),
		AppendSlice[any](ValCol(any("this is great"))),
		complexJson,
	); err != nil || res != `{"posts":[{"author":"anon","comments":["first comment","second comment","this is great"],"title":"First post"},{"author":"stedolan","comments":["third comment","fourth comment","this is great"],"title":"A well-written article"}]}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

	if res, err := Modify(
		Compose3(
			ParseString[any](),
			Filtered(Key("posts").Traverse(), Key("author").Eq("stedolan")),
			Key("comments").Array(),
		),
		AppendSlice[any](ValCol(any("terrible"))),
		complexJson,
	); err != nil || res != `{"posts":[{"author":"anon","comments":["first comment","second comment"],"title":"First post"},{"author":"stedolan","comments":["third comment","fourth comment","terrible"],"title":"A well-written article"}]}` {
		t.Fatalf("%v : %T , %v", res, res, err)
	}

}

const testJson = `
[
  {
    "sha": "588ff1874c8c394253c231733047a550efe78260",
    "node_id": "C_kwDOAE3WVdoAKDU4OGZmMTg3NGM4YzM5NDI1M2MyMzE3MzMwNDdhNTUwZWZlNzgyNjA",
    "commit": {
      "author": {
        "name": "dependabot[bot]",
        "email": "49699333+dependabot[bot]@users.noreply.github.com",
        "date": "2024-12-29T12:54:10Z"
      },
      "committer": {
        "name": "GitHub",
        "email": "noreply@github.com",
        "date": "2024-12-29T12:54:10Z"
      },
      "message": "build(deps): bump jinja2 from 3.1.4 to 3.1.5 in /docs (#3226)\n\nBumps [jinja2](https://github.com/pallets/jinja) from 3.1.4 to 3.1.5.\r\n- [Release notes](https://github.com/pallets/jinja/releases)\r\n- [Changelog](https://github.com/pallets/jinja/blob/main/CHANGES.rst)\r\n- [Commits](https://github.com/pallets/jinja/compare/3.1.4...3.1.5)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: jinja2\r\n  dependency-type: direct:production\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
      "tree": {
        "sha": "e47dd35840a08fc1fa45e364ad536ac0941c1f72",
        "url": "https://api.github.com/repos/jqlang/jq/git/trees/e47dd35840a08fc1fa45e364ad536ac0941c1f72"
      },
      "url": "https://api.github.com/repos/jqlang/jq/git/commits/588ff1874c8c394253c231733047a550efe78260",
      "comment_count": 0,
      "verification": {
        "verified": true,
        "reason": "valid",
        "signature": "-----BEGIN PGP SIGNATURE-----\n\nwsFcBAABCAAQBQJncUZyCRC1aQ7uu5UhlAAAkqsQAGDitC5DSe3ySbRtmHzPf7iv\n5Khp9i9SyBWrRkD1Craega8o9YJnmkg2tMl1ewZ0ivRs95OIhJk7RDbnNWxOupmJ\nfvLXbi2KHqz4IJb3gzHmZ0k2mlPI/iGlSbvIuOs4ot7hfZP/4gXZgIPoMMkNJUu1\n+hgb8RtiKpQujamxNytUaahW/vTIO5OR+iinfB0fDk66LDiLtNYWlZeOVoGiD5rD\nogha/Cm7InhdiqUiac/L2fezPbi1ZDZfn2iKUHhVyeHb5RVafU20vzHw7gk32mFz\n9PgQcnliHrG1+l9dLSNp+7b/Lft57VDE8p8C8IJEsjwlXXChAOJn0x5Nrl+mo8Tm\ngt/Q18DaRYyzwHUdSMX9Kp4bTnPoydMsriT7aibFGgkI/VLlK5/zAiUOSw0ffSmV\nrF/DufMxWqkUy2U5egGbyo0ZCFu+198jwcQCQC4p+09Z9xB4/YkR1jgFmoJB6ZGS\nswj8SXN6dATpI0GFyCPka9cIZU780Q6RKWd3MDBJqAfth5lRVEkGRcZB0k3wL6B4\n+O870LGWn/eqyDxniWvYEz9MKTZhuT46FGS8d7WWsOdMBy+ssO78q9/XHy96z25w\najkM4XErrpnGtxH1CtYt4SkYBB9oTdY2pgJDnXQcKiUKScTgQPoqQ4OfH0R27xmX\n1anb5H4HAqQj3CrD9qe2\n=Zja6\n-----END PGP SIGNATURE-----\n",
        "payload": "tree e47dd35840a08fc1fa45e364ad536ac0941c1f72\nparent bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195\nauthor dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com> 1735476850 +0900\ncommitter GitHub <noreply@github.com> 1735476850 +0900\n\nbuild(deps): bump jinja2 from 3.1.4 to 3.1.5 in /docs (#3226)\n\nBumps [jinja2](https://github.com/pallets/jinja) from 3.1.4 to 3.1.5.\r\n- [Release notes](https://github.com/pallets/jinja/releases)\r\n- [Changelog](https://github.com/pallets/jinja/blob/main/CHANGES.rst)\r\n- [Commits](https://github.com/pallets/jinja/compare/3.1.4...3.1.5)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: jinja2\r\n  dependency-type: direct:production\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
        "verified_at": "2024-12-29T12:54:13Z"
      }
    },
    "url": "https://api.github.com/repos/jqlang/jq/commits/588ff1874c8c394253c231733047a550efe78260",
    "html_url": "https://github.com/jqlang/jq/commit/588ff1874c8c394253c231733047a550efe78260",
    "comments_url": "https://api.github.com/repos/jqlang/jq/commits/588ff1874c8c394253c231733047a550efe78260/comments",
    "author": {
      "login": "dependabot[bot]",
      "id": 49699333,
      "node_id": "MDM6Qm90NDk2OTkzMzM=",
      "avatar_url": "https://avatars.githubusercontent.com/in/29110?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/dependabot%5Bbot%5D",
      "html_url": "https://github.com/apps/dependabot",
      "followers_url": "https://api.github.com/users/dependabot%5Bbot%5D/followers",
      "following_url": "https://api.github.com/users/dependabot%5Bbot%5D/following{/other_user}",
      "gists_url": "https://api.github.com/users/dependabot%5Bbot%5D/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/dependabot%5Bbot%5D/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/dependabot%5Bbot%5D/subscriptions",
      "organizations_url": "https://api.github.com/users/dependabot%5Bbot%5D/orgs",
      "repos_url": "https://api.github.com/users/dependabot%5Bbot%5D/repos",
      "events_url": "https://api.github.com/users/dependabot%5Bbot%5D/events{/privacy}",
      "received_events_url": "https://api.github.com/users/dependabot%5Bbot%5D/received_events",
      "type": "Bot",
      "user_view_type": "public",
      "site_admin": false
    },
    "committer": {
      "login": "web-flow",
      "id": 19864447,
      "node_id": "MDQ6VXNlcjE5ODY0NDQ3",
      "avatar_url": "https://avatars.githubusercontent.com/u/19864447?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/web-flow",
      "html_url": "https://github.com/web-flow",
      "followers_url": "https://api.github.com/users/web-flow/followers",
      "following_url": "https://api.github.com/users/web-flow/following{/other_user}",
      "gists_url": "https://api.github.com/users/web-flow/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/web-flow/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/web-flow/subscriptions",
      "organizations_url": "https://api.github.com/users/web-flow/orgs",
      "repos_url": "https://api.github.com/users/web-flow/repos",
      "events_url": "https://api.github.com/users/web-flow/events{/privacy}",
      "received_events_url": "https://api.github.com/users/web-flow/received_events",
      "type": "User",
      "user_view_type": "public",
      "site_admin": false
    },
    "parents": [
      {
        "sha": "bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195",
        "url": "https://api.github.com/repos/jqlang/jq/commits/bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195",
        "html_url": "https://github.com/jqlang/jq/commit/bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195"
      }
    ]
  },
  {
    "sha": "bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195",
    "node_id": "C_kwDOAE3WVdoAKGJjYmYyYjQ2MTY4OTNjZjJlNmE0ZjhhOTJkYzNkYjNiMWVlYjExOTU",
    "commit": {
      "author": {
        "name": "lectrical",
        "email": "14344693+lectrical@users.noreply.github.com",
        "date": "2024-12-29T12:53:16Z"
      },
      "committer": {
        "name": "GitHub",
        "email": "noreply@github.com",
        "date": "2024-12-29T12:53:16Z"
      },
      "message": "Generate provenance attestations for release artifacts and docker image (#3225)\n\nAdding https://github.com/actions/attest-build-provenance to the ci builds so\r\nthat the release assets and docker image for the next release tag generate\r\nsigned build provenance attestations for workflow artifacts.",
      "tree": {
        "sha": "3b91252273d4a46a78a74974a09cb3dd62d73223",
        "url": "https://api.github.com/repos/jqlang/jq/git/trees/3b91252273d4a46a78a74974a09cb3dd62d73223"
      },
      "url": "https://api.github.com/repos/jqlang/jq/git/commits/bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195",
      "comment_count": 0,
      "verification": {
        "verified": true,
        "reason": "valid",
        "signature": "-----BEGIN PGP SIGNATURE-----\n\nwsFcBAABCAAQBQJncUY8CRC1aQ7uu5UhlAAAzxAQABJceWBtoI3mfLmZlUa5L5s5\ne98L6+EPWucJlfTHsIGKYLZbK1gITyNVzzgLaizo9+ht5cm+2I1H+nqcbhIYg5ge\nM0w838W6EzkF8EKTMOElI6YQuVfTZEgIu4nlF9e6484kkBm4ed1d1iLk5NpLa+6a\nOjNXFM8gbBNj/4/+SZD5NJkNp/An53JrKA/NQkzw9EINmnOfkQ3Og2NKxm74fUGu\nbnJS2oA61yUZFq+4Po8hudOutNkCMpm+MpHs3t/E/lFXMdgoiTjSmnCEswll3Rcj\ngKYxfrLCtOkC3P9IsN5OuYfQs3vIZBaeKxwMOLbrnarG2W4kgMwkF/fM1ToaUkrb\nEsepBoW4SLe118bwC1g8syWJX8l8Pd0YXJlOil1Og828p9oAGhUP/v3ycUL3sYaZ\ntOFEcyEJgZJgEcLpMZ0QXySWIyp5GGLb/keMdD6HxmnjLM9cliTNYs4YbzJoveo1\n3lykmA08chqldFADpJM93pWUtLOyV2WjAF2IdDBZQ/ZCni9Ge6a8v3JXR89jPHV/\nJhbzkoSUCR40IH7+ReAiRiz2MYLUUQMS2CHrwftCZGvIcYz7wrpSXDxMHcVrSdr4\n0bjbVoDeCVzJ9IhJR705ZglFh1+KvwaquO3fT8r/6wFITHAaw9u9EmMIJN7IfMnt\ni0a+sTL6HIqy8KIwK3Pk\n=ByY+\n-----END PGP SIGNATURE-----\n",
        "payload": "tree 3b91252273d4a46a78a74974a09cb3dd62d73223\nparent 8bcdc9304ace5f2cc9bf662ab8998d75537e05f0\nauthor lectrical <14344693+lectrical@users.noreply.github.com> 1735476796 +0000\ncommitter GitHub <noreply@github.com> 1735476796 +0900\n\nGenerate provenance attestations for release artifacts and docker image (#3225)\n\nAdding https://github.com/actions/attest-build-provenance to the ci builds so\r\nthat the release assets and docker image for the next release tag generate\r\nsigned build provenance attestations for workflow artifacts.",
        "verified_at": "2024-12-29T12:53:19Z"
      }
    },
    "url": "https://api.github.com/repos/jqlang/jq/commits/bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195",
    "html_url": "https://github.com/jqlang/jq/commit/bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195",
    "comments_url": "https://api.github.com/repos/jqlang/jq/commits/bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195/comments",
    "author": {
      "login": "lectrical",
      "id": 14344693,
      "node_id": "MDQ6VXNlcjE0MzQ0Njkz",
      "avatar_url": "https://avatars.githubusercontent.com/u/14344693?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/lectrical",
      "html_url": "https://github.com/lectrical",
      "followers_url": "https://api.github.com/users/lectrical/followers",
      "following_url": "https://api.github.com/users/lectrical/following{/other_user}",
      "gists_url": "https://api.github.com/users/lectrical/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/lectrical/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/lectrical/subscriptions",
      "organizations_url": "https://api.github.com/users/lectrical/orgs",
      "repos_url": "https://api.github.com/users/lectrical/repos",
      "events_url": "https://api.github.com/users/lectrical/events{/privacy}",
      "received_events_url": "https://api.github.com/users/lectrical/received_events",
      "type": "User",
      "user_view_type": "public",
      "site_admin": false
    },
    "committer": {
      "login": "web-flow",
      "id": 19864447,
      "node_id": "MDQ6VXNlcjE5ODY0NDQ3",
      "avatar_url": "https://avatars.githubusercontent.com/u/19864447?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/web-flow",
      "html_url": "https://github.com/web-flow",
      "followers_url": "https://api.github.com/users/web-flow/followers",
      "following_url": "https://api.github.com/users/web-flow/following{/other_user}",
      "gists_url": "https://api.github.com/users/web-flow/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/web-flow/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/web-flow/subscriptions",
      "organizations_url": "https://api.github.com/users/web-flow/orgs",
      "repos_url": "https://api.github.com/users/web-flow/repos",
      "events_url": "https://api.github.com/users/web-flow/events{/privacy}",
      "received_events_url": "https://api.github.com/users/web-flow/received_events",
      "type": "User",
      "user_view_type": "public",
      "site_admin": false
    },
    "parents": [
      {
        "sha": "8bcdc9304ace5f2cc9bf662ab8998d75537e05f0",
        "url": "https://api.github.com/repos/jqlang/jq/commits/8bcdc9304ace5f2cc9bf662ab8998d75537e05f0",
        "html_url": "https://github.com/jqlang/jq/commit/8bcdc9304ace5f2cc9bf662ab8998d75537e05f0"
      }
    ]
  },
  {
    "sha": "8bcdc9304ace5f2cc9bf662ab8998d75537e05f0",
    "node_id": "C_kwDOAE3WVdoAKDhiY2RjOTMwNGFjZTVmMmNjOWJmNjYyYWI4OTk4ZDc1NTM3ZTA1ZjA",
    "commit": {
      "author": {
        "name": "Emanuele Torre",
        "email": "torreemanuele6@gmail.com",
        "date": "2024-11-25T17:59:08Z"
      },
      "committer": {
        "name": "Emanuele Torre",
        "email": "torreemanuele6@gmail.com",
        "date": "2024-12-01T10:49:44Z"
      },
      "message": "jq_next: simplify CALL_BUILTIN implementation",
      "tree": {
        "sha": "f1282cbd1dccbc84d9cfd8ffa35babf457a8ae25",
        "url": "https://api.github.com/repos/jqlang/jq/git/trees/f1282cbd1dccbc84d9cfd8ffa35babf457a8ae25"
      },
      "url": "https://api.github.com/repos/jqlang/jq/git/commits/8bcdc9304ace5f2cc9bf662ab8998d75537e05f0",
      "comment_count": 0,
      "verification": {
        "verified": false,
        "reason": "unsigned",
        "signature": null,
        "payload": null,
        "verified_at": null
      }
    },
    "url": "https://api.github.com/repos/jqlang/jq/commits/8bcdc9304ace5f2cc9bf662ab8998d75537e05f0",
    "html_url": "https://github.com/jqlang/jq/commit/8bcdc9304ace5f2cc9bf662ab8998d75537e05f0",
    "comments_url": "https://api.github.com/repos/jqlang/jq/commits/8bcdc9304ace5f2cc9bf662ab8998d75537e05f0/comments",
    "author": {
      "login": "emanuele6",
      "id": 20175435,
      "node_id": "MDQ6VXNlcjIwMTc1NDM1",
      "avatar_url": "https://avatars.githubusercontent.com/u/20175435?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/emanuele6",
      "html_url": "https://github.com/emanuele6",
      "followers_url": "https://api.github.com/users/emanuele6/followers",
      "following_url": "https://api.github.com/users/emanuele6/following{/other_user}",
      "gists_url": "https://api.github.com/users/emanuele6/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/emanuele6/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/emanuele6/subscriptions",
      "organizations_url": "https://api.github.com/users/emanuele6/orgs",
      "repos_url": "https://api.github.com/users/emanuele6/repos",
      "events_url": "https://api.github.com/users/emanuele6/events{/privacy}",
      "received_events_url": "https://api.github.com/users/emanuele6/received_events",
      "type": "User",
      "user_view_type": "public",
      "site_admin": false
    },
    "committer": {
      "login": "emanuele6",
      "id": 20175435,
      "node_id": "MDQ6VXNlcjIwMTc1NDM1",
      "avatar_url": "https://avatars.githubusercontent.com/u/20175435?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/emanuele6",
      "html_url": "https://github.com/emanuele6",
      "followers_url": "https://api.github.com/users/emanuele6/followers",
      "following_url": "https://api.github.com/users/emanuele6/following{/other_user}",
      "gists_url": "https://api.github.com/users/emanuele6/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/emanuele6/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/emanuele6/subscriptions",
      "organizations_url": "https://api.github.com/users/emanuele6/orgs",
      "repos_url": "https://api.github.com/users/emanuele6/repos",
      "events_url": "https://api.github.com/users/emanuele6/events{/privacy}",
      "received_events_url": "https://api.github.com/users/emanuele6/received_events",
      "type": "User",
      "user_view_type": "public",
      "site_admin": false
    },
    "parents": [
      {
        "sha": "0b82b3841b05faefe0ac18379bdb361b7e4e3464",
        "url": "https://api.github.com/repos/jqlang/jq/commits/0b82b3841b05faefe0ac18379bdb361b7e4e3464",
        "html_url": "https://github.com/jqlang/jq/commit/0b82b3841b05faefe0ac18379bdb361b7e4e3464"
      }
    ]
  },
  {
    "sha": "0b82b3841b05faefe0ac18379bdb361b7e4e3464",
    "node_id": "C_kwDOAE3WVdoAKDBiODJiMzg0MWIwNWZhZWZlMGFjMTgzNzliZGIzNjFiN2U0ZTM0NjQ",
    "commit": {
      "author": {
        "name": "Emanuele Torre",
        "email": "torreemanuele6@gmail.com",
        "date": "2024-11-25T06:47:14Z"
      },
      "committer": {
        "name": "Emanuele Torre",
        "email": "torreemanuele6@gmail.com",
        "date": "2024-12-01T10:49:44Z"
      },
      "message": "builtin.c: typecheck builtin cfunctions in function_list\n\nIn C23 (default C standard used by GCC 15),  jv (*fptr)();  has become\nequivalent to  jv (*fptr)(void);  so we can no longer assign builtin\nimplemenations directly to the fptr member of cfunctions without\ngenerating a compile error.\n\nSince there does not seem to be any straight-forward way to tell\nautoconf to force the compiler to use C99 short of explicitly adding\n-std=c99 to CFLAGS, it is probably a cleaner solution to just make the\ncode C23 compatible.\n\nA possible solution could have been to just redeclare  cfunction.fptr\nas void*, but then the functions' return type would not have been type\nchecked (e.g. if you tried to add a {printf, \"printf\", 2}, where printf\nis a function that does not return jv, the compiler wouldn't have\ncomplained.)\nWe were already not typechecking the arguments of the functions, so e.g.\n  {binop_plus, \"_plus\", 3},  /* instead of {f_plus, \"_plus, 3},       */\n  {f_setpath, \"setpath\", 4}, /* instead of {f_setpath, \"setpath\", 3}, */\ncompile without errors despite not having the correct prototype.\n\nSo I thought of instead improving the situation by redefining\ncfunction.fptr as a union of function pointers with the prototypes that\nthe jq bytecode interpreter can call, and use a macro to add the builtin\nfunctions to function_list using to the arity argument to assign the\nimplementation function to the appropriate union member.\n\nNow the code won't compile if the wrong arity, or an arity not supported\nby the bytecode interpreter (>5 = 1input+4arguments), or a prototype not\njallable by the bytecode interpreter (e.g. binop_plus that doesn't\nexpect a  jq_state*  argument).\n\nAlso, the code now compiles with gcc -std=c23.\n\nFixes #3206",
      "tree": {
        "sha": "cd27325cd80e838f3dfdee37effe64cd1eab0018",
        "url": "https://api.github.com/repos/jqlang/jq/git/trees/cd27325cd80e838f3dfdee37effe64cd1eab0018"
      },
      "url": "https://api.github.com/repos/jqlang/jq/git/commits/0b82b3841b05faefe0ac18379bdb361b7e4e3464",
      "comment_count": 0,
      "verification": {
        "verified": false,
        "reason": "unsigned",
        "signature": null,
        "payload": null,
        "verified_at": null
      }
    },
    "url": "https://api.github.com/repos/jqlang/jq/commits/0b82b3841b05faefe0ac18379bdb361b7e4e3464",
    "html_url": "https://github.com/jqlang/jq/commit/0b82b3841b05faefe0ac18379bdb361b7e4e3464",
    "comments_url": "https://api.github.com/repos/jqlang/jq/commits/0b82b3841b05faefe0ac18379bdb361b7e4e3464/comments",
    "author": {
      "login": "emanuele6",
      "id": 20175435,
      "node_id": "MDQ6VXNlcjIwMTc1NDM1",
      "avatar_url": "https://avatars.githubusercontent.com/u/20175435?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/emanuele6",
      "html_url": "https://github.com/emanuele6",
      "followers_url": "https://api.github.com/users/emanuele6/followers",
      "following_url": "https://api.github.com/users/emanuele6/following{/other_user}",
      "gists_url": "https://api.github.com/users/emanuele6/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/emanuele6/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/emanuele6/subscriptions",
      "organizations_url": "https://api.github.com/users/emanuele6/orgs",
      "repos_url": "https://api.github.com/users/emanuele6/repos",
      "events_url": "https://api.github.com/users/emanuele6/events{/privacy}",
      "received_events_url": "https://api.github.com/users/emanuele6/received_events",
      "type": "User",
      "user_view_type": "public",
      "site_admin": false
    },
    "committer": {
      "login": "emanuele6",
      "id": 20175435,
      "node_id": "MDQ6VXNlcjIwMTc1NDM1",
      "avatar_url": "https://avatars.githubusercontent.com/u/20175435?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/emanuele6",
      "html_url": "https://github.com/emanuele6",
      "followers_url": "https://api.github.com/users/emanuele6/followers",
      "following_url": "https://api.github.com/users/emanuele6/following{/other_user}",
      "gists_url": "https://api.github.com/users/emanuele6/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/emanuele6/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/emanuele6/subscriptions",
      "organizations_url": "https://api.github.com/users/emanuele6/orgs",
      "repos_url": "https://api.github.com/users/emanuele6/repos",
      "events_url": "https://api.github.com/users/emanuele6/events{/privacy}",
      "received_events_url": "https://api.github.com/users/emanuele6/received_events",
      "type": "User",
      "user_view_type": "public",
      "site_admin": false
    },
    "parents": [
      {
        "sha": "96e8d893c10ed2f7656ccb8cfa39a9a291663a7e",
        "url": "https://api.github.com/repos/jqlang/jq/commits/96e8d893c10ed2f7656ccb8cfa39a9a291663a7e",
        "html_url": "https://github.com/jqlang/jq/commit/96e8d893c10ed2f7656ccb8cfa39a9a291663a7e"
      }
    ]
  },
  {
    "sha": "96e8d893c10ed2f7656ccb8cfa39a9a291663a7e",
    "node_id": "C_kwDOAE3WVdoAKDk2ZThkODkzYzEwZWQyZjc2NTZjY2I4Y2ZhMzlhOWEyOTE2NjNhN2U",
    "commit": {
      "author": {
        "name": "itchyny",
        "email": "itchyny@cybozu.co.jp",
        "date": "2024-11-20T23:12:18Z"
      },
      "committer": {
        "name": "GitHub",
        "email": "noreply@github.com",
        "date": "2024-11-20T23:12:18Z"
      },
      "message": "fix: reduce/foreach state variable should not be reset each iteration (#3205)\n\nWhen the UPDATE query in a reduce or foreach syntax does not emit\r\na value, it should simply skip updating the state variable. But\r\ncurrently it resets the state variable to null.",
      "tree": {
        "sha": "98540e1dbd3e948c6e17596180133ceb5af12517",
        "url": "https://api.github.com/repos/jqlang/jq/git/trees/98540e1dbd3e948c6e17596180133ceb5af12517"
      },
      "url": "https://api.github.com/repos/jqlang/jq/git/commits/96e8d893c10ed2f7656ccb8cfa39a9a291663a7e",
      "comment_count": 0,
      "verification": {
        "verified": true,
        "reason": "valid",
        "signature": "-----BEGIN PGP SIGNATURE-----\n\nwsFcBAABCAAQBQJnPmzSCRC1aQ7uu5UhlAAAHKEQAIxmkTf8Q6SFR6/sAnNO5JLA\nHKbdObrCno4z4HwGGL+l8UkiIQvzaazbjkO8vCmkvS/KlcXigGz6uI6KliA3A2jr\n0aAcTW7ru97o5vv4ggwiebe/u9R8AID95wGEbDOzKfeVU40XOTzevAVrsHCPpfk2\nFRsu6my4v1Dxu7RzLRC/hkY7mdqtTkcM9irY6WeW0VQicP929WEzjWBCs2LPukMt\nU3RP2nTIdKzKVYl6Fp29ulJ1Kx3BqCo5MGz2k8WrOyoyDbAx/vJ1Ydkhl9br/IC+\n1SOcqSGeeVR6z+raBlMtm+kZzdQJcxxaUgsScVH2mVV/t+9e+AnDQ3CXIoa0byrQ\nyC1zeo7TbhRJkwQtp7w7UtnYFXMTH1lJJ//jDue/pUUGd6sJx1JifTveU/57ljr4\n9N875DUpj07iZNK2U3ZQo16wJ4QTW8k5iP7QY7z2hZvCyJT36qdJAitV6EhwkxOx\nWrC2tPwUTSCQP5mWXfvKSbIF1QTLiU2Er36Qb9umETmevZo6ssMzQGX29p+XIu37\nt7XfwjEhrApLK19zn+5N/g7hYR/AUILjgYUEIBFDsgPVmke+wciT1RiGy9An1H22\n4lh+Rj2X52XyZNnnBA/u8f4yzD3990VHHG8q/pYZwhqOjFe/YZwzMO/ixiuIPflP\ntqeiNDqZmMjI65EnC4kw\n=opV+\n-----END PGP SIGNATURE-----\n",
        "payload": "tree 98540e1dbd3e948c6e17596180133ceb5af12517\nparent 8619f8a8ac746887cf43abd8e2116abba253cdcc\nauthor itchyny <itchyny@cybozu.co.jp> 1732144338 +0900\ncommitter GitHub <noreply@github.com> 1732144338 +0900\n\nfix: reduce/foreach state variable should not be reset each iteration (#3205)\n\nWhen the UPDATE query in a reduce or foreach syntax does not emit\r\na value, it should simply skip updating the state variable. But\r\ncurrently it resets the state variable to null.",
        "verified_at": "2024-11-20T23:12:21Z"
      }
    },
    "url": "https://api.github.com/repos/jqlang/jq/commits/96e8d893c10ed2f7656ccb8cfa39a9a291663a7e",
    "html_url": "https://github.com/jqlang/jq/commit/96e8d893c10ed2f7656ccb8cfa39a9a291663a7e",
    "comments_url": "https://api.github.com/repos/jqlang/jq/commits/96e8d893c10ed2f7656ccb8cfa39a9a291663a7e/comments",
    "author": {
      "login": "itchyny",
      "id": 375258,
      "node_id": "MDQ6VXNlcjM3NTI1OA==",
      "avatar_url": "https://avatars.githubusercontent.com/u/375258?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/itchyny",
      "html_url": "https://github.com/itchyny",
      "followers_url": "https://api.github.com/users/itchyny/followers",
      "following_url": "https://api.github.com/users/itchyny/following{/other_user}",
      "gists_url": "https://api.github.com/users/itchyny/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/itchyny/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/itchyny/subscriptions",
      "organizations_url": "https://api.github.com/users/itchyny/orgs",
      "repos_url": "https://api.github.com/users/itchyny/repos",
      "events_url": "https://api.github.com/users/itchyny/events{/privacy}",
      "received_events_url": "https://api.github.com/users/itchyny/received_events",
      "type": "User",
      "user_view_type": "public",
      "site_admin": false
    },
    "committer": {
      "login": "web-flow",
      "id": 19864447,
      "node_id": "MDQ6VXNlcjE5ODY0NDQ3",
      "avatar_url": "https://avatars.githubusercontent.com/u/19864447?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/web-flow",
      "html_url": "https://github.com/web-flow",
      "followers_url": "https://api.github.com/users/web-flow/followers",
      "following_url": "https://api.github.com/users/web-flow/following{/other_user}",
      "gists_url": "https://api.github.com/users/web-flow/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/web-flow/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/web-flow/subscriptions",
      "organizations_url": "https://api.github.com/users/web-flow/orgs",
      "repos_url": "https://api.github.com/users/web-flow/repos",
      "events_url": "https://api.github.com/users/web-flow/events{/privacy}",
      "received_events_url": "https://api.github.com/users/web-flow/received_events",
      "type": "User",
      "user_view_type": "public",
      "site_admin": false
    },
    "parents": [
      {
        "sha": "8619f8a8ac746887cf43abd8e2116abba253cdcc",
        "url": "https://api.github.com/repos/jqlang/jq/commits/8619f8a8ac746887cf43abd8e2116abba253cdcc",
        "html_url": "https://github.com/jqlang/jq/commit/8619f8a8ac746887cf43abd8e2116abba253cdcc"
      }
    ]
  }
]	
`
