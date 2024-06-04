package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

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

// WhitelistedUsersRequestPayload defines the structure of the incoming
// JSON payload for all identity Ids associated with a server Id.
func CallToHandler(mux *http.ServeMux, db *gorm.DB) {
	database := Database{db: db}
	mux.HandleFunc("/check-whitelist", database.CheckWhitelistHandler)
}

// checkWhitelistHandler handles the /check-whitelist endpoint.
func (m *Database) CheckWhitelistHandler(w http.ResponseWriter, r *http.Request) {
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

	log.Println("Received POST request to /check-whitelist from:", reqPayload.ServerID)

	whitelisted, err := IsWhitelisted(m.db, reqPayload.ServerID, reqPayload.IdentityID)
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
func IsWhitelisted(db *gorm.DB, serverId string, identityId string) (bool, error) {
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
