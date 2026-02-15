package main

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/examples/recipes/model"
)

func MainPage[T any, RET TReturnOne, RW TReadWrite, DIR, ERR any](w fyne.Window, recipesOptic Optic[Void, T, T, []Recipe, []Recipe, RET, RW, DIR, ERR], state *T) {

	data, _, err := GetFirst(recipesOptic, *state)
	if err != nil {
		dialog.ShowError(err, w)
	}

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {

			nameLabel := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			descriptonLabel := widget.NewLabel("")

			c := container.New(layout.NewVBoxLayout(), nameLabel, descriptonLabel)
			return c
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			c.Objects[0].(*widget.Label).SetText(MustGet(O.Recipe().Name(), data[i]))
			c.Objects[1].(*widget.Label).SetText(MustGet(O.Recipe().Description(), data[i]))
		})

	list.OnSelected = func(id widget.ListItemID) {
		recipeOptic := FirstOrError(Index(Compose(recipesOptic, TraverseSlice[Recipe]()), id), errors.New("list.OnSelected id not found"))
		EditRecipe(w, recipeOptic, func() { MainPage(w, recipesOptic, state) }, state)
	}

	buttons := container.NewHBox(layout.NewSpacer(), widget.NewButton("New", func() {
		newState, err := Modify(recipesOptic, AppendSlice(ValCol(Recipe{})), *state)
		if err != nil {
			dialog.ShowError(err, w)
		}
		*state = newState
		recipeOptic := LastOrError(Compose(recipesOptic, TraverseSlice[Recipe]()), errors.New("list new element not found"))
		EditRecipe(w, recipeOptic, func() { MainPage(w, recipesOptic, state) }, state)
	}))

	page := container.NewBorder(nil, buttons, nil, nil, list)

	w.SetContent(page)
}

func EditRecipe[T any, RET TReturnOne, RW TReadWrite, DIR, ERR any](w fyne.Window, recipeOptic Optic[int, T, T, Recipe, Recipe, RET, RW, DIR, ERR], back func(), state *T) {

	optic := ORecipeOf(recipeOptic)

	recipe, err := Get(recipeOptic, *state)
	if err != nil {
		dialog.ShowError(err, w)
	}

	nameEditor := widget.NewEntryWithData(BindString(state, optic.Name()))

	descEditor := widget.NewEntryWithData(BindString(state, optic.Description()))
	descEditor.MultiLine = true
	descEditor.SetMinRowsVisible(5)
	descEditor.SetText(MustGet(O.Recipe().Description(), recipe))

	servesEditor := widget.NewEntryWithData(BindString(state, Compose(optic.Serves(), AsReverseGet(ParseInt[int](10, 32)))))

	cookTimeEditor := widget.NewEntryWithData(BindString(state, Compose(optic.CookTime(), AsReverseGet(ParseFloat[float64]('f', -1, 64)))))

	difficultyEditor := NewSelectWithBinding(
		[]string{
			"Easy",
			"Intermediate",
			"Expert",
		},
		func(err error) {
			dialog.ShowError(err, w)
		},
		state,
		optic.Difficulty(),
	)

	ingredientList := NewWidgetListWithBinding(
		ColFocusErr(optic.Ingredients()),
		func() *fyne.Container {
			qntEditor := widget.NewEntry()
			qntEditor.SetPlaceHolder("Quantity")

			unitEditor := widget.NewSelect([]string{
				"Gram",
				"Litre",
				"Slice",
				"Of",
			}, nil)

			delButton := widget.NewButton("Delete", func() {})

			ingNameEditor := widget.NewEntry()
			ingNameEditor.SetPlaceHolder("Ingredient")

			c := container.New(layout.NewGridLayout(4),
				qntEditor,
				unitEditor,
				ingNameEditor,
				delButton,
			)
			return c
		},
		func(listItem Optic[int, T, T, Ingredient, Ingredient, ReturnOne, RW, UniDir, Err], widgetIndex int) Optic[Void, T, T, any, any, ReturnOne, ReadWrite, UniDir, Err] {
			switch widgetIndex {
			case 0:
				return EErr(Ret1(Rw(Ud(Compose4(
					listItem,
					O.Ingredient().Quantity(),
					AsReverseGet(ParseFloat[float64]('f', -1, 64)),
					IsoCastE[string, any]()),
				))))
			case 1:
				return EErr(Ret1(Rw(Ud(Compose3(
					listItem,
					O.Ingredient().Unit(),
					IsoCast[Unit, any](),
				)))))
			case 2:
				return EErr(Ret1(Rw(Ud(Compose3(
					listItem,
					O.Ingredient().Name(),
					IsoCastE[string, any](),
				)))))
			case 3:
				return nil
			default:
				panic("unknown widgetIndex")
			}
		},
		func(err error) {
			dialog.ShowError(err, w)
		},
		state,
	)

	ingredientList.OnSelected = func(id widget.ListItemID) {
		ingredientList.UnselectAll()
	}

	detailsForm := container.New(layout.NewFormLayout(),
		widget.NewLabel("Name"), nameEditor,
		widget.NewLabel("Description"), descEditor,
		widget.NewLabel("Serves"), servesEditor,
		widget.NewLabel("Cooking time"), cookTimeEditor,
		widget.NewLabel("Difficulty"), difficultyEditor,
	)

	ingredientsButtons := container.NewHBox(layout.NewSpacer(),
		widget.NewButton("Back", back),
		widget.NewButton("Add ingredient", func() {
			newState, err := Modify(optic.Ingredients(), Append(Ingredient{}), *state)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			*state = newState
			ingredientList.Refresh()
		}),
	)

	ingredientsContainer := container.NewBorder(nil, ingredientsButtons, nil, nil, ingredientList)

	onErr := func(err error) {
		dialog.ShowError(err, w)
	}

	instructionsList := NewWidgetListWithBinding(
		ColFocusErr(optic.Instructions()),
		func() *fyne.Container {

			label := widget.NewLabel("")

			instrEditor := widget.NewEntry()
			instrEditor.MultiLine = true
			instrEditor.SetMinRowsVisible(3)

			delButton := widget.NewButton("Delete", func() {})

			c := container.NewBorder(
				nil,
				nil,
				label,
				delButton,
				instrEditor,
			)
			return c
		},
		func(listItem Optic[int, T, T, string, string, ReturnOne, RW, UniDir, Err], widgetIndex int) Optic[Void, T, T, any, any, ReturnOne, ReadWrite, UniDir, Err] {
			switch widgetIndex {
			case 0:
				return EErr(Ret1(Rw(Ud(
					Compose(
						listItem,
						IsoCastE[string, any](),
					),
				))))
			case 1:
				return EErr(Ret1(Rw(Ud(
					Compose5(
						WithIndex(listItem),
						ValueIIndex[int, string](),
						Add(1),
						AsReverseGet(ParseInt[int](10, 32)),
						IsoCastE[string, any](),
					),
				))))
			case 2:
				return nil //Nil means delete for buttons
			default:
				panic("unknown widgetIndex")
			}
		},
		onErr,
		state,
	)

	instructionsButtons := container.NewHBox(layout.NewSpacer(),
		widget.NewButton("Back", back),
		widget.NewButton("Add instruction", func() {
			newState, err := Modify(optic.Instructions(), Append(""), *state)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			*state = newState
			instructionsList.Refresh()
		}),
	)

	detailsButtons := container.NewHBox(layout.NewSpacer(),
		widget.NewButton("Back", back),
	)

	instructionsContainer := container.NewBorder(nil, instructionsButtons, nil, nil, instructionsList)

	detailsContainer := container.NewBorder(nil, detailsButtons, nil, nil, detailsForm)

	tabs := container.NewAppTabs(
		container.NewTabItem("Details", detailsContainer),
		container.NewTabItem("Ingredients", ingredientsContainer),
		container.NewTabItem("Instructions", instructionsContainer),
	)

	w.SetContent(tabs)
}
