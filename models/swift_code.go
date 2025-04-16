package models

// Swift Code Model
type SwiftCode struct {
	CountryISO2 string `json:"countryISO2"`
	SwiftCode   string `json:"swiftCode"`
	CodeType    string `json:"codeType"`
	Name        string `json:"Name"`
	Address     string `json:"address"`
	TownName    string `json:"townName"`
	CountryName string `json:"countryName"`
	TimeZone    string `json:"timeZone"`
}
