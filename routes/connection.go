package routes

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mophos/minifi-cli-go/models"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func GetConnections(ctx *fiber.Ctx) error {

	var dataPath = viper.GetString("data.path")
	var settingFilePath = filepath.Join(dataPath, "data/config", "setting.yml")

	confData, errReadYaml := ioutil.ReadFile(settingFilePath)

	if errReadYaml != nil {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "Configure file not found."})
	}

	var configYaml models.SettingStruct

	errConnYaml := yaml.Unmarshal([]byte(confData), &configYaml)
	if errConnYaml != nil {
		log.Println(errConnYaml)
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": errConnYaml.Error()})
	}

	connections := configYaml.Connections

	return ctx.JSON(connections)
}

func GetConnectionInfo(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")

	if len(id) == 0 {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "Connection ID not found."})
	}

	var dataPath = viper.GetString("data.path")
	var settingFilePath = filepath.Join(dataPath, "data/config", "setting.yml")

	confData, errReadYaml := ioutil.ReadFile(settingFilePath)

	if errReadYaml != nil {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": errReadYaml.Error()})
	}

	var configYaml models.SettingStruct

	errConnYaml := yaml.Unmarshal([]byte(confData), &configYaml)
	if errConnYaml != nil {
		log.Fatalf("cannot unmarshal data: %v", errConnYaml)
	}

	connections := configYaml.Connections

	var info models.SettingConnectionStruct

	for _, item := range connections {

		if item.ID == id {
			info = item
		}
	}

	if len(info.ID) > 0 {
		return ctx.JSON(info)
	}

	return ctx.JSON(fiber.Map{"ok": false, "message": "ไม่พบข้อมูลการเชื่อมต่อ"})
}

func CreateConnection(ctx *fiber.Ctx) error {

	data := new(models.CreateConnectionValidateStruct)

	if errParser := ctx.BodyParser(data); errParser != nil {
		log.Println(errParser.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": errParser.Error(),
		})
	}

	errors := models.ValidateCreateConnectionStruct(*data)
	if errors != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": errors,
		})
	}

	var dataPath = viper.GetString("data.path")

	connectionName := data.ConnectionName

	var settingFilePath = filepath.Join(dataPath, "data/config", "setting.yml")

	confData, errReadYaml := ioutil.ReadFile(settingFilePath)

	if errReadYaml != nil {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": errReadYaml.Error()})
	}

	var configYaml models.SettingStruct

	errConnYaml := yaml.Unmarshal([]byte(confData), &configYaml)
	if errConnYaml != nil {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": errConnYaml.Error()})
	}

	connections := configYaml.Connections

	var connectionExist bool

	for _, item := range connections {

		if item.Name == connectionName {
			connectionExist = true
		}
	}

	if connectionExist {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "มี Connection ชื่อนี้แล้ว"})
	}

	// update config

	confSettingData, errReadSettingYaml := ioutil.ReadFile(settingFilePath)

	if errReadSettingYaml != nil {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่พบไฟล์ setting.yml"})
	}

	var configSettingYaml models.SettingStruct

	errConnfigYaml := yaml.Unmarshal([]byte(confSettingData), &configSettingYaml)
	if errConnfigYaml != nil {
		log.Println(errConnfigYaml)
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถ decode ข้อมูลสำหรับไฟล์ setting.yml ได้"})
	}

	var setting models.SettingConnectionStruct

	setting.ID = uuid.New().String()
	setting.Name = connectionName
	setting.Hospcode = data.Hospcode
	setting.Type = data.Database.Type
	setting.HisName = data.Database.HisName
	setting.Host = data.Database.Host
	setting.Port = data.Database.Port
	setting.Database = data.Database.Name
	setting.Schema = data.Database.Schema
	setting.Username = data.Database.Username
	setting.Password = data.Database.Password
	setting.Cronjob.CronjobQuery.Dayago = data.CronjobQuery.Dayago
	setting.Cronjob.CronjobQuery.RunTime = data.CronjobQuery.RunTime
	setting.Cronjob.CronjobAll.RunEvery = data.CronjobAll.RunEvery
	setting.Cronjob.CronjobAll.RunTime = data.CronjobAll.RunTime
	setting.Broker.BootstrapServer = data.Broker.BootstrapServer
	setting.Broker.Topic = data.Broker.Topic

	configSettingYaml.Connections = append(configSettingYaml.Connections, setting)

	yamlSetting, errSettingMarshal := yaml.Marshal(&configSettingYaml)
	if errSettingMarshal != nil {
		log.Println(errSettingMarshal.Error())
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถสร้างข้อมูล yaml สำหรับไฟล์  setting.yml ได้"})
	}

	// create setting file
	errWriteSettingFile := ioutil.WriteFile(settingFilePath, yamlSetting, os.ModePerm)
	if errWriteSettingFile != nil {
		log.Println(errWriteSettingFile.Error())
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถสร้างไฟล์ setting.yml ได้"})
	}

	return ctx.Status(201).JSON(fiber.Map{"ok": true})

}

func EditConnection(ctx *fiber.Ctx) error {

	id := ctx.Params("id", "")

	data := new(models.CreateConnectionValidateStruct)

	if errParser := ctx.BodyParser(data); errParser != nil {
		log.Println(errParser.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": errParser.Error(),
		})
	}

	errors := models.ValidateCreateConnectionStruct(*data)
	if errors != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": errors,
		})
	}

	var dataPath = viper.GetString("data.path")

	connectionName := data.ConnectionName

	var settingFilePath = filepath.Join(dataPath, "data/config", "setting.yml")

	confData, errReadYaml := ioutil.ReadFile(settingFilePath)

	if errReadYaml != nil {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": errReadYaml.Error()})
	}

	var configYaml models.SettingStruct

	errConnYaml := yaml.Unmarshal([]byte(confData), &configYaml)
	if errConnYaml != nil {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": errConnYaml.Error()})
	}

	connections := configYaml.Connections

	var isDuplicated bool

	for _, item := range connections {
		if item.Name == connectionName && item.ID != id {
			isDuplicated = true
		}
	}

	if isDuplicated {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "มี Connection ชื่อนี้แล้ว"})
	}

	// update config

	confSettingData, errReadSettingYaml := ioutil.ReadFile(settingFilePath)

	if errReadSettingYaml != nil {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่พบไฟล์ setting.yml"})
	}

	var configSettingYaml models.SettingStruct

	errConnfigYaml := yaml.Unmarshal([]byte(confSettingData), &configSettingYaml)
	if errConnfigYaml != nil {
		log.Println(errConnfigYaml)
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถ decode ข้อมูลสำหรับไฟล์ setting.yml ได้"})
	}

	var setting models.SettingConnectionStruct

	setting.ID = id
	setting.Name = connectionName
	setting.Hospcode = data.Hospcode
	setting.Type = data.Database.Type
	setting.HisName = data.Database.HisName
	setting.Host = data.Database.Host
	setting.Port = data.Database.Port
	setting.Database = data.Database.Name
	setting.Schema = data.Database.Schema
	setting.Username = data.Database.Username
	setting.Password = data.Database.Password
	setting.Cronjob.CronjobQuery.Dayago = data.CronjobQuery.Dayago
	setting.Cronjob.CronjobQuery.RunTime = data.CronjobQuery.RunTime
	setting.Cronjob.CronjobAll.RunEvery = data.CronjobAll.RunEvery
	setting.Cronjob.CronjobAll.RunTime = data.CronjobAll.RunTime
	setting.Broker.BootstrapServer = data.Broker.BootstrapServer
	setting.Broker.Topic = data.Broker.Topic

	// remove old

	for idx, item := range configSettingYaml.Connections {
		if item.ID == id {
			configSettingYaml.Connections = append(configSettingYaml.Connections[:idx], configSettingYaml.Connections[idx+1:]...)
		}
	}

	// add new
	configSettingYaml.Connections = append(configSettingYaml.Connections, setting)

	yamlSetting, errSettingMarshal := yaml.Marshal(&configSettingYaml)
	if errSettingMarshal != nil {
		log.Println(errSettingMarshal.Error())
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถสร้างข้อมูล yaml สำหรับไฟล์  setting.yml ได้"})
	}

	// create setting file
	errWriteSettingFile := ioutil.WriteFile(settingFilePath, yamlSetting, os.ModePerm)
	if errWriteSettingFile != nil {
		log.Println(errWriteSettingFile.Error())
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถสร้างไฟล์ setting.yml ได้"})
	}

	return ctx.Status(200).JSON(setting)

}

func RemoveConnection(ctx *fiber.Ctx) error {

	id := ctx.Params("id")

	var dataPath = viper.GetString("data.path")

	var settingFilePath = filepath.Join(dataPath, "data/config", "setting.yml")

	confData, errReadYaml := ioutil.ReadFile(settingFilePath)

	if errReadYaml != nil {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": errReadYaml.Error()})
	}

	var configYaml models.SettingStruct

	errConnYaml := yaml.Unmarshal([]byte(confData), &configYaml)
	if errConnYaml != nil {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": errConnYaml.Error()})
	}

	confSettingData, errReadSettingYaml := ioutil.ReadFile(settingFilePath)

	if errReadSettingYaml != nil {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่พบไฟล์ setting.yml"})
	}

	var configSettingYaml models.SettingStruct

	errConnfigYaml := yaml.Unmarshal([]byte(confSettingData), &configSettingYaml)
	if errConnfigYaml != nil {
		log.Println(errConnfigYaml)
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถ decode ข้อมูลสำหรับไฟล์ setting.yml ได้"})
	}

	for idx, item := range configSettingYaml.Connections {
		if item.ID == id {
			configSettingYaml.Connections = append(configSettingYaml.Connections[:idx], configSettingYaml.Connections[idx+1:]...)
		}
	}

	// save file
	yamlSetting, errSettingMarshal := yaml.Marshal(&configSettingYaml)
	if errSettingMarshal != nil {
		log.Println(errSettingMarshal.Error())
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถสร้างข้อมูล yaml สำหรับไฟล์  setting.yml ได้"})
	}

	// create setting file
	errWriteSettingFile := ioutil.WriteFile(settingFilePath, yamlSetting, os.ModePerm)
	if errWriteSettingFile != nil {
		log.Println(errWriteSettingFile.Error())
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถสร้างไฟล์ setting.yml ได้"})
	}

	return ctx.Status(200).JSON(fiber.Map{"ok": true})
}
