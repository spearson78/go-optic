package optic

import (
	"context"
	"iter"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
	"golang.org/x/text/encoding"
)

// TraverseString returns an [Traversal] optic that focuses on the runes in a string.
//
// See:
// - [TraverseStringP] for a version that supports aliases of string and rune
func TraverseString() Optic[int, string, string, rune, rune, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseStringP[string, rune]()
}

// TraverseStringP returns an [Traversal] optic that focuses on the runes in a string.
//
// See:
// - [TraverseStringP] for a simple version that doesn't support aliases of string and rune
func TraverseStringP[S ~string, R ~rune]() Optic[int, S, S, R, R, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, int, S, S, R, R](

		func(ctx context.Context, source S) SeqIE[int, R] {
			return func(yield func(ValueIE[int, R]) bool) {
				i := 0
				for index, focus := range source {
					i++
					if !yield(ValIE(index, R(focus), ctx.Err())) {
						break
					}
				}
			}
		},
		func(ctx context.Context, source S) (int, error) {
			return len([]rune(source)), nil
		},
		func(ctx context.Context, fmap func(index int, focus R) (R, error), source S) (S, error) {
			var ret []rune
			i := 0
			for index, focus := range source {
				i++
				m, err := fmap(index, R(focus))
				err = JoinCtxErr(ctx, err)
				if err != nil {
					return "", err
				}
				ret = append(ret, rune(m))
			}

			return S(ret), ctx.Err()
		},
		func(ctx context.Context, index int, source S) SeqIE[int, R] {
			return func(yield func(ValueIE[int, R]) bool) {
				runes := []rune(source)
				if index >= 0 && index < len(runes)-1 {
					yield(ValIE(index, R(runes[index]), ctx.Err()))
				}
			}
		},
		IxMatchComparable[int](),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Traverse{
				OpticTypeExpr: ot,
			}
		}),
	)
}

type MatchIndex struct {
	Offsets  []int
	Captures []string
}

// MatchIndexOffsets returns a [Lens] that focuses the Offsets field of a [MatchIndex]
//
// See:
//   - [MatchString] for an optic using a [MatchIndex]
func MatchIndexOffsets() Optic[Void, MatchIndex, MatchIndex, []int, []int, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens(func(source *MatchIndex) *[]int { return &source.Offsets })
}

// MatchIndexOffsets returns a [Lens] that focuses the Captures field of a [MatchIndex]
//
// See:
//   - [MatchString] for an optic using a [MatchIndex]
func MatchIndexCaptures() Optic[Void, MatchIndex, MatchIndex, []string, []string, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens(func(source *MatchIndex) *[]string { return &source.Captures })
}

// MatchString returns an [Traversal] that focuses on the n substrings that match the given regexp.
//
// Use an n value of -1 to match all the focuses.
// The index provides access to the capture group offsets and their values
//
// See:
//   - [CaptureString] for a for Traversal that focuses on the indexed capture groups within a regexp..
//   - [CaptureMapString] for a for Traversal that focuses on the named capture groups within a regexp..
//   - [SplitString] for a for Traversal that splits strings based on a regexp.
func MatchString(matchOn *regexp.Regexp, n int) Optic[MatchIndex, string, string, string, string, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, MatchIndex, string, string, string, string](
		func(ctx context.Context, source string) SeqIE[MatchIndex, string] {
			return func(yield func(ValueIE[MatchIndex, string]) bool) {

				remainingSource := source

				offset := 0
				i := 0
				for {
					if n != -1 && n == i {
						break
					}

					match := matchOn.FindStringSubmatchIndex(remainingSource)
					if match == nil {
						break
					}

					focus := remainingSource[match[0]:match[1]]
					if len(focus) > 0 {
						offsets := MustModify(TraverseSlice[int](), If(Eq(-1), Const[int](-1), Add(offset)), match)
						captures := make([]string, (len(match)/2)-1)
						for c := 0; c < len(captures); c++ {
							i := 2 + c
							if match[i] != -1 {
								captures[c] = remainingSource[match[i]:match[i+1]]
							}
						}

						if !yield(ValIE(MatchIndex{offsets, captures}, focus, ctx.Err())) {
							i++
							break
						}
					}

					remainingSource = remainingSource[match[1]:]
					offset = offset + match[1]
					i++
				}

			}
		},
		nil,
		func(ctx context.Context, fmap func(index MatchIndex, focus string) (string, error), source string) (string, error) {

			var sb strings.Builder

			remainingSource := source

			offset := 0
			i := 0
			for {
				if n != -1 && n == i {
					sb.WriteString(remainingSource)
					return sb.String(), nil
				}
				match := matchOn.FindStringSubmatchIndex(remainingSource)
				if match == nil {
					if len(remainingSource) > 0 {
						sb.WriteString(remainingSource)
					}
					return sb.String(), nil
				}

				sb.WriteString(remainingSource[:match[0]])

				focus := remainingSource[match[0]:match[1]]
				if len(focus) > 0 {
					offsets := MustModify(TraverseSlice[int](), Add(offset), match)
					captures := make([]string, (len(match)/2)-1)
					for c := 0; c < len(captures); c++ {
						i := 2 + c
						captures[c] = remainingSource[match[i]:match[i+1]]
					}

					mapped, err := fmap(MatchIndex{offsets, captures}, focus)
					i++
					err = JoinCtxErr(ctx, err)
					if err != nil {
						return "", err
					}

					sb.WriteString(mapped)
				}

				remainingSource = remainingSource[match[1]:]
				offset = offset + match[1]
			}

		},
		nil,
		IxMatchDeep[MatchIndex](),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.MatchString{
				OpticTypeExpr: ot,
				Match:         matchOn,
			}
		}),
	)
}

func applyCaptureMapFmap(source string, names []string, match []int, captures map[string]string, sb *strings.Builder) {

	lastPos := 0

	for i := 0; i < len(match); i += 2 {
		if match[i] != -1 {
			name := names[i>>1]
			prefix := source[lastPos:match[i]]
			capture := captures[name]
			sb.WriteString(prefix)
			sb.WriteString(capture)
			lastPos = match[i+1]
		}
	}

	postfix := source[lastPos:]
	sb.WriteString(postfix)

}

// CaptureMapString returns an [Traversal] that focuses on the named capture groups of the given regexp.
//
// Use an n value of -1 to match all the focuses.
//
// The index is a tuple of the offsets of the capture groups and the full string match.
//
// See:
// - [MatchString] for a for Traversal that focuses on the matches of a given regexp..
// - [CaptureString] for a for Traversal that focuses on the indexed capture groups within a regexp..
// - [SplitString] for a for Traversal that splits strings based on a regexp.
func CaptureMapString(matchOn *regexp.Regexp, n int) Optic[lo.Tuple2[[]int, string], string, string, map[string]string, map[string]string, ReturnMany, ReadWrite, UniDir, Pure] {

	names := make([]string, len(matchOn.SubexpNames())-1)
	for i, v := range matchOn.SubexpNames() {
		if i == 0 {
			continue
		}

		if v == "" {
			v = strconv.Itoa(i - 1)
		}
		names[i-1] = v
	}

	return CombiTraversal[ReturnMany, ReadWrite, Pure, lo.Tuple2[[]int, string], string, string, map[string]string, map[string]string](
		func(ctx context.Context, source string) SeqIE[lo.Tuple2[[]int, string], map[string]string] {
			return func(yield func(ValueIE[lo.Tuple2[[]int, string], map[string]string]) bool) {

				remainingSource := source

				offset := 0
				i := 0
				for {
					if n != -1 && n == i {
						break
					}
					match := matchOn.FindStringSubmatchIndex(remainingSource)
					if match == nil {
						break
					}

					focus := remainingSource[match[0]:match[1]]
					if len(focus) > 0 {
						offsets := MustModify(TraverseSlice[int](), If(Eq(-1), Const[int](-1), Add(offset)), match)
						captures := make(map[string]string, len(names))
						for c := 0; c < len(names); c++ {
							i := (2 * c) + 2
							if match[i] != -1 {
								captureName := names[c]
								captures[captureName] = remainingSource[match[i]:match[i+1]]
							}
						}

						if !yield(ValIE(lo.T2(offsets, focus), captures, ctx.Err())) {
							i++
							break
						}
					}

					remainingSource = remainingSource[match[1]:]
					offset = offset + match[1]
					i++
				}

			}
		},
		nil,
		func(ctx context.Context, fmap func(index lo.Tuple2[[]int, string], focus map[string]string) (map[string]string, error), source string) (string, error) {

			var sb strings.Builder

			remainingSource := source

			offset := 0
			i := 0
			for {
				if n != -1 && n == i {
					if len(remainingSource) > 0 {
						sb.WriteString(remainingSource)
					}
					return sb.String(), nil
				}
				match := matchOn.FindStringSubmatchIndex(remainingSource)
				if match == nil {
					if len(remainingSource) > 0 {
						sb.WriteString(remainingSource)
					}
					return sb.String(), nil
				}

				sb.WriteString(remainingSource[:match[0]])

				focus := remainingSource[match[0]:match[1]]
				if len(focus) > 0 {
					offsets := MustModify(TraverseSlice[int](), Add(offset), match)
					captures := make(map[string]string, len(names))
					for c := 0; c < len(names); c++ {
						i := (2 * c) + 2
						if match[i] != -1 {
							captureName := names[c]
							captures[captureName] = remainingSource[match[i]:match[i+1]]
						}
					}

					mapped, err := fmap(lo.T2(offsets, focus), captures)
					i++
					err = JoinCtxErr(ctx, err)
					if err != nil {
						return "", err
					}

					applyCaptureMapFmap(focus, names, match[2:], mapped, &sb)
				}

				remainingSource = remainingSource[match[1]:]
				offset = offset + match[1]
			}

		},
		nil,
		IxMatchDeep[lo.Tuple2[[]int, string]](),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.CaptureString{
				OpticTypeExpr: ot,
				MatchOn:       matchOn,
			}
		}),
	)
}

func applyCaptureFmap(source string, match []int, captures []string, sb *strings.Builder) {

	lastPos := 0

	for i := 0; i < len(match); i += 2 {
		if match[i] != -1 {
			name := i >> 1
			prefix := source[lastPos:match[i]]
			capture := captures[name]
			sb.WriteString(prefix)
			sb.WriteString(capture)
			lastPos = match[i+1]
		}
	}

	postfix := source[lastPos:]
	sb.WriteString(postfix)

}

// CaptureString returns an [Traversal] that focuses on the indexed capture groups of the given regexp.
//
// Use an n value of -1 to match all the focuses.
//
// The index is a tuple of the offsets of the capture groups and the full string match.
//
// See:
// - [MatchString] for a for Traversal that focuses on the matches of a given regexp.
// - [CaptureMapString] for a for Traversal that focuses on the named capture groups within a regexp.
// - [SplitString] for a for Traversal that splits strings based on a regexp.
func CaptureString(matchOn *regexp.Regexp, n int) Optic[lo.Tuple2[[]int, string], string, string, []string, []string, ReturnMany, ReadWrite, UniDir, Pure] {

	return CombiTraversal[ReturnMany, ReadWrite, Pure, lo.Tuple2[[]int, string], string, string, []string, []string](
		func(ctx context.Context, source string) SeqIE[lo.Tuple2[[]int, string], []string] {
			return func(yield func(ValueIE[lo.Tuple2[[]int, string], []string]) bool) {

				remainingSource := source

				offset := 0
				i := 0
				for {
					if n != -1 && n == i {
						break
					}
					match := matchOn.FindStringSubmatchIndex(remainingSource)
					if match == nil {
						break
					}

					focus := remainingSource[match[0]:match[1]]
					if len(focus) > 0 {
						offsets := MustModify(TraverseSlice[int](), If(Eq(-1), Const[int](-1), Add(offset)), match)
						captures := make([]string, len(matchOn.SubexpNames())-1)
						for c := 0; c < len(captures); c++ {
							i := (2 * c) + 2
							if match[i] != -1 {
								captures[c] = remainingSource[match[i]:match[i+1]]
							}
						}

						if !yield(ValIE(lo.T2(offsets, focus), captures, ctx.Err())) {
							i++
							break
						}
					}

					remainingSource = remainingSource[match[1]:]
					offset = offset + match[1]
					i++
				}

			}
		},
		nil,
		func(ctx context.Context, fmap func(index lo.Tuple2[[]int, string], focus []string) ([]string, error), source string) (string, error) {

			var sb strings.Builder

			remainingSource := source

			offset := 0
			i := 0
			for {
				if n != -1 && n == i {
					if len(remainingSource) > 0 {
						sb.WriteString(remainingSource)
					}
					return sb.String(), nil
				}
				match := matchOn.FindStringSubmatchIndex(remainingSource)
				if match == nil {
					if len(remainingSource) > 0 {
						sb.WriteString(remainingSource)
					}
					return sb.String(), nil
				}

				sb.WriteString(remainingSource[:match[0]])

				focus := remainingSource[match[0]:match[1]]
				if len(focus) > 0 {
					offsets := MustModify(TraverseSlice[int](), Add(offset), match)
					captures := make([]string, len(matchOn.SubexpNames())-1)
					for c := 0; c < len(captures); c++ {
						i := (2 * c) + 2
						if match[i] != -1 {
							captures[c] = remainingSource[match[i]:match[i+1]]
						}
					}

					mapped, err := fmap(lo.T2(offsets, focus), captures)
					i++
					err = JoinCtxErr(ctx, err)
					if err != nil {
						return "", err
					}

					applyCaptureFmap(focus, match[2:], mapped, &sb)
				}

				remainingSource = remainingSource[match[1]:]
				offset = offset + match[1]
			}

		},
		nil,
		IxMatchDeep[lo.Tuple2[[]int, string]](),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.CaptureString{
				OpticTypeExpr: ot,
				MatchOn:       matchOn,
			}
		}),
	)
}

// SplitString returns an [Traversal] that focuses on the substrings separated by the splitOn parameter.
//
// See:
// - [MatchString] for a for a version that focuses on the match
// - [Worded] for a word splitting Traversal
// - [Lined] for line splitting Traversal
func SplitString(splitOn *regexp.Regexp) Optic[int, string, string, string, string, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, int, string, string, string, string](
		func(ctx context.Context, source string) SeqIE[int, string] {
			return func(yield func(ValueIE[int, string]) bool) {

				remainingSource := source

				i := 0
				for {
					match := splitOn.FindStringIndex(remainingSource)
					if match == nil {
						if len(remainingSource) > 0 {
							yield(ValIE(i, remainingSource, ctx.Err()))
							i++
						}
						break
					}

					focus := remainingSource[:match[0]]
					if len(focus) > 0 {
						if !yield(ValIE(i, focus, ctx.Err())) {
							i++
							break
						}
					}

					remainingSource = remainingSource[match[1]:]
					i++
				}
			}
		},
		nil,
		func(ctx context.Context, fmap func(index int, focus string) (string, error), source string) (string, error) {

			var sb strings.Builder

			remainingSource := source

			i := 0
			for {
				match := splitOn.FindStringIndex(remainingSource)
				if match == nil {
					if len(remainingSource) > 0 {
						end, err := fmap(i, remainingSource)
						i++
						err = JoinCtxErr(ctx, err)
						if err != nil {
							return "", err
						}

						sb.WriteString(end)
					}
					return sb.String(), nil
				}

				focus := remainingSource[:match[0]]
				if len(focus) > 0 {
					mapped, err := fmap(i, focus)
					i++
					err = JoinCtxErr(ctx, err)
					if err != nil {
						return "", err
					}

					sb.WriteString(mapped)
				}
				sb.WriteString(remainingSource[match[0]:match[1]])

				remainingSource = remainingSource[match[1]:]
			}

		},
		nil,
		IxMatchComparable[int](),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.SplitString{
				OpticTypeExpr: ot,
				SplitOn:       splitOn,
			}
		}),
	)
}

var wordSplit = regexp.MustCompile(`\s+`)

// Worded returns a [Traversal] over the words in the source string.
//
// See:
// - []
func Worded() Optic[int, string, string, string, string, ReturnMany, ReadWrite, UniDir, Pure] {
	return SplitString(wordSplit)
}

var lineSplit = regexp.MustCompile(`[\r\n]+`)

// Lined returns a [Traversal] over the line feed separated lines in the source string.
func Lined() Optic[int, string, string, string, string, ReturnMany, ReadWrite, UniDir, Pure] {
	return SplitString(lineSplit)
}

// StringToCol returns an [Iso] that converts a string to an [Collection]
func StringToCol() Optic[Void, string, string, Collection[int, rune, Pure], Collection[int, rune, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, string, string, Collection[int, rune, Pure], Collection[int, rune, Pure]](
		func(ctx context.Context, source string) (Collection[int, rune, Pure], error) {
			return ColI(
				func(yield func(index int, focus rune) bool) {
					f := 0
					for i, v := range source {
						f++
						if !yield(i, v) {
							break
						}
					}
				},
				func(index int) iter.Seq2[int, rune] {
					return func(yield func(index int, focus rune) bool) {
						runes := []rune(source)
						if index >= 0 && index < len(runes)-1 {
							yield(index, runes[index])
						}
					}
				},
				IxMatchComparable[int](),
				func() int {
					return len(source)
				},
			), nil
		},
		func(ctx context.Context, focus Collection[int, rune, Pure]) (string, error) {
			if focus == nil {
				return "", nil
			}
			var sb strings.Builder
			i := 0
			var retErr error
			focus.AsIter()(ctx)(func(val ValueIE[int, rune]) bool {
				_, focus, err := val.Get()
				err = JoinCtxErr(ctx, err)
				if err != nil {
					retErr = err
					return false
				}
				i++
				sb.WriteRune(focus)
				return true
			})
			return sb.String(), retErr
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ToCol{
				OpticTypeExpr: ot,
				I:             reflect.TypeFor[int](),
				A:             reflect.TypeFor[rune](),
				B:             reflect.TypeFor[rune](),
			}
		}),
	)
}

// EncodeString returns an [Iso] that focuses the string as a bytes with the given encoding.
func EncodeString(encoding encoding.Encoding) Optic[Void, string, string, []byte, []byte, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoEP[string, string, []byte, []byte](
		func(ctx context.Context, source string) ([]byte, error) {
			return encoding.NewEncoder().Bytes([]byte(source))
		},
		func(ctx context.Context, focus []byte) (string, error) {
			ret, err := encoding.NewDecoder().Bytes(focus)
			if err != nil {
				return "", err
			}
			return string(ret), nil
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.EncodeString{
				OpticTypeExpr: ot,
				Encoding:      encoding,
			}
		}),
	)
}

//go:generate ./makecolops  -nocolof -eq "EqT2[rune]()" string_ops.go "optic" "String" "" "" "" "StringColType()" "" "int" "string" "string" "rune" "rune"
func StringColType() CollectionType[int, string, string, rune, rune, Pure] {
	return ColType(
		StringToCol(),
		TraverseString(),
	)
}

// StringHasPrefix returns an [Operator] that focus the result of [strings.HasPrefix]
func StringHasPrefix(prefix string) Optic[Void, string, string, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, Void, string, string, bool, bool](
		func(ctx context.Context, source string) (Void, bool, error) {
			return Void{}, strings.HasPrefix(string(source), string(prefix)), nil
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.StringHasPrefix{
				OpticTypeExpr: ot,
				Prefix:        string(prefix),
			}
		}),
	)
}

// StringHasSuffix returns an [Operator] that focus the result of [strings.HasSuffix]
func StringHasSuffix(suffix string) Optic[Void, string, string, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, Void, string, string, bool, bool](
		func(ctx context.Context, source string) (Void, bool, error) {
			return Void{}, strings.HasSuffix(source, suffix), nil
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.StringHasSuffix{
				OpticTypeExpr: ot,
				Suffix:        suffix,
			}
		}),
	)
}

// StringBuilderReducer returns a reducer which appends the elements.
func StringBuilderReducer() ReductionP[*strings.Builder, string, string, Pure] {
	return ReducerP[*strings.Builder, string, string](
		func() *strings.Builder {
			return &strings.Builder{}
		},
		func(state *strings.Builder, appendVal string) *strings.Builder {
			state.WriteString(appendVal)
			return state
		},
		func(state *strings.Builder) string {
			return state.String()
		},
		ReducerExprDef(
			func(t expr.ReducerTypeExpr) expr.ReducerExpression {
				return expr.StringBuilderReducerExpr{
					ReducerTypeExpr: t,
				}
			},
		),
	)
}

// The StringOf combinator focuses on a string composed of the runes focused by the given optic.
//
// Under modification this string can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified string contains fewer runes the result will use runes from the original source.
// If the modified string contains more runes they will be ignored.
func StringOf[I, S, T, RETI, RW, DIR, ERR any](o Optic[I, S, T, rune, rune, RETI, RW, DIR, ERR], size int) Optic[Void, S, T, string, string, ReturnOne, RW, UniDir, ERR] {
	return CombiLens[RW, ERR, Void, S, T, string, string](
		func(ctx context.Context, source S) (Void, string, error) {
			ret := make([]rune, 0, size)
			var err error
			o.AsIter()(ctx, source)(func(val ValueIE[I, rune]) bool {
				_, a, focusErr := val.Get()
				err = JoinCtxErr(ctx, focusErr)
				if err != nil {
					return false
				}
				ret = append(ret, a)
				return true
			})

			if err != nil {
				return Void{}, "", err
			}

			return Void{}, string(ret), err
		},
		func(ctx context.Context, va string, vs S) (T, error) {
			i := 0
			runes := []rune(va)
			l := len(va)
			ret, err := o.AsModify()(ctx, func(index I, focus rune) (rune, error) {
				if i >= l {
					return focus, ctx.Err()
				} else {
					ret := runes[i]
					i++
					return ret, ctx.Err()
				}
			}, vs)

			return ret, JoinCtxErr(ctx, err)
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.StringOf{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}
