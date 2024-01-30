package handler

import (
	"database/sql"
	"db-final-project/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TransferRequest struct {
	SourceIdentifier      string  `json:"source_identifier"`
	DestinationIdentifier string  `json:"destination_identifier"`
	Amount                float64 `json:"amount"`
}

type TransferResponse struct {
	Message      string `json:"message"`
	TrackingCode string `json:"tracking_code"`
}

func TransferMoney(c echo.Context, db *sql.DB, transferType string) error {
	var requestData TransferRequest

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to bind request data"})
	}

	sourceCardNumber := requestData.SourceIdentifier
	destinationCardNumber := requestData.DestinationIdentifier
	amount := requestData.Amount

	// Perform the money transfer based on the transferType
	trackingCode, err := transferMoney(db, sourceCardNumber, destinationCardNumber, amount, transferType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Money transfer failed", "error": err.Error()})
	}

	response := TransferResponse{
		Message:      "Money transferred successfully",
		TrackingCode: trackingCode,
	}

	return c.JSON(http.StatusOK, response)
}

// Modify the TransferMoneyCardToCard function
func TransferMoneyCardToCard(c echo.Context, db *sql.DB) error {
	return TransferMoney(c, db, "CardToCard")
}

// Modify the TransferMoneySATNA function
func TransferMoneySATNA(c echo.Context, db *sql.DB) error {
	return TransferMoney(c, db, "SATNA")
}

// Modify the TransferMoneyPAYA function
func TransferMoneyPAYA(c echo.Context, db *sql.DB) error {
	return TransferMoney(c, db, "PAYA")
}

func transferMoney(db *sql.DB, sourceCardNumber, destinationCardNumber string, amount float64, transferType string) (string, error) {
	// Check if the source and destination card numbers are valid
	var sourceAccountID, destinationAccountID int

	sourceQuery := "SELECT account_id FROM accounts WHERE card_number = ?"
	err := db.QueryRow(sourceQuery, sourceCardNumber).Scan(&sourceAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("Source card number not found")
		}
		return "", err
	}

	destinationQuery := "SELECT account_id FROM accounts WHERE card_number = ?"
	err = db.QueryRow(destinationQuery, destinationCardNumber).Scan(&destinationAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("Destination card number not found")
		}
		return "", err
	}

	// Check if the source account has enough balance for the transfer
	var sourceBalance float64
	query := "SELECT balance FROM accounts WHERE account_id = ?"
	err = db.QueryRow(query, sourceAccountID).Scan(&sourceBalance)
	if err != nil {
		return "", err
	}

	if sourceBalance < amount {
		return "", fmt.Errorf("Insufficient balance in the source account")
	}

	// Check daily transaction limits based on transferType
	if exceedsDailyLimit(db, sourceCardNumber, amount, transferType) {
		return "", fmt.Errorf("Daily transaction limit exceeded for this transfer type")
	}

	// Start a database transaction
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	// Deduct the amount from the source account
	updateSourceQuery := "UPDATE accounts SET balance = balance - ? WHERE account_id = ?"
	_, err = tx.Exec(updateSourceQuery, amount, sourceAccountID)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// Add the amount to the destination account
	updateDestinationQuery := "UPDATE accounts SET balance = balance + ? WHERE account_id = ?"
	_, err = tx.Exec(updateDestinationQuery, amount, destinationAccountID)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// Generate a tracking code for the transaction (you may use a different method)
	trackingCode := util.GenerateTrackingCode()

	// Insert the transaction record into the transactions table with the correct transfer_type
	insertTransactionQuery := `
        INSERT INTO transactions (source_account_id, destination_account_id, amount, transfer_type, tracking_code, status)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err = tx.Exec(insertTransactionQuery, sourceAccountID, destinationAccountID, amount, transferType, trackingCode, 1)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	return trackingCode, nil
}

func exceedsDailyLimit(db *sql.DB, sourceCardNumber string, amount float64, transferType string) bool {
	// Calculate the total daily transaction sum for the user and transfer type
	var totalSum float64
	query := `
        SELECT COALESCE(SUM(transaction_sum), 0)
        FROM user_daily_transaction_sum
        WHERE user_id = (SELECT user_id FROM accounts WHERE card_number = ?)
            AND transaction_date = CURDATE()
            AND transfer_type = ?
    `
	err := db.QueryRow(query, sourceCardNumber, transferType).Scan(&totalSum)
	if err != nil {
		return false // Assuming no data means within limit
	}

	// Set the daily limit based on the transfer type
	var dailyLimit float64
	switch transferType {
	case "CardToCard":
		dailyLimit = 10000000.00 // 10 million for CardToCard
	case "SATNA":
		dailyLimit = 5000000.00 // 5 million for SATNA
	case "PAYA":
		dailyLimit = 5000000.00 // 5 million for PAYA
	}

	// Check if the total sum plus the current amount exceeds the daily limit
	return totalSum+amount > dailyLimit
}
