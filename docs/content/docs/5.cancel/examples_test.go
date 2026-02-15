package cancel

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/spearson78/go-optic"
)

func TestCancel(t *testing.T) {
	//BEGIN cancel
	ctx, cancel := context.WithTimeout(
		context.Background(),
		100*time.Millisecond,
	)
	defer cancel()

	result, err := GetContext(
		ctx,
		SliceOf(
			TraverseCol[int, int](),
			0,
		),
		Col(
			func(yield func(focus int) bool) {
				i := 0
				for yield(i) {
					i++
				}
			},
			nil, //default length getter
		),
	)

	fmt.Println(result, err)
	//END cancel

	expected := []any{
		//BEGIN cancel_result
		[]int(nil),
		`context deadline exceeded
optic error path:
	Traverse
	SliceOf(Traverse)
`,
		//END cancel_result
	}

	if !strings.HasPrefix(err.Error(), "context deadline exceeded") {
		t.Fatal(expected, err.Error())
	}

}
