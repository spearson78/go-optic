package otree_test

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/otree"
)

func TestTraverseMyTree(t *testing.T) {

	type mytree struct {
		Value    string
		Children []mytree
	}

	data := mytree{
		Value: "Root",
		Children: []mytree{
			{
				Value: "Child 1",
			},
			{
				Value: "Child 2",
			},
		},
	}

	value := FieldLens(func(source *mytree) *string { return &source.Value })

	children := Compose(
		FieldLens(func(source *mytree) *[]mytree { return &source.Children }),
		TraverseSlice[mytree](),
	)

	topDownVal := Compose(TopDown(children), value)

	log.Println("SliceOf")

	if result := MustGet(SliceOf(topDownVal, 10), data); !reflect.DeepEqual(result, []string{"Root", "Child 1", "Child 2"}) {
		t.Fatal("mytree traverse", result)
	}

	log.Println("Modify")

	if result := MustModify(topDownVal, Op(strings.ToUpper), data); !reflect.DeepEqual(result, mytree{
		Value: "ROOT",
		Children: []mytree{
			{
				Value:    "CHILD 1",
				Children: []mytree{},
			},
			{
				Value:    "CHILD 2",
				Children: []mytree{},
			},
		},
	}) {
		t.Fatal("mytree modify", result)
	}

	log.Println("SliceOf")

	bottomUpVal := Compose(BottomUp(children), value)

	if result := MustGet(SliceOf(bottomUpVal, 10), data); !reflect.DeepEqual(result, []string{"Child 1", "Child 2", "Root"}) {
		t.Fatal("mytree traverse bup", result)
	}

	log.Println("Modify")

	i := 0
	if result := MustModify(bottomUpVal, Op(func(focus string) string { i++; return fmt.Sprintf("%v:%v", i, focus) }), data); !reflect.DeepEqual(result, mytree{
		Value: "3:Root",
		Children: []mytree{
			{
				Value:    "1:Child 1",
				Children: []mytree{},
			},
			{
				Value:    "2:Child 2",
				Children: []mytree{},
			},
		},
	}) {
		t.Fatal("mytree modify bup-", result)
	}

}
