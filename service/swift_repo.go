package service

import (
	"database/sql"
	"errors"
	"swift-codes-project/models"
)

type SwiftRepository struct {
	DB *sql.DB
}

// GetSwiftCode returns:
// the exact row whose swift_code = requestedCode
// if that row is a head‑office, all its branch rows

func (repo *SwiftRepository) GetSwiftCode(requestedCode string) (models.SwiftCode, []models.SwiftCode, error) {
	const findByCodeSQL = `
		SELECT country_iso2, swift_code, code_type, name, address,
		       town_name, country_name, time_zone,
		       is_headquarter, hq_swift_code
		  FROM swift_codes
		 WHERE swift_code = ?;
	`

	var headOffice models.SwiftCode
	err := repo.DB.QueryRow(findByCodeSQL, requestedCode).Scan(
		&headOffice.CountryISO2, &headOffice.SwiftCode, &headOffice.CodeType,
		&headOffice.Name, &headOffice.Address, &headOffice.TownName,
		&headOffice.CountryName, &headOffice.TimeZone,
		&headOffice.IsHeadquarter, &headOffice.HqSwiftCode,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return headOffice, nil, errors.New("not found")
		}
		return headOffice, nil, err
	}

	// If the row is a head‑office, pull its branches.
	if !headOffice.IsHeadquarter {
		return headOffice, nil, nil // it’s just a branch – done.
	}

	const findBranchesSQL = `
		SELECT country_iso2, swift_code, code_type, name, address,
		       town_name, country_name, time_zone,
		       is_headquarter, hq_swift_code
		  FROM swift_codes
		 WHERE hq_swift_code = ?;
	`
	rows, err := repo.DB.Query(findBranchesSQL, headOffice.SwiftCode)
	if err != nil {
		return headOffice, nil, err
	}
	defer rows.Close()

	var branches []models.SwiftCode
	for rows.Next() {
		var branch models.SwiftCode
		if err := rows.Scan(
			&branch.CountryISO2, &branch.SwiftCode, &branch.CodeType,
			&branch.Name, &branch.Address, &branch.TownName,
			&branch.CountryName, &branch.TimeZone,
			&branch.IsHeadquarter, &branch.HqSwiftCode,
		); err != nil {
			return headOffice, nil, err
		}
		branches = append(branches, branch)
	}
	return headOffice, branches, nil
}

// GetCountrySwiftCodes returns all rows for the given ISO‑2 country code.
func (repo *SwiftRepository) GetCountrySwiftCodes(iso2 string) ([]models.SwiftCode, error) {
	const byCountrySQL = `
		SELECT country_iso2, swift_code, code_type, name, address,
		       town_name, country_name, time_zone,
		       is_headquarter, hq_swift_code
		  FROM swift_codes
		 WHERE country_iso2 = ?;
	`

	rows, err := repo.DB.Query(byCountrySQL, iso2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SwiftCode
	for rows.Next() {
		var sc models.SwiftCode
		if err := rows.Scan(
			&sc.CountryISO2, &sc.SwiftCode, &sc.CodeType,
			&sc.Name, &sc.Address, &sc.TownName,
			&sc.CountryName, &sc.TimeZone,
			&sc.IsHeadquarter, &sc.HqSwiftCode,
		); err != nil {
			return nil, err
		}
		results = append(results, sc)
	}
	return results, nil
}

// CreateSwiftCode inserts a brand‑new row. It returns an error if the PK clashes
// (duplicate swift_code) or if the SQL fails.
func (repo *SwiftRepository) CreateSwiftCode(sc models.SwiftCode) error {
	const insertSQL = `
		INSERT INTO swift_codes (
			country_iso2, swift_code, code_type, name, address,
			town_name, country_name, time_zone,
			is_headquarter, hq_swift_code
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err := repo.DB.Exec(insertSQL,
		sc.CountryISO2, sc.SwiftCode, sc.CodeType, sc.Name, sc.Address,
		sc.TownName, sc.CountryName, sc.TimeZone,
		sc.IsHeadquarter, sc.HqSwiftCode,
	)
	return err
}

// DeleteSwiftCode removes the row whose swift_code = codeToDelete.
func (repo *SwiftRepository) DeleteSwiftCode(codeToDelete string) error {
	const deleteSQL = `DELETE FROM swift_codes WHERE swift_code = ?;`
	_, err := repo.DB.Exec(deleteSQL, codeToDelete)
	return err
}
