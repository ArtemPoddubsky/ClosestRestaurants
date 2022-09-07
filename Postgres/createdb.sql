CREATE DATABASE restaurants;
\c restaurants;

CREATE TABLE IF NOT EXISTS restaurants (
    Id INT PRIMARY KEY,
    Name TEXT NOT NULL,
    Address TEXT NOT NULL,
    Phone TEXT NOT NULL,
    Longitude DOUBLE PRECISION,
    Latitude DOUBLE PRECISION
);
