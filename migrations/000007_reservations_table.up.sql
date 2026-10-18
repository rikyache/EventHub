CREATE TABLE reservations (
    reservation_id UUID PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(event_id),
    user_id UUID NOT NULL REFERENCES users(user_id),

    price BIGINT NOT NULL CHECK (price >= 0),

    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT reservations_status_check
        CHECK (status IN (
                          'pending',
                          'confirmed',
                          'expired',
                          'cancelled'
        ))
)