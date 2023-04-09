package handlers

import (
	"net/http"
	"word-count/models"
	"word-count/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

const (
	ErrorTitle = "ระบบขัดข้อง"
	ErrorMsg   = "ไม่สามารถทำรายการได้ในขณะนี้ กรุณาลองใหม่ภายหลัง"
)

func GetSummary(c echo.Context) error {

	// Get text from external url
	url := viper.GetString("external.beef.url")
	if url == "" {
		log.Error("Fail to read url from config")
		return c.JSONPretty(http.StatusInternalServerError, models.ErrorResponse{ErrorTitle: ErrorTitle, ErrorMessage: ErrorMsg}, "")
	}

	log.Infof("Calling external url: %s", url)
	external := newExternal(url)
	text, err := external.getText()
	if err != nil {
		log.Error(err)
		return c.JSONPretty(http.StatusInternalServerError, models.ErrorResponse{ErrorTitle: ErrorTitle, ErrorMessage: ErrorMsg}, "")
	}
	log.Infof("Response from %s is %s", url, text)

	response := services.WordCount(text)

	return c.JSONPretty(http.StatusOK, response, "")
}
