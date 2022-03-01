package views

import (
	"errors"
	"log"
	"os"
	"path"
	"strings"

	"github.com/strick-j/scimplistic/utils"
)

func GetRedirectUrl(referer string) string {
	var redirectUrl string
	url := strings.Split(referer, "/")

	if len(url) > 4 {
		redirectUrl = "/" + strings.Join(url[3:], "/")
	} else {
		redirectUrl = "/"
	}
	return redirectUrl
}

func GetRedirectUrlNoId(referer string) string {
	var redirectUrl string
	url := strings.Split(referer, "/")

	if len(url) > 4 {
		redirectUrl = "/" + strings.Join(url[3:], "/")
	} else {
		redirectUrl = "/"
	}

	redirectUrlNoId := path.Dir(redirectUrl)

	return redirectUrlNoId
}

func CheckTLS(CertFile string, PrivKeyFile string) (bool, error) {
	log.Println("INFO CheckTLS: Enable TLS selected, checking for required cert and private key")
	if _, err := os.Open(CertFile); errors.Is(err, os.ErrNotExist) {
		//if errors.Is(err, os.ErrNotExist) {
		log.Println("ERROR CheckTLS: Enable TLS was selected, however certificate was not found. Disabling TLS")
		return false, err
	} else if _, err := os.Open(PrivKeyFile); errors.Is(err, os.ErrNotExist) {
		log.Println("ERROR CheckTLS: Enable TLS was selected, however Private Key was not found. Disabling TLS")
		return false, err
	}
	return true, nil
}

func ReturnScimUrl() *string {
	values, err := utils.ReadConfig("settings.json")
	if err != nil {
		log.Println(err)
	}

	return &values.ScimURL
}
