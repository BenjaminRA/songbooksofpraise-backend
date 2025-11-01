package models

import (
	"database/sql"
	"time"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
)

type Church struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Phone       *string   `json:"phone"`
	Email       string    `json:"email"`
	Description *string   `json:"description"`
	Website     *string   `json:"website"`
	Established *string   `json:"established"` // Using string for DATE field
	Facebook    *string   `json:"facebook"`
	Instagram   *string   `json:"instagram"`
	YouTube     *string   `json:"youtube"`
	Spotify     *string   `json:"spotify"`
	StateID     int       `json:"state_id"`
	State       *State    `json:"state,omitempty"`    // Not in database, but used in API responses
	Elders      []Elder   `json:"elders,omitempty"`   // Not in database, but used in API responses
	Services    []Service `json:"services,omitempty"` // Not in database, but used in API responses
	Events      []Event   `json:"events,omitempty"`   // Not in database, but used in API responses
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (n *Church) GetAllChurches() ([]Church, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query(`
		SELECT id, name, address, phone, email, description, website, established, 
		       facebook, instagram, youtube, spotify, state_id 
		FROM churches 
		ORDER BY name
	`)
	if err != nil {
		return []Church{}, err
	}
	defer rows.Close()

	result := []Church{}
	for rows.Next() {
		elem := Church{}
		err := rows.Scan(
			&elem.ID, &elem.Name, &elem.Address, &elem.Phone, &elem.Email,
			&elem.Description, &elem.Website, &elem.Established, &elem.Facebook,
			&elem.Instagram, &elem.YouTube, &elem.Spotify, &elem.StateID,
		)
		if err != nil {
			continue
		}

		// Get state information
		if state, err := (&State{}).GetStateByID(elem.StateID); err == nil && state != nil {
			elem.State = state
		}

		result = append(result, elem)
	}

	return result, nil
}

func (n *Church) GetChurchByID(id int) (*Church, error) {
	db := sqlite.GetDBConnection()
	row := db.QueryRow(`
		SELECT id, name, address, phone, email, description, website, established, 
		       facebook, instagram, youtube, spotify, state_id 
		FROM churches 
		WHERE id = ?
	`, id)

	elem := Church{}
	err := row.Scan(
		&elem.ID, &elem.Name, &elem.Address, &elem.Phone, &elem.Email,
		&elem.Description, &elem.Website, &elem.Established, &elem.Facebook,
		&elem.Instagram, &elem.YouTube, &elem.Spotify, &elem.StateID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Get state information
	if state, err := (&State{}).GetStateByID(elem.StateID); err == nil && state != nil {
		elem.State = state
	}

	// Get elders for this church
	if elders, err := (&Elder{}).GetEldersByChurchID(id); err == nil {
		elem.Elders = elders
	} else {
		elem.Elders = []Elder{}
	}

	// Get services for this church
	if services, err := (&Service{}).GetServicesByChurchID(id); err == nil {
		elem.Services = services
	} else {
		elem.Services = []Service{}
	}

	// Get events for this church
	if events, err := (&Event{}).GetEventsByChurchID(id); err == nil {
		elem.Events = events
	} else {
		elem.Events = []Event{}
	}

	return &elem, nil
}

func (n *Church) CreateChurch() error {
	db := sqlite.GetDBConnection()

	n.CreatedAt = time.Now().UTC()
	n.UpdatedAt = time.Now().UTC()

	_, err := db.Exec(`
		INSERT INTO churches (name, address, phone, email, description, website, established, 
		                     facebook, instagram, youtube, spotify, state_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, n.Name, n.Address, n.Phone, n.Email, n.Description, n.Website,
		n.Established, n.Facebook, n.Instagram, n.YouTube, n.Spotify, n.StateID,
		n.CreatedAt, n.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (n *Church) UpdateChurch() error {
	db := sqlite.GetDBConnection()
	n.UpdatedAt = time.Now().UTC()

	_, err := db.Exec(`
		UPDATE churches 
		SET name = ?, address = ?, phone = ?, email = ?, description = ?, website = ?, 
		    established = ?, facebook = ?, instagram = ?, youtube = ?, spotify = ?, 
		    state_id = ?, updated_at = ?
		WHERE id = ?
	`, n.Name, n.Address, n.Phone, n.Email, n.Description, n.Website,
		n.Established, n.Facebook, n.Instagram, n.YouTube, n.Spotify,
		n.StateID, n.UpdatedAt, n.ID)
	if err != nil {
		return err
	}

	return nil
}

func (n *Church) AddOrUpdateElder(elder *Elder) error {
	db := sqlite.GetDBConnection()

	// Check if elder has an id
	if elder.ID == nil {
		// Check if the elder already exists

	}

	// Insert new elder
	_, err := db.Exec(`
		INSERT INTO elders (name, email, phone, picture, church_id)
		VALUES (?, ?, ?, ?, ?)
	`, elder.Name, elder.Email, elder.Phone, elder.Picture, n.ID)
	if err != nil {
		return err
	}

	return nil
}

func (n *Church) DeleteChurch(id int) error {
	db := sqlite.GetDBConnection()
	_, err := db.Exec("DELETE FROM churches WHERE id = ?", id)
	return err
}

func (n *Church) GetChurchesByStateID(stateID int) ([]Church, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query(`
		SELECT id, name, address, phone, email, description, website, established, 
		       facebook, instagram, youtube, spotify, state_id 
		FROM churches 
		WHERE state_id = ?
		ORDER BY name
	`, stateID)
	if err != nil {
		return []Church{}, err
	}
	defer rows.Close()

	result := []Church{}
	for rows.Next() {
		elem := Church{}
		err := rows.Scan(
			&elem.ID, &elem.Name, &elem.Address, &elem.Phone, &elem.Email,
			&elem.Description, &elem.Website, &elem.Established, &elem.Facebook,
			&elem.Instagram, &elem.YouTube, &elem.Spotify, &elem.StateID,
		)
		if err != nil {
			continue
		}

		// Get state information
		if state, err := (&State{}).GetStateByID(elem.StateID); err == nil && state != nil {
			elem.State = state
		}

		result = append(result, elem)
	}

	return result, nil
}
