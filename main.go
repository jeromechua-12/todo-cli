package main

import (
	"log"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/jeromechua-12/todo-cli/cmd"
	"github.com/jeromechua-12/todo-cli/internal/storage"
	"github.com/jeromechua-12/todo-cli/internal/service"
)

func main() {
	// establish db connection
	db, err := sql.Open("sqlite3", "todo.db")
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	taskStorage := storage.NewSqliteTaskStorage(db)
	taskService := service.NewTaskService(taskStorage)

	rootCmd := cmd.GetRootCmd()
	rootCmd.AddCommand(cmd.NewAddCmd(taskService))
	rootCmd.AddCommand(cmd.NewListCmd(taskService))
	rootCmd.AddCommand(cmd.NewDeleteCmd(taskService))
	rootCmd.AddCommand(cmd.NewUpdateTaskCmd(taskService))
	rootCmd.AddCommand(cmd.NewUpdateStatusCmd(taskService))

	cmd.Execute()
}
