package routes

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mophos/minifi-cli-go/models"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func UpdateConfig(ctx *fiber.Ctx) error {

	data := new(models.ConfigValidateStruct)

	if errParser := ctx.BodyParser(data); errParser != nil {
		log.Println(errParser.Error())
		return ctx.Status(200).JSON(fiber.Map{
			"ok":    false,
			"error": errParser.Error(),
		})
	}

	errors := models.ValidateConfigStruct(*data)
	if errors != nil {
		return ctx.Status(200).JSON(fiber.Map{
			"ok":    false,
			"error": "ข้อมูลไม่ครบ กรุณาตรวจสอบ",
		})
	}

	var settingFilePath = viper.GetString("settingFile")

	confData, errReadYaml := ioutil.ReadFile(settingFilePath)

	if errReadYaml != nil {
		return ctx.Status(200).JSON(fiber.Map{"ok": false, "error": "Configure file not found."})
	}

	var configYaml models.SettingStruct

	errConnYaml := yaml.Unmarshal([]byte(confData), &configYaml)
	if errConnYaml != nil {
		log.Println(errConnYaml)
		return ctx.Status(200).JSON(fiber.Map{"ok": false, "error": errConnYaml.Error()})
	}

	configYaml.Server.MaxConcurrentThreads = data.MaxConcurrentThreads
	configYaml.Server.KeystorePath = data.KeystorePath
	configYaml.Server.KeystorePassword = data.KeystorePassword
	configYaml.Server.TruststorePath = data.TruststorePath
	configYaml.Server.TruststorePassword = data.TruststorePassword

	// create file
	yamlData, errMarshal := yaml.Marshal(&configYaml)
	if errMarshal != nil {
		log.Println(errMarshal.Error())
		return ctx.Status(200).JSON(fiber.Map{"ok": false, "error": errMarshal.Error()})
	}

	errWriteFile := ioutil.WriteFile(settingFilePath, yamlData, os.ModePerm)
	if errWriteFile != nil {
		log.Println(errWriteFile.Error())
		return ctx.Status(200).JSON(fiber.Map{"ok": false, "error": errWriteFile.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{"ok": true})

}

func GetConfig(ctx *fiber.Ctx) error {
	var settingFile = viper.GetString("settingFile")

	confData, errReadYaml := ioutil.ReadFile(settingFile)

	if errReadYaml != nil {
		return ctx.Status(200).JSON(fiber.Map{"ok": false, "error": "Configure file not found."})
	}

	var configYaml models.SettingStruct

	errConnYaml := yaml.Unmarshal([]byte(confData), &configYaml)
	if errConnYaml != nil {
		log.Println(errConnYaml)
		return ctx.Status(200).JSON(fiber.Map{"ok": false, "error": errConnYaml.Error()})
	}

	setting := configYaml.Server

	return ctx.JSON(setting)
}
