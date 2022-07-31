package models

import "github.com/go-playground/validator"

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

type ConfigValidateStruct struct {
	MaxConcurrentThreads int    `json:"max_concurrent_threads" validate:"required"`
	KeystorePath         string `json:"keystore_path" validate:"required"`
	KeystorePassword     string `json:"keystore_password" validate:"required"`
	TruststorePath       string `json:"truststore_path" validate:"required"`
	TruststorePassword   string `json:"truststore_password" validate:"required"`
}

type CreateConnectionValidateStruct struct {
	ConnectionName string `json:"connection_name" validate:"required"`
	Hospcode       string `json:"hospcode" validate:"required"`
	CronjobAll     struct {
		RunEvery int    `json:"run_every" validate:"required"`
		RunTime  string `json:"run_time" validate:"required"`
	} `json:"cronjob_all"`
	CronjobQuery struct {
		Dayago  int    `json:"dayago" validate:"required"`
		RunTime string `json:"run_time" validate:"required"`
	} `json:"cronjob_query"`
	Database struct {
		HisName  string `json:"his_name" validate:"required"`
		Host     string `json:"host" validate:"required"`
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
		Port     int32  `json:"port" validate:"required"`
		Schema   string `json:"schema"`
		Type     string `json:"type" validate:"required"`
		Username string `json:"username" validate:"required"`
	} `json:"database"`
	Broker struct {
		BootstrapServer string `json:"bootstrap_server" validate:"required"`
		Topic           string `json:"topic" validate:"required"`
	} `json:"broker"`
}

type CreateTableQueryStruct struct {
	Tables []struct {
		Name      string `json:"name" validate:"required"`
		StartDate string `json:"start_date" validate:"required"`
		EndDate   string `json:"end_date" validate:"required"`
	} `json:"tables"`
}

func ValidateCreateTableQueryStruct(data CreateTableQueryStruct) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func ValidateCreateConnectionStruct(data CreateConnectionValidateStruct) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func ValidateConfigStruct(data ConfigValidateStruct) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
