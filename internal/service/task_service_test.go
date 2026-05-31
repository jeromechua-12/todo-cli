package service

import (
	"database/sql"
	"testing"
	"time"

	"github.com/jeromechua-12/todo-cli/internal/storage"
	_ "github.com/mattn/go-sqlite3"
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
	// in memory database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	// defer db connection closing
	t.Cleanup(func() {
		db.Close()
	})

	// create table
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

	// insert sample rows
	_, err = db.Exec(`
	INSERT INTO todo (desc, status, deadline, created_at, updated_at)
	VALUES
		("task 1", "todo", NULL, CURRENT_TIMESTAMP, NULL),
		("task 2", "todo", NULL, CURRENT_TIMESTAMP, NULL),
		("task 3", "in-progress", DATETIME("2026-12-31 23:59"), CURRENT_TIMESTAMP, NULL),
		("task 4", "done", NULL, CURRENT_TIMESTAMP, DATETIME("2026-06-01 15:00"));
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

func TestGetTaskByID(t *testing.T) {
	service := initService(t)

	validID := 1
	invalidID := 999
	negativeID := -1

	tests := []struct{
		name string
		inputID int
		expectError bool
	}{
		{
			name: "valid id",
			inputID: validID,
			expectError: false,
		},
		{
			name: "negative id",
			inputID: negativeID,
			expectError: true,
		},
		{
			name: "id not found",
			inputID: invalidID,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTask, err := service.GetTaskByID(tt.inputID)

			// test for expected errors
			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error for %q, but got none", tt.name)
				}
				return
			}

			// test for unexpected errors
			if err != nil {
				t.Fatalf("unexpected error for input %d: %v", tt.inputID, err)
			}

			// test for correct id
			if gotTask.ID != tt.inputID {
				t.Errorf("expected id %d, got %d", tt.inputID, gotTask.ID)
			}
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	service := initService(t)

	expectedNum := 4
	expectedIDs := map[int]bool {
		1: false, 
		2: false, 
		3: false, 
		4: false, 
	}

	tasks, err := service.GetAllTasks()
	if err != nil {
		t.Fatalf("unexpected error for fetching all tasks: %v", err)
	}

	// test for number of tasks fetch
	if len(tasks) != expectedNum {
		t.Errorf("expected %d rows, got %d", expectedNum, len(tasks))

	}

	// test for correct tasks id and desc
	for _, tk := range tasks {
		_, ok := expectedIDs[tk.ID]
		if !ok {
			t.Errorf("unexpected id %d fetched", tk.ID)
		} else {
			expectedIDs[tk.ID] = true
		}
	}
	for wantID, queried := range expectedIDs {
		if !queried {
			t.Errorf("expected id %d to be fethced, but not", wantID)
		}
	}
}
