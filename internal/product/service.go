package product

import (
	"GoWeb/internal/domain"
	"errors"
)

var (
	ErrAlreadyExist = errors.New("already exist")
)

type Service interface {
	Get() ([]domain.Producto, error)
	GetById(id int) (domain.Producto, error)
	GetGreaterThanPrice(price float64) ([]domain.Producto, error)
	Update(id int, product domain.Producto) (domain.Producto, error)
	Create(name, codeValue, expiration string, quantity, price float64, isPublished bool) (domain.Producto, error)
	Delete(id int) error
}

func NewService(rp Repository) Service {
	return &service{rp: rp}
}

type service struct {
	rp Repository
}

func (sv *service) Get() ([]domain.Producto, error) {
	return sv.rp.Get()
}

func (sv *service) GetById(id int) (domain.Producto, error) {
	return sv.rp.GetById(id)
}

func (sv *service) GetGreaterThanPrice(price float64) ([]domain.Producto, error) {
	return sv.rp.GetGreaterThanPrice(price)
}

func (sv *service) Update(id int, product domain.Producto) (domain.Producto, error) {
	p, err := sv.rp.GetById(id)
	if err != nil {
		return domain.Producto{}, err
	}
	if product.Name != "" {
		p.Name = product.Name
	}
	if product.CodeValue != "" {
		p.CodeValue = product.CodeValue
	}
	if product.Expiration != "" {
		p.Expiration = product.Expiration
	}
	if product.Quantity > 0 {
		p.Quantity = product.Quantity
	}
	if product.Price > 0 {
		p.Price = product.Price
	}

	p, err = sv.rp.Update(id, p)
	if err != nil {
		return domain.Producto{}, err
	}
	return p, nil
}

func (sv *service) Delete(id int) error {
	err := sv.rp.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (sv *service) Create(name, codeValue, expiration string, quantity, price float64, isPublished bool) (domain.Producto, error) {

	pr := domain.Producto{
		Name: name,
		Quantity: quantity,
		IsPublished: isPublished,
		Expiration: expiration,
		Price: price,
	}

	lastId, err := sv.rp.Create(pr)
	if err != nil {
		return domain.Producto{}, err
	}

	pr.Id = lastId

	return pr, nil
}