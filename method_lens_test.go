package optic_test

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

func TestMethodGetterExpr(t *testing.T) {

	methodLens := MethodGetter(fs.FileInfo.Name)

	expression := methodLens.AsExpr()

	if expression.(expr.MethodGetter).MethodName != "io/fs.FileInfo.Name" {
		t.Fatal(expression)
	}
}

func ExampleMethodGetter() {

	data, _ := os.Stat(".")

	optic := MethodGetter(fs.FileInfo.Name)

	name := MustGet(optic, data)
	fmt.Println(name)

	//Output:
	//.

}
