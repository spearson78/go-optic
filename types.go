package optic

import (
	"golang.org/x/exp/constraints"
)

// Type constraint for types than be operated on by basic arithmetic operations like +,-./,*
// This type constraint includes the complex types. See [Real] for a type constraint without complex types.
type Arithmetic interface {
	constraints.Complex | Real
}

// Type constraint for types than be operated on by basic arithmetic operations like +,-./,*
// This type constraint excludes the complex types. See [Arithmetic] for a type constraint with complex types.
type Real interface {
	constraints.Float | constraints.Integer
}

// Type constraint for integer types
type Integer interface {
	constraints.Integer
}

// Empty type as a placeholder in indexed optics when no index is available.
type Void struct{}
