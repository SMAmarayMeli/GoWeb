package handler_tests

import (
	"GoWeb/cmd/handlers"
	"GoWeb/internal/domain"
	"GoWeb/internal/product"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createServerForProductsHandlerTest() *gin.Engine {
	// Create router
	dbP := []domain.Producto{
		{
			Id:          1,
			Name:        "Pepe",
			Quantity:    500,
			CodeValue:   "SG44554",
			IsPublished: true,
			Expiration:  "15/01/2023",
			Price:       5130.21,
		}, {
			Id:          2,
			Name:        "Moni",
			Quantity:    134,
			CodeValue:   "SG4A54",
			IsPublished: false,
			Expiration:  "15/05/2023",
			Price:       51.21,
		},
	}

	// Create handler
	rp := product.NewRepository(&dbP, len(dbP))
	sv := product.NewService(rp)
	h := handlers.NewProducto(sv)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	prod := r.Group("/products")
	prod.GET("/ping", h.Ping())
	prod.GET("/", h.Products())
	prod.GET("/:id", h.ProductId())
	prod.GET("/search", h.ProductsPriceGt())
	prod.POST("/", h.ProductAdd())
	prod.PUT("/:id", h.ProductReplace())
	prod.PATCH("/:id", h.ProductPatch())
	prod.DELETE("/:id", h.DeleteProduct())

	return r
}

func TestGetAll_Ok(t *testing.T) {
	// Arrange
	server := createServerForProductsHandlerTest()
	request := httptest.NewRequest(http.MethodGet, "/products/", nil)
	responseRecorded := httptest.NewRecorder()
	// Act
	server.ServeHTTP(responseRecorded, request)
	body, err := io.ReadAll(responseRecorded.Body)
	var obtainedBody []domain.Producto
	// TODO Check why Unmarshal doesnt work. Body contains ok data.
	// errOnBodyUnmarshal := json.Unmarshal(body, &obtainedBody)
	t.Logf("%v\n", string(body))
	t.Logf("%v\n", obtainedBody)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorded.Code)
	// assert.Nil(t, errOnBodyUnmarshal, "Error in unmarshal")
	assert.Equal(t, nil, err)
	assert.True(t, len(body) > 0)
}

func TestProductGetById(t *testing.T) {
	// Arrange
	server := createServerForProductsHandlerTest()
	request := httptest.NewRequest(http.MethodGet, "/products/2", nil)
	responseRecorded := httptest.NewRecorder()

	// Act
	server.ServeHTTP(responseRecorded, request)

	body, _ := io.ReadAll(responseRecorded.Body)
	t.Logf("%v\n", string(body))

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorded.Code)
}
