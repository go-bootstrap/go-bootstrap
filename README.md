[![GoDoc](https://godoc.org/github.com/go-bootstrap/go-bootstrap?status.svg)](http://godoc.org/github.com/go-bootstrap/go-bootstrap)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/go-bootstrap/go-bootstrap/master/LICENSE.md)

## go-bootstrap

This is not a web framework. It generates a skeleton web project for you to kick-ass.

Feel free to use or rip-out any of its parts.


## Installation

1. `mkdir -p $GOPATH/src/github.com/go-bootstrap && cd $GOPATH/src/github.com/go-bootstrap`

2. `git clone git@github.com:go-bootstrap/go-bootstrap.git`

3. `go run main.go -dir github.com/{git-user}/{project-name}`

4. Start using it: `cd $GOPATH/src/github.com/{git-user}/{project-name} && go run main.go`


## Decisions It Made for You

This generator makes **A LOT** of decisions for you. Here's the list of things it uses for your project:

1. PostgreSQL is chosen for the database.

2. bcrypt is chosen as the password hasher.

3. Bootstrap Flatly is chosen for the UI theme.

4. Session is stored inside encrypted cookie.

5. Static directory is located under `/static`.

6. Model directory is located under `/dal` (Database Access Layer).

7. It does not use ORM nor installs one.

8. Test database is automatically created under `$GO_BOOTSTRAP_REPO_NAME-test`.

9. [github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx) is chosen to connect to a database.

10. [github.com/gorilla](https://github.com/gorilla) is chosen for a lot of the HTTP plumbings.

11. [github.com/carbocation/interpose](https://github.com/carbocation/interpose) is chosen as the middleware library.

12. [github.com/mattes/migrate](https://github.com/mattes/migrate) is chosen as the database migration tool.

13. [github.com/Sirupsen/logrus](https://github.com/Sirupsen/logrus) is chosen as the logging library.
