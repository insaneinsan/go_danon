package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	username = "root"
	password = "password"
	hostname = "127.0.0.1"
	port     = 5432
	db       = "postgres"
)

type Task struct {
	Id        int64
	Name      string
	Completed bool
}

func dbInsertTask(db *sql.DB, taskName string) error {
	_, err := db.Exec("INSERT INTO tasks(name) VALUES ($1)", taskName)
	return err
}

func dbGetAllTasks(db *sql.DB) ([]*Task, error) {
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.Id, &task.Name, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func dbUpdateTask(db *sql.DB, taskID int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE tasks SET completed = true WHERE id = $1", taskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func dbDeleteTask(db *sql.DB, taskID int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
func main() {
	DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, hostname, port, db)

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return
	}

	tasks, err := dbGetAllTasks(db)
	if err != nil {
		fmt.Println("Error getting tasks:", err)
	}

	for _, task := range tasks {
		fmt.Println("Task ID:", task.Id, " Name:", task.Name, " Completed:", task.Completed)
	}

	err = dbInsertTask(db, "Do assignment 3")
	if err != nil {
		fmt.Println("Error creating task:", err)
	}

	err = dbUpdateTask(db, 1)
	if err != nil {
		fmt.Println("Error updating task:", err)
	}

	err = dbDeleteTask(db, 1)
	if err != nil {
		fmt.Println("Error deleting task:", err)
	}
}
