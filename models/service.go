package models

import (
	"database/sql"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
)

type Service struct {
	ID          *int   `json:"id"`
	ServiceType string `json:"service_type"`
	Schedule    string `json:"schedule"`
	ChurchID    int    `json:"church_id"`
}

func (n *Service) GetServicesByChurchID(churchID int) ([]Service, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, service_type, schedule, church_id FROM church_services WHERE church_id = ? ORDER BY service_type", churchID)
	if err != nil {
		return []Service{}, err
	}
	defer rows.Close()

	result := []Service{}
	for rows.Next() {
		elem := Service{}
		err := rows.Scan(&elem.ID, &elem.ServiceType, &elem.Schedule, &elem.ChurchID)
		if err != nil {
			continue
		}
		result = append(result, elem)
	}

	return result, nil
}

func (n *Service) GetServiceByID(id int) (*Service, error) {
	db := sqlite.GetDBConnection()
	row := db.QueryRow("SELECT id, service_type, schedule, church_id FROM church_services WHERE id = ?", id)

	elem := Service{}
	err := row.Scan(&elem.ID, &elem.ServiceType, &elem.Schedule, &elem.ChurchID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &elem, nil
}

func (n *Service) CreateService() error {
	db := sqlite.GetDBConnection()

	_, err := db.Exec(`
		INSERT INTO church_services (service_type, schedule, church_id)
		VALUES (?, ?, ?)
	`, n.ServiceType, n.Schedule, n.ChurchID)

	if err != nil {
		return err
	}

	return nil
}

func (n *Service) UpdateService() error {
	db := sqlite.GetDBConnection()

	_, err := db.Exec(`
		UPDATE church_services
		SET service_type = ?, schedule = ?, church_id = ?
		WHERE id = ?
	`, n.ServiceType, n.Schedule, n.ChurchID, n.ID)

	if err != nil {
		return err
	}

	return nil
}

func (n *Service) DeleteService() error {
	db := sqlite.GetDBConnection()
	_, err := db.Exec("DELETE FROM church_services WHERE id = ?", n.ID)
	return err
}
