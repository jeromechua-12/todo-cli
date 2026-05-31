package storage

import (
	"fmt"
	"database/sql"

	"github.com/jeromechua-12/todo-cli/internal/task"
)

type SqliteTaskStorage struct {
	db *sql.DB
}

func NewSqliteTaskStorage(db *sql.DB) *SqliteTaskStorage {
	return &SqliteTaskStorage{db}
}


// insert new task into datbase
func (s *SqliteTaskStorage) Add(t task.Task) error {
	query := `
	INSERT INTO todo (desc, status, deadline, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, t.Desc, t.Status, t.Deadline, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error adding task: %v", err)
	}
	_, err = result.LastInsertId()
	return nil
}

// query task by id from database
func (s *SqliteTaskStorage) GetByID(id int) (*task.Task, error) {
	return nil, nil
}

func (s *SqliteTaskStorage) GetAll() ([]task.Task, error) {
	return nil, nil
}
func (s *SqliteTaskStorage) GetByStatus(task.Status) ([]task.Task, error) {
	return nil, nil
}
func (s *SqliteTaskStorage) UpdateTask(t *task.Task) error {
	return nil
}
func (s *SqliteTaskStorage) UpdateStatus(t *task.Task) error {
	return nil
}

func (s *SqliteTaskStorage) Delete(id int) error {
	return nil
}
