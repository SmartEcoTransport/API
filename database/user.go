package database

import (
	"API/models"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func CheckUserCredentials(email, password string) (int, error) {
	// return error if email is empty
	if email == "" {
		return 0, errors.New("email is empty")
	}

	var query string
	query = `SELECT user_id, password_hash FROM Users WHERE email = $1`
	row := DbInstance.DB.QueryRow(query, email)
	var userID int
	var passwordHash string
	if err := row.Scan(&userID, &passwordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("user not found")
		}
		log.Println("Error retrieving user:", err)
		return 0, err
	}

	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return 0, errors.New("incorrect password")
	}
	return userID, nil
}

func RegisterUserFromEmail(email, username, password string) (int, error) {

	// return error if email is empty
	if email == "" {
		return 0, errors.New("email is empty")
	}

	// return error if username is empty
	if username == "" {
		return 0, errors.New("username is empty")
	}

	// return error if password is empty
	if password == "" {
		return 0, errors.New("password is empty")
	}

	// check if email already exists
	query := `SELECT user_id FROM Users WHERE email = $1`
	row := DbInstance.DB.QueryRow(query, email)
	var userID int
	if err := row.Scan(&userID); err == nil {
		return 0, errors.New("email already exists")
	}

	// check if username already exists
	query = `SELECT user_id FROM Users WHERE username = $1`
	row = DbInstance.DB.QueryRow(query, username)
	if err := row.Scan(&userID); err == nil {
		return 0, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return 0, err
	}
	user := models.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	return CreateUser(user)
}

// CreateUser creates a new user in the database
func CreateUser(user models.User) (int, error) {
	query := `INSERT INTO Users (email, username, password_hash, google_id, github_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id`

	var userID int
	err := DbInstance.DB.QueryRow(query,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.GoogleID,
		user.GithubID,
		time.Now(),
		time.Now(),
	).Scan(&userID)
	if err != nil {
		log.Println("Error creating user:", err)
		return 0, err
	}

	return userID, nil
}

// GetUser retrieves a user by their ID
func GetUser(userID int) (*models.User, error) {
	query := `SELECT user_id, email, username, password_hash, google_id, github_id, created_at, updated_at 
		FROM Users WHERE user_id = $1`

	row := DbInstance.DB.QueryRow(query, userID)
	var user models.User
	if err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.GoogleID,
		&user.GithubID,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		log.Println("Error retrieving user:", err)
		return nil, err
	}

	return &user, nil
}

func GetUserTrips(userID int) ([]models.Trip, error) {
	// Check if user exists
	_, err := GetUser(userID)
	if err != nil {
		return nil, err
	}

	query := `SELECT trip_id, user_id, start_address, end_address, distance_km, mode_id, carbon_impact_kg, trip_date, created_at 
		FROM Trips WHERE user_id = $1`

	rows, err := DbInstance.DB.Query(query, userID)
	if err != nil {
		log.Println("Error retrieving trips:", err)
		return nil, err
	}
	defer rows.Close()

	var trips []models.Trip
	for rows.Next() {
		var trip models.Trip
		if err := rows.Scan(
			&trip.TripID,
			&trip.UserID,
			&trip.StartAddress,
			&trip.EndAddress,
			&trip.DistanceKm,
			&trip.ModeID,
			&trip.CarbonImpactKg,
			&trip.TripDate,
			&trip.CreatedAt,
		); err != nil {
			log.Println("Error scanning trip row:", err)
			return nil, err
		}
		trips = append(trips, trip)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error after iterating rows:", err)
		return nil, err
	}

	return trips, nil
}

// UpdateUser updates an existing user's details
func UpdateUser(user models.User) error {
	query := `UPDATE Users SET email = $1, username = $2, password_hash = $3, google_id = $4, github_id = $5, updated_at = $6 
		WHERE user_id = $7`

	_, err := DbInstance.DB.Exec(query,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.GoogleID,
		user.GithubID,
		time.Now(),
		user.UserID,
	)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}

	return nil
}

// DeleteUser deletes a user by their ID
func DeleteUser(userID int) error {
	query := `DELETE FROM Users WHERE user_id = $1`

	_, err := DbInstance.DB.Exec(query, userID)
	if err != nil {
		log.Println("Error deleting user:", err)
		return err
	}

	return nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers() ([]models.User, error) {
	query := `SELECT user_id, email, username, password_hash, google_id, github_id, created_at, updated_at 
		FROM Users`

	rows, err := DbInstance.DB.Query(query)
	if err != nil {
		log.Println("Error retrieving users:", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.UserID,
			&user.Email,
			&user.Username,
			&user.PasswordHash,
			&user.GoogleID,
			&user.GithubID,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			log.Println("Error scanning user row:", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error after iterating rows:", err)
		return nil, err
	}

	return users, nil
}
