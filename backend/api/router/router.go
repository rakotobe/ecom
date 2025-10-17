package router

import (
	"ecom-backend/api/handler"
	"ecom-backend/api/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// Setup creates and configures the HTTP router
func Setup(
	productHandler *handler.ProductHandler,
	basketHandler *handler.BasketHandler,
	orderHandler *handler.OrderHandler,
) *mux.Router {
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.CORS)
	r.Use(middleware.Logging)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Product routes
	api.HandleFunc("/products", productHandler.CreateProduct).Methods("POST", "OPTIONS")
	api.HandleFunc("/products", productHandler.GetAllProducts).Methods("GET", "OPTIONS")
	api.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET", "OPTIONS")
	api.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT", "OPTIONS")
	api.HandleFunc("/products/{id}/stock", productHandler.UpdateStock).Methods("PATCH", "OPTIONS")
	api.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE", "OPTIONS")

	// Basket routes
	api.HandleFunc("/baskets", basketHandler.CreateBasket).Methods("POST", "OPTIONS")
	api.HandleFunc("/baskets/{id}", basketHandler.GetBasket).Methods("GET", "OPTIONS")
	api.HandleFunc("/baskets/{id}/items", basketHandler.AddItem).Methods("POST", "OPTIONS")
	api.HandleFunc("/baskets/{id}/items/{productId}", basketHandler.RemoveItem).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/baskets/{id}/items/{productId}", basketHandler.UpdateItemQuantity).Methods("PATCH", "OPTIONS")
	api.HandleFunc("/baskets/{id}/items", basketHandler.ClearBasket).Methods("DELETE", "OPTIONS")

	// Order routes
	api.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST", "OPTIONS")
	api.HandleFunc("/orders", orderHandler.GetAllOrders).Methods("GET", "OPTIONS")
	api.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET", "OPTIONS")
	api.HandleFunc("/orders/{id}/confirm", orderHandler.ConfirmOrder).Methods("POST", "OPTIONS")
	api.HandleFunc("/orders/{id}/ship", orderHandler.ShipOrder).Methods("POST", "OPTIONS")
	api.HandleFunc("/orders/{id}/deliver", orderHandler.DeliverOrder).Methods("POST", "OPTIONS")
	api.HandleFunc("/orders/{id}/cancel", orderHandler.CancelOrder).Methods("POST", "OPTIONS")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return r
}
