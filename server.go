package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mophos/minifi-cli-go/routes"
	"github.com/spf13/viper"
)

func main() {

	app := fiber.New()

	app.Use(recover.New())

	// read configure file
	viper.SetConfigName("env")
	viper.AddConfigPath(".")

	confErr := viper.ReadInConfig()

	if confErr != nil {
		panic(confErr.Error())
	}

	viper.SetDefault("dataPth", "/opt/minifi/data/template")
	viper.SetDefault("outPath", "/opt/minifi/conf")
	viper.SetDefault("settingFile", "/opt/minifi/data/config/setting.yml")

	configRoute := app.Group("configs")
	connectionRoute := app.Group("connections")
	tableRoute := app.Group("tables")
	minifiRoute := app.Group("minifi")

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("MINIFI API service (version 1.0)")
	})

	tableRoute.Get("/query-status/:id", routes.GetTableQueryStatus)
	tableRoute.Post("/manual/create/:id", routes.CreateQueryTable)
	tableRoute.Get("/info/:id/:table", routes.GetTableQueryStatusInfo)

	configRoute.Put("/", routes.UpdateConfig)
	configRoute.Get("/", routes.GetConfig)
	configRoute.Post("/generate", routes.GenerateConfig)

	minifiRoute.Get("/flow-status", routes.GetFlowStatus)
	minifiRoute.Get("/status", routes.GetMinifiStatus)
	minifiRoute.Post("/start", routes.StartMinifi)
	minifiRoute.Post("/stop", routes.StopMinifi)
	minifiRoute.Post("/restart", routes.RestartMinifi)

	connectionRoute.Get("/", routes.GetConnections)
	connectionRoute.Post("/", routes.CreateConnection)
	connectionRoute.Put("/:id", routes.EditConnection)
	connectionRoute.Delete("/:id", routes.RemoveConnection)
	connectionRoute.Get("/:id", routes.GetConnectionInfo)

	app.Get("/version", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"version": "1.0.0"})
	})
	log.Fatal(app.Listen(":3000"))

}
