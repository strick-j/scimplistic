package views

import (
	"html/template"
	"log"
	"net/http"

	"github.com/strick-j/scimplistic/config"
	"github.com/strick-j/scimplistic/types"
)

var (
	tpl *template.Template
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func IndexReq(w http.ResponseWriter, r *http.Request) {
	values, err := config.ReadConfig("settings.json")
	if err != nil {
		log.Println("ERROR IndexReq:", err)
	}

	context := types.Context{
		Navigation: "",
	}

	if values.PrevConf {
		context.SettingsConfigured = true
	} else {
		context.SettingsConfigured = false
	}

	tpl.ExecuteTemplate(w, "index.html", context)
}
