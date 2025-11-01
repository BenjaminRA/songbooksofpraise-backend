package models

import (
	"database/sql"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
)

type Country struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	ISOAlpha2  string  `json:"iso_alpha2"`
	ISOAlpha3  string  `json:"iso_alpha3"`
	ISONumeric string  `json:"iso_numeric"`
	States     []State `json:"states,omitempty"` // Not in database, but used in API responses
}

func (c *Country) GetAllCountries() ([]Country, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, name, iso_alpha2, iso_alpha3, iso_numeric FROM countries ORDER BY name")
	if err != nil {
		return []Country{}, err
	}
	defer rows.Close()

	result := []Country{}
	for rows.Next() {
		elem := Country{}
		err := rows.Scan(&elem.ID, &elem.Name, &elem.ISOAlpha2, &elem.ISOAlpha3, &elem.ISONumeric)
		if err != nil {
			continue
		}
		result = append(result, elem)
	}

	return result, nil
}

func (c *Country) GetCountryByID(id int) (*Country, error) {
	db := sqlite.GetDBConnection()
	row := db.QueryRow("SELECT id, name, iso_alpha2, iso_alpha3, iso_numeric FROM countries WHERE id = ?", id)

	elem := Country{}
	err := row.Scan(&elem.ID, &elem.Name, &elem.ISOAlpha2, &elem.ISOAlpha3, &elem.ISONumeric)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Get states for this country
	states, err := (&State{}).GetStatesByCountryID(id)
	if err == nil {
		elem.States = states
	} else {
		elem.States = []State{}
	}

	return &elem, nil
}

func (c *Country) GetCountriesWithStates() ([]Country, error) {
	countries, err := c.GetAllCountries()
	if err != nil {
		return []Country{}, err
	}

	// Get states for each country
	for i := range countries {
		states, err := (&State{}).GetStatesByCountryID(countries[i].ID)
		if err == nil {
			countries[i].States = states
		} else {
			countries[i].States = []State{}
		}
	}

	return countries, nil
}
