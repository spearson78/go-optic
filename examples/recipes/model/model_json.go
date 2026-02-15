package model

import "encoding/json"

func (r *Recipe) MarshalJSON() ([]byte, error) {
	//MarshalJSON needs exported fields so we copy to an exported version of the struct
	return json.Marshal(struct {
		Name         string
		Description  string
		Serves       int
		CookTime     float64
		Difficulty   Difficulty
		Ingredients  []Ingredient
		Instructions []string
	}{
		Name:         r.name,
		Description:  r.description,
		Serves:       r.serves,
		CookTime:     r.cookTime,
		Difficulty:   r.difficulty,
		Ingredients:  r.ingredients,
		Instructions: r.instructions,
	})
}

func (r *Recipe) UnmarshalJSON(data []byte) error {

	//UnmarshalJSON needs exported fields so we unmarshal to an exported version of the struct. Then copy to the unexported version.
	u := struct {
		Name         string
		Description  string
		Serves       int
		CookTime     float64
		Difficulty   Difficulty
		Ingredients  []Ingredient
		Instructions []string
	}{}

	err := json.Unmarshal(data, &u)

	r.name = u.Name
	r.description = u.Description
	r.serves = u.Serves
	r.cookTime = u.CookTime
	r.difficulty = u.Difficulty
	r.ingredients = u.Ingredients
	r.instructions = u.Instructions

	return err
}

func (i *Ingredient) MarshalJSON() ([]byte, error) {
	//MarshalJSON needs exported fields so we copy to an exported version of the struct
	return json.Marshal(struct {
		Name     string
		Quantity float64
		Unit     Unit
	}{
		Name:     i.name,
		Quantity: i.quantity,
		Unit:     i.unit,
	})
}

func (i *Ingredient) UnmarshalJSON(data []byte) error {

	//UnmarshalJSON needs exported fields so we unmarshal to an exported version of the struct. Then copy to the unexported version.
	u := struct {
		Name     string
		Quantity float64
		Unit     Unit
	}{}

	err := json.Unmarshal(data, &u)

	i.name = u.Name
	i.quantity = u.Quantity
	i.unit = u.Unit

	return err
}
