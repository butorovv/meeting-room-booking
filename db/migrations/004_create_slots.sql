CREATE TABLE IF NOT EXISTS slots (
    id UUID PRIMARY KEY,
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT valid_slot_duration CHECK (end_time = start_time + INTERVAL '30 minutes'),
    CONSTRAINT unique_room_time UNIQUE (room_id, start_time)
);

CREATE INDEX idx_slots_room_start ON slots(room_id, start_time);