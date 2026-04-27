ALTER TABLE slots RENAME TO slots_deprecated;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'slots_pkey') THEN
        ALTER INDEX slots_pkey RENAME TO slots_deprecated_pkey;
    END IF;

    IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'idx_slots_room_date') THEN
        ALTER INDEX idx_slots_room_date RENAME TO idx_slots_deprecated_room_date;
    END IF;

    IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'idx_slots_room_start') THEN
        ALTER INDEX idx_slots_room_start RENAME TO idx_slots_deprecated_room_start;
    END IF;
END $$;

ALTER TABLE bookings ADD COLUMN IF NOT EXISTS room_id UUID;
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS start_time TIMESTAMPTZ;
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS end_time TIMESTAMPTZ;

UPDATE bookings
SET room_id = slots_deprecated.room_id,
    start_time = slots_deprecated.start_time,
    end_time = slots_deprecated.end_time
FROM slots_deprecated
WHERE bookings.slot_id = slots_deprecated.id;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM bookings WHERE room_id IS NULL OR start_time IS NULL OR end_time IS NULL) THEN
        RAISE EXCEPTION 'Cannot set NOT NULL: some bookings were not backfilled';
    END IF;
END $$;

ALTER TABLE bookings ALTER COLUMN room_id SET NOT NULL;
ALTER TABLE bookings ALTER COLUMN start_time SET NOT NULL;
ALTER TABLE bookings ALTER COLUMN end_time SET NOT NULL;
ALTER TABLE bookings ALTER COLUMN slot_id DROP NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_active_booking_room_time
    ON bookings(room_id, start_time)
    WHERE status = 'active';
