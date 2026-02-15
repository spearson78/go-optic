package main

import (
	"github.com/rivo/tview"
	. "github.com/spearson78/go-optic"
)

func NewBoundInputField[I, T any, RET TReturnOne, RW TReadWrite, DIR, ERR any](label string, fieldWidth int, onOk func(), onErr func(error), state *T, o Optic[I, T, T, string, string, RET, RW, DIR, ERR]) *tview.InputField {
	value, err := Get(o, *state)
	if err != nil {
		if onErr != nil {
			onErr(err)
		}
		value = ""
	}
	return tview.NewInputField().
		SetLabel(label).
		SetText(value).
		SetFieldWidth(fieldWidth).
		SetChangedFunc(func(text string) {
			newState, err := Set(o, text, *state)
			if err != nil {
				if onErr != nil {
					onErr(err)
				}
			} else {
				*state = newState
				if onOk != nil {
					onOk()
				}
			}
		})
}

func NewBoundTextArea[I, T any, RET TReturnOne, RW TReadWrite, DIR any, ERR any](label string, fieldWidth, fieldHeight, maxLength int, onErr func(error), state *T, o Optic[I, T, T, string, string, RET, RW, DIR, ERR]) *tview.TextArea {
	text, err := Get(o, *state)
	if err != nil {
		onErr(err)
		text = ""
	}

	if fieldHeight == 0 {
		fieldHeight = tview.DefaultFormFieldHeight
	}
	textArea := tview.NewTextArea().
		SetLabel(label).
		SetSize(fieldHeight, fieldWidth).
		SetMaxLength(maxLength)
	if text != "" {
		textArea.SetText(text, true)
	}

	textArea.SetChangedFunc(func() {
		newState, err := Set(o, textArea.GetText(), *state)
		if err != nil {
			onErr(err)
		} else {
			*state = newState
		}
	})

	return textArea
}

func NewBoundCheckBox[I, T any, RET TReturnOne, RW TReadWrite, DIR any, ERR any](label string, onErr func(error), state *T, o Optic[I, T, T, bool, bool, RET, RW, DIR, ERR]) *tview.Checkbox {

	checked, err := Get(o, *state)
	if err != nil {
		onErr(err)
	}

	return tview.NewCheckbox().
		SetLabel(label).
		SetChecked(checked).
		SetChangedFunc(func(checked bool) {
			newState, err := Set(o, checked, *state)
			if err != nil {
				onErr(err)
			} else {
				*state = newState
			}
		})

}
