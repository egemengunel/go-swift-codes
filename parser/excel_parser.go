package parser

import (
	"database/sql"
	"fmt"
	"strings"
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
			CountryISO2: strings.ToUpper(row[0]),
			SwiftCode:   row[1],
			CodeType:    row[2],
			Name:        row[3],
			Address:     row[4],
			TownName:    row[5],
			CountryName: strings.ToUpper(row[6]),
			TimeZone:    row[7],
		}
		isHQ := strings.HasSuffix(codeEntry.SwiftCode, "XXX")
		hqCode := ""
		if !isHQ {
			hqCode = codeEntry.SwiftCode[:8] + "XXX"
		}
		codeEntry.IsHeadquarter = isHQ
		codeEntry.HqSwiftCode = hqCode
		// Insert the SwiftCode entry into the database.
		if err := InsertSwiftCode(db, codeEntry); err != nil {
			return fmt.Errorf("failed to insert data at row %d: %v", i, err)
		}
	}
	return nil
}

// InsertSwiftCode is a placeholder function
func InsertSwiftCode(db *sql.DB, sc models.SwiftCode) error {
	const insertSQL = `
    INSERT INTO swift_codes (
        country_iso2, swift_code, code_type, name, address,
        town_name, country_name, time_zone,
        is_headquarter, hq_swift_code
    ) VALUES (
        ?, ?, ?, ?, ?,
        ?, ?, ?, ?, ?
    );`

	_, err := db.Exec(insertSQL,
		sc.CountryISO2,   // 1
		sc.SwiftCode,     // 2
		sc.CodeType,      // 3
		sc.Name,          // 4
		sc.Address,       // 5
		sc.TownName,      // 6
		sc.CountryName,   // 7
		sc.TimeZone,      // 8
		sc.IsHeadquarter, // 9  <- new
		sc.HqSwiftCode,   // 10 <- new
	)
	return err
}
