package model

//go:generate ../../../makelens model model.go model_generated.go

type Difficulty int

const (
	Easy Difficulty = iota
	Intermediate
	Expert
)

type Unit int

const (
	Gram Unit = iota
	Litre
	Slices
	Of
)

type Recipe struct {
	name         string
	description  string
	serves       int
	cookTime     float64
	difficulty   Difficulty
	ingredients  []Ingredient
	instructions []string
}

type Ingredient struct {
	name     string
	quantity float64
	unit     Unit
}
