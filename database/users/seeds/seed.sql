-- clean everything out first
TRUNCATE TABLE users RESTART IDENTITY CASCADE;

-- seed users
INSERT INTO users (name, email)
VALUES ('Alice Johnson', 'alice@example.com'),
       ('Bob Smith', 'bob@example.com'),
       ('Carol White', 'carol@example.com'),
       ('David Brown', 'david@example.com');