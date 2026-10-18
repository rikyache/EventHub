ALTER TABLE events
    ADD COLUMN reserved_count INT NOT NULL DEFAULT 0;

ALTER TABLE events
    ADD CONSTRAINT events_capacity_check
        CHECK (
            sold_count >= 0
                AND reserved_count >= 0
                AND sold_count + reserved_count <= capacity
            );