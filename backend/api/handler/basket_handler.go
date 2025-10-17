package handler

import (
	"encoding/json"
	"ecom-backend/application/dto"
	"ecom-backend/application/service"
	"net/http"

	"github.com/gorilla/mux"
)

// BasketHandler handles basket HTTP requests
type BasketHandler struct {
	basketService *service.BasketService
}

// NewBasketHandler creates a new BasketHandler
func NewBasketHandler(basketService *service.BasketService) *BasketHandler {
	return &BasketHandler{
		basketService: basketService,
	}
}

// CreateBasket handles POST /baskets
func (h *BasketHandler) CreateBasket(w http.ResponseWriter, r *http.Request) {
	basket, err := h.basketService.CreateBasket(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, basket)
}

// GetBasket handles GET /baskets/{id}
func (h *BasketHandler) GetBasket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	basket, err := h.basketService.GetBasket(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, basket)
}

// AddItem handles POST /baskets/{id}/items
func (h *BasketHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	basketID := vars["id"]

	var req dto.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	basket, err := h.basketService.AddItem(r.Context(), basketID, &req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, basket)
}

// RemoveItem handles DELETE /baskets/{id}/items/{productId}
func (h *BasketHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	basketID := vars["id"]
	productID := vars["productId"]

	basket, err := h.basketService.RemoveItem(r.Context(), basketID, productID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, basket)
}

// UpdateItemQuantity handles PATCH /baskets/{id}/items/{productId}
func (h *BasketHandler) UpdateItemQuantity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	basketID := vars["id"]
	productID := vars["productId"]

	var req dto.UpdateItemQuantityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	basket, err := h.basketService.UpdateItemQuantity(r.Context(), basketID, productID, &req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, basket)
}

// ClearBasket handles DELETE /baskets/{id}/items
func (h *BasketHandler) ClearBasket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	basketID := vars["id"]

	basket, err := h.basketService.ClearBasket(r.Context(), basketID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, basket)
}
