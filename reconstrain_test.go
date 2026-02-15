package optic_test

import (
	"fmt"

	"github.com/samber/lo"

	. "github.com/spearson78/go-optic"
)

func ExampleRet1() {
	var optic Optic[
		Void,
		lo.Tuple2[int, string],
		lo.Tuple2[int, string],
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[UniDir, BiDir],
		CompositionTree[Pure, Err],
	] = Compose(
		T2B[int, string](),
		ParseInt[int](10, 32),
	)

	var reconstrained Optic[
		Void,
		lo.Tuple2[int, string],
		lo.Tuple2[int, string],
		int,
		int,
		ReturnOne,
		ReadWrite,
		UniDir,
		Err,
	] = Ret1(Rw(Ud(EErr(optic))))

	fmt.Println(reconstrained.OpticType())

	//Output:
	//Lens
}

func ExampleRetM() {
	var optic Optic[
		Void,
		[]string,
		[]string,
		int,
		int,
		CompositionTree[ReturnMany, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[UniDir, BiDir],
		CompositionTree[Pure, Err],
	] = Compose(
		TraverseSlice[string](),
		ParseInt[int](10, 32),
	)

	var reconstrained Optic[
		Void,
		[]string,
		[]string,
		int,
		int,
		ReturnMany,
		ReadWrite,
		UniDir,
		Err,
	] = RetM(Rw(Ud(EErr(optic))))

	fmt.Println(reconstrained.OpticType())

	//Output:
	//Traversal
}

func ExampleRetL() {
	var opticM Optic[
		Void,
		[]string,
		[]string,
		int,
		int,
		CompositionTree[ReturnMany, ReturnOne], //Note ReturnMany on the left
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[UniDir, BiDir],
		CompositionTree[Pure, Err],
	] = Compose(
		TraverseSlice[string](),
		ParseInt[int](10, 32),
	)

	var reconstrainedM Optic[
		Void,
		[]string,
		[]string,
		int,
		int,
		ReturnMany, //RetL reconstrained as ReturnMany
		ReadWrite,
		UniDir,
		Err,
	] = RetL(Rw(Ud(EErr(opticM))))

	fmt.Println(reconstrainedM.OpticType())

	var optic1 Optic[
		Void,
		[]byte,
		[]byte,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne], //Note the ReturnOne on the left
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[BiDir, BiDir],
		CompositionTree[Pure, Err],
	] = Compose(
		IsoCast[[]byte, string](),
		ParseInt[int](10, 32),
	)

	var reconstrained1 Optic[
		Void,
		[]byte,
		[]byte,
		int,
		int,
		ReturnOne, //RetL reconstrained as ReturnOne
		ReadWrite,
		BiDir,
		Err,
	] = RetL(Rw(Bd(EErr(optic1))))

	fmt.Println(reconstrained1.OpticType())

	//Output:
	//Traversal
	//Iso
}

func ExampleRetR() {
	var opticM Optic[
		int,
		[]byte,
		[]byte,
		rune,
		rune,
		CompositionTree[ReturnOne, ReturnMany], //Note ReturnMany on the right
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[BiDir, UniDir],
		CompositionTree[Pure, Pure],
	] = Compose(
		IsoCast[[]byte, string](),
		TraverseString(),
	)

	var reconstrainedM Optic[
		int,
		[]byte,
		[]byte,
		rune,
		rune,
		ReturnMany, //RetL reconstrained as ReturnMany
		ReadWrite,
		UniDir,
		Err,
	] = RetR(Rw(Ud(EErr(opticM))))

	fmt.Println(reconstrainedM.OpticType())

	var optic1 Optic[
		Void,
		[]byte,
		[]byte,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne], //Note the ReturnOne on the right
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[BiDir, BiDir],
		CompositionTree[Pure, Err],
	] = Compose(
		IsoCast[[]byte, string](),
		ParseInt[int](10, 32),
	)

	var reconstrained1 Optic[
		Void,
		[]byte,
		[]byte,
		int,
		int,
		ReturnOne, //RetR reconstrained as ReturnOne
		ReadWrite,
		BiDir,
		Err,
	] = RetR(Rw(Bd(EErr(optic1))))

	fmt.Println(reconstrained1.OpticType())

	//Output:
	//Traversal
	//Iso
}

func ExampleRw() {
	var optic Optic[
		Void,
		lo.Tuple2[int, string],
		lo.Tuple2[int, string],
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[UniDir, BiDir],
		CompositionTree[Pure, Err],
	] = Compose(
		T2B[int, string](),
		ParseInt[int](10, 32),
	)

	var reconstrained Optic[
		Void,
		lo.Tuple2[int, string],
		lo.Tuple2[int, string],
		int,
		int,
		ReturnOne,
		ReadWrite,
		UniDir,
		Err,
	] = Ret1(Rw(Ud(EErr(optic))))

	fmt.Println(reconstrained.OpticType())

	//Output:
	//Lens
}

func ExampleRo() {
	var optic Optic[
		Void,
		[]int,
		[]int,
		bool,
		bool,
		CompositionTree[ReturnMany, ReturnOne],
		CompositionTree[ReadWrite, ReadOnly],
		CompositionTree[UniDir, UniDir],
		CompositionTree[Pure, Pure],
	] = Compose(
		TraverseSlice[int](),
		Eq(10),
	)

	var reconstrained Optic[
		Void,
		[]int,
		[]int,
		bool,
		bool,
		ReturnMany,
		ReadOnly,
		UniDir,
		Pure,
	] = RetM(Ro(Ud(EPure(optic))))

	fmt.Println(reconstrained.OpticType())

	//Output:
	//Iteration
}

func ExampleRwL() {
	var opticRo Optic[
		Void,
		int,
		int,
		bool,
		bool,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadOnly, ReadWrite], //Note the ReadOnly on the left
		CompositionTree[UniDir, BiDir],
		CompositionTree[Pure, Pure],
	] = Compose(
		Eq(10),
		Not(),
	)

	var reconstrainedRo Optic[
		Void,
		int,
		int,
		bool,
		bool,
		ReturnOne,
		ReadOnly, //RwL reconstrained to ReadOnly
		UniDir,
		Pure,
	] = Ret1(RwL(Ud(EPure(opticRo))))

	fmt.Println(reconstrainedRo.OpticType())

	var opticRw Optic[
		Void,
		int,
		int,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite], //Note the ReadWrite on the left
		CompositionTree[BiDir, BiDir],
		CompositionTree[Pure, Pure],
	] = Compose(
		Add(10),
		Mul(2),
	)

	var reconstrainedRw Optic[
		Void,
		int,
		int,
		int,
		int,
		ReturnOne,
		ReadWrite, //RwL reconstrained to ReadWrite
		BiDir,
		Pure,
	] = Ret1(RwL(Bd(EPure(opticRw))))

	fmt.Println(reconstrainedRw.OpticType())

	//Output:
	//Getter
	//Iso
}

func ExampleRwR() {
	var opticRo Optic[
		Void,
		int,
		int,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadOnly], //Note the ReadOnly on the right
		CompositionTree[BiDir, UniDir],
		CompositionTree[Pure, Pure],
	] = Compose(
		Mul(10),
		Mod(3),
	)

	var reconstrainedRo Optic[
		Void,
		int,
		int,
		int,
		int,
		ReturnOne,
		ReadOnly, //RwR reconstrained to ReadOnly
		UniDir,
		Pure,
	] = Ret1(RwR(Ud(EPure(opticRo))))

	fmt.Println(reconstrainedRo.OpticType())

	var opticRw Optic[
		Void,
		int,
		int,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite], //Note the ReadWrite on the right
		CompositionTree[BiDir, BiDir],
		CompositionTree[Pure, Pure],
	] = Compose(
		Add(10),
		Mul(2),
	)

	var reconstrainedRw Optic[
		Void,
		int,
		int,
		int,
		int,
		ReturnOne,
		ReadWrite, //RwL reconstrained to ReadWrite
		BiDir,
		Pure,
	] = Ret1(RwR(Bd(EPure(opticRw))))

	fmt.Println(reconstrainedRw.OpticType())

	//Output:
	//Getter
	//Iso
}

func ExampleBd() {
	var optic Optic[
		Void,
		int,
		int,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[BiDir, BiDir],
		CompositionTree[Pure, Pure],
	] = Compose(
		Add(10),
		Mul(2),
	)

	var reconstrained Optic[
		Void,
		int,
		int,
		int,
		int,
		ReturnOne,
		ReadWrite,
		BiDir,
		Pure,
	] = Ret1(Rw(Bd(EPure(optic))))

	fmt.Println(reconstrained.OpticType())

	//Output:
	//Iso
}

func ExampleUd() {
	var optic Optic[
		Void,
		lo.Tuple2[int, string],
		lo.Tuple2[int, string],
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[UniDir, BiDir],
		CompositionTree[Pure, Err],
	] = Compose(
		T2B[int, string](),
		ParseInt[int](10, 32),
	)

	var reconstrained Optic[
		Void,
		lo.Tuple2[int, string],
		lo.Tuple2[int, string],
		int,
		int,
		ReturnOne,
		ReadWrite,
		UniDir,
		Err,
	] = Ret1(Rw(Ud(EErr(optic))))

	fmt.Println(reconstrained.OpticType())

	//Output:
	//Lens
}

func ExampleDirL() {
	var opticBd Optic[
		Void,
		int,
		int,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[BiDir, BiDir], //Note the BiDir on the left
		CompositionTree[Pure, Pure],
	] = Compose(
		Mul(10),
		Add(3),
	)

	var reconstrainedBd Optic[
		Void,
		int,
		int,
		int,
		int,
		ReturnOne,
		ReadWrite,
		BiDir, //DirL reconstrained to BiDir
		Pure,
	] = Ret1(Rw(DirL(EPure(opticBd))))

	fmt.Println(reconstrainedBd.OpticType())

	var opticUd Optic[
		Void,
		int,
		int,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadOnly, ReadWrite],
		CompositionTree[UniDir, BiDir], //Note the UniDir on the right
		CompositionTree[Pure, Pure],
	] = Compose(
		Mod(10),
		Mul(2),
	)

	var reconstrainedUd Optic[
		Void,
		int,
		int,
		int,
		int,
		ReturnOne,
		ReadOnly,
		UniDir, //DirL reconstrained to BiDir
		Pure,
	] = Ret1(Ro(DirL(EPure(opticUd))))

	fmt.Println(reconstrainedUd.OpticType())

	//Output:
	//Iso
	//Getter
}

func ExampleDirR() {
	var opticUd Optic[
		Void,
		int,
		int,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadOnly],
		CompositionTree[BiDir, UniDir], //Note the ReadOnly on the right
		CompositionTree[Pure, Pure],
	] = Compose(
		Mul(10),
		Mod(3),
	)

	var reconstrainedUd Optic[
		Void,
		int,
		int,
		int,
		int,
		ReturnOne,
		ReadOnly,
		UniDir, //BdR reconstrained to UniDir
		Pure,
	] = Ret1(Ro(DirR(EPure(opticUd))))

	fmt.Println(reconstrainedUd.OpticType())

	var opticBd Optic[
		Void,
		int,
		int,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[BiDir, BiDir], //Note the BiDir on the right
		CompositionTree[Pure, Pure],
	] = Compose(
		Add(10),
		Mul(2),
	)

	var reconstrainedBd Optic[
		Void,
		int,
		int,
		int,
		int,
		ReturnOne,
		ReadWrite,
		BiDir, //BdR reconstrained to BiDir
		Pure,
	] = Ret1(Rw(DirR(EPure(opticBd))))

	fmt.Println(reconstrainedBd.OpticType())

	//Output:
	//Getter
	//Iso
}

func ExampleEPure() {
	var opticPure Optic[
		Void,
		int,
		int,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[BiDir, BiDir],
		CompositionTree[Pure, Pure], //Note the path is Pure.
	] = Compose(
		Add(10),
		Mul(2),
	)

	var reconstrainedPure Optic[
		Void,
		int,
		int,
		int,
		int,
		ReturnOne,
		ReadWrite,
		BiDir,
		Pure, //EPure reconstrained to Pure
	] = Ret1(Rw(Bd(EPure(opticPure))))

	fmt.Println(reconstrainedPure.OpticType())

	//Output:
	//Iso
}

func ExampleEErr() {
	var opticPure Optic[
		Void,
		string,
		string,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[BiDir, BiDir],
		CompositionTree[Err, Pure], //Note the path is not Pure.
	] = Compose(
		ParseInt[int](10, 0),
		Mul(2),
	)

	var reconstrainedPure Optic[
		Void,
		string,
		string,
		int,
		int,
		ReturnOne,
		ReadWrite,
		BiDir,
		Err, //EErr reconstrained to Err
	] = Ret1(Rw(Bd(EErr(opticPure))))

	fmt.Println(reconstrainedPure.OpticType())

	//Output:
	//Iso
}

func ExampleEErrL() {
	var opticPure Optic[
		Void,
		string,
		string,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[BiDir, BiDir],
		CompositionTree[Err, Pure], //Note the path is not Pure.
	] = Compose(
		ParseInt[int](10, 0),
		Mul(2),
	)

	var reconstrainedPure Optic[
		Void,
		string,
		string,
		int,
		int,
		ReturnOne,
		ReadWrite,
		BiDir,
		Err, //EErrL reconstrained to Err
	] = Ret1(Rw(Bd(EErrL(opticPure))))

	fmt.Println(reconstrainedPure.OpticType())

	//Output:
	//Iso
}

func ExampleEErrR() {
	var opticPure Optic[
		Void,
		int,
		int,
		string,
		string,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[BiDir, BiDir],
		CompositionTree[Pure, Err], //Note the path is not Pure.
	] = ComposeLeft(
		Mul(2),
		AsReverseGet(ParseInt[int](10, 0)),
	)

	var reconstrainedPure Optic[
		Void,
		int,
		int,
		string,
		string,
		ReturnOne,
		ReadWrite,
		BiDir,
		Err, //EErrR reconstrained to Err
	] = Ret1(Rw(Bd(EErrR(opticPure))))

	fmt.Println(reconstrainedPure.OpticType())

	//Output:
	//Iso
}

func ExampleEErrSwap() {
	var opticPure Optic[
		Void,
		int,
		int,
		string,
		string,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[BiDir, BiDir],
		CompositionTree[Pure, Err],
	] = Compose(
		Mul(2),
		AsReverseGet(ParseInt[int](10, 0)),
	)

	var reconstrainedPure Optic[
		Void,
		int,
		int,
		string,
		string,
		ReturnOne,
		ReadWrite,
		BiDir,
		CompositionTree[Err, Pure], //Pure and Err have been swapped.
	] = Ret1(Rw(Bd(EErrSwap(opticPure))))

	fmt.Println(reconstrainedPure.OpticType())

	//Output:
	//Iso
}

func ExampleEErrSwapL() {
	var opticPure Optic[
		Void,
		int,
		int,
		string,
		string,
		CompositionTree[CompositionTree[ReturnOne, ReturnOne], ReturnOne],
		CompositionTree[CompositionTree[ReadWrite, ReadWrite], ReadOnly],
		CompositionTree[CompositionTree[BiDir, BiDir], UniDir],
		CompositionTree[CompositionTree[Pure, Err], Pure],
	] = Compose(
		Compose(
			Mul(2),
			AsReverseGet(ParseInt[int](10, 0)),
		),
		PrependString(StringCol("X")),
	)

	var reconstrainedPure Optic[
		Void,
		int,
		int,
		string,
		string,
		ReturnOne,
		ReadOnly,
		UniDir,
		CompositionTree[CompositionTree[Err, Pure], Pure], //Pure and Err have been swapped.
	] = Ret1(Ro(Ud(EErrSwapL(opticPure))))

	fmt.Println(reconstrainedPure.OpticType())

	//Output:
	//Getter
}

func ExampleEErrMerge() {
	var opticPure Optic[
		Void,
		int,
		int,
		int,
		int,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[BiDir, BiDir],
		CompositionTree[Pure, Pure], //Note the 2 components are identical.
	] = Compose(
		Add(10),
		Mul(2),
	)

	var reconstrainedPure Optic[
		Void,
		int,
		int,
		int,
		int,
		ReturnOne,
		ReadWrite,
		BiDir,
		Pure, //EErrMerge reconstrained to Pure
	] = Ret1(Rw(Bd(EErrMerge(opticPure))))

	fmt.Println(reconstrainedPure.OpticType())

	//Output:
	//Iso
}

func ExampleEErrMergeL() {
	var opticPure Optic[
		Void,
		int,
		int,
		int,
		int,
		CompositionTree[CompositionTree[ReturnOne, ReturnOne], ReturnOne],
		CompositionTree[CompositionTree[ReadWrite, ReadWrite], ReadWrite],
		CompositionTree[CompositionTree[BiDir, BiDir], BiDir],
		CompositionTree[CompositionTree[Pure, Pure], Pure], //Note the left 2 components are identical.
	] = Compose(
		Compose(
			Add(10),
			Mul(2),
		),
		Sub(1),
	)

	var reconstrainedPure Optic[
		Void,
		int,
		int,
		int,
		int,
		ReturnOne,
		ReadWrite,
		BiDir,
		CompositionTree[Pure, Pure], //EErrMergeL reconstrained the left hand side.
	] = Ret1(Rw(Bd(EErrMergeL(opticPure))))

	fmt.Println(reconstrainedPure.OpticType())

	//Output:
	//Iso
}

func ExampleEErrTrans() {
	var opticPure Optic[
		Void,
		int,
		int,
		int,
		int,
		CompositionTree[CompositionTree[ReturnOne, ReturnOne], CompositionTree[ReturnOne, ReturnOne]],
		CompositionTree[CompositionTree[ReadWrite, ReadWrite], CompositionTree[ReadWrite, ReadWrite]],
		CompositionTree[CompositionTree[BiDir, BiDir], CompositionTree[BiDir, BiDir]],
		CompositionTree[CompositionTree[Pure, Err], CompositionTree[Pure, Err]],
	] = Compose4(
		Add(10),
		EErr(Mul(2)),
		Add(10),
		EErr(Mul(2)),
	)

	var reconstrainedPure Optic[
		Void,
		int,
		int,
		int,
		int,
		ReturnOne,
		ReadWrite,
		BiDir,
		CompositionTree[CompositionTree[Pure, Pure], CompositionTree[Err, Err]],
	] = Ret1(Rw(Bd(EErrTrans(opticPure))))

	fmt.Println(reconstrainedPure.OpticType())

	//Output:
	//Iso
}

func ExampleEErrTransL() {
	var opticPure Optic[
		Void,
		int,
		int,
		int,
		int,
		CompositionTree[CompositionTree[ReturnOne, ReturnOne], ReturnOne],
		CompositionTree[CompositionTree[ReadWrite, ReadWrite], ReadWrite],
		CompositionTree[CompositionTree[BiDir, BiDir], BiDir],
		CompositionTree[CompositionTree[Pure, Err], Pure],
	] = Compose(
		Compose(
			Add(10),
			EErr(Mul(2)),
		),
		Sub(1),
	)

	var reconstrainedPure Optic[
		Void,
		int,
		int,
		int,
		int,
		ReturnOne,
		ReadWrite,
		BiDir,
		CompositionTree[CompositionTree[Pure, Pure], Err],
	] = Ret1(Rw(Bd(EErrTransL(opticPure))))

	fmt.Println(reconstrainedPure.OpticType())

	//Output:
	//Iso
}
