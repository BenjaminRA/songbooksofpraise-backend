package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Admin     bool      `json:"admin"`
	Editor    bool      `json:"editor"`
	Moderator bool      `json:"moderator"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CheckEmailTaken(email string) bool {
	// db := mongodb.GetMongoDBConnection()
	// match := db.Collection("Users").FindOne(context.TODO(), bson.M{
	// 	"email": email,
	// })

	db := sqlite.GetDBConnection()
	result := db.QueryRow("SELECT email FROM users WHERE email = ?", email)

	if result.Err() != nil {
		return false
	}

	var existingEmail string
	if err := result.Scan(&existingEmail); err != nil {
		if err == sql.ErrNoRows {
			return false // No user with this email exists
		}
		fmt.Println("Error scanning email:", err.Error())
		return true // Error occurred, assume email is taken
	}

	return existingEmail == email // Email exists in the database

}

func CheckEmailTakenWithId(email string, user_id int) bool {
	db := sqlite.GetDBConnection()
	result := db.QueryRow("SELECT email FROM users WHERE email = ? AND id != ?", email, user_id)

	if result.Err() != nil {
		return false
	}

	var existingEmail string
	if err := result.Scan(&existingEmail); err != nil {
		if err == sql.ErrNoRows {
			return false // No user with this email exists
		}
		fmt.Println("Error scanning email:", err)
		return true // Error occurred, assume email is taken
	}

	return existingEmail == email // Email exists in the database

}

func (n *User) GetUserById(user_id int) (User, error) {
	db := sqlite.GetDBConnection()
	var user User
	err := db.QueryRow("SELECT * FROM users WHERE id = ?", user_id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Admin, &user.Editor, &user.Moderator, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (n *User) GetUserByEmail(email string) (User, error) {
	db := sqlite.GetDBConnection()
	var user User
	err := db.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Admin, &user.Editor, &user.Moderator, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (n *User) GetAllUsers() ([]User, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, first_name, last_name, email, password, admin, editor, moderator, verified, created_at, updated_at FROM users")
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()

	result := []User{}
	for rows.Next() {
		elem := User{}
		err := rows.Scan(&elem.ID, &elem.FirstName, &elem.LastName, &elem.Email, &elem.Password, &elem.Admin, &elem.Editor, &elem.Moderator, &elem.Verified, &elem.CreatedAt, &elem.UpdatedAt)
		if err != nil {
			continue
		}
		result = append(result, elem)
	}

	return result, nil
}

func (n *User) GetAllModerators() ([]User, error) {
	db := sqlite.GetDBConnection()
	rows, err := db.Query("SELECT id, first_name, last_name, email, password, admin, editor, moderator, verified, created_at, updated_at FROM users WHERE moderator = 1")
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()

	result := []User{}
	for rows.Next() {
		elem := User{}
		err := rows.Scan(&elem.ID, &elem.FirstName, &elem.LastName, &elem.Email, &elem.Password, &elem.Admin, &elem.Editor, &elem.Moderator, &elem.Verified, &elem.CreatedAt, &elem.UpdatedAt)
		if err != nil {
			continue
		}
		result = append(result, elem)
	}

	return result, nil
}

func (n *User) Register() error {
	if CheckEmailTaken(n.Email) {
		return fmt.Errorf("register.invalid.email")
	}

	db := sqlite.GetDBConnection()

	n.Admin = false
	n.Editor = true
	n.Moderator = false
	n.Verified = false
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	result, err := db.Exec("INSERT INTO users (first_name, last_name, email, password, admin, editor, moderator, verified, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		n.FirstName, n.LastName, n.Email, n.Password, n.Admin, n.Editor, n.Moderator, n.Verified, n.CreatedAt, n.UpdatedAt)
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

func (n *User) Login(email string, password string) (User, error) {
	db := sqlite.GetDBConnection()
	var user User
	err := db.QueryRow("SELECT id, first_name, last_name, email, password, admin, editor, moderator, verified, created_at, updated_at FROM users WHERE email = ? AND password = ?", email, password).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Admin, &user.Editor, &user.Moderator, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, fmt.Errorf("register.invalid.email")
	}

	return user, nil
}

func (n *User) UpdateUser() error {
	if CheckEmailTakenWithId(n.Email, n.ID) {
		return fmt.Errorf("register.invalid.email")
	}

	db := sqlite.GetDBConnection()
	n.UpdatedAt = time.Now()

	_, err := db.Exec("UPDATE users SET first_name = ?, last_name = ?, email = ?, admin = ?, editor = ?, moderator = ?, verified = ?, updated_at = ? WHERE id = ?",
		n.FirstName, n.LastName, n.Email, n.Admin, n.Editor, n.Moderator, n.Verified, n.UpdatedAt, n.ID)
	if err != nil {
		return err
	}

	return nil
}

func (n *User) DeleteUser() error {
	db := sqlite.GetDBConnection()
	_, err := db.Exec("DELETE FROM users WHERE email = ?", n.Email)
	if err != nil {
		return err
	}

	return nil
}
