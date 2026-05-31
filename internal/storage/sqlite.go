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
	var t task.Task

	query := "SELECT * FROM todo WHERE id = ?"
	row := s.db.QueryRow(query, id)
	err := row.Scan(&t.ID, &t.Desc, &t.Status, &t.Deadline, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &task.Task{}, fmt.Errorf("no task id %d", id)
		}
		return &task.Task{}, fmt.Errorf("error fetching task: %s", err)
	}
	return &t, nil
}

func (s *SqliteTaskStorage) GetAll() ([]task.Task, error) {
	var tasks []task.Task

	query := "SELECT * FROM todo"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching tasks: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var t task.Task
		err := rows.Scan(&t.ID, &t.Desc, &t.Status, &t.Deadline, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error fetching tasks: %v", err)
		}
		tasks = append(tasks, t)
	}
	if err :=  rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching tasks: %v", err)
	}
	return tasks, nil
}

func (s *SqliteTaskStorage) GetByStatus(status task.Status) ([]task.Task, error) {
	var tasks []task.Task

	query := "SELECT * FROM todo where status = ?"
	rows, err := s.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("error fetching tasks: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var t task.Task
		err := rows.Scan(&t.ID, &t.Desc, &t.Status, &t.Deadline, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error fetching tasks: %v", err)
		}
		tasks = append(tasks, t)
	}
	if err :=  rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching tasks: %v", err)
	}
	return tasks, nil
}

func (s *SqliteTaskStorage) UpdateTask(t *task.Task) error {
	query := `
	UPDATE todo
	SET
		desc = ?,
		status = ?,
		deadline = ?,
		created_at = ?,
		updated_at = ?
	WHERE id = ?;
	`
	_, err := s.db.Exec(query, t.Desc, t.Status, *t.Deadline, t.CreatedAt, *t.UpdatedAt, t.ID)
	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}
	return nil
}

func (s *SqliteTaskStorage) Delete(id int) error {
	query := "DELETE FROM todo WHERE id = ?"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no task id %d", id)
	}
	return nil
}
