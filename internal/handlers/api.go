// internal/handlers/api.go
package handlers

import (
	"context"
	"ecommerce-service/internal/models"
	"ecommerce-service/internal/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type APIHandler struct {
	authService     services.AuthService
	orderService    services.OrderService
	productService  services.ProductService
	categoryService services.CategoryService
	logger          *zerolog.Logger
}

func NewAPIHandler(
	authService services.AuthService,
	orderService services.OrderService,
	productService services.ProductService,
	categoryService services.CategoryService,
	logger *zerolog.Logger,
) *APIHandler {
	return &APIHandler{
		authService:     authService,
		orderService:    orderService,
		productService:  productService,
		categoryService: categoryService,
		logger:          logger,
	}
}

func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	// Product routes
	router.HandleFunc("/products", h.createProduct).Methods("POST")
	router.HandleFunc("/products/{id}", h.getProduct).Methods("GET")
	router.HandleFunc("/categories/{id}/average-price", h.getCategoryAveragePrice).Methods("GET")

	// Order routes
	router.HandleFunc("/orders", h.createOrder).Methods("POST")
	router.HandleFunc("/orders/{id}", h.getOrder).Methods("GET")

	// Auth middleware
	authRouter := router.PathPrefix("/").Subrouter()
	authRouter.Use(h.authMiddleware)
	authRouter.HandleFunc("/secure-route", h.secureHandler).Methods("GET")

	router.ServeHTTP(w, r)
}

func (h *APIHandler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		email, err := h.authService.Authenticate(token)
		if err != nil {
			h.logger.Error().Err(err).Msg("Authentication failed")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userEmail", email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *APIHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode product")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.productService.CreateProduct(r.Context(), &product); err != nil {
		h.logger.Error().Err(err).Msg("Failed to create product")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *APIHandler) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	product, err := h.productService.GetProduct(r.Context(), productID)
	if err != nil {
		h.logger.Error().Err(err).Str("productID", productID).Msg("Failed to get product")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *APIHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode order")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.orderService.CreateOrder(r.Context(), &order); err != nil {
		h.logger.Error().Err(err).Msg("Failed to create order")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *APIHandler) getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	order, err := h.orderService.GetOrder(r.Context(), orderID)
	if err != nil {
		h.logger.Error().Err(err).Str("orderID", orderID).Msg("Failed to get order")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *APIHandler) getCategoryAveragePrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["id"]

	avgPrice, err := h.categoryService.GetAveragePrice(r.Context(), categoryID)
	if err != nil {
		h.logger.Error().Err(err).Str("categoryID", categoryID).Msg("Failed to get average price")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]float64{"average_price": avgPrice}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *APIHandler) secureHandler(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("userEmail").(string)
	response := map[string]string{
		"message": "Authenticated successfully",
		"email":   email,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
