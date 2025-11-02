package models

import (
	"database/sql"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
)

type Elder struct {
	ID       *int    `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
	Picture  *string `json:"picture"`
	ChurchID int     `json:"church_id"`
}

func (n *Elder) GetEldersByChurchID(churchID int) ([]Elder, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, name, email, phone, picture, church_id FROM elders WHERE church_id = ? ORDER BY name", churchID)
	if err != nil {
		return []Elder{}, err
	}
	defer rows.Close()

	result := []Elder{}
	for rows.Next() {
		elem := Elder{}
		err := rows.Scan(&elem.ID, &elem.Name, &elem.Email, &elem.Phone, &elem.Picture, &elem.ChurchID)
		if err != nil {
			continue
		}
		result = append(result, elem)
	}

	return result, nil
}

func (n *Elder) GetElderByID(id int) (*Elder, error) {
	db := sqlite.GetDBConnection()
	row := db.QueryRow("SELECT id, name, email, phone, picture, church_id FROM elders WHERE id = ?", id)

	elem := Elder{}
	err := row.Scan(&elem.ID, &elem.Name, &elem.Email, &elem.Phone, &elem.Picture, &elem.ChurchID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &elem, nil
}

func (n *Elder) CreateElder() error {
	db := sqlite.GetDBConnection()

	result, err := db.Exec(`
		INSERT INTO elders (name, email, phone, picture, church_id)
		VALUES (?, ?, ?, ?, ?)
	`, n.Name, n.Email, n.Phone, n.Picture, n.ChurchID)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	idInt := int(id)
	n.ID = &idInt

	return nil
}

func (n *Elder) UpdateElder() error {
	db := sqlite.GetDBConnection()

	_, err := db.Exec(`
		UPDATE elders
		SET name = ?, email = ?, phone = ?, picture = ?, church_id = ?
		WHERE id = ?
	`, n.Name, n.Email, n.Phone, n.Picture, n.ChurchID, n.ID)

	if err != nil {
		return err
	}

	return nil
}

func (n *Elder) DeleteElder() error {
	db := sqlite.GetDBConnection()
	_, err := db.Exec("DELETE FROM elders WHERE id = ?", n.ID)
	return err
}
