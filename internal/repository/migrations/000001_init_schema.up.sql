CREATE TYPE order_type AS ENUM ('BID', 'ASK');
CREATE TYPE order_status AS ENUM ('PENDING', 'PARTIAL', 'FILLED', 'EXPIRED');

CREATE TABLE brokers (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    broker_id UUID NOT NULL REFERENCES brokers(id),
    owner_doc VARCHAR(255) NOT NULL,
    type order_type NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    price NUMERIC(20, 8) NOT NULL,
    quantity INTEGER NOT NULL,
    remaining_quantity INTEGER NOT NULL,
    valid_until TIMESTAMP WITH TIME ZONE NOT NULL,
    status order_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE trades (
    id UUID PRIMARY KEY,
    buy_order_id UUID NOT NULL REFERENCES orders(id),
    sell_order_id UUID NOT NULL REFERENCES orders(id),
    symbol VARCHAR(20) NOT NULL,
    price NUMERIC(20, 8) NOT NULL,
    quantity INTEGER NOT NULL,
    executed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_trade_pair UNIQUE (buy_order_id, sell_order_id)
);

CREATE INDEX idx_orders_ask
ON orders (symbol, type, status, price ASC, created_at ASC)
WHERE type = 'ASK' AND status IN ('PENDING', 'PARTIAL');

CREATE INDEX idx_orders_bid
ON orders (symbol, type, status, price DESC, created_at ASC)
WHERE type = 'BID' AND status IN ('PENDING', 'PARTIAL');
