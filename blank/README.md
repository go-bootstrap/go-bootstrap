## Installation

1. Install PostgreSQL 9.4.x

2. Install Go 1.4.x, git, setup `$GOPATH`, and `PATH=$PATH:$GOPATH/bin`

3. Create PostgreSQL database.
    ```
    createdb $GO_BOOTSTRAP_PROJECT_NAME
    ```

4. Get the source code.
    ```
    go get $GO_BOOTSTRAP_REPO_NAME/$GO_BOOTSTRAP_REPO_USER/$GO_BOOTSTRAP_PROJECT_NAME
    ```

5. Run the PostgreSQL migration.
    ```
    go get github.com/mattes/migrate
    cd $GOPATH/src/$GO_BOOTSTRAP_REPO_NAME/$GO_BOOTSTRAP_REPO_USER/$GO_BOOTSTRAP_PROJECT_NAME
    migrate -url postgres://$(whoami)@$localhost:5432/$GO_BOOTSTRAP_PROJECT_NAME?sslmode=disable -path ./migrations up
    ```

6. Run the server
    ```
    cd $GOPATH/src/$GO_BOOTSTRAP_REPO_NAME/$GO_BOOTSTRAP_REPO_USER/$GO_BOOTSTRAP_PROJECT_NAME
    go run main.go
    ```


## Environment Variables for Configuration

* **HTTP_ADDR:** The host and port. Default: `":8888"`

* **HTTP_CERT_FILE:** Path to cert file. Default: `""`

* **HTTP_KEY_FILE:** Path to key file. Default: `""`

* **HTTP_DRAIN_INTERVAL:** How long application will wait to drain old requests before restarting. Default: `"1s"`

* **DSN:** RDBMS database path. Default: `postgres://$(whoami)@localhost:5432/$GO_BOOTSTRAP_PROJECT_NAME?sslmode=disable`

* **COOKIE_SECRET:** Cookie secret for session. Default: Auto generated.


## Running Migrations

Migration is handled by a separate project: [github.com/mattes/migrate](https://github.com/mattes/migrate).

Here's a quick tutorial on how to use it. For more details, read the tutorial [here](https://github.com/mattes/migrate#usage-from-terminal).
```
# Installing the library
go get github.com/mattes/migrate

# Create a new migration file
migrate -url driver://url -path ./migrations create {filename}

# Migrate all the way up
migrate -url driver://url -path ./migrations up

# Migrate all the way down
migrate -url driver://url -path ./migrations down

# Roll back the most recently applied migration, then run it again.
migrate -url driver://url -path ./migrations redo

# Run down and then up command
migrate -url driver://url -path ./migrations reset

# Show the current migration version
migrate -url driver://url -path ./migrations version
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
