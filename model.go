package main

type FlowStatusStruct struct {
	ControllerServiceStatusList interface{} `json:"controllerServiceStatusList"`
	ProcessorStatusList         []struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		ProcessorHealth struct {
			RunStatus           string        `json:"runStatus"`
			HasBulletins        bool          `json:"hasBulletins"`
			ValidationErrorList []interface{} `json:"validationErrorList"`
		} `json:"processorHealth"`
		ProcessorStats struct {
			ActiveThreads     int `json:"activeThreads"`
			FlowfilesReceived int `json:"flowfilesReceived"`
			BytesRead         int `json:"bytesRead"`
			BytesWritten      int `json:"bytesWritten"`
			FlowfilesSent     int `json:"flowfilesSent"`
			Invocations       int `json:"invocations"`
			ProcessingNanos   int `json:"processingNanos"`
		} `json:"processorStats"`
		BulletinList []interface{} `json:"bulletinList"`
	} `json:"processorStatusList"`
	ConnectionStatusList         interface{}   `json:"connectionStatusList"`
	RemoteProcessGroupStatusList interface{}   `json:"remoteProcessGroupStatusList"`
	InstanceStatus               interface{}   `json:"instanceStatus"`
	SystemDiagnosticsStatus      interface{}   `json:"systemDiagnosticsStatus"`
	ReportingTaskStatusList      interface{}   `json:"reportingTaskStatusList"`
	ErrorsGeneratingReport       []interface{} `json:"errorsGeneratingReport"`
}
