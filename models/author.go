package models

// Author represents an author in the system
// This model is not currently backed by the SQLite database
type Author struct {
	Author string `json:"author"`
}

func (n *Author) GetAllAuthors() ([]Author, error) {
	// This functionality is not implemented for SQLite
	return []Author{}, nil
}

func AddAuthor(author string) (Author, error) {
	// This functionality is not implemented for SQLite
	return Author{Author: author}, nil
}
