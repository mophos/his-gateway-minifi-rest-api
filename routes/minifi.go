package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mophos/minifi-cli-go/models"
	"github.com/spf13/viper"
)

func GetFlowStatus(ctx *fiber.Ctx) error {
	cmdPath := viper.GetString("cmd")
	name := ctx.Query("name")
	query := fmt.Sprintf("processor:%s:health,stats,bulletins", name)

	cmd := exec.Command(cmdPath, "flowStatus", query)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return ctx.Status(500).JSON(err)
	}

	fmt.Printf("%q\n", out.String())

	s := strings.Replace(out.String(), "\n\n", "#", -1)
	v := strings.Split(s, "#")

	msg := v[2]
	status := models.FlowStatusStruct{}
	errJson := json.Unmarshal([]byte(msg), &status)

	if errJson != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  500,
			"message": err.Error(),
		})
	}
	return ctx.JSON(status)
}

func GetMinifiStatus(ctx *fiber.Ctx) error {
	cmdPath := viper.GetString("cmd")

	var out bytes.Buffer
	cmd := exec.Command(cmdPath, "status")
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  500,
			"message": err.Error(),
		})
	}

	fmt.Printf("%q\n", out.String())

	s := strings.Replace(out.String(), "\n\n", "#", -1)
	v := strings.Split(s, "#")

	msg := v[2]

	return ctx.JSON(fiber.Map{
		"status":  200,
		"message": msg,
	})
}

func StartMinifi(ctx *fiber.Ctx) error {
	cmdPath := viper.GetString("cmd")

	var out bytes.Buffer
	cmd := exec.Command(cmdPath, "start")
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  500,
			"message": err.Error(),
		})
	}

	fmt.Printf("%q\n", out.String())

	s := strings.Replace(out.String(), "\n\n", "#", -1)
	v := strings.Split(s, "#")

	msg := v[2]

	return ctx.JSON(fiber.Map{
		"status":  200,
		"message": msg,
	})
}

func StopMinifi(ctx *fiber.Ctx) error {
	cmdPath := viper.GetString("cmd")

	var out bytes.Buffer
	cmd := exec.Command(cmdPath, "stop")
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  500,
			"message": err.Error(),
		})
	}

	fmt.Printf("%q\n", out.String())

	s := strings.Replace(out.String(), "\n\n", "#", -1)
	v := strings.Split(s, "#")

	msg := v[2]

	return ctx.JSON(fiber.Map{
		"status":  200,
		"message": msg,
	})
}

func RestartMinifi(ctx *fiber.Ctx) error {
	cmdPath := viper.GetString("cmd")

	var out bytes.Buffer
	cmd := exec.Command(cmdPath, "restart")
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  500,
			"message": err.Error(),
		})
	}

	fmt.Printf("%q\n", out.String())

	s := strings.Replace(out.String(), "\n\n", "#", -1)
	v := strings.Split(s, "#")

	msg := v[2]

	return ctx.JSON(fiber.Map{
		"status":  200,
		"message": msg,
	})
}
