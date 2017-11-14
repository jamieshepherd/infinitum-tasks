package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"time"
	"strconv"
)

type Task struct {
	Id		int
	Item    	int
	Type		string
	User		int
	Status		string
	Quantity	int
	Start		time.Time
	End		time.Time
}

func CheckTasks(db *sql.DB) {

	stmtOut, err := db.Prepare("SELECT id, item_id, type, user_id, status, quantity, start, end FROM tasks WHERE status != 'completed' AND end <= ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(time.Now())
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.Id, &task.Item, &task.Type, &task.User, &task.Status, &task.Quantity, &task.Start, &task.End)
		go Process(db, task)
	}
}

func Process(db *sql.DB, task Task) {
	// Gather, Craft, etc.
	switch task.Type {
	case "gathering":
		fmt.Println("Gathering: #" + strconv.Itoa(task.Id) + " (Item: " + strconv.Itoa(task.Item) + ")")
		Gather(db, task)
	case "crafting":
		fmt.Println("Crafting: #" + strconv.Itoa(task.Id) + " (Item: " + strconv.Itoa(task.Item) + ")")
		Craft(db, task)
	case "company":
		fmt.Println("Company: #" + strconv.Itoa(task.Id) + " (Item: " + strconv.Itoa(task.Item) + ")")
		Produce(db, task)
	}
}

func Gather(db *sql.DB, task Task) {
	// Add to users inventory
	stmtUpd, err := db.Prepare("INSERT INTO item_user (item_id, user_id, quantity) VALUES (?,?,?) ON DUPLICATE KEY UPDATE quantity=quantity+?");
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtUpd.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtUpd.Exec(task.Item, task.User, task.Quantity, task.Quantity)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Complete it
	CompleteTask(db, task)
}

func Craft(db *sql.DB, task Task) {
	// Add to users inventory
	stmtUpd, err := db.Prepare("INSERT INTO item_user (item_id, user_id, quantity) VALUES (?,?,?) ON DUPLICATE KEY UPDATE quantity=quantity+?");
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtUpd.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtUpd.Exec(task.Item, task.User, task.Quantity, task.Quantity)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Complete it
	CompleteTask(db, task)
}

func Produce(db *sql.DB, task Task) {
	// Add to users inventory
	stmtUpd, err := db.Prepare("INSERT INTO item_user (item_id, user_id, quantity) VALUES (?,?,?) ON DUPLICATE KEY UPDATE quantity=quantity+?");
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtUpd.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtUpd.Exec(task.Item, task.User, task.Quantity, task.Quantity)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Complete it
	CompleteTask(db, task)
}

func CompleteTask(db *sql.DB, task Task) {
	// Update item, set complete
	stmtUpd, err := db.Prepare("UPDATE tasks SET status = 'completed' WHERE id = ?") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtUpd.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtUpd.Exec(task.Id)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}