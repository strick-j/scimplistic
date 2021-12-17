package views

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	config "github.com/strick-j/go-form-webserver/config"
	"github.com/strick-j/go-form-webserver/types"
)

func BuildUrl(target string, apimethod string) []byte {
	values, err := config.ReadConfig("config.json")

	if err != nil {
		fmt.Println(err)
	}

	finu := &url.URL{
		Scheme:      "https",
		Opaque:      "",
		Host:        values.ScimUrl,
		Path:        "/scim/" + target,
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}

	// Pull API Method from function call (Get, Patch, Del, etc..)
	method := apimethod

	finurl := finu.String()

	client := &http.Client{}
	req, err := http.NewRequest(method, finurl, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", "Bearer "+values.AuthToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	return body
}

func ScimAPI(target string, apimethod string, data types.PostObjectRequest, userdata types.PostUserRequest) []byte {
	values, err := config.ReadConfig("config.json")
	if err != nil {
		fmt.Println(err)
	}

	// Generate json from pass struct
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		//return
	}

	// Generate target URL from passed target
	finu := &url.URL{
		Scheme:      "https",
		Opaque:      "",
		Host:        values.ScimUrl,
		Path:        "/scim/" + target,
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}

	// Pull API Method from function call (Get, Patch, Del, etc..)
	method := apimethod
	finurl := finu.String()
	client := &http.Client{}
	// Build Request
	req, err := http.NewRequest(method, finurl, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
	}

	// Add authorization for request
	req.Header.Add("Authorization", "Bearer "+values.AuthToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	fmt.Println(res)
	// add check here for  status code. return success and response.

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		//return
	}

	return body
}
