package product

import (
	"GoWeb/internal/domain"
	"database/sql"
	"errors"
	"log"
)

var (
	ErrNotFoundMySql = errors.New("item not found")
)

type RepositoryMySql interface {
	// read
	Get() ([]domain.Producto, error)
	GetById(id int) (domain.Producto, error)
	validateCodeValue(codeValue string) bool
	ExistId(url int) bool
	// write
	Delete(id int) error
	Create(domain.Producto) (int, error)
	Update(id int, product domain.Producto) (domain.Producto, error)
}

type repositoryMySql struct {
	db *sql.DB
}

func NewRepositoryMySql(db *sql.DB) RepositoryMySql {
	return &repositoryMySql{db: db}
}

func (r *repositoryMySql) Get() ([]domain.Producto, error) {
	var products []domain.Producto
	db := r.db
	rows, err := db.Query("SELECT id, name, quantity, code_value, is_published, expiration, price FROM products")
	if err != nil {
		log.Println(err.Error())
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var product domain.Producto
		err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Quantity,
			&product.CodeValue,
			&product.IsPublished,
			&product.Expiration,
			&product.Price)
		if err != nil {
			log.Println(err.Error())
			return products, err
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		log.Println(err.Error())
		return products, err
	}
	return products, nil
}
func (r *repositoryMySql) GetById(id int) (domain.Producto, error) {

	var product domain.Producto
	db := r.db
	row := db.QueryRow("SELECT id, name, quantity, code_value, is_published, expiration, price FROM products WHERE ID = ?", id)
	if row.Err() != nil {

	}
	err := row.Scan(
		&product.Id,
		&product.Name,
		&product.Quantity,
		&product.CodeValue,
		&product.IsPublished,
		&product.Expiration,
		&product.Price)
	if err != nil {
		log.Println(err.Error())
		return product, err
	}
	return product, nil
}
func (r *repositoryMySql) validateCodeValue(codeValue string) bool {
	//for _, p := range *r.db {
	//	if p.CodeValue == codeValue {
	//		return false
	//	}
	//}
	return true
}
func (r *repositoryMySql) ExistId(id int) bool {
	//for _, p := range *r.db {
	//	if p.Id == id {
	//		return true
	//	}
	//}
	return false
}

// write
func (r *repositoryMySql) Create(product domain.Producto) (int, error) {
	db := r.db // se inicializa la base
	query := "INSERT INTO products(id, name, quantity, code_value, is_published, expiration, price) VALUES( ?, ?, ?, ?)"
	prep, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer prep.Close()

	var maxID int
	err = db.QueryRow("SELECT MAX(id) FROM products").Scan(&maxID)
	if err != nil {
		return 0, err
	}
	product.Id = maxID + 1
	_, err = prep.Exec(product.Id, product.Name, product.Quantity, product.CodeValue, product.IsPublished, product.Expiration, product.Price) // retorna un sql.Result y un error
	if err != nil {
		return 0, err
	}

	return maxID, nil
}

func (r *repositoryMySql) Delete(id int) error {

	return ErrNotFound
}

func (r *repositoryMySql) Update(id int, product domain.Producto) (domain.Producto, error) {
	db := r.db // se inicializa la base
	query := "UPDATE products SET name = ?, quantity = ?, code_value = ?, is_published = ?, expiration = ?, price = ? WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.Name, product.Quantity, product.CodeValue, product.IsPublished, product.Expiration, product.Price)
	if err != nil {
		return domain.Producto{}, err
	}
	return product, nil
}