package main

import (
	"database/sql"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/exp/osql"
)

type db struct {
}

func (d db) Task() dbTask {
	return dbTask{
		Optic: Table[int64, Task]("TASK"),
	}
}

type dbTask struct {
	Optic[Void, *sql.DB, *sql.DB, Collection[int64, Task, Err], Collection[int64, Task, Err], ReturnMany, ReadWrite, UniDir, Err]
}

func (d dbTask) Traverse() *lTask[int64, *sql.DB, *sql.DB, ReturnMany, ReadWrite, UniDir, Err] {
	return OTaskOf(
		RetM(Rw(Ud(EErr(Compose(
			d.Optic,
			TraverseColE[int64, Task, Err](),
		))))),
	)
}

var DB db
