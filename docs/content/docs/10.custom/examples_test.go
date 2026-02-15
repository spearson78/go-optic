package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"iter"
	"reflect"
	"strconv"
	"testing"

	. "github.com/spearson78/go-optic"
)

func TestCustomGetter1(t *testing.T) {
	//BEGIN getter1
	sliceLen := Getter(
		func(source []string) int {
			return len(source)
		},
		ExprCustom("sliceLen"),
	)

	result := MustGet(sliceLen, []string{"alpha", "beta"})
	fmt.Println(result)
	//END getter1

	if result != 2 {
		t.Fatal(result)
	}
}

func TestCustomGetterI(t *testing.T) {
	//BEGIN getterI
	sliceLast := GetterI[int, []string, string](
		func(source []string) (int, string) {
			lastIndex := len(source) - 1
			return lastIndex, source[lastIndex]
		},
		func(indexA, indexB int) bool {
			return indexA == indexB
		},
		ExprCustom("sliceLast"),
	)

	index, result := MustGetI(sliceLast, []string{"alpha", "beta"})
	fmt.Println(result)
	//END getterI

	if index != 1 || result != "beta" {
		t.Fatal(index, result)
	}
}

func TestCustomGetterE(t *testing.T) {
	//BEGIN getterE
	parseInt := GetterE[string, int64](
		func(ctx context.Context, source string) (int64, error) {
			return strconv.ParseInt(source, 10, 0)
		},
		ExprCustom("parseInt"),
	)

	result, err := Get(parseInt, "alpha")
	fmt.Println(result, err)
	//END getterE

	if err == nil {
		t.Fatal(result, err)
	}
}

func TestCustomIsoEP(t *testing.T) {
	//BEGIN isoep
	parseInt := IsoEP[string, int64, int64, int64](
		func(ctx context.Context, source string) (int64, error) {
			return strconv.ParseInt(source, 10, 0)
		},
		func(ctx context.Context, focus int64) (int64, error) {
			return focus, nil
		},
		ExprCustom("parseInt"),
	)

	result, err := Modify(parseInt, Add[int64](10), "1")
	fmt.Println(result, err)
	//END isoep

	if result != int64(11) {
		t.Fatal(result, err)
	}
}

func TestCustomGetterOp(t *testing.T) {
	//BEGIN getter_op
	sliceLen := Op(
		func(source []string) int {
			return len(source)
		},
	)

	result := MustGet(sliceLen, []string{"alpha", "beta"})
	fmt.Println(result)
	//END getter_op
	if result != 2 {
		t.Fatal(result)
	}
}

//BEGIN lens_types

type ExampleStruct struct {
	name    string
	address string
}

//END lens_types

func TestCustomLensFieldLens(t *testing.T) {
	//BEGIN lens_fieldlens
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
	//END lens_fieldlens

	if !reflect.DeepEqual(result,
		ExampleStruct{
			name:    "Erika Mustermann",
			address: "Musterstadt",
		},
	) {
		t.Fatal(result)
	}
}

func TestCustomLensExample(t *testing.T) {
	//BEGIN lens_customfieldlens
	customFieldLens := Lens(
		//Lens getter
		func(source ExampleStruct) string {
			return source.name
		},
		//Lens setter
		func(newValue string, source ExampleStruct) ExampleStruct {
			source.name = newValue
			return source
		},
		ExprCustom("customFieldLens"),
	)

	data := ExampleStruct{
		name:    "Max Mustermann",
		address: "Musterstadt",
	}

	result := MustSet(customFieldLens, "Erika Mustermann", data)
	fmt.Println(result)
	//END lens_customfieldlens

	if !reflect.DeepEqual(result,
		ExampleStruct{
			name:    "Erika Mustermann",
			address: "Musterstadt",
		},
	) {
		t.Fatal(result)
	}
}

func TestCustomIso(t *testing.T) {
	//BEGIN custom_iso
	celsiusToFahrenheit := Iso(
		//Iso getter
		func(celsius float64) float64 {
			return (celsius * 1.8) + 32
		},
		//Iso reverse getter
		func(fahrenheit float64) float64 {
			return (fahrenheit - 32) / 1.8
		},
		ExprCustom("celsiusToFahrenheit"),
	)

	fahrenHeit := MustGet(celsiusToFahrenheit, 32)
	fmt.Println(fahrenHeit)

	celsius := MustReverseGet(celsiusToFahrenheit, 89.6)
	fmt.Println(celsius)
	//END custom_iso

	if !reflect.DeepEqual([]any{fahrenHeit, celsius}, []any{
		89.6,
		31.999999999999996,
	},
	) {
		t.Fatal(fahrenHeit, celsius)
	}
}

func TestCustomIsoMath(t *testing.T) {
	//BEGIN custom_iso_math
	celsiusToFahrenheit := Compose(
		Mul(1.8),
		Add(32.0),
	)

	fahrenHeit := MustGet(celsiusToFahrenheit, 32)
	fmt.Println(fahrenHeit)

	celsius := MustReverseGet(celsiusToFahrenheit, 89.6)
	fmt.Println(celsius)
	//END custom_iso_math

	if !reflect.DeepEqual([]any{fahrenHeit, celsius}, []any{
		89.6,
		31.999999999999996,
	},
	) {
		t.Fatal(fahrenHeit, celsius)
	}
}

func TestCustomIteration(t *testing.T) {
	//BEGIN custom_iteration
	sliceIteration := Iteration[[]int, int](
		//Iteration function
		func(source []int) iter.Seq[int] {
			return func(yield func(focus int) bool) {
				for _, v := range source {
					if !yield(v) {
						break
					}
				}
			}
		},
		func(source []int) int {
			return len(source)
		},
		ExprCustom("sliceIteration"),
	)

	result, found := MustGetFirst(
		sliceIteration,
		[]int{1, 2, 3},
	)
	fmt.Println(result, found)
	//END custom_iteration
}

func TestCustomTraversal(t *testing.T) {
	//BEGIN custom_traversal
	sliceTraversal := Traversal[[]int, int](
		//Iteration function
		func(source []int) iter.Seq[int] {
			return func(yield func(focus int) bool) {
				for _, v := range source {
					if !yield(v) {
						break
					}
				}
			}
		},
		//Length getter
		func(source []int) int {
			return len(source)
		},
		//Modify function
		func(fmap func(focus int) int, source []int) []int {
			var modified []int
			for _, v := range source {
				modified = append(modified, fmap(v))
			}
			return modified
		},
		ExprCustom("sliceTraversal"),
	)

	result := MustModify(
		sliceTraversal,
		Mul(2),
		[]int{1, 2, 3},
	)
	fmt.Println(result)
	//END custom_traversal

	if !reflect.DeepEqual([]any{result}, []any{
		[]int{2, 4, 6},
	},
	) {
		t.Fatal(result)
	}
}

func TestCustomTraversalI(t *testing.T) {
	//BEGIN custom_traversali
	sliceTraversalI := TraversalI[int, []int, int](
		//Iteration function
		func(source []int) SeqI[int, int] {
			return func(yield func(index int, focus int) bool) {
				for i, v := range source {
					if !yield(i, v) {
						break
					}
				}
			}
		},
		//Length getter
		func(source []int) int {
			return len(source)
		},
		//Modify function
		func(fmap func(index int, focus int) int, source []int) []int {
			var modified []int
			for i, v := range source {
				modified = append(modified, fmap(i, v))
			}
			return modified
		},
		//Index getter function
		func(source []int, index int) iter.Seq2[int, int] {
			return func(yield func(index int, focus int) bool) {
				yield(index, source[index])
			}
		},
		//Ix Match
		func(indexA, indexB int) bool {
			return indexA == indexB
		},
		ExprCustom("sliceTraversalI"),
	)

	result, found := MustGetFirst(
		Index(
			sliceTraversalI,
			1,
		),
		[]int{1, 2, 3},
	)
	fmt.Println(result, found)
	//END custom_traversali

	if !reflect.DeepEqual([]any{result, found}, []any{
		2, true,
	},
	) {
		t.Fatal(result)
	}
}

func TestCustomPrism(t *testing.T) {
	//BEGIN custom_prism
	writerToBuffer := Prism(
		//Match function
		func(source io.Writer) (*bytes.Buffer, bool) {
			buf, ok := source.(*bytes.Buffer)
			return buf, ok
		},
		//Embed function
		func(focus *bytes.Buffer) io.Writer {
			return focus
		},
		ExprCustom("writerToBuffer"),
	)

	var w io.Writer = &bytes.Buffer{}

	result := MustModify(writerToBuffer, Op(func(buf *bytes.Buffer) *bytes.Buffer {
		buf.Grow(100)
		fmt.Println("buf.Grow(100)")
		return buf
	}), w)
	//END custom_prism

	if result.(*bytes.Buffer).Cap() != 112 {
		t.Fatal(result.(*bytes.Buffer).Cap())
	}

}
