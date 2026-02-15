package main

import (
	"bufio"
	"bytes"
	"context"
	"testing"
	"time"

	goscript "github.com/spearson78/go-script"
)

var scriptTests = map[string]struct {
	Script string
	Result string
}{
	"ErrorAware": {
		Script: `
package main

import (
	"context"
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

func main() {
	result, err := Modify(
		ParseInt[int](10, 0),
		Mul(2),
		"1",
	)
	fmt.Println(result, err)
}            	
`,
		Result: "2 <nil>\n",
	},
	"TestCustomFieldLens": {
		Script: `
package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
)
type ExampleStruct struct {
	name    string
	address string
}

func main() {
	customFieldLens := Lens(
		//Lens getter
		func(source ExampleStruct) string {
			return source.name
		},
		//Lens setter
		func(newValue string, source ExampleStruct) ExampleStruct {
			source.name = newValue
			return source
		},
		ExprCustom("customFieldLens"),
	)

	data := ExampleStruct{
		name:    "Max Mustermann",
		address: "Musterstadt",
	}

	result := MustSet(customFieldLens, "Erika Mustermann", data)
	fmt.Println(result)
}            	
	`,
		Result: "{Musterstadt Erika Mustermann}\n",
	},
	"TestCustomIso": {
		Script: `
package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
)

func main() {
	celsiusToFahrenheit := Iso(
		//Iso getter
		func(celsius float64) float64 {
			return (celsius * 1.8) + 32
		},
		//Iso reverse getter
		func(fahrenheit float64) float64 {
			return (fahrenheit - 32) / 1.8
		},
		ExprCustom("celsiusToFahrenheit"),
	)

	fahrenHeit := MustGet(celsiusToFahrenheit, 32)
	fmt.Println(fahrenHeit)

	celsius := MustReverseGet(celsiusToFahrenheit, 89.6)
	fmt.Println(celsius)
}
            
`,
		Result: "89.6\n31.999999999999996\n",
	},
	"TestCustomIteration": {
		Script: `
package main

import (
	"fmt"
	"iter"

	. "github.com/spearson78/go-optic"
)

func main() {
	sliceIteration := Iteration[[]int, int](
		//Iteration function
		func(source []int) iter.Seq[int] {
			return func(yield func(focus int) bool) {
				for _, v := range source {
					if !yield(v) {
						break
					}
				}
			}
		},
		func(source []int) int {
			return len(source)
		},
		ExprCustom("sliceIteration"),
	)

	result, found := MustGetFirst(
		sliceIteration,
		[]int{1, 2, 3},
	)
	fmt.Println(result, found)
}
            
`,
		Result: "1 true\n",
	},
	"TestCustomTraversal": {
		Script: `
package main

import (
	"fmt"
	"iter"

	. "github.com/spearson78/go-optic"
)

func main() {
	sliceTraversal := Traversal[[]int, int](
		//Iteration function
		func(source []int) iter.Seq[int] {
			return func(yield func(focus int) bool) {
				for _, v := range source {
					if !yield(v) {
						break
					}
				}
			}
		},
		//Length getter
		func(source []int) int {
			return len(source)
		},
		//Modify function
		func(fmap func(focus int) int, source []int) []int {
			var modified []int
			for _, v := range source {
				modified = append(modified, fmap(v))
			}
			return modified
		},
		ExprCustom("sliceTraversal"),
	)

	result := MustModify(
		sliceTraversal,
		Mul(2),
		[]int{1, 2, 3},
	)
	fmt.Println(result)
}
            
`,
		Result: "[2 4 6]\n",
	},
	"TestCustomTraversalI": {
		Script: `
package main

import (
	"fmt"
	"iter"

	. "github.com/spearson78/go-optic"
)

func main() {
	sliceTraversalI := TraversalI[int, []int, int](
		//Iteration function
		func(source []int) SeqI[int, int] {
			return func(yield func(index int, focus int) bool) {
				for i, v := range source {
					if !yield(i, v) {
						break
					}
				}
			}
		},
		//Length getter
		func(source []int) int {
			return len(source)
		},
		//Modify function
		func(fmap func(index int, focus int) int, source []int) []int {
			var modified []int
			for i, v := range source {
				modified = append(modified, fmap(i, v))
			}
			return modified
		},
		//Index getter function
		func(source []int, index int) iter.Seq2[int, int] {
			return func(yield func(index int, focus int) bool) {
				yield(index, source[index])
			}
		},
		//Ix Match
		func(indexA, indexB int) bool {
			return indexA == indexB
		},
		ExprCustom("sliceTraversalI"),
	)

	result, found := MustGetFirst(
		Index(
			sliceTraversalI,
			1,
		),
		[]int{1, 2, 3},
	)
	fmt.Println(result, found)
}
            
`,
		Result: "2 true\n",
	},
	"TestCustomPrism": {
		Script: `
package main

import (
	"bytes"
	"fmt"
	"io"
	"iter"

	. "github.com/spearson78/go-optic"
)

func main() {
	writerToBuffer := Prism(
		//Match function
		func(source io.Writer) (*bytes.Buffer, bool) {
			buf, ok := source.(*bytes.Buffer)
			return buf, ok
		},
		//Embed function
		func(focus *bytes.Buffer) io.Writer {
			return focus
		},
		ExprCustom("writerToBuffer"),
	)

	var w io.Writer = &bytes.Buffer{}

	result := MustModify(writerToBuffer, Op(func(buf *bytes.Buffer) *bytes.Buffer {
		buf.Grow(100)
		fmt.Println("buf.Grow(100)")
		return buf
	}), w)
}
            
            
`,
		Result: "buf.Grow(100)\n",
	},
	"TestCallOp": {
		Script: `
package main

import (
	"bytes"
	"fmt"
	"io"
	"iter"

	. "github.com/spearson78/go-optic"
)

func main() {
	sliceLen := Op(
		func(source []string) int {
			return len(source)
		},
	)

	result := MustGet(sliceLen, []string{"alpha", "beta"})
	fmt.Println(result)
}`,
		Result: "2\n",
	},
	"TestCallBufGrow": {
		Script: `
package main

import (
	"bytes"
	"fmt"
	"io"
	"iter"

	. "github.com/spearson78/go-optic"
)

func main() {
	w := &bytes.Buffer{}

	c := Op(func(source io.Writer) *bytes.Buffer {
		return source.(*bytes.Buffer)
	})

	buf := MustGet(c,w)
	buf.Grow(100)

	fmt.Println(buf.Cap())
}
            
            
`,
		Result: "112\n",
	},
	"TestAppendSpread": {
		Script: `
package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)
// Update the given comment in the blog with the given content and return a new immutable BlogPost
func UpdateCommentContent(source BlogPost, commentIndex int, newContent string) BlogPost {
	//We can't modify the existing comments so we have to clone them
	var updatedComments []Comment
	updatedComments = append(updatedComments, source.Comments...)

	//Comments are stored by value so we can update our copy
	updatedComments[commentIndex].Content = newContent

	//BlogPost was also passed by value so we can update its Comments directly
	source.Comments = updatedComments

	return source
}

func main() {

	data := NewBlogPost(
		"BlogPost",
		[]Comment{
			NewComment("First Comment", "comment content"),
			NewComment("Second Comment", "comment content"),
		},
	)

	result := UpdateCommentContent(data, 1, "Updated Content2")

	fmt.Println(data)
	fmt.Println(result)
}       
`,
		Result: "{BlogPost [{First Comment comment content} {Second Comment comment content}] []}\n{BlogPost [{First Comment comment content} {Second Comment Updated Content2}] []}\n",
	},
	"TestBasicLitInt": {
		Script: `
package main

import (
	"context"
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

func main() {
	result := MustModify(
		Filtered(
			TraverseSlice[int](),
			Lt(10),
		),
		Mul(2),
		[]int{1, 2, 30, 4, 5},
	)
	fmt.Println(result)
}
                 
`,
		Result: "[2 4 30 8 10]\n",
	},
	"TestOpStringsUpper": {
		Script: `
package main

import (
	"context"
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

func main() {
	result := MustModify(
		TraverseSlice[string](), //Optic
		Op(strings.ToUpper),     //Modify Operation
		[]string{"alpha", "beta"},
	)
	fmt.Println(result)
}                         
`,
		Result: "[ALPHA BETA]\n",
	},
	"TestOpStringsUpperSimplified": {
		Script: `
package main

import (
	"context"
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

func main() {
	result := MustGet(
		Op(strings.ToUpper),
		"alpha",
	)
	fmt.Println(result)
}                         
`,
		Result: "ALPHA\n",
	},
	"TestCustTypeFieldLens": {
		Script: `
package main

import (
	"context"
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

func main() {
	type ExampleStruct struct {
		name    string
		address string
	}

	nameField := FieldLens(
		func(source *ExampleStruct) *string {
			return &source.name
		},
	)

	data := ExampleStruct{
		name:    "Max Mustermann",
		address: "Musterstadt",
	}

	result := MustSet(nameField, "Erika Mustermann", data)
	fmt.Println(data)
	fmt.Println(result)
}                        
`,
		Result: "{Musterstadt Max Mustermann}\n{Musterstadt Erika Mustermann}\n",
	},
	"TestRangeSeq": {
		Script: `

package main

import (
	"context"
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

func main() {
	for v := range MustGet(
		SeqOf(
			TraverseSlice[int](),
		),
		[]int{1, 2, 3, 4, 5},
	) {
		fmt.Println(v)
	}
}
                             
`,
		Result: "1\n2\n3\n4\n5\n",
	},
	"TestDownCast": {
		Script: `


package main

import (
	"context"
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

func main() {
	result := MustGet(
		SliceOf(
			Compose(
				TraverseSlice[any](),
				DownCast[any, int](),
			),
			3, //Initial size of slice
		),
		[]any{1, "two", 3},
	)
	fmt.Println(result)
}
      
`,
		Result: "[1 3]\n",
	},
	"TestAnyfyedReflect": {
		Script: `
package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/internal/playground/data"
)

func main() {
	for commentTitle := range MustGet(
		SeqOf(
			data.O.BlogPost().Comments().Traverse().Title(),
		),
		data.NewBlogPost(
			"Content",
			[]data.Comment{
				data.NewComment("First Comment", "My comment"),
				data.NewComment("Second Comment", "Another comment"),
			},
		),
	) {
		fmt.Println(commentTitle)
	}
}       
      
`,
		Result: "First Comment\nSecond Comment\n",
	},
	"TestOpticalUpdate": {
		Script: `package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)
// Update the given comment in the blog with the given content and return a new immutable BlogPost
func OpticalUpdateCommentContent(blogPost BlogPost, commentIndex int, newContent string) BlogPost {
	return MustSet(
		O.BlogPost().Comments().Nth(commentIndex).Content(),
		newContent,
		blogPost,
	)
}

func main() {
	data := NewBlogPost(
		"BlogPost",
		[]Comment{
			NewComment("First Comment", "comment content"),
			NewComment("Second Comment", "comment content"),
		},
	)

	result := OpticalUpdateCommentContent(data, 1, "Updated Content")

	fmt.Println(result)
}
            `,
		Result: "{BlogPost [{First Comment comment content} {Second Comment Updated Content}] []}\n",
	},
	"TestOpticIdentityReorder": {
		Script: `
package main

import (
	"context"
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

func main() {
	result := MustGet(
		SliceOf(
			Ordered(
				TraverseSlice[int](),
				OrderBy[int,int](Identity[int]()),
			),
			10, //initial size of slice
		),
		[]int{1, 2, 30, 4, 5},
	)
	fmt.Println(result)
}
`,
		Result: "[1 2 4 5 30]\n",
	},
	"TestUsingOpticPredComposition": {
		Script: `
package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func main() {

	blogPostRatings := FieldLens(func(source *BlogPost) *[]Rating {
		return &source.Ratings
	})

	ratingAuthor := FieldLens(func(source *Rating) *string {
		return &source.Author
	})

	ratingStars := FieldLens(func(source *Rating) *int {
		return &source.Stars
	})
	traverseRatings := TraverseSlice[Rating]()
	eqMustermann := Eq("Max Mustermann")
	incOne := Add(1)
	// for _,rating := range blogPost.ratings
	traverseBlogPostRatings := Compose(blogPostRatings, traverseRatings)

	// rating.author == "Max Mustermann"
	ratingAuthorEqMaxMustermann := Compose(ratingAuthor, eqMustermann)

	// if rating.author == "Max Mustermann"
	ifRatingAuthorMustermann := Filtered(traverseBlogPostRatings, ratingAuthorEqMaxMustermann)

	// rating.stars = ....
	matchingRatings := Compose(
		ifRatingAuthorMustermann,
		ratingStars,
	)
	result := MustModify(
		matchingRatings,
		incOne,
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)
	fmt.Println(result)
}
`,
		Result: "{ [] [{Max Mustermann 1} {Erika Mustermann 0}]}\n",
	},
	"TestUsingOpticPredCompositionSimplified": {
		Script: `
package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func main() {

	blogPostRatings := FieldLens(func(source *BlogPost) *[]Rating {
		return &source.Ratings
	})

	result := MustGet(
		blogPostRatings,
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)
	fmt.Println(result)
}
`,
		Result: "[{Max Mustermann 0} {Erika Mustermann 0}]\n",
	},
	"TestUsingOpticPredCompositionSimplified2": {
		Script: `
package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func main() {

	blogPostRatings := FieldLens(func(source *BlogPost) *[]Rating {
		return &source.Ratings
	})

	result := MustSet(
		blogPostRatings,
		[]Rating{
			Rating{
				Author: "Max Mustermann",
				Stars:  10,
			},
			Rating{
				Author: "Erika Mustermann",
				Stars:  0,
			},
		},
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)
	fmt.Println(result)
}
`,
		Result: "{ [] [{Max Mustermann 10} {Erika Mustermann 0}]}\n",
	},
	"TestUsingOpticPredCompositionMakeLens": {
		Script: `

package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func main() {
	result := MustModify(
		Compose(
			Filtered(
				O.BlogPost().Ratings().Traverse(),
				O.Rating().Author().Eq("Max Mustermann"),
			),
			O.Rating().Stars(),
		),
		Add(1),
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)
	fmt.Println(result)
}
`,
		Result: "{ [] [{Max Mustermann 1} {Erika Mustermann 0}]}\n",
	},
	"TestUsingOpticFilterIndexed": {
		Script: `
package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func main() {
	result := MustModify(
		Compose(
			FilteredI(
				O.BlogPost().Ratings().Traverse(),
				OpOnIx[Rating](Even[int]()),
			),
			O.Rating().Stars(),
		),
		Add(1),
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)
	fmt.Println(result)
}
`,
		Result: "{ [] [{Max Mustermann 1} {Erika Mustermann 0}]}\n",
	},
	"TestUsingOpticPredicateAndOp": {
		Script: `

package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func main() {
	result := MustModify(
		Compose(
			Filtered(
				O.BlogPost().Ratings().Traverse(),
				AndOp(
					O.Rating().Author().Eq("Max Mustermann"),
					O.Rating().Stars().Lt(5),
				),
			),
			O.Rating().Stars(),
		),
		Add(1),
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)
	fmt.Println(result)
}

`,
		Result: "{ [] [{Max Mustermann 1} {Erika Mustermann 0}]}\n",
	},
	"TestUsingOpticReordered": {
		Script: `
package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func main() {
	data := []int{3, 2, 5, 4, 1}

	ordered := Ordered(
		TraverseSlice[int](),
		OrderBy[int](Identity[int]()),
	)

	result := MustGet(
		SliceOf(
			ordered,
			len(data),
		),
		data,
	)

	fmt.Println(result)
}
            

`,
		Result: "[1 2 3 4 5]\n",
	},
	"TestUsingOpticCollectionFilteredCol": {
		Script: `
package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)	

func main() {
	result := MustModify(
		O.BlogPost().Ratings(),
		FilteredCol[int](
			O.Rating().Author().Ne("Max Mustermann"),	
		),
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)

	fmt.Println(result)
}
                    

`,
		Result: "{ [] [{Erika Mustermann 0}]}\n",
	},
	"TestReducersSum": {
		Script: `

package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
)

func main() {
	data := []int{1, 2, 3, 4}

	result, ok := MustGetFirst(
		Reduce(
			TraverseSlice[int](),
			Sum[int](),
		),
		data,
	)
	fmt.Println(result, ok)
}
            
`,
		Result: "10 true\n",
	},
	"TestMakeLensSimple": {
		Script: `
package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func main() {
	result, ok := MustGetFirst(
		O.BlogPost().Content().Eq("Blog Content"),
		BlogPost{
			Content: "Blog Content",
			Comments: []Comment{
				Comment{
					Title:   "First Comment",
					Content: "Comment 1 content",
				},
				Comment{
					Title:   "Second Comment",
					Content: "Comment 2 content",
				},
			},
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  5,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  3,
				},
			},
		},
	)
	fmt.Println(result, ok)
}            
`,
		Result: "true true\n",
	},
	"TestCustomGetterE": {
		Script: `
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"iter"
	"strconv"

	. "github.com/spearson78/go-optic"
)

func main() {
	parseInt := GetterE[string, int64](
		func(ctx context.Context, source string) (int64, error) {
			return strconv.ParseInt(source, 10, 0)
		},
		ExprCustom("parseInt"),
	)

	result, err := Get(parseInt, "alpha")
	fmt.Println(result, err)
}        
`,
		Result: `0 strconv.ParseInt: parsing "alpha": invalid syntax
optic error path:
	Custom(parseInt)

`,
	},
	"TestIndexedTraverseMap": {
		Script: `
package main

import (
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

func main() {
	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
		"delta": 4,
	}
	result := MustGet(
		SeqIOf(
			TraverseMap[string, int](),
		),
		data,
	)
	for index, value := range result {
		fmt.Println(index, value)
	}
}       
`,
		Result: `alpha 1
beta 2
delta 4
gamma 3
`,
	},
	"TestIndexedLostIndex": {
		Script: `
package main

import (
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

func main() {
	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
		"delta": 4,
	}
	result, err := Get(
		SeqIOf(
			Compose(
				TraverseMap[string, int](),
				Mul(10),
			),
		),
		data,
	)

	if err != nil {
		fmt.Println(err)
	}

	for index, value := range result {
		fmt.Println(index, value)
	}
}   
`,
		Result: `{} 10
{} 20
{} 40
{} 30
`,
	},
	"TestIndexedComposeBoth": {
		Script: `
package main

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
)

func main() {
	data := []map[string]int{
		map[string]int{
			"alpha": 1,
			"beta":  2,
		},
		map[string]int{
			"gamma": 3,
			"delta": 4,
		},
	}
	result, err := Get(
		SeqIOf(
			ComposeBoth(
				TraverseSlice[map[string]int](),
				TraverseMap[string, int](),
			),
		),
		data,
	)

	if err != nil {
		fmt.Println(err)
	}

	for index, value := range result {
		fmt.Println(index, value)
	}
}
`,
		Result: `{0 alpha} 1
{0 beta} 2
{1 delta} 4
{1 gamma} 3
`,
	},
	"TestIndexedComposeI": {
		Script: `
package main

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
)

func main() {
	data := []map[string]int{
		map[string]int{
			"alpha": 1,
			"beta":  2,
		},
		map[string]int{
			"gamma": 3,
			"delta": 4,
		},
	}

	ixMap := IxMapIso[int, string, lo.Tuple2[int, string]](
		func(left int, right string) lo.Tuple2[int, string] {
			return lo.T2(left, right)
		},
		func(t1, t2 lo.Tuple2[int, string]) bool {
			return t1 == t2
		},
		func(mapped lo.Tuple2[int, string]) (int, bool, string, bool) {
			return mapped.A, true, mapped.B, true
		},
		ExprCustom("IxMapBoth"),
	)

	result, err := Get(
		SeqIOf(
			ComposeI(
				ixMap,
				TraverseSlice[map[string]int](),
				TraverseMap[string, int](),
			),
		),
		data,
	)

	if err != nil {
		fmt.Println(err)
	}

	for index, value := range result {
		fmt.Println(index, value)
	}
}
`,
		Result: `{0 alpha} 1
{0 beta} 2
{1 delta} 4
{1 gamma} 3
`,
	},
	"TestCustomeIndexAware": {
		Script: `
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"iter"
	"strconv"

	. "github.com/spearson78/go-optic"
)

func main() {
	sliceLast := GetterI[int, []string, string](
		func(source []string) (int, string) {
			lastIndex := len(source) - 1
			return lastIndex, source[lastIndex]
		},
		func(indexA, indexB int) bool {
			return indexA == indexB
		},
		ExprCustom("sliceLast"),
	)

	index, result := MustGetI(sliceLast, []string{"alpha", "beta"})
	fmt.Println(result)
}
`,
		Result: `beta
`,
	},
	"NilIxGet": {
		Script: `package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
)

func main() {
	sliceTraversalI := TraversalI[int, []int, int](
		//Iteration function
		func(source []int) SeqI[int, int] {
			return func(yield func(index int, focus int) bool) {
				for i, v := range source {
					if !yield(i, v) {
						break
					}
				}
			}
		},
		//Length getter
		func(source []int) int {
			return len(source)
		},
		//Modify function
		func(fmap func(index int, focus int) int, source []int) []int {
			var modified []int
			for i, v := range source {
				modified = append(modified, fmap(i, v))
			}
			return modified
		},
		//Index getter function
		nil, //IxGet
		//IxMatch
		func(indexA, indexB int) bool {
			return indexA == indexB
		},
		ExprCustom("sliceTraversalI"),
	)

	result, found := MustGetFirst(
		Index(
			sliceTraversalI,
			1,
		),
		[]int{1, 2, 3},
	)
	fmt.Println(result, found)
	//Output:
	//2 true
}
`,
		Result: `2 true
`,
	},
	"OrderByI": {
		Script: `package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func main() {
	data := map[int]string{
		4: "alpha",
		3: "beta",
		2: "gamma",
		1: "delta",
	}

	optic := OrderedI(
		TraverseMap[int, string](),
		OrderByI[int,string,int](
			ValueIIndex[int, string](),
		),
	)

	result := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(result)
}
`,
		Result: `[delta gamma beta alpha]
`,
	},
	"WithMetrics": {
		Script: `
package main

import (
	"fmt"

	. "github.com/spearson78/go-optic"
)

func main() {
	data := []string{"alpha", "beta", "gamma", "delta"}

	sortStringSlice := Ordered(
		TraverseSlice[string](),
		OrderBy[string, string](
			Identity[string](),
		),
	)

	//Attach metrics to the sort
	var m Metrics
	sortStringSlice = WithMetrics(sortStringSlice, &m)

	res := MustGet(
		SliceOf(
			sortStringSlice,
			4,
		),
		data,
	)

	fmt.Println(res)
	fmt.Println(m)
}            
`,
		Result: `[alpha beta delta gamma]
metrics[F:4 A:1 I:0 L:0 Custom:map[comparisons:6]]
`,
	},
	"EqT2": {
		Script: `package main

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
)

func main() {
	data := []string{
		"alpha",
		"beta",
		"gamma",
		"delta",
	}

	optic := ReIndexed(
		TraverseSlice[string](),
		Compose(
			FormatInt[int](10),
			Op(func(focus string) string {
				return "KEY:" + focus
			}),
		),
		EqT2[string](),
	)

	//Focus on elements with an index > 4 where the index is now the string length of the key
	for index, focus := range MustGet(SeqIOf(optic), data) {
		fmt.Println(index, focus)
	}
}`,
		Result: `KEY:0 alpha
KEY:1 beta
KEY:2 gamma
KEY:3 delta
`,
	},
	"SelfIndex": {
		Script: `package main

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
)

func main() {
	data := []lo.Tuple2[string, int]{
		lo.T2("Max Mustermann", 42),
		lo.T2("Erika Mustermann", 37),
	}

	optic := ReIndexed(
		//Self index will set the index to the lo.Tuple2[string,int]
		SelfIndex(
			TraverseSlice[lo.Tuple2[string, int]](),
			EqT2[lo.Tuple2[string, int]](),
		),
		//Reindexed will then set the index to element A of the tuple.
		T2A[string, int](),
		EqT2[string](),
	)

	//We can then build a map the element A of the tuple as key and the fully tuple as the element.
	res, err := Get(MapOf(optic, 10), data)
	fmt.Println(res, err)
}
            `,
		Result: `map[Erika Mustermann:{Erika Mustermann 37} Max Mustermann:{Max Mustermann 42}] <nil>
`,
	},
}

func TestScripts(t *testing.T) {
	for name, test := range scriptTests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(
				context.Background(),
				time.Millisecond*1000,
			)
			defer cancel()

			testScript(ctx, t, test.Script, test.Result)
		})
	}
}

func testScript(ctx context.Context, t *testing.T, script, result string) {

	b := bytes.Buffer{}
	x := bufio.NewWriter(&b)
	goscript.Out = x
	err := goscript.Run(ctx, script)
	if err != nil {
		t.Fatal(err)
	}

	x.Flush()
	res := b.String()

	if res != result {
		t.Fatal(res)
	}
}

func TestSingleScript(t *testing.T) {
	test := scriptTests["SelfIndex"]
	testScript(context.Background(), t, test.Script, test.Result)
}
