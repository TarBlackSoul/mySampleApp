package lib

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

/*
	For managing and routing http request to their corresponding defined handlers, gorrilla/mux library is used.
	This library supports middlewares to router, which receive the request at first, act on it and then pass it to the handler.
*/

// VARIABLES DEFINITION

//Define token for request authorization checking
const token = "token"

//Define a specific format for different types of logs, as info, debug, error or warning
var (
	Warning *log.Logger
	Info    *log.Logger
	Debug   *log.Logger
	Error   *log.Logger
)

//Init logs definition
func InitLogs(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {
	Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(warningHandle, "\033[33mDEBUG: \033[39m", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "\033[33mWARNING: \033[39m", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "\033[31mERROR: \033[39m", log.Ldate|log.Ltime|log.Lshortfile)
}

func DefineLogs() {
	InitLogs(os.Stdout, os.Stdout, os.Stdout)
}

//Init logs format
func init() {
	InitLogs(os.Stdout, os.Stdout, os.Stdout)
}

//Starting server and managing incoming APIs
func StartServer() {
	Info.Println("Starting server....")

	//Starting server
	r := mux.NewRouter()
	//Activating middleware requests check for authentication verification
	r.Use(Middleware)

	//Handle request for all applications, no authentication required
	//It will render all applications in json format
	r.HandleFunc("/applications/", GetHandler).
		Methods("GET").
		Schemes("http")

	//Handle request for one application given in the api, no authentication required
	//It will render one application in json format
	r.HandleFunc("/applications/{id}", GetOneAppHandler).
		Methods("GET").
		Schemes("http")

	//Add an application as a json payload, defining the structure as :
	// i.e:  {
	//        "name": "app1",
	//        "version": "v0.0.0",
	//        "author": "auth"
	//     }
	//Authorization token will be verified first, it should be in the format : "Authorization: Bearer token"
	r.HandleFunc("/applications/", PostHandler).
		Methods("POST").
		Schemes("http")

	//Delete an application by giving its name
	//Authorization token will be verified first
	r.HandleFunc("/applications/{id}", DeleteHandler).
		Methods("DELETE").
		Schemes("http")

	//Modify an application from db by passing a json structure as in POST requests
	//Authorization token will be verified first
	r.HandleFunc("/applications/{id}", PutHandler).
		//Headers("Authorization", "Bearer " + token).
		Methods("PUT").
		Schemes("http")

	//Server will start on localhost and will listen on port 8080
	log.Fatal(http.ListenAndServe("localhost:9000", r))
}

/*
Definition of middleware function which will intercept all incoming requests and check for Authorization Header
Header should be passed as in the example : -H "Authorization: Bearer token"
*/

//Verification of token sent in Authorization header
//In this case token value is defined as "token"
//If given token is not correct a 401 Not authorized response will be returned
func simpleAuth(w http.ResponseWriter, r *http.Request) (err string) {
	header := strings.Fields(r.Header.Get("Authorization"))
	if len(header) < 2 {
		Info.Println("No token given, not able to process request")
		return "ko"
	}
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if s[1] == token {
		Debug.Println("Token is eligible")
	} else {
		http.Error(w, "Not authorized", 401)

	}
	return ""
}

//Middleware function, which will be called for each request
//If request of type GET, token will not be verified
//If request of type POST/PUT/DELETE token will be verified
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Debug.Println("Method received is :", r.Method)
		switch m := r.Method; string(m) {
		//On get method token verification is skipped
		case "GET":
			Info.Println("On Get Method, no authentication needed")
			//Authorization process on middleware not needed, head on to request handlers
			next.ServeHTTP(w, r)
		//Verify token for the following methods
		case "POST", "PUT", "DELETE":
			sentToken := r.Header.Get("Authorization")
			Info.Println("Checking token", sentToken)
			authError := simpleAuth(w, r)
			if authError != "" {
				NotAuthorized(w, r)
				Info.Println("Not token given ! ")
			} else {
				//Authorization process on middleware level done, head on to request handlers
				next.ServeHTTP(w, r)
			}
		//In case a non defined http method is sent
		default:
			fmt.Print("not here ")
			w.Write([]byte("Not a supported http method"))
			//Error.Fatal("Not a supported method ! ")
		}
		Info.Println("End of middle level...")
	})
}
