package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
    "database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "log"
	"fmt"
)

type Task struct {
	Images string `json:"images"`
	Task string `json:"task"`
	BoardId string `json:"boardId"`
	ListId string `json:"listId"`
	Id string `json:"id"`
}

func main() {
    db, err := sql.Open("sqlite3", ":memory:")
	checkErr(err, "", true)
	defer db.Close()

	_, err = db.Exec("create table tasks (id text not null primary key, boardId text, listId text, task text, images text)")
	checkErr(err, "", true)

	_, err = db.Exec("insert into tasks values ('1', 'hello', 'list 1', '1. world', 'https://via.placeholder.com/150'), ('2', 'hello', 'list 1', '2. scott', ''), ('3', 'hello', 'list 1a', '3. scott', ''), ('4', 'hello', 'list 2', '4. scott', ''), ('5', 'hello', 'list 3', '5. scawer awerott', ''), ('6', 'hello', 'list 4', '6. sco awer awer tt', ''), ('7', 'hello', 'list 5', 'scott awwerawerawer', ''), ('8', 'test', 'list 2', 'awesome', '')")
	checkErr(err, "", true)

	app := fiber.New()
	app.Use(cors.New(cors.Config{ AllowOrigins: "*", AllowHeaders: "Content-Type", }))
	app.Static("/", "./index.html")
	app.Post("/api/task", func(c *fiber.Ctx) error { return postTask(c, db) })
	app.Patch("/api/task/:id", func(c *fiber.Ctx) error { return patchTask(c, db) })
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

	_, err := db.Exec("insert into tasks values (?, ?, ?, ?, ?, ?)", uuid.New(), body.BoardId, body.ListId, body.Task, body.Images)
	checkErr(err, "coulnt't insert task into board", false)

	return c.SendStatus(fiber.StatusOK)
}

func patchTask(c *fiber.Ctx, db *sql.DB) error {
	var id = c.Params("id")

	body := new(Task)
	if err := c.BodyParser(body); err != nil {
		checkErr(err, "couln't parse body", false)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var listIdPatch string
	if body.ListId != "" {
		listIdPatch = fmt.Sprintf("listId = '%s'", body.ListId)
	}

	if listIdPatch != "" {
		var execStr = fmt.Sprintf("update tasks set %s where id = '%s'", listIdPatch, id);
		_, err2 := db.Exec(execStr)
		checkErr(err2, "coulnt't update task", false)
	}

	return c.SendStatus(fiber.StatusOK)
}

func getTasksFromBoard(c *fiber.Ctx, db *sql.DB) error {
	var boardId = c.Params("boardId")
	rows, err := db.Query("select * from tasks where boardId = ?", boardId)
	checkErr(err, "coulnt't select board id", false)
	defer rows.Close()

	var tasks []Task
	var task Task
	for rows.Next() {
		rows.Scan(&task.Id, &task.BoardId, &task.ListId, &task.Task, &task.Images)
		tasks = append(tasks, task)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{ "tasks": tasks })
}
