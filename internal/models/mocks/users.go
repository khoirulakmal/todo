package mocks

import (
	"todo.khoirulakmal.dev/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "mockemail@gmail.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Auth(email, password string) (int, error) {
	if email == "mockemail@gmail.com" && password == "mock123" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}
