SELECT 
    s.id, s.room_id, s.start_time, s.end_time, s.created_at,
    CASE WHEN b.id IS NOT NULL AND b.status = 'active' THEN true ELSE false END as is_booked
FROM slots s
LEFT JOIN bookings b ON b.slot_id = s.id AND b.status = 'active'
WHERE s.room_id = $1 AND DATE(s.start_time) = $2::date
ORDER BY s.start_time
