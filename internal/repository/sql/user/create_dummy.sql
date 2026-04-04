INSERT INTO users (id, email, role) VALUES ($1, $2, $3)
RETURNING id, email, role, created_at