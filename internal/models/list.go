package models

import (
	"database/sql"
	"time"
)

type ListInterface interface {
	Insert(content string, status string) (int, error)
	Get(id int) (*List, error)
	GetRows() ([]*List, error)
	Delete(id int64) (bool, error)
	Done(id int64) (bool, error)
}

type List struct {
	ID      int
	Content string
	Status  string
	Date    time.Time
}

type TodoModel struct {
	DB *sql.DB
}

func (m *TodoModel) Insert(content string, status string) (int, error) {
	statement := "INSERT INTO lists (content, status, date) values(?, ?, UTC_TIMESTAMP())"
	result, err := m.DB.Exec(statement, content, status)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), err
}

func (m *TodoModel) Get(id int) (*List, error) {
	todoList := &List{}
	statement := "SELECT id, content, status, date FROM lists WHERE id = ?"
	err := m.DB.QueryRow(statement, id).Scan(&todoList.ID, &todoList.Content, &todoList.Status, &todoList.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return todoList, nil
}

func (m *TodoModel) GetRows() ([]*List, error) {
	var lists []*List
	statement := "SELECT * FROM lists"
	result, err := m.DB.Query(statement)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		var temp List
		if err := result.Scan(&temp.ID, &temp.Content, &temp.Date, &temp.Status); err != nil {
			return lists, err
		}
		lists = append(lists, &temp)
	}
	err = result.Err()
	if err != nil {
		return lists, err
	}
	return lists, nil
}

func (m *TodoModel) Delete(id int64) (bool, error) {
	statement := "DELETE FROM lists WHERE id = ?"
	_, err := m.DB.Exec(statement, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *TodoModel) Done(id int64) (bool, error) {
	statement := "UPDATE lists SET status = 'Done' WHERE ID = ?"
	_, err := m.DB.Exec(statement, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
