CREATE UNIQUE INDEX IF NOT EXISTS assembly_orders_source_order_id_key
    ON assembly_orders (source_order_id);
