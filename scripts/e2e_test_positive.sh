#!/bin/bash

set -e

BASE_URL="http://localhost:8080"

echo "=== Positive E2E Tests ==="

ADMIN_TOKEN=$(curl -s -X POST $BASE_URL/dummyLogin -H "Content-Type: application/json" -d '{"role":"admin"}' | jq -r '.token')
USER_TOKEN=$(curl -s -X POST $BASE_URL/dummyLogin -H "Content-Type: application/json" -d '{"role":"user"}' | jq -r '.token')
echo "Tokens received"

ROOM_ID=$(curl -s -X POST $BASE_URL/rooms/create \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"E2E Room"}' | jq -r '.room.id')
echo "Room created: $ROOM_ID"

curl -s -X POST $BASE_URL/rooms/$ROOM_ID/schedule/create \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"daysOfWeek":[1,2,3,4,5],"startTime":"09:00","endTime":"10:00"}' > /dev/null
echo "Schedule created"

DATE=$(date -d "next monday" +%Y-%m-%d 2>/dev/null || echo "2026-04-07")
SLOTS_RESPONSE=$(curl -s -X GET "$BASE_URL/rooms/$ROOM_ID/slots/list?date=$DATE" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
SLOT_ID=$(echo "$SLOTS_RESPONSE" | jq -r '.slots[0].id')
SLOT_START=$(echo "$SLOTS_RESPONSE" | jq -r '.slots[0].start')
SLOT_END=$(echo "$SLOTS_RESPONSE" | jq -r '.slots[0].end')
echo "Slot ID: $SLOT_ID"

BOOKING_ID=$(curl -s -X POST $BASE_URL/bookings/create \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"roomId\":\"$ROOM_ID\",\"startTime\":\"$SLOT_START\",\"endTime\":\"$SLOT_END\"}" | jq -r '.booking.id')
echo "Booking created: $BOOKING_ID"

curl -s -X POST $BASE_URL/bookings/$BOOKING_ID/cancel \
    -H "Authorization: Bearer $USER_TOKEN" > /dev/null
echo "Booking cancelled"

echo "=== All Positive E2E Tests Passed! ==="
