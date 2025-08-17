package models

import (
	"time"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
)

type Category struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	ParentCategoryID *int       `json:"parent_category_id"`
	Parent           *Category  `json:"parent"`   // Not in database, but used in API responses
	Children         []Category `json:"children"` // Not in database, but used in API responses
	Songs            []Song     `json:"songs"`    // Not in database, but used in API responses
	SongbookID       *int       `json:"songbook_id"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	SongCount int `json:"song_count"` // Not in database, but used in API responses
}

func (n *Category) songbookUpdatedAt() error {
	db := sqlite.GetDBConnection()
	_, err := db.Exec("UPDATE songbooks SET updated_at = ? WHERE id = ?", time.Now(), n.SongbookID)
	return err
}

func (n *Category) GetAllCategories() ([]Category, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, name, parent_category_id, songbook_id, created_at, updated_at FROM categories WHERE parent_category_id IS NULL")
	if err != nil {
		return []Category{}, err
	}
	defer rows.Close()

	result := []Category{}
	for rows.Next() {
		elem := Category{}
		err := rows.Scan(&elem.ID, &elem.Name, &elem.ParentCategoryID, &elem.SongbookID, &elem.CreatedAt, &elem.UpdatedAt)
		if err != nil {
			continue
		}

		// Load children recursively
		children, err := elem.loadChildrenRecursive()
		if err == nil {
			elem.Children = children
		}

		result = append(result, elem)
	}

	return result, nil
}

func (n *Category) GetCategoriesBySongbookID(songbookID int) ([]Category, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, name, parent_category_id, songbook_id, created_at, updated_at FROM categories WHERE parent_category_id IS NULL AND songbook_id = ? ORDER BY name ASC", songbookID)
	if err != nil {
		return []Category{}, err
	}
	defer rows.Close()

	result := []Category{}
	for rows.Next() {
		elem := Category{}
		err := rows.Scan(&elem.ID, &elem.Name, &elem.ParentCategoryID, &elem.SongbookID, &elem.CreatedAt, &elem.UpdatedAt)
		if err != nil {
			continue
		}

		// Load children recursively
		children, err := elem.loadChildrenRecursive()
		if err == nil {
			elem.Children = children
		}

		elem.SongCount = 0

		count, err := db.Query("SELECT COUNT(*) FROM song_categories WHERE category_id = ?", elem.ID)
		if err == nil {
			defer count.Close()
			if count.Next() {
				count.Scan(&elem.SongCount)
			}
		}

		result = append(result, elem)
	}

	return result, nil
}

func (n *Category) GetCategoryById(id int) (Category, error) {
	db := sqlite.GetDBConnection()
	var result Category
	err := db.QueryRow("SELECT id, name, parent_category_id, songbook_id, created_at, updated_at FROM categories WHERE id = ?", id).Scan(
		&result.ID, &result.Name, &result.ParentCategoryID, &result.SongbookID, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return Category{}, err
	}

	// Load parent if it exists
	if result.ParentCategoryID != nil {
		parent, err := result.loadParent()
		if err == nil {
			result.Parent = parent
		}
	}

	// Load children recursively
	children, err := result.loadChildrenRecursive()
	if err == nil {
		result.Children = children
	}

	return result, nil
}

func (n *Category) GetChildren() ([]Category, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, name, parent_category_id, songbook_id, created_at, updated_at FROM categories WHERE parent_category_id = ?", n.ID)
	if err != nil {
		return []Category{}, err
	}
	defer rows.Close()

	result := []Category{}
	for rows.Next() {
		elem := Category{}
		err := rows.Scan(&elem.ID, &elem.Name, &elem.ParentCategoryID, &elem.SongbookID, &elem.CreatedAt, &elem.UpdatedAt)
		if err != nil {
			continue
		}
		result = append(result, elem)
	}

	return result, nil
}

func (n *Category) CreateCategory() error {
	db := sqlite.GetDBConnection()

	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	result, err := db.Exec("INSERT INTO categories (name, parent_category_id, songbook_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		n.Name, n.ParentCategoryID, n.SongbookID, n.CreatedAt, n.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	n.ID = int(id)

	// Update songbook updated_at field
	if err := n.songbookUpdatedAt(); err != nil {
		return err
	}

	return nil
}

func (n *Category) UpdateCategory() error {
	db := sqlite.GetDBConnection()
	n.UpdatedAt = time.Now()

	_, err := db.Exec("UPDATE categories SET name = ?, parent_category_id = ?, updated_at = ? WHERE id = ?",
		n.Name, n.ParentCategoryID, n.UpdatedAt, n.ID)
	if err != nil {
		return err
	}

	// Update songbook updated_at field
	if err := n.songbookUpdatedAt(); err != nil {
		return err
	}

	return nil
}

func (n *Category) DeleteCategory() error {
	db := sqlite.GetDBConnection()

	// Delete related records first
	_, err := db.Exec("DELETE FROM song_categories WHERE category_id = ?", n.ID)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE categories SET parent_category_id = ? WHERE parent_category_id = ?", n.ParentCategoryID, n.ID)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM categories WHERE id = ?", n.ID)
	if err != nil {
		return err
	}

	// Update songbook updated_at field
	if err := n.songbookUpdatedAt(); err != nil {
		return err
	}

	return nil
}

// loadParent loads the parent category
func (n *Category) loadParent() (*Category, error) {
	if n.ParentCategoryID == nil {
		return nil, nil
	}

	db := sqlite.GetDBConnection()
	parent := &Category{}
	err := db.QueryRow("SELECT id, name, parent_category_id, songbook_id, created_at, updated_at FROM categories WHERE id = ?", *n.ParentCategoryID).Scan(
		&parent.ID, &parent.Name, &parent.ParentCategoryID, &parent.SongbookID, &parent.CreatedAt, &parent.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return parent, nil
}

// loadChildrenRecursive loads all children categories recursively
func (n *Category) loadChildrenRecursive() ([]Category, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, name, parent_category_id, songbook_id, created_at, updated_at FROM categories WHERE parent_category_id = ? ORDER BY name ASC", n.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []Category
	for rows.Next() {
		child := Category{}
		err := rows.Scan(&child.ID, &child.Name, &child.ParentCategoryID, &child.SongbookID, &child.CreatedAt, &child.UpdatedAt)
		if err != nil {
			continue
		}

		child.SongCount = 0

		count, err := db.Query("SELECT COUNT(*) FROM song_categories WHERE category_id = ?", child.ID)
		if err == nil {
			defer count.Close()
			if count.Next() {
				count.Scan(&child.SongCount)
			}
		}

		// Recursively load children of this child
		grandChildren, err := child.loadChildrenRecursive()
		if err == nil {
			child.Children = grandChildren
		}

		children = append(children, child)
	}

	return children, nil
}
