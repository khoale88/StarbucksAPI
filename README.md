# Multi-Tenant-Starbuck backend with Golang

## Software Dependency: 

### Local environment:

Language: golang:1.7

Required packages: github.com/gorilla/mux, gopkg.in/mgo.v2

### Docker environment:

Software tool: Docker

## Database dependency:

Before the back-end sever is run, in addition to satisfy software requirements, make sure the following conditions are also met:

1. Have at least one mongo database available with public access.

2. If a mongo cluster is desired, make sure they have proper replica configuration.

3. Modify three variables in server.go, mongo1, mongo2, mongo3, to have those addresses.

4. For more detail how to create a mongo cluster, visit the page:

   https://github.com/azhadm/multi-tenant-starbucks/blob/master/Mongo/README.md

## Test/Run application:

### Run in local environment: 

1. After installing Golang and required packages, go to directory where sever.go is located and run the following command : 

        go run sever.go

2. Visit design document to see available apis:

        https://github.com/azhadm/multi-tenant-starbucks/blob/master/Khoa.Restbucks/design/starbucks.pdf

3. Test can be done using Curl or Postman tool with address localhost:9090.

### Run in Docker:

1. After installing Docker, go to directory where Dockerfile is located and build an image.

######   syntax:
        docker build -t <image_name> .

######   eg:
        docker build -t KhoaRestBuck .

2. To check if the image is built successfully, run the following command and look for your <image_name>.

        docker images

3. Run the docker image:

######   syntax:
        docker run -it --rm -p 9090:9090 <image_name>

######   eg:
        docker build -it --rm -p 9090:9090 KhoaRestBuck

4. Test can be done using Curl or Postman tool with address localhost:9090.
