package handlers

import (
	"database/sql"

	"github.com/osag1e/product-crud/internal/query"
)

type ProductHandler struct {
	DB           *sql.DB
	productsRepo query.ProductRepository
}

func NewProductHandler(productsRepo query.ProductRepository) *ProductHandler {
	return &ProductHandler{
		productsRepo: productsRepo,
	}
}
