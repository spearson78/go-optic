package osql_test

import (
	"github.com/samber/mo"
	"github.com/spearson78/go-optic"
)
type lHobby[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Hobby,Hobby,RET,RW,DIR,ERR]
}
func (s *lHobby[I,S,T,RET,RW,DIR,ERR])Id()optic.MakeLensRealOps[I,S,T,int,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Hobby)*int{
return &x.id
})))))))
}
func (s *lHobby[I,S,T,RET,RW,DIR,ERR])Name()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Hobby)*string{
return &x.name
})))))))
}
type sHobby[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Hobby, optic.Pure],optic.Collection[int,Hobby, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Hobby,[]Hobby,RET,RW,DIR,ERR]
}
func (s *sHobby[I,S,T,RET,RW,DIR,ERR])Traverse()*lHobby[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OHobbyOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Hobby]()))))))
}
func (s *sHobby[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lHobby[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OHobbyOf(optic.Index(s.Traverse(),index))
}
type mHobby[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Hobby, optic.Pure],optic.Collection[I,Hobby, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Hobby,map[I]Hobby,RET,RW,DIR,ERR]
}
func (s *mHobby[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lHobby[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OHobbyOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Hobby]()))))))
}
func (s *mHobby[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lHobby[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OHobbyOf(optic.Index(s.Traverse(),index))
}
type oHobby[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Hobby,*Hobby,RET,RW,DIR,ERR]
}
func (s *oHobby[I,S,T,RET,RW,DIR,ERR])Id()optic.Optic[I,S,T,int,int,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Hobby)*int{
return &x.id
}))))))
}
func (s *oHobby[I,S,T,RET,RW,DIR,ERR])Name()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Hobby)*string{
return &x.name
}))))))
}
func (s *oHobby[I,S,T,RET,RW,DIR,ERR])Some()*lHobby[optic.Void,mo.Option[Hobby],mo.Option[Hobby],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OHobbyOf(optic.Some[Hobby]())
}
func (s *oHobby[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Hobby],mo.Option[Hobby],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Hobby]())))))
}
type lPerson[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Person,Person,RET,RW,DIR,ERR]
}
func (s *lPerson[I,S,T,RET,RW,DIR,ERR])Id()optic.MakeLensRealOps[I,S,T,int,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Person)*int{
return &x.id
})))))))
}
func (s *lPerson[I,S,T,RET,RW,DIR,ERR])Name()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Person)*string{
return &x.name
})))))))
}
func (s *lPerson[I,S,T,RET,RW,DIR,ERR])Age()optic.MakeLensRealOps[I,S,T,int,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Person)*int{
return &x.age
})))))))
}
type sPerson[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Person, optic.Pure],optic.Collection[int,Person, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Person,[]Person,RET,RW,DIR,ERR]
}
func (s *sPerson[I,S,T,RET,RW,DIR,ERR])Traverse()*lPerson[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPersonOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Person]()))))))
}
func (s *sPerson[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lPerson[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPersonOf(optic.Index(s.Traverse(),index))
}
type mPerson[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Person, optic.Pure],optic.Collection[I,Person, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Person,map[I]Person,RET,RW,DIR,ERR]
}
func (s *mPerson[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lPerson[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPersonOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Person]()))))))
}
func (s *mPerson[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lPerson[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPersonOf(optic.Index(s.Traverse(),index))
}
type oPerson[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Person,*Person,RET,RW,DIR,ERR]
}
func (s *oPerson[I,S,T,RET,RW,DIR,ERR])Id()optic.Optic[I,S,T,int,int,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Person)*int{
return &x.id
}))))))
}
func (s *oPerson[I,S,T,RET,RW,DIR,ERR])Name()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Person)*string{
return &x.name
}))))))
}
func (s *oPerson[I,S,T,RET,RW,DIR,ERR])Age()optic.Optic[I,S,T,int,int,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Person)*int{
return &x.age
}))))))
}
func (s *oPerson[I,S,T,RET,RW,DIR,ERR])Some()*lPerson[optic.Void,mo.Option[Person],mo.Option[Person],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OPersonOf(optic.Some[Person]())
}
func (s *oPerson[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Person],mo.Option[Person],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Person]())))))
}
type lPersonHobby[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,PersonHobby,PersonHobby,RET,RW,DIR,ERR]
}
func (s *lPersonHobby[I,S,T,RET,RW,DIR,ERR])Id()optic.MakeLensRealOps[I,S,T,int,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *PersonHobby)*int{
return &x.id
})))))))
}
func (s *lPersonHobby[I,S,T,RET,RW,DIR,ERR])PersonId()optic.MakeLensRealOps[I,S,T,int,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *PersonHobby)*int{
return &x.personId
})))))))
}
func (s *lPersonHobby[I,S,T,RET,RW,DIR,ERR])HobbyId()optic.MakeLensRealOps[I,S,T,int,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *PersonHobby)*int{
return &x.hobbyId
})))))))
}
type sPersonHobby[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,PersonHobby, optic.Pure],optic.Collection[int,PersonHobby, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]PersonHobby,[]PersonHobby,RET,RW,DIR,ERR]
}
func (s *sPersonHobby[I,S,T,RET,RW,DIR,ERR])Traverse()*lPersonHobby[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPersonHobbyOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[PersonHobby]()))))))
}
func (s *sPersonHobby[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lPersonHobby[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPersonHobbyOf(optic.Index(s.Traverse(),index))
}
type mPersonHobby[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,PersonHobby, optic.Pure],optic.Collection[I,PersonHobby, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]PersonHobby,map[I]PersonHobby,RET,RW,DIR,ERR]
}
func (s *mPersonHobby[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lPersonHobby[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPersonHobbyOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,PersonHobby]()))))))
}
func (s *mPersonHobby[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lPersonHobby[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPersonHobbyOf(optic.Index(s.Traverse(),index))
}
type oPersonHobby[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*PersonHobby,*PersonHobby,RET,RW,DIR,ERR]
}
func (s *oPersonHobby[I,S,T,RET,RW,DIR,ERR])Id()optic.Optic[I,S,T,int,int,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *PersonHobby)*int{
return &x.id
}))))))
}
func (s *oPersonHobby[I,S,T,RET,RW,DIR,ERR])PersonId()optic.Optic[I,S,T,int,int,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *PersonHobby)*int{
return &x.personId
}))))))
}
func (s *oPersonHobby[I,S,T,RET,RW,DIR,ERR])HobbyId()optic.Optic[I,S,T,int,int,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *PersonHobby)*int{
return &x.hobbyId
}))))))
}
func (s *oPersonHobby[I,S,T,RET,RW,DIR,ERR])Some()*lPersonHobby[optic.Void,mo.Option[PersonHobby],mo.Option[PersonHobby],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OPersonHobbyOf(optic.Some[PersonHobby]())
}
func (s *oPersonHobby[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[PersonHobby],mo.Option[PersonHobby],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[PersonHobby]())))))
}
type ldb[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,db,db,RET,RW,DIR,ERR]
}
type sdb[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,db, optic.Pure],optic.Collection[int,db, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]db,[]db,RET,RW,DIR,ERR]
}
func (s *sdb[I,S,T,RET,RW,DIR,ERR])Traverse()*ldb[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OdbOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[db]()))))))
}
func (s *sdb[I,S,T,RET,RW,DIR,ERR])Nth(index int)*ldb[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OdbOf(optic.Index(s.Traverse(),index))
}
type mdb[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,db, optic.Pure],optic.Collection[I,db, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]db,map[I]db,RET,RW,DIR,ERR]
}
func (s *mdb[I,J,S,T,RET,RW,DIR,ERR])Traverse()*ldb[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OdbOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,db]()))))))
}
func (s *mdb[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*ldb[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OdbOf(optic.Index(s.Traverse(),index))
}
type odb[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*db,*db,RET,RW,DIR,ERR]
}
func (s *odb[I,S,T,RET,RW,DIR,ERR])Some()*ldb[optic.Void,mo.Option[db],mo.Option[db],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OdbOf(optic.Some[db]())
}
func (s *odb[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[db],mo.Option[db],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[db]())))))
}
type ldbHobby[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,dbHobby,dbHobby,RET,RW,DIR,ERR]
}
type sdbHobby[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,dbHobby, optic.Pure],optic.Collection[int,dbHobby, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]dbHobby,[]dbHobby,RET,RW,DIR,ERR]
}
func (s *sdbHobby[I,S,T,RET,RW,DIR,ERR])Traverse()*ldbHobby[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OdbHobbyOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[dbHobby]()))))))
}
func (s *sdbHobby[I,S,T,RET,RW,DIR,ERR])Nth(index int)*ldbHobby[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OdbHobbyOf(optic.Index(s.Traverse(),index))
}
type mdbHobby[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,dbHobby, optic.Pure],optic.Collection[I,dbHobby, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]dbHobby,map[I]dbHobby,RET,RW,DIR,ERR]
}
func (s *mdbHobby[I,J,S,T,RET,RW,DIR,ERR])Traverse()*ldbHobby[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OdbHobbyOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,dbHobby]()))))))
}
func (s *mdbHobby[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*ldbHobby[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OdbHobbyOf(optic.Index(s.Traverse(),index))
}
type odbHobby[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*dbHobby,*dbHobby,RET,RW,DIR,ERR]
}
func (s *odbHobby[I,S,T,RET,RW,DIR,ERR])Some()*ldbHobby[optic.Void,mo.Option[dbHobby],mo.Option[dbHobby],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OdbHobbyOf(optic.Some[dbHobby]())
}
func (s *odbHobby[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[dbHobby],mo.Option[dbHobby],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[dbHobby]())))))
}
type ldbPerson[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,dbPerson,dbPerson,RET,RW,DIR,ERR]
}
type sdbPerson[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,dbPerson, optic.Pure],optic.Collection[int,dbPerson, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]dbPerson,[]dbPerson,RET,RW,DIR,ERR]
}
func (s *sdbPerson[I,S,T,RET,RW,DIR,ERR])Traverse()*ldbPerson[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OdbPersonOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[dbPerson]()))))))
}
func (s *sdbPerson[I,S,T,RET,RW,DIR,ERR])Nth(index int)*ldbPerson[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OdbPersonOf(optic.Index(s.Traverse(),index))
}
type mdbPerson[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,dbPerson, optic.Pure],optic.Collection[I,dbPerson, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]dbPerson,map[I]dbPerson,RET,RW,DIR,ERR]
}
func (s *mdbPerson[I,J,S,T,RET,RW,DIR,ERR])Traverse()*ldbPerson[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OdbPersonOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,dbPerson]()))))))
}
func (s *mdbPerson[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*ldbPerson[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OdbPersonOf(optic.Index(s.Traverse(),index))
}
type odbPerson[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*dbPerson,*dbPerson,RET,RW,DIR,ERR]
}
func (s *odbPerson[I,S,T,RET,RW,DIR,ERR])Some()*ldbPerson[optic.Void,mo.Option[dbPerson],mo.Option[dbPerson],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OdbPersonOf(optic.Some[dbPerson]())
}
func (s *odbPerson[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[dbPerson],mo.Option[dbPerson],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[dbPerson]())))))
}
type o struct {
}
func (s *o)Hobby()*lHobby[optic.Void,Hobby,Hobby,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OHobbyOf[optic.Void,Hobby,Hobby,optic.ReturnOne](optic.Identity[Hobby]())
}
func (s *o)Person()*lPerson[optic.Void,Person,Person,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OPersonOf[optic.Void,Person,Person,optic.ReturnOne](optic.Identity[Person]())
}
func (s *o)PersonHobby()*lPersonHobby[optic.Void,PersonHobby,PersonHobby,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OPersonHobbyOf[optic.Void,PersonHobby,PersonHobby,optic.ReturnOne](optic.Identity[PersonHobby]())
}
func (s *o)db()*ldb[optic.Void,db,db,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OdbOf[optic.Void,db,db,optic.ReturnOne](optic.Identity[db]())
}
func (s *o)dbHobby()*ldbHobby[optic.Void,dbHobby,dbHobby,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OdbHobbyOf[optic.Void,dbHobby,dbHobby,optic.ReturnOne](optic.Identity[dbHobby]())
}
func (s *o)dbPerson()*ldbPerson[optic.Void,dbPerson,dbPerson,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OdbPersonOf[optic.Void,dbPerson,dbPerson,optic.ReturnOne](optic.Identity[dbPerson]())
}
func OHobbyOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Hobby,Hobby,RET,RW,DIR,ERR])*lHobby[I,S,T,RET,RW,DIR,ERR]{
return &lHobby[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func OPersonOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Person,Person,RET,RW,DIR,ERR])*lPerson[I,S,T,RET,RW,DIR,ERR]{
return &lPerson[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func OPersonHobbyOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,PersonHobby,PersonHobby,RET,RW,DIR,ERR])*lPersonHobby[I,S,T,RET,RW,DIR,ERR]{
return &lPersonHobby[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func OdbOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,db,db,RET,RW,DIR,ERR])*ldb[I,S,T,RET,RW,DIR,ERR]{
return &ldb[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func OdbHobbyOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,dbHobby,dbHobby,RET,RW,DIR,ERR])*ldbHobby[I,S,T,RET,RW,DIR,ERR]{
return &ldbHobby[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func OdbPersonOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,dbPerson,dbPerson,RET,RW,DIR,ERR])*ldbPerson[I,S,T,RET,RW,DIR,ERR]{
return &ldbPerson[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
var O  = o{
}
