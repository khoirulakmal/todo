package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int
	UserName  string
	Email     string
	Password  []byte
	Timestamp time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	statement := "INSERT INTO users (name, email, password, created) values(?, ?, ?, UTC_TIMESTAMP())"
	_, err = m.DB.Exec(statement, name, email, hashed)
	if err != nil {
		// If this returns an error, we use the errors.As() function to check
		// whether the error has the type *mysql.MySQLError. If it does, the
		// error will be assigned to the mySQLError variable. We can then check
		// whether or not the error relates to our users_uc_email key by
		// checking if the error code equals 1062 and the contents of the error
		// message string. If it does, we return an ErrDuplicateEmail error.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Auth(email, password string) (int, error) {
	var user User
	var hashed []byte
	statement := "SELECT id, password FROM users WHERE email = ?"
	err := m.DB.QueryRow(statement, email).Scan(&user.ID, &hashed)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	return user.ID, nil
}

// func (m *UserModel) Exist(id int) (bool, error) {

// }
