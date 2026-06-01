package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/jeromechua-12/todo-cli/internal/service"
)

func NewUpdateTaskCmd (svc *service.TaskService) *cobra.Command {
	var desc, deadline string

	var cmd = &cobra.Command{
		Use: "update",
		Short: "Update task's description and/or deadline",
		Long: `Update a task's description or deadline or both.

You are required to pass in the ID of the task you want to update.

You can then pass in a new description and/or deadline to update.
Deadline has to be in yyyy-mm-dd hh:mm or yyyy-mm-dd format

Examples:
	todo-cli update 1 -d "New description" -l "New deadline"
	todo-cli update 2 -d "New description"
	todo-cli update 3 -l "New deadline"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, convErr := strconv.Atoi(args[0])
			if convErr != nil {
				return convErr
			}
			var descPtr *string
			var deadlinePtr *string

			if cmd.Flags().Changed("desc") {
				descPtr = &desc
			}
			if cmd.Flags().Changed("deadline") {
				deadlinePtr = &deadline
			}
			err := svc.UpdateTask(id, descPtr, deadlinePtr)
			if err != nil {
				return err
			}
			fmt.Printf("Successfully updated task %d\n", id)
			return nil
		},
	}

	// desc and deadline flags
	cmd.Flags().StringVarP(&desc, "desc", "d", "", "New task description") 
	cmd.Flags().StringVarP(&deadline, "deadline", "l", "", "New task deadline") 

	return cmd
}
