DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
            CREATE TYPE user_role AS ENUM ('user', 'manager', 'warehouse_admin', 'admin');
        END IF;
    END
$$;

DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'custom_order_status') THEN
            CREATE TYPE custom_order_status AS ENUM (
                'CREATED',
                'APPROVED',
                'WAITING_PAYMENT',
                'PAID',
                'WAITING_DELIVERY',
                'READY_FOR_PICK_UP',
                'COMPLETED',
                'CANCELED'
                );
        END IF;
    END
$$;

DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'in_stock_order_status') THEN
            CREATE TYPE in_stock_order_status AS ENUM (
                'CREATED',
                'APPROVED',
                'WAITING_PAYMENT',
                'PAID',
                'READY_FOR_PICK_UP',
                'COMPLETED',
                'CANCELED'
                );
        END IF;
    END
$$;

DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'detail_type') THEN
            CREATE TYPE detail_type AS ENUM (
                'wheels',
                'transmission',
                'engine',
                'interior',
                'steering'
                );
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS users
(
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    role       user_role NOT NULL,
    created_at TIMESTAMPTZ      DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ      DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS car_models
(
    id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    price      NUMERIC(12, 0) NOT NULL CHECK (price >= 0),
    brand VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS cars
(
    id           UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    car_model_id UUID           NOT NULL REFERENCES car_models (id),
    year         INT            NOT NULL CHECK (year > 1900),
    price        NUMERIC(12, 0) NOT NULL CHECK (price >= 0),
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS custom_orders
(
    id         UUID PRIMARY KEY             DEFAULT gen_random_uuid(),
    manager_id UUID                NOT NULL REFERENCES users (id),
    client_id  UUID                NOT NULL REFERENCES users (id),
    car_id     UUID                NOT NULL REFERENCES cars (id),
    status     custom_order_status NOT NULL,
    price      NUMERIC(12, 0)      NOT NULL CHECK (price >= 0),
    created_at TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS in_stock_orders
(
    id         UUID PRIMARY KEY               DEFAULT gen_random_uuid(),
    manager_id UUID                  NOT NULL REFERENCES users (id),
    client_id  UUID                  NOT NULL REFERENCES users (id),
    car_id     UUID                  NOT NULL REFERENCES cars (id),
    status     in_stock_order_status NOT NULL,
    price      NUMERIC(12, 0)        NOT NULL CHECK (price >= 0),
    created_at TIMESTAMPTZ           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS test_drive_requests
(
    id             UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    client_id      UUID        NOT NULL REFERENCES users (id),
    car_id         UUID        NOT NULL REFERENCES cars (id),
    scheduled_time TIMESTAMPTZ NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS details
(
    id         UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    type       detail_type    NOT NULL,
    price      NUMERIC(12, 0) NOT NULL CHECK (price >= 0),
    created_at TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS part_compatibility
(
    detail_id      UUID NOT NULL REFERENCES details (id),
    car_model_id UUID NOT NULL REFERENCES car_models (id),
    PRIMARY KEY (detail_id, car_model_id)
);

CREATE TABLE IF NOT EXISTS car_parts
(
    car_id    UUID NOT NULL REFERENCES cars (id),
    detail_id UUID NOT NULL REFERENCES details (id),
    PRIMARY KEY (car_id, detail_id)
);


CREATE TABLE outbox (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    topic       VARCHAR(64) NOT NULL,
    payload     JSONB       NOT NULL,
    published   BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
