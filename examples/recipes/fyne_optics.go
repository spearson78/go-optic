package main

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	. "github.com/spearson78/go-optic"
)

//This is a very basic Optical binding for the Fyne ui toolkit

// BindString provides a [binding.String] using the given optic to focus on a string in the current state.
func BindString[I, S any, RET TReturnOne, RW TReadWrite, DIR any, ERR any](state *S, o Optic[I, S, S, string, string, RET, RW, DIR, ERR]) binding.String {
	return &opticBoundString[I, S, RET, DIR, ERR]{
		state: state,
		o:     Rw(o),
	}
}

// opticBoundString implements [binding.String]
type opticBoundString[I, S any, RET TReturnOne, DIR, ERR any] struct {
	state *S
	o     Optic[I, S, S, string, string, RET, ReadWrite, DIR, ERR]
}

func (p *opticBoundString[I, S, RET, DIR, ERR]) AddListener(l binding.DataListener) {
	//Optics focus immutable data we only need to fire the data changed event once to cause fyne to read the current value.
	l.DataChanged()
}

func (p *opticBoundString[I, S, RET, DIR, ERR]) RemoveListener(binding.DataListener) {
	//Optics focus immutable data no listener was added
}

func (p *opticBoundString[I, S, RET, DIR, ERR]) Get() (string, error) {
	//Get the current value
	r, err := Get(p.o, *p.state)
	return r, err
}

func (p *opticBoundString[I, S, RET, DIR, ERR]) Set(str string) error {
	//Set the new value and update the shared state
	r, err := Set(p.o, str, *p.state)
	if err == nil {
		*p.state = r
	}
	return err
}

// bindSelect binds the selected option in the widget to the the int focused by the optic in the state
func bindSelect[I, S any, A ~int, RET TReturnOne, RW TReadWrite, DIR any, ERR any](w *widget.Select, onErr func(error), state *S, selectedOptionOptic Optic[I, S, S, A, A, RET, RW, DIR, ERR]) {

	w.OnChanged = nil

	//Get the initial selected option
	i, _, err := GetFirst(selectedOptionOptic, *state)
	if err != nil {
		onErr(err)
		return
	}

	w.SetSelectedIndex(int(i))

	w.OnChanged = func(s string) {
		//Find thr selected option's index
		for i, v := range w.Options {
			if s == v {
				//Set it in the current state
				r, err := Set(selectedOptionOptic, A(i), *state)
				if err != nil {
					onErr(err)
					return
				}
				*state = r
				break
			}
		}
	}
}

// NewSelectWithBinding creates a [widget.Select] with the selected option bound to the int focused by the optic.
func NewSelectWithBinding[I, S any, A ~int, RET TReturnOne, RW TReadWrite, DIR, ERR any](options []string, onErr func(error), state *S, o Optic[I, S, S, A, A, RET, RW, DIR, ERR]) *widget.Select {
	editor := widget.NewSelect(options, nil)
	bindSelect(
		editor,
		onErr,
		state,
		o,
	)
	return editor
}

// NewWidgetListWithBinding creates a [widget.List] with each sub widget of the container bound to the corresponding bindings optic. The o Optic focuses te value to be passed to the binsings for a given [widget.listItemID]
func NewWidgetListWithBinding[I, J, S, A any, RET any, RETB TReturnOne, RW, RWB TReadWrite, DIR, DIRB, ERR, ERRB any](
	o Optic[I, S, S, Collection[int, A, ERR], Collection[int, A, ERR], RET, RW, DIR, ERR],
	createItem func() *fyne.Container,
	bind func(listItem Optic[int, S, S, A, A, ReturnOne, RW, UniDir, Err], listItemId widget.ListItemID) Optic[J, S, S, any, any, RETB, RWB, DIRB, ERRB],
	onErr func(error),
	state *S,
) *widget.List {
	var ingredientList *widget.List
	ingredientList = widget.NewList(
		//List length function
		func() int {
			l, err := Get(Length(Compose(o, TraverseColE[int, A, ERR]())), *state)
			if err != nil {
				onErr(err)
				return 0
			}
			return l
		},
		//Create template widget
		func() fyne.CanvasObject {
			return createItem()
		},
		//Bind the template widget to the given ListItemID
		func(listItemId widget.ListItemID, co fyne.CanvasObject) {

			//Iterate through the sub widgets in the container
			for bindingIndex, c := range co.(*fyne.Container).Objects {

				//Focus the correct value for this sub widget.
				binding := bind(
					RwL(
						FirstOrError(
							//Find the focused value from the o Optic
							Index(
								Compose(
									o,
									TraverseColE[int, A, ERR](),
								),
								listItemId,
							),
							errors.New("item with id not found"),
						),
					),
					bindingIndex,
				)

				if btn, ok := c.(*widget.Button); ok {

					if binding != nil {
						panic("non nil bindings for buttons not supported")
					}

					//nil binding means delete action
					btn.OnTapped = func() {
						newState, err := Modify(
							o,
							//Remove the element with this listItemId
							FilteredColI(
								CombiEErr[ERR](EqI[A](listItemId)),
								IxMatchComparable[widget.ListItemID](),
							),
							*state,
						)
						if err != nil {
							onErr(err)
							return
						}

						*state = newState
						ingredientList.Refresh()
					}
				} else {

					//Get the value for this sub widget
					str, err := Get(
						binding,
						*state,
					)
					if err != nil {
						onErr(err)
					}

					//Bind each widget to the correct binding
					switch w := c.(type) {
					case *widget.Label:
						//Set the value into the sub widget
						w.SetText(fmt.Sprintf("%v", str))

					case *widget.Entry:

						//Clear the previous onChanged before setting the text
						w.OnChanged = nil
						//Set the value into the sub widget
						w.SetText(str.(string))
						//Handle any value changes using the binding optic
						w.OnChanged = func(s string) {
							//Set the new value
							r, err := Set(binding, any(s), *state)
							w.SetValidationError(err) //Display errors as validation errors
							if err != nil {
								onErr(err)
								return
							}
							*state = r
						}

					case *widget.Select:
						bindSelect(w, onErr, state, Compose(binding, IsoCast[any, int]()))
					default:
						panic(fmt.Errorf("unsupported widget type %T", w))
					}
				}
			}
		},
	)

	return ingredientList

}
