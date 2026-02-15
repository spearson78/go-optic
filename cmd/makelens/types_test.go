package main_test

import (
	"reflect"
	"testing"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/ojson"
)

//go:generate ../../makelens main_test types_test.go types_generated_test.go

type Database struct {
	Drawings []Drawing
}

type Drawing struct {
	Title string
	Pages []Page
	Meta  *MetaData
}

type MetaData struct {
	Author string
	Date   string
	Params map[string]string
}

type Page struct {
	Title string
	Lines []Line
	Css   []string
}

type Line struct {
	Start Point
	End   Point
}

type Point struct {
	X float64
	Y float64
}

func TestDatabaseTypes(t *testing.T) {

	d := Drawing{
		Title: "Test Drawing",
		Pages: []Page{
			{
				Title: "Page 1",
				Lines: []Line{
					{
						Start: Point{
							X: 1,
							Y: 2,
						},
						End: Point{
							X: 3,
							Y: 4,
						},
					},
					{
						Start: Point{
							X: 5,
							Y: 6,
						},
						End: Point{
							X: 7,
							Y: 8,
						},
					},
				},
			},
			{
				Title: "Page 2",
				Lines: []Line{
					{
						Start: Point{
							X: 23,
							Y: 22,
						},
						End: Point{
							X: 21,
							Y: 24,
						},
					},
					{
						Start: Point{
							X: 25,
							Y: 26,
						},
						End: Point{
							X: 27,
							Y: 28,
						},
					},
				},
			},
		},
	}

	expectedNestedModify := Drawing{
		Title: "Test Drawing",
		Pages: []Page{
			{
				Title: "Page 1",
				Lines: []Line{
					{
						Start: Point{
							X: 1,
							Y: 4,
						},
						End: Point{
							X: 3,
							Y: 4,
						},
					},
					{
						Start: Point{
							X: 5,
							Y: 6,
						},
						End: Point{
							X: 7,
							Y: 8,
						},
					},
				},
			},
			{
				Title: "Page 2",
				Lines: []Line{
					{
						Start: Point{
							X: 23,
							Y: 44,
						},
						End: Point{
							X: 21,
							Y: 24,
						},
					},
					{
						Start: Point{
							X: 25,
							Y: 52,
						},
						End: Point{
							X: 27,
							Y: 28,
						},
					},
				},
			},
		},
	}

	nestedModify := MustModify(O.Drawing().Pages().Traverse().Lines().Traverse(), Op(func(focus Line) Line {
		if MustGet(O.Line().Start().X(), focus) == 5 {
			return focus
		} else {
			return MustModify(O.Line().Start().Y(), Mul(2.0), focus)
		}
	}), d)

	if !reflect.DeepEqual(nestedModify, expectedNestedModify) {
		t.Fatalf(`nestedModify : %v`, nestedModify)
	}

	expectedAppendModify := Drawing{
		Title: "Test Drawing",
		Pages: []Page{
			{
				Title: "Page 1",
				Lines: []Line{
					{
						Start: Point{
							X: 1,
							Y: 2,
						},
						End: Point{
							X: 3,
							Y: 4,
						},
					},
					{
						Start: Point{
							X: 5,
							Y: 6,
						},
						End: Point{
							X: 7,
							Y: 8,
						},
					},
					{
						Start: Point{5, 2},
						End:   Point{12, 13},
					},
				},
			},
			{
				Title: "Page 2",
				Lines: []Line{
					{
						Start: Point{
							X: 23,
							Y: 22,
						},
						End: Point{
							X: 21,
							Y: 24,
						},
					},
					{
						Start: Point{
							X: 25,
							Y: 26,
						},
						End: Point{
							X: 27,
							Y: 28,
						},
					},
				},
			},
		},
	}

	//O.Drawing().Pages().Nth(0).Lines()
	x := O.Drawing().Pages().Nth(0).Lines()
	appendModify := MustModify(x, Append(Line{
		Start: Point{5, 2},
		End:   Point{12, 13},
	}), d)

	if !reflect.DeepEqual(appendModify, expectedAppendModify) {
		t.Fatalf(`appendModify : %v`, nestedModify)
	}

	//O.Drawing().Pages().Nth(0).Lines().Nth(1).Start().X()
	startX, ok := MustGetFirst(O.Drawing().Pages().Nth(0).Lines().Nth(1).Start().X(), d)
	if !ok || startX != 5 {
		t.Fatalf(`startX := MustFirstOf(LDrawing.Pages().Index(0).Lines().Index(1).Start().X(), d) : %v`, startX)
	}
	//O.DrawingFrom(ojson.Struct[Drawing]()).Title.X()

	jsonDrawing := ODrawingOf(ojson.Parse[Drawing]())
	jd := []byte("{}")
	njd, errNjd := Set(jsonDrawing.Title(), "Hello World", jd)
	if errNjd != nil || string(njd) != `{"Title":"Hello World","Pages":null,"Meta":null}` {
		t.Fatalf(`Set(jsonDrawing.Title(), "Hello World", jd) : %v : %v`, errNjd, string(njd))
	}

	//O.DrawingFrom(ojson.Struct[Drawing]()).Meta().Option()
	metaRes, err := Set(jsonDrawing.Meta(), &MetaData{
		Author: "alice",
		Date:   "date",
		Params: map[string]string{"a": "p1", "b": "p2"},
	}, njd)

	if err != nil || string(metaRes) != `{"Title":"Hello World","Pages":null,"Meta":{"Author":"alice","Date":"date","Params":{"a":"p1","b":"p2"}}}` {
		t.Fatalf(`Set(jsonDrawing.Meta(), mo.Some(MetaData{ : %v , %v`, string(metaRes), err)
	}

	if keyRes, found, err := GetFirst(jsonDrawing.Meta().Params().Key("a"), metaRes); err != nil || !found || keyRes != "p1" {
		t.Fatalf(`MustPreView(jsonDrawing.Meta().Params().Key("a"), metaRes) : %v : %v `, string(metaRes), err)
	}

	paramsRes, err := Set(jsonDrawing.Meta().Params(), MapCol(map[string]string{"ea": "e1", "eb": "e2"}), metaRes)

	if err != nil || string(paramsRes) != `{"Title":"Hello World","Pages":null,"Meta":{"Author":"alice","Date":"date","Params":{"ea":"e1","eb":"e2"}}}` {
		t.Fatalf(`MustSet(jsonDrawing.Meta().Params(), Map2Seq(map[string]string{"ea": "e1", "eb": "e2"}), metaRes) : %v %v`, string(paramsRes), err)
	}

	metaNilRes, metaNilErr := Set(jsonDrawing.Meta(), nil, metaRes)

	if metaNilErr != nil || string(metaNilRes) != `{"Title":"Hello World","Pages":null,"Meta":null}` {
		t.Fatalf(`Set(jsonDrawing.Meta(), mo.None[MetaData](), metaRes) : %v , %v`, metaNilErr, string(metaNilRes))
	}

	//Pages where there exists a line where the start x is higher than the starty.
	fb := Filtered(O.Drawing().Pages().Traverse().Lines().Traverse(), GtOp(O.Line().Start().X(), O.Line().Start().Y()))
	if res := MustGet(SliceOf(fb, 1), d); len(res) != 1 {
		t.Fatalf(`Compose(LDrawing.Pages().Traverse(), FilteredBy(Filtered(Gter(LLine.Start().X(), LLine.Start().Y()), LPage.Lines().Traverse()))) : %v `, res)
	}

	//Filter Pages where there exists a line where the start x is higher than the starty.
	filtered := FilteredCol[int](EPure(Any(O.Page().Lines().Traverse(), GtOp(O.Line().Start().X(), O.Line().Start().Y()))))
	if res := MustModify(O.Drawing().Pages(), filtered, d); len(res.Pages) != 1 || res.Pages[0].Title != "Page 2" {
		t.Fatalf(`FilteredCol[int](EPure(Compose(Any(O.Page().Lines().Traverse(), GtPred(O.Line().Start().X(), O.Line().Start().Y())), Not()))): %v : %v `, len(res.Pages), res.Pages)
	}
}
