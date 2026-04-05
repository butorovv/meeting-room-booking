#!/bin/bash

set -e

BASE_URL="http://localhost:8080"

echo "=== E2E Tests ==="

# 1. получаем токен admin
echo "1. Getting admin token..."
ADMIN_TOKEN=$(curl -s -X POST $BASE_URL/dummyLogin \
  -H "Content-Type: application/json" \
  -d '{"role":"admin"}' | jq -r '.token')
echo "Admin token: ${ADMIN_TOKEN:0:20}..."

# 2. получаем токен user
echo "2. Getting user token..."
USER_TOKEN=$(curl -s -X POST $BASE_URL/dummyLogin \
  -H "Content-Type: application/json" \
  -d '{"role":"user"}' | jq -r '.token')
echo "User token: ${USER_TOKEN:0:20}..."

# 3. создаем комнату (admin)
echo "3. Creating room..."
ROOM_RESPONSE=$(curl -s -X POST $BASE_URL/rooms/create \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"E2E Test Room","description":"Test","capacity":5}')
ROOM_ID=$(echo $ROOM_RESPONSE | jq -r '.id')
echo "Room ID: $ROOM_ID"

# 4. создаем расписание (admin)
echo "4. Creating schedule..."
curl -s -X POST $BASE_URL/rooms/$ROOM_ID/schedule/create \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"daysOfWeek":[1,2,3,4,5],"startTime":"09:00","endTime":"10:00"}' | jq '.'

# 5. получаем слоты
echo "5. Getting slots..."
DATE=$(date -d "tomorrow" +%Y-%m-%d)
SLOTS_RESPONSE=$(curl -s -X GET "$BASE_URL/rooms/$ROOM_ID/slots/list?date=$DATE" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
SLOT_ID=$(echo $SLOTS_RESPONSE | jq -r '.slots[0].id')
echo "Slot ID: $SLOT_ID"

# 6. создаем бронь (user)
echo "6. Creating booking..."
BOOKING_RESPONSE=$(curl -s -X POST $BASE_URL/bookings/create \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"slotId\":\"$SLOT_ID\"}")
BOOKING_ID=$(echo $BOOKING_RESPONSE | jq -r '.id')
echo "Booking ID: $BOOKING_ID"

# 7. проверяем список броней пользователя
echo "7. Getting user bookings..."
curl -s -X GET $BASE_URL/bookings/my \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.'

# 8. отменяем бронь
echo "8. Cancelling booking..."
curl -s -X POST $BASE_URL/bookings/$BOOKING_ID/cancel \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.'

# 9. проверяем список всех броней (admin)
echo "9. Getting all bookings (admin)..."
curl -s -X GET "$BASE_URL/bookings/list?page=1&pageSize=10" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'

echo "=== E2E Tests Completed ==="