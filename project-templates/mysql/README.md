## Installation

1. Install MySQL

2. Install Go 1.5.x, git, setup `$GOPATH`, and `PATH=$PATH:$GOPATH/bin`

3. Setup MySQL database.
    ```
    go get github.com/mattes/migrate
    cd $GOPATH/src/$GO_BOOTSTRAP_REPO_NAME/$GO_BOOTSTRAP_REPO_USER/$GO_BOOTSTRAP_PROJECT_NAME
    ./scripts/db-bootstrap
    ```

4. Run the server
    ```
    cd $GOPATH/src/$GO_BOOTSTRAP_REPO_NAME/$GO_BOOTSTRAP_REPO_USER/$GO_BOOTSTRAP_PROJECT_NAME
    go run main.go
    ```


## Environment Variables for Configuration

* **HTTP_ADDR:** The host and port. Default: `":8888"`

* **HTTP_CERT_FILE:** Path to cert file. Default: `""`

* **HTTP_KEY_FILE:** Path to key file. Default: `""`

* **HTTP_DRAIN_INTERVAL:** How long application will wait to drain old requests before restarting. Default: `"1s"`

* **DSN:** RDBMS database path. Default: `$(whoami)@tcp(localhost:3306)/$GO_BOOTSTRAP_PROJECT_NAME?parseTime=true`

* **COOKIE_SECRET:** Cookie secret for session. Default: Auto generated.


## Running Migrations

Migration is handled by a separate project: [github.com/mattes/migrate](https://github.com/mattes/migrate).

Here's a quick tutorial on how to use it. For more details, read the tutorial [here](https://github.com/mattes/migrate#usage-from-terminal).
```
# Installing the library
go get github.com/mattes/migrate

# create new migration file in path
migrate -url driver://url -path ./migrations create migration_file_xyz

# apply all available migrations
migrate -url driver://url -path ./migrations up

# roll back all migrations
migrate -url driver://url -path ./migrations down

# roll back the most recently applied migration, then run it again.
migrate -url driver://url -path ./migrations redo
```


## Vendoring Dependencies

Vendoring is handled by a separate project: [github.com/tools/godep](https://github.com/tools/godep).

Here's a quick tutorial on how to use it. For more details, read the readme [here](https://github.com/tools/godep#godep).
```
# Save all your dependencies after running go get ./...
godep save ./...

# Building with godep
godep go build

# Running tests with godep
godep go test ./...
```


## Running in Vagrant

There are two potential gotchas you need to know when running in Vagrant:

1. `GOPATH` is not defined when you ssh into Vagrant. To fix the problem, do `export GOPATH=/go` immediately after ssh.

2. MySQL is not installed inside Vagrant. You must connect to your host MySQL. Here's an example on how to run your application inside vagrant while connecting to your host MySQL:
```
GOPATH=/go DSN=mysql://$(whoami)@tcp($(netstat -rn | grep "^0.0.0.0 " | cut -d " " -f10):5432)/$PROJECT_NAME?parseTime=true go run main.go
```
