package models

import (
	"fmt"

	"github.com/usvc/go-password"
)

type User struct {
	Id           uint64
	Nickname     string
	FirstName    string
	Email        string
	CreationDate string
	Password     string
}

var customPolicy = password.Policy{
	MaximumLength:         64,
	MinimumLength:         10,
	MinimumLowercaseCount: 1,
	MinimumUppercaseCount: 1,
	MinimumNumericCount:   1,
	MinimumSpecialCount:   1,
	CustomSpecial:         []byte("`!@"),
}

func (db *UsersDatabase) AddNewUser(u User) (int64, error) {
	if password.Validate(u.Password, customPolicy) != nil {
		return -1, fmt.Errorf("Password doesn't meet requirements")
	}

	result, err := db.dbHandle.Exec(`INSERT INTO users 
		(nickname, first_name, email, creation_date, password) VALUES (?, ?, ?, ?, ?)`,
		u.Nickname, u.FirstName, u.Email, u.CreationDate, u.Password)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}

func (db *UsersDatabase) LoginExists(login string) (bool, error) {
	return db.exists("nickname", login)
}

func (db *UsersDatabase) EmailExists(email string) (bool, error) {
	return db.exists("email", email)
}

func (db *UsersDatabase) exists(field string, fieldValue string) (bool, error) {
	rows, err := db.dbHandle.Query(`
		SELECT EXISTS 
			(SELECT * 
			 FROM users 
			 WHERE `+field+"=?);", fieldValue)

	if err != nil {
		return false, err
	}
	defer rows.Close()

	if !rows.Next() {
		return false, fmt.Errorf("EXISTS Query has incorrect output")
	}

	var exists int
	err = rows.Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists != 0 {
		return true, nil
	}

	return false, nil
}

func (db *UsersDatabase) VerifyLogin(login string, passwordChecked string) bool {
	return db.verify(login, "nickname", passwordChecked)
}

func (db *UsersDatabase) VerifyEmail(login string, passwordChecked string) bool {
	return db.verify(login, "email", passwordChecked)
}

func (db *UsersDatabase) verify(usernameOrEmail string, field string, passwordChecked string) bool {
	rows, err := db.dbHandle.Query(`SELECT password FROM users WHERE `+field+"=?", usernameOrEmail)
	if err != nil {
		return false
	}
	defer rows.Close()

	if !rows.Next() {
		return false
	}

	var passwordOfUser string
	err = rows.Scan(&passwordOfUser)
	if err != nil {
		return false
	}

	return passwordOfUser == passwordChecked
}

func (db *UsersDatabase) GetUserIdForLogin(username string) (uint64, error) {
	return db.getUserId(username, "nickname")
}

func (db *UsersDatabase) GetUserIdForEmail(email string) (uint64, error) {
	return db.getUserId(email, "email")
}

func (db *UsersDatabase) getUserId(usernameOrEmail string, field string) (uint64, error) {
	rows, err := db.dbHandle.Query(`SELECT id FROM users WHERE `+field+"=?", usernameOrEmail)

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, fmt.Errorf("No ID found for user, does it exist?")
	}

	var id uint64
	err = rows.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
