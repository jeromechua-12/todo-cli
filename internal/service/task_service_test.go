package service

import (
	"database/sql"
	"testing"
	"time"

	"github.com/jeromechua-12/todo-cli/internal/task"
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

func TestGetByStatus(t *testing.T) {
	service := initService(t)

	invalidStatus := "completed"
	emptyStatus := ""

	tests := []struct{
		name string
		inputStatus string
		expectedNums int
		expectedIDs map[int]bool
		expectError bool
	}{
		{
			name: "valid todo",
			inputStatus: "todo",
			expectedNums: 2,
			expectedIDs: map[int]bool{1: false, 2: false},
			expectError: false,

		},
		{
			name: "valid inProgress",
			inputStatus: "in-progress",
			expectedNums: 1,
			expectedIDs: map[int]bool{3: false},
			expectError: false,
		},
		{
			name: "valid done",
			inputStatus: "done",
			expectedNums: 1,
			expectedIDs: map[int]bool{4: false},
			expectError: false,
		},
		{
			name: "empty status",
			expectedNums: 0,
			expectedIDs: nil,
			inputStatus: emptyStatus,
			expectError: true,
		},
		{
			name: "invalid status",
			expectedNums: 0,
			expectedIDs: nil,
			inputStatus: invalidStatus,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := service.GetTasksByStatus(tt.inputStatus)

			// test for expected errors
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, but got none", tt.name)
				}
				return
			}

			// test for unexpected errors
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.name, err)
			}

			// test for correct of rows queried
			if len(tasks) != tt.expectedNums {
				t.Errorf("expected %d rows, got %d", tt.expectedNums, len(tasks))
			}

			// test for correct ids queried
			for _, tk := range tasks {
				_, ok := tt.expectedIDs[tk.ID]
				if !ok {
					t.Errorf("unexpected id %d fetched", tk.ID)
				} else {
					tt.expectedIDs[tk.ID] = true
				}
			}
			for wantID, queried := range tt.expectedIDs {
				if !queried {
					t.Errorf("expected id %d to be fethced, but not", wantID)
				}
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	service := initService(t)

	validDesc := "new task name"
	validDeadline := "2026-10-31 10:00"
	validParsedDeadline := time.Date(2026, 10, 31, 10, 0, 0, 0, time.Local)
	validStatus := "todo"

	invalidID := 999
	invalidStatus := "completed"

	var emptyDesc *string
	var emptyDeadline *string


	tests := []struct{
		name string
		idInput int
		descInput *string
		deadlineInput *string
		statusInput *string
		parsedDeadline time.Time
		expectError bool
	}{
		{
			name: "valid update with both desc and deadline",
			idInput: 1,
			descInput: &validDesc,
			deadlineInput: &validDeadline,
			parsedDeadline: validParsedDeadline,
			expectError: false,
		},
		{
			name: "valid update without desc",
			idInput: 2,
			descInput: emptyDesc,
			deadlineInput: &validDeadline,
			parsedDeadline: validParsedDeadline,
			expectError: false,
		},
		{
			name: "valid update without deadline",
			idInput: 3,
			descInput: &validDesc,
			deadlineInput: emptyDeadline,
			parsedDeadline: validParsedDeadline,
			expectError: false,
		},
		{
			name: "valid status update",
			idInput: 1,
			statusInput: &validStatus,
			expectError: false,
		},
		{
			name: "invalid id",
			idInput: invalidID,
			expectError: true,
		},
		{
			name: "invalid update with no fields",
			idInput: 1,
			descInput: emptyDesc,
			deadlineInput: emptyDeadline,
			expectError: true,
		},
		{
			name: "invalid status",
			idInput: 1,
			statusInput: &invalidStatus,
			expectError: true,
		},
	}

	for _, tt := range tests {
		var err error

		// 2 possible methods to run: update desc/deadline or update status
		t.Run(tt.name, func(t *testing.T) {
			if tt.statusInput != nil {
				err = service.UpdateTaskStatus(tt.idInput, *tt.statusInput)
			} else {
				err = service.UpdateTask(tt.idInput, tt.descInput, tt.deadlineInput)
			}

			// test for expected errors
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, but got none", tt.name)
				}
				return
			}

			// test for unexpected errors
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.name, err)
			}

			// test if fields actually got updated
			tk, err := service.GetTaskByID(tt.idInput)
			if err != nil {
				t.Fatalf("unexpected error fetching task %d: %v", tt.idInput, err)
			}
			// test for desc and deadline update
			if tt.descInput != nil {
				if tk.Desc != *tt.descInput {
					t.Errorf("expected desc %q, got %q", *tt.descInput, tk.Desc)
				}
			}
			if tt.deadlineInput != nil {
				if !tk.Deadline.Equal(tt.parsedDeadline) {
					t.Errorf("expected deadline %q, got %q", tt.parsedDeadline, tk.Deadline)
				}
			}
			// test for status update
			if tt.statusInput != nil {
				if tk.Status != task.Status(*tt.statusInput) {
					t.Errorf("expected status %q, got %q", *tt.statusInput, tk.Status)
				}
			}
			// test if updateTime was updated
			if tk.UpdatedAt == nil {
					t.Errorf("expected updatedAt to be non-nil")
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	service := initService(t)

	tests := []struct{
		name string
		inputID int
		expectError bool
	}{
		{
			name: "valid id deleted",
			inputID: 4,
			expectError: false,
		},
		{
			name: "no id to delete",
			inputID: 999,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteTask(tt.inputID)

			// test for expected errors
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for input id %d, but got none", tt.inputID)
				}
				return
			}

			// test for unexpected errors
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// check if actually deleted
			_, err = service.GetTaskByID(tt.inputID)
			if err == nil {
				t.Errorf("expected error when finding deleted task id but gone none")
			}
		})
	}
}
