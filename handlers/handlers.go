package handlers

import (
	"GoWeb/globals"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Ping(c *gin.Context) {
	c.String(200, "pong")
}

func Products(c *gin.Context) {
	productos := globals.Productos
	c.JSON(http.StatusOK, gin.H{
		"products": productos,
	})
}

func ProductId(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to parse id",
			"data":    nil,
		})
		return
	}
	var searched globals.Producto
	for _, a := range globals.Productos {
		if a.Id == n {
			searched = a
			break
		}
	}

	if searched.Id != 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully found ID",
			"data":    searched,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id not found",
			"data":    nil,
		})
		return
	}
}

func ProductsPriceGt(c *gin.Context) {
	priceQuery, err := strconv.ParseFloat(c.Query("price"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Couldnt parse number",
			"data":    nil,
		})
		return
	}
	var productosQueried = make([]globals.Producto, 0)
	for _, w := range globals.Productos {
		if priceQuery != 0 && w.Price > priceQuery {
			productosQueried = append(productosQueried, w)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "okay",
		"data":    productosQueried,
	})
}
