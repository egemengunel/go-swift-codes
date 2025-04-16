package parser

import (
	"database/sql"
	"fmt"

	"swift-codes-project/models"

	"github.com/xuri/excelize/v2"
)

//Function to open excel and parse each row, convert to SwiftCode objects and store each entry in db

func ParseExcelAndStore(db *sql.DB, filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to open excel file %v", err)
	}
	defer f.Close()

	//since data is on the first sheet
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)

	if err != nil {
		return fmt.Errorf("unable to get rows %v", err)
	}
	if len(rows) < 2 {
		return fmt.Errorf("not enough rows")
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) < 8 {
			continue
		}

		// Map the columns to the SwiftCode model fields.
		codeEntry := models.SwiftCode{
			CountryISO2: row[0],
			SwiftCode:   row[1],
			CodeType:    row[2],
			Name:        row[3],
			Address:     row[4],
			TownName:    row[5],
			CountryName: row[6],
			TimeZone:    row[7],
		}
		// Insert the SwiftCode entry into the database.
		if err := InsertSwiftCode(db, codeEntry); err != nil {
			return fmt.Errorf("failed to insert data at row %d: %v", i, err)
		}
	}
	return nil
}

// InsertSwiftCode is a placeholder function
func InsertSwiftCode(db *sql.DB, sc models.SwiftCode) error {
	query := `
        INSERT INTO swift_codes (country_iso2, swift_code, code_type, name, address, town_name, country_name, time_zone)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(query,
		sc.CountryISO2,
		sc.SwiftCode,
		sc.CodeType,
		sc.Name,
		sc.Address,
		sc.TownName,
		sc.CountryName,
		sc.TimeZone,
	)
	return err
}
