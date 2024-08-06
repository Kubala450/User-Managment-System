package database

import (
	"database/sql"
	"fmt"
	"strings"

	"main/user"
	"main/utils"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type UserDataBase struct {
	db *sql.DB
}

func NewUserDataBase(dbPath string) (*UserDataBase, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		user_id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &UserDataBase{db: db}, nil
}

func (u *UserDataBase) AddNewUser(username, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	if err := utils.ValidateUserInput(username, password); err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	var storedUsername string
	if err := u.db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&storedUsername); err == nil {
		return fmt.Errorf("user already exists")
	} else if err != sql.ErrNoRows {
		return err
	}

	userID := uuid.New().String()

	insertUserQuery := `
	INSERT INTO users (user_id, username, password)
	VALUES (?, ?, ?);`

	_, err = u.db.Exec(insertUserQuery, userID, username, hashedPassword)
	if err != nil {
		return err
	}

	fmt.Printf("New user added: %s, UUID: %s\n", username, userID)
	return nil
}

func (u *UserDataBase) AuthenticateAdmin(username, password string) bool {
	var storedPassword string
	err := u.db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		fmt.Printf("Error retrieving password for user %s: %s\n", username, err)
		return false
	}

	err = utils.CompareHashAndPassword(storedPassword, password)
	if err != nil {
		fmt.Printf("Error comparing password for user %s: %s\n", username, err)
		return false
	}

	return true
}

func (u *UserDataBase) EditUser(userID, newUsername, newPassword string) error {
	newUsername = strings.TrimSpace(newUsername)
	newPassword = strings.TrimSpace(newPassword)

	if err := utils.ValidateUserInput(newUsername, newPassword); err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	updateUserQuery := `
	UPDATE users
	SET username = ?, password = ?
	WHERE user_id = ?;`

	_, err = u.db.Exec(updateUserQuery, newUsername, hashedPassword, userID)
	if err != nil {
		return err
	}

	fmt.Printf("User with UUID %s updated\n", userID)
	return nil
}

func (u *UserDataBase) ListUsers() ([]user.User, error) {
	rows, err := u.db.Query("SELECT user_id, username, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var user user.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UserDataBase) Close() {
	u.db.Close()
}
