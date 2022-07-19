package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	app := fiber.New()

	app.Use(recover.New())

	cmdPath := "/opt/minifi/bin/minifi.sh"

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("MINIFI API service (version 1.0)")
	})

	app.Get("/flow-status", func(c *fiber.Ctx) error {
		name := c.Query("name")
		query := fmt.Sprintf("processor:%s:health,stats,bulletins", name)

		cmd := exec.Command(cmdPath, "flowStatus", query)

		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()

		if err != nil {
			return c.Status(500).JSON(err)
		}

		fmt.Printf("%q\n", out.String())

		s := strings.Replace(out.String(), "\n\n", "#", -1)
		v := strings.Split(s, "#")

		msg := v[2]
		status := FlowStatusStruct{}
		errJson := json.Unmarshal([]byte(msg), &status)

		if errJson != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  500,
				"message": err.Error(),
			})
		}
		return c.JSON(status)
	})

	app.Post("/status", func(c *fiber.Ctx) error {
		cmd := exec.Command(cmdPath, "start")

		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  500,
				"message": err.Error(),
			})
		}

		fmt.Printf("%q\n", out.String())

		s := strings.Replace(out.String(), "\n\n", "#", -1)
		v := strings.Split(s, "#")

		msg := v[2]

		return c.JSON(fiber.Map{
			"status":  200,
			"message": msg,
		})
	})

	app.Post("/start", func(c *fiber.Ctx) error {
		cmd := exec.Command(cmdPath, "start")

		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  500,
				"message": err.Error(),
			})
		}

		fmt.Printf("%q\n", out.String())

		s := strings.Replace(out.String(), "\n\n", "#", -1)
		v := strings.Split(s, "#")

		msg := v[2]

		return c.JSON(fiber.Map{
			"status":  200,
			"message": msg,
		})
	})

	app.Post("/stop", func(c *fiber.Ctx) error {
		cmd := exec.Command(cmdPath, "stop")

		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  500,
				"message": err.Error(),
			})
		}

		fmt.Printf("%q\n", out.String())

		s := strings.Replace(out.String(), "\n\n", "#", -1)
		v := strings.Split(s, "#")

		msg := v[2]

		return c.JSON(fiber.Map{
			"status":  200,
			"message": msg,
		})
	})

	log.Fatal(app.Listen(":3000"))

}
