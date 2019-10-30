package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type chillybinClient struct {
	base string
}

func (c chillybinClient) pickV2(ingredient string, quantity int, customer string) error {
	input, err := json.Marshal(struct {
		Ingredient string
		Quantity   int
		Customer   string
	}{
		Ingredient: ingredient,
		Quantity:   quantity,
		Customer:   customer,
	})
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/v2/pick", c.base), "application/json", bytes.NewBuffer(input))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Could not retrieve ingredients from chilly bin - '%s'", resp.Status)
		return errors.New("Could not retrieve ingredients from chilly bin")
	}
	return nil
}
