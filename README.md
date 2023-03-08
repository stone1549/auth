# auth
Example REST service for managing and authenticating users, written in Go using Chi.

## Configuration

### Environment Variables

##### AUTH_SERVICE_ENVIRONMENT

Controls log levels and configuration defaults. 

* DEV
* PRE_PROD
* PROD
 
##### AUTH_SERVICE_REPO_TYPE

* IN_MEMORY
* POSTGRESQL
    * AUTH_SERVICE_PG_URL - Full connection string for PG

##### AUTH_SERVICE_TIMEOUT

Incoming request timeout value in seconds.

##### AUTH_SERVICE_PORT

Port to run service on.


## Run

```go run main.go```

