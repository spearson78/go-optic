package benchmark

import (
	"log"
	"testing"

	"github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic"
)

func BenchmarkSlice(b *testing.B) {

	const sliceLen = 1000000

	genData := func() []int {
		data := make([]int, sliceLen)
		for i := 0; i < len(data); i++ {
			data[i] = i
		}
		return data
	}

	testResult := func(b *testing.B, data []int) {

		if len(data) != sliceLen {
			b.Fatal("len(data) != sliceLen")
		}

		for i, v := range data {
			if v != i*2 {
				b.Fatalf("v != i*2 : v=%v , i=%v", v, i)
			}
		}

	}

	testNotMutated := func(b *testing.B, data []int) {

		if len(data) != sliceLen {
			b.Fatal("len(data) != sliceLen")
		}

		for i, v := range data {
			if v != i {
				b.Fatalf("v != i : v=%v , i=%v", v, i)
			}
		}

	}

	b.Run("go mutable", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			data := genData()
			b.StartTimer()

			for i, v := range data {
				data[i] = v * 2
			}

			b.StopTimer()
			testResult(b, data)
			b.StartTimer()
		}
	})

	b.Run("go immutable", func(b *testing.B) {
		data := genData()
		b.ResetTimer()

		var result []int
		for n := 0; n < b.N; n++ {
			result = nil
			for _, v := range data {
				result = append(result, v*2)
			}
		}

		b.StopTimer()

		testNotMutated(b, data)
		testResult(b, result)
	})

	b.Run("go immutable pre", func(b *testing.B) {
		data := genData()
		b.ResetTimer()

		var result []int
		for n := 0; n < b.N; n++ {
			result = make([]int, 0, sliceLen)
			for _, v := range data {
				result = append(result, v*2)
			}
		}

		b.StopTimer()

		testNotMutated(b, data)
		testResult(b, result)
	})

	b.Run("optic", func(b *testing.B) {
		data := genData()
		b.ResetTimer()

		var result []int
		for n := 0; n < b.N; n++ {
			result = MustModify(
				optic.TraverseSlice[int](),
				Mul(2),
				data,
			)
		}

		b.StopTimer()

		testNotMutated(b, data)
		testResult(b, result)
	})

	b.Run("optic execonly", func(b *testing.B) {
		data := genData()
		o := TraverseSlice[int]()
		f := Mul(2)
		b.ResetTimer()

		var result []int
		for n := 0; n < b.N; n++ {
			result = MustModify(
				o,
				f,
				data,
			)
		}

		b.StopTimer()

		testNotMutated(b, data)
		testResult(b, result)
	})

}

func TestOverSliceAllocs(t *testing.T) {

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	o := optic.TraverseSlice[int]()
	f := OpToOpI[int](Mul(2))

	for i := 0; i < 100; i++ {
		result := MustModifyI(
			o,
			f,
			data,
		)

		log.Println(result)
	}

}

func TestIterateSliceAllocs(t *testing.T) {

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	o := optic.TraverseSlice[int]()

	result := MustGet(
		SliceOf(
			o,
			len(data),
		),
		data,
	)

	log.Println(result)

}

func TestIterateNestedSliceAllocs(t *testing.T) {

	data := [][]int{
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	o := Compose(optic.TraverseSlice[[]int](), optic.TraverseSlice[int]())

	result := MustGet(
		SliceOf(
			o,
			len(data)*len(data[0]),
		),
		data,
	)

	log.Println(result)

}
