package storage

import (
	"database/sql"
	"fmt"
	"log"
)

func ConnectToDatabase() (*sql.DB, error) {
	// Configure the database connection (always check and edit this part)
	const (
		username = "root"
		password = "" // replace with your password (if you've set one)
		hostname = "127.0.0.1:3306"
		dbname   = "your_db_name" // replace with your storage name
	)

	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)

	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database!")
	return db, nil
}

func CreateTables(db *sql.DB) {
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
