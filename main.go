package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
    "database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "log"
	"fmt"
)

type Task struct {
	Images  string `json:"images"`
	Task    string `json:"task"`
	BoardId string `json:"boardId"`
}

func main() {
    db, err := sql.Open("sqlite3", ":memory:")
	checkErr(err, "", true)
	defer db.Close()

	_, err = db.Exec("create table tasks (boardId text, task text, images text)")
	checkErr(err, "", true)

	_, err = db.Exec("insert into tasks values ('hello', 'world', 'https://via.placeholder.com/150'), ('hello', 'scott', ''), ('test', 'awesome', '')")
	checkErr(err, "", true)

	app := fiber.New()
	app.Use(cors.New(cors.Config{ AllowOrigins: "*", AllowHeaders: "Content-Type", }))
	app.Static("/", "./index.html")
	app.Post("/api/task", func(c *fiber.Ctx) error { return postTask(c, db) })
	app.Get("/api/tasks/:boardId", func(c *fiber.Ctx) error { return getTasksFromBoard(c, db) })
	app.Listen(":8080")
}

func checkErr(err error, msg string, shouldPanic bool) {
    if err != nil {
		fmt.Println(err, msg)
		if (shouldPanic) {
			panic(err)
		}
    }
}

/* client test:
	var range = window.getSelection().getRangeAt(0);
	var frag = range.cloneContents();
	var string = window.getSelection().toString();
	console.log(string, Array.from(frag.querySelectorAll('img')).map(a => a.getAttribute('src')));
*/

/* test:
	fetch('http://localhost:8080/api/task', {
		method: 'POST',
		mode: 'cors',
		headers: {
		  'Content-Type': 'application/json'
		},
		body: JSON.stringify({ task: 'hello task', boardId: 'hello' })
	})
*/
func postTask(c *fiber.Ctx, db *sql.DB) error {
	body := new(Task)
	if err := c.BodyParser(body); err != nil {
		checkErr(err, "couln't parse body", false)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	_, err := db.Exec("insert into tasks values (?, ?, ?)", body.BoardId, body.Task, body.Images)
	checkErr(err, "coulnt't insert task into board", false)

	return c.SendStatus(fiber.StatusOK)
}

/* test:
	fetch('http://localhost:8080/api/tasks/hello').then(a => a.json()).then(a => console.log(a))
*/
func getTasksFromBoard(c *fiber.Ctx, db *sql.DB) error {
	var boardId = c.Params("boardId")
	rows, err := db.Query("select * from tasks where boardId = ?", boardId)
	checkErr(err, "coulnt't select board id", false)
	defer rows.Close()

	var tasks []Task
	var task Task
	for rows.Next() {
		rows.Scan(&task.BoardId, &task.Task, &task.Images)
		tasks = append(tasks, task)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{ "tasks": tasks })
}
