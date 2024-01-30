// handler/transfer_handler.go

package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Transaction struct {
	TransactionID        int       `json:"transaction_id"`
	SourceAccountID      int       `json:"source_account_id"`
	DestinationAccountID int       `json:"destination_account_id"`
	Amount               float64   `json:"amount"`
	TransferType         string    `json:"transfer_type"`
	TransactionTime      time.Time `json:"transaction_time"`
	TrackingCode         string    `json:"tracking_code"`
	Status               int       `json:"status"`
}

func RetrieveLastNTransactions(c echo.Context, db *sql.DB) error {
	// Parse the account identifier (SHEBA or card number) and n from the request
	accountIdentifier := c.Param("accountIdentifier")
	nStr := c.Param("n")

	// Convert n to an integer
	n, err := strconv.Atoi(nStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid value for n"})
	}

	// Retrieve the user's account ID based on SHEBA or card number
	var accountID int
	query := `
		SELECT account_id
		FROM accounts
		WHERE sheba_number = ? OR card_number = ?
	`
	err = db.QueryRow(query, accountIdentifier, accountIdentifier).Scan(&accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Account not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve account"})
	}

	// Retrieve the last n transactions for the account by account ID
	var transactions []Transaction
	query = `
		SELECT transaction_id, source_account_id, destination_account_id, amount, transfer_type, transaction_time, tracking_code, status
		FROM transactions
		WHERE source_account_id = ? OR destination_account_id = ?
		ORDER BY transaction_time DESC
		LIMIT ?
	`
	rows, err := db.Query(query, accountID, accountID, n)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve transactions"})
	}
	defer rows.Close()

	for rows.Next() {
		var transaction Transaction
		var transactionTimeStr string
		err := rows.Scan(
			&transaction.TransactionID,
			&transaction.SourceAccountID,
			&transaction.DestinationAccountID,
			&transaction.Amount,
			&transaction.TransferType,
			&transactionTimeStr,
			&transaction.TrackingCode,
			&transaction.Status,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan transaction"})
		}

		// Convert transactionTimeStr to time.Time
		transaction.TransactionTime, err = time.Parse("2006-01-02 15:04:05", transactionTimeStr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse transaction time"})
		}

		transactions = append(transactions, transaction)
	}

	return c.JSON(http.StatusOK, transactions)
}

func VerifyTransaction(c echo.Context, db *sql.DB) error {
	// Parse the tracking code from the request
	trackingCode := c.Param("trackingCode")

	// Validate the tracking code
	if trackingCode == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tracking code"})
	}

	// Retrieve the transaction details based on the tracking code
	var transaction Transaction
	var transactionTimeStr string
	query := `
		SELECT transaction_id, source_account_id, destination_account_id, amount, transfer_type, transaction_time, tracking_code, status
		FROM transactions
		WHERE tracking_code = ?
	`
	err := db.QueryRow(query, trackingCode).Scan(
		&transaction.TransactionID,
		&transaction.SourceAccountID,
		&transaction.DestinationAccountID,
		&transaction.Amount,
		&transaction.TransferType,
		&transactionTimeStr,
		&transaction.TrackingCode,
		&transaction.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Transaction not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve transaction"})
	}

	// Convert transactionTimeStr to time.Time
	transaction.TransactionTime, err = time.Parse("2006-01-02 15:04:05", transactionTimeStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse transaction time"})
	}

	return c.JSON(http.StatusOK, transaction)
}
