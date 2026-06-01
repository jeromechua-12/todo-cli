package cmd

import (
	"github.com/spf13/cobra"
	"github.com/jeromechua-12/todo-cli/internal/task"
	"github.com/jeromechua-12/todo-cli/internal/service"
)

func NewListCmd(svc *service.TaskService) *cobra.Command {
	var cmd = &cobra.Command{
		Use: "list",
		Short: "Lists tasks in todo list",
		Long: `Lists all tasks in todo list or by status.

Accepted status values:
	- todo
	- in-progress
	- done

Examples:
	todo-cli list: lists all task
	todo-cli todo: lists all task that have status "todo"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var tasks []task.Task
			var err error
			if len(args) == 1 {
				tasks, err = svc.GetTasksByStatus(args[0])
			} else {
				tasks, err = svc.GetAllTasks()
			}
			if err != nil {
				return err
			}
			task.Print(tasks)
			return nil
		},
	}
	return cmd
}
