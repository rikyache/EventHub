ALTER TABLE events
ADD COLUMN sold_count INT NOT NULL DEFAULT 0,
    ADD CONSTRAINT events_sold_count_check
        CHECK (sold_count >= 0 AND sold_count <= capacity);