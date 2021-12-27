package views

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	types "github.com/strick-j/go-form-webserver/types"
	"github.com/strick-j/go-form-webserver/utils"
)

// Generate Struct for the forms required by Safe Functions
// Form Data is used by several functions (add, get, update, etc...)
var safeFormData = types.CreateForm{
	FormAction: "/safeaddreq/",
	FormMethod: "POST",
	FormLegend: "Add Safe Form",
	FormFields: []types.FormFields{
		{
			FieldLabel:      "SafeName",
			FieldLabelText:  "Safe Name",
			FieldInputType:  "Text",
			FieldRequired:   true,
			FieldInputName:  "FormSafeName",
			FieldInFeedback: "Safe Name is Required.",
			FieldIdNum:      1,
		},
		{
			FieldLabel:     "DisplayName",
			FieldLabelText: "Safe Display Name",
			FieldInputType: "Text",
			FieldRequired:  false,
			FieldInputName: "FormSafeDisplayName",
			FieldDescBy:    "displayHelp",
			FieldHelp:      "Optional",
			FieldIdNum:     2,
		},
		{
			FieldLabel:     "SafeDescription",
			FieldLabelText: "Description",
			FieldInputType: "Text",
			FieldRequired:  false,
			FieldInputName: "FormSafeDescription",
			FieldDescBy:    "descHelp",
			FieldHelp:      "Optional",
			FieldIdNum:     3,
		},
	},
}

//SafeAllReq is the function for requesting user info for collecting data to add a new user
func SafeAllReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Retrieve byte based object via BuildUrl function
	log.Println("INFO SafeAllReq: Attempting to obtain Safe Data from SCIM API.")
	res, err := BuildUrl("Containers", "GET")
	if err != nil {
		log.Println("ERROR SafeAllReq:", err)
		return
	} else {
		log.Println("INFO SafeAllReq: Safe Information Recieved")
	}

	// Declare and unmarshal byte based response
	var bodyObject types.Safes
	json.Unmarshal(res, &bodyObject)

	// Creating a unique Safe Id for use in safe form creation
	for i := 0; i < len(bodyObject.Resources); i++ {
		bodyObject.Resources[i].UniqueSafeId = strconv.Itoa(i)
	}

	// Establish context for populating allinfo template
	context := types.Context{
		Navigation: "Safes",
		CreateForm: safeFormData,
		Safes:      bodyObject,
	}

	tpl.ExecuteTemplate(w, "objectallinfo.html", context)
}

// SafeAddForm is the form for collecting data to add a new Safe
func SafeAddForm(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO SafeAddForm: Initializing Add Safe Form")

	// Establish context for populating add safe template
	context := types.Context{
		Navigation: "Add Safe",
		CreateForm: safeFormData,
	}

	// Pass form data to form template to dynamically build form
	tpl.ExecuteTemplate(w, "objectaddform.html", context)
}

// SafeDelFunc is the function for deleting a Safe
// Safe deletions are based off of safe Name not ID like users or groups
func SafeDelFunc(w http.ResponseWriter, r *http.Request) {
	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := utils.GetRedirectUrl(r.Referer())

	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	log.Println("INFO SafeDelFunc: Starting User Delete Process")

	// Retrieve USerID from URL to send to Del Function
	safeName := r.URL.Path[len("/userdel/"):]

	if safeName == "" {
		log.Println("ERROR SafeDelFunc: Could not determine Safe Name for Deletion.")
		return
	} else {
		log.Println("INFO SafeDelFunc: Safe Name to Delete:", safeName)
	}

	// Create Struct for passing data to SCIM API Delete Function
	delObjectData := types.DelObjectRequest{
		ResourceType: "containers",
		ID:           safeName,
	}

	// Delete Safe and recieve response from Delete Function
	res, err := ScimApiDel(delObjectData)
	if res == 204 {
		log.Println("INFO SafeDelFunc: Safe Deleted:", safeName)
		log.Println("INFO SafeDelFunc: Valid HTTP StatusCode Recieved:", res)
		log.Println("SUCCESS SafeDelFunc: UserDelete Process Complete.")
	} else {
		log.Println(err)
		log.Println("ERROR SafeDelFunc: Invalid Http StatusCode Recieved:", res)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

//SafeAddReq is used to add users from the /useraddreq URL
func SafeAddReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := utils.GetRedirectUrl(r.Referer())

	log.Println("INFO SafeAddReq: Reading Data from Safe Add Form")
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

	// Required as placeholder
	blankstruct := types.PostUserRequest{}

	log.Println("INFO SafeAddReq: Sending POST to SCIM server for Safe addition.")
	res, code, err := ScimAPI("Containers", "POST", addSafeData, blankstruct)
	if code != 201 {
		log.Println("ERROR SafeAddReq: Error Adding User - Response StatusCode:", code)
		log.Println("ERROR SafeAddReq: Error Adding User", err)
		return
	} else {
		log.Println("INFO SafeAddReq: Recieved SCIM Response - Valid HTTP StatusCode:", code)
	}

	// Declare and unmarshal byte based response
	var bodyObject types.PostSafeResponse
	err = json.Unmarshal(res, &bodyObject)
	if err != nil {
		log.Println("ERROR SafeAddReq: ", err)
		return
	} else {
		log.Println("INFO SafeAddReq: Safe Display Name and ID: ", bodyObject.DisplayName, "-", bodyObject.ID)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
