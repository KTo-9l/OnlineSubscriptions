CREATE TABLE IF NOT EXISTS subscription(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    service_name text NOT NULL,
    price integer NOT NULL,
    user_id uuid NOT NULL,
    start_date date NOT NULL,
    end_date date,
    is_deleted boolean NOT NULL DEFAULT false,
    CONSTRAINT subscription_pkey PRIMARY KEY (id)
);