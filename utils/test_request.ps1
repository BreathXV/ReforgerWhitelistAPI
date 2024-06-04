$headers = @{
  "Content-Type" = "application/json"
}

$body = @{
  server_id = "1cdfa108-0ba6-45fc-9756-22e76304e8fa"
  identity_id = "465c3a56-743b-4755-bad0-2c60c625a779"
  player_id = 123
  player_name = "Kieran"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/check-whitelist" -Method Post -Headers $headers -Body $body

$response | ConvertTo-Json