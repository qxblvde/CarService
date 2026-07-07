CREATE TABLE IF NOT EXISTS cars
(
    id        UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    brand     VARCHAR(255)   NOT NULL,
    model     VARCHAR(255)   NOT NULL,
    year      INT            NOT NULL CHECK (year > 1900),
    price     NUMERIC(12, 0) NOT NULL CHECK (price >= 0),
    available BOOLEAN        NOT NULL DEFAULT TRUE
);

INSERT INTO cars (id, brand, model, year, price, available)
VALUES ('33333333-0000-0000-0000-000000000001', 'Toyota', 'Camry', 2022, 2700000, TRUE),
       ('33333333-0000-0000-0000-000000000002', 'BMW', 'X5', 2023, 5500000, TRUE),
       ('33333333-0000-0000-0000-000000000003', 'Mercedes', 'E-Class', 2021, 4200000, TRUE),
       ('33333333-0000-0000-0000-000000000004', 'Audi', 'A6', 2020, 3800000, FALSE)
ON CONFLICT (id) DO NOTHING;
