DO $$
BEGIN
    IF to_regclass('slots_deprecated') IS NOT NULL THEN
        ALTER TABLE slots_deprecated RENAME TO slots;

        IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'slots_deprecated_pkey') THEN
            ALTER INDEX slots_deprecated_pkey RENAME TO slots_pkey;
        END IF;

        IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'idx_slots_deprecated_room_date') THEN
            ALTER INDEX idx_slots_deprecated_room_date RENAME TO idx_slots_room_date;
        END IF;

        IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'idx_slots_deprecated_room_start') THEN
            ALTER INDEX idx_slots_deprecated_room_start RENAME TO idx_slots_room_start;
        END IF;
    END IF;
END $$;
