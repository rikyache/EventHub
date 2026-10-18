CREATE TABLE balance_transactions (
            transaction_id UUID PRIMARY KEY,
            user_id UUID NOT NULL REFERENCES users(user_id),
            amount BIGINT NOT NULL,
            type VARCHAR(30) NOT NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

            CONSTRAINT balance_transactions_type_check
                CHECK (type IN (
                    'deposit',
                        'ticket_purchase',
                        'refund'
                ))
);