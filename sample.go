package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
    "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"log"
	"fmt"
)

type Task struct {
	Task    string `json:"task"`
	BoardId string `json:"boardId"`
}

type taskRequest struct {
	Info    string       `json:"info"`
	BoardId string 		 `json:"boardId"`
}

type urlRequest struct {
	URL    string        `json:"url"`
	Expiry time.Duration `json:"expiry"`
}

type urlResponse struct {
	URL    string        `json:"url"`
	Short  string   	 `json:"short"`
	Expiry time.Duration `json:"expiry"`
}

func main() {
    db, err := sql.Open("sqlite3", ":memory:")
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("create table tasks (boardId text, task text)")
	checkErr(err)

	_, err = db.Exec("insert into tasks values ('hello', 'world'), ('hello', 'scott'), ('test', 'awesome')")
	checkErr(err)

	app := fiber.New()
	app.Use(cors.New(cors.Config{ AllowOrigins: "*", AllowHeaders: "Content-Type", }))
	app.Static("/", "./index.html")
	app.Post("/api/url", ShortenURL)
	app.Post("/api/task", func(c *fiber.Ctx) error { return postTask(c, db) })
	app.Get("/api/tasks/:boardId", func(c *fiber.Ctx) error { return getTasksFromBoard(c, db) })
	app.Listen(":8080")
}

func checkErr(err error) {
    if err != nil {
		log.Println(err)
        panic(err)
    }
}

func hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

/* test:
fetch('http://localhost:8080/api/task', {
	method: 'POST',
	mode: 'cors',
	headers: {
	  'Content-Type': 'application/json'
	},
	body: JSON.stringify({ info: 'hello task', boardId: 'hello board' })
})
*/
func postTask(c *fiber.Ctx, db *sql.DB) error {
	body := new(taskRequest)
	if err := c.BodyParser(body); err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	db.Exec("insert into tasks values (?, ?)", body.BoardId, body.Info)

	log.Println(body.Info)
	log.Println(body.BoardId)
	return c.SendStatus(fiber.StatusOK)
}

/* test:
	fetch('http://localhost:8080/api/tasks/hello').then(a => a.json()).then(a => console.log(a))
*/
func getTasksFromBoard(c *fiber.Ctx, db *sql.DB) error {
	var boardId = c.Params("boardId")
	rows, err := db.Query("select * from tasks where boardId = ?", boardId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tasks []Task
	var task Task
	for rows.Next() {
		rows.Scan(&task.BoardId, &task.Task)
		fmt.Println("testing")
		tasks = append(tasks, task)
	}

	fmt.Println(tasks)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{ "tasks": tasks })
}

/* test:
fetch('http://localhost:8080/api/url', {
	method: 'POST',
	mode: 'no-cors',
	headers: {
	  'Content-Type': 'application/json'
	},
	body: JSON.stringify({ url: 'https://gnu.org', expiry: 10 })
})
*/
func ShortenURL(c *fiber.Ctx) error {
	body := new(urlRequest)
	if err := c.BodyParser(&body); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	r := urlResponse{
		URL: body.URL,
		Short: "aaa",
		Expiry: body.Expiry,
	}

	log.Println(body.URL)
	log.Println(body.Expiry)

	return c.Status(fiber.StatusOK).JSON(r)
}
