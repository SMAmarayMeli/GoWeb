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

	prod := router.Group("/products")

	router.GET("/ping", handlers.Ping)
	prod.GET("/", handlers.Products)
	prod.GET("/:id", handlers.ProductId)
	prod.GET("/search", handlers.ProductsPriceGt)
	prod.POST("/", handlers.ProductAdd)

	router.Run()
}
