package expr

import (
	"fmt"

	"github.com/spearson78/go-optic/internal/util"
)

type TraverseMembers struct {
	OpticTypeExpr
}

func (e TraverseMembers) Short() string {
	return "TraverseMembers"
}

func (e TraverseMembers) String() string {
	return fmt.Sprintf("TraverseMembers(%v,%v)", util.FullTypeName(e.opticS), util.FullTypeName(e.opticA))
}
