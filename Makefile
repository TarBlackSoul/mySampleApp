MYGODOCKER="mygolang:1.15.8"
URL=http://localhost:9000/applications/
AUTH_HEADER="Authorization: Header token"
JSON_TYPE="Content-Type: application/json"



###Build and start server as a container

.PHONY: build run

build: ##build application and access it as a docker container
	docker build -f Dockerfile.builder -t myapp_image:latest . --no-cache


run: ##run docker app listening on port 9000, mount the projet
	docker run --net=host -d -p 9000:9000 --name myapp_container myapp_image	


stop:##stop running container
	docker stop myapp_container
	docker rm myapp_container

debug:#get logs
	docker logs myapp_container


###Add sample values

.PHONY: get init

get:
	curl -sv -X GET $(URL)


init: get ##add sample apps in database
	curl -sv -X POST  -H $(AUTH_HEADER) -H $(JSON_TYPE) -d '{"name":"first_app","version":"v1.0.0","author":"first_author"}' $(URL)
	curl -sv -X POST  -H $(AUTH_HEADER) -H $(JSON_TYPE) -d '{"name":"second_app","version":"v2.0.0","author":"second_author"}' $(URL)
	curl -sv -X POST  -H $(AUTH_HEADER) -H $(JSON_TYPE) -d '{"name":"third_app","version":"v3.0.0","author":"third_author"}' $(URL)
	make -s get
	make -s cleanupinit


cleanupinit:
	curl -sv -X DELETE -H $(AUTH_HEADER) $(URL)first_app
	curl -sv -X DELETE -H $(AUTH_HEADER) $(URL)second_app
	curl -sv -X DELETE -H $(AUTH_HEADER) $(URL)third_app
	make -s get


####Before building the code Linting and unit testing

.PHONY: fmt buildmygo style format

fmt: buildmygo style format


buildmygo: ##build golang image with required libraries
	 docker build -f Dockerfile.golang -t $(MYGODOCKER) .


style: ## Run golint and go vet
	docker run $(MYGODOCKER) go version
	docker run $(MYGODOCKER) golint ./... >> myAppReport.txt
	#docker run $(MYGODOCKER) go vet ./...

format: #change go files to the right format
	docker run $(MYGODOCKER) goimports -w




####Testing

.PHONY: test test-unit test-functional test-integration

test: test-unit test-functional test-integration


test-unit:


test-integration:
	make -s stop
	make -s run

	#add one application called test_first_app in the following json format, authentify first
	curl -sv -X POST  -H $(AUTH_HEADER) -H $(JSON_TYPE) -d '{"name":"test_first_app","version":"v1.0.0","author":"test_first_author"}' $(URL) |egrep -e "Data Added Successfully" 
        #get that app, no need to authentify
	curl -sv -X GET $(URL)test_first_app |jq .name |egrep -e test_first_app
	#get all apps
	curl -sv -X GET $(URL) |jq .

	#if json is not complete, empty values will be put in database, in this case author is missing form json
	curl -sv -X POST  -H $(AUTH_HEADER) -H $(JSON_TYPE) -d '{"name":"test_second_app","version":"v2.0.0"}' $(URL) |egrep -e "Data Added Successfully"
	#get that app, no need to authentify
	curl -sv -X GET $(URL)test_first_app |jq .
	

	##change a value, like author, authentify first
	curl -sv -X PUT -H $(AUTH_HEADER) -H $(JSON_TYPE) -d '{"author":"test_second_author"}' $(URL)test_first_app |egrep -e "Data Modified Successfully"
	#get that app to check if auth has changed
	curl -sv -X GET $(URL)test_first_app |jq .author |egrep -e test_second_author

	#imposible to replace a value from database with an empty one
	curl -sv -X PUT -H $(AUTH_HEADER) -H $(JSON_TYPE) -d '{"author":""}' $(URL)test_first_app |egrep -e "Empty values sent, nothing to change"
	#get that app to check if auth has changed
	curl -sv -X GET $(URL)test_first_app |jq .author |egrep -e test_second_author
	
	#imposible to replace the key, but the other existing values yes
	curl -sv -X PUT -H $(AUTH_HEADER) -H $(JSON_TYPE) -d '{"name":"change_first_app","version":"v1.0.1"}' $(URL)test_first_app |egrep -e "Data Modified Successfully"
	#get that app to check if auth has changed
	curl -sv -X GET $(URL)test_first_app |jq .name |egrep -e test_first_app
	curl -sv -X GET $(URL)test_first_app |jq .version |egrep -e "v1.0.1"

	#delete added app,need to authentify
	curl -sv -X DELETE -H $(AUTH_HEADER) $(URL)test_first_app |grep "Data Deleted Successfully"
	#get all apps
	curl -sv -X GET $(URL) |jq .


	


test-functional: buildmygo
	docker run $(MYGODOCKER) go test




####Security

.PHONY: security-sast

security-sast:
	docker run $(MYGODOCKER) gosec ./...



####Complexity

.PHONY: complexity

complexity:
	 docker run $(MYGODOCKER) gocyclo .
