CREATE TYPE assembly_order_status AS ENUM ('CREATED', 'ASSEMBLED', 'FAIL');
CREATE TYPE source_order_type AS ENUM ('IN_STOCK', 'CUSTOM');

CREATE TABLE assembly_orders (
    id              UUID                  PRIMARY KEY DEFAULT gen_random_uuid(),
    source_order_id UUID                  NOT NULL,
    order_type      source_order_type     NOT NULL,
    status          assembly_order_status NOT NULL DEFAULT 'CREATED',
    removed         BOOLEAN               NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ           NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE outbox (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    topic       VARCHAR(64) NOT NULL,
    payload     JSONB       NOT NULL,
    published   BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
