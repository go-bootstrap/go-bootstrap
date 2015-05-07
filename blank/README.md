### Installation

1. Install PostgreSQL 9.4.x

2. Install Go 1.4.x, git, setup $GOPATH, and PATH=$PATH:$GOPATH/bin

3. Create PostgreSQL database.
    ```
    createdb $GO_BOOTSTRAP_REPO_NAME
    ```

4. Get the source code.
    ```
    go get github.com/$GO_BOOTSTRAP_REPO_USER/$GO_BOOTSTRAP_REPO_NAME
    ```

5. Run the PostgreSQL migration.
    ```
    go get github.com/mattes/migrate
    cd $GOPATH/src/github.com/$GO_BOOTSTRAP_REPO_USER/$GO_BOOTSTRAP_REPO_NAME
    migrate -url postgres://$(whoami)@$localhost:5432/$GO_BOOTSTRAP_REPO_NAME?sslmode=disable -path ./migrations up
    ```

6. Run the server
    ```
    cd $GOPATH/src/github.com/$GO_BOOTSTRAP_REPO_USER/$GO_BOOTSTRAP_REPO_NAME
    go run main.go
    ```

### Environment variables for Configuration

* **HTTP_ADDR:** The host and port. Default: `":8888"`

* **HTTP_CERT_FILE:** Path to cert file. Default: `""`

* **HTTP_KEY_FILE:** Path to key file. Default: `""`

* **HTTP_DRAIN_INTERVAL:** How long application will wait to drain old requests before restarting. Default: `"1s"`

* **DSN:** RDBMS database path. Default: `postgres://$(whoami)@localhost:5432/$GO_BOOTSTRAP_REPO_NAME?sslmode=disable`

* **COOKIE_SECRET:** Cookie secret for session. Default: See the source code.
