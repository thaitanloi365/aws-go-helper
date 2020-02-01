package handlers

import (
	"aws-go-helper/config"
	"net/http"

	"github.com/labstack/echo"
)

// GetConfigHandler get config
func GetConfigHandler(c echo.Context) error {
	var response = config.Instance

	return c.JSON(http.StatusOK, response)
}

// GetPublicInfoHandler get info
func GetPublicInfoHandler(c echo.Context) error {
	var response = map[string]string{
		"name":    config.Instance.AppName,
		"version": config.Instance.Version,
	}

	return c.JSON(http.StatusOK, response)
}
