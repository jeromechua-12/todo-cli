package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/jeromechua-12/todo-cli/internal/service"
)

func NewDeleteCmd(svc *service.TaskService) *cobra.Command {
	var cmd = &cobra.Command{
		Use: "delete",
		Short: "Delete task",
		Long: `Deletes a task. 

You are required to pass in the id of the task you want to delete.

Example:
todo-cli delete 3  -- deletes task id 3`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args[]string) error {
			id, convErr := strconv.Atoi(args[0])
			if convErr != nil {
				return convErr
			}
			err := svc.DeleteTask(id)
			if err != nil {
				return err
			}
			fmt.Printf("Successfully deleted task %d\n", id)
			return nil
		},
	}
	return cmd
}
