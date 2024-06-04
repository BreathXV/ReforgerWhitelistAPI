package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// WhitelistedUsers defines the structure of the rows queried when getting all whitelisted players.
type client struct {
	gorm.Model
	ID         uint
	ServerID   string
	IdentityID string
}

// UserWhitelistRequestPayload defines the structure of the incoming JSON payload.
type UserWhitelistRequestPayload struct {
	ServerID   string `json:"server_id"`
	IdentityID string `json:"identity_id"`
	PlayerID   int32  `json:"player_id"`
	PlayerName string `json:"player_name"`
}

// UserWhitelistResponsePayload defines the structure of the JSON response.
type UserWhitelistResponsePayload struct {
	Whitelisted bool `json:"whitelisted"`
}

// WhitelistedUsersRequestPayload defines the structure of the incoming JSON payload for all identity Ids associated with a server Id.
type WhitelistedUsersRequestPayload struct {
	ServerID string `json:"server_id"`
}

// WhitelistedUsersResponsePayload defines the structure of the outgoing JSON payload for all identity Ids associated with a server Id.
type WhitelistedUsersResponsePayload struct {
	ServerID   string `json:"server_id"`
	IdentityID string `json:"identity_id"`
}

type Database struct {
	db *gorm.DB
}

func main() {
	db, err := gorm.Open(sqlite.Open("components\\db\\dev_database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&client{})

	log.Println("Connected to database successfully.")

	devDatabase(*db)

	database := Database{db: db}

	mux := http.NewServeMux()

	mux.HandleFunc("/check-whitelist", database.checkWhitelistHandler)
	/*
	   TODO: Setup methods used by API
	   http.HandleFunc("/get-whitelisted-users", getWhitelistedIds)
	*/
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}

// func (m *Database) databaseConnect() {
// 	database_conn, err := gorm.Open(sqlite.Open("dev_database.db"), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	m.db = database_conn
// }

func devDatabase(db gorm.DB) {
	devIdentityId := "465c3a56-743b-4755-bad0-2c60c625a779"
	devServerId := "1cdfa108-0ba6-45fc-9756-22e76304e8fa"

	result := db.Create(&client{ServerID: devServerId, IdentityID: devIdentityId})
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	log.Println("Added dev user to database.")
}

// checkWhitelistHandler handles the /check-whitelist endpoint.
func (m *Database) checkWhitelistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var reqPayload UserWhitelistRequestPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqPayload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Println("Received POST request to /check-whitelist from: ", reqPayload.ServerID)

	whitelisted, err := isWhitelisted(m.db, reqPayload.ServerID, reqPayload.IdentityID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resPayload := UserWhitelistResponsePayload{Whitelisted: whitelisted}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resPayload)
}

// TODO: Add a route/func to get all whitelisted identity Ids for that server Id.
// ----------------------------------------------------------------------------------------------
// getWhitelistedIds collects and stores all identity Ids associated with the server Id provided.
// func getWhitelistedIds(serverId string) ([]WhitelistedRow, error) {
// 	rows, err := db.Query(`SELECT identity_id FROM users WHERE server_id = ?`, serverId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()
// }

// isWhitelisted checks if the identity ID is whitelisted under the server ID.
func isWhitelisted(db *gorm.DB, serverId string, identityId string) (bool, error) {
	var whitelisted client
	result := db.Where("server_id = ? AND identity_id = ?", serverId, identityId).First(&whitelisted)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
