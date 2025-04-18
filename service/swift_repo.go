package service

import (
	"database/sql"
	"fmt"
	"swift-codes-project/models"
)

type SwiftRepo struct{ DB *sql.DB }

// returns: (requested code row, []branchRows, error)
func (r *SwiftRepo) GetByCode(code string) (models.SwiftCode, []models.SwiftCode, error) {
	var sc models.SwiftCode
	q := `SELECT country_iso2, swift_code, code_type, name, address,
	             town_name, country_name, time_zone,
	             is_headquarter, hq_swift_code
	      FROM swift_codes
	      WHERE swift_code = ?`
	if err := r.DB.QueryRow(q, code).Scan(
		&sc.CountryISO2, &sc.SwiftCode, &sc.CodeType, &sc.Name, &sc.Address,
		&sc.TownName, &sc.CountryName, &sc.TimeZone,
		&sc.IsHeadquarter, &sc.HqSwiftCode,
	); err != nil {
		if err == sql.ErrNoRows {
			return sc, nil, fmt.Errorf("not found")
		}
		return sc, nil, err
	}
	// if it’s a head‑office fetch branches
	var branches []models.SwiftCode
	if sc.IsHeadquarter {
		rows, err := r.DB.Query(
			`SELECT country_iso2, swift_code, code_type, name, address,
						town_name, country_name, time_zone,
						is_headquarter, hq_swift_code
				 FROM swift_codes
				 WHERE hq_swift_code = ?`, sc.SwiftCode)
		if err != nil {
			return sc, nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var br models.SwiftCode
			if err := rows.Scan(
				&br.CountryISO2, &br.SwiftCode, &br.CodeType, &br.Name, &br.Address,
				&br.TownName, &br.CountryName, &br.TimeZone,
				&br.IsHeadquarter, &br.HqSwiftCode,
			); err != nil {
				return sc, nil, err
			}
			branches = append(branches, br)
		}
	}
	return sc, branches, nil
}
