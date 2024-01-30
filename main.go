package main

import (
	"db-final-project/handler"
	"db-final-project/storage"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	// Enable CORS middleware
	e.Use(middleware.CORS())

	db, err := storage.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storage.CreateTables(db)

	// Define HTTP endpoints
	e.POST("/register", func(c echo.Context) error {
		return handler.RegisterNewUser(c, db)
	})

	e.POST("/login", func(c echo.Context) error {
		return handler.LoginUser(c, db)
	})

	e.POST("/create-account", func(c echo.Context) error {
		return handler.CreateNewAccount(c, db)
	})

	e.POST("/transfer-card-to-card", func(c echo.Context) error {
		return handler.TransferMoneyCardToCard(c, db)
	})

	e.POST("/transfer-satna", func(c echo.Context) error {
		return handler.TransferMoneySATNA(c, db)
	})

	e.POST("/transfer-paya", func(c echo.Context) error {
		return handler.TransferMoneyPAYA(c, db)
	})

	e.GET("/transactions/:accountIdentifier/:n", func(c echo.Context) error {
		return handler.RetrieveLastNTransactions(c, db)
	})

	e.GET("/verify-transaction/:trackingCode", func(c echo.Context) error {
		return handler.VerifyTransaction(c, db)
	})

	port := ":8080"
	fmt.Printf("Server is running on port %s...\n", port)
	e.Logger.Fatal(e.Start(port))
}
