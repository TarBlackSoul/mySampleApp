package main

import (
	"awesomeProject1/db"
	"awesomeProject1/lib"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/steinfletcher/apitest"
)

//TEST CASES

//Delete test app without asserting the response, so in case testing app already exists it will deleted
func TestPreDeleteApplication(t *testing.T) {
	fmt.Println("pre-delete test...")
	// Test deleting applications API
	req, err := http.NewRequest("DELETE", "/applications/test_app", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer token")
	//Chck for more info on errors
	getErrorCauseStackTrace(err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	//Redirect the request to indicated path
	router.HandleFunc("/applications/{id}", lib.DeleteHandler)
	router.ServeHTTP(rr, req)
}

func TestPostApplication(t *testing.T) {
	// Test add new applications API
	fmt.Println("post test")
	testAddApp := db.Application{Name: "test_app", Version: "test_version", Author: "test_auth"}
	jsonBytes, err := json.Marshal(testAddApp)
	req, err := http.NewRequest("POST", "/applications/", bytes.NewReader(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	//show more information on the error
	getErrorCauseStackTrace(err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(lib.PostHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect
	assert(rr.Code, http.StatusCreated, t)
	// Check the response header is what we expect
	assert(rr.Header().Get("Content-Type"), "text/plain", t)
	//Check the response body is what we expect
	assert(rr.Body.String(), "\nData Added Successfully\n", t)

	//Generate report for post request without asserting result
	apitest.New().
		Report(apitest.SequenceDiagram()).
		HandlerFunc(lib.PostHandler).
		Post("/applications/").
		Body(`"{"name":"test_app","version":"test_version","author":"test_auth"}"`).
		Header("Authorization", "Bearer token").
		Header("Content-Type", "application/json")
}

func TestAllApps(t *testing.T) {
	// Test get one application API
	fmt.Println("Get all applications test")
	req, err := http.NewRequest("GET", "/applications/", nil)
	if err != nil {
		t.Fatal(err)
	}
	getErrorCauseStackTrace(err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(lib.GetHandler)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	assert(rr.Code, http.StatusOK, t)
	//Check the response content-type
	expected := "application/json"
	assert(rr.Header().Get("Content-Type"), expected, t)
	fmt.Println(rr.Body.String())

	//Generate report for tests
	apitest.New().
		Report(apitest.SequenceDiagram()).
		HandlerFunc(lib.GetHandler).
		Get("/applications/")
}

func TestOneApp(t *testing.T) {
	// Test get all applications API
	fmt.Print("Get only one app test")

	req, err := http.NewRequest("GET", "/applications/test_app", nil)
	if err != nil {
		t.Fatal(err)
	}
	getErrorCauseStackTrace(err)
	println(req.Host)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/applications/{id}", lib.GetOneAppHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert(rr.Code, http.StatusOK, t)
	//Check the response content-type
	expected := "application/json"
	assert(rr.Header().Get("Content-Type"), expected, t)
	assert(rr.Body.String(), "{\"name\":\"test_app\",\"version\":\"test_version\",\"author\":\"test_auth\"}", t)

	//Generate report for test
	apitest.New().
		Report(apitest.SequenceDiagram()).
		HandlerFunc(lib.GetOneAppHandler).
		Get("/applications/test_app")
}

func TestModifyApplication(t *testing.T) {
	// Test add new applications API
	testAddApp := db.Application{Name: "test_app", Version: "new_test_version", Author: ""}
	jsonBytes, err := json.Marshal(testAddApp)
	req, err := http.NewRequest("PUT", "/applications/test_app", bytes.NewReader(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")
	getErrorCauseStackTrace(err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	//Redirect the request to indicated path
	router.HandleFunc("/applications/{id}", lib.PutHandler)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert(rr.Code, http.StatusCreated, t)
	//Check the response content-type
	assert(rr.Header().Get("Content-Type"), "text/plain", t)
	assert(rr.Body.String(), "\nData Modified Successfully\n", t)

	//Generate report for post request without asserting result
	apitest.New().
		Report(apitest.SequenceDiagram()).
		HandlerFunc(lib.PutHandler).
		Post("/applications/test_app").
		Body(`"{"name":"test_app","version":"new_test_version","author":""}"`).
		Header("Authorization", "Bearer token").
		Header("Content-Type", "application/json")
}

func TestEndDeleteApplication(t *testing.T) {
	fmt.Println("end-delete test...")
	// Test deleting applications API
	req, err := http.NewRequest("DELETE", "/applications/test_app", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer token")
	//Chck for more info on errors
	getErrorCauseStackTrace(err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	//Redirect the request to indicated path
	router.HandleFunc("/applications/{id}", lib.DeleteHandler)
	router.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	assert(rr.Code, http.StatusOK, t)
	// Check the response header is what we expect
	assert(rr.Header().Get("Content-Type"), "text/plain", t)
	//Check the response body is what we expect
	assert(rr.Body.String(), "\nData Deleted Successfully\n", t)

	//Generate report for post request without asserting result
	apitest.New().
		Report(apitest.SequenceDiagram()).
		HandlerFunc(lib.DeleteHandler).
		Post("/applications/test_app").
		Header("Authorization", "Bearer token").
		Header("Content-Type", "application/json")
}

//END OF TEST CASES

func assert(actualResult interface{}, expectedResult interface{}, t *testing.T) {
	if actualResult != expectedResult {
		t.Fatalf("Expected '%v' but got %v", expectedResult, actualResult)
	}
}

func getErrorCauseStackTrace(err error) errors.StackTrace {
	// This code is inspired by github.com/pkg/errors.Cause().
	var st errors.StackTrace
	for err != nil {
		s := getErrorStackTrace(err)
		if s != nil {
			st = s
		}
		err = getErrorCause(err)
	}
	return st
}

func getErrorStackTrace(err error) errors.StackTrace {
	ster, ok := err.(interface {
		StackTrace() errors.StackTrace
	})
	if !ok {
		return nil
	}
	return ster.StackTrace()
}

func getErrorCause(err error) error {
	cer, ok := err.(interface {
		Cause() error
	})
	if !ok {
		return nil
	}
	return cer.Cause()
}
