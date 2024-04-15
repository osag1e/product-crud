package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/osag1e/product-crud/internal/model"
)

func (p *ProductHandler) HandleCreateProducts(w http.ResponseWriter, req *http.Request) {
	var products []model.Products
	if err := json.NewDecoder(req.Body).Decode(&products); err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	maxConcurrent := 11
	concurrencyLimiter := make(chan struct{}, maxConcurrent)
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	tx, err := p.productsRepo.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	errorChannel := make(chan error, len(products))
	var wg sync.WaitGroup

	for _, product := range products {
		if product.Brand == "" || product.Price <= 0 {
			http.Error(w, "Create Product Request Failed: Invalid input data", http.StatusBadRequest)
			return
		}

		concurrencyLimiter <- struct{}{}
		wg.Add(1)
		go func(product model.Products) {
			defer func() { <-concurrencyLimiter }()
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				if _, err := p.productsRepo.InsertProduct(&product); err != nil {
					errorChannel <- err
				}
			}
		}(product)
	}

	wg.Wait()
	close(errorChannel)

	var errMsgs []string
	for err := range errorChannel {
		errMsgs = append(errMsgs, err.Error())
	}

	if len(errMsgs) == 0 {
		if err := tx.Commit(); err != nil {
			http.Error(w, "Transaction Commit Failed", http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Products have been created"})
	}
}
