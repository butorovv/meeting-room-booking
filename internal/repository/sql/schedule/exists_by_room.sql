SELECT EXISTS(SELECT 1 FROM schedules WHERE room_id = $1)
