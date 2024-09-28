CREATE TABLE
    IF NOT EXISTS songs (
        id SERIAL PRIMARY KEY,
        "group" VARCHAR(255),
        song VARCHAR(255)
    );