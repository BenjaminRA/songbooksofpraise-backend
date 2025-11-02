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
	State       *State    `json:"state,omitempty"`          // Not in database, but used in API responses
	Elders      []Elder   `json:"elders,omitempty"`         // Not in database, but used in API responses
	Services    []Service `json:"church_services,omitempty"` // Not in database, but used in API responses
	Events      []Event   `json:"church_events,omitempty"`   // Not in database, but used in API responses
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

	result, err := db.Exec(`
		INSERT INTO churches (name, address, phone, email, description, website, established, 
		                     facebook, instagram, youtube, spotify, state_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, n.Name, n.Address, n.Phone, n.Email, n.Description, n.Website,
		n.Established, n.Facebook, n.Instagram, n.YouTube, n.Spotify, n.StateID,
		n.CreatedAt, n.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	n.ID = int(id)

	// Create elders
	for i := range n.Elders {
		n.Elders[i].ChurchID = n.ID
		if err := n.Elders[i].CreateElder(); err != nil {
			return err
		}
	}

	// Create services
	for i := range n.Services {
		n.Services[i].ChurchID = n.ID
		if err := n.Services[i].CreateService(); err != nil {
			return err
		}
	}

	// Create events
	for i := range n.Events {
		n.Events[i].ChurchID = n.ID
		if err := n.Events[i].CreateEvent(); err != nil {
			return err
		}
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

	// Handle Elders: Update/Create/Delete
	if err := n.syncElders(); err != nil {
		return err
	}

	// Handle Services: Update/Create/Delete
	if err := n.syncServices(); err != nil {
		return err
	}

	// Handle Events: Update/Create/Delete
	if err := n.syncEvents(); err != nil {
		return err
	}

	return nil
}

func (n *Church) syncElders() error {
	db := sqlite.GetDBConnection()

	// Get existing elders from database
	existingElders, err := (&Elder{}).GetEldersByChurchID(n.ID)
	if err != nil {
		return err
	}

	// Create a map of existing elder IDs
	existingIDs := make(map[int]bool)
	for _, elder := range existingElders {
		if elder.ID != nil {
			existingIDs[*elder.ID] = true
		}
	}

	// Track which IDs are in the new data
	newIDs := make(map[int]bool)

	// Update or create elders
	for i := range n.Elders {
		n.Elders[i].ChurchID = n.ID
		if n.Elders[i].ID != nil && *n.Elders[i].ID > 0 {
			// Update existing elder
			newIDs[*n.Elders[i].ID] = true
			if err := n.Elders[i].UpdateElder(); err != nil {
				return err
			}
		} else {
			// Create new elder
			if err := n.Elders[i].CreateElder(); err != nil {
				return err
			}
		}
	}

	// Delete elders that are not in the new data
	for id := range existingIDs {
		if !newIDs[id] {
			_, err := db.Exec("DELETE FROM elders WHERE id = ?", id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *Church) syncServices() error {
	db := sqlite.GetDBConnection()

	// Get existing services from database
	existingServices, err := (&Service{}).GetServicesByChurchID(n.ID)
	if err != nil {
		return err
	}

	// Create a map of existing service IDs
	existingIDs := make(map[int]bool)
	for _, service := range existingServices {
		if service.ID != nil {
			existingIDs[*service.ID] = true
		}
	}

	// Track which IDs are in the new data
	newIDs := make(map[int]bool)

	// Update or create services
	for i := range n.Services {
		n.Services[i].ChurchID = n.ID
		if n.Services[i].ID != nil && *n.Services[i].ID > 0 {
			// Update existing service
			newIDs[*n.Services[i].ID] = true
			if err := n.Services[i].UpdateService(); err != nil {
				return err
			}
		} else {
			// Create new service
			if err := n.Services[i].CreateService(); err != nil {
				return err
			}
		}
	}

	// Delete services that are not in the new data
	for id := range existingIDs {
		if !newIDs[id] {
			_, err := db.Exec("DELETE FROM church_services WHERE id = ?", id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *Church) syncEvents() error {
	db := sqlite.GetDBConnection()

	// Get existing events from database
	existingEvents, err := (&Event{}).GetEventsByChurchID(n.ID)
	if err != nil {
		return err
	}

	// Create a map of existing event IDs
	existingIDs := make(map[int]bool)
	for _, event := range existingEvents {
		if event.ID != nil {
			existingIDs[*event.ID] = true
		}
	}

	// Track which IDs are in the new data
	newIDs := make(map[int]bool)

	// Update or create events
	for i := range n.Events {
		n.Events[i].ChurchID = n.ID
		if n.Events[i].ID != nil && *n.Events[i].ID > 0 {
			// Update existing event
			newIDs[*n.Events[i].ID] = true
			if err := n.Events[i].UpdateEvent(); err != nil {
				return err
			}
		} else {
			// Create new event
			if err := n.Events[i].CreateEvent(); err != nil {
				return err
			}
		}
	}

	// Delete events that are not in the new data
	for id := range existingIDs {
		if !newIDs[id] {
			_, err := db.Exec("DELETE FROM church_events WHERE id = ?", id)
			if err != nil {
				return err
			}
		}
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
