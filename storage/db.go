package storage

import (
	"database/sql"
	"fmt"
	"log"
)

func ConnectToDatabase() (*sql.DB, error) {
	const (
		username = "root"
		password = "yourpassword"
		hostname = "127.0.0.1"
		dbname   = "your_db_name"
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3307)/%s", username, password, hostname, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database!")
	return db, nil
}

func CreateTables(db *sql.DB) {
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

		`CREATE TABLE IF NOT EXISTS balance_change_log (
 			log_id INT AUTO_INCREMENT PRIMARY KEY,
 			account_id INT NOT NULL,
 			old_balance DECIMAL(12, 2) NOT NULL,
 			new_balance DECIMAL(12, 2) NOT NULL,
 			change_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 			FOREIGN KEY (account_id) REFERENCES accounts(account_id)
			)`,

		`CREATE TRIGGER IF NOT EXISTS balance_change_trigger BEFORE UPDATE ON accounts
       FOR EACH ROW
       BEGIN
           INSERT INTO balance_change_log (account_id, old_balance, new_balance)
           VALUES (OLD.account_id, OLD.balance, NEW.balance);
       END;`,
	}

	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			log.Fatal(err)
		}
	}
}
