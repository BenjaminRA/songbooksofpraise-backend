package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
	helpers "github.com/BenjaminRA/himnario-backend/helpers"
)

type Songbook struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Verified       bool      `json:"verified"`
	InVerification bool      `json:"in_verification"`
	Rejected       bool      `json:"rejected"`
	OwnerID        int       `json:"owner_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Editors    []*SongbookEditor `json:"editors"`    // Not in database, but used in API responses
	Categories []Category        `json:"categories"` // Not in database, but used in API responses
	SongCount  *int              `json:"song_count"` // Not in database, but used in API responses
}

type SongbookEditor struct {
	ID         int       `json:"id"`
	SongbookID int       `json:"songbook_id"`
	UserID     int       `json:"user_id"`
	User       User      `json:"user"` // Not in database, but used in API responses
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (n *Songbook) GetAllSongbooks() ([]Songbook, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, title, verified, in_verification, rejected, owner_id, created_at, updated_at FROM songbooks")
	if err != nil {
		return []Songbook{}, err
	}
	defer rows.Close()

	result := []Songbook{}
	for rows.Next() {
		elem := Songbook{}
		err := rows.Scan(&elem.ID, &elem.Title, &elem.Verified, &elem.InVerification, &elem.Rejected, &elem.OwnerID, &elem.CreatedAt, &elem.UpdatedAt)
		if err != nil {
			continue
		}

		// Populate editors for this songbook
		editors, err := GetSongbookEditors(elem.ID)
		if err == nil {
			elem.Editors = editors
		} else {
			elem.Editors = []*SongbookEditor{} // Empty slice if error
		}

		// Populate categories for this songbook
		categories, err := (&Category{}).GetCategoriesBySongbookID(elem.ID)
		if err == nil {
			elem.Categories = categories
		} else {
			elem.Categories = []Category{} // Empty slice if error
		}

		// Get SongCount
		count, err := db.Query("SELECT COUNT(*) FROM songs WHERE songbook_id = ?", elem.ID)
		if err == nil {
			defer count.Close()
			if count.Next() {
				count.Scan(&elem.SongCount)
			}
		} else {
			elem.SongCount = nil
		}

		result = append(result, elem)
	}

	return result, nil
}

func (n *Songbook) GetAllSongbooksApp() ([]Songbook, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, title, verified, in_verification, rejected, owner_id, created_at, updated_at FROM songbooks")
	if err != nil {
		return []Songbook{}, err
	}
	defer rows.Close()

	result := []Songbook{}
	for rows.Next() {
		elem := Songbook{}
		err := rows.Scan(&elem.ID, &elem.Title, &elem.Verified, &elem.InVerification, &elem.Rejected, &elem.OwnerID, &elem.CreatedAt, &elem.UpdatedAt)
		if err != nil {
			continue
		}

		// Get SongCount
		count, err := db.Query("SELECT COUNT(*) FROM songs WHERE songbook_id = ?", elem.ID)
		if err == nil {
			defer count.Close()
			if count.Next() {
				count.Scan(&elem.SongCount)
			}
		} else {
			elem.SongCount = nil
		}

		result = append(result, elem)
	}

	return result, nil
}

func (n *Songbook) GetSongs(id int) ([]Song, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, title, lyrics, music_sheet, music, music_only, youtube_url, description, number, voices_all, voices_soprano, voices_contralto, voices_tenor, voices_bass, transpose, scroll_speed, songbook_id, created_at, updated_at FROM songs WHERE songbook_id = ?", id)
	if err != nil {
		return []Song{}, err
	}
	defer rows.Close()

	result := []Song{}
	for rows.Next() {
		elem := Song{}
		err := rows.Scan(&elem.ID, &elem.Title, &elem.Lyrics, &elem.MusicSheet, &elem.Music, &elem.MusicOnly, &elem.YouTubeURL, &elem.Description, &elem.Number, &elem.VoicesAll, &elem.VoicesSoprano, &elem.VoicesContralto, &elem.VoicesTenor, &elem.VoicesBass, &elem.Transpose, &elem.ScrollSpeed, &elem.SongbookID, &elem.CreatedAt, &elem.UpdatedAt)
		if err != nil {
			continue
		}
		result = append(result, elem)
	}

	return result, nil
}

func (n *Songbook) GetSongbookByID(id int) (Songbook, error) {
	db := sqlite.GetDBConnection()
	var result Songbook
	err := db.QueryRow("SELECT id, title, verified, in_verification, rejected, owner_id, created_at, updated_at FROM songbooks WHERE id = ?", id).Scan(
		&result.ID, &result.Title, &result.Verified, &result.InVerification, &result.Rejected, &result.OwnerID, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return Songbook{}, err
	}

	// Populate editors for this songbook
	editors, err := GetSongbookEditors(result.ID)
	if err == nil {
		result.Editors = editors
	} else {
		result.Editors = []*SongbookEditor{} // Empty slice if error
	}

	// Populate categories for this songbook
	categories, err := (&Category{}).GetCategoriesBySongbookID(result.ID)
	if err == nil {
		result.Categories = categories
	} else {
		result.Categories = []Category{} // Empty slice if error
	}

	return result, nil
}

func (n *Songbook) CreateSongbook() error {
	db := sqlite.GetDBConnection()

	n.Verified = false
	n.InVerification = false
	n.Rejected = false
	n.Editors = []*SongbookEditor{} // Initialize empty editors slice
	n.Categories = []Category{}     // Initialize empty categories slice
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	result, err := db.Exec("INSERT INTO songbooks (title, verified, in_verification, rejected, owner_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		n.Title, n.Verified, n.InVerification, n.Rejected, n.OwnerID, n.CreatedAt, n.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	n.ID = int(id)

	return nil
}

func (n *Songbook) DeleteSongbook() error {
	db := sqlite.GetDBConnection()

	// Delete related records first
	_, err := db.Exec("DELETE FROM songs WHERE songbook_id = ?", n.ID)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM categories WHERE songbook_id = ?", n.ID)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM songbook_editors WHERE songbook_id = ?", n.ID)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM songbooks WHERE id = ?", n.ID)
	if err != nil {
		return err
	}

	return nil
}

func (n *Songbook) UpdateSongbook() error {
	db := sqlite.GetDBConnection()
	n.UpdatedAt = time.Now()

	_, err := db.Exec("UPDATE songbooks SET title = ?, updated_at = ? WHERE id = ?",
		n.Title, n.UpdatedAt, n.ID)
	if err != nil {
		return err
	}

	// Refresh editors to ensure the struct has current data
	err = n.RefreshEditors()
	if err != nil {
		return err
	}

	// Refresh categories to ensure the struct has current data
	err = n.RefreshCategories()
	if err != nil {
		return err
	}

	return nil
}

func (n *Songbook) AddEditor(editor string) error {
	db := sqlite.GetDBConnection()

	createdAt := time.Now()
	updatedAt := time.Now()

	user, err := (&User{}).GetUserByEmail(editor)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO songbook_editors (songbook_id, user_id, created_at, updated_at) VALUES (?, ?, ?, ?)",
		n.ID, user.ID, createdAt, updatedAt)
	if err != nil {
		return err
	}

	// Refresh the editors list
	editors, err := GetSongbookEditors(n.ID)
	if err == nil {
		n.Editors = editors
	}

	// Refresh categories to ensure consistency
	err = n.RefreshCategories()
	if err != nil {
		return err
	}

	return nil
}

func (n *Songbook) RemoveEditor(editor string) error {
	db := sqlite.GetDBConnection()

	user, err := (&User{}).GetUserByEmail(editor)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM songbook_editors WHERE songbook_id = ? AND user_id = ?", n.ID, user.ID)
	if err != nil {
		return err
	}

	// Refresh the editors list
	editors, err := GetSongbookEditors(n.ID)
	if err == nil {
		n.Editors = editors
	}

	// Refresh categories to ensure consistency
	err = n.RefreshCategories()
	if err != nil {
		return err
	}

	return nil
}

func (n *Songbook) RemoveAllEditors() error {
	db := sqlite.GetDBConnection()

	_, err := db.Exec("DELETE FROM songbook_editors WHERE songbook_id = ?", n.ID)
	if err != nil {
		return err
	}

	// Clear the editors list
	n.Editors = []*SongbookEditor{}

	return nil
}

// RefreshEditors updates the Editors field with current data from database
func (n *Songbook) RefreshEditors() error {
	editors, err := GetSongbookEditors(n.ID)
	if err != nil {
		return err
	}
	n.Editors = editors
	return nil
}

// RefreshCategories updates the Categories field with current data from database
func (n *Songbook) RefreshCategories() error {
	categories, err := (&Category{}).GetCategoriesBySongbookID(n.ID)
	if err != nil {
		return err
	}
	n.Categories = categories
	return nil
}

func GetSongbookEditors(songbookID int) ([]*SongbookEditor, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, songbook_id, user_id, created_at, updated_at FROM songbook_editors WHERE songbook_id = ?", songbookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*SongbookEditor{}
	for rows.Next() {
		elem := SongbookEditor{}
		err := rows.Scan(&elem.ID, &elem.SongbookID, &elem.UserID, &elem.CreatedAt, &elem.UpdatedAt)
		if err != nil {
			continue
		}

		// Fetch user details for each editor
		user, err := (&User{}).GetUserById(elem.UserID)
		if err != nil {
			continue // Skip this editor if user not found
		}
		elem.User = user

		result = append(result, &elem)
	}

	return result, nil
}

// SetSongbookVerificationStatus updates the verification status of a songbook
func SetSongbookVerificationStatus(id int, verified bool, inVerification bool, rejected bool, setTime bool) error {
	db := sqlite.GetDBConnection()

	query := "UPDATE songbooks SET verified = ?, in_verification = ?, rejected = ?"
	args := []interface{}{verified, inVerification, rejected}

	if setTime {
		query += ", updated_at = ?"
		args = append(args, time.Now())
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err := db.Exec(query, args...)
	return err
}

// Export a single songbook as SQL script
func (n *Songbook) ExportSongbookSQL() (string, error) {
	db := sqlite.GetDBConnection()
	var sql strings.Builder

	// Add header comment
	sql.WriteString("-- Songbook Export\n")
	sql.WriteString(fmt.Sprintf("-- Generated: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	sql.WriteString("-- \n\n")

	// Create table statements
	sql.WriteString("-- Create Tables\n")
	sql.WriteString(`CREATE TABLE IF NOT EXISTS songbooks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    parent_category_id INTEGER,
    songbook_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_category_id) REFERENCES categories(id) ON DELETE SET NULL,
    FOREIGN KEY (songbook_id) REFERENCES songbooks(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS songs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
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
    transpose INTEGER DEFAULT 0,
    scroll_speed REAL DEFAULT 1.0,
    songbook_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (songbook_id) REFERENCES songbooks(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS song_categories (
    song_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (song_id, category_id),
    FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

`)

	// Begin transaction
	sql.WriteString("BEGIN TRANSACTION;\n\n")

	// Export the songbook
	songbook, err := n.GetSongbookByID(n.ID)
	if err != nil {
		return "", err
	}

	sql.WriteString("-- Songbook\n")
	sql.WriteString(fmt.Sprintf("INSERT OR REPLACE INTO songbooks (id, title, created_at, updated_at) VALUES (%d, %s, %s, %s);\n\n",
		songbook.ID,
		helpers.SqlEscape(songbook.Title),
		helpers.SqlEscape(songbook.CreatedAt.Format("2006-01-02 15:04:05")),
		helpers.SqlEscape(songbook.UpdatedAt.Format("2006-01-02 15:04:05"))))

	// Export categories
	if len(songbook.Categories) > 0 {
		sql.WriteString("-- Categories\n")
		for _, category := range songbook.Categories {
			parentID := "NULL"
			if category.ParentCategoryID != nil {
				parentID = fmt.Sprintf("%d", *category.ParentCategoryID)
			}
			songbookID := "NULL"
			if category.SongbookID != nil {
				songbookID = fmt.Sprintf("%d", *category.SongbookID)
			}

			sql.WriteString(fmt.Sprintf("INSERT OR REPLACE INTO categories (id, name, parent_category_id, songbook_id, created_at, updated_at) VALUES (%d, %s, %s, %s, %s, %s);\n",
				category.ID,
				helpers.SqlEscape(category.Name),
				parentID,
				songbookID,
				helpers.SqlEscape(category.CreatedAt.Format("2006-01-02 15:04:05")),
				helpers.SqlEscape(category.UpdatedAt.Format("2006-01-02 15:04:05"))))
		}
		sql.WriteString("\n")
	}

	// Export songs (using existing GetSongs method)
	songs, err := n.GetSongs(n.ID)
	if err != nil {
		return "", err
	}

	if len(songs) > 0 {
		sql.WriteString("-- Songs\n")
		for _, song := range songs {
			sql.WriteString(fmt.Sprintf("INSERT OR REPLACE INTO songs (id, title, lyrics, music_sheet, music, music_only, youtube_url, description, number, voices_all, voices_soprano, voices_contralto, voices_tenor, voices_bass, transpose, scroll_speed, songbook_id, created_at, updated_at) VALUES (%d, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %d, %s, %s);\n",
				song.ID,
				helpers.SqlEscape(song.Title),
				helpers.SqlEscape(song.Lyrics),
				helpers.SqlEscapeNullString(song.MusicSheet),
				helpers.SqlEscapeNullString(song.Music),
				helpers.SqlEscapeNullString(song.MusicOnly),
				helpers.SqlEscapeNullString(song.YouTubeURL),
				helpers.SqlEscapeNullString(song.Description),
				helpers.SqlEscapeNullInt(song.Number),
				helpers.SqlEscapeNullString(song.VoicesAll),
				helpers.SqlEscapeNullString(song.VoicesSoprano),
				helpers.SqlEscapeNullString(song.VoicesContralto),
				helpers.SqlEscapeNullString(song.VoicesTenor),
				helpers.SqlEscapeNullString(song.VoicesBass),
				helpers.SqlEscapeNullInt(song.Transpose),
				helpers.SqlEscapeNullInt(song.ScrollSpeed),
				song.SongbookID,
				helpers.SqlEscape(song.CreatedAt.Format("2006-01-02 15:04:05")),
				helpers.SqlEscape(song.UpdatedAt.Format("2006-01-02 15:04:05"))))
		}
		sql.WriteString("\n")

		// Export song-category relationships
		sql.WriteString("-- Song-Category Relationships\n")
		for _, song := range songs {
			// Get categories for this song
			rows, err := db.Query("SELECT category_id FROM song_categories WHERE song_id = ?", song.ID)
			if err != nil {
				continue
			}

			for rows.Next() {
				var categoryID int
				if rows.Scan(&categoryID) == nil {
					sql.WriteString(fmt.Sprintf("INSERT OR REPLACE INTO song_categories (song_id, category_id) VALUES (%d, %d);\n",
						song.ID, categoryID))
				}
			}
			rows.Close()
		}
		sql.WriteString("\n")
	}

	// Commit transaction
	sql.WriteString("COMMIT;\n")

	return sql.String(), nil
}
