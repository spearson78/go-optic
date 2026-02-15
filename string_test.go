package optic_test

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/samber/lo"
	u "golang.org/x/text/encoding/unicode"

	. "github.com/spearson78/go-optic"
)

func FuzzWorded(f *testing.F) {
	f.Add("Hello World", "lorem")
	f.Add("", "lorem")
	f.Add(" Hello World ", "lorem")

	f.Fuzz(func(t *testing.T, source string, newVal string) {
		ValidateOpticTest(t, Worded(), source, newVal)
	})
}

func ExampleTraverseString() {

	element2, _ := MustGetFirst(Index(TraverseString(), 1), "example")
	fmt.Println(string(element2))

	//unicode.ToUpper operates on individual runes.
	result := MustModify(TraverseString(), Op(unicode.ToUpper), "example")
	fmt.Println(result)

	//Output:x
	//EXAMPLE
}

func ExampleMatchString() {

	data := "Lorem 10 ipsum 20 dolor"

	var sliceResult []string = MustGet(
		SliceOf(
			MatchString(
				regexp.MustCompile(`\d+`),
				-1,
			),
			10,
		),
		data,
	)
	fmt.Println(sliceResult)

	var overResult string
	overResult, err := Modify(
		Compose(
			MatchString(
				regexp.MustCompile(`\d+`),
				-1,
			),
			ParseInt[int](10, 32)),
		Mul(2),
		data,
	)
	fmt.Println(overResult, err)

	//Output:
	//[10 20]
	//Lorem 20 ipsum 40 dolor <nil>
}

func ExampleCaptureString() {

	data := "25/03/2025"

	captures, ok := MustGetFirst(
		CaptureString(
			regexp.MustCompile(`(\d+)/(\d+)/(\d+)`),
			-1,
		),
		data,
	)
	fmt.Println(captures, ok)

	var modifyResult string
	modifyResult, err := Modify(
		Compose3(
			CaptureString(
				regexp.MustCompile(`(\d+)/(\d+)/(\d+)`),
				-1,
			),
			TraverseSlice[string](),
			ParseInt[int](10, 0),
		),
		Add(1),
		data,
	)
	fmt.Println(modifyResult, err)

	//Output:
	//[25 03 2025] true
	//26/4/2026 <nil>
}

func ExampleSplitString() {

	data := "Lorem  Ipsum    Dolor    Sit    Amet"

	splitWhitespace := SplitString(regexp.MustCompile(`\s+`))

	numWords := MustGet(Length(splitWhitespace), data)
	fmt.Println(numWords)

	result := MustModifyI(splitWhitespace, OpI(func(index int, word string) string {
		if index%2 == 0 {
			return strings.ToUpper(word) //Convert even indexed words to uppercase
		} else {
			return strings.ToLower(word) //Convert odd indexed words to lowercase
		}

	}), data)
	fmt.Println(result)

	//Output:5
	//LOREM  ipsum    DOLOR    sit    AMET
}

func ExampleWorded() {

	numWords := MustGet(Length(Worded()), "Lorem  Ipsum    Dolor    Sit    Amet")
	fmt.Println(numWords)

	//Note that all white space is converted to a single space rune in the result
	result := MustModifyI(Worded(), OpI(func(index int, word string) string {
		if index%2 == 0 {
			return strings.ToUpper(word) //Convert even indexed words to uppercase
		} else {
			return strings.ToLower(word) //Convert odd indexed words to lowercase
		}

	}), "Lorem  Ipsum    Dolor    Sit    Amet")
	fmt.Println(result)

	//Output:5
	//LOREM  ipsum    DOLOR    sit    AMET
}

func ExampleLined() {

	numLines := MustGet(Length(Lined()), "Lorem Ipsum\n\nDolor\nSit Amet")
	fmt.Println(numLines)

	result := MustModifyI(Lined(), OpI(func(index int, word string) string {
		if index%2 == 0 {
			return strings.ToUpper(word) //Convert even indexed lines to uppercase
		} else {
			return strings.ToLower(word) //Convert odd indexed lines to lowercase
		}

	}), "Lorem Ipsum\n\nDolor\nSit Amet")
	fmt.Println(result)

	//Output:3
	//LOREM IPSUM
	//
	//dolor
	//SIT AMET
}

func ExampleStringToCol() {

	data := "Lorem ipsum dolor sit amet"

	//See: [FilteredString] for a more convenient string filter function.
	var getRes string = MustGet(
		Compose3(
			StringToCol(),
			FilteredCol[int](In('a', 'e', 'i', 'o', 'u', ' ')),
			AsReverseGet(StringToCol()),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes string = MustModify(
		StringToCol(),
		FilteredCol[int](In('a', 'e', 'i', 'o', 'u', ' ')),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//oe iu oo i ae
	//oe iu oo i ae
}

func ExampleEncodeString() {

	data := "\u03bc"

	res, err := Get(
		Length(
			Compose(
				EncodeString(u.UTF8), //import u "golang.org/x/text/encoding/unicode"
				TraverseSlice[byte](),
			),
		),
		data,
	)
	fmt.Println(res, err)

	//Output:
	//2 <nil>
}

func ExampleStringHasPrefix() {

	data := "alphabet"

	res := MustGet(
		StringHasPrefix("alpha"),
		data,
	)
	fmt.Println(res)

	//Output:
	//true
}

func ExampleStringHasSuffix() {

	data := "alphabet"

	res := MustGet(
		StringHasSuffix("bet"),
		data,
	)
	fmt.Println(res)

	//Output:
	//true
}

func ExampleStringBuilderReducer() {

	data := []string{"alpha", "bet"}

	res, ok := MustGetFirst(
		Reduce(
			TraverseSlice[string](),
			StringBuilderReducer(),
		),
		data,
	)
	fmt.Println(res, ok)

	//Output:
	//alphabet true
}

func ExampleStringColType() {

	data := lo.T2(1, "Lorem ipsum dolor sit amet")

	//See: [FilteredString] for a more convenient string filter function.
	var getRes string = MustGet(
		Compose(
			T2B[int, string](),
			ColTypeOp(StringColType(), FilteredCol[int](In('a', 'e', 'i', 'o', 'u', ' '))),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, string] = MustModify(
		T2B[int, string](),
		ColTypeOp(StringColType(), FilteredCol[int](In('a', 'e', 'i', 'o', 'u', ' '))),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//oe iu oo i ae
	//{1 oe iu oo i ae}
}
func ExampleStringOf() {

	data := map[int]rune{
		1: 'a',
		2: 'l',
		3: 'p',
		4: 'h',
		5: 'a',
	}

	optic := StringOf(TraverseMap[int, rune](), 10)

	var viewResult string = MustGet(optic, data)
	fmt.Println(viewResult)

	var modifyResult map[int]rune = MustModify(optic, Op(strings.ToUpper), data)
	fmt.Println(MustModify(TraverseMapP[int, rune, string](), UpCast[rune, string](), modifyResult))

	//Output:
	//alpha
	//map[1:A 2:L 3:P 4:H 5:A]
}

func TestCaptureString(t *testing.T) {
	// Off-by-one error case: index 2 (third group)
	matchOn := regexp.MustCompile("a(b*)")
	n := 2
	captureStringOptic := SliceOf(CaptureString(matchOn, n), 0)
	result := MustGet(captureStringOptic, "abac")
	if !reflect.DeepEqual(
		result,
		[][]string{
			[]string{
				"b",
			},
			[]string{
				"",
			},
		},
	) {
		t.Errorf("Incorrect result for match off-by-one case 1: %v", result)
	}

	// Off-by-one error case: index 0 (first group)
	matchOn = regexp.MustCompile("(b*)a")
	n = 0
	captureStringOptic = SliceOf(CaptureString(matchOn, n), 0)
	result = MustGet(captureStringOptic, "ba")
	if !reflect.DeepEqual(
		result,
		[][]string{},
	) {
		t.Errorf("Incorrect result for match off-by-one case 2: %v", result)
	}

	// Index out of bounds
	matchOn = regexp.MustCompile("(b*)a")
	n = 3 // There are only two groups (0 and 1), so n is out of bounds.
	captureStringOptic = SliceOf(CaptureString(matchOn, n), 0)
	result = MustGet(captureStringOptic, "ba")
	if !reflect.DeepEqual(
		result,
		[][]string{
			[]string{"b"},
		},
	) {
		t.Errorf("Incorrect result for index out of bounds case 3: %v", result)
	}

	// Empty input
	matchOn = regexp.MustCompile("(b*)a")
	n = 0
	captureStringOptic = SliceOf(CaptureString(matchOn, n), 0)
	result = MustGet(captureStringOptic, "")
	if !reflect.DeepEqual(
		result,
		[][]string{},
	) {
		t.Errorf("Incorrect result for empty input case 4: %v", result)
	}

	// No matches
	matchOn = regexp.MustCompile("(b*)a")
	n = 0
	captureStringOptic = SliceOf(CaptureString(matchOn, n), 0)
	result = MustGet(captureStringOptic, "abc")
	if !reflect.DeepEqual(
		result,
		[][]string{},
	) {
		t.Errorf("Incorrect result for no matches case 5: %v", result)
	}

	// MatchAll with empty string
	matchOn = regexp.MustCompile("(b*)a")
	n = -1
	captureStringOptic = SliceOf(CaptureString(matchOn, n), 0)
	result = MustGet(captureStringOptic, "")
	if !reflect.DeepEqual(
		result,
		[][]string{},
	) {
		t.Errorf("Incorrect result for MatchAll with empty string 6: %v", result)
	}

	// MatchAll with a single match
	matchOn = regexp.MustCompile("(b*)a")
	n = -1
	captureStringOptic = SliceOf(CaptureString(matchOn, n), 0)
	result = MustGet(captureStringOptic, "ba")
	if !reflect.DeepEqual(
		result,
		[][]string{
			[]string{"b"},
		},
	) {
		t.Errorf("Incorrect result for MatchAll with a single match 7: %v", result)
	}

	// MatchAll with multiple matches
	matchOn = regexp.MustCompile("(b*)a")
	n = -1
	captureStringOptic = SliceOf(CaptureString(matchOn, n), 0)
	result = MustGet(captureStringOptic, "baaab")
	if !reflect.DeepEqual(
		result,
		[][]string{
			[]string{"b"},
			[]string{""},
			[]string{""},
		},
	) {
		t.Errorf("Incorrect result for MatchAll with multiple matches 8: %v", result)
	}

}
