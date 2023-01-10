package product

import (
	"GoWeb/internal/domain"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ReadJson() {
	data, err := ioutil.ReadFile("products.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(data, &domain.Productos)
	if err != nil {
		fmt.Println(err)
		return
	}
}