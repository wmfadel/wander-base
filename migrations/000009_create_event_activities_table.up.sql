CREATE TABLE IF NOT EXISTS event_activities (
    event_id INT NOT NULL,
    activity_id INT NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (activity_id) REFERENCES activities(id),
    PRIMARY KEY (event_id, activity_id)
);