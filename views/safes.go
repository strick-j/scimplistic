package views

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/strick-j/scimplistic/types"
)

// Generate Struct for the forms required by Safe Functions
// Form Data is used by several functions (add, get, update, etc...)
var safeFormData = types.CreateForm{
	FormAction: "/safes/add",
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

///////////////////////// Safe Default Handler /////////////////////////

// SafesHandler is the function for displaying the basic Safes page.
func SafesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Retrieve byte based object via BuildUrl function
	log.Println("INFO SafesHandler: Attempting to obtain Safe Data from SCIM API.")
	res, err := BuildUrl("Containers", "GET")
	if err != nil {
		log.Println("ERROR SafesHandler:", err)
		return
	} else {
		log.Println("INFO SafesHandler: Safe Information Recieved")
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

///////////////////////// Safe Action Handlers /////////////////////////

// SafesAction Handler will decide what is required based on the action provided.
// add: Proceed to SafeAddHandler
// del: Proceed to SafeDelHandler
// update: Proceed to SafeUpdateHandler
// review: Proceed to SafeReviewHandler
func SafesActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// Extract action from URL using mux.Vars.
	vars := mux.Vars(r)
	action := vars["action"]

	// Switch to appropriate handler based on action type.
	switch action {
	case "add":
		log.Println("INFO SafesActionHandler: Calling SafeAddHandler")
		SafeAddHandler(w, r)
	case "del":
		log.Println("INFO SafesActionHandler: Calling SafeDelHandler")
		SafeDelHandler(w, r)
	case "update":
		log.Println("INFO SafesActionHandler: Calling SafeUpdateHandler")
		SafeUpdateHandler(w, r)
	case "review":
		log.Println("INFO SafesActionHandler: Calling SafeReviewHandler")
		SafeUpdateHandler(w, r)
	}
}

// SafeAddHandler reads in form data from the modal that is triggered
// when the add safe button is pressed. This function calls the SCIM
// function which submits the Delete action to the SCIM Endpoint.
func SafeAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := GetRedirectUrl(r.Referer())

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
		log.Println("ERROR SafeAddReq: Error Adding Safe - Response StatusCode:", code)
		log.Println("ERROR SafeAddReq: Error Adding Safe", err)
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

// SafeDelHandler reads in form data from the modal that is triggered
// when the delete button for a particular safe is pressed. This function
// calls the SCIM function which submits the Delete action to the SCIM Endpoint.
// Note: Safe deletions are based off of safe Name not ID like users or groups
func SafeDelHandler(w http.ResponseWriter, r *http.Request) {
	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := GetRedirectUrl(r.Referer())

	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	log.Println("INFO SafeDelFunc: Starting Safe Delete Process")

	// Retrieve USerID from URL to send to Del Function
	vars := mux.Vars(r)
	safeName := vars["id"]
	log.Println("INFO SafeDelFunc: Safe ID to Delete:", safeName)

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
		log.Println("SUCCESS SafeDelFunc: Safe Delete Process Complete.")
	} else {
		log.Println(err)
		log.Println("ERROR SafeDelFunc: Invalid Http StatusCode Recieved:", res)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// SafeUpdateHandler reads in form data from the modal that is triggered
// when the Update button for a particular safe is pressed. This function
// calls the SCIM function which submits the UPDATE action to the SCIM Endpoint.
func SafeUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Read information from Update Form

	// TODO: Parse info and form PUT/POST for Group Update

	// TODO: Execute PUT/POST

	// TODO: Return to Groups page after success
	tpl.ExecuteTemplate(w, "/", nil)
}

// SafeReviewHandler reads in safe id information when the Review button
// for a particular safe is pressed. This function then provides the user with
// more information about a particular group.
func SafeReviewHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Read information from Update Form

	// TODO: Parse info and form PUT/POST for Group Update

	// TODO: Execute PUT/POST

	// TODO: Return to Groups page after success
	tpl.ExecuteTemplate(w, "/", nil)
}

///////////////////////// Safe Form Handlers /////////////////////////

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
