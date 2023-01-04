package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type producto struct {
	Id           int
	Name         string
	Quantity     float64
	Code_value   string
	Is_published bool
	Expiration   string
	Price        float64
}

var Productos = make([]producto, 0)

func main() {
	data, err := ioutil.ReadFile("products.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(data, &Productos)
	if err != nil {
		fmt.Println(err)
		return
	}

}
