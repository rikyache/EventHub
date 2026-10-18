ALTER TABLE events
    DROP CONSTRAINT IF EXISTS events_capacity_check;

ALTER TABLE events
    DROP COLUMN IF EXISTS reserved_count;