package server

import (
	"API/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"os"
	"strconv"
	"time"
)

func StartAndInitializeServer() {

	// Initialize Fiber app
	app := fiber.New()

	// Register routes
	app.Post("/register", registerHandler)

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/login", loginHandler)
	auth.Post("/login/cookie", loginCookieHandler)

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
	// Get JWT from cookie
	jwtCookie := c.Cookies("jwt")
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
