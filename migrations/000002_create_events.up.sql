CREATE TABLE events(
    event_id UUID PRIMARY KEY,
    organizer_id UUID NOT NULL REFERENCES users(user_id),

    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(30) NOT NULL,
    location VARCHAR(255) NOT NULL,

    capacity INT NOT NULL CHECK (capacity > 0),
    price BIGINT NOT NULL CHECK (price >= 0),

    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    create_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    status VARCHAR(20) NOT NULL DEFAULT 'draft',

    CONSTRAINT events_type_check
        CHECK (type IN (
            'concert',
            'theatre',
            'cinema',
            'festival',
            'conference',
            'sport',
            'exhibition'
    )),

    CONSTRAINT events_status_check
        CHECK (status IN (
            'draft',
            'published',
            'canceled',
            'completed'
        )),

        CONSTRAINT events_date_check
            CHECK (end_date > start_date)
);

CREATE INDEX idx_events_organizer_id ON events(organizer_id);
