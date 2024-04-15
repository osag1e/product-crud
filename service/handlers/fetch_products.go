package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

func (p *ProductHandler) HandleFetchProducts(w http.ResponseWriter, r *http.Request) {
	productTypes, err := p.productsRepo.GetProducts(0, 11)
	if err != nil {
		http.Error(w, "Get Product Request Failed", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{"data": productTypes})
}

func (p *ProductHandler) HandleFetchProductByID(w http.ResponseWriter, r *http.Request) {
	productIDstr := r.URL.Query().Get("productID")

	productID, err := uuid.Parse(productIDstr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	product, err := p.productsRepo.GetProductByID(productID)
	if err != nil {
		http.Error(w, "Get Product Request Failed", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{"data": product})
}
