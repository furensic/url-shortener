CREATE TABLE IF NOT EXISTS "shortened_uri" (
    "id" SERIAL PRIMARY KEY,
    "origin_uri" TEXT NOT NULL
);