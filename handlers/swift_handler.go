package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"swift-codes-project/models"
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
	if !headOfficeRow.IsHeadquarter {
		branchPayload := brachResponsePayload{
			Address:       headOfficeRow.Address,
			BankName:      headOfficeRow.Name,
			CountryISO2:   headOfficeRow.CountryISO2,
			CountryName:   headOfficeRow.CountryName,
			IsHeadquarter: false,
			SwiftCode:     headOfficeRow.SwiftCode,
		}
		json.NewEncoder(responseWriter).Encode(branchPayload)
		return
	}
	//case 2 the requested row is a head office
	headOfficePayload := headOfficeResponsePayload{
		Address:       headOfficeRow.Address,
		BankName:      headOfficeRow.Name,
		CountryISO2:   headOfficeRow.CountryISO2,
		CountryName:   headOfficeRow.CountryName,
		IsHeadquarter: true,
		SwiftCode:     headOfficeRow.SwiftCode,
	}
	for _, branchRow := range branchRows {
		headOfficePayload.Branches = append(headOfficePayload.Branches, branchResponsePayload{
			Address:       branchRow.Address,
			BankName:      branchRow.Name,
			CountryISO2:   branchRow.CountryISO2,
			CountryName:   branchRow.CountryName,
			IsHeadquarter: false,
			SwiftCode:     branchRow.SwiftCode,
		})
	}
	json.NewEncoder(responseWriter).Encode(headOfficePayload)
}

// GET /v1/swift-codes/country/{iso2}

func (httpHandler *SwiftHTTPHandler) GetCountrySwiftCodes(responseWriter http.ResponseWriter, incomingRequest http.Request) {
	pathVariables := mux.Vars(&incomingRequest)
	requestedISO2 := strings.ToUpper(pathVariables["iso2"])
	allRows, queryError :=
		httpHandler.DataStore.GetCountrySwiftCodes(requestedISO2)
	if queryError != nil || len(allRows) == 0 {
		http.Error(responseWriter, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	var listPayload []brachResponsePayload
	for _, row := range allRows {
		listPayload = append(listPayload, branchResponsePayload{
			Address:       row.Address,
			BankName:      row.Name,
			CountryISO2:   row.CountryISO2,
			CountryName:   row.CountryName,
			IsHeadquarter: row.IsHeadquarter,
			SwiftCode:     row.SwiftCode,
		})
	}
	countryPayload := countryResponsePayload{
		CountryISO2: requestedISO2,
		CountryName: allRows[0].CountryName,
		SwiftCodes:  listPayload,
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(countryPayload)
}

// POST /v1/swift-codes
func (httpHandler *SwiftHTTPHandler) CreateSwiftCode(
	responseWriter http.ResponseWriter,
	incomingRequest *http.Request,
) {
	var incomingBody models.SwiftCode
	if err := json.NewDecoder(incomingRequest.Body).Decode(&incomingBody); err != nil {
		http.Error(responseWriter, `{"error":"bad json"}`, http.StatusBadRequest)
		return
	}

	if err := httpHandler.DataStore.CreateSwiftCode(incomingBody); err != nil {
		http.Error(responseWriter, `{"error":"cannot insert"}`, http.StatusConflict)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusCreated)
	responseWriter.Write([]byte(`{"message":"swift code created"}`))
}

// DELETE /v1/swift-codes/{code}

func (httpHandler *SwiftHTTPHandler) DeleteSwiftCode(
	responseWriter http.ResponseWriter,
	incomingRequest *http.Request,
) {
	requestedSwiftCode := strings.ToUpper(mux.Vars(incomingRequest)["code"])

	if err := httpHandler.DataStore.DeleteSwiftCode(requestedSwiftCode); err != nil {
		http.Error(responseWriter, `{"error":"db failure"}`,
			http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write([]byte(`{"message":"swift code deleted"}`))
}
