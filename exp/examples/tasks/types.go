package main

//go:generate ../../../makelens main types.go types_generated.go

type Task struct {
	Id          int64  `db:"ID" osql:"PK"`
	Title       string `db:"TITLE"`
	Description string `db:"DESCRIPTION"`
	Category    string `db:"CATEGORY"`
	DueDate     int64  `db:"DUEDATE"`
	Completed   bool   `db:"COMPLETED"`
}
