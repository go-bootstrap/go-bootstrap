## Installation

1. Install Go 1.4.x or greater, git, setup `$GOPATH`, and `PATH=$PATH:$GOPATH/bin`

2. Run the server
    ```
    cd $GOPATH/src/$GO_BOOTSTRAP_REPO_NAME/$GO_BOOTSTRAP_REPO_USER/$GO_BOOTSTRAP_PROJECT_NAME
    go run main.go
    ```


## Environment Variables for Configuration

* **HTTP_ADDR:** The host and port. Default: `":8888"`

* **HTTP_CERT_FILE:** Path to cert file. Default: `""`

* **HTTP_KEY_FILE:** Path to key file. Default: `""`

* **HTTP_DRAIN_INTERVAL:** How long application will wait to drain old requests before restarting. Default: `"1s"`

* **COOKIE_SECRET:** Cookie secret for session. Default: Auto generated.


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

2. PostgreSQL is not installed inside Vagrant. You must connect to your host PostgreSQL. Here's an example on how to run your application inside vagrant while connecting to your host PostgreSQL:
```
GOPATH=/go DSN=postgres://$(whoami)@$(netstat -rn | grep "^0.0.0.0 " | cut -d " " -f10):5432/$PROJECT_NAME?sslmode=disable go run main.go
```
