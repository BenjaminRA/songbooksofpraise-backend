I have a database with the following structure:

```
CREATE TABLE `temas` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `tema` VARCHAR(255), `createdAt` DATETIME NOT NULL, `updatedAt` DATETIME NOT NULL);
CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE `parrafos` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `himno_id` INTEGER REFERENCES `himnos` (`id`), `coro` TINYINT(1), `parrafo` TEXT, `createdAt` DATETIME NOT NULL, `updatedAt` DATETIME NOT NULL, acordes TEXT);
CREATE TABLE `sub_temas` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `sub_tema` VARCHAR(255), `tema_id` INTEGER REFERENCES `temas` (`id`), `createdAt` DATETIME NOT NULL, `updatedAt` DATETIME NOT NULL);
CREATE TABLE `sub_tema_himnos` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `himno_id` [object Object] REFERENCES `himnos` (`id`), `sub_tema_id` [object Object] REFERENCES `sub_temas` (`id`), `createdAt` DATETIME NOT NULL, `updatedAt` DATETIME NOT NULL);
CREATE TABLE `tema_himnos` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `himno_id` [object Object] REFERENCES `himnos` (`id`), `tema_id` [object Object] REFERENCES `temas` (`id`), `createdAt` DATETIME NOT NULL, `updatedAt` DATETIME NOT NULL);
CREATE TABLE IF NOT EXISTS "himnos"
(
  id INTEGER
          primary key autoincrement,
  titulo VARCHAR(255),
  createdAt DATETIME not null,
  updatedAt DATETIME not null
, transpose INTEGER, scroll_speed INTEGER);
CREATE TABLE `coros` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `titulo` VARCHAR(255), `tono` VARCHAR(255), `coro` TEXT, `createdAt` DATETIME NOT NULL, `updatedAt` DATETIME NOT NULL);
CREATE TABLE visitas (himno_id INTEGER, date DATETIME DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY (himno_id) REFERENCES himnos(id));
```

I need to migrate the database schema to a new version. The migration should include the following changes:

- There should be a new table called `users` with the following columns:
  - `id` as the primary key.
  - `first_name` for the user's first name.
  - `last_name` for the user's last name.
  - `email` for the user's email.
  - `password` for the user's password (hashed).
  - `admin` as a boolean to indicate if the user is an admin.
  - `editor` as a boolean to indicate if the user is an editor.
  - `moderator` as a boolean to indicate if the user is a moderator.
  - `verified` as a boolean to indicate if the user has verified their email.
  - `createdAt` for the creation date.
  - `updatedAt` for the last update date.
- There should be a new table called `songbooks` with the following columns:
  - `id` as the primary key.
  - `title` for the title of the songbook.
  - `verified` as a boolean to indicate if the songbook has been verified.
  - `in_verification` as a boolean to indicate if the songbook is currently in verification.
  - `owner_id` as a foreign key that references the `users` table.
  - `createdAt` for the creation date.
  - `updatedAt` for the last update date.
- There should be a new table called `songbook_editors` with the following columns:
  - `id` as the primary key.
  - `songbook_id` as a foreign key that references the `songbooks` table.
  - `user_id` as a foreign key that references the `users` table.
  - `createdAt` for the creation date.
  - `updatedAt` for the last update date.
- `himnos` should be renamed to `songs`.
- `parrafos` should be a column in `songs` instead of a separate table. This column should be of type TEXT and should have the `parrafos` concatenated to be in one big string. When a song has multiple paragraphs, they should be separated by a newline character. Each of the `parrafos` should be concatenated in the order they were originally in the `parrafos` table. Before each `parrafo` there should be a new line indicating if it's a verse or a chorus, e.g., "Verse 1:", "Chorus:", etc.
- `coros` should be deleted.
- `temas` and `sub_temas` should be merged into one single table called `categories`. There should be a column in this new table called `parent_category_id` as a foreign_key that references the same `categories` table to allow for subcategories. When a song does not have a `parent_category_id`, it should be NULL and it means is a top-level category.
- `tema_himnos` and `sub_tema_himnos` should be merged into one table called `song_categories`. This table should have a foreign key to the `songs` table and a foreign key to the `categories` table.
- The `visitas` table should be deleted.
- `songs` should have the following additional columns:
  - `music_sheet` for the url of the music sheet.
  - `music` for the url of the music file.
  - `music_only` for the url of the music file without vocals.
  - `youtube_url` for the url of the YouTube video.
  - `description` for a text description of the song.
  - `number` for the song number, nullable.
  - `voices_all` for the url of the all voices music file.
  - `voices_soprano` for the url of the soprano voice music file.
  - `voices_contralto` for the url of the alto voice music file.
  - `voices_tenor` for the url of the tenor voice music file.
  - `voices_bass` for the url of the bass voice music file.
  - `songbook_id` as a foreign key that references the `songbooks` table.
- The `users` table should be populated with the following initial data:
  - An admin user with the following details:
    - `first_name`: "Benjamín"
    - `last_name`: "Rodríguez"
    - `email`: "benjamin.gra720@gmail.com"
    - `password`: "admin" (hashed)
    - `admin`: true
    - `editor`: true
    - `moderator`: true
    - `verified`: true
    - `createdAt`: current timestamp
    - `updatedAt`: current timestamp
- All the information in the current `himnos_coros.sqlite` database should be migrated to the new schema.
- There should two songbooks created:
  - "Himnos y Cánticos del Evangelio" with:
    - `verified` set to true.
    - `in_verification` set to false.
    - `owner_id` set to the admin user created above.
    - `createdAt` set to current timestamp.
    - `updatedAt` set to current timestamp.
  - "Coros" with:
    - `verified` set to true.
    - `in_verification` set to false.
    - `owner_id` set to the admin user created above.
    - `createdAt` set to current timestamp.
    - `updatedAt` set to current timestamp.
- The `himnos` with id <= 517 should be added to the songbook "Himnos y Cánticos del Evangelio". The rest should be added to the songbook "Coros".
- All the categories and songs should be migrated to the new schema, ensuring that the relationships between songs and categories are preserved.

Generate a python script that performs this migration. The script should connect to the `himnos_coros.sqlite` database, perform the necessary transformations, and create a new SQLite database file called `songbooks_of_praise.sqlite` with the updated schema and data. The script should handle any potential errors gracefully and log the progress of the migration.

Work inside the `migration` directory
