package server

import (
	"API/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
	"strconv"
	"time"
)

func StartAndInitializeServer() {

	// Initialize Fiber app
	app := fiber.New()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// Register routes
	app.Post("/register", registerHandler)

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/login", loginHandler)
	auth.Post("/login/cookie", loginCookieHandler)

	trips := app.Group("/trips")
	trips.Use(AuthMiddleware)
	trips.Get("/", tripsHandler)
	trips.Post("/", createTripHandler)
	trips.Get("/impactgraphday", tripsImpactGraphDayHandler)
	trips.Get("/impactgraphmonth", tripsImpactGraphMonthHandler)
	trips.Get("/aggregation", tripsAggregationHandler)
	trips.Get("/impact", totalImpactHandler)

	// Transport modes routes
	transportation := app.Group("/transportation")
	transportation.Get("/", transportationModesHandler)
	transportation.Get("/:mode_id", transportationModeHandler)

	// read port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))

}

type Point struct {
	X int     `json:"x"`
	Y float64 `json:"y"`
}

func tripsImpactGraphDayHandler(c *fiber.Ctx) error {
	// 1 year graph with 1 datapoint per day

	temp := c.Locals("user").(float64)
	userID := int(temp)
	trips, err := database.GetUserTrips(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	points := make([]Point, 366)
	for _, trip := range trips {
		// if the trip date is within the last year
		if trip.TripDate.After(time.Now().AddDate(-1, 0, 0)) {
			points[trip.TripDate.YearDay()].Y += *trip.CarbonImpactKg
			points[trip.TripDate.YearDay()].X = trip.TripDate.YearDay()

		}
	}

	// now that we have the impact for each day cascade the values to have a cumulative impact
	newPoints := make([]Point, 366)
	var sum float64
	for i, point := range points {
		sum += point.Y
		newPoints[i].Y = sum
		newPoints[i].X = i + 1
	}

	return c.JSON(fiber.Map{"points": newPoints})

}

func tripsImpactGraphMonthHandler(c *fiber.Ctx) error {
	// 1 year graph with 1 datapoint per month

	temp := c.Locals("user").(float64)
	userID := int(temp)

	trips, err := database.GetUserTrips(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	points := make([]Point, 13)
	for _, trip := range trips {
		// if the trip date is within the last year
		if trip.TripDate.After(time.Now().AddDate(-1, 0, 0)) {
			points[trip.TripDate.Month()].Y += *trip.CarbonImpactKg
			points[trip.TripDate.Month()].X = int(trip.TripDate.Month())
		}
	}
	// now that we have the impact for each month cascade the values to have a cumulative impact
	newPoints := make([]Point, 13)
	var sum float64
	for i, point := range points {
		sum += point.Y
		newPoints[i].Y = sum
		newPoints[i].X = i + 1
	}

	return c.JSON(fiber.Map{"points": newPoints})
}

func totalImpactHandler(c *fiber.Ctx) error {
	// Get user ID from JWT
	temp := c.Locals("user").(float64)
	userID := int(temp)

	if userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	// Get total carbon impact for the user
	totalImpact, err := database.TotalCarbonImpact(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"total_impact": totalImpact})
}
func createTripHandler(c *fiber.Ctx) error {
	// Parse request body
	var req struct {
		StartAddress string  `json:"start_address"`
		EndAddress   string  `json:"end_address"`
		CarBrand     string  `json:"car_brand"`
		CarModel     string  `json:"car_model"`
		DistanceKm   float64 `json:"distance_km"`
		ModeID       int     `json:"mode_id"`
		TripDate     string  `json:"trip_date"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Get user ID from JWT
	temp := c.Locals("user").(float64)
	userID := int(temp)

	if userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	// if distance is 0, and start and end address not provided return error
	if req.DistanceKm == 0 && (req.StartAddress == "" && req.EndAddress == "") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no distance or address provided"})
	}

	// if mode ID is 0 return error
	if req.ModeID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "mode_id is required"})
	}

	// Register trip in the database
	err := database.RegisterTrip(req.StartAddress, req.EndAddress, req.CarBrand, req.CarModel, req.DistanceKm, req.ModeID, userID, req.TripDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "trip registered"})

}

func tripsAggregationHandler(c *fiber.Ctx) error {
	// Get user ID from JWT
	temp := c.Locals("user").(float64)
	userID := int(temp)

	if userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	// Get aggregated trips for the user
	trips, err := database.AggregateUserTripsByMode(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"trips": trips})

}

func tripsHandler(c *fiber.Ctx) error {
	// Get user ID from JWT
	temp := c.Locals("user").(float64)
	userID := int(temp)

	if userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	// Get all trips for the user
	trips, err := database.GetUserTrips(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"trips": trips})

}

func transportationModeHandler(c *fiber.Ctx) error {
	// Get mode ID from URL
	modeID := c.Params("mode_id")
	if modeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "mode_id is required"})
	}

	// Convert mode ID to integer
	modeIDInt, err := strconv.Atoi(modeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid mode_id"})
	}

	// Get transportation mode by ID
	mode, err := database.GetTransportationModeByID(modeIDInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"mode": mode})

}

func transportationModesHandler(c *fiber.Ctx) error {
	// Get all transportation modes
	modes, err := database.GetAllTransportationModes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"modes": modes})
}

func AuthMiddleware(c *fiber.Ctx) error {
	// Get JWT from cookie or Authorization header
	jwtCookie := c.Cookies("jwt")
	if jwtCookie == "" {
		// remove the bearer from the Authorization header
		jwtCookie = c.Get("Authorization")
		if jwtCookie != "" {
			jwtCookie = jwtCookie[len("Bearer "):]
		}
	}
	if jwtCookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Parse JWT
	token, err := jwt.Parse(jwtCookie, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Check if token is valid
	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	c.Locals("user", token.Claims.(jwt.MapClaims)["user_id"])

	return c.Next()
}

var jwtSecret = []byte("your_secret_key")

func registerHandler(c *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Register user in the database
	userID, err := database.RegisterUserFromEmail(req.Email, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user registered", "user_id": userID})
}

func loginHandler(c *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Validate user credentials
	userID, err := database.CheckUserCredentials(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}

func loginCookieHandler(c *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Validate user credentials
	userID, err := database.CheckUserCredentials(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	// Set JWT as a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})
	// redirect to the home page
	return c.Redirect("/")
}
