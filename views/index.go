package views

import (
	"html/template"
	"log"
	"net/http"

	config "github.com/strick-j/go-form-webserver/config"
	types "github.com/strick-j/go-form-webserver/types"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func IndexReq(w http.ResponseWriter, r *http.Request) {
	values, err := config.ReadConfig("config.json")
	if err != nil {
		log.Println("ERROR IndexReq:", err)
	}

	context := types.Context{
		Navigation: "",
	}

	if values.ScimURL != "" {
		context.SettingsConfigured = true
	} else {
		context.SettingsConfigured = false
	}

	tpl.ExecuteTemplate(w, "index.html", context)
}
