#!/bin/bash
set -e

BASE_URL="http://localhost:8080"
TMP=$(mktemp)
RESULT_DIR=$(mktemp -d)
trap "rm -f $TMP; rm -rf $RESULT_DIR" EXIT

echo "=========================================="
echo "       RACE CONDITION TEST"
echo "=========================================="

echo "1. Получаем токены..."
USER_TOKEN=$(curl -s -X POST "$BASE_URL/dummyLogin" \
  -H "Content-Type: application/json" \
  -d '{"role":"user"}' | jq -r '.token')

ADMIN_TOKEN=$(curl -s -X POST "$BASE_URL/dummyLogin" \
  -H "Content-Type: application/json" \
  -d '{"role":"admin"}' | jq -r '.token')

if [ -z "$USER_TOKEN" ] || [ "$USER_TOKEN" = "null" ]; then
  echo "❌ Не удалось получить user token"
  exit 1
fi

if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" = "null" ]; then
  echo "❌ Не удалось получить admin token"
  exit 1
fi

echo "✅ Токены получены"

echo "2. Создаём комнату..."
ROOM_ID=$(curl -s -X POST "$BASE_URL/rooms/create" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Race Test Room","capacity":10}' | jq -r '.room.id')

if [ -z "$ROOM_ID" ] || [ "$ROOM_ID" = "null" ]; then
  echo "❌ Не удалось создать комнату"
  exit 1
fi

echo "✅ Room ID: $ROOM_ID"

echo "3. Создаём расписание..."
curl -s -X POST "$BASE_URL/rooms/$ROOM_ID/schedule/create" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"daysOfWeek":[1,2,3,4,5,6,7],"startTime":"09:00","endTime":"17:00"}' > /tmp/race_schedule.json

SCHEDULE_ID=$(jq -r '.schedule.id' /tmp/race_schedule.json)

if [ -z "$SCHEDULE_ID" ] || [ "$SCHEDULE_ID" = "null" ]; then
  echo "❌ Не удалось создать расписание"
  cat /tmp/race_schedule.json
  exit 1
fi

echo "✅ Schedule ID: $SCHEDULE_ID"

echo "4. Получаем первый доступный слот..."
TOMORROW=$(date -d "tomorrow" +%Y-%m-%d)

SLOTS_RESPONSE=$(curl -s "$BASE_URL/rooms/$ROOM_ID/slots/list?date=$TOMORROW" \
  -H "Authorization: Bearer $USER_TOKEN")

SLOT_START=$(echo "$SLOTS_RESPONSE" | jq -r '.slots[0].start')
SLOT_END=$(echo "$SLOTS_RESPONSE" | jq -r '.slots[0].end')

if [ -z "$SLOT_START" ] || [ "$SLOT_START" = "null" ]; then
  echo "❌ Не удалось получить слот"
  echo "$SLOTS_RESPONSE" | jq .
  exit 1
fi

echo "✅ Slot: $SLOT_START — $SLOT_END"

echo "5. Запускаем 50 параллельных бронирований одного слота..."

for i in $(seq 1 50); do
  (
    curl -s -o "$RESULT_DIR/body_$i.json" -w "%{http_code}\n" \
      -X POST "$BASE_URL/bookings/create" \
      -H "Authorization: Bearer $USER_TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"roomId\":\"$ROOM_ID\",
        \"startTime\":\"$SLOT_START\",
        \"endTime\":\"$SLOT_END\"
      }" > "$RESULT_DIR/code_$i.txt"
  ) &
done

wait

cat "$RESULT_DIR"/code_*.txt > "$TMP"

TOTAL=$(wc -l < "$TMP" | tr -d ' ')
SUCCESS=$(grep -cE '^(200|201)$' "$TMP" || true)
CONFLICT=$(grep -c '^409$' "$TMP" || true)
SERVER_ERRORS=$(grep -cE '^5[0-9][0-9]$' "$TMP" || true)
OTHER=$(grep -vE '^(200|201|409)$' "$TMP" | wc -l | tr -d ' ')

echo
echo "Результаты:"
echo "Successes:     $SUCCESS"
echo "Conflicts:     $CONFLICT"
echo "Server errors: $SERVER_ERRORS"
echo "Other:         $OTHER"
echo "Total:         $TOTAL"
echo
echo "HTTP codes:"
sort "$TMP" | uniq -c

if [ "$TOTAL" -eq 50 ] && [ "$SUCCESS" -eq 1 ] && [ "$CONFLICT" -eq 49 ] && [ "$SERVER_ERRORS" -eq 0 ] && [ "$OTHER" -eq 0 ]; then
  echo
  echo "✅ PASS: защита от двойного бронирования работает"
  exit 0
else
  echo
  echo "❌ FAIL: ожидалось 1 success, 49 conflicts, 0 server errors"
  exit 1
fi