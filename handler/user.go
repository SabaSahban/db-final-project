package handler

import (
	"database/sql"
	"db-final-project/util"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterNewUser(c echo.Context, db *sql.DB) error {
	type User struct {
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		NationalID string `json:"national_id"`
	}

	var user User

	if err := c.Bind(&user); err != nil {
		return err
	}

	// Check if the national ID and username are unique
	var count int
	query := "SELECT COUNT(*) FROM users WHERE national_id = ? OR username = ?"
	err := db.QueryRow(query, user.NationalID, user.Username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "National ID or username already exists"})
	}

	// Hash the user's password
	hashedPassword, err := util.HashPassword(user.Password)

	// Insert the new user into the database
	insertQuery := "INSERT INTO users (first_name, last_name, username, password_hash, national_id) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(insertQuery, user.FirstName, user.LastName, user.Username, hashedPassword, user.NationalID)
	if err != nil {
		return err
	}

	return c.String(http.StatusCreated, "User registration successful")
}

func LoginUser(c echo.Context, db *sql.DB) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginReq LoginRequest

	if err := c.Bind(&loginReq); err != nil {
		return err
	}

	// Check if the username exists in the database and retrieve the user ID
	var userID int
	var storedPasswordHash string

	query := "SELECT user_id, password_hash FROM users WHERE username = ?"
	err := db.QueryRow(query, loginReq.Username).Scan(&userID, &storedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Username not found"})
		}
		return err
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(loginReq.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid password"})
	}

	// Set the user ID in the context
	c.Set("user_id", userID)

	return c.String(http.StatusOK, "Login successful")
}

func CreateNewAccount(c echo.Context, db *sql.DB) error {
	type AccountRequest struct {
		UserID         int     `json:"user_id"`
		InitialBalance float64 `json:"initial_balance"`
	}

	var accountReq AccountRequest

	if err := c.Bind(&accountReq); err != nil {
		return err
	}

	// Generate random card number and SHEBA number (for simplicity, you may want to use a better method)
	rand.Seed(time.Now().UnixNano())
	cardNumber := fmt.Sprintf("%016d", rand.Int63n(9999999999999999)) // Card number with up to 16 digits

	maxInt64 := int64(9223372036854775807)                       // Maximum int64 value
	shebaNumber := fmt.Sprintf("IR%022d", rand.Int63n(maxInt64)) // SHEBA number with up to 22 digits

	// Insert the new account into the database
	insertQuery := "INSERT INTO accounts (user_id, card_number, sheba_number, balance) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(insertQuery, accountReq.UserID, cardNumber, shebaNumber, accountReq.InitialBalance)
	if err != nil {
		return err
	}

	response := map[string]string{
		"card_number":  cardNumber,
		"sheba_number": shebaNumber,
	}

	return c.JSON(http.StatusCreated, response)
}
