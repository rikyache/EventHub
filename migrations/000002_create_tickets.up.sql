CREATE TABLE IF NOT EXISTS tickets
(
    ticket_id   UUID PRIMARY KEY,
    event_id    UUID        NOT NULL,
    user_id     UUID        not null,

    price       BIGINT      NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',

    bought_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_until TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,

    CONSTRAINT tickets_event_fk
        FOREIGN KEY (event_id)
            REFERENCES events (event_id),

    CONSTRAINT tickets_event_fk
        FOREIGN KEY (user_id)
            REFERENCES users (user_id),

    CONSTRAINT tickets_price_check
        CHECK (price >= 0),

    CONSTRAINT tickets_status_check
        CHECK (status IN ('active', 'used', 'cancelled', 'refunded'))
);

    CREATE INDEX idx_tickets_user_id ON tickets(user_id);
    CREATE INDEX idx_tickets_event_id ON tickets(event_id);
