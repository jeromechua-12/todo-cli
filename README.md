# Todo CLI
[![Go Reference](https://pkg.go.dev/badge/github.com/jeromechua-12/todo-cli.svg)](https://pkg.go.dev/github.com/jeromechua-12/todo-cli)

A simple command-line interface (CLI) application for managing tasks directly from your terminal.

## Overview 

Todo CLI supports the following operations to help you manage your tasks:


### Add
Create a new task with an optional deadline.
```bash
todo-cli add <description> [--deadline]
```

### List
View all your tasks, with an optional filter by status.
```bash
todo-cli list [status]
```

### Update
Modify an existing task's description, deadline, or both.
```bash
todo-cli update <task_id> [--desc] [--deadline]
```


### Mark
Change the status of a specific task (e.g., from "todo" to "done").
```bash
todo-cli mark <task_id> <status>
```

### Delete
Permanently remove a task.
```bash
todo-cli delete <task_id>
```

## Database Initialisation

A SQLite database is required for the CLI to function correctly.

### 1. Install SQLite (If required)

If you do not already have SQLite installed on your machine, you can install it from the [SQLite download page](https://www.sqlite.org/download.html).

### 2. Initialise the database

Once SQLite is installed, run the following command from the project root in your terminal:

```bash
sqlite3 todo.db < ./scripts/init.sql
```
