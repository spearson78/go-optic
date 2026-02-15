package main

import (
	"github.com/samber/mo"
	"github.com/spearson78/go-optic"
)
type lTask[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Task,Task,RET,RW,DIR,ERR]
}
func (s *lTask[I,S,T,RET,RW,DIR,ERR])Id()optic.MakeLensRealOps[I,S,T,int64,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Task)*int64{
return &x.Id
})))))))
}
func (s *lTask[I,S,T,RET,RW,DIR,ERR])Title()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Task)*string{
return &x.Title
})))))))
}
func (s *lTask[I,S,T,RET,RW,DIR,ERR])Description()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Task)*string{
return &x.Description
})))))))
}
func (s *lTask[I,S,T,RET,RW,DIR,ERR])Category()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Task)*string{
return &x.Category
})))))))
}
func (s *lTask[I,S,T,RET,RW,DIR,ERR])DueDate()optic.MakeLensRealOps[I,S,T,int64,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Task)*int64{
return &x.DueDate
})))))))
}
func (s *lTask[I,S,T,RET,RW,DIR,ERR])Completed()optic.MakeLensCmpOps[I,S,T,bool,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensCmpOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Task)*bool{
return &x.Completed
})))))))
}
type sTask[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Task, optic.Pure],optic.Collection[int,Task, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Task,[]Task,RET,RW,DIR,ERR]
}
func (s *sTask[I,S,T,RET,RW,DIR,ERR])Traverse()*lTask[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OTaskOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Task]()))))))
}
func (s *sTask[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lTask[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OTaskOf(optic.Index(s.Traverse(),index))
}
type mTask[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Task, optic.Pure],optic.Collection[I,Task, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Task,map[I]Task,RET,RW,DIR,ERR]
}
func (s *mTask[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lTask[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OTaskOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Task]()))))))
}
func (s *mTask[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lTask[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OTaskOf(optic.Index(s.Traverse(),index))
}
type oTask[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Task,*Task,RET,RW,DIR,ERR]
}
func (s *oTask[I,S,T,RET,RW,DIR,ERR])Id()optic.Optic[I,S,T,int64,int64,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Task)*int64{
return &x.Id
}))))))
}
func (s *oTask[I,S,T,RET,RW,DIR,ERR])Title()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Task)*string{
return &x.Title
}))))))
}
func (s *oTask[I,S,T,RET,RW,DIR,ERR])Description()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Task)*string{
return &x.Description
}))))))
}
func (s *oTask[I,S,T,RET,RW,DIR,ERR])Category()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Task)*string{
return &x.Category
}))))))
}
func (s *oTask[I,S,T,RET,RW,DIR,ERR])DueDate()optic.Optic[I,S,T,int64,int64,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Task)*int64{
return &x.DueDate
}))))))
}
func (s *oTask[I,S,T,RET,RW,DIR,ERR])Completed()optic.Optic[I,S,T,bool,bool,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Task)*bool{
return &x.Completed
}))))))
}
func (s *oTask[I,S,T,RET,RW,DIR,ERR])Some()*lTask[optic.Void,mo.Option[Task],mo.Option[Task],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OTaskOf(optic.Some[Task]())
}
func (s *oTask[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Task],mo.Option[Task],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Task]())))))
}
type o struct {
}
func (s *o)Task()*lTask[optic.Void,Task,Task,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OTaskOf[optic.Void,Task,Task,optic.ReturnOne](optic.Identity[Task]())
}
func OTaskOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Task,Task,RET,RW,DIR,ERR])*lTask[I,S,T,RET,RW,DIR,ERR]{
return &lTask[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
var O  = o{
}
