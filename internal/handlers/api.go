// internal/handlers/api.go
package handlers

import (
	"context"
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

		// Add email to context for downstream handlers
		ctx := context.WithValue(r.Context(), "userEmail", email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *APIHandler) getCategoryAveragePrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["id"]

	avgPrice, err := h.categoryService.GetAveragePrice(categoryID)
	if err != nil {
		h.logger.Error().Err(err).Str("categoryID", categoryID).Msg("Failed to get average price")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]float64{"average_price": avgPrice}
	json.NewEncoder(w).Encode(response)
}

// Other handler methods would follow similar patterns...
