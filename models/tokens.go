package models

import "fmt"

type Token struct {
	UserId         int
	Hash           string
	ExpirationDate string
	Permissions    string
}

func (db *UsersDatabase) AddNewTokenForUser(token Token) error {

	result, err := db.dbHandle.Exec(`INSERT INTO tokens 
		VALUES (?, ?, ?, ?)
	`, token.UserId, token.Hash, token.ExpirationDate, token.Permissions)

	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	return err
}

func (db *UsersDatabase) UserTokenExists(user string) (bool, error) {
	rows, err := db.dbHandle.Query(`
		SELECT EXISTS 
			(SELECT * 
			 FROM tokens 
			 WHERE user_id=?`, user)

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
