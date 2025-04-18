package handler

import (
	"net/http"
	"strings"

	"swift-codes-project/service"

	"github.com/gorilla/mux"
)

// this is returned if the requested code is the branch
type brachResponsePayload struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type headOfficeResponsePayload struct {
	Address       string                  `json:"address"`
	BankName      string                  `json:"bankName"`
	CountryISO2   string                  `json:"countryISO2"`
	CountryName   string                  `json:"countryName"`
	IsHeadquarter bool                    `json:"isHeadquarter"`
	SwiftCode     string                  `json:"swiftCode"`
	Branches      []branchResponsePayload `json:"branches"`
}

// this is returned with “list by country” endpoint
type countryResponsePayload struct {
	CountryISO2 string                  `json:"countryISO2"`
	CountryName string                  `json:"countryName"`
	SwiftCodes  []branchResponsePayload `json:"swiftCodes"`
}

type SwiftHTTPHandler struct {
	DataStore *service.SwiftRepository
}

// GET /v1/swift-codes/{code}
func (httpHandler *SwiftHTTPHandler) GetSwiftCode(
	responseWriter http.ResponseWriter,
	incomingRequest *http.Request,
) {
	pathVariables := mux.Vars(incomingRequest)
	requestedSwiftCode := strings.ToUpper(pathVariables["code"])

	headOfficeRow, branchRows, queryError := httpHandler.DataStore.GetSwiftCode(requestedSwiftCode)

	if queryError != nil {
		http.Error(responseWriter, `{"error": "not found"}`, http.StatusNotFound)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")

	//case 1, the requested row itself is a branch
}
