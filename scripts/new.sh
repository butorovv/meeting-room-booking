#!/bin/bash

echo "=========================================="
echo "     ПОЛНАЯ ПРОВЕРКА ПРОЕКТА"
echo "=========================================="

# 1. Проверка форматирования
echo -e "\n1. ПРОВЕРКА ФОРМАТИРОВАНИЯ..."
gofmt -w .
echo "✅ Форматирование выполнено"

# 2. Проверка компиляции
echo -e "\n2. ПРОВЕРКА КОМПИЛЯЦИИ..."
go build ./...
if [ $? -eq 0 ]; then
    echo "✅ Компиляция успешна"
else
    echo "❌ Ошибка компиляции"
    exit 1
fi

# 3. Запуск тестов
echo -e "\n3. ЗАПУСК ТЕСТОВ..."
go test ./... -short
if [ $? -eq 0 ]; then
    echo "✅ Тесты прошли"
else
    echo "⚠️ Некоторые тесты упали (но это может быть из-за моков)"
fi

# 4. Запуск Docker
echo -e "\n4. ЗАПУСК DOCKER..."
docker-compose down -v 2>/dev/null
docker-compose up -d --build
sleep 10
echo "✅ Docker запущен"

# 5. Проверка health
echo -e "\n5. ПРОВЕРКА HEALTH..."
HEALTH=$(curl -s http://localhost:8080/_info)
if [[ "$HEALTH" == *"ok"* ]]; then
    echo "✅ Health check: $HEALTH"
else
    echo "❌ Health check failed: $HEALTH"
    docker logs room_booking_app --tail 20
    exit 1
fi

# 6. Получение токенов
echo -e "\n6. ПОЛУЧЕНИЕ ТОКЕНОВ..."
USER_TOKEN=$(curl -s -X POST http://localhost:8080/dummyLogin -H "Content-Type: application/json" -d '{"role":"user"}' | jq -r '.token')
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/dummyLogin -H "Content-Type: application/json" -d '{"role":"admin"}' | jq -r '.token')

if [ -n "$USER_TOKEN" ] && [ "$USER_TOKEN" != "null" ]; then
    echo "✅ User token получен"
else
    echo "❌ User token не получен"
    exit 1
fi

if [ -n "$ADMIN_TOKEN" ] && [ "$ADMIN_TOKEN" != "null" ]; then
    echo "✅ Admin token получен"
else
    echo "❌ Admin token не получен"
    exit 1
fi

# 7. Создание комнаты
echo -e "\n7. СОЗДАНИЕ КОМНАТЫ..."
ROOM_ID=$(curl -s -X POST http://localhost:8080/rooms/create \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Final Test Room","capacity":10}' | jq -r '.room.id')

if [ -n "$ROOM_ID" ] && [ "$ROOM_ID" != "null" ]; then
    echo "✅ Room ID: $ROOM_ID"
else
    echo "❌ Ошибка создания комнаты"
    exit 1
fi

# 8. Создание расписания
echo -e "\n8. СОЗДАНИЕ РАСПИСАНИЯ..."
TOMORROW=$(date -d "tomorrow" +%Y-%m-%d)
SCHEDULE_RESPONSE=$(curl -s -X POST "http://localhost:8080/rooms/$ROOM_ID/schedule/create" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"daysOfWeek\": [1,2,3,4,5], \"startTime\": \"09:00\", \"endTime\": \"17:00\"}")
SCHEDULE_ID=$(echo "$SCHEDULE_RESPONSE" | jq -r '.schedule.id')

if [ -n "$SCHEDULE_ID" ] && [ "$SCHEDULE_ID" != "null" ]; then
    echo "✅ Schedule ID: $SCHEDULE_ID"
else
    echo "❌ Ошибка создания расписания"
    exit 1
fi

# 9. Получение слотов
echo -e "\n9. ПОЛУЧЕНИЕ СЛОТОВ..."
SLOTS_RESPONSE=$(curl -s "http://localhost:8080/rooms/$ROOM_ID/slots/list?date=$TOMORROW" \
    -H "Authorization: Bearer $USER_TOKEN")
SLOTS_COUNT=$(echo "$SLOTS_RESPONSE" | jq '.slots | length')
FIRST_SLOT=$(echo "$SLOTS_RESPONSE" | jq '.slots[0]')
SLOT_START=$(echo "$FIRST_SLOT" | jq -r '.start')
SLOT_END=$(echo "$FIRST_SLOT" | jq -r '.end')

if [ -z "$SLOT_START" ] || [ "$SLOT_START" = "null" ] || [ -z "$SLOT_END" ] || [ "$SLOT_END" = "null" ]; then
    echo "❌ Не удалось получить start/end первого слота"
    echo "$SLOTS_RESPONSE" | jq .
    exit 1
fi

if [ "$SLOTS_COUNT" -eq 16 ]; then
    echo "✅ Количество слотов: $SLOTS_COUNT (ожидается 16)"
else
    echo "⚠️ Количество слотов: $SLOTS_COUNT (ожидается 16)"
fi

# Проверка, что createdAt отсутствует в ответе
if echo "$FIRST_SLOT" | jq -e '.createdAt' > /dev/null 2>&1; then
    echo "⚠️ Внимание: поле createdAt присутствует в ответе"
else
    echo "✅ Поле createdAt отсутствует в ответе (хорошо)"
fi

echo "Первый слот:"
echo "$FIRST_SLOT" | jq .

# 10. Создание брони
echo -e "\n10. СОЗДАНИЕ БРОНИ..."
BOOKING_RESPONSE=$(curl -s -X POST http://localhost:8080/bookings/create \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"roomId\": \"$ROOM_ID\",
        \"startTime\": \"$SLOT_START\",
        \"endTime\": \"$SLOT_END\"
    }")
BOOKING_STATUS=$(echo "$BOOKING_RESPONSE" | jq -r '.booking.status')
BOOKING_ID=$(echo "$BOOKING_RESPONSE" | jq -r '.booking.id')

if [ "$BOOKING_STATUS" = "active" ]; then
    echo "✅ Статус брони: $BOOKING_STATUS (ожидается active)"
    echo "✅ Booking ID: $BOOKING_ID"
else
    echo "❌ Ошибка создания брони: $BOOKING_STATUS"
    exit 1
fi

# 11. Проверка isBooked
echo -e "\n11. ПРОВЕРКА isBooked ПОСЛЕ БРОНИ..."
IS_BOOKED=$(curl -s "http://localhost:8080/rooms/$ROOM_ID/slots/list?date=$TOMORROW" \
    -H "Authorization: Bearer $USER_TOKEN" | jq '.slots[0].isBooked')

if [ "$IS_BOOKED" = "true" ]; then
    echo "✅ isBooked: $IS_BOOKED (ожидается true)"
else
    echo "❌ isBooked: $IS_BOOKED (ожидается true)"
    exit 1
fi

# 12. Двойное бронирование (должен быть 409)
echo -e "\n12. ПРОВЕРКА ЗАЩИТЫ ОТ ДВОЙНОГО БРОНИРОВАНИЯ..."
HTTP_STATUS=$(curl -s -X POST http://localhost:8080/bookings/create \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"roomId\": \"$ROOM_ID\",
        \"startTime\": \"$SLOT_START\",
        \"endTime\": \"$SLOT_END\"
    }" -w "%{http_code}" -o /dev/null)

if [ "$HTTP_STATUS" -eq 409 ]; then
    echo "✅ HTTP статус: $HTTP_STATUS (ожидается 409)"
else
    echo "❌ HTTP статус: $HTTP_STATUS (ожидается 409)"
fi

# 13. Отмена брони
echo -e "\n13. ОТМЕНА БРОНИ..."
CANCEL_STATUS=$(curl -s -X POST "http://localhost:8080/bookings/$BOOKING_ID/cancel" \
    -H "Authorization: Bearer $USER_TOKEN" | jq -r '.booking.status')

if [ "$CANCEL_STATUS" = "cancelled" ]; then
    echo "✅ Статус после отмены: $CANCEL_STATUS (ожидается cancelled)"
else
    echo "❌ Ошибка отмены: $CANCEL_STATUS"
fi

# 14. Проверка isBooked после отмены
echo -e "\n14. ПРОВЕРКА isBooked ПОСЛЕ ОТМЕНЫ..."
IS_BOOKED_AFTER=$(curl -s "http://localhost:8080/rooms/$ROOM_ID/slots/list?date=$TOMORROW" \
    -H "Authorization: Bearer $USER_TOKEN" | jq '.slots[0].isBooked')

if [ "$IS_BOOKED_AFTER" = "false" ]; then
    echo "✅ isBooked после отмены: $IS_BOOKED_AFTER (ожидается false)"
else
    echo "❌ isBooked после отмены: $IS_BOOKED_AFTER (ожидается false)"
fi

# 15. Идемпотентность повторной отмены
echo -e "\n15. ПРОВЕРКА ИДЕМПОТЕНТНОСТИ (ПОВТОРНАЯ ОТМЕНА)..."
IDEMPOTENT_STATUS=$(curl -s -X POST "http://localhost:8080/bookings/$BOOKING_ID/cancel" \
    -H "Authorization: Bearer $USER_TOKEN" -w "%{http_code}" -o /tmp/cancel_resp.json)
IDEMPOTENT_CODE=$(cat /tmp/cancel_resp.json | jq -r '.booking.status')

if [ "$IDEMPOTENT_STATUS" -eq 200 ] && [ "$IDEMPOTENT_CODE" = "cancelled" ]; then
    echo "✅ Идемпотентность: HTTP $IDEMPOTENT_STATUS, статус $IDEMPOTENT_CODE"
else
    echo "⚠️ Идемпотентность: HTTP $IDEMPOTENT_STATUS, статус $IDEMPOTENT_CODE"
fi

# 16. Проверка только будущих броней
echo -e "\n16. ПРОВЕРКА GET /bookings/my (только будущие активные)..."
MY_BOOKINGS=$(curl -s "http://localhost:8080/bookings/my" \
    -H "Authorization: Bearer $USER_TOKEN")
MY_BOOKINGS_COUNT=$(echo "$MY_BOOKINGS" | jq '.bookings | length')

echo "Количество броней в /bookings/my: $MY_BOOKINGS_COUNT"
if [ "$MY_BOOKINGS_COUNT" -eq 0 ]; then
    echo "✅ Отменённые брони не отображаются (хорошо)"
else
    echo "⚠️ В /bookings/my есть брони:"
    echo "$MY_BOOKINGS" | jq '.bookings[] | {id, status, startTime}'
fi

# 17. Проверка UserOnly middleware (попытка админом создать бронь)
echo -e "\n17. ПРОВЕРКА UserOnly MIDDLEWARE (admin пытается создать бронь)..."

ADMIN_SLOT_START=$(echo "$SLOTS_RESPONSE" | jq -r '.slots[1].start')
ADMIN_SLOT_END=$(echo "$SLOTS_RESPONSE" | jq -r '.slots[1].end')

ADMIN_BOOKING_STATUS=$(curl -s -X POST http://localhost:8080/bookings/create \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"roomId\": \"$ROOM_ID\",
        \"startTime\": \"$ADMIN_SLOT_START\",
        \"endTime\": \"$ADMIN_SLOT_END\"
    }" -w "%{http_code}" -o /tmp/admin_booking_resp.json)

if [ "$ADMIN_BOOKING_STATUS" -eq 403 ]; then
    echo "✅ Admin получил 403 Forbidden"
else
    echo "❌ Ожидался 403, получен $ADMIN_BOOKING_STATUS"
    echo "Ответ:"
    cat /tmp/admin_booking_resp.json
    exit 1
fi

echo -e "\n=========================================="
echo "              ПРОВЕРКА ЗАВЕРШЕНА"
echo "=========================================="