module main

import vweb
import json
import components

// WhitelistedRow defines the structure of the rows queried when getting all whitelisted players.
struct WhitelistedRow {
        id          int
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
        identity_id ?[]string @[json: identity_id]
}

struct App {
        vweb.Context
}

fn main() {
        mut app := App{}
        result := components.dev_database()
        if !result {
                println('Failed to connect to database component. ${@LOCATION}')
                return
        }

        vweb.run[App](app, 8080)
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

@['/whitelisted-players'; post]
pub fn (mut app App) whitelisted_players() vweb.Result {
        req_payload := json.decode(WhitelistedUsersRequestPayload, app.req.data) or {
                return app.json('{"error": "Bad request"}')
        }

        whitelisted := components.players_whitelisted(req_payload.server_id) or {
                return app.json('{"error": "Failed to fetch whitelisted players"}')
        }

        mut identity_ids := []string{}
        for player in whitelisted {
                identity_ids << player.str()
        }

        res_payload := WhitelistedUsersResponsePayload{
                identity_id: identity_ids
        }

        return app.json(json.encode(res_payload))
}
