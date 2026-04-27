SELECT b.id, b.slot_id, b.room_id, b.start_time, b.end_time, b.user_id, b.status, b.conference_link, b.created_at, b.updated_at
FROM bookings b
WHERE b.user_id = $1 AND b.start_time > NOW()
ORDER BY b.start_time
