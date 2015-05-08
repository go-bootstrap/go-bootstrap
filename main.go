// Package main generates web project.
package main

import (
	"flag"
	"fmt"
	"github.com/go-bootstrap/go-bootstrap/helpers"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {
	dir := flag.String("dir", "", "directory of project relative to $GOPATH/src/")
	flag.Parse()

	if *dir == "" {
		log.Fatal("dir option is missing.")
	}

	fullpath := os.ExpandEnv(filepath.Join("$GOPATH", "src", *dir))
	dirChunks := strings.Split(*dir, "/")
	repoName := dirChunks[len(dirChunks)-3]
	repoUser := dirChunks[len(dirChunks)-2]
	projectName := dirChunks[len(dirChunks)-1]
	dbName := projectName
	testDbName := projectName + "-test"

	// 1. Create target directory
	log.Print("Creating " + fullpath + "...")
	err := os.MkdirAll(fullpath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Copy everything under blank directory to target directory.
	log.Print("Copying a blank project to " + fullpath + "...")
	if output, err := exec.Command("cp", "-rf", "./blank/.", fullpath).CombinedOutput(); err != nil {
		log.Fatal(string(output))
	}

	// 3. Interpolate placeholder variables on the new project.
	log.Print("Replacing placeholder variables to " + repoUser + "/" + projectName + "...")
	replacers := make(map[string]string)
	replacers["$GO_BOOTSTRAP_REPO_NAME"] = repoName
	replacers["$GO_BOOTSTRAP_REPO_USER"] = repoUser
	replacers["$GO_BOOTSTRAP_PROJECT_NAME"] = projectName
	replacers["$GO_BOOTSTRAP_COOKIE_SECRET"] = helpers.RandString(16)
	if err := helpers.RecursiveSearchReplaceFiles(fullpath, replacers); err != nil {
		log.Fatal(err)
	}

	// 4. Create PostgreSQL databases.
	for _, name := range []string{dbName, testDbName} {
		log.Print("Creating a database named " + name + "...")
		if exec.Command("createdb", name).Run() != nil {
			log.Print("Unable to create PostgreSQL database: " + name)
		}
	}

	// 5.a. go get github.com/mattes/migrate.
	log.Print("Installing github.com/mattes/migrate...")
	if output, err := exec.Command("go", "get", "github.com/mattes/migrate").CombinedOutput(); err != nil {
		log.Fatal(string(output))
	}

	// 5.b. Run migrations on localhost:5432.
	u, _ := user.Current()
	for _, name := range []string{dbName, testDbName} {
		pgDSN := fmt.Sprintf("postgres://%v@localhost:5432/%v?sslmode=disable", u.Username, name)

		log.Print("Running database migrations on " + pgDSN + "...")
		if output, err := exec.Command("migrate", "-url", pgDSN, "-path", filepath.Join(fullpath, "migrations"), "up").CombinedOutput(); err != nil {
			log.Fatal(string(output))
		}
	}

	// 6. Get all application dependencies.
	log.Print("Running go get ./...")
	cmd := exec.Command("go", "get", "./...")
	cmd.Dir = fullpath
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatal(string(output))
	}

	// 7. Run tests on newly generated app.
	log.Print("Running go test ./...")
	cmd = exec.Command("go", "test", "./...")
	cmd.Dir = fullpath
	output, _ := cmd.CombinedOutput()
	log.Print(string(output))
}
