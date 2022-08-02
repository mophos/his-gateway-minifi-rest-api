package routes

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mophos/minifi-cli-go/models"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func GenerateConfig(ctx *fiber.Ctx) error {

	var dataPath = viper.GetString("dataPath")
	var outPath = viper.GetString("outPath")
	var triggerFile = "/opt/minifi/data/update.txt"
	var connectionsPath = filepath.Join(dataPath, "connections")

	var settingFilePath = viper.GetString("settingFile")

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

	//Template generate
	templateDir := viper.GetString("templatePath")
	tmpDir := filepath.Join(dataPath, "tmp")
	//file main.yml
	mainFlowFile := filepath.Join(templateDir, "main.yml")

	// Create tmp directory
	err := os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		log.Println("Create tmp directory: ", err.Error())
	}

	// Create tmp directory
	errCreateConnectionDir := os.MkdirAll(connectionsPath, os.ModePerm)
	if errCreateConnectionDir != nil {
		log.Println("Create connection directory: ", errCreateConnectionDir.Error())
	}

	mainFlowFileData := models.MainFlowTemplateDataStruct{
		MAXCONCURRENTTHREADS: configYaml.Server.MaxConcurrentThreads,
		KEYSTORE_PATH:        configYaml.Server.KeystorePath,
		KEYSTORE_PASSWORD:    configYaml.Server.KeystorePassword,
		TRUSTSTORE_PATH:      configYaml.Server.TruststorePath,
		TRUSTSTORE_PASSWORD:  configYaml.Server.TruststorePassword,
	}

	mainFlowFileTemplate, errParseMainFlowFile := template.ParseFiles(mainFlowFile)

	if errParseMainFlowFile != nil {
		log.Println("Read main.yml template: ", errParseMainFlowFile.Error())
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "message": "ไม่สามารถอ่านไฟล์ main.yml ได้", "error": errParseMainFlowFile.Error()})
	}

	// tmp file
	tmpMainFlowFilePath := filepath.Join(tmpDir, "main.yml")
	tmpMainFlowFile, errWriteTmpMainFlowFile := os.Create(tmpMainFlowFilePath)

	if errWriteTmpMainFlowFile != nil {
		log.Println("create file: ", errWriteTmpMainFlowFile.Error())
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "message": "ไม่สามารถสร้างไฟล์ tmp/main.yml ได้", "error": errWriteTmpMainFlowFile.Error()})
	}

	errWriteMainFlowFileTemplate := mainFlowFileTemplate.Execute(tmpMainFlowFile, mainFlowFileData)

	if errWriteMainFlowFileTemplate != nil {
		log.Println("create file: ", errWriteMainFlowFileTemplate.Error())
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "message": "ไม่สามารถแก้ไขข้อมูล template tmp/main.yml ได้", "error": errWriteMainFlowFileTemplate.Error()})
	}

	// Get all conections
	connections := configYaml.Connections

	if len(connections) == 0 {
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "message": "ไม่พบ Connection"})
	}

	// read file main.yml
	mainFlowTmpPath, err := ioutil.ReadFile(tmpMainFlowFilePath)

	if err != nil {
		log.Println(err)
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถอ่านไฟล์ tmp/main.yml ได้"})
	}

	var mainFlowData models.FlowMainStruct

	errFlowMainYaml := yaml.Unmarshal([]byte(mainFlowTmpPath), &mainFlowData)
	if errFlowMainYaml != nil {
		log.Println(errFlowMainYaml.Error())
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถ decode ค่าในไฟล์ main.yml"})
	}

	// connection loop
	for _, conn := range connections {

		connectionName := conn.Name
		connectionId := conn.ID
		databaseType := conn.Type

		connectionTmpDir := filepath.Join(tmpDir, connectionName)
		// configure connection.yml/flow.yml
		err := os.MkdirAll(connectionTmpDir, os.ModePerm)
		if err != nil {
			log.Println("Create tmp directory: ", err.Error())
			return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถ tmp directory ได้"})
		}

		// template path
		connectionFilePath := filepath.Join(templateDir, databaseType, "connection.yml")
		flowFilePath := filepath.Join(templateDir, databaseType, "flow.yml")

		// parse template file
		connectionFileTemplate, errParseConnectionFileTemplate := template.ParseFiles(connectionFilePath)

		if errParseConnectionFileTemplate != nil {
			log.Println("Read connection.yml file: ", errParseConnectionFileTemplate.Error())
			return ctx.Status(500).JSON(fiber.Map{"ok": false, "message": "ไม่สามารถอ่านไฟล์ connection.yml ได้", "error": errParseConnectionFileTemplate.Error()})
		}

		flowFileTemplate, errParseFlowFileTemplate := template.ParseFiles(flowFilePath)

		if errParseFlowFileTemplate != nil {
			log.Println("Read connection.yml file: ", errParseFlowFileTemplate.Error())
			return ctx.Status(500).JSON(fiber.Map{"ok": false, "message": "ไม่สามารถอ่านไฟล์ flow.yml ได้", "error": errParseFlowFileTemplate.Error()})
		}

		connData := models.ConnectoionTemplateStruct{
			CONNECTION_UUID: connectionId,
			CONNECTION_NAME: connectionName,
			HOST:            conn.Host,
			PORT:            conn.Port,
			NAME:            conn.Database,
			USERNAME:        conn.Username,
			PASSWORD:        conn.Password,
		}

		cronQuery := strings.Split(conn.Cronjob.CronjobQuery.RunTime, ":")
		queryCronTab := fmt.Sprintf("%s %s * * *", cronQuery[1], cronQuery[0])

		cronAll := strings.Split(conn.Cronjob.CronjobAll.RunTime, ":")
		allCronTab := fmt.Sprintf("%s %s * * *", cronAll[1], cronAll[0])

		manualPath := filepath.Join(dataPath, connectionsPath, connectionId, "table_manual")

		flowID := uuid.NewString()
		flowData := models.FlowTemplateStruct{
			FLOW_UUID:        flowID,
			CONNECTION_UUID:  connectionId,
			CONNECTION_NAME:  connectionName,
			MANUAL_PATH:      manualPath,
			HIS_NAME:         conn.HisName,
			DAYAGO:           conn.Cronjob.Dayago,
			TOPIC:            conn.Broker.Topic,
			HOSPCODE:         conn.Hospcode,
			BOOTSTRAP_SERVER: conn.Broker.BootstrapServer,
			CRON_QUERY:       queryCronTab,
			CRON_ALL:         allCronTab,
			CONNECTION_PATH:  connectionsPath,
		}

		// tmp connection file
		tmpConnectionFlowFilePath := filepath.Join(connectionTmpDir, "connection.yml")
		tmpConnectionFlowFile, errWriteTmpConnectionFlowFile := os.Create(tmpConnectionFlowFilePath)

		if errWriteTmpConnectionFlowFile != nil {
			log.Println("apply template data (connection.yml): ", errWriteTmpConnectionFlowFile.Error())
			return ctx.Status(500).JSON(fiber.Map{"ok": false, "message": "ไม่สามารถ apply template data (connection.yml)", "error": errWriteTmpConnectionFlowFile.Error()})
		}

		errWriteConnectionFlowFileTemplate := connectionFileTemplate.Execute(tmpConnectionFlowFile, connData)

		if errWriteConnectionFlowFileTemplate != nil {
			log.Println("create file: ", errWriteConnectionFlowFileTemplate.Error())
			return ctx.Status(500).JSON(fiber.Map{"ok": false, "message": "ไม่สามารถแก้ไขข้อมูล template connnection.yml ได้", "error": errWriteConnectionFlowFileTemplate.Error()})
		}

		// tmp flow file
		tmpFlowFilePath := filepath.Join(connectionTmpDir, "flow.yml")
		tmpFlowFile, errWriteTmpFlowFile := os.Create(tmpFlowFilePath)

		if errWriteTmpFlowFile != nil {
			log.Println("apply template data (flow.yml): ", errWriteTmpFlowFile.Error())
			return ctx.Status(500).JSON(fiber.Map{"ok": false, "message": "ไม่สามารถ apply template data (flow.yml)", "error": errWriteTmpFlowFile.Error()})
		}

		errWriteFlowFileTemplate := flowFileTemplate.Execute(tmpFlowFile, flowData)

		if errWriteFlowFileTemplate != nil {
			log.Println("create file: ", errWriteFlowFileTemplate.Error())
			return ctx.Status(500).JSON(fiber.Map{"ok": false, "message": "ไม่สามารถแก้ไขข้อมูล template connnection.yml ได้", "error": errWriteFlowFileTemplate.Error()})
		}

		// read connection/flow file

		var _connData models.ControllerServiceStruct
		var _flowData models.FlowStruct

		connYamlData, errReadConnection := ioutil.ReadFile(tmpConnectionFlowFilePath)

		flowYamlData, errReadFlow := ioutil.ReadFile(tmpFlowFilePath)

		if errReadConnection != nil {
			log.Println(errReadConnection)
			return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถอ่านไฟล์ tmp/connection.yml ได้"})
		}

		if errReadFlow != nil {
			log.Println(errReadFlow)
			return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถอ่านไฟล์ tmp/flow.yml ได้"})
		}

		errConnsYaml := yaml.Unmarshal([]byte(connYamlData), &_connData)
		if errConnYaml != nil {
			log.Println(errConnsYaml.Error())
			return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถ encode ค่าสำหรับไฟล์ connection.yml"})
		}

		errFlowsYaml := yaml.Unmarshal([]byte(flowYamlData), &_flowData)
		if errFlowsYaml != nil {
			log.Println(errFlowsYaml)
			return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": errFlowsYaml, "message": "ไม่สามารถ encode ค่าสำหรับไฟล์ flow.yml"})
		}

		mainFlowData.ControllerServices = append(mainFlowData.ControllerServices, _connData)
		mainFlowData.ProcessGroups = append(mainFlowData.ProcessGroups, _flowData)

	} // end connection loop

	// create file
	yamlData, errMarshal := yaml.Marshal(&mainFlowData)
	if errMarshal != nil {
		log.Println(errMarshal.Error())
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": errMarshal.Error()})
	}

	configPath := filepath.Join(outPath, "config.yml")
	errWriteFile := ioutil.WriteFile(configPath, yamlData, os.ModePerm)
	if errWriteFile != nil {
		log.Println(errWriteFile.Error())
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": errWriteFile.Error()})
	}

	// remove all tmp
	errSettingfile := os.RemoveAll(tmpDir)
	if errSettingfile != nil {
		log.Println(errSettingfile)
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "error": "ไม่สามารถลบโฟลเดอร์ tmp ได้"})
	}

	// trigger minifi reload configure

	f, errCreateStatusFile := os.Create(triggerFile)

	if errCreateStatusFile != nil {
		log.Println(errCreateStatusFile)
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "message": "ไม่สามารถสร้างไฟล์ update.txt ได้", "error": errCreateStatusFile.Error()})
	}

	defer f.Close()

	updateTime := time.Now().Unix()
	txtMessage := fmt.Sprintf("%v\n", updateTime)
	_, errWriteUpdate := f.WriteString(txtMessage)

	if errWriteUpdate != nil {
		log.Println(errWriteUpdate)
		return ctx.Status(500).JSON(fiber.Map{"ok": false, "message": "ไม่สามารถเขียนไฟล์ update.txt ได้", "error": errWriteUpdate.Error()})
	}

	return ctx.Status(201).JSON(fiber.Map{"ok": true})

}
