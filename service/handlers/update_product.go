package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/osag1e/product-crud/internal/model"
)

func (p *ProductHandler) HandleUpdateProductByID(w http.ResponseWriter, r *http.Request) {
	productIDstr := r.URL.Query().Get("productID")

	productID, err := uuid.Parse(productIDstr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var product *model.Products
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	_, err = p.productsRepo.UpdateProductByID(productID, product)
	if err != nil {
		http.Error(w, "Product Not Updated", http.StatusBadRequest)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Product has been updated"})
}
