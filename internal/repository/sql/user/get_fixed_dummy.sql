INSERT INTO users (id, email, role, created_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (id) DO UPDATE SET
    email = EXCLUDED.email,
    role = EXCLUDED.role
RETURNING id, email, role, created_at