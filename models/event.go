package models

import (
	"database/sql"
	"time"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
)

type Event struct {
	ID          *int       `json:"id"`
	Name        string     `json:"name"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Location    *string    `json:"location"`
	Image       *string    `json:"image"`
	Color       *string    `json:"color"`
	Recurrent   bool       `json:"recurrent"`
	Description *string    `json:"description"`
	ChurchID    int        `json:"church_id"`
}

func (n *Event) GetEventsByChurchID(churchID int) ([]Event, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query(`
		SELECT id, name, start_date, end_date, location, image, color, recurrent, description, church_id 
		FROM church_events 
		WHERE church_id = ? 
		ORDER BY start_date DESC
	`, churchID)
	if err != nil {
		return []Event{}, err
	}
	defer rows.Close()

	result := []Event{}
	for rows.Next() {
		elem := Event{}
		err := rows.Scan(
			&elem.ID, &elem.Name, &elem.StartDate, &elem.EndDate, &elem.Location,
			&elem.Image, &elem.Color, &elem.Recurrent, &elem.Description, &elem.ChurchID,
		)
		if err != nil {
			continue
		}
		result = append(result, elem)
	}

	return result, nil
}

func (n *Event) GetEventByID(id int) (*Event, error) {
	db := sqlite.GetDBConnection()
	row := db.QueryRow(`
		SELECT id, name, start_date, end_date, location, image, color, recurrent, description, church_id 
		FROM church_events 
		WHERE id = ?
	`, id)

	elem := Event{}
	err := row.Scan(
		&elem.ID, &elem.Name, &elem.StartDate, &elem.EndDate, &elem.Location,
		&elem.Image, &elem.Color, &elem.Recurrent, &elem.Description, &elem.ChurchID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &elem, nil
}

func (n *Event) CreateEvent() error {
	db := sqlite.GetDBConnection()

	_, err := db.Exec(`
		INSERT INTO church_events (name, start_date, end_date, location, image, color, recurrent, description, church_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, n.Name, n.StartDate, n.EndDate, n.Location, n.Image, n.Color, n.Recurrent, n.Description, n.ChurchID)

	if err != nil {
		return err
	}

	return nil
}

func (n *Event) UpdateEvent() error {
	db := sqlite.GetDBConnection()

	_, err := db.Exec(`
		UPDATE church_events
		SET name = ?, start_date = ?, end_date = ?, location = ?, image = ?, color = ?, recurrent = ?, description = ?, church_id = ?
		WHERE id = ?
	`, n.Name, n.StartDate, n.EndDate, n.Location, n.Image, n.Color, n.Recurrent, n.Description, n.ChurchID, n.ID)

	if err != nil {
		return err
	}

	return nil
}

func (e *Event) DeleteEvent(id int) error {
	db := sqlite.GetDBConnection()
	_, err := db.Exec("DELETE FROM church_events WHERE id = ?", id)
	return err
}
