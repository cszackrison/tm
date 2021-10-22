package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
    "database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "log"
	"fmt"
    "strings"
	"net/url"
	"flag"
)

type Task struct {
	Images string `json:"images"`
	Task string `json:"task"`
	BoardId string `json:"boardId"`
	ListId string `json:"listId"`
	Id string `json:"id"`
	Priority string `json:"priority"`
}

const MEM string = ":memory:";

func main() {
	dbPath := flag.String("db", MEM, "a path to a sqlite db")
	port := flag.String("port", ":8080", "a port to run on")
	flag.Parse()
	fmt.Println(*dbPath)

    db, err := sql.Open("sqlite3", *dbPath)
	checkErr(err, "", true)
	defer db.Close()

	_, err = db.Exec("create table if not exists tasks (id text not null primary key, boardId text, listId text, task text, images text, priority real)")
	checkErr(err, "", true)

	if *dbPath == MEM {
		_, err = db.Exec(`insert into tasks values
			('1', 'hello', 'list 1', '1. world', 'https://via.placeholder.com/150', '1'),
			('2', 'hello', 'list 1', '2. scott', '', '2'),
			('3', 'hello', 'list 1a', '3. scott', 'https://via.placeholder.com/300', '3'),
			('4', 'hello', 'list 2', '4. scott', '', '4'),
			('5', 'hello', 'list 3', '5. scawer awerott', '', '5'),
			('6', 'hello', 'list 4', '6. sco awer awer tt', '', '6'),
			('7', 'hello', 'list 5', 'scott awwerawerawer', '', '7'),
			('8', 'test', 'list 2', 'awesome', '', '8')
		`)
		checkErr(err, "", true)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{ AllowOrigins: "*", AllowHeaders: "Content-Type", }))
	app.Static("/", "./index.html")
	app.Delete("/api/task/:id", func(c *fiber.Ctx) error { return deleteTask(c, db) })
	app.Post("/api/task", func(c *fiber.Ctx) error { return postTask(c, db) })
	app.Patch("/api/task/:id", func(c *fiber.Ctx) error { return patchTask(c, db) })
	app.Patch("/api/board/:boardId/list/:listId", func(c *fiber.Ctx) error { return patchList(c, db) })
	app.Get("/api/tasks/:boardId", func(c *fiber.Ctx) error { return getTasksFromBoard(c, db) })
	app.Listen(*port)
}

func checkErr(err error, msg string, shouldPanic bool) {
    if err != nil {
		fmt.Println(err, msg)
		if shouldPanic {
			panic(err)
		}
    }
}

func deleteTask(c *fiber.Ctx, db *sql.DB) error {
	var id = c.Params("id")
	_, err2 := db.Exec(fmt.Sprintf("delete from tasks where id = '%s'", id))
	checkErr(err2, "cannot delete task", false)
	return c.SendStatus(fiber.StatusOK)
}
/* client test:
	var range = window.getSelection().getRangeAt(0);
	var frag = range.cloneContents();
	var string = window.getSelection().toString();
	console.log(string, Array.from(frag.querySelectorAll('img')).map(a => a.getAttribute('src')));
*/
func postTask(c *fiber.Ctx, db *sql.DB) error {
	body := new(Task)
	if err := c.BodyParser(body); err != nil {
		checkErr(err, "cannot parse body", false)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	newId := uuid.New()
	_, err := db.Exec("insert into tasks values (?, ?, ?, ?, ?, ?)", newId, body.BoardId, body.ListId, body.Task, body.Images, body.Priority)
	checkErr(err, "cannot insert task into board", false)

	body.Id = newId.String()
	return c.Status(fiber.StatusOK).JSON(body)
}

func patchList(c *fiber.Ctx, db *sql.DB) error {
	listId, err3 := url.QueryUnescape(c.Params("listId"))
	checkErr(err3, "cannot decode listId", false)
	boardId, err4 := url.QueryUnescape(c.Params("boardId"))
	checkErr(err4, "cannot decode boardId", false)

	body := new(Task)
	if err := c.BodyParser(body); err != nil {
		checkErr(err, "cannot parse body", false)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	values := []string{}
	if body.ListId != "" {
	    values = append(values, fmt.Sprintf("listId = '%s'", body.ListId))
	}
	result := strings.Join(values, ", ")

	if result != "" {
		_, err2 := db.Exec(fmt.Sprintf("update tasks set %s where listId = '%s' and boardId = '%s'", result, listId, boardId))
		checkErr(err2, "cannot update list", false)
	}

	return c.SendStatus(fiber.StatusOK)
}

func patchTask(c *fiber.Ctx, db *sql.DB) error {
	var id = c.Params("id")

	body := new(Task)
	if err := c.BodyParser(body); err != nil {
		checkErr(err, "cannot parse body", false)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	values := []string{}
	if body.ListId != "" {
	    values = append(values, fmt.Sprintf("listId = '%s'", body.ListId))
	}
	if body.Task != "" {
	    values = append(values, fmt.Sprintf("task = '%s'", body.Task))
	}
	if body.Priority != "" {
	    values = append(values, fmt.Sprintf("priority = '%s'", body.Priority))
	}
	result := strings.Join(values, ", ")

	if result != "" {
		_, err2 := db.Exec(fmt.Sprintf("update tasks set %s where id = '%s'", result, id))
		checkErr(err2, "cannot update task", false)
	}

	return c.SendStatus(fiber.StatusOK)
}

func getTasksFromBoard(c *fiber.Ctx, db *sql.DB) error {
	var boardId = c.Params("boardId")
	rows, err := db.Query("select * from tasks where boardId = ?", boardId)
	checkErr(err, "cannot select board id", false)
	defer rows.Close()

	var tasks []Task
	var task Task
	for rows.Next() {
		rows.Scan(&task.Id, &task.BoardId, &task.ListId, &task.Task, &task.Images, &task.Priority)
		tasks = append(tasks, task)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{ "tasks": tasks })
}
