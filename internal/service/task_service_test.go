package service

import (
	"testing"
	"time"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/jeromechua-12/todo-cli/internal/task"
	"github.com/jeromechua-12/todo-cli/internal/storage"
)

func TestParseDateString(t *testing.T) {
	tests := []struct {
		name string
		input string
		expected time.Time
		expectError bool
	}{
		{
			name: "valid date without time",
			input: "2006-01-02",
			expected: time.Date(2006, 1, 2, 23, 59, 0, 0, time.Local),
			expectError: false,
		},
		{
			name: "valid date and time",
			input: "2006-01-02 15:04",
			expected: time.Date(2006, 1, 2, 15, 4, 0, 0, time.Local),
			expectError: false,
		},
		{
			name: "invalid format using slashes",
			input: "01/02/06 15:04",
			expectError: true,
		},
		{
			name: "invalid format with seconds",
			input: "2006-01-02 15:04:05",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDeadline, err := parseDateString(&tt.input)

			// test for expected errors
			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error for %q, but got none", tt.input)
				}
				return
			}

			// test for unexpected errors
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tt.input, err)
			}

			// test for correct deadline parsed
			if *gotDeadline != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, *gotDeadline)
			}
		})
	}
}

func initService(t *testing.T) *TaskService {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	// defer db connection closing
	t.Cleanup(func() {
		db.Close()
	})

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS todo (
		id INTEGER PRIMARY KEY,
		desc TEXT NOT NULL,
		status TEXT NOT NULL CHECK (status in ('todo', 'in-progress', 'done')),
		deadline DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME
	)
	`)
	if err != nil {
		t.Fatal(err)
	}
	store := storage.NewSqliteTaskStorage(db)
	return NewTaskService(store)
}

func TestAddTask(t *testing.T) {
	service := initService(t)

	var validDesc = "task 1"
	var validDeadline = "2006-01-02 15:04"
	var invalidDeadline = "2 Jan 2006"

	tests := []struct {
		name string
		descInput string
		deadlineInput *string
		expectError bool
	}{
		{
			name: "valid desc without deadline",
			descInput: validDesc,
			deadlineInput: nil,
			expectError: false,
		},
		{
			name: "valid desc and deadline",
			descInput: validDesc,
			deadlineInput: &validDeadline,
			expectError: false,
		},
		{
			name: "empty desc",
			descInput: "",
			deadlineInput: &validDeadline,
			expectError: true,
		},
		{
			name: "invalid deadline",
			descInput: "Task 4",
			deadlineInput: &invalidDeadline,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.AddTask(tt.descInput, tt.deadlineInput)

			// test for expected errors
			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error for %q, but got none", tt.name)

				}
				return
			}

			// test for unexpected errors
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.name, err)
			}
		})
	}
}
