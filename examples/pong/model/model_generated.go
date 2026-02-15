package model

import (
	"github.com/samber/mo"
	"github.com/spearson78/go-optic"
)
type lPoint[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Point,Point,RET,RW,DIR,ERR]
}
func (s *lPoint[I,S,T,RET,RW,DIR,ERR])X()optic.MakeLensRealOps[I,S,T,float64,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Point)*float64{
return &x.x
})))))))
}
func (s *lPoint[I,S,T,RET,RW,DIR,ERR])Y()optic.MakeLensRealOps[I,S,T,float64,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Point)*float64{
return &x.y
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
return &x.x
}))))))
}
func (s *oPoint[I,S,T,RET,RW,DIR,ERR])Y()optic.Optic[I,S,T,float64,float64,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Point)*float64{
return &x.y
}))))))
}
func (s *oPoint[I,S,T,RET,RW,DIR,ERR])Some()*lPoint[optic.Void,mo.Option[Point],mo.Option[Point],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OPointOf(optic.Some[Point]())
}
func (s *oPoint[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Point],mo.Option[Point],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Point]())))))
}
type lPong[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Pong,Pong,RET,RW,DIR,ERR]
}
func (s *lPong[I,S,T,RET,RW,DIR,ERR])BallPos()*lPoint[I,S,T,RET,RW,optic.UniDir,ERR]{
return &lPoint[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Pong)*Point{
return &x.ballPos
})))))),
}
}
func (s *lPong[I,S,T,RET,RW,DIR,ERR])BallSpeed()*lVector[I,S,T,RET,RW,optic.UniDir,ERR]{
return &lVector[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Pong)*Vector{
return &x.ballSpeed
})))))),
}
}
func (s *lPong[I,S,T,RET,RW,DIR,ERR])Paddle1()optic.MakeLensRealOps[I,S,T,float64,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Pong)*float64{
return &x.paddle1
})))))))
}
func (s *lPong[I,S,T,RET,RW,DIR,ERR])Paddle2()optic.MakeLensRealOps[I,S,T,float64,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Pong)*float64{
return &x.paddle2
})))))))
}
func (s *lPong[I,S,T,RET,RW,DIR,ERR])Score1()optic.MakeLensRealOps[I,S,T,int,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Pong)*int{
return &x.score1
})))))))
}
func (s *lPong[I,S,T,RET,RW,DIR,ERR])Score2()optic.MakeLensRealOps[I,S,T,int,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Pong)*int{
return &x.score2
})))))))
}
func (s *lPong[I,S,T,RET,RW,DIR,ERR])KeyUpPressed()optic.MakeLensCmpOps[I,S,T,bool,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensCmpOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Pong)*bool{
return &x.keyUpPressed
})))))))
}
func (s *lPong[I,S,T,RET,RW,DIR,ERR])KeyDownPressed()optic.MakeLensCmpOps[I,S,T,bool,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensCmpOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Pong)*bool{
return &x.keyDownPressed
})))))))
}
type sPong[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Pong, optic.Pure],optic.Collection[int,Pong, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Pong,[]Pong,RET,RW,DIR,ERR]
}
func (s *sPong[I,S,T,RET,RW,DIR,ERR])Traverse()*lPong[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPongOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Pong]()))))))
}
func (s *sPong[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lPong[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPongOf(optic.Index(s.Traverse(),index))
}
type mPong[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Pong, optic.Pure],optic.Collection[I,Pong, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Pong,map[I]Pong,RET,RW,DIR,ERR]
}
func (s *mPong[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lPong[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPongOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Pong]()))))))
}
func (s *mPong[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lPong[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OPongOf(optic.Index(s.Traverse(),index))
}
type oPong[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Pong,*Pong,RET,RW,DIR,ERR]
}
func (s *oPong[I,S,T,RET,RW,DIR,ERR])BallPos()*lPoint[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &lPoint[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Pong)*Point{
return &x.ballPos
})))))),
}
}
func (s *oPong[I,S,T,RET,RW,DIR,ERR])BallSpeed()*lVector[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &lVector[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Pong)*Vector{
return &x.ballSpeed
})))))),
}
}
func (s *oPong[I,S,T,RET,RW,DIR,ERR])Paddle1()optic.Optic[I,S,T,float64,float64,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Pong)*float64{
return &x.paddle1
}))))))
}
func (s *oPong[I,S,T,RET,RW,DIR,ERR])Paddle2()optic.Optic[I,S,T,float64,float64,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Pong)*float64{
return &x.paddle2
}))))))
}
func (s *oPong[I,S,T,RET,RW,DIR,ERR])Score1()optic.Optic[I,S,T,int,int,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Pong)*int{
return &x.score1
}))))))
}
func (s *oPong[I,S,T,RET,RW,DIR,ERR])Score2()optic.Optic[I,S,T,int,int,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Pong)*int{
return &x.score2
}))))))
}
func (s *oPong[I,S,T,RET,RW,DIR,ERR])KeyUpPressed()optic.Optic[I,S,T,bool,bool,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Pong)*bool{
return &x.keyUpPressed
}))))))
}
func (s *oPong[I,S,T,RET,RW,DIR,ERR])KeyDownPressed()optic.Optic[I,S,T,bool,bool,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Pong)*bool{
return &x.keyDownPressed
}))))))
}
func (s *oPong[I,S,T,RET,RW,DIR,ERR])Some()*lPong[optic.Void,mo.Option[Pong],mo.Option[Pong],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OPongOf(optic.Some[Pong]())
}
func (s *oPong[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Pong],mo.Option[Pong],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Pong]())))))
}
type lVector[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Vector,Vector,RET,RW,DIR,ERR]
}
func (s *lVector[I,S,T,RET,RW,DIR,ERR])U()optic.MakeLensRealOps[I,S,T,float64,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Vector)*float64{
return &x.u
})))))))
}
func (s *lVector[I,S,T,RET,RW,DIR,ERR])V()optic.MakeLensRealOps[I,S,T,float64,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Vector)*float64{
return &x.v
})))))))
}
type sVector[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Vector, optic.Pure],optic.Collection[int,Vector, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Vector,[]Vector,RET,RW,DIR,ERR]
}
func (s *sVector[I,S,T,RET,RW,DIR,ERR])Traverse()*lVector[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OVectorOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Vector]()))))))
}
func (s *sVector[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lVector[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OVectorOf(optic.Index(s.Traverse(),index))
}
type mVector[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Vector, optic.Pure],optic.Collection[I,Vector, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Vector,map[I]Vector,RET,RW,DIR,ERR]
}
func (s *mVector[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lVector[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OVectorOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Vector]()))))))
}
func (s *mVector[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lVector[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OVectorOf(optic.Index(s.Traverse(),index))
}
type oVector[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Vector,*Vector,RET,RW,DIR,ERR]
}
func (s *oVector[I,S,T,RET,RW,DIR,ERR])U()optic.Optic[I,S,T,float64,float64,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Vector)*float64{
return &x.u
}))))))
}
func (s *oVector[I,S,T,RET,RW,DIR,ERR])V()optic.Optic[I,S,T,float64,float64,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Vector)*float64{
return &x.v
}))))))
}
func (s *oVector[I,S,T,RET,RW,DIR,ERR])Some()*lVector[optic.Void,mo.Option[Vector],mo.Option[Vector],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OVectorOf(optic.Some[Vector]())
}
func (s *oVector[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Vector],mo.Option[Vector],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Vector]())))))
}
type o struct {
}
func (s *o)Point()*lPoint[optic.Void,Point,Point,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OPointOf[optic.Void,Point,Point,optic.ReturnOne](optic.Identity[Point]())
}
func (s *o)Pong()*lPong[optic.Void,Pong,Pong,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OPongOf[optic.Void,Pong,Pong,optic.ReturnOne](optic.Identity[Pong]())
}
func (s *o)Vector()*lVector[optic.Void,Vector,Vector,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OVectorOf[optic.Void,Vector,Vector,optic.ReturnOne](optic.Identity[Vector]())
}
func OPointOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Point,Point,RET,RW,DIR,ERR])*lPoint[I,S,T,RET,RW,DIR,ERR]{
return &lPoint[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func OPongOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Pong,Pong,RET,RW,DIR,ERR])*lPong[I,S,T,RET,RW,DIR,ERR]{
return &lPong[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func OVectorOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Vector,Vector,RET,RW,DIR,ERR])*lVector[I,S,T,RET,RW,DIR,ERR]{
return &lVector[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
var O  = o{
}
