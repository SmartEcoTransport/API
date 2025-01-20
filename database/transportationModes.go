package database

import (
	"API/models"
	"database/sql"
	"fmt"
)

func GetTransportationModeByID(modeID int) (*models.TransportationMode, error) {
	query := `SELECT mode_id, mode_name, description FROM transportationmodes WHERE mode_id = $1`
	row := DbInstance.DB.QueryRow(query, modeID)
	mode := &models.TransportationMode{}
	if err := row.Scan(&mode.ModeID, &mode.ModeName, &mode.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("mode not found")
		}
		return nil, fmt.Errorf("failed to get mode: %w", err)
	}
	return mode, nil
}

func GetAllTransportationModes() ([]*models.TransportationMode, error) {
	query := `SELECT mode_id, mode_name, description FROM transportationmodes`
	rows, err := DbInstance.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get modes: %w", err)
	}
	defer rows.Close()
	var modes []*models.TransportationMode
	for rows.Next() {
		mode := &models.TransportationMode{}
		if err := rows.Scan(&mode.ModeID, &mode.ModeName, &mode.Description); err != nil {
			return nil, fmt.Errorf("failed to get mode: %w", err)
		}
		modes = append(modes, mode)
	}
	// if no modes are found return an empty slice
	if len(modes) == 0 {
		return []*models.TransportationMode{}, nil
	}
	return modes, nil
}
