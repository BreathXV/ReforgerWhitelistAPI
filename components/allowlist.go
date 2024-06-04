package components

import (
	"gorm.io/gorm"
)

// WhitelistedUsers defines the structure of the rows queried when getting all whitelisted players.
type client struct {
	gorm.Model
	ID         uint
	ServerID   string
	IdentityID string
}

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
