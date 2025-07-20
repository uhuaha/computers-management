CREATE TABLE computers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    mac_address TEXT NOT NULL UNIQUE,
    employee_abbreviation TEXT,
    description TEXT
);