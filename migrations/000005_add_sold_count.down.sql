ALTER TABLE events
    DROP CONSTRAINT events_sold_count_check,
    DROP COLUMN sold_count;