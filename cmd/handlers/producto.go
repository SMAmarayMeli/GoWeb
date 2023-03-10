package handlers

import (
	"GoWeb/internal/domain"
	"GoWeb/internal/product"
	"GoWeb/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
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
		c.JSON(200, response.Response{Data: "pong"})
	}
}

func (p *Producto) Products() gin.HandlerFunc {
	return func(c *gin.Context) {
		productos, err := p.sv.Get()
		if err != nil {
			c.JSON(500, nil)
			return
		}
		c.JSON(http.StatusOK, response.Response{Data: productos})
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

		c.JSON(http.StatusFound, response.Response{Data: searched})
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

		c.JSON(http.StatusOK, response.Response{Data: productosQueried})
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

		token := c.GetHeader("token")
		if token != os.Getenv("SECRET"){
			c.JSON(http.StatusUnauthorized, "Token invalido")
			return
		}

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

		c.JSON(http.StatusCreated, response.Response{Data: pr})
	}

}

func (p *Producto) ProductReplace() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("token")
		if token != os.Getenv("SECRET"){
			c.JSON(http.StatusUnauthorized, "Token invalido")
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "fail to parse id",
				"error": err,
			})
			return
		}

		var replaceRequest requestProduct
		if err := c.ShouldBindJSON(&replaceRequest); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		if errVacios := verificarVacios(replaceRequest.Price, replaceRequest.Name, replaceRequest.Expiration, replaceRequest.CodeValue, replaceRequest.Quantity); errVacios != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error vacios": errVacios.Error(),
			})
			return
		}

		if errFecha := verificarFecha(replaceRequest.Expiration); errFecha != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error fecha": errFecha,
			})
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

		prodNew, errReplace := p.sv.Update(id, replaceProduct)
		if err != nil {
			switch errReplace {
			case product.ErrNotFound:
				c.JSON(http.StatusNotFound, "No se encontro")
			default:
				c.JSON(http.StatusInternalServerError, nil)
			}
		}

		c.JSON(http.StatusOK, response.Response{Data: prodNew})
	}
}

func (p *Producto) DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != os.Getenv("SECRET"){
			c.JSON(http.StatusUnauthorized, "Token invalido")
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "fail to parse id",
				"error": err,
			})
			return
		}
		err = p.sv.Delete(id)
		if err != nil {
			switch err{
			case product.ErrNotFound:
				c.JSON(http.StatusNotFound, err)
			default:
				c.JSON(http.StatusInternalServerError, "")
			}
		}
		c.JSON(http.StatusOK, response.Response{Data: "Product deleted"})
	}
}

func (p *Producto) ProductPatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != os.Getenv("SECRET"){
			c.JSON(http.StatusUnauthorized, "Token invalido")
			return
		}

		var r requestProduct

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "fail to parse id",
				"error": err,
			})
			return
		}
		if err = c.ShouldBindJSON(&r); err != nil {
			c.JSON(http.StatusBadRequest, err)
		}

		update := domain.Producto{
			Name:        r.Name,
			Quantity:    r.Quantity,
			CodeValue:   r.CodeValue,
			IsPublished: r.IsPublished,
			Expiration:  r.Expiration,
			Price:       r.Price,
		}

		if update.Expiration != "" {
			if errFecha := verificarFecha(update.Expiration); errFecha != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error fecha": errFecha,
				})
				return
			}
		}

		lastP, errUpdate := p.sv.Update(id, update)
		if errUpdate != nil {
			c.JSON(http.StatusConflict, err)
		}
		c.JSON(http.StatusOK, response.Response{Data: lastP})
	}
}