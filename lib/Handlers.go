package lib

import (
	"awesomeProject1/db"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//Specifying handlers for each request method

//Handler to get all applications API
//Returns all the data from the database
//Request defined as /applications/
func GetHandler(w http.ResponseWriter, r *http.Request) {
	Debug.Println("On get all apps handler : ")
	//Calling function to get applications from database
	apps, err := db.GetApps()
	if err != nil {
		Error.Fatal("could not get apps from db, ", err)
		//If error on retrieving data, 500 Internal error message is returned
		InternalError(w, r)
	}
	if apps != nil {
		//Show response if data retrieved
		RenderAllApps(w, r, apps)
	} else {
		//In case database is empty
		NotFound(w, r)
	}
}

//Handler to get only one applications API
//Returns information corresponding to one key/value from the database
//Request defined as /applications/{id}
func GetOneAppHandler(w http.ResponseWriter, r *http.Request) {
	//Get value sent in the API
	vars := mux.Vars(r)
	id := vars["id"]
	Debug.Println("Extracting data on db for : ", id)
	//Calling function to get only one application from database
	app, err := db.GetApp(id)
	if err != nil {
		Error.Println("could not get app from db, ", err)
		//If error on retrieving data, 500 Internal error message is returned
		InternalError(w, r)
	}
	var a db.Application
	_ = json.Unmarshal(app, &db.Application{})
	Debug.Println("Extracted data from db : ", a)
	if app != nil {
		//Show response if data retrieved
		Debug.Println("App to render")
		RenderApp(w, r, app)
	} else {
		//In case database is empty
		Debug.Println("Nothing to render")
		NotFound(w, r)
	}
}

//Handler to add a new application, passed as a json, using header -H "Content-Type: application/json"
func PostHandler(w http.ResponseWriter, r *http.Request) {
	//check if content-type of request is application-json
	defer r.Body.Close()
	if r.Header.Get("Content-type") != "application/json" {
		//If no header passed, 415 error message will be returned
		NotSupported(w, r)
	}
	//check if json sent is in the right format as the structure Application defined
	var a db.Application
	e := json.NewDecoder(r.Body).Decode(&a)
	if e != nil {
		Error.Println("json sent not in the right format, ", e)
		InternalError(w, r)
		return
	}
	Debug.Println("Adding values to database", a.Name, a.Author, a.Version)
	//Check if the application already exists on the database
	app, _ := db.GetApp(a.Name)
	if app != nil {
		//If data exists, it will not be re-added
		var curr db.Application
		_ = json.Unmarshal(app, &curr)
		Warning.Println("Application already exists on database, it won't be added : ", curr)
		ActionDone(w, r, "data exists, nothing to do")
		return
	} else {
		//Call function to add a new instance of Application in database
		err := db.AddApp(a)
		if err != nil {
			Error.Println("could not add app to db, ", err)
			InternalError(w, r)
			return
		}
		//Message sent on response
		Added(w, r, "\nData Added Successfully\n")
	}
}

//Handler to delete an instance from the database, based on key (application name)
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	//Get value sent in the API
	vars := mux.Vars(r)
	id := vars["id"]
	//check if data exists
	Debug.Println("Extracting data on db for : ", id)
	//Getting data from database
	app, err := db.GetApp(id)
	if err != nil {
		Error.Println("could not get app from db, ", err)
		InternalError(w, r)
	}
	//If data does not exit, it cannot be deleted
	if app == nil {
		NotFound(w, r)
	} else {
		//Data exists, deleting it
		Debug.Println("Deleting data : ", id)
		//Call function to delete instance from database
		err = db.DeleteApp(id)
		if err != nil {
			Error.Println("could not get app from db, ", err)
			InternalError(w, r)
		}
		var a db.Application
		_ = json.Unmarshal(app, &a)
		Debug.Println("Extracted data from db : ", a)
		ActionDone(w, r, "\nData Deleted Successfully\n")
	}
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//check if content-type of request is application/json
	if r.Header.Get("Content-type") != "application/json" {
		NotSupported(w, r)
		return
	}
	//Get value sent in the API
	vars := mux.Vars(r)
	id := vars["id"]

	//check if json sent is in the right format
	var sentApp db.Application
	e := json.NewDecoder(r.Body).Decode(&sentApp)
	if e != nil {
		Error.Println("json sent not in the right format, ", e)
		InternalError(w, r)
		return
	}
	Debug.Println("Values sent to be changed to the database", sentApp.Author, sentApp.Version)

	//Check if value passed in api exists in database, so it can be modified or not
	var existingApp db.Application
	getExistingApp, _ := db.GetApp(id)
	//checking if data exists or not
	if getExistingApp == nil {
		Failed(w, r, "\nData does not exist, create it first \n")
	} else {
		//check values which exist on database corresponding to the key sent
		//if a value like author or version exists on database, it will not be replaced by an empty value
		err := json.Unmarshal(getExistingApp, &existingApp)
		if err != nil {
			Error.Println("value sent from database not in the right format, ", err)
			InternalError(w, r)
			return
		}
		Debug.Println("Extracted data from db : ", existingApp.Name, existingApp.Author, existingApp.Version)

		//nothing will happen if empty values are sent
		if sentApp.Author == "" && sentApp.Version == "" {
			Warning.Println("values sent to are empty to change ")
			Failed(w, r, "\nEmpty values sent, nothing to change\n")
			return
		}

		//We don't want to overwrite existing values by empty ones
		var authorValueToAdd string
		if sentApp.Author == "" && existingApp.Author != "" {
			Debug.Println("value Author sent is empty ", sentApp.Author, " it will not replace exiting values ", existingApp.Author)
			Warning.Println("value sent for Author field is empty ", sentApp.Author)
			authorValueToAdd = existingApp.Author
		} else {
			authorValueToAdd = sentApp.Author
			Debug.Println("Existing author value :", existingApp.Author, " will be replaced by :", authorValueToAdd)
		}

		var versionValueToAdd string
		if sentApp.Version == "" && existingApp.Version != "" {
			Debug.Println("value Version sent is empty ", sentApp.Version, "it will not replace exiting values ", existingApp.Version)
			Warning.Println("value sent for Value field is empty ", sentApp.Version)
			versionValueToAdd = existingApp.Version
		} else {
			Debug.Println("Existing Version value : ", existingApp.Version, " will be replaced by : ", authorValueToAdd)
			versionValueToAdd = sentApp.Version
		}

		//Change values on database
		//It is impossible to change the key only its corresponding values
		Info.Println("Calling add to database function with : ", id, authorValueToAdd, versionValueToAdd)
		changeApp := db.Application{Name: id, Version: versionValueToAdd, Author: authorValueToAdd}
		err = db.AddApp(changeApp)
		if err != nil {
			Error.Println("could not add app to db ", err)
			InternalError(w, r)
		}
		Added(w, r, "\nData Modified Successfully\n")
		return
	}
}
