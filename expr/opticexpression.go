package expr

import (
	"fmt"
	"reflect"

	"github.com/spearson78/go-optic/internal/util"
)

type OpticExpression interface {
	OpticI() reflect.Type
	OpticS() reflect.Type
	OpticT() reflect.Type
	OpticA() reflect.Type
	OpticB() reflect.Type

	// Provides this optics internal type. This indicates which methods have efficient implementations.
	OpticType() OpticType

	Pure() bool

	//Provides a string representation of this Optic for diagnosis/logging purposes.
	String() string
	Short() string
}

type OpticType byte

const (
	OpticTypeReturnManyFlag OpticType = 1 << iota
	OpticTypeReadWriteFlag
	OpticTypeBiDirFlag
	OpticTypeIdentityFlag

	opticTypeReturnOneFlag OpticType = 0
	opticTypeReadOnlyFlag  OpticType = 0
	opticTypeUniDirFlag    OpticType = 0

	OpticTypeGetter    = opticTypeReturnOneFlag | opticTypeReadOnlyFlag | opticTypeUniDirFlag
	OpticTypeLens      = opticTypeReturnOneFlag | OpticTypeReadWriteFlag | opticTypeUniDirFlag
	OpticTypeIteration = OpticTypeReturnManyFlag | opticTypeReadOnlyFlag | opticTypeUniDirFlag
	OpticTypeTraversal = OpticTypeReturnManyFlag | OpticTypeReadWriteFlag | opticTypeUniDirFlag
	OpticTypeIso       = opticTypeReturnOneFlag | OpticTypeReadWriteFlag | OpticTypeBiDirFlag
	OpticTypePrism     = OpticTypeReturnManyFlag | OpticTypeReadWriteFlag | OpticTypeBiDirFlag
	OpticTypeIdentity  = opticTypeReturnOneFlag | OpticTypeReadWriteFlag | OpticTypeBiDirFlag | OpticTypeIdentityFlag
)

func (o OpticType) String() string {

	switch o {
	case OpticTypeGetter:
		return "Getter"
	case OpticTypeLens:
		return "Lens"
	case OpticTypeIteration:
		return "Iteration"
	case OpticTypeTraversal:
		return "Traversal"
	case OpticTypeIso:
		return "Iso"
	case OpticTypePrism:
		return "Prism"
	case OpticTypeIdentity:
		return "Identity"
	default:
		ret := ""

		if o&OpticTypeReturnManyFlag != 0 {
			ret += "ReturnMany"
		} else {
			ret += "ReturnOne"
		}

		if o&OpticTypeReadWriteFlag != 0 {
			ret += "ReadWrite"
		} else {
			ret += "ReadOnly"
		}

		if o&OpticTypeBiDirFlag != 0 {
			ret += "BiDir"
		} else {
			ret += "Uni"
		}

		return ret
	}
}

type OpticTypeExpr struct {
	opticI    reflect.Type
	opticS    reflect.Type
	opticT    reflect.Type
	opticA    reflect.Type
	opticB    reflect.Type
	opticType OpticType
	pure      bool
}

func (o OpticTypeExpr) OpticI() reflect.Type {
	return o.opticI
}

func (o OpticTypeExpr) OpticS() reflect.Type {
	return o.opticS
}

func (o OpticTypeExpr) OpticT() reflect.Type {
	return o.opticT
}
func (o OpticTypeExpr) OpticA() reflect.Type {
	return o.opticA
}

func (o OpticTypeExpr) OpticB() reflect.Type {
	return o.opticB
}

// Provides this optics internal type. This indicates which methods have efficient implementations.
func (o OpticTypeExpr) OpticType() OpticType {
	return o.opticType
}

func (o OpticTypeExpr) Pure() bool {
	return o.pure
}

func (o OpticTypeExpr) Signature() string {
	err := "Err"
	if o.pure {
		err = "Pure"
	}

	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v", util.FullTypeName(o.opticI), util.FullTypeName(o.opticS), util.FullTypeName(o.opticT), util.FullTypeName(o.opticA), util.FullTypeName(o.opticB), o.opticType, err)
}

func GetOpticType[RET, RW, DIR any]() OpticType {
	var opticType OpticType

	if !reflect.TypeFor[RET]().Comparable() {
		opticType |= OpticTypeReturnManyFlag
	}

	if reflect.TypeFor[RW]().Comparable() {
		opticType |= OpticTypeReadWriteFlag
	}

	if reflect.TypeFor[DIR]().Comparable() {
		opticType |= OpticTypeBiDirFlag
	}

	return opticType
}

func GetPurity[ERR any]() bool {
	return reflect.TypeFor[ERR]().Comparable()
}

func NewOpticTypeExpr[I, S, T, A, B, RET, RW, DIR, ERR any]() OpticTypeExpr {

	return OpticTypeExpr{
		opticI:    reflect.TypeFor[I](),
		opticS:    reflect.TypeFor[S](),
		opticT:    reflect.TypeFor[T](),
		opticA:    reflect.TypeFor[A](),
		opticB:    reflect.TypeFor[B](),
		opticType: GetOpticType[RET, RW, DIR](),
		pure:      GetPurity[ERR](),
	}
}
