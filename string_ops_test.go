package optic_test

import (
	"fmt"

	"github.com/samber/lo"

	. "github.com/spearson78/go-optic"
)

func ExampleStringOp() {

	data := lo.T2(1, "Lorem ipsum dolor sit amet")

	optic := StringOp(FilteredCol[int](In('a', 'e', 'i', 'o', 'u', ' ')))

	var getRes string = MustGet(
		Compose(
			T2B[int, string](),
			optic,
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, string] = MustModify(
		T2B[int, string](),
		optic,
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//oe iu oo i ae
	//{1 oe iu oo i ae}
}

func ExampleFilteredString() {

	data := lo.T2(1, "Lorem ipsum dolor sit amet")

	var getRes string = MustGet(
		Compose(
			T2B[int, string](),
			FilteredString(In('a', 'e', 'i', 'o', 'u', ' ')),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, string] = MustModify(
		T2B[int, string](),
		FilteredString(In('a', 'e', 'i', 'o', 'u', ' ')),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//oe iu oo i ae
	//{1 oe iu oo i ae}
}

func ExampleFilteredStringI() {

	data := lo.T2(1, "Lorem ipsum dolor sit amet")

	var getRes string = MustGet(
		Compose(
			T2B[int, string](),
			FilteredStringI(OpOnIx[rune](Odd[int]())),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, string] = MustModify(
		T2B[int, string](),
		FilteredStringI(OpOnIx[rune](Odd[int]())),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//oe pu oo i mt
	//{1 oe pu oo i mt}
}

func ExampleAppendString() {

	data := lo.T2(1, "Lorem ipsum dolor sit")

	var getRes string = MustGet(
		Compose(
			T2B[int, string](),
			AppendString(StringCol(" amet")),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, string] = MustModify(
		T2B[int, string](),
		AppendString(StringCol(" amet")),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//Lorem ipsum dolor sit amet
	//{1 Lorem ipsum dolor sit amet}
}

func ExampleReversedString() {

	data := lo.T2(1, "Lorem ipsum dolor sit amet")

	var getRes string = MustGet(
		Compose(
			T2B[int, string](),
			ReversedString(),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, string] = MustModify(
		T2B[int, string](),
		ReversedString(),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//tema tis rolod muspi meroL
	//{1 tema tis rolod muspi meroL}
}

func ExamplePrependString() {

	data := lo.T2(1, "ipsum dolor sit amet")

	var getRes string = MustGet(
		Compose(
			T2B[int, string](),
			PrependString(StringCol("Lorem ")),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, string] = MustModify(
		T2B[int, string](),
		PrependString(StringCol("Lorem ")),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//Lorem ipsum dolor sit amet
	//{1 Lorem ipsum dolor sit amet}
}

func ExampleOrderedString() {

	data := lo.T2(1, "Lorem ipsum dolor sit amet")

	var getRes string = MustGet(
		Compose(
			T2B[int, string](),
			OrderedString(OrderBy(Identity[rune]())),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, string] = MustModify(
		T2B[int, string](),
		OrderedString(OrderBy(Identity[rune]())),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//Ladeeiilmmmoooprrssttu
	//{1     Ladeeiilmmmoooprrssttu}
}

func ExampleOrderedStringI() {

	data := lo.T2(1, "Lorem ipsum dolor sit amet")

	var getRes string = MustGet(
		Compose(
			T2B[int, string](),
			OrderedStringI(
				Desc(
					OrderByI(
						ValueIIndex[int, rune](),
					),
				),
			),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, string] = MustModify(
		T2B[int, string](),
		OrderedStringI(Desc(OrderByI(ValueIIndex[int, rune]()))),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//tema tis rolod muspi meroL
	//{1 tema tis rolod muspi meroL}
}

func ExampleSubString() {

	data := lo.T2(1, "Lorem ipsum dolor sit amet")

	var getRes string = MustGet(
		Compose(
			T2B[int, string](),
			SubString(1, -1),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, string] = MustModify(
		T2B[int, string](),
		SubString(1, -1),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//orem ipsum dolor sit ame
	//{1 orem ipsum dolor sit ame}
}
