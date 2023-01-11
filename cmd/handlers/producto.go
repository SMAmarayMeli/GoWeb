package handlers

import (
	"GoWeb/internal/domain"
	"GoWeb/internal/product"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type requestProduct struct {
	Name        string  `json:"name"`
	Quantity    float64 `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published" `
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type Producto struct {
	sv product.Service
}

func NewProducto(sv product.Service) *Producto {
	return &Producto{sv: sv}
}

func (p *Producto) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "pong")
	}
}

func (p *Producto) Products() gin.HandlerFunc {
	return func(c *gin.Context) {
		productos, err := p.sv.Get()
		if err != nil {
			c.JSON(500, nil)
			return
		}
		c.JSON(http.StatusOK, productos)
	}
}

func (p *Producto) ProductId() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "fail to parse id",
				"error": err,
			})
			return
		}

		searched, err := p.sv.GetById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, err)
			return
		}

		c.JSON(http.StatusFound, *searched)
	}
}

func (p *Producto) ProductsPriceGt() gin.HandlerFunc {
	return func(c *gin.Context) {
		priceQuery, err := strconv.ParseFloat(c.Query("price"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Couldn't parse number",
			})
			return
		}

		productosQueried, err1 := p.sv.GetGreaterThanPrice(priceQuery)
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, "")
			return
		}

		c.JSON(http.StatusOK, productosQueried)
	}
}

func verificarFecha(date string) error {
	layout := "02/01/2006"
	_, err := time.Parse(layout, date)
	return err
}

func verificarVacios(price float64, name string, expiration string, codeValue string, quantity float64) error {
	if price == 0 {
		return errors.New("Price no puede estar vacio")
	}
	if name == "" {
		return errors.New("Name no puede estar vacio")
	}
	if expiration == "" {
		return errors.New("Expiration no puede estar vacio")
	}
	if codeValue == "" {
		return errors.New("CodeValue no puede estar vacio")
	}
	if quantity == 0 {
		return errors.New("Quantity no puede estar vacio")
	}
	return nil
}

func (p *Producto) ProductAdd() gin.HandlerFunc{
	return func(c *gin.Context) {

		var r requestProduct

		if err := c.ShouldBindJSON(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error json": err,
			})
		}

		if errFecha := verificarFecha(r.Expiration); errFecha != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error fecha": errFecha,
			})
			return
		}

		if errVacios := verificarVacios(r.Price, r.Name, r.Expiration, r.CodeValue, r.Quantity); errVacios != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error vacios": errVacios.Error(),
			})
			return
		}

		pr, errCreate := p.sv.Create(r.Name, r.CodeValue, r.Expiration, r.Quantity, r.Price, r.IsPublished)
		if errCreate != nil {
			c.JSON(http.StatusInternalServerError, errCreate)
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Created ok",
			"data":    pr,
		})
	}

}

func (p *Producto) ProductReplace() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "fail to parse code",
				"error": err,
			})
			return
		}

		var replaceRequest requestProduct
		if err := c.ShouldBindJSON(&replaceRequest); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		var replaceProduct = domain.Producto{
			Id: id,
			Name: replaceRequest.Name,
			Quantity: replaceRequest.Quantity,
			CodeValue: replaceRequest.CodeValue,
			IsPublished: replaceRequest.IsPublished,
			Expiration: replaceRequest.Expiration,
			Price: replaceRequest.Price,
		}

		err = p.sv.Update(id, replaceProduct)
		if err != nil {
			switch err {
			case product.ErrNotFound:
				c.JSON(http.StatusNotFound, "No se encontro")
			default:
				c.JSON(http.StatusInternalServerError, nil)
			}
		}

		c.JSON(http.StatusOK, replaceProduct)
	}
}