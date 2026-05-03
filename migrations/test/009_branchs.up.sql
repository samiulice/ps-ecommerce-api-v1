CREATE TABLE IF NOT EXISTS branches (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    country VARCHAR(100) NOT null DEFAULT '',
    city VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    mobile VARCHAR(50) NOT null DEFAULT '',
    telephone VARCHAR(50) NOT null DEFAULT '',
    email VARCHAR(255) NOT null DEFAULT '',
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- INSERT INTO branches (
--     name,
--     country,
--     city,
--     address,
--     mobile,
--     telephone,
--     email,
--     latitude,
--     longitude
-- ) VALUES
--     (
--         'Noor Super Mart Main Branch',
--         'Bangladesh',
--         'Feni',
--         'Noor Super Mart Main Outlet',
--         '01800000001',
--         '03310000001',
--         'main@noorsupermart.com',
--         23.1698248,
--         91.201172
--     ),
--     (
--         'Noor Super Mart Sub Branch',
--         'Bangladesh',
--         'Feni',
--         'Noor Super Mart Sub Outlet',
--         '01800000002',
--         '03310000002',
--         'sub@noorsupermart.com',
--         23.1710248,
--         91.203172
--     )
-- ON CONFLICT (name) DO NOTHING;