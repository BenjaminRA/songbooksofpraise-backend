package models

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BenjaminRA/himnario-backend/aws"
	"github.com/BenjaminRA/himnario-backend/db/sqlite"
	"github.com/BenjaminRA/himnario-backend/helpers"
	"github.com/oleiade/reflections"
)

type Song struct {
	ID              int        `json:"id"`
	Title           string     `json:"title"`
	Lyrics          string     `json:"lyrics"`
	MusicSheet      *string    `json:"music_sheet"`
	Music           *string    `json:"music"`
	MusicOnly       *string    `json:"music_only"`
	YouTubeURL      *string    `json:"youtube_url"`
	Description     *string    `json:"description"`
	Number          *int       `json:"number"`
	VoicesAll       *string    `json:"voices_all"`
	VoicesSoprano   *string    `json:"voices_soprano"`
	VoicesContralto *string    `json:"voices_contralto"`
	VoicesTenor     *string    `json:"voices_tenor"`
	VoicesBass      *string    `json:"voices_bass"`
	Transpose       *int       `json:"transpose"`
	ScrollSpeed     *int       `json:"scroll_speed"`
	SongbookID      int        `json:"songbook_id"`
	Categories      []Category `json:"categories"`    // Not in database, but used in API responses
	CategoriesID    []int      `json:"categories_id"` // Not in database, but used in API responses
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type SongCategory struct {
	ID         int       `json:"id"`
	SongID     int       `json:"song_id"`
	CategoryID int       `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (n *Song) songbookUpdatedAt() error {
	db := sqlite.GetDBConnection()
	_, err := db.Exec("UPDATE songbooks SET updated_at = ? WHERE id = ?", time.Now(), n.SongbookID)
	return err
}

func (n *Song) GetAllSongs() ([]Song, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, title, lyrics, music_sheet, music, music_only, youtube_url, description, number, voices_all, voices_soprano, voices_contralto, voices_tenor, voices_bass, transpose, scroll_speed, songbook_id, created_at, updated_at FROM songs")
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

func (n *Song) GetSongByID(id int) (Song, error) {
	db := sqlite.GetDBConnection()
	var result Song
	err := db.QueryRow("SELECT id, title, lyrics, music_sheet, music, music_only, youtube_url, description, number, voices_all, voices_soprano, voices_contralto, voices_tenor, voices_bass, transpose, scroll_speed, songbook_id, created_at, updated_at FROM songs WHERE id = ?", id).Scan(
		&result.ID, &result.Title, &result.Lyrics, &result.MusicSheet, &result.Music, &result.MusicOnly, &result.YouTubeURL, &result.Description, &result.Number, &result.VoicesAll, &result.VoicesSoprano, &result.VoicesContralto, &result.VoicesTenor, &result.VoicesBass, &result.Transpose, &result.ScrollSpeed, &result.SongbookID, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return Song{}, err
	}

	var categories []Category
	// Load categories for this song
	rows, err := db.Query("SELECT id, name, parent_category_id, songbook_id, created_at, updated_at FROM categories WHERE id IN (SELECT category_id FROM song_categories WHERE song_id = ?)", id)
	if err != nil {
		return Song{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.ParentCategoryID, &category.SongbookID, &category.CreatedAt, &category.UpdatedAt); err != nil {
			continue
		}
		categories = append(categories, category)
	}
	result.Categories = categories

	return result, nil
}

func (n *Song) createS3File(fieldName string, bucket string) error {
	actualFieldName, err := reflections.GetFieldNameByTagValue(n, "json", fieldName)
	if err != nil {
		return err
	}

	valueRaw, err := reflections.GetField(n, actualFieldName)
	if err != nil {
		return err
	}

	value, ok := valueRaw.(*string)
	if !ok {
		return fmt.Errorf("field %s is not a string pointer", actualFieldName)
	}

	if value != nil && *value == "__same__" {
		return nil
	}

	filename := helpers.GetFilenameFromPath(*value)
	key := strings.ReplaceAll(filename, "new", fmt.Sprintf("%d", n.ID))

	url, err := aws.S3UploadFile(*value, key, os.Getenv(bucket))
	if err != nil {
		fmt.Printf("Failed to upload file to S3: %v\n", err)
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Delete file after uploading to S3
	if err := os.Remove(*value); err != nil {
		fmt.Printf("Failed to delete local file after upload: %s\n", err.Error())
	}

	value = &url // Update the new song's field with the S3 URL

	reflections.SetField(n, actualFieldName, &url)

	_, err = sqlite.GetDBConnection().Exec("UPDATE songs SET "+fieldName+" = ?, updated_at = ? WHERE id = ?", url, time.Now(), n.ID)
	if err != nil {
		return fmt.Errorf("failed to update song in database: %w", err)
	}

	return nil
}

func (n *Song) CreateSong() error {
	db := sqlite.GetDBConnection()

	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	result, err := db.Exec("INSERT INTO songs (title, lyrics, youtube_url, description, number, transpose, scroll_speed, songbook_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		n.Title, n.Lyrics, n.YouTubeURL, n.Description, n.Number, n.Transpose, n.ScrollSpeed, n.SongbookID, n.CreatedAt, n.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	n.ID = int(id)

	n.createS3File("music_sheet", "AWS_S3_MUSIC_SHEET_BUCKET")
	n.createS3File("music", "AWS_S3_MUSIC_BUCKET")
	n.createS3File("music_only", "AWS_S3_MUSIC_ONLY_BUCKET")
	n.createS3File("voices_all", "AWS_S3_VOICES_BUCKET")
	n.createS3File("voices_soprano", "AWS_S3_VOICES_BUCKET")
	n.createS3File("voices_contralto", "AWS_S3_VOICES_BUCKET")
	n.createS3File("voices_tenor", "AWS_S3_VOICES_BUCKET")
	n.createS3File("voices_bass", "AWS_S3_VOICES_BUCKET")

	for _, categoryID := range n.CategoriesID {
		_, err := db.Exec("INSERT INTO song_categories (song_id, category_id, created_at, updated_at) VALUES (?, ?, ?, ?)",
			n.ID, categoryID, time.Now(), time.Now())
		if err != nil {
			return err
		}
	}

	// Update songbook updated_at field
	if err := n.songbookUpdatedAt(); err != nil {
		return err
	}

	return nil
}

func (n *Song) replaceFile(newSong *Song, fieldName string, bucket string) error {
	actualFieldName, err := reflections.GetFieldNameByTagValue(n, "json", fieldName)
	if err != nil {
		return err
	}

	originalValueRaw, err := reflections.GetField(n, actualFieldName)
	if err != nil {
		return err
	}

	originalValue, ok := originalValueRaw.(*string)
	if !ok {
		return fmt.Errorf("field %s is not a string pointer", actualFieldName)
	}

	newValueRaw, err := reflections.GetField(newSong, actualFieldName)
	if err != nil {
		return err
	}

	newValue, ok := newValueRaw.(*string)
	if !ok {
		return fmt.Errorf("field %s is not a string pointer", actualFieldName)
	}

	// Replacing music sheet file if it has changed
	if newValue != nil {
		// Do nothing if the value is "__same__"
		if *newValue == "__same__" {
			return nil
		}

		url, err := aws.S3UploadFile(*newValue, helpers.GetFilenameFromPath(*newValue), os.Getenv(bucket))
		if err != nil {
			fmt.Printf("Failed to upload file to S3: %v\n", err)
			return fmt.Errorf("failed to upload file to S3: %w", err)
		}

		// Delete file after uploading to S3
		if err := os.Remove(*newValue); err != nil {
			return fmt.Errorf("failed to delete local file after upload: %w", err)
		}

		newValue = &url // Update the new song's field with the S3 URL

		reflections.SetField(newSong, actualFieldName, &url)
	}

	// Delete old music sheet if it was changed
	if newValue == nil && originalValue != nil {
		filename := helpers.GetFilenameFromPath(*originalValue)
		if err := aws.S3DeleteFile(filename, os.Getenv(bucket)); err != nil {
			return err
		}
	}

	db := sqlite.GetDBConnection()

	_, err = db.Exec("UPDATE songs SET "+fieldName+" = ?, updated_at = ? WHERE id = ?", newValue, time.Now(), newSong.ID)
	if err != nil {
		return err
	}

	return nil
}

func (n *Song) UpdateSong() error {
	db := sqlite.GetDBConnection()
	n.UpdatedAt = time.Now()

	originalSong, err := n.GetSongByID(n.ID)
	if err != nil {
		return err
	}

	originalSong.replaceFile(n, "music_sheet", "AWS_S3_MUSIC_SHEET_BUCKET")
	originalSong.replaceFile(n, "music", "AWS_S3_MUSIC_BUCKET")
	originalSong.replaceFile(n, "music_only", "AWS_S3_MUSIC_ONLY_BUCKET")
	originalSong.replaceFile(n, "voices_all", "AWS_S3_VOICES_BUCKET")
	originalSong.replaceFile(n, "voices_soprano", "AWS_S3_VOICES_BUCKET")
	originalSong.replaceFile(n, "voices_contralto", "AWS_S3_VOICES_BUCKET")
	originalSong.replaceFile(n, "voices_tenor", "AWS_S3_VOICES_BUCKET")
	originalSong.replaceFile(n, "voices_bass", "AWS_S3_VOICES_BUCKET")

	db.Exec("DELETE FROM song_categories WHERE song_id = ?", n.ID)

	for _, categoryID := range n.CategoriesID {
		_, err := db.Exec("INSERT INTO song_categories (song_id, category_id, created_at, updated_at) VALUES (?, ?, ?, ?)",
			n.ID, categoryID, time.Now(), time.Now())
		if err != nil {
			return err
		}
	}

	_, err = db.Exec("UPDATE songs SET title = ?, lyrics = ?, youtube_url = ?, description = ?, number = ?, songbook_id = ?, transpose = ?, scroll_speed = ?, updated_at = ? WHERE id = ?",
		n.Title, n.Lyrics, n.YouTubeURL, n.Description, n.Number, n.SongbookID, n.Transpose, n.ScrollSpeed, n.UpdatedAt, n.ID)
	if err != nil {
		return err
	}

	// Update songbook updated_at field
	if err := n.songbookUpdatedAt(); err != nil {
		return err
	}

	return nil
}

func (n *Song) deleteS3Files() error {
	// Delete music sheet file if it exists
	if n.MusicSheet != nil {
		filename := helpers.GetFilenameFromPath(*n.MusicSheet)
		if err := aws.S3DeleteFile(filename, os.Getenv("AWS_S3_MUSIC_SHEET_BUCKET")); err != nil {
			return err
		}
	}

	// Delete music file if it exists
	if n.Music != nil {
		filename := helpers.GetFilenameFromPath(*n.Music)
		if err := aws.S3DeleteFile(filename, os.Getenv("AWS_S3_MUSIC_BUCKET")); err != nil {
			return err
		}
	}

	// Delete music only file if it exists
	if n.MusicOnly != nil {
		filename := helpers.GetFilenameFromPath(*n.MusicOnly)
		if err := aws.S3DeleteFile(filename, os.Getenv("AWS_S3_MUSIC_ONLY_BUCKET")); err != nil {
			return err
		}
	}

	// Delete voices files if they exist
	if n.VoicesAll != nil {
		filename := helpers.GetFilenameFromPath(*n.VoicesAll)
		if err := aws.S3DeleteFile(filename, os.Getenv("AWS_S3_VOICES_BUCKET")); err != nil {
			return err
		}
	}

	if n.VoicesSoprano != nil {
		filename := helpers.GetFilenameFromPath(*n.VoicesSoprano)
		if err := aws.S3DeleteFile(filename, os.Getenv("AWS_S3_VOICES_BUCKET")); err != nil {
			return err
		}
	}

	if n.VoicesContralto != nil {
		filename := helpers.GetFilenameFromPath(*n.VoicesContralto)
		if err := aws.S3DeleteFile(filename, os.Getenv("AWS_S3_VOICES_BUCKET")); err != nil {
			return err
		}
	}

	if n.VoicesTenor != nil {
		filename := helpers.GetFilenameFromPath(*n.VoicesTenor)
		if err := aws.S3DeleteFile(filename, os.Getenv("AWS_S3_VOICES_BUCKET")); err != nil {
			return err
		}
	}

	if n.VoicesBass != nil {
		filename := helpers.GetFilenameFromPath(*n.VoicesBass)
		if err := aws.S3DeleteFile(filename, os.Getenv("AWS_S3_VOICES_BUCKET")); err != nil {
			return err
		}
	}

	return nil
}

func (n *Song) DeleteSong() error {
	db := sqlite.GetDBConnection()

	err := n.deleteS3Files()
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM song_categories WHERE song_id = ?", n.ID)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM songs WHERE id = ?", n.ID)
	if err != nil {
		return err
	}

	// Update songbook updated_at field
	if err := n.songbookUpdatedAt(); err != nil {
		return err
	}

	return nil
}

func (n *SongCategory) AddSongToCategory() error {
	db := sqlite.GetDBConnection()

	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	result, err := db.Exec("INSERT INTO song_categories (song_id, category_id, created_at, updated_at) VALUES (?, ?, ?, ?)",
		n.SongID, n.CategoryID, n.CreatedAt, n.UpdatedAt)
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

func (n *SongCategory) RemoveSongFromCategory() error {
	db := sqlite.GetDBConnection()

	_, err := db.Exec("DELETE FROM song_categories WHERE song_id = ? AND category_id = ?", n.SongID, n.CategoryID)
	if err != nil {
		return err
	}

	return nil
}

func GetSongsByCategoryID(categoryID int) ([]Song, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query(`
		SELECT s.id, s.title, s.number, s.created_at, s.updated_at 
		FROM songs s 
		JOIN song_categories sc ON s.id = sc.song_id 
		WHERE sc.category_id = ?`, categoryID)
	if err != nil {
		return []Song{}, err
	}
	defer rows.Close()

	result := []Song{}
	for rows.Next() {
		elem := Song{}

		err := rows.Scan(&elem.ID, &elem.Title, &elem.Number, &elem.CreatedAt, &elem.UpdatedAt)
		if err != nil {
			continue
		}
		result = append(result, elem)
	}

	return result, nil
}

func GetSongsBySongbookID(songbookID int) ([]Song, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, title, number, created_at, updated_at FROM songs WHERE songbook_id = ?", songbookID)
	if err != nil {
		return []Song{}, err
	}
	defer rows.Close()

	result := []Song{}
	for rows.Next() {
		elem := Song{}
		err := rows.Scan(&elem.ID, &elem.Title, &elem.Number, &elem.CreatedAt, &elem.UpdatedAt)
		if err != nil {
			panic(err)
		}
		result = append(result, elem)
	}

	return result, nil
}
