package main

import (
	"GoWeb/cmd/routes"
	"GoWeb/internal/domain"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
)

func readJson(dbP *[]domain.Producto) {
	data, err := ioutil.ReadFile("../GoWeb/internal/product/products.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(data, dbP)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	dbP := []domain.Producto{}
	readJson(&dbP)

	en := gin.Default()
	rt := routes.NewRouter(en, &dbP)
	rt.SetRoutes()

	if err := en.Run(); err != nil {
		log.Fatal(err)
	}

}
