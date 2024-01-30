package main

import (
	"db-final-project/handler"
	"db-final-project/storage"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := storage.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storage.CreateTables(db)

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
