package task

import (
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	validTime := time.Date(2026, 1, 1, 15, 16, 0, 0, time.Local)
	timeNow := time.Now()
	zeroTime := time.Time{}

	// 1. Define the table of test cases
	tests := []struct {
		name string
		desc string
		deadline *time.Time
		expectError bool
	}{
		{
			name: "valid task without deadline",
			desc: "task 1",
			deadline: nil,
			expectError: false,
		},
		{
			name: "valid task with deadline",
			desc: "task 2",
			deadline: &validTime,
			expectError: false,
		},
		{
			name: "invalid task with empty description",
			desc: "",
			deadline: &timeNow,
			expectError: true,
		},
		{
			name: "invalid task with zero-value deadline",
			desc: "Fail task 2",
			deadline: &zeroTime,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask(tt.desc, tt.deadline)

			// test for expected errors
			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error for %q, but got none", tt.name)
				}
				return 
			}

			// test for unexpected errors
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// test for correct desc
			if task.Desc != tt.desc {
				t.Errorf("expected desc %q, got %q", tt.desc, task.Desc)
			}

			// test for correct deadline
			if tt.deadline == nil {
				if task.Deadline != nil {
					t.Errorf("expected deadline to be nil, but got %v", task.Deadline)
				}
			} else {
				if task.Deadline == nil || !task.Deadline.Equal(*tt.deadline) {
					t.Errorf("expected deadline %v, got %v", tt.deadline, task.Deadline)
				}
			}
		})
	}
}

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name string
		input string
		expectedStatus Status
		expectError bool
	}{
		{
			name: "valid status todo",
			input: "todo",
			expectedStatus: ToDo,
			expectError: false,
		},
		{
			name: "valid status in-progress",
			input: "in-progress",
			expectedStatus: InProgress,
			expectError: false,
		},
		{
			name: "valid status done",
			input: "done",
			expectedStatus: Done,
			expectError: false,
		},
		{
			name: "invalid status completed",
			input: "completed",
			expectedStatus: Status(""),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := ParseStatus(tt.input)

			// test for expected errors
			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error for %q, but got none", tt.input)
				}
				return
			}

			// test for unexpected error
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.input, err)
			}
			
			// test for correct status
			if status != tt.expectedStatus{
				t.Errorf("expected status %q, got %q", tt.expectedStatus, status)
			}
		})
	}
}
