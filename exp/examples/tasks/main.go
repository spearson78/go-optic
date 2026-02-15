package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic"
	_ "modernc.org/sqlite"
)

func openDB() (*sql.DB, error) {

	db, err := sql.Open("sqlite", "tasks.db")
	if err != nil {
		return nil, err
	}

	row := db.QueryRow("SELECT 1 FROM sqlite_master WHERE type='table' AND name='TASK'")
	var r int
	err = row.Scan(&r)
	if errors.Is(err, sql.ErrNoRows) {

		_, err := db.Exec(`create table TASK (ID INTEGER PRIMARY KEY,TITLE TEXT,DESCRIPTION TEXT,CATEGORY TEXT,DUEDATE INTEGER,COMPLETED INTEGER);`)
		if err != nil {
			return nil, err
		}

		task := Task{
			Title:       "Example Task",
			Description: "This is an example task",
			DueDate:     time.Now().Add(time.Hour * 24).Unix(),
			Category:    "Example",
			Completed:   false,
		}

		db, err = Modify(DB.Task(), AppendCol(ColErr(ValColI[int64](IxMatchComparable[int64](), ValI(int64(0), task)))), db)

		return db, err

	} else {
		return db, err
	}
}

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	app := tview.NewApplication().EnableMouse(true).EnablePaste(true)

	showTaskList(app, db)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

var UnixTimeIso = optic.Iso[int64, time.Time](
	func(focus int64) time.Time {
		return time.Unix(focus, 0)
	},
	func(source time.Time) int64 {
		return source.Unix()
	},
	ExprCustom("UnixTimeIso"),
)

var DateStringIso = optic.IsoE[time.Time, string](
	func(ctx context.Context, source time.Time) (string, error) {
		return source.Format("01/02/2006"), nil
	},
	func(ctx context.Context, focus string) (time.Time, error) {
		return time.Parse("01/02/2006", focus)
	},
	ExprCustom("DateStringIso"),
)

func editTask(app *tview.Application, db *sql.DB, task Task) {

	dateValid := true

	form := tview.NewForm()

	onErr := func(err error) {
		if err != nil {
			modal := tview.NewModal().
				SetText("Error: " + err.Error()).
				AddButtons([]string{"Ok"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					app.SetRoot(form, true)
				})
			app.SetRoot(modal, false)
		}
	}

	form.
		AddFormItem(NewBoundInputField("Title", 64, nil, nil, &task, O.Task().Title())).
		AddFormItem(NewBoundTextArea("Description", 64, 5, 0, onErr, &task, O.Task().Description())).
		AddFormItem(NewBoundInputField("Category", 64, nil, nil, &task, O.Task().Category())).
		AddFormItem(NewBoundInputField("Due Date", 64, func() { dateValid = true }, func(err error) {
			dateValid = err == nil
		}, &task, O.Task().DueDateString())).
		AddFormItem(NewBoundCheckBox("Completed", onErr, &task, O.Task().Completed())).
		AddButton("Save", func() {
			if dateValid {
				if task.Id == 0 {

					_, err := Modify(DB.Task(), AppendCol(ColErr(ValColI[int64](IxMatchComparable[int64](), ValI(int64(0), task)))), db)
					if err != nil {
						modal := tview.NewModal().
							SetText(err.Error()).
							AddButtons([]string{"Ok"}).
							SetDoneFunc(func(buttonIndex int, buttonLabel string) {
								app.SetRoot(form, true)
							})
						app.SetRoot(modal, false)
					} else {
						showTaskList(app, db)
					}

				} else {

					_, err := Set(Index(DB.Task().Traverse(), task.Id), task, db)
					if err != nil {
						modal := tview.NewModal().
							SetText(err.Error()).
							AddButtons([]string{"Ok"}).
							SetDoneFunc(func(buttonIndex int, buttonLabel string) {
								app.SetRoot(form, true)
							})
						app.SetRoot(modal, false)
					} else {
						showTaskList(app, db)
					}
				}
			} else {
				_, dateErr := Get(O.Task().DueDateString(), task)
				modal := tview.NewModal().
					SetText("Invalid Date : " + dateErr.Error()).
					AddButtons([]string{"Ok"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						app.SetRoot(form, true)
					})
				app.SetRoot(modal, false)
			}
		}).
		AddButton("Cancel", func() {
			showTaskList(app, db)
		})

	form.SetBorder(true).SetTitle(" Optic Task Example ")

	app.SetRoot(form, true)
}

func showTaskList(app *tview.Application, db *sql.DB) {

	list := tview.NewList().ShowSecondaryText(false)
	list.SetBorder(true).SetTitle(" Optic Task Example ")

	tasks, err := Get(SliceOf(DB.Task().Traverse(), 10), db)
	if err != nil {
		log.Fatal(err)
	}

	for _, task := range tasks {
		list.AddItem(task.Title, "", 0, func() { editTask(app, db, task) })
	}

	frame := tview.NewFrame(list).
		SetBorders(0, 1, 0, 0, 0, 0).
		AddText("(q)uit | (n)ew task | (d)elete task", false, tview.AlignLeft, tcell.ColorWhite)

	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			app.Stop()
			return nil
		case 'n':
			editTask(app, db, Task{})
			return nil
		case 'd':
			id, ok := MustGetFirst(OTaskOf(Index((TraverseSlice[Task]()), list.GetCurrentItem())).Id(), tasks)
			if !ok {
				log.Fatal("delete: task not found")
			}

			db, err := Modify(DB.Task(), FilteredCol[int64](EErr(O.Task().Id().Ne(id))), db)
			if err != nil {
				log.Fatal(err)
			}

			showTaskList(app, db)
			return nil
		default:
			return event
		}
	})

	app.SetRoot(frame, true)
}
