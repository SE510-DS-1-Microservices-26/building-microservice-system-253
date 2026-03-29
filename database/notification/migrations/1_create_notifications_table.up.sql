-- create the notifications table
CREATE TABLE notifications
(
    id            SERIAL PRIMARY KEY,
    event_id      UUID      NOT NULL UNIQUE,
    core_item_id  INT       NOT NULL,
    owner_user_id INT,
    summary       TEXT      NOT NULL,
    occurred_at   TIMESTAMP NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);