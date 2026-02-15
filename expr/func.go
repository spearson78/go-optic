package expr

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/spearson78/go-optic/internal/util"
)

type GoFuncExpr struct {
	OpticTypeExpr
	Func reflect.Value
}

func (e GoFuncExpr) Short() string {
	f := runtime.FuncForPC(e.Func.Pointer())
	if f != nil {
		name := f.Name()
		split := strings.Split(name, ".")
		return split[len(split)-1]
	} else {
		return "func(...){}"
	}
}

func (e GoFuncExpr) String() string {
	return fmt.Sprintf("GoFunc(%v,%v)", util.FullTypeName(e.opticS), util.FullTypeName(e.opticA))
}
