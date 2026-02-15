package optic_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/samber/mo"

	. "github.com/spearson78/go-optic"
)

func ExamplePrism() {

	//This prism converts an any to a string
	//the Cast function implements a generic version of this example
	anyToString := Prism[any, string](
		func(source any) (string, bool) {
			str, ok := source.(string)
			return str, ok
		},
		func(focus string) any {
			return focus
		},
		ExprCustom("ExamplePrism"),
	)

	var str string
	var isMatch bool
	str, isMatch = MustGetFirst(anyToString, "Lorem ipsum")
	fmt.Println(str, isMatch)

	str, isMatch = MustGetFirst(anyToString, nil)
	fmt.Println(isMatch)

	//If the given value is a string then convert it to upper case
	var resStr any = MustModify(anyToString, Op(strings.ToUpper), "Lorem ipsum")
	fmt.Println(resStr)

	//If the value is not a string the the prism returns the original value back
	var resInt any = MustModify(anyToString, Op(strings.ToUpper), 1)
	fmt.Println(resInt)

	//Output:
	//Lorem ipsum true
	//false
	//LOREM IPSUM
	//1
}

func ExamplePrismP() {

	//This prism converts an any to a string and an int to an any
	//the CastP function implements a generic version of this example
	anyToString := PrismP[any, any, string, int](
		func(source any) mo.Either[any, string] {
			if str, ok := source.(string); ok {
				return mo.Right[any, string](str)
			} else {
				return mo.Left[any, string](source)
			}

		},
		func(focus int) any {
			return focus
		},
		ExprCustom("ExamplePrism"),
	)

	var str string
	var isMatch bool
	str, isMatch = MustGetFirst(anyToString, "Lorem ipsum")
	fmt.Println(str, isMatch)

	//If the given value is a string then parse it and multiply it by 2
	var resStr any
	resStr, err := Modify(Compose(anyToString, ParseIntP[int](10, 32)), Mul(2), "100")
	fmt.Printf("%v : %T : %v\n", resStr, resStr, err)

	//If the value is not a string the the prism returns the original value back
	var resInt any
	resInt, err = Modify(Compose(anyToString, ParseIntP[int](10, 32)), Mul(2), true)
	fmt.Printf("%v : %T : %v\n", resInt, resInt, err)

	//Output:
	//Lorem ipsum true
	//200 : int : <nil>
	//true : bool : <nil>

}

func ExamplePrismE() {

	//This prism umarshals json and only focuses on json objects
	//other json types (string,array,number) are ignored
	//Json parse errors are returned
	jsonPrism := PrismE[string, map[string]any](
		func(ctx context.Context, source string) (map[string]any, bool, error) {
			var a any
			err := json.Unmarshal([]byte(source), &a)
			if err != nil {
				return nil, true, err
			}

			switch t := a.(type) {
			case map[string]any:
				return t, true, nil
			default:
				return nil, false, nil
			}
		},
		func(ctx context.Context, focus map[string]any) (string, error) {
			bytes, err := json.Marshal(focus)
			return string(bytes), err
		},
		ExprCustom("ExamplePrism"),
	)

	var jsonObj map[string]any
	var isMatch bool
	jsonObj, isMatch, err := GetFirst(jsonPrism, `{"key":"value"}`)
	fmt.Println(jsonObj, isMatch, err)

	//Parse the json obj, Traverse over the properties, cast them to string and convert them to uppercase
	var resStr string
	resStr, err = Modify(
		Compose3(
			jsonPrism,
			TraverseMap[string, any](),
			DownCast[any, string](),
		),
		Op(strings.ToUpper),
		`{"key":"value"}`,
	)
	fmt.Println(resStr, err)

	//If the value is valid json but is not a json obj then the prism returns the original value back
	resStr, err = Modify(
		Compose3(
			jsonPrism,
			TraverseMap[string, any](),
			DownCast[any, string](),
		),
		Op(strings.ToUpper),
		`1`,
	)
	fmt.Println(resStr, err)

	//If the source is invalid json then an error is returned
	_, err = Modify(
		Compose3(
			jsonPrism,
			TraverseMap[string, any](),
			DownCast[any, string](),
		),
		Op(strings.ToUpper),
		`invalid json`,
	)
	fmt.Println(err)

	//If the source is valid but we create an invalid json object then the marshal error is returned
	_, err = Set(
		Compose(
			jsonPrism,
			TraverseMap[string, any](),
		),
		any(func() {}),
		`{"key":"value"}`,
	)
	fmt.Println(err)

	//Output:
	//map[key:value] true <nil>
	//{"key":"VALUE"} <nil>
	//1 <nil>
	//invalid character 'i' looking for beginning of value
	//optic error path:
	//	Custom(ExamplePrism)
	//
	//json: unsupported type: func()
	//optic error path:
	//	Custom(ExamplePrism)
}

func ExamplePrismIEP() {

	//This is a rather contrived example.
	//This Prism unmarshals a json string and casts it to a map[string]any
	//other json types (string,array,number) are ignored
	//Json parse errors are returned
	//Under modification the prism expects a string instead of a map[string]any
	//And returns an []byte instead of a string.
	jsonPrism := PrismIEP[Void, string, []byte, map[string]any, string](
		func(ctx context.Context, source string) (mo.Either[[]byte, ValueI[Void, map[string]any]], error) {
			var a any
			err := json.Unmarshal([]byte(source), &a)
			if err != nil {
				return mo.Left[[]byte, ValueI[Void, map[string]any]]([]byte(source)), err
			}

			switch t := a.(type) {
			case map[string]any:
				return mo.Right[[]byte, ValueI[Void, map[string]any]](ValI(Void{}, t)), nil
			default:
				return mo.Left[[]byte, ValueI[Void, map[string]any]]([]byte(source)), err
			}
		},
		func(ctx context.Context, focus string) ([]byte, error) {
			bytes, err := json.Marshal(focus)
			return bytes, err
		},
		IxMatchVoid(),
		ExprCustom("ExamplePrism"),
	)

	var jsonObj map[string]any
	var isMatch bool
	jsonObj, isMatch, err := GetFirst(jsonPrism, `{"key":"value"}`)
	fmt.Println(jsonObj, isMatch, err)

	var resStr []byte
	resStr, err = Modify( //return type is []byte
		jsonPrism,
		Op(func(focus map[string]any) string { //convert from map to string
			return fmt.Sprint(focus)
		}),
		`{"key":"value"}`, //String input
	)

	//The result is a byte array containing the go string representation of a map embedded inside a json string.....
	fmt.Println(string(resStr), err)

	//Output:
	//map[key:value] true <nil>
	//"map[key:value]" <nil>
}
