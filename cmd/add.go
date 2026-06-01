package cmd

import (
	"fmt"

	"github.com/jeromechua-12/todo-cli/internal/service"
	"github.com/spf13/cobra"
)

func NewAddCmd(svc *service.TaskService) *cobra.Command {
	var deadline string

	var cmd = &cobra.Command{
		Use: "add",
		Short: "Add a new task",
		Long: `Add a new task to your todo lists.

You are required to pass in a description of the task.

Additionally you can add an optional deadline.
Deadline format has to be in yyyy-mm-dd or yyyy-mm-dd hh:mm format.

Examples:
	todo-cli add "Task 1"
	todo-cli add "Task 2 -l "2026-10-31"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			desc := args[0]
			var deadlinePtr *string
			if deadline != "" {
				deadlinePtr = &deadline
			}
			id, err := svc.AddTask(desc, deadlinePtr)
			if err != nil {
				return err
			}
			fmt.Printf("Successfully added task (ID: %d)\n", id)
			return nil
		},
	}

	// optional deadline flag
	cmd.Flags().StringVarP(&deadline, "deadline", "l", "", "Task deadline")

	return cmd 
}
