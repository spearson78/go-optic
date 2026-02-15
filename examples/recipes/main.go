package main

import (
	"context"
	"os"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/examples/recipes/model"
	"github.com/spearson78/go-optic/ojson"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func FileContents() Optic[Void, string, string, string, string, ReturnOne, ReadWrite, UniDir, Err] {
	return LensE[string, string](
		func(ctx context.Context, source string) (string, error) {
			b, err := os.ReadFile(source)
			if err != nil {
				return "", err
			}
			return string(b), nil
		},
		func(ctx context.Context, focus, source string) (string, error) {
			f, err := os.Create(source)
			if err != nil {
				return source, err
			}
			_, err = f.Write([]byte(focus))
			return source, err
		},
		ExprCustom("FileContents"),
	)
}

func main() {

	fileName := "recipes.json"

	recipesJson := Compose(
		FileContents(),
		ojson.ParseString[[]Recipe](),
	)

	a := app.New()
	w := a.NewWindow("Optics Recipe Example")
	w.Resize(fyne.NewSize(480, 640))

	MainPage(w, recipesJson, &fileName)

	w.ShowAndRun()
}
