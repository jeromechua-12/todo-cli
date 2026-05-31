package service

import (
	"fmt"
	"time"
	"strings"

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

func parseDateString(date *string) (*time.Time, error) {
	if date == nil {
		return nil, nil
	}

	// parse date. Only accept yyyy-mm-dd or yyyy-mm-dd hh:mm format
	*date = strings.TrimSpace(*date)
	if len(*date) != 10 && len(*date) != 16 {
		return nil, fmt.Errorf("invalid format for deadline; expected yyyy-mm-dd or yyyy-mm-dd hh:mm")
	}

	var formattedDeadline *time.Time
	validDateLayouts := []string{"2006-01-02", "2006-01-02 15:04"}
	for _, layout := range validDateLayouts {
		dl, err := time.ParseInLocation(layout, *date, time.Local)
		if err == nil {
			// if yyyy-mm-dd, set time as 23:59
			if len(layout) == 10 {
				dl = dl.Add(time.Hour * 23).Add(time.Minute * 59)
			}
			formattedDeadline = &dl
			break
		}
	}
	if formattedDeadline == nil {
		return nil, fmt.Errorf("invalid format for deadline; expected yyyy-mm-dd or yyyy-mm-dd hh:mm")
	}
	return formattedDeadline, nil
}

// adds new task
func (s *TaskService) AddTask(desc string, deadline *string) error {
	parsedDeadline, err := parseDateString(deadline)
	if err != nil {
		return err
	}
	task, err := task.NewTask(desc, parsedDeadline)
	if err != nil {
		return err
	}
	return s.store.Add(task)
}

// get task by ID
func (s *TaskService) GetTaskByID(id int) (*task.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("task id cannot be less than 1")
	}
	return s.store.GetByID(id)
}

// get all tasks
func (s *TaskService) GetAllTasks() ([]task.Task, error) {
	return s.store.GetAll()
}
