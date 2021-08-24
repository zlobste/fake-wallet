package data

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 50
)

type User struct {
	ID       int64  `db:"id" structs:"-"`
	Name     string `db:"name" structs:"name"`
	Surname  string `db:"surname" structs:"surname"`
	Email    string `db:"email" structs:"email"`
	Password string `db:"password" structs:"password"`
}

func (u *User) EncryptPassword() error {
	passwordLength := len(u.Password)
	if passwordLength < MinPasswordLength || passwordLength > MaxPasswordLength {
		return errors.New("Wrong password length!")
	}

	enc, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	u.Password = string(enc)
	return nil
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password)) != nil
}
