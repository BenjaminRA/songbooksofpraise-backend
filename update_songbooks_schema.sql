-- Script to update column names from createdAt/updatedAt to created_at/updated_at
-- in songbooks_of_praise.sqlite database

-- Since SQLite doesn't support ALTER COLUMN directly, we need to:
-- 1. Create new tables with updated column names
-- 2. Copy data from old tables to new tables
-- 3. Drop old tables
-- 4. Rename new tables to original names

BEGIN TRANSACTION;

-- Update users table
CREATE TABLE users_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    admin BOOLEAN DEFAULT FALSE,
    editor BOOLEAN DEFAULT FALSE,
    moderator BOOLEAN DEFAULT FALSE,
    verified BOOLEAN DEFAULT FALSE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

INSERT INTO users_new (id, first_name, last_name, email, password, admin, editor, moderator, verified, created_at, updated_at)
SELECT id, first_name, last_name, email, password, admin, editor, moderator, verified, createdAt, updatedAt FROM users;

DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

-- Update songbooks table
CREATE TABLE songbooks_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    in_verification BOOLEAN DEFAULT FALSE,
    rejected BOOLEAN DEFAULT FALSE,
    owner_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

INSERT INTO songbooks_new (id, title, verified, in_verification, rejected, owner_id, created_at, updated_at)
SELECT id, title, verified, in_verification, rejected, owner_id, createdAt, updatedAt FROM songbooks;

DROP TABLE songbooks;
ALTER TABLE songbooks_new RENAME TO songbooks;

-- Update songbook_editors table
CREATE TABLE songbook_editors_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    songbook_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (songbook_id) REFERENCES songbooks(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO songbook_editors_new (id, songbook_id, user_id, created_at, updated_at)
SELECT id, songbook_id, user_id, createdAt, updatedAt FROM songbook_editors;

DROP TABLE songbook_editors;
ALTER TABLE songbook_editors_new RENAME TO songbook_editors;

-- Update categories table
CREATE TABLE categories_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    parent_category_id INTEGER,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (parent_category_id) REFERENCES categories(id)
);

INSERT INTO categories_new (id, name, parent_category_id, created_at, updated_at)
SELECT id, name, parent_category_id, createdAt, updatedAt FROM categories;

DROP TABLE categories;
ALTER TABLE categories_new RENAME TO categories;

-- Update songs table
CREATE TABLE songs_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    lyrics TEXT,
    music_sheet TEXT,
    music TEXT,
    music_only TEXT,
    youtube_url TEXT,
    description TEXT,
    number INTEGER,
    voices_all TEXT,
    voices_soprano TEXT,
    voices_contralto TEXT,
    voices_tenor TEXT,
    voices_bass TEXT,
    transpose INTEGER,
    scroll_speed INTEGER,
    songbook_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (songbook_id) REFERENCES songbooks(id)
);

INSERT INTO songs_new (id, title, lyrics, music_sheet, music, music_only, youtube_url, description, number, voices_all, voices_soprano, voices_contralto, voices_tenor, voices_bass, transpose, scroll_speed, songbook_id, created_at, updated_at)
SELECT id, title, lyrics, music_sheet, music, music_only, youtube_url, description, number, voices_all, voices_soprano, voices_contralto, voices_tenor, voices_bass, transpose, scroll_speed, songbook_id, createdAt, updatedAt FROM songs;

DROP TABLE songs;
ALTER TABLE songs_new RENAME TO songs;

-- Update song_categories table
CREATE TABLE song_categories_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    song_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (song_id) REFERENCES songs(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

INSERT INTO song_categories_new (id, song_id, category_id, created_at, updated_at)
SELECT id, song_id, category_id, createdAt, updatedAt FROM song_categories;

DROP TABLE song_categories;
ALTER TABLE song_categories_new RENAME TO song_categories;

COMMIT;
