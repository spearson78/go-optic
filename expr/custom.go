package expr

import "fmt"

type CustomOptic struct {
	OpticTypeExpr
	Id string
}

func (e CustomOptic) Short() string {
	return fmt.Sprintf("Custom(%v)", e.Id)
}

func (e CustomOptic) String() string {
	return fmt.Sprintf("Custom(%v)", e.Id)
}

type Unknown struct{}

func Custom(id string) func(ot OpticTypeExpr) OpticExpression {
	return func(ot OpticTypeExpr) OpticExpression {
		return CustomOptic{
			OpticTypeExpr: ot,
			Id:            id,
		}
	}
}

func TODO(id string) func(ot OpticTypeExpr) OpticExpression {
	return func(ot OpticTypeExpr) OpticExpression {
		return CustomOptic{
			OpticTypeExpr: ot,
			Id:            id,
		}
	}
}
