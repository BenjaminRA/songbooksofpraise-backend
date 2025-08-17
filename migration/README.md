# Database Migration Script

This script migrates the database schema from `himnos_coros.sqlite` to the new `songbooks_of_praise.sqlite` format.

## Overview

The migration script transforms the database structure as follows:

### Schema Changes

1. **New Tables Created:**

   - `users` - User accounts with roles (admin, editor, moderator)
   - `songbooks` - Collections of songs
   - `songbook_editors` - Many-to-many relationship for songbook editors
   - `categories` - Unified category system (merges `temas` and `sub_temas`)
   - `songs` - Renamed from `himnos` with additional fields
   - `song_categories` - Many-to-many relationship between songs and categories

2. **Tables Removed:**

   - `coros` - Data integrated into songs
   - `parrafos` - Data integrated into songs as lyrics
   - `temas` and `sub_temas` - Merged into `categories`
   - `tema_himnos` and `sub_tema_himnos` - Merged into `song_categories`
   - `visitas` - Removed as requested

3. **Data Transformations:**
   - `himnos` → `songs` with additional multimedia fields
   - `parrafos` → Concatenated lyrics in songs with proper formatting
   - `temas`/`sub_temas` → Hierarchical categories structure
   - Songs distributed into two songbooks based on ID (≤517 = "Himnos y Cánticos del Evangelio", >517 = "Coros")

### Initial Data

The script creates:

- Admin user: Benjamín Rodríguez (benjamin.gra720@gmail.com)
- Two songbooks: "Himnos y Cánticos del Evangelio" and "Coros"

## Prerequisites

- Python 3.6 or higher
- The source database `himnos_coros.sqlite` must exist in the parent directory

## Usage

1. Navigate to the migration directory:

   ```bash
   cd migration
   ```

2. Run the migration script:

   ```bash
   python migrate_database.py
   ```

   Or make it executable and run directly:

   ```bash
   chmod +x migrate_database.py
   ./migrate_database.py
   ```

## Output

- **Target Database:** `songbooks_of_praise.sqlite` (created in parent directory)
- **Log File:** `migration.log` (created in migration directory)
- **Console Output:** Progress information and verification results

## Script Features

- **Error Handling:** Graceful error handling with detailed logging
- **Progress Tracking:** Real-time progress updates
- **Data Validation:** Verification of migrated data
- **Backup Safety:** Removes existing target database before creating new one
- **Detailed Logging:** Both file and console logging

## Migration Process

1. **Database Connection:** Connects to source and creates target database
2. **Schema Creation:** Creates all new tables with proper foreign keys
3. **Initial Data:** Creates admin user and default songbooks
4. **Category Migration:** Migrates and merges tema/sub_tema data
5. **Song Migration:** Transforms himnos to songs with concatenated lyrics
6. **Relationship Migration:** Preserves song-category relationships
7. **Verification:** Validates the migration results

## Troubleshooting

### Common Issues

1. **Source database not found:**

   - Ensure `himnos_coros.sqlite` exists in the parent directory
   - Check file permissions

2. **Permission errors:**

   - Ensure write permissions in the target directory
   - Run with appropriate user privileges

3. **Migration fails partway:**
   - Check the `migration.log` file for detailed error information
   - The script will remove any partially created target database

### Log Analysis

The migration log contains:

- Timestamp for each operation
- Record counts for verification
- Warning messages for any data inconsistencies
- Error details if migration fails

## Database Schema Reference

### New Tables Structure

```sql
-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    admin BOOLEAN DEFAULT FALSE,
    editor BOOLEAN DEFAULT FALSE,
    moderator BOOLEAN DEFAULT FALSE,
    verified BOOLEAN DEFAULT FALSE,
    createdAt DATETIME NOT NULL,
    updatedAt DATETIME NOT NULL
);

-- Songbooks table
CREATE TABLE songbooks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    in_verification BOOLEAN DEFAULT FALSE,
    owner_id INTEGER NOT NULL,
    createdAt DATETIME NOT NULL,
    updatedAt DATETIME NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

-- Categories table (hierarchical)
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    parent_category_id INTEGER,
    createdAt DATETIME NOT NULL,
    updatedAt DATETIME NOT NULL,
    FOREIGN KEY (parent_category_id) REFERENCES categories(id)
);

-- Songs table (enhanced)
CREATE TABLE songs (
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
    createdAt DATETIME NOT NULL,
    updatedAt DATETIME NOT NULL,
    FOREIGN KEY (songbook_id) REFERENCES songbooks(id)
);
```

## Post-Migration

After successful migration:

1. Verify the `songbooks_of_praise.sqlite` database was created
2. Check the migration log for any warnings
3. Test the new database with your application
4. Consider backing up the original `himnos_coros.sqlite` file
