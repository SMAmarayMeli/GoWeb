package routes

import (
	"GoWeb/cmd/handlers"
	"GoWeb/internal/domain"
	"GoWeb/internal/product"
	"github.com/gin-gonic/gin"
)

type Router struct {
	db *[]domain.Producto
	en *gin.Engine
}
func NewRouter(en *gin.Engine, db *[]domain.Producto) *Router {
	return &Router{en: en, db: db}
}

func (r *Router) SetRoutes() {
	r.SetWebsite()
}
// website
func (r *Router) SetWebsite() {
	// instances
	rp := product.NewRepository(r.db, 3)
	sv := product.NewService(rp)
	h := handlers.NewProducto(sv)

	prod := r.en.Group("/products")

	prod.GET("/ping", h.Ping())
	prod.GET("/", h.Products())
	prod.GET("/:id", h.ProductId())
	prod.GET("/search", h.ProductsPriceGt())
	prod.POST("/", h.ProductAdd())
}