DROP INDEX IF EXISTS idx_unique_active_booking;
CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_unique_active_booking 
ON bookings(room_id, start_time, end_time) 
WHERE status = 'active';