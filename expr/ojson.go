package expr

import (
	"fmt"

	"github.com/spearson78/go-optic/internal/util"
)

type ParseJson struct {
	OpticTypeExpr
}

func (e ParseJson) Short() string {
	return "ParseJson"
}

func (e ParseJson) String() string {
	return fmt.Sprintf("ParseJson(%v,%v)", util.FullTypeName(e.opticS), util.FullTypeName(e.opticA))
}

type JqObjectsOf struct {
	OpticTypeExpr
	Fields OpticExpression
}

func (e JqObjectsOf) Short() string {
	return fmt.Sprintf("JqObjectsOf(%v)", e.Fields.Short())
}

func (e JqObjectsOf) String() string {
	return fmt.Sprintf("JqObjectsOf(%v)", e.Fields.String())
}

type JqPick struct {
	OpticTypeExpr
	Paths [][]any
}

func (e JqPick) Short() string {
	return fmt.Sprintf("JqPick(%v)", e.Paths)
}

func (e JqPick) String() string {
	return fmt.Sprintf("JqPick(%v)", e.Paths)
}

type JqOrderBy struct {
	OpticTypeExpr
}

func (e JqOrderBy) Short() string {
	return "JqOrderBy"
}

func (e JqOrderBy) String() string {
	return "JqOrderBy"
}

type JqContains struct {
	OpticTypeExpr
}

func (e JqContains) Short() string {
	return "JqContains"
}

func (e JqContains) String() string {
	return "JqContains"
}

type Children struct {
	OpticTypeExpr
}

func (e Children) Short() string {
	return "Children"
}

func (e Children) String() string {
	return "Children"
}
