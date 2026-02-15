package optic

import (
	"fmt"
)

type Tree struct {
	Value    string
	Children []*Tree
}

func (t *Tree) String() string {
	return fmt.Sprintf("{%v [%v]}", t.Value, t.Children)
}

func ExamplePtrFieldLensE() {

	childrenField := PtrFieldLensE(func(source *Tree) *[]*Tree { return &source.Children })
	childrenTraversal := Compose(childrenField, TraverseSlice[*Tree]())

	data := &Tree{
		Value: "root",
		Children: []*Tree{
			&Tree{
				Value: "root/first",
			},
			&Tree{
				Value: "root/second",
			},
		},
	}

	var found bool
	var firstChild *Tree
	firstChild, found, err := GetFirst(Index(childrenTraversal, 0), data)
	fmt.Println(firstChild, found, err)

	//Note: a new root node is returned with the new child node
	var addChild *Tree
	addChild, err = Modify(childrenField, AppendSlice(ValCol(&Tree{Value: "root/third"})), data)
	fmt.Println(addChild, err)

	//The original tree has not been modified.
	fmt.Println(data)

	//Output:
	//{root/first [[]]} true <nil>
	//{root [[{root/first [[]]} {root/second [[]]} {root/third [[]]}]]} <nil>
	//{root [[{root/first [[]]} {root/second [[]]}]]}
}
