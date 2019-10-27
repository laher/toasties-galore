package main

import "testing"

func TestValidate(t *testing.T) {

	var cases = []struct {
		name        string
		ingredients []string
		err         error
	}{
		{"not enough ingredients", []string{}, notEnoughIngredients},
		{"no cheese", []string{"parsnips"}, noCheese},
		{"ok", []string{"cheese", "radish"}, nil},
	}
	for _, c := range cases {
		err := validate(c.ingredients)
		if err != c.err {
			t.Error("Error didnt match")
		}
	}
}
