package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type chillybinClient struct {
	base string
}

func (c chillybinClient) pick(ingredient string, quantity int) error {
	resp, err := http.Get(fmt.Sprintf("%s/pick?name=%s&quantity=%d", c.base, ingredient, quantity))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Could not retrieve ingredients from chilly bin - '%s'", resp.Status)
		return errors.New("Could not retrieve ingredients from chilly bin")
	}
	return nil
}
