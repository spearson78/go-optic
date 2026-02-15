package main

import (
	"fmt"
	"iter"
	"log"
	"os"
	"reflect"
	"strings"

	svg "github.com/ajstarks/svgo"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
	"github.com/spearson78/go-optic/internal/util"
	"github.com/spearson78/go-optic/oreflect"
	"github.com/spearson78/go-optic/otree"
)

type OpticTree struct {
	Name      string
	Type      string
	I         string
	S         string
	T         string
	A         string
	B         string
	SubOptics []*OpticTree
	Width     int
	Height    int
}

type options struct {
	unitWidth  int
	unitHeight int
	biDir      bool
}

func (o *OpticTree) String() string {
	return fmt.Sprintf("%v[%v,%v,%v,%v,%v](%v)%v:%v", o.Name, o.I, o.S, o.T, o.A, o.B, o.SubOptics, o.Width, o.Height)
}

func supressCompose(o *OpticTree) bool {

	if o.Type != "Compose" {
		return false
	}

	if isVoid(o.I) {
		return true
	}

	if len(o.SubOptics) == 0 {
		return true
	}

	if o.I == o.SubOptics[len(o.SubOptics)-1].I {
		return true
	}

	return false
}

func buildOpticTree(exp expr.OpticExpression, opt *options) *OpticTree {

	if bindIsoOp, ok := exp.(expr.IsoOpT2BindExpr); ok {
		return buildOpticTree(bindIsoOp.Op, opt)
	}

	var o OpticTree

	switch t := exp.(type) {
	case expr.CustomOptic:
		o.Name = t.Id
	case expr.Compose:
		o.Name = fmt.Sprintf("Compose %v", t.IxMap)
	case expr.TupleOf:
		var sb strings.Builder

		sb.WriteString("TupleOf(")

		for i, v := range t.Elements {
			if i != 0 {
				sb.WriteString(" , ")
			}
			sb.WriteString(v.Short())
		}

		sb.WriteString(")")
		o.Name = sb.String()
	default:
		o.Name = exp.Short()
	}

	o.Type = reflect.TypeOf(exp).Name()

	if o.Name == "TYPESP" {
		o.Name = "Optic"
		o.I = "I"
		o.S = "S"
		o.T = "T"
		o.A = "A"
		o.B = "B"
	} else if o.Name == "TYPES" {
		o.Name = "Optic"
		o.I = "I"
		o.S = "S"
		o.T = "S"
		o.A = "A"
		o.B = "A"
	} else {

		o.I = util.FullTypeName(exp.OpticI())
		o.S = util.FullTypeName(exp.OpticS())
		o.T = util.FullTypeName(exp.OpticT())
		o.A = util.FullTypeName(exp.OpticA())
		o.B = util.FullTypeName(exp.OpticB())
	}

	if binaryExpr, ok := exp.(expr.BinaryExpr); ok {
		o.Name = binaryExpr.Op
		o.S = util.FullTypeName(binaryExpr.L)
		o.T = util.FullTypeName(binaryExpr.L)
	}

	trav := FilteredI(
		oreflect.TraverseMembers[expr.OpticExpression, expr.OpticExpression](),
		OpOnIx[expr.OpticExpression](Op(func(m *otree.PathNode[oreflect.Member]) bool {
			lastMember := m.Value()
			if lastMember.StructField.Name == "IxMap" {
				return false
			}
			if lastMember.StructField.Name == "Pred" {
				return false
			}

			return true
		})),
	)

	maxSubTreeHeight := 1

	for _, subExp := range MustGet(
		SeqIOf(
			trav,
		),
		exp,
	) {

		//fmt.Println(index)

		subTree := buildOpticTree(subExp, opt)
		o.SubOptics = append(o.SubOptics, subTree)
		o.Width += subTree.Width

		if subTree.Height > maxSubTreeHeight {
			maxSubTreeHeight = subTree.Height
		}

	}

	if supressCompose(&o) {
		o.Height = maxSubTreeHeight
	} else {
		o.Width += 2
		o.Height = maxSubTreeHeight + 1

	}

	if opt.biDir {
		o.Height += 1
	}

	return &o
}

const (
	fill          = `fill="black"`
	stroke        = `stroke="black"`
	fontSize      = "15pt"
	unitWidthStd  = 185
	unitHeightStd = 40
)

func connector(canvas *svg.SVG, text string, opt *options) {
	canvas.Polygon([]int{0, 10, 0, opt.unitWidth, opt.unitWidth + 10, opt.unitWidth}, []int{0, opt.unitHeight / 2, opt.unitHeight, opt.unitHeight, opt.unitHeight / 2, 0}, `fill="none"`, stroke, `stroke-width="1"`)
	canvas.Text((opt.unitWidth/2)+5, (opt.unitHeight/2)+5, text, `text-anchor="middle"`, `font-size="`+fontSize+`"`, fill)
}

func revConnector(canvas *svg.SVG, text string, opt *options) {
	canvas.Polygon([]int{0, -10, 0, opt.unitWidth, opt.unitWidth - 10, opt.unitWidth}, []int{0, opt.unitHeight / 2, opt.unitHeight, opt.unitHeight, opt.unitHeight / 2, 0}, `fill="none"`, stroke, `stroke-width="1"`)
	canvas.Text((opt.unitWidth/2)-5, (opt.unitHeight/2)+5, text, `text-anchor="middle"`, `font-size="`+fontSize+`"`, fill)
}

func opNode(canvas *svg.SVG, op string, opt *options) {

	canvas.Polygon([]int{0, 10, -10, 0, opt.unitWidth, opt.unitWidth}, []int{0, opt.unitHeight / 2, opt.unitHeight + (opt.unitHeight / 2), 2 * opt.unitHeight, 2 * opt.unitHeight, 0}, `fill="none"`, stroke, `stroke-width="1"`)
	canvas.Text(opt.unitWidth/2, (opt.unitHeight)+5, op, `text-anchor="middle"`, `font-size="`+fontSize+`"`, fill)

}

func isVoid(t string) bool {
	return t == "Void" || strings.Contains(t, "optic.Void")
}

func treeNode(canvas *svg.SVG, opticTree *OpticTree, opt *options) {

	if opticTree.Name == "HIDDEN" {
		return
	}

	if supressCompose(opticTree) {
		pos := 0
		for _, child := range opticTree.SubOptics {
			extraHeight := opticTree.Height - child.Height
			canvas.Translate(pos*opt.unitWidth, extraHeight*opt.unitHeight)
			treeNode(canvas, child, opt)
			canvas.Gend()
			pos += child.Width
		}
	} else {

		//canvas.Rect(0, 0, opticTree.Width*opt.unitWidth, opt.unitHeight, `fill="none"`, stroke, `stroke-width="1"`)
		canvas.Text((opticTree.Width*opt.unitWidth)/2, (opt.unitHeight/2)+5, opticTree.Name, `text-anchor="middle"`, `font-size="`+fontSize+`"`, fill)

		dockSize := 1
		if opt.biDir {
			dockSize = 2
		}

		canvas.Polygon([]int{0, 0, opt.unitWidth, opt.unitWidth, (opticTree.Width - 1) * opt.unitWidth, (opticTree.Width - 1) * opt.unitWidth, opticTree.Width * opt.unitWidth, opticTree.Width * opt.unitWidth}, []int{0, (opticTree.Height - dockSize) * opt.unitHeight, (opticTree.Height - dockSize) * opt.unitHeight, opt.unitHeight, opt.unitHeight, (opticTree.Height - dockSize) * opt.unitHeight, (opticTree.Height - dockSize) * opt.unitHeight, 0}, `fill="none"`, stroke, `stroke-width="1"`)

		dockingHeight := (opticTree.Height - dockSize) * opt.unitHeight

		canvas.Translate(0, dockingHeight)
		connector(canvas, opticTree.S, opt)
		canvas.Gend()
		canvas.Translate((opticTree.Width-1)*opt.unitWidth, dockingHeight)
		if !isVoid(opticTree.I) {
			connector(canvas, "["+opticTree.I+"] "+opticTree.A, opt)
		} else {
			connector(canvas, opticTree.A, opt)
		}
		canvas.Gend()

		if opt.biDir {
			canvas.Translate(0, dockingHeight+opt.unitHeight)
			revConnector(canvas, opticTree.T, opt)
			canvas.Gend()
			canvas.Translate((opticTree.Width-1)*opt.unitWidth, dockingHeight+opt.unitHeight)
			revConnector(canvas, opticTree.B, opt)
			canvas.Gend()
		}

		pos := 0
		for _, child := range opticTree.SubOptics {

			extraHeight := (opticTree.Height - 1) - child.Height

			canvas.Translate(opt.unitWidth+(pos*opt.unitWidth), opt.unitHeight+(extraHeight*opt.unitHeight))
			treeNode(canvas, child, opt)
			canvas.Gend()
			pos += child.Width
		}
	}

}

func exportDiagram[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], fileName string) {
	exportOpDiagram[ReturnOne, any](o, fileName, nil, options{
		unitWidth:  unitWidthStd,
		unitHeight: unitHeightStd,
		biDir:      false,
	})
}

func exportOpDiagram[ORET TReturnOne, OERR any, I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], fileName string, op Operation[A, B, ORET, OERR], opt options) {

	opticTree := buildOpticTree(o.AsExpr(), &opt)

	out, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	width := opticTree.Width*opt.unitWidth + 20

	if op != nil {
		width += opt.unitWidth * 2
	}

	height := opticTree.Height * opt.unitHeight
	canvas := svg.New(out)
	canvas.Start(width, height)
	canvas.Translate(10, 0)
	treeNode(canvas, opticTree, &opt)

	dockingSize := 1
	if opt.biDir {
		dockingSize = 2
	}

	if op != nil {
		canvas.Translate(opticTree.Width*opt.unitWidth, (opticTree.Height-dockingSize)*opt.unitHeight)

		if o.AsExpr().Short() == "Custom(TYPESP)" {
			opNode(canvas, "op", &opt)
		} else {
			opNode(canvas, op.AsExpr().Short(), &opt)
		}

		canvas.Gend()
	}

	canvas.Gend()
	canvas.End()
}

type BlogPost struct {
	author   string
	title    string
	content  string
	comments []Comment
	ratings  []Rating
}

type Comment struct {
	author  string
	title   string
	content string
}

type Rating struct {
	author string
	stars  int
}

// V2 data types
type BlogPostV2 struct {
	author   string
	title    string
	content  string
	comments []Comment
	ratings  []RatingV2
}

type RatingV2 struct {
	author string
	stars  float64
}

func main() {

	exportDiagram(
		TraverseSlice[int](),
		"../../../docs/static/concepts_1.svg",
	)

	exportDiagram(
		Index(
			Iteration[[]int, int](
				func(source []int) iter.Seq[int] {
					return func(yield func(int) bool) {}
				},
				func(source []int) int {
					return 0
				},
				ExprCustom("HIDDEN"),
			),
			3,
		),
		"../../../docs/static/concepts_2.svg",
	)

	exportDiagram(
		Index(
			TraverseSlice[int](),
			3,
		),
		"../../../docs/static/concepts_3.svg",
	)

	blogComments := FieldLens(func(source *BlogPost) *[]Comment {
		return &source.comments
	})

	blogRatings := FieldLens(func(source *BlogPost) *[]Rating {
		return &source.ratings
	})

	ratingStars := FieldLens(func(source *Rating) *int {
		return &source.stars
	})

	ratingAuthor := FieldLens(func(source *Rating) *string {
		return &source.author
	})

	commentTitle := FieldLens(func(source *Comment) *string {
		return &source.title
	})

	blogCommentTitles := Compose3(
		blogComments,
		TraverseSlice[Comment](),
		commentTitle,
	)

	exportDiagram(
		blogCommentTitles,
		"../../../docs/static/concepts_4.svg",
	)

	exportOpDiagram(
		blogCommentTitles,
		"../../../docs/static/concepts_5.svg",
		Op(strings.ToUpper),
		options{
			unitWidth:  unitWidthStd,
			unitHeight: unitHeightStd,
			biDir:      true,
		},
	)

	exportOpDiagram[ReturnOne, any](
		Iteration[[]int, int](
			func(source []int) iter.Seq[int] {
				return func(yield func(int) bool) {}
			},
			func(source []int) int {
				return 0
			},
			ExprCustom("TYPESP"),
		),
		"../../../docs/static/using_1.svg",
		nil,
		options{
			unitWidth:  unitWidthStd,
			unitHeight: unitHeightStd,
			biDir:      true,
		},
	)

	exportOpDiagram[ReturnOne, any](
		Iteration[[]int, int](
			func(source []int) iter.Seq[int] {
				return func(yield func(int) bool) {}
			},
			func(source []int) int {
				return 0
			},
			ExprCustom("TYPES"),
		),
		"../../../docs/static/using_2.svg",
		nil,
		options{
			unitWidth:  unitWidthStd,
			unitHeight: unitHeightStd,
			biDir:      true,
		},
	)

	exportOpDiagram(
		Lens[int, int](
			func(source int) int {
				return source
			},
			func(focus, source int) int {
				return focus
			},
			ExprCustom("TYPESP"),
		),
		"../../../docs/static/using_3.svg",
		Add(1),
		options{
			unitWidth:  unitWidthStd,
			unitHeight: unitHeightStd,
			biDir:      true,
		},
	)

	exportOpDiagram(
		Compose3(
			blogRatings,
			TraverseSlice[Rating](),
			ratingStars,
		),
		"../../../docs/static/poly_1.svg",
		Add(1),
		options{
			unitWidth:  unitWidthStd,
			unitHeight: unitHeightStd,
			biDir:      true,
		},
	)

	upgradeBlogRatings := LensP[BlogPost, BlogPostV2, []Rating, []RatingV2](
		func(source BlogPost) []Rating {
			return source.ratings
		},
		func(focus []RatingV2, source BlogPost) BlogPostV2 {
			return BlogPostV2{
				//Retain the fields that don't need an upgrade
				author:   source.author,
				title:    source.title,
				content:  source.content,
				comments: source.comments,

				//Use the upgraded focus for ratings
				ratings: focus,
			}
		},
		ExprCustom("upgradeBlogRatings"),
	)

	upgradeTraverseSlice := TraverseSliceP[Rating, RatingV2]()

	upgradeRatingRating := LensP[Rating, RatingV2, int, float64](
		func(source Rating) int {
			return source.stars
		},
		func(focus float64, source Rating) RatingV2 {
			return RatingV2{
				author: source.author,
				stars:  focus,
			}
		},
		ExprCustom("upgradeRatingStars"),
	)

	exportOpDiagram(
		Compose3(
			upgradeBlogRatings,
			upgradeTraverseSlice,
			upgradeRatingRating,
		),
		"../../../docs/static/poly_2.svg",
		IsoCast[int, float64](),
		options{
			unitWidth:  unitWidthStd,
			unitHeight: unitHeightStd,
			biDir:      true,
		},
	)

	exportOpDiagram(
		Compose(
			Filtered(
				Compose(blogRatings, TraverseSlice[Rating]()),
				Compose(ratingAuthor, Eq("Mustermann")),
			),
			ratingStars,
		),
		"../../../docs/static/using_4.svg",
		Add(1),
		options{
			unitWidth:  unitWidthStd,
			unitHeight: unitHeightStd,
			biDir:      true,
		},
	)

	exportDiagram(
		blogComments,
		"../../../docs/static/concepts_7.svg",
	)

	exportDiagram(
		commentTitle,
		"../../../docs/static/concepts_8.svg",
	)

	exportOpDiagram[ReturnOne, any](
		blogCommentTitles,
		"../../../docs/static/concepts_9.svg",
		nil,
		options{
			unitWidth:  unitWidthStd,
			unitHeight: unitHeightStd,
			biDir:      false,
		},
	)

	exportOpDiagram[ReturnOne, any](
		ComposeLeft(
			TraverseSlice[Comment](),
			commentTitle,
		),
		"../../../docs/static/concepts_10.svg",
		nil,
		options{
			unitWidth:  unitWidthStd,
			unitHeight: unitHeightStd,
			biDir:      false,
		},
	)

	exportDiagram(
		Filtered(
			TraverseSlice[int](),
			AndOp(
				Gt(10),
				Lt(40),
			),
		),
		"../../../docs/static/concepts_11.svg",
	)

	exportOpDiagram[ReturnOne, Pure](
		ParseIntP[int64](10, 0),
		"../../../docs/static/custom_isoep.svg",
		Add[int64](10),
		options{
			unitWidth:  unitWidthStd,
			unitHeight: unitHeightStd,
			biDir:      true,
		},
	)

}
