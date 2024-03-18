package server

import (
	"encoding/json"
	_ "encoding/json"
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

func StartServer() {
	app := fiber.New()
	app.Post("/api/user/:username", createUser)
	app.Delete("/api/user/:username", deleteUser)
	app.Get("/api/user/", getUsers)
	app.Listen(":3000")
}
