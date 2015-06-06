[![GoDoc](https://godoc.org/github.com/go-bootstrap/go-bootstrap?status.svg)](http://godoc.org/github.com/go-bootstrap/go-bootstrap)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/go-bootstrap/go-bootstrap/master/LICENSE.md)

## go-bootstrap

This is not a web framework. It generates a skeleton web project for you to kick-ass.

Feel free to use or rip-out any of its parts.


## Prerequisites

1. PostgreSQL

2. Go programming language, version 1.3.x or newer.

## Installation

1. `go get github.com/go-bootstrap/go-bootstrap`

2. `$GOPATH/bin/go-bootstrap -dir github.com/{git-user}/{project-name}`

3. Start using it: `cd $GOPATH/src/github.com/{git-user}/{project-name} && go run main.go`


## PostgreSQL Environment Variables

If you have `PGUSER`, `PGPASSWORD`, `PGHOST`, `PGPORT`, `PGSSLMODE` environment variables set,
they will be used to generate and bootstrap the database.


## Decisions Made for You

This generator makes **A LOT** of decisions for you. Here's the list of things it uses for your project:

1. PostgreSQL is chosen for the database.

2. bcrypt is chosen as the password hasher.

3. Bootstrap Flatly is chosen for the UI theme.

4. Session is stored inside encrypted cookie.

5. Static directory is located under `/static`.

6. Model directory is located under `/dal` (Database Access Layer).

7. It does not use ORM nor installs one.

8. Test database is automatically created under `$GO_BOOTSTRAP_PROJECT_NAME-test`.

9. A minimal Dockerfile is provided.

10. A minimal Vagrantfile is provided.

11. [github.com/tools/godep](https://github.com/tools/godep) is chosen to manage dependencies.

12. [github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx) is chosen to connect to a database.

13. [github.com/gorilla](https://github.com/gorilla) is chosen for a lot of the HTTP plumbings.

14. [github.com/carbocation/interpose](https://github.com/carbocation/interpose) is chosen as the middleware library.

15. [github.com/tylerb/graceful](https://github.com/tylerb/graceful) is chosen to enable graceful shutdown.

16. [github.com/rnubel/pgmgr](https://github.com/rnubel/pgmgr) is chosen as the database migration and management tool.

17. [github.com/Sirupsen/logrus](https://github.com/Sirupsen/logrus) is chosen as the logging library.
