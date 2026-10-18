ALTER TABLE events
    DROP CONSTRAINT IF EXISTS events_ticket_capacity_check;

ALTER TABLE events
    DROP COLUMN IF EXISTS reserved_count;