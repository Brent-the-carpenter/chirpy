

#!/usr/bin/env bash

# Function to print headers and pretty-print the JSON body
print_response() {
    local response="$1"
    
    # Extract headers (everything before JSON body)
    local headers=$(echo "$response" | awk '/^\s*{/ {exit} {print}')
    # Extract JSON body (everything from JSON body start onwards)
    local body=$(echo "$response" | awk '/^\s*{/{flag=1} flag' | jq .)
    
    # Print headers and JSON body
    echo "Headers:"
    echo "$headers"
    echo "Body:"
    echo "$body"
    echo "==========="
}

# Send the first POST request and capture the response
response=$(curl -i -s -X POST http://localhost:8080/admin/reset)

# Send the second POST request, capturing headers and JSON body separately
response2=$(curl -i -s -X POST -H "Content-Type: application/json" \
    -d '{"email":"mloneusk@example.co"}' \
    http://localhost:8080/api/users)

# Extract user_id from response2 JSON
userID1=$(echo "$response2" | awk '/^\s*{/{flag=1} flag' | jq -r '.id')

# Send the third POST request, capturing headers and JSON body separately
response3=$(curl -i -s -X POST -H "Content-Type: application/json" \
    -d '{"email":"dackjorsey@emaple.co"}' \
    http://localhost:8080/api/users)

# Extract user_id from response3 JSON
userID2=$(echo "$response3" | awk '/^\s*{/{flag=1} flag' | jq -r '.id')

# Use userID1 (or userID2) to create a new chirp
response4=$(curl -i -s -X POST -H "Content-Type: application/json" \
    -d "{
  \"body\": \"If you're committed enough, you can make any story work.\",
  \"user_id\": \"$userID1\"
}" \
    http://localhost:8080/api/chirps)

# Print responses with headers and JSON body for each request
echo "Response from /admin/reset:"
print_response "$response"

echo "Response from first /api/users request:"
print_response "$response2"

echo "Response from second /api/users request:"
print_response "$response3"

echo "Response from /api/chirps request:"
print_response "$response4"
