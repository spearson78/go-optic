package data

import (
	"github.com/samber/mo"
	"github.com/spearson78/go-optic"
)
type lBlogPost[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,BlogPost,BlogPost,RET,RW,DIR,ERR]
}
func (s *lBlogPost[I,S,T,RET,RW,DIR,ERR])Content()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *BlogPost)*string{
return &x.Content
})))))))
}
func (s *lBlogPost[I,S,T,RET,RW,DIR,ERR])Comments()*sComment[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sComment[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *BlogPost)*[]Comment{
return &x.Comments
})))))),optic.SliceToCol[Comment]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *BlogPost)*[]Comment{
return &x.Comments
})))))),
}
}
func (s *lBlogPost[I,S,T,RET,RW,DIR,ERR])Ratings()*sRating[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sRating[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *BlogPost)*[]Rating{
return &x.Ratings
})))))),optic.SliceToCol[Rating]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *BlogPost)*[]Rating{
return &x.Ratings
})))))),
}
}
type sBlogPost[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,BlogPost, optic.Pure],optic.Collection[int,BlogPost, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]BlogPost,[]BlogPost,RET,RW,DIR,ERR]
}
func (s *sBlogPost[I,S,T,RET,RW,DIR,ERR])Traverse()*lBlogPost[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OBlogPostOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[BlogPost]()))))))
}
func (s *sBlogPost[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lBlogPost[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OBlogPostOf(optic.Index(s.Traverse(),index))
}
type mBlogPost[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,BlogPost, optic.Pure],optic.Collection[I,BlogPost, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]BlogPost,map[I]BlogPost,RET,RW,DIR,ERR]
}
func (s *mBlogPost[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lBlogPost[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OBlogPostOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,BlogPost]()))))))
}
func (s *mBlogPost[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lBlogPost[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OBlogPostOf(optic.Index(s.Traverse(),index))
}
type oBlogPost[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*BlogPost,*BlogPost,RET,RW,DIR,ERR]
}
func (s *oBlogPost[I,S,T,RET,RW,DIR,ERR])Content()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *BlogPost)*string{
return &x.Content
}))))))
}
func (s *oBlogPost[I,S,T,RET,RW,DIR,ERR])Comments()*sComment[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sComment[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *BlogPost)*[]Comment{
return &x.Comments
})))))),optic.SliceToCol[Comment]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *BlogPost)*[]Comment{
return &x.Comments
})))))),
}
}
func (s *oBlogPost[I,S,T,RET,RW,DIR,ERR])Ratings()*sRating[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sRating[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *BlogPost)*[]Rating{
return &x.Ratings
})))))),optic.SliceToCol[Rating]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *BlogPost)*[]Rating{
return &x.Ratings
})))))),
}
}
func (s *oBlogPost[I,S,T,RET,RW,DIR,ERR])Some()*lBlogPost[optic.Void,mo.Option[BlogPost],mo.Option[BlogPost],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OBlogPostOf(optic.Some[BlogPost]())
}
func (s *oBlogPost[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[BlogPost],mo.Option[BlogPost],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[BlogPost]())))))
}
type lComment[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Comment,Comment,RET,RW,DIR,ERR]
}
func (s *lComment[I,S,T,RET,RW,DIR,ERR])Title()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Comment)*string{
return &x.Title
})))))))
}
func (s *lComment[I,S,T,RET,RW,DIR,ERR])Content()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Comment)*string{
return &x.Content
})))))))
}
type sComment[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Comment, optic.Pure],optic.Collection[int,Comment, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Comment,[]Comment,RET,RW,DIR,ERR]
}
func (s *sComment[I,S,T,RET,RW,DIR,ERR])Traverse()*lComment[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OCommentOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Comment]()))))))
}
func (s *sComment[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lComment[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OCommentOf(optic.Index(s.Traverse(),index))
}
type mComment[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Comment, optic.Pure],optic.Collection[I,Comment, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Comment,map[I]Comment,RET,RW,DIR,ERR]
}
func (s *mComment[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lComment[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OCommentOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Comment]()))))))
}
func (s *mComment[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lComment[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OCommentOf(optic.Index(s.Traverse(),index))
}
type oComment[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Comment,*Comment,RET,RW,DIR,ERR]
}
func (s *oComment[I,S,T,RET,RW,DIR,ERR])Title()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Comment)*string{
return &x.Title
}))))))
}
func (s *oComment[I,S,T,RET,RW,DIR,ERR])Content()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Comment)*string{
return &x.Content
}))))))
}
func (s *oComment[I,S,T,RET,RW,DIR,ERR])Some()*lComment[optic.Void,mo.Option[Comment],mo.Option[Comment],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OCommentOf(optic.Some[Comment]())
}
func (s *oComment[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Comment],mo.Option[Comment],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Comment]())))))
}
type lRating[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Rating,Rating,RET,RW,DIR,ERR]
}
func (s *lRating[I,S,T,RET,RW,DIR,ERR])Author()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Rating)*string{
return &x.Author
})))))))
}
func (s *lRating[I,S,T,RET,RW,DIR,ERR])Stars()optic.MakeLensRealOps[I,S,T,int,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Rating)*int{
return &x.Stars
})))))))
}
type sRating[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Rating, optic.Pure],optic.Collection[int,Rating, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Rating,[]Rating,RET,RW,DIR,ERR]
}
func (s *sRating[I,S,T,RET,RW,DIR,ERR])Traverse()*lRating[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ORatingOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Rating]()))))))
}
func (s *sRating[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lRating[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ORatingOf(optic.Index(s.Traverse(),index))
}
type mRating[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Rating, optic.Pure],optic.Collection[I,Rating, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Rating,map[I]Rating,RET,RW,DIR,ERR]
}
func (s *mRating[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lRating[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ORatingOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Rating]()))))))
}
func (s *mRating[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lRating[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ORatingOf(optic.Index(s.Traverse(),index))
}
type oRating[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Rating,*Rating,RET,RW,DIR,ERR]
}
func (s *oRating[I,S,T,RET,RW,DIR,ERR])Author()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Rating)*string{
return &x.Author
}))))))
}
func (s *oRating[I,S,T,RET,RW,DIR,ERR])Stars()optic.Optic[I,S,T,int,int,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Rating)*int{
return &x.Stars
}))))))
}
func (s *oRating[I,S,T,RET,RW,DIR,ERR])Some()*lRating[optic.Void,mo.Option[Rating],mo.Option[Rating],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return ORatingOf(optic.Some[Rating]())
}
func (s *oRating[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Rating],mo.Option[Rating],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Rating]())))))
}
type o struct {
}
func (s *o)BlogPost()*lBlogPost[optic.Void,BlogPost,BlogPost,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OBlogPostOf[optic.Void,BlogPost,BlogPost,optic.ReturnOne](optic.Identity[BlogPost]())
}
func (s *o)Comment()*lComment[optic.Void,Comment,Comment,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OCommentOf[optic.Void,Comment,Comment,optic.ReturnOne](optic.Identity[Comment]())
}
func (s *o)Rating()*lRating[optic.Void,Rating,Rating,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return ORatingOf[optic.Void,Rating,Rating,optic.ReturnOne](optic.Identity[Rating]())
}
func OBlogPostOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,BlogPost,BlogPost,RET,RW,DIR,ERR])*lBlogPost[I,S,T,RET,RW,DIR,ERR]{
return &lBlogPost[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func OCommentOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Comment,Comment,RET,RW,DIR,ERR])*lComment[I,S,T,RET,RW,DIR,ERR]{
return &lComment[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func ORatingOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Rating,Rating,RET,RW,DIR,ERR])*lRating[I,S,T,RET,RW,DIR,ERR]{
return &lRating[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
var O  = o{
}
