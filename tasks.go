package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"time"
)

type Task struct {
	Id			int
	Item    	sql.NullInt64
	Product    	sql.NullInt64
	Type		string
	User		int
	Status		string
	Reference	sql.NullInt64
	Quantity	int
	Repeating	int
	Start		time.Time
	End			time.Time
}

func CheckTasks(db *sql.DB) {

	stmtOut, err := db.Prepare("SELECT id, type, product_id, item_id, user_id, reference, status, quantity, repeating, start, end FROM tasks WHERE status = 'processing' AND end <= ?")
	if err != nil {
		fmt.Println(err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(time.Now())
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.Id, &task.Type, &task.Product, &task.Item, &task.User, &task.Reference, &task.Status, &task.Quantity, &task.Repeating, &task.Start, &task.End)
		go Process(db, task)
	}
}

func Process(db *sql.DB, task Task) {
	// Gather, Craft, etc.
	switch task.Type {
	case "gathering":
		fmt.Printf("Gathering: #%d - (Item: %d)\n", task.Id, task.Item)
		Gather(db, task)
	case "crafting":
		fmt.Printf("Crafting: #%d - (Item: %d)\n", task.Id, task.Item)
		Craft(db, task)
	case "company":
		fmt.Printf("Company: #%d - (Item: %d)\n", task.Id, task.Item)
		Produce(db, task)
	case "wage":
		fmt.Printf("Wage: #%d (Employee: %d, cost: %d)\n", task.Id, task.Reference.Int64, task.Quantity)
		Wage(db, task)
	}
}

func Gather(db *sql.DB, task Task) {
	// Add to users inventory
	stmtUpd, err := db.Prepare("INSERT INTO item_user (item_id, user_id, quantity) VALUES (?,?,?) ON DUPLICATE KEY UPDATE quantity=quantity+?")
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

func Wage(db *sql.DB, task Task) {
	// Add to users inventory
	stmtUpd, err := db.Prepare("UPDATE users SET bank = bank - ? WHERE id = ?");
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtUpd.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtUpd.Exec(task.Quantity, task.User)
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

	// Execute update
	_, err = stmtUpd.Exec(task.Id)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Was it a repeating task?
	if task.Repeating == 1 {
		stmtIns, err := db.Prepare("INSERT INTO tasks (type, user_id, reference, quantity, repeating, end) VALUES(?, ?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

		// Execute insert
		_, err = stmtIns.Exec(task.Type, task.User, task.Reference, task.Quantity, task.Repeating, task.End.AddDate(0,0,1))
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
}