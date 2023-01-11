package product

import (
	"GoWeb/internal/domain"
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("item not found")
)

type Repository interface {
	// read
	Get() ([]domain.Producto, error)
	GetById(id int) (domain.Producto, error)
	validateCodeValue(codeValue string) bool
	ExistId(url int) bool
	GetGreaterThanPrice(price float64) ([]domain.Producto, error)
	// write
	Delete(id int) error
	Create(domain.Producto) (int, error)
	Update(id int, product domain.Producto) (domain.Producto, error)
}

type repository struct {
	db *[]domain.Producto
	// config
	lastID	int
}

func NewRepository(db *[]domain.Producto, lastID int) Repository {
	return &repository{db: db, lastID: lastID}
}

func (r *repository) Get() ([]domain.Producto, error) {
	return *r.db, nil
}
func (r *repository) GetById(id int) (domain.Producto, error) {
	for _, p := range *r.db {
		if p.Id == id {
			return p, nil
		}
	}

	return domain.Producto{}, fmt.Errorf("%w. %s", ErrNotFound, "product does not exist")
}
func (r *repository) validateCodeValue(codeValue string) bool {
	for _, p := range *r.db {
		if p.CodeValue == codeValue {
			return false
		}
	}

	return true
}
func (r *repository) ExistId(id int) bool {
	for _, p := range *r.db {
		if p.Id == id {
			return true
		}
	}
	return false
}

func (r *repository) GetGreaterThanPrice(price float64) ([]domain.Producto, error)  {
	var productosQueried = make([]domain.Producto, 0)
	for _, w := range *r.db{
		if price != 0 && w.Price > price {
			productosQueried = append(productosQueried, w)
		}
	}
	return productosQueried, nil
}

// write
func (r *repository) Create(product domain.Producto) (int, error) {
	r.lastID++
	product.Id = r.lastID
	*r.db = append(*r.db, product)

	return r.lastID, nil
}

func (r *repository) Delete(id int) error {
	for i, product := range *r.db {
		if product.Id == id {
			*r.db = append((*r.db)[:i], (*r.db)[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (r *repository) Update(id int, product domain.Producto) (domain.Producto, error) {
	for i, p := range *r.db{
		if p.Id == id {
			if !r.validateCodeValue(product.CodeValue) && product.CodeValue != p.CodeValue {
				return domain.Producto{}, ErrAlreadyExist
			}
			(*r.db)[i] = product
			return product, nil
		}
	}
	return domain.Producto{}, ErrNotFound
}