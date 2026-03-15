-- migrations/init.sql

CREATE TABLE IF NOT EXISTS users (
    id        SERIAL PRIMARY KEY,
    name      VARCHAR(100) NOT NULL,
    email     VARCHAR(100) UNIQUE NOT NULL,
    gender    VARCHAR(10)  NOT NULL,
    birthdate DATE         NOT NULL
);

CREATE TABLE IF NOT EXISTS user_friends (
    user_id   INTEGER REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, friend_id),
    CHECK (user_id <> friend_id)
);

-- 20 users
INSERT INTO users (name, email, gender, birthdate) VALUES
('Alice Johnson',  'alice@mail.com',   'female', '1995-03-12'),
('Bob Smith',      'bob@mail.com',     'male',   '1993-07-22'),
('Carol White',    'carol@mail.com',   'female', '1998-01-05'),
('David Brown',    'david@mail.com',   'male',   '1990-11-30'),
('Eva Green',      'eva@mail.com',     'female', '1997-06-18'),
('Frank Miller',   'frank@mail.com',   'male',   '1992-09-14'),
('Grace Lee',      'grace@mail.com',   'female', '1999-04-25'),
('Henry Wilson',   'henry@mail.com',   'male',   '1994-02-08'),
('Iris Davis',     'iris@mail.com',    'female', '1996-08-17'),
('Jack Taylor',    'jack@mail.com',    'male',   '1991-12-03'),
('Karen Moore',    'karen@mail.com',   'female', '2000-05-20'),
('Leo Anderson',   'leo@mail.com',     'male',   '1993-10-11'),
('Mia Thomas',     'mia@mail.com',     'female', '1998-07-29'),
('Noah Jackson',   'noah@mail.com',    'male',   '1995-03-06'),
('Olivia Harris',  'olivia@mail.com',  'female', '1997-01-15'),
('Paul Martinez',  'paul@mail.com',    'male',   '1989-11-22'),
('Quinn Robinson', 'quinn@mail.com',   'female', '2001-09-09'),
('Ryan Clark',     'ryan@mail.com',    'male',   '1992-04-18'),
('Sophia Lewis',   'sophia@mail.com',  'female', '1996-06-27'),
('Tom Walker',     'tom@mail.com',     'male',   '1994-08-01');

INSERT INTO user_friends (user_id, friend_id) VALUES
(1, 3),(3, 1),
(1, 4),(4, 1),
(1, 5),(5, 1),
(2, 3),(3, 2),
(2, 4),(4, 2),
(2, 5),(5, 2),
(1, 2),(2, 1),
(6, 7),(7, 6),
(8, 9),(9, 8),
(10,11),(11,10),
(12,13),(13,12),
(14,15),(15,14);
