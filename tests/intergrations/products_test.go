package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ecommerce-service/internal/handlers"
	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"
	"ecommerce-service/internal/services"
	"ecommerce-service/tests"

	"github.com/stretchr/testify/assert"
)

func TestProductAPI(t *testing.T) {
	db := tests.SetupTestDB(t)

	// Initialize services
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productService := services.NewProductService(productRepo, categoryRepo)

	// Setup handler
	h := handlers.NewAPIHandler(
		nil, // auth service
		nil, // order service
		productService,
		nil, // category service
		nil, // logger
	)

	t.Run("Create Product", func(t *testing.T) {
		product := models.Product{
			Name:       "Integration Test Product",
			Price:      29.99,
			CategoryID: "test-category",
		}

		body, _ := json.Marshal(product)
		req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Product
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response.ID)
	})
}
