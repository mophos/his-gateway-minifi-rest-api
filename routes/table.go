package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/mophos/minifi-cli-go/models"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func CreateQueryTable(ctx *fiber.Ctx) error {

	connectionId := ctx.Params("id")

	data := new(models.CreateTableQueryStruct)

	if errParser := ctx.BodyParser(data); errParser != nil {
		log.Println(errParser.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "ไม่สามารถตรวจสอบ parameters ได้",
			"error":   errParser.Error(),
		})
	}

	errors := models.ValidateCreateTableQueryStruct(*data)
	if errors != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"error":   errors,
			"message": "ข้อมูลไม่ครบ กรุณาตรวจสอบ",
		})
	}

	var connectionPath = viper.GetString("data.connections")
	var connectionPathForCreate = filepath.Join(connectionPath, connectionId, "table_manual")

	// create connections directory
	errCreateConnectDir := os.MkdirAll(connectionPathForCreate, os.ModePerm)
	if errCreateConnectDir != nil {
		log.Println("Create connections directory: ", errCreateConnectDir.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"error":   errCreateConnectDir.Error(),
			"message": "ไม่สามารถสร้าง Directory connections ได้",
		})
	}

	var isErrorWrite bool
	var errorWriteJsonFileResponse []error

	for _, table := range data.Tables {
		tableFile := fmt.Sprintf("%s.json", filepath.Join(connectionPathForCreate, table.Name))
		// create file with json
		jsonData, _ := json.MarshalIndent(table, "", " ")
		errWriteJsonFile := ioutil.WriteFile(tableFile, jsonData, os.ModePerm)
		if errWriteJsonFile != nil {
			log.Println(errWriteJsonFile.Error())
			errorWriteJsonFileResponse = append(errorWriteJsonFileResponse, errWriteJsonFile)

			isErrorWrite = true
		}
	}

	if isErrorWrite {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"error":   errorWriteJsonFileResponse,
			"message": "ไม่สามารถสร้าง table json ได้",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{"ok": true})

}

func GetTableQueryStatus(ctx *fiber.Ctx) error {

	connectionId := ctx.Params("id")

	var dataPath = viper.GetString("dataPath")
	var connectionsPath = filepath.Join(dataPath, "connections")

	tableStatusPath := filepath.Join(connectionsPath, connectionId, "tables")
	// Read directory
	files, err := ioutil.ReadDir(tableStatusPath)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"error":   err,
			"message": "ไม่สามารถสร้าง table json ได้",
		})
	}

	// var jsonData []models.StatusFileJsonStruct

	var jsonTables []models.StatusFileJsonStruct

	for _, f := range files {
		log.Println(f.Name())
		if f.IsDir() {
			tablePath := filepath.Join(tableStatusPath, f.Name())

			tableFile, _ := ioutil.ReadDir(tablePath)

			for _, t := range tableFile {
				if !t.IsDir() {
					filename := filepath.Join(tablePath, t.Name())

					statusData, _ := ioutil.ReadFile(filename)

					var statusYaml models.StatusFileJsonStruct
					statusYaml.Table = f.Name()

					_ = yaml.Unmarshal([]byte(statusData), &statusYaml)
					jsonTables = append(jsonTables, statusYaml)
				}

			}

			return ctx.JSON(jsonTables)

		}
	}

	return ctx.Status(200).JSON(fiber.Map{"ok": true})

}

func GetTableQueryStatusInfo(ctx *fiber.Ctx) error {

	connectionId := ctx.Params("id")
	tableName := ctx.Params("table")

	var dataPath = viper.GetString("dataPath")
	var connectionsPath = filepath.Join(dataPath, "connections")

	tableStatusPath := filepath.Join(connectionsPath, connectionId, "tables", tableName)
	var jsonTables []models.StatusFileJsonStruct

	statusFiles, errReadDir := ioutil.ReadDir(tableStatusPath)

	if errReadDir != nil {
		log.Println(errReadDir)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"error":   errReadDir,
			"message": "ไม่สามารถอ่านไฟล์ได้",
		})
	}
	for _, t := range statusFiles {
		if !t.IsDir() {
			filename := filepath.Join(tableStatusPath, t.Name())
			statusData, _ := ioutil.ReadFile(filename)

			var statusYaml models.StatusFileJsonStruct
			statusYaml.Table = tableName

			_ = yaml.Unmarshal([]byte(statusData), &statusYaml)
			jsonTables = append(jsonTables, statusYaml)
		}

	}

	if len(jsonTables) == 0 {
		return ctx.JSON(fiber.Map{"ok": false, "message": "ไม่พบข้อมูล"})
	}

	return ctx.JSON(jsonTables)

}
