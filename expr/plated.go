package expr

import (
	"fmt"
)

type TopDownMatch struct {
	OpticTypeExpr
	Children OpticExpression
	Matcher  OpticExpression
}

func (e TopDownMatch) Short() string {
	return fmt.Sprintf("Deep(%v)", e.Matcher.Short())
}

func (e TopDownMatch) String() string {
	return fmt.Sprintf("Deep(%v)", e.Matcher.String())
}

type Rewrite struct {
	OpticTypeExpr

	Children  OpticExpression
	RewriteOp OpticExpression
}

func (e Rewrite) Short() string {
	return fmt.Sprintf("Rewrite(%v)", e.RewriteOp.Short())
}

func (e Rewrite) String() string {
	return fmt.Sprintf("Rewrite(%v)", e.RewriteOp.String())
}
