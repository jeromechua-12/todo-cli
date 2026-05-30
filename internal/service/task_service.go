package service

import (
	"github.com/jeromechua-12/todo-cli/internal/task"
)

type TaskStorage interface {
	GetByID(id int) (*task.Task, error)
	Add(t task.Task) error
	GetAll() ([]task.Task, error)
	GetByStatus(task.Status) ([]task.Task, error)
	UpdateTask(t *task.Task) error
	UpdateStatus(t *task.Task) error
	Delete(id int) error
}

type TaskService struct {
	store TaskStorage
}

// constructor for TaskService
func NewTaskService(s TaskStorage) *TaskService {
	return &TaskService{store: s}
}

