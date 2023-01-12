package domain

type Producto struct {
	Id          int     `json:"id,-"`
	Name        string  `json:"name"`
	Quantity    float64 `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published" `
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}