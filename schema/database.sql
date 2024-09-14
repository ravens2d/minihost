CREATE TABLE IF NOT EXISTS post_entries (
    uuid BLOB(16) PRIMARY KEY,
    author_uuid BLOB(16) NOT NULL,
    posts BLOB NOT NULL,
    created_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    uuid BLOB(16) PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash BLOB NOT NULL
);
