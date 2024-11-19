CREATE TABLE dl (
                    id SERIAL PRIMARY KEY,
                    code VARCHAR(64) NOT NULL UNIQUE,
                    title VARCHAR(64) NOT NULL UNIQUE,
                    version INT NOT NULL DEFAULT 1
);

CREATE TABLE sl (
                    id SERIAL PRIMARY KEY,
                    code VARCHAR(64) NOT NULL UNIQUE,
                    title VARCHAR(64) NOT NULL UNIQUE,
                    is_detail BOOLEAN NOT NULL,
                    version INT NOT NULL DEFAULT 1
);
