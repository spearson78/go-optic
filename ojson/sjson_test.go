package ojson_test

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/ojson"
)

func TestGJSON(t *testing.T) {

	//These tests are derived from https://github.com/tidwall/gjson

	const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

	if value, ok, err := GetFirst(Compose(ParseString[any](), Key("name").Key("last").String()), json); !ok || err != nil || value != "Prichard" {
		t.Fatalf("lastName %v %v", value, ok)
	}

	//Path Syntax
	const pathSyntaxJson = `{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}  
`

	if value, ok, err := GetFirst(Compose(ParseString[any](), Key("name").Key("last").String()), pathSyntaxJson); err != nil || !ok || value != "Anderson" {
		t.Fatalf("name.last %v %v %v", value, ok, err)
	}

	if value, ok, err := GetFirst(Compose(ParseString[any](), Key("age").Float()), pathSyntaxJson); err != nil || !ok || value != 37 {
		t.Fatalf("age %v %v %v", value, ok, err)
	}

	if value, ok, err := GetFirst(Compose(ParseString[any](), Key("children").Array()), pathSyntaxJson); err != nil || !ok || !reflect.DeepEqual(value, []any{"Sara", "Alex", "Jack"}) {
		t.Fatalf("children %v %v %v", value, ok, err)
	}

	if value, err := Get(Length(Compose(ParseString[any](), Key("children").Traverse())), pathSyntaxJson); err != nil || value != 3 {
		t.Fatalf("children.# %v %v", value, err)
	}

	if value, ok, err := GetFirst(Compose(ParseString[any](), Key("children").Nth(1).String()), pathSyntaxJson); err != nil || !ok || value != "Alex" {
		t.Fatalf("children %v %v %v", value, ok, err)
	}

	//"child*.2"           >> "Jack"
	if value, ok, err := GetFirst(Compose(ParseString[any](), KeyLike(regexp.MustCompile("child.*")).Nth(2)), pathSyntaxJson); err != nil || !ok || value != "Jack" {
		t.Fatalf("child*.2 %v %v", value, ok)
	}

	//"c?ildren.0"         >> "Sara"
	if value, ok, err := GetFirst(Compose(ParseString[any](), KeyLike(regexp.MustCompile("c.ildren")).Nth(0)), pathSyntaxJson); err != nil || !ok || value != "Sara" {
		t.Fatalf("c?ildren.0 %v %v", value, ok)
	}

	//"fav\.movie"         >> "Deer Hunter"
	if value, ok, err := GetFirst(Compose(ParseString[any](), Key("fav.movie")), pathSyntaxJson); err != nil || !ok || value != "Deer Hunter" {
		t.Fatalf("fav\\.movie %v %v", value, ok)
	}

	//"friends.#.first"    >> ["Dale","Roger","Jane"]
	if value, err := Get(SliceOf(Compose(ParseString[any](), Key("friends").Traverse().Key("first")), 10), pathSyntaxJson); err != nil || !reflect.DeepEqual(value, []any{"Dale", "Roger", "Jane"}) {
		t.Fatalf("friends.#.first %v %v", value, err)
	}

	//"friends.1.last"     >> "Craig"
	if value, ok, err := GetFirst(Compose(ParseString[any](), Key("friends").Nth(1).Key("last")), pathSyntaxJson); err != nil || !ok || value != "Craig" {
		t.Fatalf("friends.#.first %v %v", value, err)
	}

	//friends.#(last=="Murphy").first    >> "Dale"
	if value, ok, err := GetFirst(
		Compose3(
			ParseString[any](),
			Filtered(
				KeyE("friends").TraverseE(),
				KeyE("last").Eq("Murphy"),
			),
			KeyE("first"),
		), pathSyntaxJson); err != nil || !ok || value != "Dale" {
		t.Fatalf("friends.#(last==Murphy).first %v %v", value, err)
	}

	//friends.#(last=="Murphy")#.first   >> ["Dale","Jane"]
	if value, err := Get(SliceOf(Compose(ParseString[any](), Compose(Filtered(Key("friends").Traverse(), Key("last").Eq("Murphy")), Key("first"))), 10), pathSyntaxJson); err != nil || !reflect.DeepEqual(value, []any{"Dale", "Jane"}) {
		t.Fatalf("friends.#(last==Murphy)#.first %v %v", value, err)
	}

	//friends.#(age>45)#.last            >> ["Craig","Murphy"]
	if value, err := Get(SliceOf(Compose(ParseString[any](), Compose(Filtered(Key("friends").Traverse(), Key("age").Gt(45)), Key("last"))), 10), pathSyntaxJson); err != nil || !reflect.DeepEqual(value, []any{"Craig", "Murphy"}) {
		t.Fatalf("friends.#(age>45)#.last %v %v", value, err)
	}

	//friends.#(first%"D*").last         >> "Murphy"
	if value, ok, err := GetFirst(Compose(ParseString[any](), Compose(Filtered(Key("friends").Traverse(), Key("first").Like(regexp.MustCompile("D.*"))), Key("last"))), pathSyntaxJson); err != nil || !ok || value != "Murphy" {
		t.Fatalf(`friends.#(first%%"D*").last  %v`, value)
	}

	//friends.#(first!%"D*").last        >> "Craig"
	if value, ok, err := GetFirst(Compose(ParseString[any](), Compose(Filtered(Key("friends").Traverse(), Compose(Key("first").Like(regexp.MustCompile("D.*")), Not())), Key("last"))), pathSyntaxJson); err != nil || !ok || value != "Craig" {
		t.Fatalf(`friends.#(first%%"D*").last  %v`, value)
	}

	//friends.#(nets.#(=="fb"))#.first   >> ["Dale","Roger"]
	if value, err := Get(SliceOf(Compose(ParseString[any](), Compose(Filtered(Key("friends").Traverse(), Any(Key("nets").Traverse(), Eq[any]("fb"))), Key("first"))), 10), pathSyntaxJson); err != nil || !reflect.DeepEqual(value, []any{"Dale", "Roger"}) {
		t.Fatalf(`friends.#(nets.#(=="fb"))#.first  %v`, value)
	}

	//"children|@reverse"           >> ["Jack","Alex","Sara"]
	if value, err := Get(SliceOf(Reversed(Compose(ParseString[any](), Key("children").Traverse())), 10), pathSyntaxJson); err != nil || !reflect.DeepEqual(value, []any{"Jack", "Alex", "Sara"}) {
		t.Fatalf(`children|@reverse"  %v`, value)
	}

	//children|@reverse|0
	if value, err := Get(SliceOf(Element(Reversed(Compose(ParseString[any](), Key("children").Traverse())), 0), 10), pathSyntaxJson); err != nil || !reflect.DeepEqual(value, []any{"Jack"}) {
		t.Fatalf(`children|@reverse|0"  %v %v`, value, err)
	}

	//"children|@case:upper"           >> ["SARA","ALEX","JACK"]
	if value, err := Get(SliceOf(Compose(Compose(ParseString[any](), Key("children").Traverse().String()), Op(strings.ToUpper)), 10), pathSyntaxJson); err != nil || !reflect.DeepEqual(value, []string{"SARA", "ALEX", "JACK"}) {
		t.Fatalf(`children|@case:upper"  %v %v`, value, err)
	}

	//"children|@case:lower|@reverse"  >> ["jack","alex","sara"]
	if value, err := Get(SliceOf(Reversed(Compose(Compose(ParseString[any](), Key("children").Traverse().String()), Op(strings.ToLower))), 10), pathSyntaxJson); err != nil || !reflect.DeepEqual(value, []string{"jack", "alex", "sara"}) {
		t.Fatalf(`children|@case:lower|@reverse  %v`, value)
	}

}

func TestSJSON(t *testing.T) {

	//These tests are derived from https://github.com/tidwall/sjson

	const json1 = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	if newJson, err := Set(Compose(ParseString[any](), Key("name").Key("last")), "Anderson", json1); err != nil || newJson == `{"name":{"first":"Janet","last":"Anderson"},"age":47}` {
		t.Fatalf(`Set(json, "name.last", "Anderson") %v %v`, newJson, err)
	}
	if newJson, err := Set(Compose(ParseString[any](), Key("name")), "Tom", "{}"); err != nil || newJson != `{"name":"Tom"}` {
		t.Fatalf(`Set("", "name", "Tom") %v %v`, newJson, err)
	}

	if newJson, err := Set(Compose(ParseString[any](), Key("name").Key("first")), "Sara", `{"name":{"last":"Anderson"}}`); err != nil || newJson != `{"name":{"first":"Sara","last":"Anderson"}}` {
		t.Fatalf(`"name.first", "Sara" %v %v`, newJson, err)
	}

	if newJson, err := Set(Compose(ParseString[any](), Key("name").Key("last")), "Smith", `{"name":{"last":"Anderson"}}`); err != nil || newJson != `{"name":{"last":"Smith"}}` {
		t.Fatalf(`"name.last", "Smith" %v %v`, newJson, err)
	}

	if newJson, err := Set(Compose(ParseString[any](), Key("friends").Nth(2)), "Sara", `{"friends":["Andy","Carol"]}`); err != nil || newJson != `{"friends":["Andy","Carol","Sara"]}` {
		t.Fatalf(`"friends.2", "Sara" %v %v`, newJson, err)
	}

	if newJson, err := Modify(Compose(ParseString[any](), Key("friends").Array()), AppendSlice[any](ValCol(any("Sara"))), `{"friends":["Andy","Carol"]}`); err != nil || newJson != `{"friends":["Andy","Carol","Sara"]}` {
		t.Fatalf(`"friends.-1", "Sara" %v %v`, newJson, err)
	}

	//value, _ := sjson.Set(`{"friends":["Andy","Carol"]}`, "friends.4", "Sara")

	if newJson, err := Set(Compose(ParseString[any](), Key("friends").Nth(4)), "Sara", `{"friends":["Andy","Carol"]}`); err != nil || newJson != `{"friends":["Andy","Carol",null,null,"Sara"]}` {
		t.Fatalf(`"friends.4", "Sara" %v %v`, newJson, err)
	}

	//value, _ := sjson.Delete(`{"name":{"first":"Sara","last":"Anderson"}}`, "name.first")

	if value, err := Modify(Compose(ParseString[any](), Key("name").Object()), FilteredMapI[string, any](NeI[any]("first")), `{"name":{"first":"Sara","last":"Anderson"}}`); err != nil || value != `{"name":{"last":"Anderson"}}` {
		t.Fatalf(`"Delete name.first %v %v`, value, err)
	}

	//value, _ := sjson.Delete(`{"friends":["Andy","Carol"]}`, "friends.1")

	if value, err := Modify(Compose(ParseString[any](), Key("friends").Array()), FilteredSliceI[any](NeI[any](1)), `{"friends":["Andy","Carol"]}`); err != nil || value != `{"friends":["Andy"]}` {
		t.Fatalf(`"Delete name.first %v %v`, value, err)
	}

	//value, _ := sjson.Delete(`{"friends":["Andy","Carol"]}`, "friends.-1")

	if value, err := Modify(Compose(ParseString[any](), Key("friends").Array()), SubSlice[any](0, -1), `{"friends":["Andy","Carol"]}`); err != nil || value != `{"friends":["Andy"]}` {
		t.Fatalf(`"friends.-1" %v %v`, value, err)
	}

}
