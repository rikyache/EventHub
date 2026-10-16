-- 1. Создаём events

CREATE TABLE events (
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
                                            'concert', 'theatre', 'cinema', 'festival',
                                            'conference', 'sport', 'exhibition'
                                )),

                        CONSTRAINT events_status_check
                            CHECK (status IN (
                                              'draft', 'published', 'canceled', 'completed'
                                )),

                        CONSTRAINT events_date_check
                            CHECK (end_date > start_date)
);

CREATE INDEX idx_events_organizer_id
    ON events(organizer_id);


-- 2. Создаём tickets

CREATE TABLE tickets (
                         ticket_id UUID PRIMARY KEY,
                         event_id UUID NOT NULL,
                         user_id UUID NOT NULL,

                         price BIGINT NOT NULL,
                         status VARCHAR(20) NOT NULL DEFAULT 'active',

                         bought_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                         valid_until TIMESTAMPTZ NOT NULL,
                         used_at TIMESTAMPTZ,

                         CONSTRAINT tickets_event_fk
                             FOREIGN KEY (event_id) REFERENCES events(event_id),

                         CONSTRAINT tickets_user_fk
                             FOREIGN KEY (user_id) REFERENCES users(user_id),

                         CONSTRAINT tickets_price_check
                             CHECK (price >= 0),

                         CONSTRAINT tickets_status_check
                             CHECK (status IN (
                                               'active', 'used', 'cancelled', 'refunded'
                                 ))
);

CREATE INDEX idx_tickets_user_id ON tickets(user_id);
CREATE INDEX idx_tickets_event_id ON tickets(event_id);