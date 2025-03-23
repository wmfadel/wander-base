-- Enable PostGIS extension (if not already enabled in an earlier migration)
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE destinations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    location GEOMETRY(POINT, 4326) NOT NULL, -- SRID 4326 for WGS 84 lat/long
    description TEXT
);