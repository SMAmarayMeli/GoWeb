package main

import (
	"GoWeb/handlers"
	"GoWeb/internal/domain"
	"GoWeb/internal/product"
	"github.com/gin-gonic/gin"
)



func main() {
	product.ReadJson()
	domain.LastId = len(domain.Productos)

	router := gin.Default()

	prod := router.Group("/products")

	router.GET("/ping", handlers.Ping)
	prod.GET("/", handlers.Products)
	prod.GET("/:id", handlers.ProductId)
	prod.GET("/search", handlers.ProductsPriceGt)
	prod.POST("/", handlers.ProductAdd)

	router.Run()
}
