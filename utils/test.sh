#!/bin/sh

# Start the application in the background
./rwapi &
RWAPI_PID=$!

# Wait for the application to start
sleep 5

# Run the tests
response=$(curl -s -X POST http://localhost:8080/check-whitelist -H "Content-Type: application/json" -d '{"ServerID": "1cdfa108-0ba6-45fc-9756-22e76304e8fa", "IdentityID": "465c3a56-743b-4755-bad0-2c60c625a779"}')
echo "$response" | jq .
echo "$response" | jq -e '.whitelisted == "true"' > /dev/null

# Capture the result of the test
TEST_RESULT=$?

# Stop the application
kill $RWAPI_PID

# Exit with the test result
exit $TEST_RESULT
