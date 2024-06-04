package main

import (
	"log"
	"net/http"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	cmp "github.com/BreathXV/ReforgerWhitelistAPI/handlers"
)

// WhitelistedUsers defines the structure of the rows queried when getting all whitelisted players.
type client struct {
	gorm.Model
	ID         uint
	ServerID   string
	IdentityID string
}

// main is the entry point of the application. It initializes the database connection,
// performs database migration and starts listening on port 8080 for incoming requests.
func main() {
	db, err := gorm.Open(sqlite.Open(os.Getenv(("DATABASE_PATH"))), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&client{})

	log.Println("Connected to database successfully.")

	devDatabase(*db)

	mux := http.NewServeMux()

	cmp.CallToHandler(mux, db)
	/*
	   TODO: Setup methods used by API
	   http.HandleFunc("/get-whitelisted-users", getWhitelistedIds)
	*/
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}

// devDatabase adds a dev user to the database for testing purposes.
func devDatabase(db gorm.DB) {
	devIdentityId := "465c3a56-743b-4755-bad0-2c60c625a779"
	devServerId := "1cdfa108-0ba6-45fc-9756-22e76304e8fa"

	result := db.Create(&client{ServerID: devServerId, IdentityID: devIdentityId})
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	log.Println("Added dev user to database.")
}
