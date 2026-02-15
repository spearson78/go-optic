package yield

import (
	"fmt"
	"testing"

	. "github.com/spearson78/go-optic"
)

func TestYieldAfterBreak(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				if err.Error() == "yield called after break" {
					//Ignore and let the error be returned
					return
				}
			}
			panic(r)
		}
	}()

	//BEGIN yieldAfterBreak
	result, err := Get(
		SliceOf(
			Taking(
				TraverseCol[int, int](),
				1,
			),
			0,
		),
		Col(
			func(yield func(focus int) bool) {
				i := 0
				for {
					//false return from yield is ignored here.
					yield(i)
					i++
				}
			},
			nil, //default Length getter
		),
	)

	fmt.Println(result, err)
	//END yieldAfterBreak

	t.Fatal("yieldAfterBreak panics expected")
}
