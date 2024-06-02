module main

import vweb
import json
import components

// WhitelistedRow defines the structure of the rows queried when getting all whitelisted players.
struct WhitelistedRow {
	id         int
	server_id   string
	identity_id string
}

// UserWhitelistRequestPayload defines the structure of the incoming JSON payload.
struct UserWhitelistRequestPayload {
	server_id   string @[json: server_id]
	identity_id string @[json: identity_id]
	player_id   int    @[json: player_id]
	player_name string @[json: player_name]
}

// UserWhitelistResponsePayload defines the structure of the JSON response.
struct UserWhitelistResponsePayload {
	whitelisted bool @[json: whitelisted]
}

// WhitelistedUsersRequestPayload defines the structure of the incoming JSON payload for all identity Ids associated with a server Id.
struct WhitelistedUsersRequestPayload {
	server_id string @[json: server_id]
}

// WhitelistedUsersResponsePayload defines the structure of the outgoing JSON payload for all identity Ids associated with a server Id.
struct WhitelistedUsersResponsePayload {
	server_id   string @[json: server_id]
	identity_id string @[json: identity_id]
}

struct App {
	vweb.Context
}

fn main() {
	mut app := App{}
	result := components.dev_database()
	if !result {
		println("Failed to connect to database component. ${@LOCATION}")
		return
    }

	vweb.run<App>(app, 8080)
}

pub fn (mut app App) before_request() {
	// Handler for common tasks before each request
}

fn (mut app App) after_request() {
	// Handler for common tasks after each request
}

// Handler for /check-whitelist endpoint
@['/check-whitelist'; post]
pub fn (mut app App) check_whitelist() vweb.Result {
    req_payload := json.decode(UserWhitelistRequestPayload, app.req.data) or {
        return app.json('{"error": "Bad request"}')
    }

    whitelisted := components.is_whitelisted(req_payload.server_id, req_payload.identity_id)

	res_payload := UserWhitelistResponsePayload{
		whitelisted: whitelisted
	}

    return app.json(json.encode(res_payload))
}


// TODO: Add a route/func to get all whitelisted identity Ids for that server Id.
// ----------------------------------------------------------------------------------------------
// get_whitelisted_ids collects and stores all identity Ids associated with the server Id provided.
// pub fn (mut app App) get_whitelisted_ids(server_id string) ?[]WhitelistedRow {
// 	rows := db.exec('SELECT identity_id FROM users WHERE server_id = ?', server_id) or {
// 		return err
// 	}
// }

fn is_whitelisted(db sqlite.DB, server_id string, identity_id string) ?bool {
	query := 'SELECT EXISTS(SELECT 1 FROM users WHERE server_id = ? AND identity_id = ? LIMIT 1)'
	exists := db.q_int(query, server_id, identity_id) or {
		return err
	}
	return exists == 1
}
