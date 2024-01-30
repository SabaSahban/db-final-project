package handler

import (
	"database/sql"
	"db-final-project/util"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Transaction struct {
	TransactionID        int
	SourceAccountID      int
	DestinationAccountID int
	Amount               float64
	TransferType         string
	TransactionTime      time.Time
	TrackingCode         string
	Status               int
}

func RegisterNewUser(db *sql.DB) error {
	var firstName, lastName, username, password, nationalID string

	fmt.Print("Enter first name: ")
	fmt.Scanln(&firstName)

	fmt.Print("Enter last name: ")
	fmt.Scanln(&lastName)

	fmt.Print("Enter username: ")
	fmt.Scanln(&username)

	fmt.Print("Enter password: ")
	fmt.Scanln(&password)

	fmt.Print("Enter national ID: ")
	fmt.Scanln(&nationalID)

	// Check if the national ID and username are unique
	var count int
	query := "SELECT COUNT(*) FROM users WHERE national_id = ? OR username = ?"
	err := db.QueryRow(query, nationalID, username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("National ID or username already exists")
	}

	// Hash the user's password
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return err
	}

	// Insert the new user into the database
	insertQuery := "INSERT INTO users (first_name, last_name, username, password_hash, national_id) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(insertQuery, firstName, lastName, username, hashedPassword, nationalID)
	if err != nil {
		return err
	}

	return nil // User registration successful
}

func LoginUser(db *sql.DB) error {
	var username, password string

	fmt.Print("Enter username: ")
	fmt.Scanln(&username)

	fmt.Print("Enter password: ")
	fmt.Scanln(&password)

	// Check if the username exists in the database
	var storedPasswordHash string
	query := "SELECT password_hash FROM users WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&storedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Username not found")
		}
		return err
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(password))
	if err != nil {
		return fmt.Errorf("Invalid password")
	}

	return nil // Login successful
}

func CreateNewAccount(db *sql.DB) error {
	var userID int
	var initialBalance float64

	fmt.Print("Enter user ID: ")
	_, err := fmt.Scanln(&userID)
	if err != nil {
		return err
	}

	fmt.Print("Enter initial balance: ")
	_, err = fmt.Scanln(&initialBalance)
	if err != nil {
		return err
	}

	// Generate random card number and SHEBA number (for simplicity, you may want to use a better method)
	rand.Seed(time.Now().UnixNano())
	cardNumber := fmt.Sprintf("%016d", rand.Int63n(9999999999999999)) // Card number with up to 16 digits

	maxInt64 := int64(9223372036854775807)                       // Maximum int64 value
	shebaNumber := fmt.Sprintf("IR%022d", rand.Int63n(maxInt64)) // SHEBA number with up to 22 digits

	// Insert the new account into the database
	insertQuery := "INSERT INTO accounts (user_id, card_number, sheba_number, balance) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(insertQuery, userID, cardNumber, shebaNumber, initialBalance)
	if err != nil {
		return err
	}

	fmt.Println("card number:", cardNumber, "sheba number:", shebaNumber)

	return nil // Account creation successful
}

func transferMoney(db *sql.DB, sourceIdentifier, destinationIdentifier string, amount float64, transferType string) error {
	// Check if the source account identifier (card number or SHEBA number) exists in the accounts table
	var sourceAccountID, destinationAccountID int
	var sourceColumnName, destinationColumnName string

	if transferType == "CardToCard" {
		sourceColumnName = "card_number"
		destinationColumnName = "card_number"
	} else if transferType == "SATNA" || transferType == "PAYA" {
		sourceColumnName = "card_number"
		destinationColumnName = "sheba_number"
	} else {
		return fmt.Errorf("Invalid transfer type")
	}

	sourceQuery := "SELECT account_id FROM accounts WHERE " + sourceColumnName + " = ?"
	err := db.QueryRow(sourceQuery, sourceIdentifier).Scan(&sourceAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Source account not found")
		}
		return err
	}

	destinationQuery := "SELECT account_id FROM accounts WHERE " + destinationColumnName + " = ?"
	err = db.QueryRow(destinationQuery, destinationIdentifier).Scan(&destinationAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Destination account not found")
		}
		return err
	}

	// Check if the source account has enough balance for the transfer
	var sourceBalance float64
	query := "SELECT balance FROM accounts WHERE account_id = ?"
	err = db.QueryRow(query, sourceAccountID).Scan(&sourceBalance)
	if err != nil {
		return err
	}

	if sourceBalance < amount {
		return fmt.Errorf("Insufficient balance in the source account")
	}

	// Start a database transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Deduct the amount from the source account
	updateSourceQuery := "UPDATE accounts SET balance = balance - ? WHERE account_id = ?"
	_, err = tx.Exec(updateSourceQuery, amount, sourceAccountID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Add the amount to the destination account
	updateDestinationQuery := "UPDATE accounts SET balance = balance + ? WHERE account_id = ?"
	_, err = tx.Exec(updateDestinationQuery, amount, destinationAccountID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Generate a tracking code for the transaction (you may use a different method)
	trackingCode := util.GenerateTrackingCode()

	// Insert the transaction record into the transactions table
	insertTransactionQuery := `
		INSERT INTO transactions (source_account_id, destination_account_id, amount, transfer_type, tracking_code, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(insertTransactionQuery, sourceAccountID, destinationAccountID, amount, transferType, trackingCode, 1)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil // Money transfer successful
}

func TransferMoneyCardToCard(db *sql.DB, sourceCardNumber, destinationCardNumber string, amount float64) error {
	return transferMoney(db, sourceCardNumber, destinationCardNumber, amount, "CardToCard")
}

func TransferMoneySATNA(db *sql.DB, sourceCardNumber, destinationSHEBANumber string, amount float64) error {
	return transferMoney(db, sourceCardNumber, destinationSHEBANumber, amount, "SATNA")
}

func TransferMoneyPAYA(db *sql.DB, sourceCardNumber, destinationSHEBANumber string, amount float64) error {
	return transferMoney(db, sourceCardNumber, destinationSHEBANumber, amount, "PAYA")
}

func RetrieveLastNTransactions(db *sql.DB, accountIdentifier string, n int) ([]Transaction, error) {
	// Check if the account identifier (card number or SHEBA number) exists in the accounts table
	var accountID int
	var accountColumnName string

	if strings.Contains(accountIdentifier, "-") {
		accountColumnName = "sheba_number"
	} else {
		accountColumnName = "card_number"
	}

	query := "SELECT account_id FROM accounts WHERE " + accountColumnName + " = ?"
	err := db.QueryRow(query, accountIdentifier).Scan(&accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Account not found")
		}
		return nil, err
	}

	// Retrieve the last n transactions for the account
	query = `
		SELECT transaction_id, source_account_id, destination_account_id, amount, transfer_type, transaction_time, tracking_code, status
		FROM transactions
		WHERE source_account_id = ? OR destination_account_id = ?
		ORDER BY transaction_time DESC
		LIMIT ?
	`
	rows, err := db.Query(query, accountID, accountID, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var transaction Transaction
		var transactionTimeStr string // Temporary variable to hold transaction_time as string
		err := rows.Scan(
			&transaction.TransactionID,
			&transaction.SourceAccountID,
			&transaction.DestinationAccountID,
			&transaction.Amount,
			&transaction.TransferType,
			&transactionTimeStr, // Scan transaction_time as string
			&transaction.TrackingCode,
			&transaction.Status,
		)
		if err != nil {
			return nil, err
		}

		// Convert transactionTimeStr to time.Time
		transaction.TransactionTime, err = time.Parse("2006-01-02 15:04:05", transactionTimeStr)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func VerifyTransaction(db *sql.DB, trackingCode string) (Transaction, error) {
	// Retrieve the transaction details based on the tracking code
	query := `
		SELECT transaction_id, source_account_id, destination_account_id, amount, transfer_type, transaction_time, tracking_code, status
		FROM transactions
		WHERE tracking_code = ?
	`
	var transaction Transaction
	var transactionTimeStr string // Temporary variable to hold transaction_time as string
	err := db.QueryRow(query, trackingCode).Scan(
		&transaction.TransactionID,
		&transaction.SourceAccountID,
		&transaction.DestinationAccountID,
		&transaction.Amount,
		&transaction.TransferType,
		&transactionTimeStr, // Scan transaction_time as string
		&transaction.TrackingCode,
		&transaction.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Transaction{}, fmt.Errorf("Transaction with tracking code %s not found", trackingCode)
		}
		return Transaction{}, err
	}

	// Parse transactionTimeStr into transaction.TransactionTime
	parsedTime, err := time.Parse("2006-01-02 15:04:05", transactionTimeStr)
	if err != nil {
		return Transaction{}, err
	}
	transaction.TransactionTime = parsedTime

	return transaction, nil
}
