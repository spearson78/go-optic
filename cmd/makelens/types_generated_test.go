package main_test

import (
	"github.com/samber/mo"
	"github.com/spearson78/go-optic"
)
type lDatabase[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Database,Database,RET,RW,DIR,ERR]
}
func (s *lDatabase[I,S,T,RET,RW,DIR,ERR])Drawings()*sDrawing[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sDrawing[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Database)*[]Drawing{
return &x.Drawings
})))))),optic.SliceToCol[Drawing]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Database)*[]Drawing{
return &x.Drawings
})))))),
}
}
type sDatabase[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Database, optic.Pure],optic.Collection[int,Database, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Database,[]Database,RET,RW,DIR,ERR]
}
func (s *sDatabase[I,S,T,RET,RW,DIR,ERR])Traverse()*lDatabase[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ODatabaseOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Database]()))))))
}
func (s *sDatabase[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lDatabase[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ODatabaseOf(optic.Index(s.Traverse(),index))
}
type mDatabase[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Database, optic.Pure],optic.Collection[I,Database, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Database,map[I]Database,RET,RW,DIR,ERR]
}
func (s *mDatabase[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lDatabase[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ODatabaseOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Database]()))))))
}
func (s *mDatabase[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lDatabase[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ODatabaseOf(optic.Index(s.Traverse(),index))
}
type oDatabase[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Database,*Database,RET,RW,DIR,ERR]
}
func (s *oDatabase[I,S,T,RET,RW,DIR,ERR])Drawings()*sDrawing[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sDrawing[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Database)*[]Drawing{
return &x.Drawings
})))))),optic.SliceToCol[Drawing]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Database)*[]Drawing{
return &x.Drawings
})))))),
}
}
func (s *oDatabase[I,S,T,RET,RW,DIR,ERR])Some()*lDatabase[optic.Void,mo.Option[Database],mo.Option[Database],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return ODatabaseOf(optic.Some[Database]())
}
func (s *oDatabase[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Database],mo.Option[Database],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Database]())))))
}
type lDrawing[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Drawing,Drawing,RET,RW,DIR,ERR]
}
func (s *lDrawing[I,S,T,RET,RW,DIR,ERR])Title()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Drawing)*string{
return &x.Title
})))))))
}
func (s *lDrawing[I,S,T,RET,RW,DIR,ERR])Pages()*sPage[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sPage[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Drawing)*[]Page{
return &x.Pages
})))))),optic.SliceToCol[Page]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Drawing)*[]Page{
return &x.Pages
})))))),
}
}
func (s *lDrawing[I,S,T,RET,RW,DIR,ERR])Meta()*oMetaData[I,S,T,RET,RW,optic.UniDir,ERR]{
return &oMetaData[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Drawing)**MetaData{
return &x.Meta
})))))),
}
}
type sDrawing[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Drawing, optic.Pure],optic.Collection[int,Drawing, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Drawing,[]Drawing,RET,RW,DIR,ERR]
}
func (s *sDrawing[I,S,T,RET,RW,DIR,ERR])Traverse()*lDrawing[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ODrawingOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Drawing]()))))))
}
func (s *sDrawing[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lDrawing[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ODrawingOf(optic.Index(s.Traverse(),index))
}
type mDrawing[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Drawing, optic.Pure],optic.Collection[I,Drawing, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Drawing,map[I]Drawing,RET,RW,DIR,ERR]
}
func (s *mDrawing[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lDrawing[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ODrawingOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Drawing]()))))))
}
func (s *mDrawing[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lDrawing[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ODrawingOf(optic.Index(s.Traverse(),index))
}
type oDrawing[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Drawing,*Drawing,RET,RW,DIR,ERR]
}
func (s *oDrawing[I,S,T,RET,RW,DIR,ERR])Title()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Drawing)*string{
return &x.Title
}))))))
}
func (s *oDrawing[I,S,T,RET,RW,DIR,ERR])Pages()*sPage[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sPage[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Drawing)*[]Page{
return &x.Pages
})))))),optic.SliceToCol[Page]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Drawing)*[]Page{
return &x.Pages
})))))),
}
}
func (s *oDrawing[I,S,T,RET,RW,DIR,ERR])Meta()*oMetaData[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &oMetaData[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Drawing)**MetaData{
return &x.Meta
})))))),
}
}
func (s *oDrawing[I,S,T,RET,RW,DIR,ERR])Some()*lDrawing[optic.Void,mo.Option[Drawing],mo.Option[Drawing],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return ODrawingOf(optic.Some[Drawing]())
}
func (s *oDrawing[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Drawing],mo.Option[Drawing],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Drawing]())))))
}
type lLine[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Line,Line,RET,RW,DIR,ERR]
}
func (s *lLine[I,S,T,RET,RW,DIR,ERR])Start()*lPoint[I,S,T,RET,RW,optic.UniDir,ERR]{
return &lPoint[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Line)*Point{
return &x.Start
})))))),
}
}
func (s *lLine[I,S,T,RET,RW,DIR,ERR])End()*lPoint[I,S,T,RET,RW,optic.UniDir,ERR]{
return &lPoint[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Line)*Point{
return &x.End
})))))),
}
}
type sLine[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Line, optic.Pure],optic.Collection[int,Line, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Line,[]Line,RET,RW,DIR,ERR]
}
func (s *sLine[I,S,T,RET,RW,DIR,ERR])Traverse()*lLine[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OLineOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Line]()))))))
}
func (s *sLine[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lLine[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OLineOf(optic.Index(s.Traverse(),index))
}
type mLine[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Line, optic.Pure],optic.Collection[I,Line, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Line,map[I]Line,RET,RW,DIR,ERR]
}
func (s *mLine[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lLine[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OLineOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Line]()))))))
}
func (s *mLine[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lLine[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OLineOf(optic.Index(s.Traverse(),index))
}
type oLine[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Line,*Line,RET,RW,DIR,ERR]
}
func (s *oLine[I,S,T,RET,RW,DIR,ERR])Start()*lPoint[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &lPoint[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Line)*Point{
return &x.Start
})))))),
}
}
func (s *oLine[I,S,T,RET,RW,DIR,ERR])End()*lPoint[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &lPoint[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Line)*Point{
return &x.End
})))))),
}
}
func (s *oLine[I,S,T,RET,RW,DIR,ERR])Some()*lLine[optic.Void,mo.Option[Line],mo.Option[Line],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OLineOf(optic.Some[Line]())
}
func (s *oLine[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Line],mo.Option[Line],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Line]())))))
}
type lMetaData[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,MetaData,MetaData,RET,RW,DIR,ERR]
}
func (s *lMetaData[I,S,T,RET,RW,DIR,ERR])Author()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *MetaData)*string{
return &x.Author
})))))))
}
func (s *lMetaData[I,S,T,RET,RW,DIR,ERR])Date()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *MetaData)*string{
return &x.Date
})))))))
}
func (s *lMetaData[I,S,T,RET,RW,DIR,ERR])Params()optic.MakeLensMap[string,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensMap(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s,optic.FieldLens(func(x *MetaData)*map[string]string{
return &x.Params
})))))))
}
type sMetaData[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,MetaData, optic.Pure],optic.Collection[int,MetaData, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]MetaData,[]MetaData,RET,RW,DIR,ERR]
}
func (s *sMetaData[I,S,T,RET,RW,DIR,ERR])Traverse()*lMetaData[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OMetaDataOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[MetaData]()))))))
}
func (s *sMetaData[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lMetaData[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OMetaDataOf(optic.Index(s.Traverse(),index))
}
type mMetaData[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,MetaData, optic.Pure],optic.Collection[I,MetaData, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]MetaData,map[I]MetaData,RET,RW,DIR,ERR]
}
func (s *mMetaData[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lMetaData[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OMetaDataOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,MetaData]()))))))
}
func (s *mMetaData[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lMetaData[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OMetaDataOf(optic.Index(s.Traverse(),index))
}
type oMetaData[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*MetaData,*MetaData,RET,RW,DIR,ERR]
}
func (s *oMetaData[I,S,T,RET,RW,DIR,ERR])Author()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *MetaData)*string{
return &x.Author
}))))))
}
func (s *oMetaData[I,S,T,RET,RW,DIR,ERR])Date()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *MetaData)*string{
return &x.Date
}))))))
}
func (s *oMetaData[I,S,T,RET,RW,DIR,ERR])Params()optic.MakeLensMap[string,S,T,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensMap(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s,optic.PtrFieldLens(func(x *MetaData)*map[string]string{
return &x.Params
})))))))
}
func (s *oMetaData[I,S,T,RET,RW,DIR,ERR])Some()*lMetaData[optic.Void,mo.Option[MetaData],mo.Option[MetaData],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OMetaDataOf(optic.Some[MetaData]())
}
func (s *oMetaData[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[MetaData],mo.Option[MetaData],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[MetaData]())))))
}
type lPage[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Page,Page,RET,RW,DIR,ERR]
}
func (s *lPage[I,S,T,RET,RW,DIR,ERR])Title()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Page)*string{
return &x.Title
})))))))
}
func (s *lPage[I,S,T,RET,RW,DIR,ERR])Lines()*sLine[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sLine[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Page)*[]Line{
return &x.Lines
})))))),optic.SliceToCol[Line]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Page)*[]Line{
return &x.Lines
})))))),
}
}
func (s *lPage[I,S,T,RET,RW,DIR,ERR])Css()optic.MakeLensSlice[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Page)*[]string{
return &x.Css
})))))))
}
type sPage[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Page, optic.Pure],optic.Collection[int,Page, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Page,[]Page,RET,RW,DIR,ERR]
}
func (s *sPage[I,S,T,RET,RW,DIR,ERR])Traverse()*lPage[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPageOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Page]()))))))
}
func (s *sPage[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lPage[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPageOf(optic.Index(s.Traverse(),index))
}
type mPage[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Page, optic.Pure],optic.Collection[I,Page, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Page,map[I]Page,RET,RW,DIR,ERR]
}
func (s *mPage[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lPage[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPageOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Page]()))))))
}
func (s *mPage[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lPage[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPageOf(optic.Index(s.Traverse(),index))
}
type oPage[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Page,*Page,RET,RW,DIR,ERR]
}
func (s *oPage[I,S,T,RET,RW,DIR,ERR])Title()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Page)*string{
return &x.Title
}))))))
}
func (s *oPage[I,S,T,RET,RW,DIR,ERR])Lines()*sLine[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sLine[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Page)*[]Line{
return &x.Lines
})))))),optic.SliceToCol[Line]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Page)*[]Line{
return &x.Lines
})))))),
}
}
func (s *oPage[I,S,T,RET,RW,DIR,ERR])Css()optic.MakeLensSlice[I,S,T,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Page)*[]string{
return &x.Css
})))))))
}
func (s *oPage[I,S,T,RET,RW,DIR,ERR])Some()*lPage[optic.Void,mo.Option[Page],mo.Option[Page],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OPageOf(optic.Some[Page]())
}
func (s *oPage[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Page],mo.Option[Page],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Page]())))))
}
type lPoint[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Point,Point,RET,RW,DIR,ERR]
}
func (s *lPoint[I,S,T,RET,RW,DIR,ERR])X()optic.MakeLensRealOps[I,S,T,float64,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Point)*float64{
return &x.X
})))))))
}
func (s *lPoint[I,S,T,RET,RW,DIR,ERR])Y()optic.MakeLensRealOps[I,S,T,float64,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Point)*float64{
return &x.Y
})))))))
}
type sPoint[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Point, optic.Pure],optic.Collection[int,Point, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Point,[]Point,RET,RW,DIR,ERR]
}
func (s *sPoint[I,S,T,RET,RW,DIR,ERR])Traverse()*lPoint[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPointOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Point]()))))))
}
func (s *sPoint[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lPoint[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPointOf(optic.Index(s.Traverse(),index))
}
type mPoint[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Point, optic.Pure],optic.Collection[I,Point, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Point,map[I]Point,RET,RW,DIR,ERR]
}
func (s *mPoint[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lPoint[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPointOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Point]()))))))
}
func (s *mPoint[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lPoint[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPointOf(optic.Index(s.Traverse(),index))
}
type oPoint[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Point,*Point,RET,RW,DIR,ERR]
}
func (s *oPoint[I,S,T,RET,RW,DIR,ERR])X()optic.Optic[I,S,T,float64,float64,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Point)*float64{
return &x.X
}))))))
}
func (s *oPoint[I,S,T,RET,RW,DIR,ERR])Y()optic.Optic[I,S,T,float64,float64,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Point)*float64{
return &x.Y
}))))))
}
func (s *oPoint[I,S,T,RET,RW,DIR,ERR])Some()*lPoint[optic.Void,mo.Option[Point],mo.Option[Point],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OPointOf(optic.Some[Point]())
}
func (s *oPoint[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Point],mo.Option[Point],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Point]())))))
}
type o struct {
}
func (s *o)Database()*lDatabase[optic.Void,Database,Database,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return ODatabaseOf[optic.Void,Database,Database,optic.ReturnOne](optic.Identity[Database]())
}
func (s *o)Drawing()*lDrawing[optic.Void,Drawing,Drawing,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return ODrawingOf[optic.Void,Drawing,Drawing,optic.ReturnOne](optic.Identity[Drawing]())
}
func (s *o)Line()*lLine[optic.Void,Line,Line,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OLineOf[optic.Void,Line,Line,optic.ReturnOne](optic.Identity[Line]())
}
func (s *o)MetaData()*lMetaData[optic.Void,MetaData,MetaData,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OMetaDataOf[optic.Void,MetaData,MetaData,optic.ReturnOne](optic.Identity[MetaData]())
}
func (s *o)Page()*lPage[optic.Void,Page,Page,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OPageOf[optic.Void,Page,Page,optic.ReturnOne](optic.Identity[Page]())
}
func (s *o)Point()*lPoint[optic.Void,Point,Point,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OPointOf[optic.Void,Point,Point,optic.ReturnOne](optic.Identity[Point]())
}
func ODatabaseOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Database,Database,RET,RW,DIR,ERR])*lDatabase[I,S,T,RET,RW,DIR,ERR]{
return &lDatabase[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func ODrawingOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Drawing,Drawing,RET,RW,DIR,ERR])*lDrawing[I,S,T,RET,RW,DIR,ERR]{
return &lDrawing[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func OLineOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Line,Line,RET,RW,DIR,ERR])*lLine[I,S,T,RET,RW,DIR,ERR]{
return &lLine[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func OMetaDataOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,MetaData,MetaData,RET,RW,DIR,ERR])*lMetaData[I,S,T,RET,RW,DIR,ERR]{
return &lMetaData[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func OPageOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Page,Page,RET,RW,DIR,ERR])*lPage[I,S,T,RET,RW,DIR,ERR]{
return &lPage[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func OPointOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Point,Point,RET,RW,DIR,ERR])*lPoint[I,S,T,RET,RW,DIR,ERR]{
return &lPoint[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
var O  = o{
}
