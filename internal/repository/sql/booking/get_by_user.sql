SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at, b.updated_at
FROM bookings b
JOIN slots s ON b.slot_id = s.id
WHERE b.user_id = $1 AND s.start_time > NOW()
ORDER BY s.start_time