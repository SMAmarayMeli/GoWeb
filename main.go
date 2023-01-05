package main

import (
	"GoWeb/globals"
	"GoWeb/handlers"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func readJson() {
	data, err := ioutil.ReadFile("products.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(data, &globals.Productos)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	readJson()
	globals.LastId = len(globals.Productos)

	router := gin.Default()

	router.GET("/ping", handlers.Ping)
	router.GET("/products", handlers.Products)
	router.GET("/products/:id", handlers.ProductId)
	router.GET("/products/search", handlers.ProductsPriceGt)
	router.POST("/products", handlers.ProductAdd)

	router.Run()
}
