package routes

import (
	"GoWeb/cmd/handlers"
	"GoWeb/internal/product"
	"database/sql"
	"github.com/gin-gonic/gin"
)

type Router struct {
	db *sql.DB
	en *gin.Engine
}

func NewRouter(db *sql.DB, en *gin.Engine) *Router {
	return &Router{en: en, db: db}
}

func (r *Router) SetRoutes() {
	r.SetWebsite()
}

// website
func (r *Router) SetWebsite() {
	// instances
	rp := product.NewRepository(r.db)
	sv := product.NewService(rp)
	h := handlers.NewProducto(sv)

	prod := r.en.Group("/products")

	prod.GET("/ping", h.Ping())
	prod.GET("/", h.Products())
	prod.GET("/:id", h.ProductId())
	prod.GET("/search", h.ProductsPriceGt())
	prod.POST("/", h.ProductAdd())
	prod.PUT("/:id", h.ProductReplace())
	prod.PATCH("/:id", h.ProductPatch())
	prod.DELETE("/:id", h.DeleteProduct())
}
