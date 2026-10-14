CREATE TABLE users (
    user_id UUID PRIMARY KEY,

    username VARCHAR(30) NOT NULL,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,

    role VARCHAR(20) NOT NULL DEFAULT 'user',
    balance BIGINT NOT NULL DEFAULT 0,

    email VARCHAR(254) NOT NULL,
    password_hash TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_username_unique UNIQUE (username),
    CONSTRAINT users_email_unique UNIQUE (email),

    CONSTRAINT users_balance_check CHECK (balance >= 0),
    CONSTRAINT users_role_check
        CHECK (role IN ('user', 'organizer', 'admin'))
);
