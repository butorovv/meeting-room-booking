SELECT id, room_id, days_mask, start_time, end_time, created_at
FROM schedules WHERE room_id = $1
