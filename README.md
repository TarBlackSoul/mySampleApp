# Application Description

This is a REST API Application written in Golang. Same techical details :
- The http methods implemented are : GET/DELETE/POST/PUT. The corresponding implemented golang library is : https://github.com/gorilla/mux
- A persistant key/value lightweight database type is implemented in the backend. The corresponding golang library is : https://github.com/boltdb/bolt

## Use cases
This application allows you to :
- Add a new instance in the database. The structure called Application contains following fields : 
  1.name Key 
  2.value: name: version: author:
- To add a new application you need to give a simple authentification header : "Authorization: Bearer token", give a content-type as application/json and give a json structure as defined above

