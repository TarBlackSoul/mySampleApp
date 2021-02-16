# Application Description

This is a REST API Application written in Golang. Same techical details :
- The http methods implemented are : GET/DELETE/POST/PUT. The corresponding implemented golang library is : https://github.com/gorilla/mux
- A persistant key/value lightweight database type is implemented in the backend. The corresponding golang library is : https://github.com/boltdb/bolt
- Server will start on localhost:9000
- Basic url will be http://localhost:9000/applications/

## Use cases
This application allows you to :

### POST
- Add a new instance in the database. The structure called Application contains following fields : 
  1.name Key 
  2.value: name: version: author:
- To add a new application you need to give a simple authentification header : "Authorization: Bearer token", give a content-type as application/json and give a json structure as defined above, i.e   {
    "name": "app",
    "version": "v0.0.0",
    "author": "me"
  }
- You can give a non complete json stucture i.e: {"name":"myapp"} and the instance will be created with only the key value (mandatory) and empty values for the rest of the structure

### PUT
- You can modify data by sending a PUT request, be aware, existing values will not be overwritten by empty ones
- PUT method works as POST one, you have to give a content-type json, authorization header and give the key you want to modify in the url i.e. applications/muyapp
- Give a json with the values you need to change, if values are empty or identical, database will not be modified and you will be noticed on http response

### DELETE 
- You can Delete you application by adding the key in the URL i.e. applications/myapp
- You have to give authorization header to delete

### GET
- You can either get all apps in json format
- Or get only one app by giving the key in the url
- You don't need to give a authorization header for get methods


# Packaging
- This golang application is built using go modules. All dependencies are defined in go.mod file
- A builder dockerfile (Dockerfile.builder) is used to build the application and allows to run it as a docker images
- A golang environement dockerfile (Dockerfile.golang) is used to get all dependencies and run commands using this container
- A Makefile containing all steps for building, formatting, testing 
- Only functional and integration tests are added
- A helm chart using the image built

## Microk8s setup 
- Enable helm3: microk8s enable helm3
- Enable  dns : microk8s.enable dns
- Enable docker registry : microk8s.enable registry
- Push the docker image on local registry : docker tag myappimage localhost:32000/myapp_image && docker push localhost:32000/myapp_image
