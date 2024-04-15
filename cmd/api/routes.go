package main

import (
	"database/sql"
	"net/http"

	"github.com/osag1e/product-crud/db/health"
	"github.com/osag1e/product-crud/internal/query"
	"github.com/osag1e/product-crud/service/handlers"
)

func initializeRouter(dbConn *sql.DB) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /health", health.PostgreSQLHealthCheckHandler(dbConn))

	productRepo := query.NewProductStore(dbConn)
	productHandler := handlers.NewProductHandler(productRepo)

	router.HandleFunc("POST /products", productHandler.HandleCreateProducts)
	router.HandleFunc("GET /products", productHandler.HandleFetchProducts)
	router.HandleFunc("GET /product", productHandler.HandleFetchProductByID)
	router.HandleFunc("PUT /product", productHandler.HandleUpdateProductByID)
	router.HandleFunc("DELETE /product", productHandler.HandleDeleteProductByID)

	return router
}
