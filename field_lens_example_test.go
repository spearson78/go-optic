package optic_test

import (
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

type Person struct {
	Name    string
	Age     int `example:"age"`
	Hobbies []string
}

func ExampleFieldLens() {

	nameField := FieldLens(func(source *Person) *string { return &source.Name })
	ageField := FieldLens(func(source *Person) *int { return &source.Age })
	hobbiesField := FieldLens(func(source *Person) *[]string { return &source.Hobbies })

	data := Person{
		Name:    "Max Mustermann",
		Age:     46,
		Hobbies: []string{"eating", "sleeping"},
	}

	name := MustGet(nameField, data)
	age := MustGet(ageField, data)
	hobbies := MustGet(hobbiesField, data)

	fmt.Println(name, age, hobbies)

	olderPerson := MustSet(ageField, 47, data)
	fmt.Println(olderPerson)

	//Note: the return type is a person with hobbies converted to upper case
	var upperHobbies Person = MustModify(Compose(hobbiesField, TraverseSlice[string]()), Op(strings.ToUpper), data)
	fmt.Println(upperHobbies)

	//Output:
	//Max Mustermann 46 [eating sleeping]
	//{Max Mustermann 47 [eating sleeping]}
	//{Max Mustermann 46 [EATING SLEEPING]}
}
