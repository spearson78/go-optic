package main

import (
	. "github.com/spearson78/go-optic"
)

func (s *lTask[I, S, T, RET, RW, DIR, ERR]) DueDateString() Optic[Void, S, T, string, string, RET, CompositionTree[RW, ReadWrite], UniDir, Err] {
	return EErr(RetL(Ud(Compose(s.DueDate(), Ret1(Rw(Compose(UnixTimeIso, DateStringIso)))))))
}
