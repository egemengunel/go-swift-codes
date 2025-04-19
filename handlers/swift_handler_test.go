package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"swift-codes-project/models"
)

type stubSwiftRepository struct{}

func (stub *stubSwiftRepository) GetSwiftCode(requestedCode string) (models.SwiftCode, []models.SwiftCode, error) {
	headOfficeData := models.SwiftCode{
		Address:       "HQ Address",
		Name:          "HQ Bank",
		CountryISO2:   "ZZ",
		CountryName:   "ZELAND",
		IsHeadquarter: true,
		SwiftCode:     requestedCode,
	}
	branchData := models.SwiftCode{
		Address:       "Branch Address",
		Name:          "Branch Bank",
		CountryISO2:   "ZZ",
		CountryName:   "ZELAND",
		IsHeadquarter: false,
		SwiftCode:     requestedCode + "001",
	}
	return headOfficeData, []models.SwiftCode{branchData}, nil
}

// GetCountrySwiftCodes returns a single entry for the requested ISO‑2 code.
func (stub *stubSwiftRepository) GetCountrySwiftCodes(requestedISO2 string) ([]models.SwiftCode, error) {
	entry := models.SwiftCode{
		Address:       "Some Address",
		Name:          "Some Bank",
		CountryISO2:   requestedISO2,
		CountryName:   "COUNTRY NAME",
		IsHeadquarter: true,
		SwiftCode:     requestedISO2 + "BANKXXX",
	}
	return []models.SwiftCode{entry}, nil
}

// CreateSwiftCode always succeeds.
func (stub *stubSwiftRepository) CreateSwiftCode(newCode models.SwiftCode) error {
	return nil
}

// DeleteSwiftCode always succeeds.
func (stub *stubSwiftRepository) DeleteSwiftCode(codeToDelete string) error {
	return nil
}

// TestGetSwiftCodeHandler_Success tests the GET /v1/swift-codes/{code} handler.
func TestGetSwiftCodeHandler_Success(t *testing.T) {
	testRequest := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/ZZBANKXXX", nil)
	testRequest = mux.SetURLVars(testRequest, map[string]string{"code": "ZZBANKXXX"})
	responseRecorder := httptest.NewRecorder()

	handlerInstance := &SwiftHTTPHandler{DataStore: &stubSwiftRepository{}}
	handlerInstance.GetSwiftCode(responseRecorder, testRequest)

	response := responseRecorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", response.StatusCode)
	}

	var decodedPayload headOfficeResponsePayload
	if decodeError := json.NewDecoder(response.Body).Decode(&decodedPayload); decodeError != nil {
		t.Fatalf("Failed to decode JSON response: %v", decodeError)
	}

	if decodedPayload.SwiftCode != "ZZBANKXXX" {
		t.Errorf("Expected SwiftCode 'ZZBANKXXX', got '%s'", decodedPayload.SwiftCode)
	}
	if len(decodedPayload.Branches) != 1 {
		t.Errorf("Expected 1 branch, got %d", len(decodedPayload.Branches))
	}
}

// TestGetCountrySwiftCodesHandler_Success tests the GET /v1/swift-codes/country/{iso2} handler.
func TestGetCountrySwiftCodesHandler_Success(t *testing.T) {
	testRequest := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/country/ZZ", nil)
	testRequest = mux.SetURLVars(testRequest, map[string]string{"iso2": "ZZ"})
	responseRecorder := httptest.NewRecorder()

	handlerInstance := &SwiftHTTPHandler{DataStore: &stubSwiftRepository{}}
	handlerInstance.GetCountrySwiftCodes(responseRecorder, testRequest)

	response := responseRecorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", response.StatusCode)
	}

	var decodedPayload countryResponsePayload
	if decodeError := json.NewDecoder(response.Body).Decode(&decodedPayload); decodeError != nil {
		t.Fatalf("Failed to decode JSON response: %v", decodeError)
	}

	if decodedPayload.CountryISO2 != "ZZ" {
		t.Errorf("Expected CountryISO2 'ZZ', got '%s'", decodedPayload.CountryISO2)
	}
	if len(decodedPayload.SwiftCodes) != 1 {
		t.Errorf("Expected 1 entry in SwiftCodes, got %d", len(decodedPayload.SwiftCodes))
	}
}

// TestCreateSwiftCodeHandler_Success tests the POST /v1/swift-codes handler.
func TestCreateSwiftCodeHandler_Success(t *testing.T) {
	exampleCode := models.SwiftCode{
		Address:       "123 Test Ave",
		Name:          "Test Bank",
		CountryISO2:   "ZZ",
		CountryName:   "ZELAND",
		IsHeadquarter: false,
		SwiftCode:     "ZZTEST001",
	}
	requestBodyBytes, _ := json.Marshal(exampleCode)
	testRequest := httptest.NewRequest(
		http.MethodPost,
		"/v1/swift-codes",
		bytes.NewBuffer(requestBodyBytes),
	)
	testRequest.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()

	handlerInstance := &SwiftHTTPHandler{DataStore: &stubSwiftRepository{}}
	handlerInstance.CreateSwiftCode(responseRecorder, testRequest)

	response := responseRecorder.Result()
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d", response.StatusCode)
	}
}

// TestDeleteSwiftCodeHandler_Success tests the DELETE /v1/swift-codes/{code} handler.
func TestDeleteSwiftCodeHandler_Success(t *testing.T) {
	testRequest := httptest.NewRequest(http.MethodDelete, "/v1/swift-codes/ZZTEST001", nil)
	testRequest = mux.SetURLVars(testRequest, map[string]string{"code": "ZZTEST001"})
	responseRecorder := httptest.NewRecorder()

	handlerInstance := &SwiftHTTPHandler{DataStore: &stubSwiftRepository{}}
	handlerInstance.DeleteSwiftCode(responseRecorder, testRequest)

	response := responseRecorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", response.StatusCode)
	}
}
