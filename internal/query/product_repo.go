package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/osag1e/product-crud/internal/model"
)

type ProductRepository interface {
	InsertProduct(product *model.Products) (*model.Products, error)
	GetProducts(int, int) ([]model.Products, error)
	GetProductByID(productID uuid.UUID) (*model.Products, error)
	UpdateProductByID(productID uuid.UUID, updatedProduct *model.Products) (*model.Products, error)
	DeleteProductByID(productID uuid.UUID) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type ProductStore struct {
	DB *sql.DB
}

func NewProductStore(db *sql.DB) ProductRepository {
	return &ProductStore{DB: db}
}

func (ps *ProductStore) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return ps.DB.BeginTx(ctx, opts)
}

func (ps *ProductStore) InsertProduct(product *model.Products) (*model.Products, error) {
	query := `
	INSERT INTO store.products (id, brand, description, colour, size, price, sku) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	product.ID = model.NewUUID()

	_, err := ps.DB.Exec(query, product.ID, product.Brand, product.Description, product.Colour, product.Size, product.Price, product.SKU)
	if err != nil {
		return nil, fmt.Errorf("failed to insert product: %v", err)
	}
	return product, nil
}

func (ps *ProductStore) GetProducts(limit int, offset int) ([]model.Products, error) {
	query := `
	          SELECT id, brand, description, colour, size, price, sku 
	          FROM store.products ORDER BY id OFFSET $1 LIMIT $2
			  `
	rows, err := ps.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productTypes []model.Products

	for rows.Next() {
		var product model.Products
		if err := rows.Scan(&product.ID, &product.Brand, &product.Description, &product.Colour, &product.Size, &product.Price, &product.SKU); err != nil {
			return nil, err
		}
		productTypes = append(productTypes, product)
	}

	return productTypes, nil
}

func (ps *ProductStore) GetProductByID(productID uuid.UUID) (*model.Products, error) {
	var product model.Products
	query := ` 
	SELECT id, brand, description, colour, size, price, sku 
	FROM store.products 
	WHERE id = $1
	`
	stmt, err := ps.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(productID).Scan(
		&product.ID,
		&product.Brand,
		&product.Description,
		&product.Colour,
		&product.Size,
		&product.Price,
		&product.SKU,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no default config found with ID: %s", productID)
		}
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	return &product, nil
}

func (ps *ProductStore) UpdateProductByID(productID uuid.UUID, updatedProduct *model.Products) (*model.Products, error) {
	query := `
	UPDATE store.products SET brand = $1, description = $2, colour = $3, size = $4, price = $5, sku = $6 
	WHERE id = $7
	`
	result, err := ps.DB.Exec(query, updatedProduct.Brand, updatedProduct.Description, updatedProduct.Colour,
		updatedProduct.Size, updatedProduct.Price, updatedProduct.SKU, productID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no product found with ID %s", productID)
	}

	return updatedProduct, err
}

func (ps *ProductStore) DeleteProductByID(productID uuid.UUID) error {
	query := `
	DELETE FROM store.products WHERE id = $1	
	`
	_, err := ps.DB.Exec(query, productID)
	return err
}
