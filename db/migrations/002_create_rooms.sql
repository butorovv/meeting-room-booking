CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    capacity INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);