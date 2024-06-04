#!/bin/bash

headers="Content-Type: application/json"
body=$(cat <<EOF
{
  "server_id": "1cdfa108-0ba6-45fc-9756-22e76304e8fa",
  "identity_id": "465c3a56-743b-4755-bad0-2c60c625a779",
  "player_id": 123,
  "player_name": "Kieran"
}
EOF
)

response=$(curl -s -X POST -H "$headers" -d "$body" http://localhost:8080/check-whitelist)

echo $response | jq