package ojson_test

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/ojson"
)

type Line struct {
	Start Point
	End   Point
}

type Point struct {
	X float64
	Y float64
}

var LineStartField = FieldLens(func(s *Line) *Point { return &s.Start })
var LineEndField = FieldLens(func(s *Line) *Point { return &s.End })

var PointXField = FieldLens(func(s *Point) *float64 { return &s.X })
var PointYField = FieldLens(func(s *Point) *float64 { return &s.Y })

func TestGet(t *testing.T) {

	jsonLine := Parse[Line]()

	b := []byte(`
	{
		"Start": {
			"X": 1,
			"Y": 2
		},
		"End": {
			"X": 3,
			"Y": 4
		}
	}
`)

	line, err := Get(jsonLine, b)

	if err != nil {
		t.Fatalf("error decoding valid json %v", err)
	}

	if line.Start.X != 1 {
		t.Fatalf("incorrect start x value %v", line.Start.X)
	}

	if line.Start.Y != 2 {
		t.Fatalf("incorrect start y value %v", line.Start.Y)
	}

	if line.End.X != 3 {
		t.Fatalf("incorrect start x value %v", line.End.X)
	}

	if line.End.Y != 4 {
		t.Fatalf("incorrect start y value %v", line.End.Y)
	}

}

func TestReverseGet(t *testing.T) {

	jsonLine := AsReverseGet(Parse[Line]())

	line := Line{
		Start: Point{
			X: 1,
			Y: 2,
		},
		End: Point{
			X: 3,
			Y: 4,
		},
	}

	json, err := Get(jsonLine, line)

	if err != nil {
		t.Fatalf("error encoding valid object %v", err)
	}

	if string(json) != "{\"Start\":{\"X\":1,\"Y\":2},\"End\":{\"X\":3,\"Y\":4}}" {
		t.Fatalf("Incorrect encoded json %v", string(json))
	}

}

func TestComposition(t *testing.T) {

	jsonLine := Parse[Line]()

	json := "{\"Start\":{\"X\":1,\"Y\":2},\"End\":{\"X\":3,\"Y\":4}}"

	c1 := Compose(LineStartField, PointXField)
	c2 := Compose(jsonLine, c1)
	sb1 := IsoCast[string, []byte]()

	LineStartXField := Compose(sb1, c2)

	getVal, getErr := Get(LineStartXField, json)

	if getErr != nil {
		t.Fatalf("error getting from valid json %v", getErr)
	}

	if getVal != 1 {
		t.Fatalf("incorrect start x value %v", getVal)
	}

	newJson, newErr := Set(LineStartXField, 1.2, json)

	if newErr != nil {
		t.Fatalf("error getting from valid json %v", getErr)
	}

	if newJson != "{\"Start\":{\"X\":1.2,\"Y\":2},\"End\":{\"X\":3,\"Y\":4}}" {
		t.Fatalf("incorrect new json %v", newJson)
	}

}

func TestStringComposition(t *testing.T) {

	jsonLine := Compose(IsoCast[string, []byte](), Parse[Line]())

	json := "{\"Start\":{\"X\":1,\"Y\":2},\"End\":{\"X\":3,\"Y\":4}}"

	LineStartXField := Compose3(jsonLine, LineStartField, PointXField)

	getVal, getErr := Get(LineStartXField, json)

	if getErr != nil {
		t.Fatalf("error getting from valid json %v", getErr)
	}

	if getVal != 1 {
		t.Fatalf("incorrect start x value %v", getVal)
	}

	newJson, newErr := Set(LineStartXField, 1.2, json)

	if newErr != nil {
		t.Fatalf("error getting from valid json %v", getErr)
	}

	if newJson != "{\"Start\":{\"X\":1.2,\"Y\":2},\"End\":{\"X\":3,\"Y\":4}}" {
		t.Fatalf("incorrect new json %v", newJson)
	}

}

func TestObject(t *testing.T) {

	jsonObject := `
{
	"name": "Jack Sparrow",
	"rank": "Captain"	
}	
	`

	if r, err := Get(Compose(ParseString[any](), ObjectE()), jsonObject); err != nil || !reflect.DeepEqual(r, map[string]any{"name": "Jack Sparrow", "rank": "Captain"}) {
		t.Fatalf("Object 1 %v : %v", r, err)
	}

	jsonArray := `
[
	"North",
	"East",
	"South",
	"West"
]	
		`

	if r, err := Get(Compose(ParseString[any](), ArrayE()), jsonArray); err != nil || !reflect.DeepEqual(r, []any{"North", "East", "South", "West"}) {
		t.Fatalf("Array 1 %v : %v", r, err)
	}

}

func TestStructure(t *testing.T) {

	jsonObject := `
{
	"name": "Black Pearl",
	"crew": [
		{
			"name": "Jack Sparrow",
			"rank": "Captain"
		},
		{
			"name": "Will Turner",
			"rank": "First Mate"
		}
	]
}	
	`

	path := Compose(ParseString[any](), Key("crew").Nth(0).Key("name").String())

	if r, found, err := GetFirst(path, jsonObject); !found || err != nil || r != "Jack Sparrow" {
		t.Fatalf("Object 1 %v : %v", r, err)
	}

	if r, err := Set(path, "BOB", jsonObject); err != nil || r != `{"crew":[{"name":"BOB","rank":"Captain"},{"name":"Will Turner","rank":"First Mate"}],"name":"Black Pearl"}` {
		t.Fatalf("Object 2 %v : %v", r, err)
	}
}

func TestMissingKeys(t *testing.T) {
	//This looks strange but Key uses AtMap which enables keys to be created.
	const json1 = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	if res, ok, err := GetFirst(Compose(ParseString[any](), Key("name").Key("missing").Key("missing2")), json1); !ok || err != nil || res != nil {
		t.Fatalf(`Get(json, "name.missing.missing2") %v %v %v`, res, err, ok)
	}

	if res, err := Set(Compose(ParseString[any](), Key("name").Key("missing").Key("missing2")), "created", json1); err != nil || res != `{"age":47,"name":{"first":"Janet","last":"Prichard","missing":{"missing2":"created"}}}` {
		t.Fatalf(`Set(json, "name.missing.missing2" , "created") %v %v `, res, err)
	}

	//KeyR is also able to create keys
	if res, ok, err := GetFirst(Compose(ParseString[any](), KeyE("name").KeyE("missing").KeyE("missing2")), json1); !ok || err != nil || res != nil {
		t.Fatalf(`Get(json, "name.missing.missing2") %v %v %v`, res, err, ok)
	}

	if res, err := Set(Compose(ParseString[any](), KeyE("name").KeyE("missing").KeyE("missing2")), "created", json1); err != nil || res != `{"age":47,"name":{"first":"Janet","last":"Prichard","missing":{"missing2":"created"}}}` {
		t.Fatalf(`Set(json, "name.missing.missing2" , "created") %v %v `, res, err)
	}

	const json2 = `["Janet Prichard", 47]`
	//Nothing will be focused as the any is not a map.
	if res, ok, err := GetFirst(Compose(ParseString[any](), Key("name").Key("missing").Key("missing2")), json2); ok || err != nil || res != nil {
		t.Fatalf(`Get(json2, "name.missing.missing2") %v %v %v`, res, err, ok)
	}

	//As nothing is focused the result is not modified
	if res, err := Set(Compose(ParseString[any](), Key("name").Key("missing").Key("missing2")), "created", json2); err != nil || res != `["Janet Prichard",47]` {
		t.Fatalf(`Set(json2, "name.missing.missing2" , "created") %v %v `, res, err)
	}

	var castErr *ErrCastOf
	if res, ok, err := GetFirst(Compose(ParseString[any](), KeyE("name").KeyE("missing").KeyE("missing2")), json2); ok || !errors.As(err, &castErr) || res != nil {
		t.Fatalf(`Get(json2, "name.missing.missing2") %v %v %v`, res, err, ok)
	}

	castErr = nil
	if res, err := Set(Compose(ParseString[any](), KeyE("name").KeyE("missing").KeyE("missing2")), "created", json2); !errors.As(err, &castErr) {
		t.Fatalf(`Set(json2, "name.missing.missing2" , "created") %v %v `, res, err)
	}
}
