package codegen

import (
	"github.com/samber/mo"
	"github.com/spearson78/go-optic"
)
type lAddressExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,AddressExpr,AddressExpr,RET,RW,DIR,ERR]
}
func (s *lAddressExpr[I,S,T,RET,RW,DIR,ERR])Target()optic.Optic[I,S,T,Expression,Expression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *AddressExpr)*Expression{
return &x.Target
}))))))
}
type sAddressExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,AddressExpr, optic.Pure],optic.Collection[int,AddressExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]AddressExpr,[]AddressExpr,RET,RW,DIR,ERR]
}
func (s *sAddressExpr[I,S,T,RET,RW,DIR,ERR])Traverse()*lAddressExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGAddressExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[AddressExpr]()))))))
}
func (s *sAddressExpr[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lAddressExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGAddressExprOf(optic.Index(s.Traverse(),index))
}
type mAddressExpr[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,AddressExpr, optic.Pure],optic.Collection[I,AddressExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]AddressExpr,map[I]AddressExpr,RET,RW,DIR,ERR]
}
func (s *mAddressExpr[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lAddressExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGAddressExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,AddressExpr]()))))))
}
func (s *mAddressExpr[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lAddressExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGAddressExprOf(optic.Index(s.Traverse(),index))
}
type oAddressExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*AddressExpr,*AddressExpr,RET,RW,DIR,ERR]
}
func (s *oAddressExpr[I,S,T,RET,RW,DIR,ERR])Target()optic.Optic[I,S,T,Expression,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *AddressExpr)*Expression{
return &x.Target
}))))))
}
func (s *oAddressExpr[I,S,T,RET,RW,DIR,ERR])Some()*lAddressExpr[optic.Void,mo.Option[AddressExpr],mo.Option[AddressExpr],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGAddressExprOf(optic.Some[AddressExpr]())
}
func (s *oAddressExpr[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[AddressExpr],mo.Option[AddressExpr],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[AddressExpr]())))))
}
type lAssignField[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,AssignField,AssignField,RET,RW,DIR,ERR]
}
func (s *lAssignField[I,S,T,RET,RW,DIR,ERR])Name()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *AssignField)*string{
return &x.Name
})))))))
}
func (s *lAssignField[I,S,T,RET,RW,DIR,ERR])Value()optic.Optic[I,S,T,Expression,Expression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *AssignField)*Expression{
return &x.Value
}))))))
}
type sAssignField[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,AssignField, optic.Pure],optic.Collection[int,AssignField, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]AssignField,[]AssignField,RET,RW,DIR,ERR]
}
func (s *sAssignField[I,S,T,RET,RW,DIR,ERR])Traverse()*lAssignField[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGAssignFieldOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[AssignField]()))))))
}
func (s *sAssignField[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lAssignField[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGAssignFieldOf(optic.Index(s.Traverse(),index))
}
type mAssignField[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,AssignField, optic.Pure],optic.Collection[I,AssignField, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]AssignField,map[I]AssignField,RET,RW,DIR,ERR]
}
func (s *mAssignField[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lAssignField[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGAssignFieldOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,AssignField]()))))))
}
func (s *mAssignField[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lAssignField[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGAssignFieldOf(optic.Index(s.Traverse(),index))
}
type oAssignField[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*AssignField,*AssignField,RET,RW,DIR,ERR]
}
func (s *oAssignField[I,S,T,RET,RW,DIR,ERR])Name()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *AssignField)*string{
return &x.Name
}))))))
}
func (s *oAssignField[I,S,T,RET,RW,DIR,ERR])Value()optic.Optic[I,S,T,Expression,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *AssignField)*Expression{
return &x.Value
}))))))
}
func (s *oAssignField[I,S,T,RET,RW,DIR,ERR])Some()*lAssignField[optic.Void,mo.Option[AssignField],mo.Option[AssignField],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGAssignFieldOf(optic.Some[AssignField]())
}
func (s *oAssignField[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[AssignField],mo.Option[AssignField],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[AssignField]())))))
}
type lAssignVar[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,AssignVar,AssignVar,RET,RW,DIR,ERR]
}
func (s *lAssignVar[I,S,T,RET,RW,DIR,ERR])Declare()optic.MakeLensCmpOps[I,S,T,bool,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensCmpOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *AssignVar)*bool{
return &x.Declare
})))))))
}
func (s *lAssignVar[I,S,T,RET,RW,DIR,ERR])Vars()optic.MakeLensSlice[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *AssignVar)*[]string{
return &x.Vars
})))))))
}
func (s *lAssignVar[I,S,T,RET,RW,DIR,ERR])Value()optic.Optic[I,S,T,Expression,Expression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *AssignVar)*Expression{
return &x.Value
}))))))
}
type sAssignVar[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,AssignVar, optic.Pure],optic.Collection[int,AssignVar, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]AssignVar,[]AssignVar,RET,RW,DIR,ERR]
}
func (s *sAssignVar[I,S,T,RET,RW,DIR,ERR])Traverse()*lAssignVar[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGAssignVarOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[AssignVar]()))))))
}
func (s *sAssignVar[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lAssignVar[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGAssignVarOf(optic.Index(s.Traverse(),index))
}
type mAssignVar[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,AssignVar, optic.Pure],optic.Collection[I,AssignVar, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]AssignVar,map[I]AssignVar,RET,RW,DIR,ERR]
}
func (s *mAssignVar[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lAssignVar[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGAssignVarOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,AssignVar]()))))))
}
func (s *mAssignVar[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lAssignVar[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGAssignVarOf(optic.Index(s.Traverse(),index))
}
type oAssignVar[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*AssignVar,*AssignVar,RET,RW,DIR,ERR]
}
func (s *oAssignVar[I,S,T,RET,RW,DIR,ERR])Declare()optic.Optic[I,S,T,bool,bool,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *AssignVar)*bool{
return &x.Declare
}))))))
}
func (s *oAssignVar[I,S,T,RET,RW,DIR,ERR])Vars()optic.MakeLensSlice[I,S,T,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *AssignVar)*[]string{
return &x.Vars
})))))))
}
func (s *oAssignVar[I,S,T,RET,RW,DIR,ERR])Value()optic.Optic[I,S,T,Expression,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *AssignVar)*Expression{
return &x.Value
}))))))
}
func (s *oAssignVar[I,S,T,RET,RW,DIR,ERR])Some()*lAssignVar[optic.Void,mo.Option[AssignVar],mo.Option[AssignVar],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGAssignVarOf(optic.Some[AssignVar]())
}
func (s *oAssignVar[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[AssignVar],mo.Option[AssignVar],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[AssignVar]())))))
}
type lBinaryExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,BinaryExpr,BinaryExpr,RET,RW,DIR,ERR]
}
func (s *lBinaryExpr[I,S,T,RET,RW,DIR,ERR])Op()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *BinaryExpr)*string{
return &x.Op
})))))))
}
func (s *lBinaryExpr[I,S,T,RET,RW,DIR,ERR])Left()optic.Optic[I,S,T,Expression,Expression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *BinaryExpr)*Expression{
return &x.Left
}))))))
}
func (s *lBinaryExpr[I,S,T,RET,RW,DIR,ERR])Right()optic.Optic[I,S,T,Expression,Expression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *BinaryExpr)*Expression{
return &x.Right
}))))))
}
type sBinaryExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,BinaryExpr, optic.Pure],optic.Collection[int,BinaryExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]BinaryExpr,[]BinaryExpr,RET,RW,DIR,ERR]
}
func (s *sBinaryExpr[I,S,T,RET,RW,DIR,ERR])Traverse()*lBinaryExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGBinaryExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[BinaryExpr]()))))))
}
func (s *sBinaryExpr[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lBinaryExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGBinaryExprOf(optic.Index(s.Traverse(),index))
}
type mBinaryExpr[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,BinaryExpr, optic.Pure],optic.Collection[I,BinaryExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]BinaryExpr,map[I]BinaryExpr,RET,RW,DIR,ERR]
}
func (s *mBinaryExpr[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lBinaryExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGBinaryExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,BinaryExpr]()))))))
}
func (s *mBinaryExpr[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lBinaryExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGBinaryExprOf(optic.Index(s.Traverse(),index))
}
type oBinaryExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*BinaryExpr,*BinaryExpr,RET,RW,DIR,ERR]
}
func (s *oBinaryExpr[I,S,T,RET,RW,DIR,ERR])Op()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *BinaryExpr)*string{
return &x.Op
}))))))
}
func (s *oBinaryExpr[I,S,T,RET,RW,DIR,ERR])Left()optic.Optic[I,S,T,Expression,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *BinaryExpr)*Expression{
return &x.Left
}))))))
}
func (s *oBinaryExpr[I,S,T,RET,RW,DIR,ERR])Right()optic.Optic[I,S,T,Expression,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *BinaryExpr)*Expression{
return &x.Right
}))))))
}
func (s *oBinaryExpr[I,S,T,RET,RW,DIR,ERR])Some()*lBinaryExpr[optic.Void,mo.Option[BinaryExpr],mo.Option[BinaryExpr],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGBinaryExprOf(optic.Some[BinaryExpr]())
}
func (s *oBinaryExpr[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[BinaryExpr],mo.Option[BinaryExpr],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[BinaryExpr]())))))
}
type lCallExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,CallExpr,CallExpr,RET,RW,DIR,ERR]
}
func (s *lCallExpr[I,S,T,RET,RW,DIR,ERR])Func()optic.Optic[I,S,T,Expression,Expression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *CallExpr)*Expression{
return &x.Func
}))))))
}
func (s *lCallExpr[I,S,T,RET,RW,DIR,ERR])TypeParams()optic.MakeLensSlice[I,S,T,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *CallExpr)*[]TypeExpression{
return &x.TypeParams
})))))))
}
func (s *lCallExpr[I,S,T,RET,RW,DIR,ERR])Params()optic.MakeLensSlice[I,S,T,Expression,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *CallExpr)*[]Expression{
return &x.Params
})))))))
}
type sCallExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,CallExpr, optic.Pure],optic.Collection[int,CallExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]CallExpr,[]CallExpr,RET,RW,DIR,ERR]
}
func (s *sCallExpr[I,S,T,RET,RW,DIR,ERR])Traverse()*lCallExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGCallExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[CallExpr]()))))))
}
func (s *sCallExpr[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lCallExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGCallExprOf(optic.Index(s.Traverse(),index))
}
type mCallExpr[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,CallExpr, optic.Pure],optic.Collection[I,CallExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]CallExpr,map[I]CallExpr,RET,RW,DIR,ERR]
}
func (s *mCallExpr[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lCallExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGCallExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,CallExpr]()))))))
}
func (s *mCallExpr[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lCallExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGCallExprOf(optic.Index(s.Traverse(),index))
}
type oCallExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*CallExpr,*CallExpr,RET,RW,DIR,ERR]
}
func (s *oCallExpr[I,S,T,RET,RW,DIR,ERR])Func()optic.Optic[I,S,T,Expression,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *CallExpr)*Expression{
return &x.Func
}))))))
}
func (s *oCallExpr[I,S,T,RET,RW,DIR,ERR])TypeParams()optic.MakeLensSlice[I,S,T,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *CallExpr)*[]TypeExpression{
return &x.TypeParams
})))))))
}
func (s *oCallExpr[I,S,T,RET,RW,DIR,ERR])Params()optic.MakeLensSlice[I,S,T,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *CallExpr)*[]Expression{
return &x.Params
})))))))
}
func (s *oCallExpr[I,S,T,RET,RW,DIR,ERR])Some()*lCallExpr[optic.Void,mo.Option[CallExpr],mo.Option[CallExpr],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGCallExprOf(optic.Some[CallExpr]())
}
func (s *oCallExpr[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[CallExpr],mo.Option[CallExpr],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[CallExpr]())))))
}
type lDeRef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,DeRef,DeRef,RET,RW,DIR,ERR]
}
func (s *lDeRef[I,S,T,RET,RW,DIR,ERR])Value()optic.Optic[I,S,T,Expression,Expression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *DeRef)*Expression{
return &x.Value
}))))))
}
type sDeRef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,DeRef, optic.Pure],optic.Collection[int,DeRef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]DeRef,[]DeRef,RET,RW,DIR,ERR]
}
func (s *sDeRef[I,S,T,RET,RW,DIR,ERR])Traverse()*lDeRef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGDeRefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[DeRef]()))))))
}
func (s *sDeRef[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lDeRef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGDeRefOf(optic.Index(s.Traverse(),index))
}
type mDeRef[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,DeRef, optic.Pure],optic.Collection[I,DeRef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]DeRef,map[I]DeRef,RET,RW,DIR,ERR]
}
func (s *mDeRef[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lDeRef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGDeRefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,DeRef]()))))))
}
func (s *mDeRef[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lDeRef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGDeRefOf(optic.Index(s.Traverse(),index))
}
type oDeRef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*DeRef,*DeRef,RET,RW,DIR,ERR]
}
func (s *oDeRef[I,S,T,RET,RW,DIR,ERR])Value()optic.Optic[I,S,T,Expression,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *DeRef)*Expression{
return &x.Value
}))))))
}
func (s *oDeRef[I,S,T,RET,RW,DIR,ERR])Some()*lDeRef[optic.Void,mo.Option[DeRef],mo.Option[DeRef],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGDeRefOf(optic.Some[DeRef]())
}
func (s *oDeRef[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[DeRef],mo.Option[DeRef],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[DeRef]())))))
}
type lFieldDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,FieldDef,FieldDef,RET,RW,DIR,ERR]
}
func (s *lFieldDef[I,S,T,RET,RW,DIR,ERR])Name()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FieldDef)*string{
return &x.Name
})))))))
}
func (s *lFieldDef[I,S,T,RET,RW,DIR,ERR])Type()optic.Optic[I,S,T,TypeExpression,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FieldDef)*TypeExpression{
return &x.Type
}))))))
}
type sFieldDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,FieldDef, optic.Pure],optic.Collection[int,FieldDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]FieldDef,[]FieldDef,RET,RW,DIR,ERR]
}
func (s *sFieldDef[I,S,T,RET,RW,DIR,ERR])Traverse()*lFieldDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFieldDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[FieldDef]()))))))
}
func (s *sFieldDef[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lFieldDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFieldDefOf(optic.Index(s.Traverse(),index))
}
type mFieldDef[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,FieldDef, optic.Pure],optic.Collection[I,FieldDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]FieldDef,map[I]FieldDef,RET,RW,DIR,ERR]
}
func (s *mFieldDef[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lFieldDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFieldDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,FieldDef]()))))))
}
func (s *mFieldDef[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lFieldDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFieldDefOf(optic.Index(s.Traverse(),index))
}
type oFieldDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*FieldDef,*FieldDef,RET,RW,DIR,ERR]
}
func (s *oFieldDef[I,S,T,RET,RW,DIR,ERR])Name()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FieldDef)*string{
return &x.Name
}))))))
}
func (s *oFieldDef[I,S,T,RET,RW,DIR,ERR])Type()optic.Optic[I,S,T,TypeExpression,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FieldDef)*TypeExpression{
return &x.Type
}))))))
}
func (s *oFieldDef[I,S,T,RET,RW,DIR,ERR])Some()*lFieldDef[optic.Void,mo.Option[FieldDef],mo.Option[FieldDef],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGFieldDefOf(optic.Some[FieldDef]())
}
func (s *oFieldDef[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[FieldDef],mo.Option[FieldDef],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[FieldDef]())))))
}
type lFileDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,FileDef,FileDef,RET,RW,DIR,ERR]
}
func (s *lFileDef[I,S,T,RET,RW,DIR,ERR])Package()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FileDef)*string{
return &x.Package
})))))))
}
func (s *lFileDef[I,S,T,RET,RW,DIR,ERR])Imports()optic.MakeLensSlice[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FileDef)*[]string{
return &x.Imports
})))))))
}
func (s *lFileDef[I,S,T,RET,RW,DIR,ERR])DotImports()optic.MakeLensSlice[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FileDef)*[]string{
return &x.DotImports
})))))))
}
func (s *lFileDef[I,S,T,RET,RW,DIR,ERR])Structs()*sStructDef[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sStructDef[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FileDef)*[]StructDef{
return &x.Structs
})))))),optic.SliceToCol[StructDef]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FileDef)*[]StructDef{
return &x.Structs
})))))),
}
}
func (s *lFileDef[I,S,T,RET,RW,DIR,ERR])Vars()*sVarDef[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sVarDef[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FileDef)*[]VarDef{
return &x.Vars
})))))),optic.SliceToCol[VarDef]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FileDef)*[]VarDef{
return &x.Vars
})))))),
}
}
func (s *lFileDef[I,S,T,RET,RW,DIR,ERR])Funcs()*sFuncDef[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sFuncDef[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FileDef)*[]FuncDef{
return &x.Funcs
})))))),optic.SliceToCol[FuncDef]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FileDef)*[]FuncDef{
return &x.Funcs
})))))),
}
}
type sFileDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,FileDef, optic.Pure],optic.Collection[int,FileDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]FileDef,[]FileDef,RET,RW,DIR,ERR]
}
func (s *sFileDef[I,S,T,RET,RW,DIR,ERR])Traverse()*lFileDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFileDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[FileDef]()))))))
}
func (s *sFileDef[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lFileDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFileDefOf(optic.Index(s.Traverse(),index))
}
type mFileDef[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,FileDef, optic.Pure],optic.Collection[I,FileDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]FileDef,map[I]FileDef,RET,RW,DIR,ERR]
}
func (s *mFileDef[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lFileDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFileDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,FileDef]()))))))
}
func (s *mFileDef[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lFileDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFileDefOf(optic.Index(s.Traverse(),index))
}
type oFileDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*FileDef,*FileDef,RET,RW,DIR,ERR]
}
func (s *oFileDef[I,S,T,RET,RW,DIR,ERR])Package()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FileDef)*string{
return &x.Package
}))))))
}
func (s *oFileDef[I,S,T,RET,RW,DIR,ERR])Imports()optic.MakeLensSlice[I,S,T,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FileDef)*[]string{
return &x.Imports
})))))))
}
func (s *oFileDef[I,S,T,RET,RW,DIR,ERR])DotImports()optic.MakeLensSlice[I,S,T,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FileDef)*[]string{
return &x.DotImports
})))))))
}
func (s *oFileDef[I,S,T,RET,RW,DIR,ERR])Structs()*sStructDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sStructDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FileDef)*[]StructDef{
return &x.Structs
})))))),optic.SliceToCol[StructDef]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FileDef)*[]StructDef{
return &x.Structs
})))))),
}
}
func (s *oFileDef[I,S,T,RET,RW,DIR,ERR])Vars()*sVarDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sVarDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FileDef)*[]VarDef{
return &x.Vars
})))))),optic.SliceToCol[VarDef]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FileDef)*[]VarDef{
return &x.Vars
})))))),
}
}
func (s *oFileDef[I,S,T,RET,RW,DIR,ERR])Funcs()*sFuncDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sFuncDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FileDef)*[]FuncDef{
return &x.Funcs
})))))),optic.SliceToCol[FuncDef]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FileDef)*[]FuncDef{
return &x.Funcs
})))))),
}
}
func (s *oFileDef[I,S,T,RET,RW,DIR,ERR])Some()*lFileDef[optic.Void,mo.Option[FileDef],mo.Option[FileDef],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGFileDefOf(optic.Some[FileDef]())
}
func (s *oFileDef[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[FileDef],mo.Option[FileDef],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[FileDef]())))))
}
type lFuncDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,FuncDef,FuncDef,RET,RW,DIR,ERR]
}
func (s *lFuncDef[I,S,T,RET,RW,DIR,ERR])Docs()optic.MakeLensSlice[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FuncDef)*[]string{
return &x.Docs
})))))))
}
func (s *lFuncDef[I,S,T,RET,RW,DIR,ERR])Name()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FuncDef)*string{
return &x.Name
})))))))
}
func (s *lFuncDef[I,S,T,RET,RW,DIR,ERR])TypeParams()optic.MakeLensSlice[I,S,T,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FuncDef)*[]TypeExpression{
return &x.TypeParams
})))))))
}
func (s *lFuncDef[I,S,T,RET,RW,DIR,ERR])Params()*sParam[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sParam[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FuncDef)*[]Param{
return &x.Params
})))))),optic.SliceToCol[Param]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FuncDef)*[]Param{
return &x.Params
})))))),
}
}
func (s *lFuncDef[I,S,T,RET,RW,DIR,ERR])ReturnTypes()optic.MakeLensSlice[I,S,T,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FuncDef)*[]TypeExpression{
return &x.ReturnTypes
})))))))
}
func (s *lFuncDef[I,S,T,RET,RW,DIR,ERR])Body()optic.MakeLensSlice[I,S,T,Statement,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FuncDef)*[]Statement{
return &x.Body
})))))))
}
type sFuncDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,FuncDef, optic.Pure],optic.Collection[int,FuncDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]FuncDef,[]FuncDef,RET,RW,DIR,ERR]
}
func (s *sFuncDef[I,S,T,RET,RW,DIR,ERR])Traverse()*lFuncDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFuncDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[FuncDef]()))))))
}
func (s *sFuncDef[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lFuncDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFuncDefOf(optic.Index(s.Traverse(),index))
}
type mFuncDef[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,FuncDef, optic.Pure],optic.Collection[I,FuncDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]FuncDef,map[I]FuncDef,RET,RW,DIR,ERR]
}
func (s *mFuncDef[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lFuncDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFuncDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,FuncDef]()))))))
}
func (s *mFuncDef[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lFuncDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFuncDefOf(optic.Index(s.Traverse(),index))
}
type oFuncDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*FuncDef,*FuncDef,RET,RW,DIR,ERR]
}
func (s *oFuncDef[I,S,T,RET,RW,DIR,ERR])Docs()optic.MakeLensSlice[I,S,T,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FuncDef)*[]string{
return &x.Docs
})))))))
}
func (s *oFuncDef[I,S,T,RET,RW,DIR,ERR])Name()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FuncDef)*string{
return &x.Name
}))))))
}
func (s *oFuncDef[I,S,T,RET,RW,DIR,ERR])TypeParams()optic.MakeLensSlice[I,S,T,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FuncDef)*[]TypeExpression{
return &x.TypeParams
})))))))
}
func (s *oFuncDef[I,S,T,RET,RW,DIR,ERR])Params()*sParam[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sParam[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FuncDef)*[]Param{
return &x.Params
})))))),optic.SliceToCol[Param]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FuncDef)*[]Param{
return &x.Params
})))))),
}
}
func (s *oFuncDef[I,S,T,RET,RW,DIR,ERR])ReturnTypes()optic.MakeLensSlice[I,S,T,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FuncDef)*[]TypeExpression{
return &x.ReturnTypes
})))))))
}
func (s *oFuncDef[I,S,T,RET,RW,DIR,ERR])Body()optic.MakeLensSlice[I,S,T,Statement,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FuncDef)*[]Statement{
return &x.Body
})))))))
}
func (s *oFuncDef[I,S,T,RET,RW,DIR,ERR])Some()*lFuncDef[optic.Void,mo.Option[FuncDef],mo.Option[FuncDef],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGFuncDefOf(optic.Some[FuncDef]())
}
func (s *oFuncDef[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[FuncDef],mo.Option[FuncDef],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[FuncDef]())))))
}
type lFuncType[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,FuncType,FuncType,RET,RW,DIR,ERR]
}
func (s *lFuncType[I,S,T,RET,RW,DIR,ERR])Params()*sParam[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sParam[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FuncType)*[]Param{
return &x.Params
})))))),optic.SliceToCol[Param]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FuncType)*[]Param{
return &x.Params
})))))),
}
}
func (s *lFuncType[I,S,T,RET,RW,DIR,ERR])ReturnTypes()optic.MakeLensSlice[I,S,T,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *FuncType)*[]TypeExpression{
return &x.ReturnTypes
})))))))
}
type sFuncType[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,FuncType, optic.Pure],optic.Collection[int,FuncType, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]FuncType,[]FuncType,RET,RW,DIR,ERR]
}
func (s *sFuncType[I,S,T,RET,RW,DIR,ERR])Traverse()*lFuncType[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFuncTypeOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[FuncType]()))))))
}
func (s *sFuncType[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lFuncType[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFuncTypeOf(optic.Index(s.Traverse(),index))
}
type mFuncType[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,FuncType, optic.Pure],optic.Collection[I,FuncType, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]FuncType,map[I]FuncType,RET,RW,DIR,ERR]
}
func (s *mFuncType[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lFuncType[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFuncTypeOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,FuncType]()))))))
}
func (s *mFuncType[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lFuncType[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGFuncTypeOf(optic.Index(s.Traverse(),index))
}
type oFuncType[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*FuncType,*FuncType,RET,RW,DIR,ERR]
}
func (s *oFuncType[I,S,T,RET,RW,DIR,ERR])Params()*sParam[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sParam[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FuncType)*[]Param{
return &x.Params
})))))),optic.SliceToCol[Param]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FuncType)*[]Param{
return &x.Params
})))))),
}
}
func (s *oFuncType[I,S,T,RET,RW,DIR,ERR])ReturnTypes()optic.MakeLensSlice[I,S,T,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *FuncType)*[]TypeExpression{
return &x.ReturnTypes
})))))))
}
func (s *oFuncType[I,S,T,RET,RW,DIR,ERR])Some()*lFuncType[optic.Void,mo.Option[FuncType],mo.Option[FuncType],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGFuncTypeOf(optic.Some[FuncType]())
}
func (s *oFuncType[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[FuncType],mo.Option[FuncType],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[FuncType]())))))
}
type lGenFuncExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,GenFuncExpr,GenFuncExpr,RET,RW,DIR,ERR]
}
func (s *lGenFuncExpr[I,S,T,RET,RW,DIR,ERR])Params()*sParam[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sParam[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *GenFuncExpr)*[]Param{
return &x.Params
})))))),optic.SliceToCol[Param]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *GenFuncExpr)*[]Param{
return &x.Params
})))))),
}
}
func (s *lGenFuncExpr[I,S,T,RET,RW,DIR,ERR])Body()optic.MakeLensSlice[I,S,T,Statement,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *GenFuncExpr)*[]Statement{
return &x.Body
})))))))
}
func (s *lGenFuncExpr[I,S,T,RET,RW,DIR,ERR])ReturnTypes()optic.MakeLensSlice[I,S,T,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *GenFuncExpr)*[]TypeExpression{
return &x.ReturnTypes
})))))))
}
type sGenFuncExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,GenFuncExpr, optic.Pure],optic.Collection[int,GenFuncExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]GenFuncExpr,[]GenFuncExpr,RET,RW,DIR,ERR]
}
func (s *sGenFuncExpr[I,S,T,RET,RW,DIR,ERR])Traverse()*lGenFuncExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGGenFuncExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[GenFuncExpr]()))))))
}
func (s *sGenFuncExpr[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lGenFuncExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGGenFuncExprOf(optic.Index(s.Traverse(),index))
}
type mGenFuncExpr[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,GenFuncExpr, optic.Pure],optic.Collection[I,GenFuncExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]GenFuncExpr,map[I]GenFuncExpr,RET,RW,DIR,ERR]
}
func (s *mGenFuncExpr[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lGenFuncExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGGenFuncExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,GenFuncExpr]()))))))
}
func (s *mGenFuncExpr[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lGenFuncExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGGenFuncExprOf(optic.Index(s.Traverse(),index))
}
type oGenFuncExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*GenFuncExpr,*GenFuncExpr,RET,RW,DIR,ERR]
}
func (s *oGenFuncExpr[I,S,T,RET,RW,DIR,ERR])Params()*sParam[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sParam[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *GenFuncExpr)*[]Param{
return &x.Params
})))))),optic.SliceToCol[Param]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *GenFuncExpr)*[]Param{
return &x.Params
})))))),
}
}
func (s *oGenFuncExpr[I,S,T,RET,RW,DIR,ERR])Body()optic.MakeLensSlice[I,S,T,Statement,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *GenFuncExpr)*[]Statement{
return &x.Body
})))))))
}
func (s *oGenFuncExpr[I,S,T,RET,RW,DIR,ERR])ReturnTypes()optic.MakeLensSlice[I,S,T,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *GenFuncExpr)*[]TypeExpression{
return &x.ReturnTypes
})))))))
}
func (s *oGenFuncExpr[I,S,T,RET,RW,DIR,ERR])Some()*lGenFuncExpr[optic.Void,mo.Option[GenFuncExpr],mo.Option[GenFuncExpr],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGGenFuncExprOf(optic.Some[GenFuncExpr]())
}
func (s *oGenFuncExpr[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[GenFuncExpr],mo.Option[GenFuncExpr],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[GenFuncExpr]())))))
}
type lIfStmnt[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,IfStmnt,IfStmnt,RET,RW,DIR,ERR]
}
func (s *lIfStmnt[I,S,T,RET,RW,DIR,ERR])Condition()optic.Optic[I,S,T,Expression,Expression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *IfStmnt)*Expression{
return &x.Condition
}))))))
}
func (s *lIfStmnt[I,S,T,RET,RW,DIR,ERR])OnTrue()optic.MakeLensSlice[I,S,T,Statement,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *IfStmnt)*[]Statement{
return &x.OnTrue
})))))))
}
func (s *lIfStmnt[I,S,T,RET,RW,DIR,ERR])OnFalse()optic.MakeLensSlice[I,S,T,Statement,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *IfStmnt)*[]Statement{
return &x.OnFalse
})))))))
}
type sIfStmnt[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,IfStmnt, optic.Pure],optic.Collection[int,IfStmnt, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]IfStmnt,[]IfStmnt,RET,RW,DIR,ERR]
}
func (s *sIfStmnt[I,S,T,RET,RW,DIR,ERR])Traverse()*lIfStmnt[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGIfStmntOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[IfStmnt]()))))))
}
func (s *sIfStmnt[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lIfStmnt[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGIfStmntOf(optic.Index(s.Traverse(),index))
}
type mIfStmnt[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,IfStmnt, optic.Pure],optic.Collection[I,IfStmnt, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]IfStmnt,map[I]IfStmnt,RET,RW,DIR,ERR]
}
func (s *mIfStmnt[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lIfStmnt[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGIfStmntOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,IfStmnt]()))))))
}
func (s *mIfStmnt[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lIfStmnt[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGIfStmntOf(optic.Index(s.Traverse(),index))
}
type oIfStmnt[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*IfStmnt,*IfStmnt,RET,RW,DIR,ERR]
}
func (s *oIfStmnt[I,S,T,RET,RW,DIR,ERR])Condition()optic.Optic[I,S,T,Expression,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *IfStmnt)*Expression{
return &x.Condition
}))))))
}
func (s *oIfStmnt[I,S,T,RET,RW,DIR,ERR])OnTrue()optic.MakeLensSlice[I,S,T,Statement,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *IfStmnt)*[]Statement{
return &x.OnTrue
})))))))
}
func (s *oIfStmnt[I,S,T,RET,RW,DIR,ERR])OnFalse()optic.MakeLensSlice[I,S,T,Statement,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *IfStmnt)*[]Statement{
return &x.OnFalse
})))))))
}
func (s *oIfStmnt[I,S,T,RET,RW,DIR,ERR])Some()*lIfStmnt[optic.Void,mo.Option[IfStmnt],mo.Option[IfStmnt],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGIfStmntOf(optic.Some[IfStmnt]())
}
func (s *oIfStmnt[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[IfStmnt],mo.Option[IfStmnt],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[IfStmnt]())))))
}
type lMapDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,MapDef,MapDef,RET,RW,DIR,ERR]
}
func (s *lMapDef[I,S,T,RET,RW,DIR,ERR])Key()optic.Optic[I,S,T,TypeExpression,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *MapDef)*TypeExpression{
return &x.Key
}))))))
}
func (s *lMapDef[I,S,T,RET,RW,DIR,ERR])Type()optic.Optic[I,S,T,TypeExpression,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *MapDef)*TypeExpression{
return &x.Type
}))))))
}
type sMapDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,MapDef, optic.Pure],optic.Collection[int,MapDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]MapDef,[]MapDef,RET,RW,DIR,ERR]
}
func (s *sMapDef[I,S,T,RET,RW,DIR,ERR])Traverse()*lMapDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGMapDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[MapDef]()))))))
}
func (s *sMapDef[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lMapDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGMapDefOf(optic.Index(s.Traverse(),index))
}
type mMapDef[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,MapDef, optic.Pure],optic.Collection[I,MapDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]MapDef,map[I]MapDef,RET,RW,DIR,ERR]
}
func (s *mMapDef[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lMapDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGMapDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,MapDef]()))))))
}
func (s *mMapDef[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lMapDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGMapDefOf(optic.Index(s.Traverse(),index))
}
type oMapDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*MapDef,*MapDef,RET,RW,DIR,ERR]
}
func (s *oMapDef[I,S,T,RET,RW,DIR,ERR])Key()optic.Optic[I,S,T,TypeExpression,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *MapDef)*TypeExpression{
return &x.Key
}))))))
}
func (s *oMapDef[I,S,T,RET,RW,DIR,ERR])Type()optic.Optic[I,S,T,TypeExpression,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *MapDef)*TypeExpression{
return &x.Type
}))))))
}
func (s *oMapDef[I,S,T,RET,RW,DIR,ERR])Some()*lMapDef[optic.Void,mo.Option[MapDef],mo.Option[MapDef],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGMapDefOf(optic.Some[MapDef]())
}
func (s *oMapDef[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[MapDef],mo.Option[MapDef],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[MapDef]())))))
}
type lMethodCallExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,MethodCallExpr,MethodCallExpr,RET,RW,DIR,ERR]
}
func (s *lMethodCallExpr[I,S,T,RET,RW,DIR,ERR])Receiver()optic.Optic[I,S,T,Expression,Expression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *MethodCallExpr)*Expression{
return &x.Receiver
}))))))
}
func (s *lMethodCallExpr[I,S,T,RET,RW,DIR,ERR])Name()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *MethodCallExpr)*string{
return &x.Name
})))))))
}
func (s *lMethodCallExpr[I,S,T,RET,RW,DIR,ERR])Params()optic.MakeLensSlice[I,S,T,Expression,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *MethodCallExpr)*[]Expression{
return &x.Params
})))))))
}
type sMethodCallExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,MethodCallExpr, optic.Pure],optic.Collection[int,MethodCallExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]MethodCallExpr,[]MethodCallExpr,RET,RW,DIR,ERR]
}
func (s *sMethodCallExpr[I,S,T,RET,RW,DIR,ERR])Traverse()*lMethodCallExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGMethodCallExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[MethodCallExpr]()))))))
}
func (s *sMethodCallExpr[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lMethodCallExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGMethodCallExprOf(optic.Index(s.Traverse(),index))
}
type mMethodCallExpr[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,MethodCallExpr, optic.Pure],optic.Collection[I,MethodCallExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]MethodCallExpr,map[I]MethodCallExpr,RET,RW,DIR,ERR]
}
func (s *mMethodCallExpr[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lMethodCallExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGMethodCallExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,MethodCallExpr]()))))))
}
func (s *mMethodCallExpr[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lMethodCallExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGMethodCallExprOf(optic.Index(s.Traverse(),index))
}
type oMethodCallExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*MethodCallExpr,*MethodCallExpr,RET,RW,DIR,ERR]
}
func (s *oMethodCallExpr[I,S,T,RET,RW,DIR,ERR])Receiver()optic.Optic[I,S,T,Expression,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *MethodCallExpr)*Expression{
return &x.Receiver
}))))))
}
func (s *oMethodCallExpr[I,S,T,RET,RW,DIR,ERR])Name()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *MethodCallExpr)*string{
return &x.Name
}))))))
}
func (s *oMethodCallExpr[I,S,T,RET,RW,DIR,ERR])Params()optic.MakeLensSlice[I,S,T,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *MethodCallExpr)*[]Expression{
return &x.Params
})))))))
}
func (s *oMethodCallExpr[I,S,T,RET,RW,DIR,ERR])Some()*lMethodCallExpr[optic.Void,mo.Option[MethodCallExpr],mo.Option[MethodCallExpr],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGMethodCallExprOf(optic.Some[MethodCallExpr]())
}
func (s *oMethodCallExpr[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[MethodCallExpr],mo.Option[MethodCallExpr],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[MethodCallExpr]())))))
}
type lMethodDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,MethodDef,MethodDef,RET,RW,DIR,ERR]
}
func (s *lMethodDef[I,S,T,RET,RW,DIR,ERR])Docs()optic.MakeLensSlice[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *MethodDef)*[]string{
return &x.Docs
})))))))
}
func (s *lMethodDef[I,S,T,RET,RW,DIR,ERR])Name()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *MethodDef)*string{
return &x.Name
})))))))
}
func (s *lMethodDef[I,S,T,RET,RW,DIR,ERR])Params()*sParam[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sParam[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *MethodDef)*[]Param{
return &x.Params
})))))),optic.SliceToCol[Param]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *MethodDef)*[]Param{
return &x.Params
})))))),
}
}
func (s *lMethodDef[I,S,T,RET,RW,DIR,ERR])ReturnTypes()optic.MakeLensSlice[I,S,T,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *MethodDef)*[]TypeExpression{
return &x.ReturnTypes
})))))))
}
func (s *lMethodDef[I,S,T,RET,RW,DIR,ERR])Body()optic.MakeLensSlice[I,S,T,Statement,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *MethodDef)*[]Statement{
return &x.Body
})))))))
}
type sMethodDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,MethodDef, optic.Pure],optic.Collection[int,MethodDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]MethodDef,[]MethodDef,RET,RW,DIR,ERR]
}
func (s *sMethodDef[I,S,T,RET,RW,DIR,ERR])Traverse()*lMethodDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGMethodDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[MethodDef]()))))))
}
func (s *sMethodDef[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lMethodDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGMethodDefOf(optic.Index(s.Traverse(),index))
}
type mMethodDef[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,MethodDef, optic.Pure],optic.Collection[I,MethodDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]MethodDef,map[I]MethodDef,RET,RW,DIR,ERR]
}
func (s *mMethodDef[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lMethodDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGMethodDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,MethodDef]()))))))
}
func (s *mMethodDef[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lMethodDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGMethodDefOf(optic.Index(s.Traverse(),index))
}
type oMethodDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*MethodDef,*MethodDef,RET,RW,DIR,ERR]
}
func (s *oMethodDef[I,S,T,RET,RW,DIR,ERR])Docs()optic.MakeLensSlice[I,S,T,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *MethodDef)*[]string{
return &x.Docs
})))))))
}
func (s *oMethodDef[I,S,T,RET,RW,DIR,ERR])Name()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *MethodDef)*string{
return &x.Name
}))))))
}
func (s *oMethodDef[I,S,T,RET,RW,DIR,ERR])Params()*sParam[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sParam[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *MethodDef)*[]Param{
return &x.Params
})))))),optic.SliceToCol[Param]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *MethodDef)*[]Param{
return &x.Params
})))))),
}
}
func (s *oMethodDef[I,S,T,RET,RW,DIR,ERR])ReturnTypes()optic.MakeLensSlice[I,S,T,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *MethodDef)*[]TypeExpression{
return &x.ReturnTypes
})))))))
}
func (s *oMethodDef[I,S,T,RET,RW,DIR,ERR])Body()optic.MakeLensSlice[I,S,T,Statement,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *MethodDef)*[]Statement{
return &x.Body
})))))))
}
func (s *oMethodDef[I,S,T,RET,RW,DIR,ERR])Some()*lMethodDef[optic.Void,mo.Option[MethodDef],mo.Option[MethodDef],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGMethodDefOf(optic.Some[MethodDef]())
}
func (s *oMethodDef[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[MethodDef],mo.Option[MethodDef],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[MethodDef]())))))
}
type lParam[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Param,Param,RET,RW,DIR,ERR]
}
func (s *lParam[I,S,T,RET,RW,DIR,ERR])Name()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Param)*string{
return &x.Name
})))))))
}
func (s *lParam[I,S,T,RET,RW,DIR,ERR])Type()optic.Optic[I,S,T,TypeExpression,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Param)*TypeExpression{
return &x.Type
}))))))
}
type sParam[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Param, optic.Pure],optic.Collection[int,Param, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Param,[]Param,RET,RW,DIR,ERR]
}
func (s *sParam[I,S,T,RET,RW,DIR,ERR])Traverse()*lParam[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGParamOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Param]()))))))
}
func (s *sParam[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lParam[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGParamOf(optic.Index(s.Traverse(),index))
}
type mParam[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Param, optic.Pure],optic.Collection[I,Param, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Param,map[I]Param,RET,RW,DIR,ERR]
}
func (s *mParam[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lParam[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGParamOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Param]()))))))
}
func (s *mParam[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lParam[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGParamOf(optic.Index(s.Traverse(),index))
}
type oParam[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Param,*Param,RET,RW,DIR,ERR]
}
func (s *oParam[I,S,T,RET,RW,DIR,ERR])Name()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Param)*string{
return &x.Name
}))))))
}
func (s *oParam[I,S,T,RET,RW,DIR,ERR])Type()optic.Optic[I,S,T,TypeExpression,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Param)*TypeExpression{
return &x.Type
}))))))
}
func (s *oParam[I,S,T,RET,RW,DIR,ERR])Some()*lParam[optic.Void,mo.Option[Param],mo.Option[Param],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGParamOf(optic.Some[Param]())
}
func (s *oParam[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Param],mo.Option[Param],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Param]())))))
}
type lReturnStmnt[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,ReturnStmnt,ReturnStmnt,RET,RW,DIR,ERR]
}
func (s *lReturnStmnt[I,S,T,RET,RW,DIR,ERR])Values()optic.MakeLensSlice[I,S,T,Expression,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *ReturnStmnt)*[]Expression{
return &x.Values
})))))))
}
type sReturnStmnt[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,ReturnStmnt, optic.Pure],optic.Collection[int,ReturnStmnt, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]ReturnStmnt,[]ReturnStmnt,RET,RW,DIR,ERR]
}
func (s *sReturnStmnt[I,S,T,RET,RW,DIR,ERR])Traverse()*lReturnStmnt[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGReturnStmntOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[ReturnStmnt]()))))))
}
func (s *sReturnStmnt[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lReturnStmnt[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGReturnStmntOf(optic.Index(s.Traverse(),index))
}
type mReturnStmnt[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,ReturnStmnt, optic.Pure],optic.Collection[I,ReturnStmnt, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]ReturnStmnt,map[I]ReturnStmnt,RET,RW,DIR,ERR]
}
func (s *mReturnStmnt[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lReturnStmnt[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGReturnStmntOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,ReturnStmnt]()))))))
}
func (s *mReturnStmnt[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lReturnStmnt[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGReturnStmntOf(optic.Index(s.Traverse(),index))
}
type oReturnStmnt[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*ReturnStmnt,*ReturnStmnt,RET,RW,DIR,ERR]
}
func (s *oReturnStmnt[I,S,T,RET,RW,DIR,ERR])Values()optic.MakeLensSlice[I,S,T,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *ReturnStmnt)*[]Expression{
return &x.Values
})))))))
}
func (s *oReturnStmnt[I,S,T,RET,RW,DIR,ERR])Some()*lReturnStmnt[optic.Void,mo.Option[ReturnStmnt],mo.Option[ReturnStmnt],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGReturnStmntOf(optic.Some[ReturnStmnt]())
}
func (s *oReturnStmnt[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[ReturnStmnt],mo.Option[ReturnStmnt],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[ReturnStmnt]())))))
}
type lSliceDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,SliceDef,SliceDef,RET,RW,DIR,ERR]
}
func (s *lSliceDef[I,S,T,RET,RW,DIR,ERR])Type()optic.Optic[I,S,T,TypeExpression,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *SliceDef)*TypeExpression{
return &x.Type
}))))))
}
type sSliceDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,SliceDef, optic.Pure],optic.Collection[int,SliceDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]SliceDef,[]SliceDef,RET,RW,DIR,ERR]
}
func (s *sSliceDef[I,S,T,RET,RW,DIR,ERR])Traverse()*lSliceDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGSliceDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[SliceDef]()))))))
}
func (s *sSliceDef[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lSliceDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGSliceDefOf(optic.Index(s.Traverse(),index))
}
type mSliceDef[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,SliceDef, optic.Pure],optic.Collection[I,SliceDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]SliceDef,map[I]SliceDef,RET,RW,DIR,ERR]
}
func (s *mSliceDef[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lSliceDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGSliceDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,SliceDef]()))))))
}
func (s *mSliceDef[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lSliceDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGSliceDefOf(optic.Index(s.Traverse(),index))
}
type oSliceDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*SliceDef,*SliceDef,RET,RW,DIR,ERR]
}
func (s *oSliceDef[I,S,T,RET,RW,DIR,ERR])Type()optic.Optic[I,S,T,TypeExpression,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *SliceDef)*TypeExpression{
return &x.Type
}))))))
}
func (s *oSliceDef[I,S,T,RET,RW,DIR,ERR])Some()*lSliceDef[optic.Void,mo.Option[SliceDef],mo.Option[SliceDef],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGSliceDefOf(optic.Some[SliceDef]())
}
func (s *oSliceDef[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[SliceDef],mo.Option[SliceDef],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[SliceDef]())))))
}
type lSliceExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,SliceExpr,SliceExpr,RET,RW,DIR,ERR]
}
func (s *lSliceExpr[I,S,T,RET,RW,DIR,ERR])Type()*lTypeDef[I,S,T,RET,RW,optic.UniDir,ERR]{
return &lTypeDef[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *SliceExpr)*TypeDef{
return &x.Type
})))))),
}
}
func (s *lSliceExpr[I,S,T,RET,RW,DIR,ERR])Values()optic.MakeLensSlice[I,S,T,Expression,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *SliceExpr)*[]Expression{
return &x.Values
})))))))
}
type sSliceExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,SliceExpr, optic.Pure],optic.Collection[int,SliceExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]SliceExpr,[]SliceExpr,RET,RW,DIR,ERR]
}
func (s *sSliceExpr[I,S,T,RET,RW,DIR,ERR])Traverse()*lSliceExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGSliceExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[SliceExpr]()))))))
}
func (s *sSliceExpr[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lSliceExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGSliceExprOf(optic.Index(s.Traverse(),index))
}
type mSliceExpr[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,SliceExpr, optic.Pure],optic.Collection[I,SliceExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]SliceExpr,map[I]SliceExpr,RET,RW,DIR,ERR]
}
func (s *mSliceExpr[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lSliceExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGSliceExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,SliceExpr]()))))))
}
func (s *mSliceExpr[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lSliceExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGSliceExprOf(optic.Index(s.Traverse(),index))
}
type oSliceExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*SliceExpr,*SliceExpr,RET,RW,DIR,ERR]
}
func (s *oSliceExpr[I,S,T,RET,RW,DIR,ERR])Type()*lTypeDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &lTypeDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *SliceExpr)*TypeDef{
return &x.Type
})))))),
}
}
func (s *oSliceExpr[I,S,T,RET,RW,DIR,ERR])Values()optic.MakeLensSlice[I,S,T,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *SliceExpr)*[]Expression{
return &x.Values
})))))))
}
func (s *oSliceExpr[I,S,T,RET,RW,DIR,ERR])Some()*lSliceExpr[optic.Void,mo.Option[SliceExpr],mo.Option[SliceExpr],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGSliceExprOf(optic.Some[SliceExpr]())
}
func (s *oSliceExpr[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[SliceExpr],mo.Option[SliceExpr],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[SliceExpr]())))))
}
type lStar[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Star,Star,RET,RW,DIR,ERR]
}
func (s *lStar[I,S,T,RET,RW,DIR,ERR])Type()optic.Optic[I,S,T,TypeExpression,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Star)*TypeExpression{
return &x.Type
}))))))
}
type sStar[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Star, optic.Pure],optic.Collection[int,Star, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Star,[]Star,RET,RW,DIR,ERR]
}
func (s *sStar[I,S,T,RET,RW,DIR,ERR])Traverse()*lStar[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGStarOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Star]()))))))
}
func (s *sStar[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lStar[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGStarOf(optic.Index(s.Traverse(),index))
}
type mStar[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Star, optic.Pure],optic.Collection[I,Star, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Star,map[I]Star,RET,RW,DIR,ERR]
}
func (s *mStar[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lStar[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGStarOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Star]()))))))
}
func (s *mStar[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lStar[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGStarOf(optic.Index(s.Traverse(),index))
}
type oStar[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Star,*Star,RET,RW,DIR,ERR]
}
func (s *oStar[I,S,T,RET,RW,DIR,ERR])Type()optic.Optic[I,S,T,TypeExpression,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Star)*TypeExpression{
return &x.Type
}))))))
}
func (s *oStar[I,S,T,RET,RW,DIR,ERR])Some()*lStar[optic.Void,mo.Option[Star],mo.Option[Star],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGStarOf(optic.Some[Star]())
}
func (s *oStar[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Star],mo.Option[Star],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Star]())))))
}
type lStructDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,StructDef,StructDef,RET,RW,DIR,ERR]
}
func (s *lStructDef[I,S,T,RET,RW,DIR,ERR])Name()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *StructDef)*string{
return &x.Name
})))))))
}
func (s *lStructDef[I,S,T,RET,RW,DIR,ERR])TypeParams()optic.MakeLensSlice[I,S,T,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *StructDef)*[]TypeExpression{
return &x.TypeParams
})))))))
}
func (s *lStructDef[I,S,T,RET,RW,DIR,ERR])Fields()*sFieldDef[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sFieldDef[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *StructDef)*[]FieldDef{
return &x.Fields
})))))),optic.SliceToCol[FieldDef]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *StructDef)*[]FieldDef{
return &x.Fields
})))))),
}
}
func (s *lStructDef[I,S,T,RET,RW,DIR,ERR])Methods()*sMethodDef[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sMethodDef[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *StructDef)*[]MethodDef{
return &x.Methods
})))))),optic.SliceToCol[MethodDef]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *StructDef)*[]MethodDef{
return &x.Methods
})))))),
}
}
type sStructDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,StructDef, optic.Pure],optic.Collection[int,StructDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]StructDef,[]StructDef,RET,RW,DIR,ERR]
}
func (s *sStructDef[I,S,T,RET,RW,DIR,ERR])Traverse()*lStructDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGStructDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[StructDef]()))))))
}
func (s *sStructDef[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lStructDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGStructDefOf(optic.Index(s.Traverse(),index))
}
type mStructDef[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,StructDef, optic.Pure],optic.Collection[I,StructDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]StructDef,map[I]StructDef,RET,RW,DIR,ERR]
}
func (s *mStructDef[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lStructDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGStructDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,StructDef]()))))))
}
func (s *mStructDef[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lStructDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGStructDefOf(optic.Index(s.Traverse(),index))
}
type oStructDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*StructDef,*StructDef,RET,RW,DIR,ERR]
}
func (s *oStructDef[I,S,T,RET,RW,DIR,ERR])Name()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *StructDef)*string{
return &x.Name
}))))))
}
func (s *oStructDef[I,S,T,RET,RW,DIR,ERR])TypeParams()optic.MakeLensSlice[I,S,T,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *StructDef)*[]TypeExpression{
return &x.TypeParams
})))))))
}
func (s *oStructDef[I,S,T,RET,RW,DIR,ERR])Fields()*sFieldDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sFieldDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *StructDef)*[]FieldDef{
return &x.Fields
})))))),optic.SliceToCol[FieldDef]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *StructDef)*[]FieldDef{
return &x.Fields
})))))),
}
}
func (s *oStructDef[I,S,T,RET,RW,DIR,ERR])Methods()*sMethodDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sMethodDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *StructDef)*[]MethodDef{
return &x.Methods
})))))),optic.SliceToCol[MethodDef]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *StructDef)*[]MethodDef{
return &x.Methods
})))))),
}
}
func (s *oStructDef[I,S,T,RET,RW,DIR,ERR])Some()*lStructDef[optic.Void,mo.Option[StructDef],mo.Option[StructDef],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGStructDefOf(optic.Some[StructDef]())
}
func (s *oStructDef[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[StructDef],mo.Option[StructDef],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[StructDef]())))))
}
type lStructExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,StructExpr,StructExpr,RET,RW,DIR,ERR]
}
func (s *lStructExpr[I,S,T,RET,RW,DIR,ERR])Type()*lTypeDef[I,S,T,RET,RW,optic.UniDir,ERR]{
return &lTypeDef[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *StructExpr)*TypeDef{
return &x.Type
})))))),
}
}
func (s *lStructExpr[I,S,T,RET,RW,DIR,ERR])Fields()*sAssignField[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sAssignField[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *StructExpr)*[]AssignField{
return &x.Fields
})))))),optic.SliceToCol[AssignField]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *StructExpr)*[]AssignField{
return &x.Fields
})))))),
}
}
type sStructExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,StructExpr, optic.Pure],optic.Collection[int,StructExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]StructExpr,[]StructExpr,RET,RW,DIR,ERR]
}
func (s *sStructExpr[I,S,T,RET,RW,DIR,ERR])Traverse()*lStructExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGStructExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[StructExpr]()))))))
}
func (s *sStructExpr[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lStructExpr[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGStructExprOf(optic.Index(s.Traverse(),index))
}
type mStructExpr[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,StructExpr, optic.Pure],optic.Collection[I,StructExpr, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]StructExpr,map[I]StructExpr,RET,RW,DIR,ERR]
}
func (s *mStructExpr[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lStructExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGStructExprOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,StructExpr]()))))))
}
func (s *mStructExpr[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lStructExpr[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGStructExprOf(optic.Index(s.Traverse(),index))
}
type oStructExpr[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*StructExpr,*StructExpr,RET,RW,DIR,ERR]
}
func (s *oStructExpr[I,S,T,RET,RW,DIR,ERR])Type()*lTypeDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &lTypeDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *StructExpr)*TypeDef{
return &x.Type
})))))),
}
}
func (s *oStructExpr[I,S,T,RET,RW,DIR,ERR])Fields()*sAssignField[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sAssignField[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *StructExpr)*[]AssignField{
return &x.Fields
})))))),optic.SliceToCol[AssignField]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *StructExpr)*[]AssignField{
return &x.Fields
})))))),
}
}
func (s *oStructExpr[I,S,T,RET,RW,DIR,ERR])Some()*lStructExpr[optic.Void,mo.Option[StructExpr],mo.Option[StructExpr],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGStructExprOf(optic.Some[StructExpr]())
}
func (s *oStructExpr[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[StructExpr],mo.Option[StructExpr],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[StructExpr]())))))
}
type lTypeDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,TypeDef,TypeDef,RET,RW,DIR,ERR]
}
func (s *lTypeDef[I,S,T,RET,RW,DIR,ERR])Name()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *TypeDef)*string{
return &x.Name
})))))))
}
func (s *lTypeDef[I,S,T,RET,RW,DIR,ERR])TypeParams()optic.MakeLensSlice[I,S,T,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *TypeDef)*[]TypeExpression{
return &x.TypeParams
})))))))
}
func (s *lTypeDef[I,S,T,RET,RW,DIR,ERR])Constraint()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *TypeDef)*string{
return &x.Constraint
})))))))
}
type sTypeDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,TypeDef, optic.Pure],optic.Collection[int,TypeDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]TypeDef,[]TypeDef,RET,RW,DIR,ERR]
}
func (s *sTypeDef[I,S,T,RET,RW,DIR,ERR])Traverse()*lTypeDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGTypeDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[TypeDef]()))))))
}
func (s *sTypeDef[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lTypeDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGTypeDefOf(optic.Index(s.Traverse(),index))
}
type mTypeDef[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,TypeDef, optic.Pure],optic.Collection[I,TypeDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]TypeDef,map[I]TypeDef,RET,RW,DIR,ERR]
}
func (s *mTypeDef[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lTypeDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGTypeDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,TypeDef]()))))))
}
func (s *mTypeDef[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lTypeDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGTypeDefOf(optic.Index(s.Traverse(),index))
}
type oTypeDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*TypeDef,*TypeDef,RET,RW,DIR,ERR]
}
func (s *oTypeDef[I,S,T,RET,RW,DIR,ERR])Name()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *TypeDef)*string{
return &x.Name
}))))))
}
func (s *oTypeDef[I,S,T,RET,RW,DIR,ERR])TypeParams()optic.MakeLensSlice[I,S,T,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *TypeDef)*[]TypeExpression{
return &x.TypeParams
})))))))
}
func (s *oTypeDef[I,S,T,RET,RW,DIR,ERR])Constraint()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *TypeDef)*string{
return &x.Constraint
}))))))
}
func (s *oTypeDef[I,S,T,RET,RW,DIR,ERR])Some()*lTypeDef[optic.Void,mo.Option[TypeDef],mo.Option[TypeDef],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGTypeDefOf(optic.Some[TypeDef]())
}
func (s *oTypeDef[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[TypeDef],mo.Option[TypeDef],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[TypeDef]())))))
}
type lVarDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,VarDef,VarDef,RET,RW,DIR,ERR]
}
func (s *lVarDef[I,S,T,RET,RW,DIR,ERR])Name()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *VarDef)*string{
return &x.Name
})))))))
}
func (s *lVarDef[I,S,T,RET,RW,DIR,ERR])Type()optic.Optic[I,S,T,TypeExpression,TypeExpression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *VarDef)*TypeExpression{
return &x.Type
}))))))
}
func (s *lVarDef[I,S,T,RET,RW,DIR,ERR])Value()optic.Optic[I,S,T,Expression,Expression,RET,RW,optic.UniDir,ERR]{
return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *VarDef)*Expression{
return &x.Value
}))))))
}
type sVarDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,VarDef, optic.Pure],optic.Collection[int,VarDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]VarDef,[]VarDef,RET,RW,DIR,ERR]
}
func (s *sVarDef[I,S,T,RET,RW,DIR,ERR])Traverse()*lVarDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGVarDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[VarDef]()))))))
}
func (s *sVarDef[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lVarDef[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGVarDefOf(optic.Index(s.Traverse(),index))
}
type mVarDef[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,VarDef, optic.Pure],optic.Collection[I,VarDef, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]VarDef,map[I]VarDef,RET,RW,DIR,ERR]
}
func (s *mVarDef[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lVarDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGVarDefOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,VarDef]()))))))
}
func (s *mVarDef[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lVarDef[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return CGVarDefOf(optic.Index(s.Traverse(),index))
}
type oVarDef[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*VarDef,*VarDef,RET,RW,DIR,ERR]
}
func (s *oVarDef[I,S,T,RET,RW,DIR,ERR])Name()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *VarDef)*string{
return &x.Name
}))))))
}
func (s *oVarDef[I,S,T,RET,RW,DIR,ERR])Type()optic.Optic[I,S,T,TypeExpression,TypeExpression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *VarDef)*TypeExpression{
return &x.Type
}))))))
}
func (s *oVarDef[I,S,T,RET,RW,DIR,ERR])Value()optic.Optic[I,S,T,Expression,Expression,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *VarDef)*Expression{
return &x.Value
}))))))
}
func (s *oVarDef[I,S,T,RET,RW,DIR,ERR])Some()*lVarDef[optic.Void,mo.Option[VarDef],mo.Option[VarDef],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGVarDefOf(optic.Some[VarDef]())
}
func (s *oVarDef[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[VarDef],mo.Option[VarDef],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[VarDef]())))))
}
type cg struct {
}
func (s *cg)AddressExpr()*lAddressExpr[optic.Void,AddressExpr,AddressExpr,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGAddressExprOf[optic.Void,AddressExpr,AddressExpr,optic.ReturnOne](optic.Identity[AddressExpr]())
}
func (s *cg)AssignField()*lAssignField[optic.Void,AssignField,AssignField,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGAssignFieldOf[optic.Void,AssignField,AssignField,optic.ReturnOne](optic.Identity[AssignField]())
}
func (s *cg)AssignVar()*lAssignVar[optic.Void,AssignVar,AssignVar,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGAssignVarOf[optic.Void,AssignVar,AssignVar,optic.ReturnOne](optic.Identity[AssignVar]())
}
func (s *cg)BinaryExpr()*lBinaryExpr[optic.Void,BinaryExpr,BinaryExpr,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGBinaryExprOf[optic.Void,BinaryExpr,BinaryExpr,optic.ReturnOne](optic.Identity[BinaryExpr]())
}
func (s *cg)CallExpr()*lCallExpr[optic.Void,CallExpr,CallExpr,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGCallExprOf[optic.Void,CallExpr,CallExpr,optic.ReturnOne](optic.Identity[CallExpr]())
}
func (s *cg)DeRef()*lDeRef[optic.Void,DeRef,DeRef,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGDeRefOf[optic.Void,DeRef,DeRef,optic.ReturnOne](optic.Identity[DeRef]())
}
func (s *cg)FieldDef()*lFieldDef[optic.Void,FieldDef,FieldDef,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGFieldDefOf[optic.Void,FieldDef,FieldDef,optic.ReturnOne](optic.Identity[FieldDef]())
}
func (s *cg)FileDef()*lFileDef[optic.Void,FileDef,FileDef,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGFileDefOf[optic.Void,FileDef,FileDef,optic.ReturnOne](optic.Identity[FileDef]())
}
func (s *cg)FuncDef()*lFuncDef[optic.Void,FuncDef,FuncDef,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGFuncDefOf[optic.Void,FuncDef,FuncDef,optic.ReturnOne](optic.Identity[FuncDef]())
}
func (s *cg)FuncType()*lFuncType[optic.Void,FuncType,FuncType,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGFuncTypeOf[optic.Void,FuncType,FuncType,optic.ReturnOne](optic.Identity[FuncType]())
}
func (s *cg)GenFuncExpr()*lGenFuncExpr[optic.Void,GenFuncExpr,GenFuncExpr,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGGenFuncExprOf[optic.Void,GenFuncExpr,GenFuncExpr,optic.ReturnOne](optic.Identity[GenFuncExpr]())
}
func (s *cg)IfStmnt()*lIfStmnt[optic.Void,IfStmnt,IfStmnt,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGIfStmntOf[optic.Void,IfStmnt,IfStmnt,optic.ReturnOne](optic.Identity[IfStmnt]())
}
func (s *cg)MapDef()*lMapDef[optic.Void,MapDef,MapDef,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGMapDefOf[optic.Void,MapDef,MapDef,optic.ReturnOne](optic.Identity[MapDef]())
}
func (s *cg)MethodCallExpr()*lMethodCallExpr[optic.Void,MethodCallExpr,MethodCallExpr,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGMethodCallExprOf[optic.Void,MethodCallExpr,MethodCallExpr,optic.ReturnOne](optic.Identity[MethodCallExpr]())
}
func (s *cg)MethodDef()*lMethodDef[optic.Void,MethodDef,MethodDef,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGMethodDefOf[optic.Void,MethodDef,MethodDef,optic.ReturnOne](optic.Identity[MethodDef]())
}
func (s *cg)Param()*lParam[optic.Void,Param,Param,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGParamOf[optic.Void,Param,Param,optic.ReturnOne](optic.Identity[Param]())
}
func (s *cg)ReturnStmnt()*lReturnStmnt[optic.Void,ReturnStmnt,ReturnStmnt,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGReturnStmntOf[optic.Void,ReturnStmnt,ReturnStmnt,optic.ReturnOne](optic.Identity[ReturnStmnt]())
}
func (s *cg)SliceDef()*lSliceDef[optic.Void,SliceDef,SliceDef,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGSliceDefOf[optic.Void,SliceDef,SliceDef,optic.ReturnOne](optic.Identity[SliceDef]())
}
func (s *cg)SliceExpr()*lSliceExpr[optic.Void,SliceExpr,SliceExpr,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGSliceExprOf[optic.Void,SliceExpr,SliceExpr,optic.ReturnOne](optic.Identity[SliceExpr]())
}
func (s *cg)Star()*lStar[optic.Void,Star,Star,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGStarOf[optic.Void,Star,Star,optic.ReturnOne](optic.Identity[Star]())
}
func (s *cg)StructDef()*lStructDef[optic.Void,StructDef,StructDef,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGStructDefOf[optic.Void,StructDef,StructDef,optic.ReturnOne](optic.Identity[StructDef]())
}
func (s *cg)StructExpr()*lStructExpr[optic.Void,StructExpr,StructExpr,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGStructExprOf[optic.Void,StructExpr,StructExpr,optic.ReturnOne](optic.Identity[StructExpr]())
}
func (s *cg)TypeDef()*lTypeDef[optic.Void,TypeDef,TypeDef,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGTypeDefOf[optic.Void,TypeDef,TypeDef,optic.ReturnOne](optic.Identity[TypeDef]())
}
func (s *cg)VarDef()*lVarDef[optic.Void,VarDef,VarDef,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return CGVarDefOf[optic.Void,VarDef,VarDef,optic.ReturnOne](optic.Identity[VarDef]())
}
func CGAddressExprOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,AddressExpr,AddressExpr,RET,RW,DIR,ERR])*lAddressExpr[I,S,T,RET,RW,DIR,ERR]{
return &lAddressExpr[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGAssignFieldOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,AssignField,AssignField,RET,RW,DIR,ERR])*lAssignField[I,S,T,RET,RW,DIR,ERR]{
return &lAssignField[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGAssignVarOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,AssignVar,AssignVar,RET,RW,DIR,ERR])*lAssignVar[I,S,T,RET,RW,DIR,ERR]{
return &lAssignVar[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGBinaryExprOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,BinaryExpr,BinaryExpr,RET,RW,DIR,ERR])*lBinaryExpr[I,S,T,RET,RW,DIR,ERR]{
return &lBinaryExpr[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGCallExprOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,CallExpr,CallExpr,RET,RW,DIR,ERR])*lCallExpr[I,S,T,RET,RW,DIR,ERR]{
return &lCallExpr[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGDeRefOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,DeRef,DeRef,RET,RW,DIR,ERR])*lDeRef[I,S,T,RET,RW,DIR,ERR]{
return &lDeRef[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGFieldDefOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,FieldDef,FieldDef,RET,RW,DIR,ERR])*lFieldDef[I,S,T,RET,RW,DIR,ERR]{
return &lFieldDef[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGFileDefOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,FileDef,FileDef,RET,RW,DIR,ERR])*lFileDef[I,S,T,RET,RW,DIR,ERR]{
return &lFileDef[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGFuncDefOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,FuncDef,FuncDef,RET,RW,DIR,ERR])*lFuncDef[I,S,T,RET,RW,DIR,ERR]{
return &lFuncDef[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGFuncTypeOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,FuncType,FuncType,RET,RW,DIR,ERR])*lFuncType[I,S,T,RET,RW,DIR,ERR]{
return &lFuncType[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGGenFuncExprOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,GenFuncExpr,GenFuncExpr,RET,RW,DIR,ERR])*lGenFuncExpr[I,S,T,RET,RW,DIR,ERR]{
return &lGenFuncExpr[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGIfStmntOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,IfStmnt,IfStmnt,RET,RW,DIR,ERR])*lIfStmnt[I,S,T,RET,RW,DIR,ERR]{
return &lIfStmnt[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGMapDefOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,MapDef,MapDef,RET,RW,DIR,ERR])*lMapDef[I,S,T,RET,RW,DIR,ERR]{
return &lMapDef[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGMethodCallExprOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,MethodCallExpr,MethodCallExpr,RET,RW,DIR,ERR])*lMethodCallExpr[I,S,T,RET,RW,DIR,ERR]{
return &lMethodCallExpr[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGMethodDefOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,MethodDef,MethodDef,RET,RW,DIR,ERR])*lMethodDef[I,S,T,RET,RW,DIR,ERR]{
return &lMethodDef[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGParamOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Param,Param,RET,RW,DIR,ERR])*lParam[I,S,T,RET,RW,DIR,ERR]{
return &lParam[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGReturnStmntOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,ReturnStmnt,ReturnStmnt,RET,RW,DIR,ERR])*lReturnStmnt[I,S,T,RET,RW,DIR,ERR]{
return &lReturnStmnt[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGSliceDefOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,SliceDef,SliceDef,RET,RW,DIR,ERR])*lSliceDef[I,S,T,RET,RW,DIR,ERR]{
return &lSliceDef[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGSliceExprOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,SliceExpr,SliceExpr,RET,RW,DIR,ERR])*lSliceExpr[I,S,T,RET,RW,DIR,ERR]{
return &lSliceExpr[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGStarOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Star,Star,RET,RW,DIR,ERR])*lStar[I,S,T,RET,RW,DIR,ERR]{
return &lStar[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGStructDefOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,StructDef,StructDef,RET,RW,DIR,ERR])*lStructDef[I,S,T,RET,RW,DIR,ERR]{
return &lStructDef[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGStructExprOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,StructExpr,StructExpr,RET,RW,DIR,ERR])*lStructExpr[I,S,T,RET,RW,DIR,ERR]{
return &lStructExpr[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGTypeDefOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,TypeDef,TypeDef,RET,RW,DIR,ERR])*lTypeDef[I,S,T,RET,RW,DIR,ERR]{
return &lTypeDef[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func CGVarDefOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,VarDef,VarDef,RET,RW,DIR,ERR])*lVarDef[I,S,T,RET,RW,DIR,ERR]{
return &lVarDef[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
var CG  = cg{
}
