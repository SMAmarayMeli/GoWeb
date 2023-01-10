package handlers

import (
	"GoWeb/internal/domain"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func Ping(c *gin.Context) {
	c.String(200, "pong")
}

func Products(c *gin.Context) {
	productos := domain.Productos
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
	var searched domain.Producto
	for _, a := range domain.Productos {
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
	var productosQueried = make([]domain.Producto, 0)
	for _, w := range domain.Productos {
		if priceQuery != 0 && w.Price > priceQuery {
			productosQueried = append(productosQueried, w)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "okay",
		"data":    productosQueried,
	})
}

func existeCodeValue(pro domain.Producto) error {
	for _, w := range domain.Productos {
		if w.CodeValue == pro.CodeValue {
			return errors.New("Ya existe ese CodeValue")
		}
	}
	return nil
}

func verificarFecha(date string) error {
	layout := "02/01/2006"
	_, err := time.Parse(layout, date)
	return err
}

func verificarVacios(pro domain.Producto) error {
	if pro.Price == 0 {
		return errors.New("Price no puede estar vacio")
	}
	if pro.Name == "" {
		return errors.New("Name no puede estar vacio")
	}
	if pro.Expiration == "" {
		return errors.New("Expiration no puede estar vacio")
	}
	if pro.CodeValue == "" {
		return errors.New("CodeValue no puede estar vacio")
	}
	if pro.Quantity == 0 {
		return errors.New("Quantity no puede estar vacio")
	}
	return nil
}

func ProductAdd(c *gin.Context) {
	var pro domain.Producto

	if err := c.ShouldBindJSON(&pro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	if errCodeValue := existeCodeValue(pro); errCodeValue != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": errCodeValue.Error(),
		})
		return
	}

	if errFecha := verificarFecha(pro.Expiration); errFecha != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errFecha,
		})
		return
	}

	if errVacios := verificarVacios(pro); errVacios != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errVacios.Error(),
		})
		return
	}

	domain.LastId++
	pro.Id = domain.LastId
	domain.Productos = append(domain.Productos, pro)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Created ok",
		"data":    pro,
	})
}
