package views

import (
	"html/template"
	"log"
	"net/http"

	"github.com/strick-j/scimplistic/types"
	"github.com/strick-j/scimplistic/utils"
)

var (
	tpl *template.Template
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func IndexReq(w http.ResponseWriter, r *http.Request) {
	values, err := utils.ReadConfig("settings.json")
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
