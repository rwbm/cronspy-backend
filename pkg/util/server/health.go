package server

import (
	"net/http"
	"os"

	"github.com/astropay/go-tools/files"
	"github.com/labstack/echo/v4"
)

// Local vars
var (
	AppVersion          string
	HostName            string
	AppName             = "Mercury API"
	ComponentName       = "mercury-api"
	MaintenanceFilePath = "/etc/app-mode/maintenance"
	StatusOperational   = "Operational"
	StatusMaintenance   = "Maintenance"
)

type appStatus struct {
	Status    string `json:"status"`
	Component string `json:"component"`
	Version   string `json:"version"`
	Server    string `json:"server"`
}

func healthCheckHandler(c echo.Context) error {

	responseCode := http.StatusOK
	if files.Exists(MaintenanceFilePath) {
		responseCode = http.StatusFound
	}

	// load application info
	status := getStatus(responseCode)
	return c.JSON(responseCode, status)
}

func getStatus(responseCode int) *appStatus {
	status := new(appStatus)
	status.Status = StatusOperational

	if responseCode != http.StatusOK {
		status.Status = StatusMaintenance
	}

	status.Component = AppName
	status.Version = AppVersion

	// get server's name
	if host, err := os.Hostname(); err == nil {
		status.Server = host
	}

	return status
}
