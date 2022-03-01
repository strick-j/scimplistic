package views

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

// Generate Struct for actions required by User functions
var safeObject = Object{
	Type: "containers",
}

///////////////////////// Safe Default Handler /////////////////////////

// SafesHandler is the function for displaying the basic Safes page.
func SafesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafesHandler"}).Info("Starting Safe retrieval Process")

	safeObject.Method = "GET"

	// Retrieve byte based object via BuildUrl function
	log.Println("INFO SafesHandler: Attempting to obtain Safe Data from SCIM API.")
	//res, err := BuildUrl("Containers", "GET")
	res, _, err := safeObject.ScimType2Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafesHandler"}).Error(err)
		return
	}

	for i := 0; i < len(res.Resources); i++ {
		res.Resources[i].UniqueSafeId = strconv.Itoa(i)
	}

	// Establish context for populating allinfo template
	context := types.Context{
		Navigation: "Safes",
		CreateForm: safeFormData,
		Safes:      *res,
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafesHandler"}).Info("Safe retrieval process completed  without error")

	tpl.ExecuteTemplate(w, "objectallinfo.html", context)
}

///////////////////////// Safe Action Handlers /////////////////////////

func SafeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// For best UX we want the user to be returned to the page making
	// the delete transaction, we use the r.Referer() function to get the link.
	redirectURL := GetRedirectUrl(r.Referer())

	// Parse Vars from URL Variables
	vars := mux.Vars(r)

	safeObject.Method = "GET"
	safeObject.Id = vars["id"]

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeHandler"}).Info("Derived Request Id: ", safeObject.Id)

	// Retrieve byte based object via BuildUrl function
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountHandler"}).Trace("Calling ScimType2Api to retrieve Safe info")
	_, res, err := safeObject.ScimType2Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeHandler"}).Error(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	fmt.Println(res)

	http.Redirect(w, r, redirectURL, http.StatusFound)
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

	/*	log.Println("INFO SafeAddReq: Reading Data from Safe Add Form")
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
		} */

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

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	log.Println("INFO SafeDelFunc: Starting Safe Delete Process")

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountDelHandler"}).Info("Starting Safe Delete Process")

	// Retrieve USerID from URL to send to Del Function
	vars := mux.Vars(r)
	safeObject.Id = vars["id"]
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeDelHandler"}).Trace("Safe Id: ", safeObject.Id)

	// Delete Safe and recieve response from Delete Function
	_, _, err := safeObject.ScimType2Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeDelHandler"}).Error("Account Delete Process completed finished Error")
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeDelHandler"}).Info("Account Delete Process completed finished Error")

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
