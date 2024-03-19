package server

import (
	"encoding/json"
	"strconv"
	"fmt"
	"shitblog-server/db"
	"shitblog-server/utils"
	"github.com/gofiber/fiber/v2"
)

func createUser(c *fiber.Ctx) error {
	username := c.Params("username")
	res := db.CreateUser(username)
	if res == "BAN" {
		return c.SendStatus(400)
	}
	if res == "username taken" {
		return c.SendStatus(406)
	}
	return c.SendString(res)
}

func deleteUser(c *fiber.Ctx) error {
	username := c.Params("username")
	var body map[string]string
	err := json.Unmarshal(c.Body(), &body)
	if err != nil {
		return c.SendStatus(400)
	}
	res := db.DeleteUser(username, body["token"])
	if res == 2 {
		return c.SendStatus(401)
	}
	if res == 1 {
		return c.SendStatus(400)
	}
	return c.SendStatus(200)
}

func getUsers(c *fiber.Ctx) error {
	arr := db.GetUsers()
	bytes, err := json.Marshal(arr)
	utils.PanicIfError(err)
	return c.Send(bytes)
}

func createPost(c *fiber.Ctx) error {
	var body map[string]string
	err := json.Unmarshal(c.Body(), &body)
	if err != nil {
		fmt.Println("unmarshal error")
		return c.SendStatus(400)
	}
	if body["token"] == "" || body["title"] == "" || body["text"] == "" {
		fmt.Println("not all fields")
		return c.SendStatus(400)
	}
	res := db.CreatePost(body["token"], body["title"], body["text"])
	if res == 2 {
		return c.SendStatus(401)
	}
	if res == 1 {
		fmt.Println("DB error")
		return c.SendStatus(400)
	}
	return c.SendStatus(200)
}

func deletePost(c *fiber.Ctx) error {
	strid := c.Params("id")
	id, err := strconv.ParseUint(strid, 10, 64)
	if err != nil {
		return c.SendStatus(400)
	}
	var body map[string]string
	err = json.Unmarshal(c.Body(), &body)
	if err != nil || body["token"] == "" {
		return c.SendStatus(400)
	}
	res := db.DeletePost(body["token"], id)
	if res == 2 {
		return c.SendStatus(401)
	}
	return c.SendStatus(200)
}

func getPosts(c *fiber.Ctx) error {
	var body map[string]string
	err := json.Unmarshal(c.Body(), &body)
	if err != nil {
		return c.SendStatus(400)
	}
	if body["count"] == "" {
		return c.SendStatus(400)
	}
	count, err := strconv.Atoi(body["count"])
	if err != nil {
		return c.SendStatus(400)
	}
	posts := db.GetPosts(count, body["author"])
	resp, err := json.Marshal(posts)
	utils.PanicIfError(err)
	return c.Send(resp)
}

func getPost(c *fiber.Ctx) error {
	strid := c.Params("id")
	id, err := strconv.ParseUint(strid, 10, 64)
	utils.PanicIfError(err)
	resp, err := json.Marshal(db.GetPost(id))
	utils.PanicIfError(err)
	return c.Send(resp)
}

func StartServer() {
	app := fiber.New()
	app.Post("/api/user/:username", createUser)
	app.Post("/api/post/", createPost)
	app.Delete("/api/user/:username", deleteUser)
	app.Delete("/api/post/:id", deletePost)
	app.Get("/api/user/", getUsers)
	app.Get("/api/post/", getPosts)
	app.Get("/api/post/:id", getPost)
	app.Listen(":3000")
}
