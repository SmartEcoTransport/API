package database

import (
	"API/models"
	"API/utils"
	"database/sql"
	"fmt"
	"time"
)

func RegisterTrip(startAddress, endAddress, carBrand, carModel string, distanceKm float64, modeID int, user_id int, tripDate string) error {
	tripTime := time.Now()
	if tripDate != "" {
		// convert the date string to a time.Time
		var err error
		tripTime, err = utils.ConvertStringToTime(tripDate)
		if err != nil {
			return fmt.Errorf("failed to convert trip date: %w", err)
		}
	}
	trip := &models.Trip{
		UserID:   user_id,
		ModeID:   modeID,
		TripDate: tripTime,
	}

	// if the distance is 0 then use the address to calculate the distance
	if distanceKm == 0 {
		d, err := utils.CalculateDistance(startAddress, endAddress)
		if err != nil {
			return fmt.Errorf("failed to calculate distance: %w", err)
		}
		trip.DistanceKm = &d
		trip.StartAddress = &startAddress
		trip.EndAddress = &endAddress
	} else {
		trip.DistanceKm = &distanceKm
	}

	carbonImpactKg, err := utils.GetCarbonImpactByMode(modeID, *trip.DistanceKm)
	if err != nil {
		return fmt.Errorf("failed to get carbon impact: %w", err)
	}
	trip.CarbonImpactKg = &carbonImpactKg

	return CreateTrip(trip)
}

func TotalCarbonImpact(userID int) (float64, error) {
	trips, err := GetUserTrips(userID)
	if err != nil {
		return 0, err
	}
	var total float64
	for _, trip := range trips {
		total += *trip.CarbonImpactKg
	}
	return total, nil
}

func AggregateUserTripsByMode(userID int) ([]models.TripsByMode, error) {
	// get all the trips for the user
	trips, err := GetUserTrips(userID)
	if err != nil {
		return nil, err
	}
	// get all the transportation modes used by the user
	var modes map[int]*models.TripsByMode = make(map[int]*models.TripsByMode)
	for _, trip := range trips {
		// if the mode is not in the map then add it
		if _, ok := modes[trip.ModeID]; !ok {
			mode, err := GetTransportationModeByID(trip.ModeID)
			if err != nil {
				return nil, err
			}
			modes[trip.ModeID] = &models.TripsByMode{
				ModeID:        mode.ModeID,
				TotalTrips:    1,
				TotalImpact:   *trip.CarbonImpactKg,
				TotalDistance: *trip.DistanceKm,
			}
		} else {
			tripsByMode := modes[trip.ModeID]
			tripsByMode.TotalTrips++
			tripsByMode.TotalImpact += *trip.CarbonImpactKg
			tripsByMode.TotalDistance += *trip.DistanceKm
		}
	}
	// convert the map to a slice
	var tripsByMode []models.TripsByMode
	for _, mode := range modes {
		tripsByMode = append(tripsByMode, *mode)
	}
	return tripsByMode, nil
}

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
