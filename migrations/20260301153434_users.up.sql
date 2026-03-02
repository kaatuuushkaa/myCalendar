CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE users (
                       id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
                   username VARCHAR(30) UNIQUE NOT NULL,
                       password TEXT NOT NULL,
                       email VARCHAR(100) NOT NULL,
                    name VARCHAR(25) NOT NULL,
                    surname VARCHAR(55) NOT NULL,
                    birth DATE NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       deleted_at TIMESTAMP DEFAULT NULL
);