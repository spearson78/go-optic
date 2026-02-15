package oio_test

import (
	"fmt"

	"github.com/gdamore/encoding"
	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/exp/oio"
)

func ExampleDecodeFile() {

	data := "Lorem\nIpsum"

	optic := Compose3(
		IsoCast[string, []byte](),
		AsReverseGet(Bytes()),
		DecodeFile(encoding.UTF8),
	)

	res, err := Get(
		optic,
		data,
	)

	fmt.Println(res, err)

	modRes, err := Modify(
		optic,
		ColSourceFocusErr(FilteredCol[LinePosition](Ne('m'))),
		data,
	)

	fmt.Println(modRes, err)

	//Output:
	//Col[{0 0}:76 {0 1}:111 {0 2}:114 {0 3}:101 {0 4}:109 {0 5}:10 {1 0}:73 {1 1}:112 {1 2}:115 {1 3}:117 {1 4}:109] <nil>
	//Lore
	//Ipsu <nil>

}
