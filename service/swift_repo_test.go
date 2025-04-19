package service

import (
	"testing"

	"swift-codes-project/db"
	"swift-codes-project/models"

	_ "github.com/mattn/go-sqlite3"
)

func TestGetSwiftCodeReturnsHeadOfficeAndBranches(t *testing.T) {
	testDatabse, initError := db.InitDB("file::memory:?cache=shared&_fk=1")

	if initError != nil {
		t.Fatalf("Failed to initialize in memory databse: %v", initError)
	}
	defer testDatabse.Close()

	repository := &SwiftRepository{DB: testDatabse}

	headOfficeEntry := models.SwiftCode{
		CountryISO2:   "ZZ",
		SwiftCode:     "ZZBANKXXX",
		CodeType:      "BANK",
		Name:          "ZELAND NATIONAL BANK",
		Address:       "1 MAIN PLAZA",
		TownName:      "CAPITAL",
		CountryName:   "ZELAND",
		TimeZone:      "UTC+00:00",
		IsHeadquarter: true,
		HqSwiftCode:   "",
	}
	firstBranchEntry := models.SwiftCode{
		CountryISO2:   "ZZ",
		SwiftCode:     "ZZBANK001",
		CodeType:      "BRANCH",
		Name:          "ZELAND NATIONAL BANK",
		Address:       "2 SECOND AVE",
		TownName:      "CAPITAL",
		CountryName:   "ZELAND",
		TimeZone:      "UTC+00:00",
		IsHeadquarter: false,
		HqSwiftCode:   "ZZBANKXXX",
	}
	secondBranchEntry := models.SwiftCode{
		CountryISO2:   "ZZ",
		SwiftCode:     "ZZBANK002",
		CodeType:      "BRANCH",
		Name:          "ZELAND NATIONAL BANK",
		Address:       "3 THIRD STREET",
		TownName:      "CAPITAL",
		CountryName:   "ZELAND",
		TimeZone:      "UTC+00:00",
		IsHeadquarter: false,
		HqSwiftCode:   "ZZBANKXXX",
	}
	if insertError := repository.CreateSwiftCode(headOfficeEntry); insertError != nil {
		t.Fatalf("Unexpected error inserting head office: %v", insertError)
	}
	if insertError := repository.CreateSwiftCode(firstBranchEntry); insertError != nil {
		t.Fatalf("Unexpected error inserting first branch: %v", insertError)
	}
	if insertError := repository.CreateSwiftCode(secondBranchEntry); insertError != nil {
		t.Fatalf("Unexpected error inserting second branch: %v", insertError)
	}

	returnedHeadOffice, returnedBranches, queryError := repository.GetSwiftCode("ZZBANKXXX")
	if queryError != nil {
		t.Fatalf("Expected no error querying head office, got: %v", queryError)
	}
	if !returnedHeadOffice.IsHeadquarter {
		t.Errorf("Expected IsHeadquarter=True for head office, got false")
	}
	if len(returnedBranches) != 2 {
		t.Errorf("Expected branches 2 instead got %d", len(returnedBranches))
	}
}

// TestGetCountrySwiftCodesFiltersCorrectly seeds two countries and asserts
// that only the requested ISO‑2 code is returned.
func TestGetCountrySwiftCodesFiltersCorrectly(t *testing.T) {
	testDatabase, initError := db.InitDB("file::memory:?cache=shared&_fk=1")
	if initError != nil {
		t.Fatalf("Failed to initialize in-memory database: %v", initError)
	}
	defer testDatabase.Close()

	repository := &SwiftRepository{DB: testDatabase}

	entryForCountryAA := models.SwiftCode{
		CountryISO2:   "AA",
		SwiftCode:     "AABANKXXX",
		CodeType:      "BANK",
		Name:          "BANK AA",
		Address:       "1 AA STREET",
		TownName:      "AA TOWN",
		CountryName:   "AALAND",
		TimeZone:      "UTC+01:00",
		IsHeadquarter: true,
		HqSwiftCode:   "",
	}
	entryForCountryBB := models.SwiftCode{
		CountryISO2:   "BB",
		SwiftCode:     "BBBANKXXX",
		CodeType:      "BANK",
		Name:          "BANK BB",
		Address:       "1 BB STREET",
		TownName:      "BB TOWN",
		CountryName:   "BBLAND",
		TimeZone:      "UTC+02:00",
		IsHeadquarter: true,
		HqSwiftCode:   "",
	}

	repository.CreateSwiftCode(entryForCountryAA)
	repository.CreateSwiftCode(entryForCountryBB)

	codesForCountryAA, queryError := repository.GetCountrySwiftCodes("AA")
	if queryError != nil {
		t.Fatalf("Expected no error querying country AA, got: %v", queryError)
	}
	if len(codesForCountryAA) != 1 {
		t.Errorf("Expected 1 code for country AA, got %d", len(codesForCountryAA))
	}
	if codesForCountryAA[0].CountryISO2 != "AA" {
		t.Errorf("Expected CountryISO2 'AA', got '%s'", codesForCountryAA[0].CountryISO2)
	}
}
