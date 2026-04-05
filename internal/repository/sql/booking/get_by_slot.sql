SELECT id, slot_id, user_id, status, conference_link, created_at, updated_at
FROM bookings WHERE slot_id = $1
