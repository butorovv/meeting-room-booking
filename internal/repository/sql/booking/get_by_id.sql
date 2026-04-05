SELECT id, slot_id, user_id, status, conference_link, created_at, updated_at
FROM bookings
WHERE id = $1