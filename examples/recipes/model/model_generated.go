package model

import (
	"github.com/samber/mo"
	"github.com/spearson78/go-optic"
)
type lIngredient[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Ingredient,Ingredient,RET,RW,DIR,ERR]
}
func (s *lIngredient[I,S,T,RET,RW,DIR,ERR])Name()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Ingredient)*string{
return &x.name
})))))))
}
func (s *lIngredient[I,S,T,RET,RW,DIR,ERR])Quantity()optic.MakeLensRealOps[I,S,T,float64,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Ingredient)*float64{
return &x.quantity
})))))))
}
func (s *lIngredient[I,S,T,RET,RW,DIR,ERR])Unit()optic.MakeLensRealOps[I,S,T,Unit,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Ingredient)*Unit{
return &x.unit
})))))))
}
type sIngredient[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Ingredient, optic.Pure],optic.Collection[int,Ingredient, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Ingredient,[]Ingredient,RET,RW,DIR,ERR]
}
func (s *sIngredient[I,S,T,RET,RW,DIR,ERR])Traverse()*lIngredient[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OIngredientOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Ingredient]()))))))
}
func (s *sIngredient[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lIngredient[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OIngredientOf(optic.Index(s.Traverse(),index))
}
type mIngredient[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Ingredient, optic.Pure],optic.Collection[I,Ingredient, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Ingredient,map[I]Ingredient,RET,RW,DIR,ERR]
}
func (s *mIngredient[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lIngredient[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OIngredientOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Ingredient]()))))))
}
func (s *mIngredient[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lIngredient[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return OIngredientOf(optic.Index(s.Traverse(),index))
}
type oIngredient[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Ingredient,*Ingredient,RET,RW,DIR,ERR]
}
func (s *oIngredient[I,S,T,RET,RW,DIR,ERR])Name()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Ingredient)*string{
return &x.name
}))))))
}
func (s *oIngredient[I,S,T,RET,RW,DIR,ERR])Quantity()optic.Optic[I,S,T,float64,float64,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Ingredient)*float64{
return &x.quantity
}))))))
}
func (s *oIngredient[I,S,T,RET,RW,DIR,ERR])Unit()optic.Optic[I,S,T,Unit,Unit,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Ingredient)*Unit{
return &x.unit
}))))))
}
func (s *oIngredient[I,S,T,RET,RW,DIR,ERR])Some()*lIngredient[optic.Void,mo.Option[Ingredient],mo.Option[Ingredient],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OIngredientOf(optic.Some[Ingredient]())
}
func (s *oIngredient[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Ingredient],mo.Option[Ingredient],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Ingredient]())))))
}
type lRecipe[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,Recipe,Recipe,RET,RW,DIR,ERR]
}
func (s *lRecipe[I,S,T,RET,RW,DIR,ERR])Name()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Recipe)*string{
return &x.name
})))))))
}
func (s *lRecipe[I,S,T,RET,RW,DIR,ERR])Description()optic.MakeLensOrdOps[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Recipe)*string{
return &x.description
})))))))
}
func (s *lRecipe[I,S,T,RET,RW,DIR,ERR])Serves()optic.MakeLensRealOps[I,S,T,int,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Recipe)*int{
return &x.serves
})))))))
}
func (s *lRecipe[I,S,T,RET,RW,DIR,ERR])CookTime()optic.MakeLensRealOps[I,S,T,float64,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Recipe)*float64{
return &x.cookTime
})))))))
}
func (s *lRecipe[I,S,T,RET,RW,DIR,ERR])Difficulty()optic.MakeLensRealOps[I,S,T,Difficulty,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Recipe)*Difficulty{
return &x.difficulty
})))))))
}
func (s *lRecipe[I,S,T,RET,RW,DIR,ERR])Ingredients()*sIngredient[I,S,T,RET,RW,optic.UniDir,ERR]{
return &sIngredient[I,S,T,RET,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Recipe)*[]Ingredient{
return &x.ingredients
})))))),optic.SliceToCol[Ingredient]()))))),
	o:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Recipe)*[]Ingredient{
return &x.ingredients
})))))),
}
}
func (s *lRecipe[I,S,T,RET,RW,DIR,ERR])Instructions()optic.MakeLensSlice[I,S,T,string,RET,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.FieldLens(func(x *Recipe)*[]string{
return &x.instructions
})))))))
}
type sRecipe[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,optic.Collection[int,Recipe, optic.Pure],optic.Collection[int,Recipe, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[I,S,T,[]Recipe,[]Recipe,RET,RW,DIR,ERR]
}
func (s *sRecipe[I,S,T,RET,RW,DIR,ERR])Traverse()*lRecipe[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ORecipeOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseSlice[Recipe]()))))))
}
func (s *sRecipe[I,S,T,RET,RW,DIR,ERR])Nth(index int)*lRecipe[int,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ORecipeOf(optic.Index(s.Traverse(),index))
}
type mRecipe[I comparable,J any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[J,S,T,optic.Collection[I,Recipe, optic.Pure],optic.Collection[I,Recipe, optic.Pure],RET,RW,DIR,ERR]
	o optic.Optic[J,S,T,map[I]Recipe,map[I]Recipe,RET,RW,DIR,ERR]
}
func (s *mRecipe[I,J,S,T,RET,RW,DIR,ERR])Traverse()*lRecipe[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ORecipeOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o,optic.TraverseMap[I,Recipe]()))))))
}
func (s *mRecipe[I,J,S,T,RET,RW,DIR,ERR])Key(index I)*lRecipe[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return ORecipeOf(optic.Index(s.Traverse(),index))
}
type oRecipe[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,*Recipe,*Recipe,RET,RW,DIR,ERR]
}
func (s *oRecipe[I,S,T,RET,RW,DIR,ERR])Name()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Recipe)*string{
return &x.name
}))))))
}
func (s *oRecipe[I,S,T,RET,RW,DIR,ERR])Description()optic.Optic[I,S,T,string,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Recipe)*string{
return &x.description
}))))))
}
func (s *oRecipe[I,S,T,RET,RW,DIR,ERR])Serves()optic.Optic[I,S,T,int,int,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Recipe)*int{
return &x.serves
}))))))
}
func (s *oRecipe[I,S,T,RET,RW,DIR,ERR])CookTime()optic.Optic[I,S,T,float64,float64,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Recipe)*float64{
return &x.cookTime
}))))))
}
func (s *oRecipe[I,S,T,RET,RW,DIR,ERR])Difficulty()optic.Optic[I,S,T,Difficulty,Difficulty,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Recipe)*Difficulty{
return &x.difficulty
}))))))
}
func (s *oRecipe[I,S,T,RET,RW,DIR,ERR])Ingredients()*sIngredient[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
return &sIngredient[I,S,T,optic.ReturnMany,RW,optic.UniDir,ERR]{
	Optic:	optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Recipe)*[]Ingredient{
return &x.ingredients
})))))),optic.SliceToCol[Ingredient]()))))),
	o:	optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Recipe)*[]Ingredient{
return &x.ingredients
})))))),
}
}
func (s *oRecipe[I,S,T,RET,RW,DIR,ERR])Instructions()optic.MakeLensSlice[I,S,T,string,optic.ReturnMany,RW,optic.UniDir,ERR]{
return optic.NewMakeLensSlice(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s,optic.PtrFieldLens(func(x *Recipe)*[]string{
return &x.instructions
})))))))
}
func (s *oRecipe[I,S,T,RET,RW,DIR,ERR])Some()*lRecipe[optic.Void,mo.Option[Recipe],mo.Option[Recipe],optic.ReturnMany,optic.ReadWrite,optic.BiDir,optic.Pure]{
return ORecipeOf(optic.Some[Recipe]())
}
func (s *oRecipe[I,S,T,RET,RW,DIR,ERR])Option()optic.Optic[I,S,T,mo.Option[Recipe],mo.Option[Recipe],RET,RW,DIR,ERR]{
return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s,optic.PtrOption[Recipe]())))))
}
type o struct {
}
func (s *o)Ingredient()*lIngredient[optic.Void,Ingredient,Ingredient,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return OIngredientOf[optic.Void,Ingredient,Ingredient,optic.ReturnOne](optic.Identity[Ingredient]())
}
func (s *o)Recipe()*lRecipe[optic.Void,Recipe,Recipe,optic.ReturnOne,optic.ReadWrite,optic.BiDir,optic.Pure]{
return ORecipeOf[optic.Void,Recipe,Recipe,optic.ReturnOne](optic.Identity[Recipe]())
}
func OIngredientOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Ingredient,Ingredient,RET,RW,DIR,ERR])*lIngredient[I,S,T,RET,RW,DIR,ERR]{
return &lIngredient[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
func ORecipeOf[I any,S any,T any,RET any,RW any,DIR any,ERR any](l optic.Optic[I,S,T,Recipe,Recipe,RET,RW,DIR,ERR])*lRecipe[I,S,T,RET,RW,DIR,ERR]{
return &lRecipe[I,S,T,RET,RW,DIR,ERR]{
	Optic:	l,
}
}
var O  = o{
}
