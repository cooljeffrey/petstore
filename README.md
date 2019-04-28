# PetStore API
Implemented by Go.

API schema reference : https://petstore.swagger.io/

Schema saved to local file at : `schema/petstore.json`, modified default host and set schema to http only.

## Prerequisite

    go version go1.12.3 or above with [go module](https://github.com/golang/go/wiki/Modules) support enabled
    node v10 or above and yarn 1.12.3 or above for integration test
    docker 18.09.2 or above for running mongodb
    
## Start MongoDB

    docker run -d -p 27017:27017 --name petstore-mongo mongo
    
## Unit test

Please make sure mongodb is running before running the following test

    go test ./...
    
The produced data in mongodb will be cleared upon test complete.
    
## Integration test

Please make sure mongodb is running before running the following test.

Start API server

    go run main.go
    
Run jest tests in a seperate terminal console:

    yarn install
    yarn test
    
The produced data in mongodb will *NOT* be cleared upon test complete.
    
## Known Issues

 * Unit tests not done for services
 * Integration test covered simple ( eg. successful calls ) only.
 * Swagger Client uploading empty file so the test case is commented out while the API is working as expected and can be tested using curl : 
    
    `curl -X POST "http://localhost:8080/v2/pet/2/uploadImage" -H "accept: application/json" -H "Content-Type: multipart/form-data" -F "additionalMetadata=test" -F "file=@__tests__/test.jpg"`
    


    