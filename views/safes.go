package views

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	types "github.com/strick-j/go-form-webserver/types"
)

//SafeAllReq is the function for requesting user info for collecting data to add a new user
func SafeAllReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res := BuildUrl("Containers", "GET")

	var bodyObject types.Safes
	json.Unmarshal(res, &bodyObject)

	for i := 0; i < len(bodyObject.Resources); i++ {
		bodyObject.Resources[i].UniqueSafeId = strconv.Itoa(i)
	}

	tpl.ExecuteTemplate(w, "safeallinfo.html", bodyObject)
}

//UserAddForm is the form for collecting data to add a new user
func SafeAddForm(w http.ResponseWriter, r *http.Request) {

	log.Printf("Initializing Add Safe Form")

	safeFormData := types.CreateForm{
		FormAction: "/safeaddreq",
		FormMethod: "GET",
		FormLegend: "Add Safe Form",
		FormFields: []types.FormFields{
			{
				FieldLabel:     "SafeName",
				FieldLabelText: "Safe Name",
				FieldInputType: "Text",
				FieldRequired:  true,
				FieldInputName: "FormSafeName",
				FieldIdNum:     1,
			},
			{
				FieldLabel:     "DisplayName",
				FieldLabelText: "Group Display Name",
				FieldInputType: "Text",
				FieldRequired:  true,
				FieldInputName: "FormSafeDisplayName",
				FieldIdNum:     2,
			},
			{
				FieldLabel:     "SafeDescription",
				FieldLabelText: "Description",
				FieldInputType: "Text",
				FieldRequired:  true,
				FieldInputName: "FormSafeDescription",
				FieldIdNum:     3,
			},
		},
	}

	// Pass form data to form template to dynamically build form
	tpl.ExecuteTemplate(w, "objectaddform.html", safeFormData)
}

//UserDelForm is the form for deleting a user
func SafeDelForm(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "safedelform.html", nil)
}

//SafeAddReq is used to add users from the /useraddreq URL
func SafeAddReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Println("Reading Data from Safe Add Form")
	safeName := r.FormValue("FormSafeName")
	displayName := r.FormValue("FormSafeDisplayName")
	description := r.FormValue("FormSafeDescription")
	scimschema := []string{"urn:ietf:params:scim:schemas:pam:1.0:Container"}

	addSafeData := types.PostObjectRequest{
		Name:        safeName,
		DisplayName: displayName,
		Description: description,
		Schemas:     scimschema,
	}

	res := ScimAPI("Containers", "POST", addSafeData)

	fmt.Println(string(res))
	//tpl.ExecuteTemplate(w, "safeallinfo.html", bodyObject)
}
