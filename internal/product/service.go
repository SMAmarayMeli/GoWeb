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
	GetById(id int) (*domain.Producto, error)
	GetGreaterThanPrice(price float64) ([]domain.Producto, error)
	Update(id int, product domain.Producto) error
	Create(name, codeValue, expiration string, quantity, price float64, isPublished bool) (domain.Producto, error)
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

func (sv *service) GetById(id int) (*domain.Producto, error) {
	return sv.rp.GetById(id)
}

func (sv *service) GetGreaterThanPrice(price float64) ([]domain.Producto, error) {
	return sv.rp.GetGreaterThanPrice(price)
}

func (sv *service) Update(id int, product domain.Producto) error {
	return sv.rp.Update(id, product)
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