package views

import (
	"encoding/json"
	"fmt"
	"log"

	//"log"
	"net/http"

	//db "github.com/strick-j/go-form-webserver/db"
	types "github.com/strick-j/go-form-webserver/types"
)

//GroupAllReq is the function for requesting user info for collecting data to add a new user
func GroupAllReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Retrieve byte based object via BuildUrl function
	res := BuildUrl("Groups", "GET")

	// Declare and unmarshal byte based response
	var bodyObject types.Group
	json.Unmarshal(res, &bodyObject)

	tpl.ExecuteTemplate(w, "groupallinfo.html", bodyObject)
}

//GroupDelForm is the form for collecting data to add a new user
func GroupDelForm(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "groupdelform.html", nil)
}

//GroupAddForm is the form for deleting a user
func GroupAddForm(w http.ResponseWriter, r *http.Request) {

	log.Printf("Initializing Add Group Form")

	groupFormData := types.CreateForm{
		FormAction: "/groupaddreq",
		FormMethod: "GET",
		FormLegend: "Add Group Form",
		FormFields: []types.FormFields{
			{
				FieldLabel:     "DisplayName",
				FieldLabelText: "Group Display Name",
				FieldInputType: "Text",
				FieldRequired:  true,
				FieldInputName: "FormGroupDisplayName",
				FieldIdNum:     1,
			},
		},
	}

	// Pass form data to form template to dynamically build form
	tpl.ExecuteTemplate(w, "objectaddform.html", groupFormData)
}

//GroupAddReq is used to add users from the /useraddreq URL
func GroupAddReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Println("Reading Data from Group Add Form")
	displayName := r.FormValue("FormGroupDisplayName")
	scimschema := []string{"urn:ietf:params:scim:schemas:core:2.0:Group"}

	addGroupData := types.PostObjectRequest{
		DisplayName: displayName,
		Schemas:     scimschema,
	}

	res := ScimAPI("Groups", "POST", addGroupData)

	fmt.Println(string(res))
	tpl.ExecuteTemplate(w, "/", nil)
}
