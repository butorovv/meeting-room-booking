-- Добавляем тестовые комнаты
INSERT INTO rooms (id, name, description, capacity, created_at) VALUES
('11111111-1111-1111-1111-111111111111', 'Комната А', 'Большая переговорка', 10, NOW()),
('22222222-2222-2222-2222-222222222222', 'Комната Б', 'Малая переговорка', 4, NOW());

-- Добавляем расписание для комнаты А
INSERT INTO schedules (id, room_id, days_mask, start_time, end_time, created_at) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 31, '09:00', '18:00', NOW());
-- days_mask = 31 = пн-пт (1+2+4+8+16)