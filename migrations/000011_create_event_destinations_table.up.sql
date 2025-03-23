CREATE TABLE event_destinations (
    event_id INTEGER NOT NULL,
    destination_id INTEGER NOT NULL,
    datetime TIMESTAMP NOT NULL,
    PRIMARY KEY (event_id, destination_id),
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (destination_id) REFERENCES destinations(id) ON DELETE CASCADE
);