package main

import (
	"database/sql"
	"db-final-project/handler"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Define database connection
const (
	dbDriver = "mysql"
	dbUser   = "bank"
	dbPass   = "bank"
	dbName   = "bank"
)

func main() {
	// Configure the database connection (always check and edit this part)
	const (
		username = "root"
		password = "" // replace with your password (if you've set one)
		hostname = "127.0.0.1:3306"
		dbname   = "your_db_name" // replace with your db name
	)

	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)

	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ping the database to verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected!")

	// Initialize the database tables (you can add this as a separate setup function)
	createTables(db)

	// Menu for user interactions
	for {
		fmt.Println("1. User login")
		fmt.Println("2. Register new user")
		fmt.Println("3. Create new account")
		fmt.Println("4. Money transfer via card to card")
		fmt.Println("5. Money transfer via SATNA")
		fmt.Println("6. Money transfer via PAYA")
		fmt.Println("7. Retrieve last n transactions of an account")
		fmt.Println("8. Verify transaction with tracking code")
		fmt.Println("9. Exit")
		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			// User login
			err := handler.LoginUser(db)
			if err != nil {
				log.Println("Login failed. Error:", err)
			} else {
				fmt.Println("Login successful!")
			}

		case 2:
			// Register new user
			err := handler.RegisterNewUser(db)
			if err != nil {
				log.Println("Error:", err)
			} else {
				fmt.Println("Registration successful!")
			}
		case 3:
			// Create new account
			err := handler.CreateNewAccount(db)
			if err != nil {
				log.Println("Error:", err)
			} else {
				fmt.Println("Account created successfully!")
			}

		case 4:
			// Money transfer via card to card
			fmt.Print("Enter source card number: ")
			var sourceCardNumber string
			_, err := fmt.Scanln(&sourceCardNumber)
			if err != nil {
				log.Println("Error:", err)
				break
			}

			fmt.Print("Enter destination card number: ")
			var destinationCardNumber string
			_, err = fmt.Scanln(&destinationCardNumber)
			if err != nil {
				log.Println("Error:", err)
				break
			}

			fmt.Print("Enter transfer amount: ")
			var amount float64
			_, err = fmt.Scanln(&amount)
			if err != nil {
				log.Println("Error:", err)
				break
			}

			err = handler.TransferMoneyCardToCard(db, sourceCardNumber, destinationCardNumber, amount)
			if err != nil {
				log.Println("Error:", err)
			} else {
				fmt.Println("Money transferred successfully!")
			}

		case 5:
			// Money transfer via SATNA
			fmt.Print("Enter source card number: ")
			var sourceCardNumber string
			_, err := fmt.Scanln(&sourceCardNumber)
			if err != nil {
				log.Println("Error:", err)
				break
			}

			fmt.Print("Enter destination SHEBA number: ")
			var destinationSHEBANumber string
			_, err = fmt.Scanln(&destinationSHEBANumber)
			if err != nil {
				log.Println("Error:", err)
				break
			}

			fmt.Print("Enter transfer amount: ")
			var amount float64
			_, err = fmt.Scanln(&amount)
			if err != nil {
				log.Println("Error:", err)
				break
			}

			err = handler.TransferMoneySATNA(db, sourceCardNumber, destinationSHEBANumber, amount)
			if err != nil {
				log.Println("Error:", err)
			} else {
				fmt.Println("Money transferred successfully!")
			}

		case 6:
			// Money transfer via PAYA
			fmt.Print("Enter source card number: ")
			var sourceCardNumber string
			_, err := fmt.Scanln(&sourceCardNumber)
			if err != nil {
				log.Println("Error:", err)
				break
			}

			fmt.Print("Enter destination SHEBA number: ")
			var destinationSHEBANumber string
			_, err = fmt.Scanln(&destinationSHEBANumber)
			if err != nil {
				log.Println("Error:", err)
				break
			}

			fmt.Print("Enter transfer amount: ")
			var amount float64
			_, err = fmt.Scanln(&amount)
			if err != nil {
				log.Println("Error:", err)
				break
			}

			err = handler.TransferMoneyPAYA(db, sourceCardNumber, destinationSHEBANumber, amount)
			if err != nil {
				log.Println("Error:", err)
			} else {
				fmt.Println("Money transferred successfully!")
			}

		case 7:
			// Retrieve last n transactions of an account
			fmt.Print("Enter account identifier (card number or SHEBA number): ")
			var accountIdentifier string
			_, err := fmt.Scanln(&accountIdentifier)
			if err != nil {
				log.Println("Error:", err)
				break
			}

			fmt.Print("Enter the number of transactions to retrieve: ")
			var n int
			_, err = fmt.Scanln(&n)
			if err != nil {
				log.Println("Error:", err)
				break
			}

			transactions, err := handler.RetrieveLastNTransactions(db, accountIdentifier, n)
			if err != nil {
				log.Println("Error:", err)
			} else {
				fmt.Println("Last", n, "transactions for account", accountIdentifier, ":")
				for _, transaction := range transactions {
					fmt.Printf("Transaction ID: %d\n", transaction.TransactionID)
					fmt.Printf("Source Account ID: %d\n", transaction.SourceAccountID)
					fmt.Printf("Destination Account ID: %d\n", transaction.DestinationAccountID)
					fmt.Printf("Amount: %.2f\n", transaction.Amount)
					fmt.Printf("Transfer Type: %s\n", transaction.TransferType)
					fmt.Printf("Transaction Time: %s\n", transaction.TransactionTime)
					fmt.Printf("Tracking Code: %s\n", transaction.TrackingCode)
					fmt.Printf("Status: %d\n", transaction.Status)
					fmt.Println("------------------------")
				}
			}

		case 8:
			// Verify transaction with tracking code
			fmt.Print("Enter tracking code: ")
			var trackingCode string
			_, err := fmt.Scanln(&trackingCode)
			if err != nil {
				log.Println("Error:", err)
				break
			}

			transaction, err := handler.VerifyTransaction(db, trackingCode)
			if err != nil {
				log.Println("Error:", err)
			} else {
				fmt.Println("Transaction details for tracking code", trackingCode, ":")
				fmt.Printf("Transaction ID: %d\n", transaction.TransactionID)
				fmt.Printf("Source Account ID: %d\n", transaction.SourceAccountID)
				fmt.Printf("Destination Account ID: %d\n", transaction.DestinationAccountID)
				fmt.Printf("Amount: %.2f\n", transaction.Amount)
				fmt.Printf("Transfer Type: %s\n", transaction.TransferType)
				fmt.Printf("Transaction Time: %s\n", transaction.TransactionTime)
				fmt.Printf("Tracking Code: %s\n", transaction.TrackingCode)
				fmt.Printf("Status: %d\n", transaction.Status)
			}

		case 9:
			// Exit the program
			fmt.Println("Exiting program...")
			return

		default:
			fmt.Println("Invalid choice. Please select a valid option.")
		}
	}
}

func createTables(db *sql.DB) {
	// Create tables if they do not exist
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
            user_id INT AUTO_INCREMENT PRIMARY KEY,
            first_name VARCHAR(50) NOT NULL,
            last_name VARCHAR(50) NOT NULL,
            username VARCHAR(50) UNIQUE NOT NULL,
            national_id VARCHAR(10) UNIQUE NOT NULL,
            password_hash VARCHAR(128) NOT NULL
        )`,

		`CREATE TABLE IF NOT EXISTS accounts (
            account_id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT,
            card_number VARCHAR(16) UNIQUE NOT NULL,
            sheba_number VARCHAR(24) NOT NULL,
            balance DECIMAL(12, 2) DEFAULT 0.00,
            FOREIGN KEY (user_id) REFERENCES users(user_id)
        )`,

		`CREATE TABLE IF NOT EXISTS transactions (
            transaction_id INT AUTO_INCREMENT PRIMARY KEY,
            source_account_id INT,
            destination_account_id INT,
            amount DECIMAL(12, 2) NOT NULL,
            transfer_type ENUM('CardToCard', 'SATNA', 'PAYA') NOT NULL,
            transaction_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            tracking_code VARCHAR(20) UNIQUE NOT NULL,
            status INT DEFAULT 0, -- 0 for failure, 1 for success
            FOREIGN KEY (source_account_id) REFERENCES accounts(account_id),
            FOREIGN KEY (destination_account_id) REFERENCES accounts(account_id)
        )`,
	}

	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			log.Fatal(err)
		}
	}
}
