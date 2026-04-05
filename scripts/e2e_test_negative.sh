#!/bin/bash

set -e

BASE_URL="http://localhost:8080"

echo "=== Negative E2E Tests ==="

# Get tokens
ADMIN_TOKEN=$(curl -s -X POST $BASE_URL/dummyLogin -H "Content-Type: application/json" -d '{"role":"admin"}' | jq -r '.token')
USER_TOKEN=$(curl -s -X POST $BASE_URL/dummyLogin -H "Content-Type: application/json" -d '{"role":"user"}' | jq -r '.token')

# Test 1: Invalid role
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/dummyLogin -H "Content-Type: application/json" -d '{"role":"superadmin"}')
if [ "$HTTP_CODE" -eq 400 ]; then
    echo "✓ Invalid role -> 400"
else
    echo "✗ Invalid role -> $HTTP_CODE (expected 400)"
    exit 1
fi

# Test 2: User tries to create room
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/rooms/create \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Test"}')
if [ "$HTTP_CODE" -eq 403 ]; then
    echo "✓ User create room -> 403"
else
    echo "✗ User create room -> $HTTP_CODE (expected 403)"
    exit 1
fi

# Test 3: Create room without name
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/rooms/create \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"capacity":10}')
if [ "$HTTP_CODE" -eq 400 ]; then
    echo "✓ Create room without name -> 400"
else
    echo "✗ Create room without name -> $HTTP_CODE (expected 400)"
    exit 1
fi

# Test 4: No token
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET $BASE_URL/rooms/list)
if [ "$HTTP_CODE" -eq 401 ]; then
    echo "✓ No token -> 401"
else
    echo "✗ No token -> $HTTP_CODE (expected 401)"
    exit 1
fi

# Test 5: Wrong method
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET $BASE_URL/dummyLogin)
if [ "$HTTP_CODE" -eq 405 ]; then
    echo "✓ Wrong method -> 405"
else
    echo "✗ Wrong method -> $HTTP_CODE (expected 405)"
    exit 1
fi

echo "=== All Negative E2E Tests Passed! ==="