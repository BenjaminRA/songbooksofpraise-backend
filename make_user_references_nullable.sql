-- Make owner_id nullable to allow songbooks to persist after user deletion
-- This allows resources to remain in the system even when the creator is deleted

-- Step 1: Create new table with nullable owner_id and ON DELETE SET NULL
CREATE TABLE songbooks_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image TEXT,
    owner_id INTEGER,
    status VARCHAR(50) DEFAULT 'pending',
    rejected BOOLEAN DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Step 2: Copy data from old table
INSERT INTO songbooks_new SELECT * FROM songbooks;

-- Step 3: Drop old table
DROP TABLE songbooks;

-- Step 4: Rename new table
ALTER TABLE songbooks_new RENAME TO songbooks;

-- Step 5: Recreate indexes
CREATE INDEX IF NOT EXISTS idx_songbooks_owner ON songbooks(owner_id);
CREATE INDEX IF NOT EXISTS idx_songbooks_status ON songbooks(status);


-- Update songbook_editors to SET NULL on user deletion
-- Step 1: Create new table
CREATE TABLE songbook_editors_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    songbook_id INTEGER NOT NULL,
    user_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (songbook_id) REFERENCES songbooks(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Step 2: Copy data
INSERT INTO songbook_editors_new SELECT * FROM songbook_editors;

-- Step 3: Drop old table
DROP TABLE songbook_editors;

-- Step 4: Rename
ALTER TABLE songbook_editors_new RENAME TO songbook_editors;

-- Step 5: Recreate indexes
CREATE INDEX IF NOT EXISTS idx_songbook_editors_songbook ON songbook_editors(songbook_id);
CREATE INDEX IF NOT EXISTS idx_songbook_editors_user ON songbook_editors(user_id);
