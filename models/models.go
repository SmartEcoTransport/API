package models

import "time"

// User represents the Users table
type User struct {
	UserID       int       `json:"user_id" db:"user_id"`
	Email        string    `json:"email" db:"email"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	GoogleID     *string   `json:"google_id,omitempty" db:"google_id"`
	GithubID     *string   `json:"github_id,omitempty" db:"github_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// TransportationMode represents the TransportationModes table
type TransportationMode struct {
	ModeID      int     `json:"mode_id" db:"mode_id"`
	ModeName    string  `json:"mode_name" db:"mode_name"`
	Description *string `json:"description,omitempty" db:"description"`
}

// Trip represents the Trips table
type Trip struct {
	TripID         int       `json:"trip_id" db:"trip_id"`
	UserID         int       `json:"user_id" db:"user_id"`
	StartAddress   *string   `json:"start_address,omitempty" db:"start_address"`
	EndAddress     *string   `json:"end_address,omitempty" db:"end_address"`
	DistanceKm     *float64  `json:"distance_km,omitempty" db:"distance_km"`
	ModeID         int       `json:"mode_id" db:"mode_id"`
	CarbonImpactKg *float64  `json:"carbon_impact_kg,omitempty" db:"carbon_impact_kg"`
	TripDate       time.Time `json:"trip_date" db:"trip_date"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// Challenge represents the Challenges table
type Challenge struct {
	ChallengeID int       `json:"challenge_id" db:"challenge_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	StartDate   time.Time `json:"start_date" db:"start_date"`
	EndDate     time.Time `json:"end_date" db:"end_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ChallengeParticipation represents the ChallengeParticipation table
type ChallengeParticipation struct {
	ParticipationID int      `json:"participation_id" db:"participation_id"`
	UserID          int      `json:"user_id" db:"user_id"`
	ChallengeID     int      `json:"challenge_id" db:"challenge_id"`
	Progress        *float64 `json:"progress,omitempty" db:"progress"`
	Completed       bool     `json:"completed" db:"completed"`
}

// Recommendation represents the Recommendations table
type Recommendation struct {
	RecommendationID int       `json:"recommendation_id" db:"recommendation_id"`
	UserID           int       `json:"user_id" db:"user_id"`
	Message          string    `json:"message" db:"message"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}
