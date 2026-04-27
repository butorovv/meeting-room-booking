SELECT id, slot_id, room_id, start_time, end_time, user_id, status, conference_link, created_at, updated_at
FROM bookings
WHERE id = $1
