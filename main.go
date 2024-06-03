package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// WhitelistedUsers defines the structure of the rows queried when getting all whitelisted players.
type WhitelistedRow struct {
	gorm.Model
	ID         int
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

var db *sql.DB

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database successfully.")

	if !createDatabase() {
		return
	}

	_devDatabase()

	http.HandleFunc("/check-whitelist", checkWhitelistHandler)
	/*
		TODO: Setup methods used by API
		http.HandleFunc("/get-whitelisted-users", getWhitelistedIds)
	*/
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func _devDatabase() {
	devIdentityId := "465c3a56-743b-4755-bad0-2c60c625a779"
	devServerId := "1cdfa108-0ba6-45fc-9756-22e76304e8fa"
	devQuery := `
	INSERT INTO users(server_id, identity_id) VALUES(?, ?)
	`
	_, err := db.Exec(devQuery, devServerId, devIdentityId)
	if err != nil {
		log.Fatal(err)
	}
}

// createDatabase creates the database table if it does not exist.
func createDatabase() bool {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		server_id TEXT,
		identity_id TEXT
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// checkWhitelistHandler handles the /check-whitelist endpoint.
func checkWhitelistHandler(w http.ResponseWriter, r *http.Request) {
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

	whitelisted, err := isWhitelisted(reqPayload.ServerID, reqPayload.IdentityID)
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
func isWhitelisted(serverId, identityId string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE server_id = ? AND identity_id = ? LIMIT 1)`
	err := db.QueryRow(query, serverId, identityId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
