package models

import (
	"database/sql"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
)

type State struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	CountryID int      `json:"country_id"`
	Country   *Country `json:"country,omitempty"`  // Not in database, but used in API responses
	Churches  []Church `json:"churches,omitempty"` // Not in database, but used in API responses
}

func (s *State) GetAllStates() ([]State, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, name, country_id FROM states ORDER BY name")
	if err != nil {
		return []State{}, err
	}
	defer rows.Close()

	result := []State{}
	for rows.Next() {
		elem := State{}
		err := rows.Scan(&elem.ID, &elem.Name, &elem.CountryID)
		if err != nil {
			continue
		}
		result = append(result, elem)
	}

	return result, nil
}

func (s *State) GetStateByID(id int) (*State, error) {
	db := sqlite.GetDBConnection()
	row := db.QueryRow("SELECT id, name, country_id FROM states WHERE id = ?", id)

	elem := State{}
	err := row.Scan(&elem.ID, &elem.Name, &elem.CountryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Get country for this state
	country, err := (&Country{}).GetCountryByID(elem.CountryID)
	if err == nil && country != nil {
		elem.Country = country
	}

	return &elem, nil
}

func (s *State) GetStatesByCountryID(countryID int) ([]State, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, name, country_id FROM states WHERE country_id = ? ORDER BY name", countryID)
	if err != nil {
		return []State{}, err
	}
	defer rows.Close()

	result := []State{}
	for rows.Next() {
		elem := State{}
		err := rows.Scan(&elem.ID, &elem.Name, &elem.CountryID)
		if err != nil {
			continue
		}
		result = append(result, elem)
	}

	return result, nil
}
