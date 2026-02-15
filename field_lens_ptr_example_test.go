package optic_test

import (
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

type PtrPerson struct {
	Name    string
	Age     int `example:"age"`
	Hobbies []string
}

func ExamplePtrFieldLens() {

	nameField := PtrFieldLens(func(source *Person) *string { return &source.Name })
	ageField := PtrFieldLens(func(source *Person) *int { return &source.Age })
	hobbiesField := PtrFieldLens(func(source *Person) *[]string { return &source.Hobbies })

	data := &Person{
		Name:    "Max Mustermann",
		Age:     46,
		Hobbies: []string{"eating", "sleeping"},
	}

	name, _ := MustGetFirst(nameField, data)
	age, _ := MustGetFirst(ageField, data)
	hobbies, _ := MustGetFirst(hobbiesField, data)

	fmt.Println(name, age, hobbies)

	olderPerson := MustSet(ageField, 47, data)
	fmt.Println(olderPerson)

	//Note: the return type is a person with hobbies converted to upper case
	var upperHobbies *Person = MustModify(Compose(hobbiesField, TraverseSlice[string]()), Op(strings.ToUpper), data)
	fmt.Println(upperHobbies)

	//Output:
	//Max Mustermann 46 [eating sleeping]
	//&{Max Mustermann 47 [eating sleeping]}
	//&{Max Mustermann 46 [EATING SLEEPING]}
}
