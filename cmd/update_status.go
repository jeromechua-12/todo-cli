package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/jeromechua-12/todo-cli/internal/service"
)

func NewUpdateStatusCmd(svc *service.TaskService) *cobra.Command {
	var cmd = &cobra.Command{
		Use: "mark",
		Short: "Update task status",
		Long: `Update task status

You are required to pass in the ID of the task you want to update,
followed by the new status.

Accepted status values:
	- todo
	- in-progress
	- done

Examples:
	- todo-cli mark 1 done
	- todo-cli mark 2 in-progress`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, convErr := strconv.Atoi(args[0])
			if convErr != nil {
				return convErr
			}
			status := args[1]
			err := svc.UpdateTaskStatus(id, status)
			if err != nil {
				return err
			}
			fmt.Printf("Successfully updated task %d status\n", id)
			return nil
		},
	}
	return cmd
}
