package task

import (
	"fmt"
	"strings"
	"strconv"
	"time"
	"os"

	"github.com/olekukonko/tablewriter"
)

// enum for task status
type Status string 

const (
	ToDo Status = "todo" 
	InProgress Status = "in-progress"
	Done Status = "done"
)

// core data structure
type Task struct {
	ID int
	Desc string
	Status Status 
	Deadline *time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// constructor for Task
func NewTask(desc string, deadline *time.Time) (Task, error) {
	t := Task{Desc: desc, Status: ToDo, Deadline: deadline, CreatedAt: time.Now()}
	err := ValidateTask(t); if err != nil {
		return Task{}, err
	}
	return t, nil
}

// validates the fields in Task
func ValidateTask(t Task) error {
	if t.Desc == "" {
		return fmt.Errorf("task cannot have an empty description")
	}
	if !isValidStatus(t.Status) {
		return fmt.Errorf("invalid Status %s. Should be one of: %s, %s, %s",
		t.Status, ToDo, InProgress, Done)
	}
	if t.Deadline != nil && t.Deadline.IsZero() {
		return fmt.Errorf("deadline should not be zero time")
	}
	if t.CreatedAt.IsZero() {
		return fmt.Errorf("createdAt should not be zero time")
	}
	if t.UpdatedAt != nil && t.UpdatedAt.IsZero() {
		return fmt.Errorf("updatedAt should not be zero time")
	}
	return nil
}

// prints out Task fields separated by a '|'
func Print(tasks []Task) {
	table := tablewriter.NewTable(os.Stdout)
	table.Header("ID", "Description", "Status", "Deadline", "Created At", "Updated At")

	// format and append tasks to table
	for _, t := range tasks {
		// format date to yyyy-mm-dd hh:mm layout
		var deadline, createdAt, updatedAt string
		if t.Deadline != nil {
			deadline = t.Deadline.Format("2006-01-02 15:04")
		} else {
			deadline = "nil"
		}
		if t.UpdatedAt!= nil {
			updatedAt = t.UpdatedAt.Format("2006-01-02 15:04")
		} else {
			updatedAt = "nil"
		}
		createdAt = t.CreatedAt.Format("2006-01-02 15:04")

		table.Append([]string{
			strconv.Itoa(t.ID),
			t.Desc,
			string(t.Status),
			deadline,
			createdAt,
			updatedAt,
		})
	}

	// render table
	fmt.Println()
	table.Render()
	fmt.Println()
}

// parses a string to Status enum
func ParseStatus(input string) (Status, error) {
	// format str to lower case and trim spaces
	formattedStr := strings.ToLower(strings.TrimSpace(input))
	status := Status(formattedStr)

	if !isValidStatus(status) {
		return "", fmt.Errorf("invalid Status %s. Should be one of: %s, %s, %s",
		input, ToDo, InProgress, Done)
	}
	return status, nil
}

func isValidStatus(s Status) bool {
	switch s {
	case ToDo:
		return true
	case InProgress:
		return true
	case Done:
		return true
	}
	return false
}
