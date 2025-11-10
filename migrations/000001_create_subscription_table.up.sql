CREATE TABLE IF NOT EXISTS subscription(
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    service_name text NOT NULL,
    price integer NOT NULL,
    user_id uuid NOT NULL,
    start_date date NOT NULL,
    end_date date,
    is_deleted boolean NOT NULL DEFAULT false
);