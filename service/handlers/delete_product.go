package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

func (p *ProductHandler) HandleDeleteProductByID(w http.ResponseWriter, r *http.Request) {
	productIDstr := r.URL.Query().Get("productID")

	productID, err := uuid.Parse(productIDstr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	err = p.productsRepo.DeleteProductByID(productID)
	if err != nil {
		http.Error(w, "Product Not Deleted", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
