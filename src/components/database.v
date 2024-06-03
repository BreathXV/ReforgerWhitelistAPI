module components

import db.sqlite
import json

@[table: 'clients']
pub struct Client {
	id				int			@[primary; sql: 'serial']
	server_id		string		
	identity_id		string		
}


pub const db := sqlite.connect('clients.db') or { panic(err) }

// Creates the database as soon as the module is called.
fn init () {
	sql db {
		create table Client
	} or { panic('${err} at ${@LOCATION}') }
}

// Inserts development data into the database for testing.
// ### ...
// ## Returns
	// - bool
		// > Whether the data was inserted successfully (true) or not (false).
pub fn dev_database() (bool) {
	dev_data := Client{
        server_id: "test",
        identity_id: "test"
	}

    sql db {
		insert dev_data into Client
	} or {
		panic(err)
		return false
	}
	return true
}


// Checks if a client/player is whitelisted.
// ### ...
// ## Parameters
    // - serverId : string
		// > The ID of the game server.
    // - identityId : string
	    // > The identity ID of the player.
// ## Returns
	// - bool
	    // > Whether the player is whitelisted or not (true) or not (false).
pub fn is_whitelisted(serverId string, identityId string) (bool) {
	exists := sql db {
        select from Client where server_id == serverId && identity_id == identityId
    } or {
        panic(err)
    }

	if exists.len == 0 {
		return false
	} else {
		return true
	}
}


pub fn players_whitelisted(serverId string) ?[]Client {
    al_players := sql db {
        select from Client where server_id == serverId
    } or { panic(err) }

    return al_players
}
