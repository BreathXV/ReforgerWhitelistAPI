package components

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
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
