package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/cybre/re-partners-assignment/internal/api"
	"github.com/cybre/re-partners-assignment/internal/api/testdata"
	"github.com/cybre/re-partners-assignment/internal/models"
	"github.com/labstack/echo/v4"
)

func TestGetPackSizesHandler_Success(t *testing.T) {
	expectedPackSizes := []models.PackSize{
		{MaxItems: 250},
		{MaxItems: 500},
		{MaxItems: 1000},
		{MaxItems: 2000},
		{MaxItems: 5000},
	}
	mockPackingService := &testdata.MockPackingService{
		PackSizes: expectedPackSizes,
	}

	handler := api.GetPackSizesHandler(mockPackingService)

	// Create a new Echo context for testing
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/pack-sizes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err := handler(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Check the response body
	var packSizes []models.PackSize
	_ = json.Unmarshal(rec.Body.Bytes(), &packSizes)
	if !reflect.DeepEqual(packSizes, expectedPackSizes) {
		t.Errorf("Expected pack sizes %+v, got %+v", expectedPackSizes, packSizes)
	}
}

func TestGetPackSizesHandler_ServiceError(t *testing.T) {
	mockPackingService := &testdata.MockPackingService{
		Error: errors.New("service error"),
	}

	handler := api.GetPackSizesHandler(mockPackingService)

	// Create a new Echo context for testing
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/pack-sizes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err := handler(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestUpdatePackSizesHandler_Success(t *testing.T) {
	mockPackingService := &testdata.MockPackingService{}

	handler := api.UpdatePackSizesHadler(mockPackingService)

	// Create a new Echo context for testing
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/pack-sizes", bytes.NewReader([]byte(`[{"maxItems": 250}]`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err := handler(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestUpdatePackSizesHandler_ServiceError(t *testing.T) {
	mockPackingService := &testdata.MockPackingService{
		Error: errors.New("service error"),
	}

	handler := api.UpdatePackSizesHadler(mockPackingService)

	// Create a new Echo context for testing
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/pack-sizes", bytes.NewReader([]byte(`[{"maxItems": 250}]`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err := handler(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestUpdatePackSizesHandler_BadInputError(t *testing.T) {
	mockPackingService := &testdata.MockPackingService{}

	handler := api.UpdatePackSizesHadler(mockPackingService)

	// Create a new Echo context for testing
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/pack-sizes", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err := handler(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestPackOrderHandler_Success(t *testing.T) {
	expectedResponse := map[int]int{
		500: 1,
	}
	mockPackingService := &testdata.MockPackingService{
		Packs: expectedResponse,
	}

	handler := api.PackOrderHandler(mockPackingService)

	// Create a new Echo context for testing
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/pack-order", bytes.NewReader([]byte(`{"itemQty": 251}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err := handler(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	var orderPacks map[int]int
	_ = json.Unmarshal(rec.Body.Bytes(), &orderPacks)
	if !reflect.DeepEqual(orderPacks, expectedResponse) {
		t.Errorf("Expected order packs %+v, got %+v", expectedResponse, orderPacks)
	}
}

func TestPackOrderHandler_BadInputError(t *testing.T) {
	mockPackingService := &testdata.MockPackingService{}

	handler := api.PackOrderHandler(mockPackingService)

	// Create a new Echo context for testing
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/pack-order", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err := handler(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestPackOrderHandler_ServiceError(t *testing.T) {
	mockPackingService := &testdata.MockPackingService{
		Error: errors.New("service error"),
	}

	handler := api.PackOrderHandler(mockPackingService)

	// Create a new Echo context for testing
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/pack-order", bytes.NewReader([]byte(`{"itemQty": 251}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err := handler(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
