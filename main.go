package main

import (
	lib "awesomeProject1/lib"
)

/*
	This application will manipulate the database via Rest API methods : DELETE/GET/PUT/POST

*/

func main() {
	//Running the first functions to start the server and handle incoming requests
	lib.StartServer()
}
