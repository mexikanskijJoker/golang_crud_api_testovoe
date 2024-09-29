CREATE TABLE
    IF NOT EXISTS songs (
        id SERIAL PRIMARY KEY,
        "group" VARCHAR(255),
        song VARCHAR(255),
        releaseDate VARCHAR(255) DEFAULT '1970-01-01',
        link VARCHAR(255) DEFAULT 'LINK',
        "text" TEXT DEFAULT 'TEXT'
    );