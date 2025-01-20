package database

import (
	"API/models"
	"database/sql"
	"fmt"
)

func CreateTrip(trip *models.Trip) error {
	query := `INSERT INTO trips (user_id, start_address, end_address, distance_km, mode_id, carbon_impact_kg, trip_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := DbInstance.DB.Exec(query, trip.UserID, trip.StartAddress, trip.EndAddress, trip.DistanceKm, trip.ModeID, trip.CarbonImpactKg, trip.TripDate, trip.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create trip: %w", err)
	}
	return nil
}

func GetTripByID(tripID int) (*models.Trip, error) {
	query := `SELECT trip_id, user_id, start_address, end_address, distance_km, mode_id, carbon_impact_kg, trip_date, created_at FROM trips WHERE trip_id = $1`
	row := DbInstance.DB.QueryRow(query, tripID)
	trip := &models.Trip{}
	if err := row.Scan(&trip.TripID, &trip.UserID, &trip.StartAddress, &trip.EndAddress, &trip.DistanceKm, &trip.ModeID, &trip.CarbonImpactKg, &trip.TripDate, &trip.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trip not found")
		}
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}
	return trip, nil
}

func UpdateTrip(trip *models.Trip) error {
	query := `UPDATE trips SET user_id = $1, start_address = $2, end_address = $3, distance_km = $4, mode_id = $5, carbon_impact_kg = $6, trip_date = $7, created_at = $8 WHERE trip_id = $9`
	_, err := DbInstance.DB.Exec(query, trip.UserID, trip.StartAddress, trip.EndAddress, trip.DistanceKm, trip.ModeID, trip.CarbonImpactKg, trip.TripDate, trip.CreatedAt, trip.TripID)
	if err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}
	return nil
}

func DeleteTrip(tripID int) error {
	query := `DELETE FROM trips WHERE trip_id = $1`
	_, err := DbInstance.DB.Exec(query, tripID)
	if err != nil {
		return fmt.Errorf("failed to delete trip: %w", err)
	}
	return nil
}
