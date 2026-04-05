#!/bin/bash

set -e

BASE_URL="http://localhost:8080"

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Функция для вывода результатов
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
        exit 1
    fi
}

echo -e "${YELLOW}=== E2E Tests for Room Booking Service ===${NC}\n"

# 1. Получение токенов
echo -e "${YELLOW}1. Auth Tests${NC}"

ADMIN_TOKEN=$(curl -s -X POST $BASE_URL/dummyLogin -H "Content-Type: application/json" -d '{"role":"admin"}' | jq -r '.token')
if [ "$ADMIN_TOKEN" != "null" ] && [ -n "$ADMIN_TOKEN" ]; then
    print_result 0 "Admin token received"
else
    print_result 1 "Admin token failed"
fi

USER_TOKEN=$(curl -s -X POST $BASE_URL/dummyLogin -H "Content-Type: application/json" -d '{"role":"user"}' | jq -r '.token')
if [ "$USER_TOKEN" != "null" ] && [ -n "$USER_TOKEN" ]; then
    print_result 0 "User token received"
else
    print_result 1 "User token failed"
fi

# Негативные тесты Auth
echo -e "\n${YELLOW}2. Negative Auth Tests${NC}"

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/dummyLogin -H "Content-Type: application/json" -d '{"role":"superadmin"}')
if [ "$HTTP_CODE" -eq 400 ]; then
    print_result 0 "Invalid role returns 400"
else
    print_result 1 "Invalid role returns $HTTP_CODE (expected 400)"
fi

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/dummyLogin -H "Content-Type: text/plain" -d '{"role":"admin"}')
if [ "$HTTP_CODE" -eq 415 ]; then
    print_result 0 "Wrong Content-Type returns 415"
else
    print_result 1 "Wrong Content-Type returns $HTTP_CODE (expected 415)"
fi

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET $BASE_URL/dummyLogin)
if [ "$HTTP_CODE" -eq 405 ]; then
    print_result 0 "Wrong method returns 405"
else
    print_result 1 "Wrong method returns $HTTP_CODE (expected 405)"
fi

# 3. Rooms
echo -e "\n${YELLOW}3. Rooms Tests${NC}"

ROOM_RESPONSE=$(curl -s -X POST $BASE_URL/rooms/create \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"E2E Test Room","description":"Test Description","capacity":10}')
ROOM_ID=$(echo $ROOM_RESPONSE | jq -r '.room.id')

if [ "$ROOM_ID" != "null" ] && [ -n "$ROOM_ID" ]; then
    print_result 0 "Room created with ID: $ROOM_ID"
else
    print_result 1 "Room creation failed"
fi

# Проверка, что room обёрнут в объект
if echo $ROOM_RESPONSE | jq -e '.room' > /dev/null; then
    print_result 0 "Room response wrapped in 'room' object"
else
    print_result 1 "Room response missing 'room' wrapper"
fi

# Негативные тесты Rooms
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/rooms/create \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"capacity":10}')
if [ "$HTTP_CODE" -eq 400 ]; then
    print_result 0 "Create room without name returns 400"
else
    print_result 1 "Create room without name returns $HTTP_CODE (expected 400)"
fi

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/rooms/create \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Test"}')
if [ "$HTTP_CODE" -eq 403 ]; then
    print_result 0 "User cannot create room (403)"
else
    print_result 1 "User create room returns $HTTP_CODE (expected 403)"
fi

# Список комнат
ROOMS_LIST=$(curl -s -X GET $BASE_URL/rooms/list -H "Authorization: Bearer $ADMIN_TOKEN")
if echo $ROOMS_LIST | jq -e '.rooms' > /dev/null; then
    print_result 0 "Rooms list wrapped in 'rooms' object"
else
    print_result 1 "Rooms list missing 'rooms' wrapper"
fi

# 4. Schedules
echo -e "\n${YELLOW}4. Schedules Tests${NC}"

SCHEDULE_RESPONSE=$(curl -s -X POST $BASE_URL/rooms/$ROOM_ID/schedule/create \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"daysOfWeek":[1,2,3,4,5],"startTime":"09:00","endTime":"18:00"}')

# Проверяем, что вернулось (может быть 201 с обёрткой или 409)
if echo $SCHEDULE_RESPONSE | jq -e '.schedule' > /dev/null; then
    print_result 0 "Schedule created and wrapped in 'schedule' object"
elif echo $SCHEDULE_RESPONSE | jq -e '.error.code' | grep -q "SCHEDULE_EXISTS"; then
    print_result 0 "Schedule already exists (409) - OK"
else
    print_result 1 "Schedule creation failed: $SCHEDULE_RESPONSE"
fi

# Повторное создание расписания (должно вернуть 409)
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/rooms/$ROOM_ID/schedule/create \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"daysOfWeek":[1,2,3],"startTime":"10:00","endTime":"19:00"}')
if [ "$HTTP_CODE" -eq 409 ]; then
    print_result 0 "Duplicate schedule returns 409"
else
    print_result 1 "Duplicate schedule returns $HTTP_CODE (expected 409)"
fi

# 5. Slots
echo -e "\n${YELLOW}5. Slots Tests${NC}"

DATE=$(date -d "next monday" +%Y-%m-%d 2>/dev/null || echo "2026-04-07")
SLOTS_RESPONSE=$(curl -s -X GET "$BASE_URL/rooms/$ROOM_ID/slots/list?date=$DATE" \
    -H "Authorization: Bearer $ADMIN_TOKEN")

SLOT_ID=$(echo $SLOTS_RESPONSE | jq -r '.slots[0].id')
if [ "$SLOT_ID" != "null" ] && [ -n "$SLOT_ID" ]; then
    print_result 0 "Slots received, first slot ID: $SLOT_ID"
else
    print_result 1 "No slots found for date $DATE"
fi

# Без даты
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/rooms/$ROOM_ID/slots/list" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
if [ "$HTTP_CODE" -eq 400 ]; then
    print_result 0 "Missing date returns 400"
else
    print_result 1 "Missing date returns $HTTP_CODE (expected 400)"
fi

# 6. Bookings
echo -e "\n${YELLOW}6. Bookings Tests${NC}"

# Создание брони
BOOKING_RESPONSE=$(curl -s -X POST $BASE_URL/bookings/create \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"slotId\":\"$SLOT_ID\"}")

BOOKING_ID=$(echo $BOOKING_RESPONSE | jq -r '.booking.id')
if [ "$BOOKING_ID" != "null" ] && [ -n "$BOOKING_ID" ]; then
    print_result 0 "Booking created with ID: $BOOKING_ID"
else
    print_result 1 "Booking creation failed"
fi

# Проверка wrapper
if echo $BOOKING_RESPONSE | jq -e '.booking' > /dev/null; then
    print_result 0 "Booking response wrapped in 'booking' object"
else
    print_result 1 "Booking response missing 'booking' wrapper"
fi

# Повторная бронь (должна быть 409)
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/bookings/create \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"slotId\":\"$SLOT_ID\"}")
if [ "$HTTP_CODE" -eq 409 ]; then
    print_result 0 "Duplicate booking returns 409"
else
    print_result 1 "Duplicate booking returns $HTTP_CODE (expected 409)"
fi

# Список броней пользователя
MY_BOOKINGS=$(curl -s -X GET $BASE_URL/bookings/my -H "Authorization: Bearer $USER_TOKEN")
if echo $MY_BOOKINGS | jq -e '.bookings' > /dev/null; then
    print_result 0 "User bookings wrapped in 'bookings' object"
else
    print_result 1 "User bookings missing 'bookings' wrapper"
fi

# Отмена брони
CANCEL_RESPONSE=$(curl -s -X POST $BASE_URL/bookings/$BOOKING_ID/cancel \
    -H "Authorization: Bearer $USER_TOKEN")
CANCEL_STATUS=$(echo $CANCEL_RESPONSE | jq -r '.booking.status')
if [ "$CANCEL_STATUS" = "cancelled" ]; then
    print_result 0 "Booking cancelled successfully"
else
    print_result 1 "Booking cancellation failed"
fi

# Идемпотентность (повторная отмена)
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/bookings/$BOOKING_ID/cancel \
    -H "Authorization: Bearer $USER_TOKEN")
if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Idempotent cancel returns 200"
else
    print_result 1 "Idempotent cancel returns $HTTP_CODE (expected 200)"
fi

# 7. Admin: все брони
echo -e "\n${YELLOW}7. Admin Bookings List Tests${NC}"

ALL_BOOKINGS=$(curl -s -X GET "$BASE_URL/bookings/list?page=1&pageSize=10" \
    -H "Authorization: Bearer $ADMIN_TOKEN")

if echo $ALL_BOOKINGS | jq -e '.pagination' > /dev/null; then
    print_result 0 "Bookings list has pagination"
else
    print_result 1 "Bookings list missing pagination"
fi

# Проверка пагинации (pageSize > 100)
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/bookings/list?page=1&pageSize=200" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
if [ "$HTTP_CODE" -eq 400 ]; then
    print_result 0 "pageSize > 100 returns 400"
else
    print_result 1 "pageSize > 100 returns $HTTP_CODE (expected 400)"
fi

# User не может получить список всех броней
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/bookings/list" \
    -H "Authorization: Bearer $USER_TOKEN")
if [ "$HTTP_CODE" -eq 403 ]; then
    print_result 0 "User cannot access admin bookings list (403)"
else
    print_result 1 "User access admin bookings returns $HTTP_CODE (expected 403)"
fi

# 8. Health check
echo -e "\n${YELLOW}8. Health Check${NC}"

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/_info)
if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Health check returns 200"
else
    print_result 1 "Health check returns $HTTP_CODE (expected 200)"
fi

echo -e "\n${GREEN}=== All E2E Tests Passed! ===${NC}"