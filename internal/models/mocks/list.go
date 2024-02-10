package mocks

import (
	"database/sql"
	"time"

	"todo.khoirulakmal.dev/internal/models"
)

var ListMock = &models.List{
	ID:      1,
	Content: "Content mock for test",
	Status:  "Status mock for test",
	Date:    time.Now(),
}

type TodoModel struct{}

func (m *TodoModel) Insert(content string, status string) (int, error) {
	return 2, nil
}

func (m *TodoModel) Get(id int) (*models.List, error) {
	switch id {
	case 1:
		return ListMock, nil
	default:
		return nil, sql.ErrNoRows
	}
}

func (m *TodoModel) GetRows() ([]*models.List, error) {
	return []*models.List{ListMock}, nil
}

func (m *TodoModel) Delete(id int64) (bool, error) {
	if id == 1 {
		return true, nil
	}
	return false, models.ErrNoRecord
}

func (m *TodoModel) Done(id int64) (bool, error) {
	if id == 1 {
		return true, nil
	}
	return false, models.ErrNoRecord
}
