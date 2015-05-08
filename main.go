// Package main generates web project.
package main

import (
	"flag"
	"github.com/go-bootstrap/go-bootstrap/helpers"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func chDir(dir string) {
	err := os.Chdir(dir)
	exitOnError(err, "")
}

func exitOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s\n%s", msg, err.Error())
	}
}

func main() {
	dir := flag.String("dir", "", "Project directory relative to $GOPATH/src/")
	flag.Parse()

	if *dir == "" {
		log.Fatal("dir option is missing.")
	}

	// There can be more than one path, separated by colon.
	gopaths := strings.Split(os.ExpandEnv("$GOPATH"), ":")
	gopath := gopaths[0]

	fullpath := filepath.Join(gopath, "src", *dir)
	migrationsPath := filepath.Join(fullpath, "migrations")
	dirChunks := strings.Split(*dir, "/")
	repoName := dirChunks[len(dirChunks)-3]
	repoUser := dirChunks[len(dirChunks)-2]
	projectName := dirChunks[len(dirChunks)-1]
	dbName := projectName
	testDbName := projectName + "-test"
	currentUser, _ := user.Current()

	// 1. Create target directory
	log.Print("Creating " + fullpath + "...")
	err := os.MkdirAll(fullpath, 0755)
	exitOnError(err, "")

	// 2. Copy everything under blank directory to target directory.
	log.Print("Copying a blank project to " + fullpath + "...")
	blankDir := os.ExpandEnv(filepath.Join("$GOPATH", "src", "github.com", "go-bootstrap", "go-bootstrap", "blank"))
	currDir, err := os.Getwd()
	exitOnError(err, "Can't get current path!")

	chDir(blankDir)

	output, err := exec.Command("cp", "-rf", ".", fullpath).CombinedOutput()
	exitOnError(err, string(output))

	chDir(currDir)

	// 3. Interpolate placeholder variables on the new project.
	log.Print("Replacing placeholder variables on " + repoUser + "/" + projectName + "...")
	replacers := make(map[string]string)
	replacers["$GO_BOOTSTRAP_REPO_NAME"] = repoName
	replacers["$GO_BOOTSTRAP_REPO_USER"] = repoUser
	replacers["$GO_BOOTSTRAP_PROJECT_NAME"] = projectName
	replacers["$GO_BOOTSTRAP_COOKIE_SECRET"] = helpers.RandString(16)
	replacers["$GO_BOOTSTRAP_CURRENT_USER"] = currentUser.Username
	replacers["$GO_BOOTSTRAP_DOCKERFILE_DSN"] = helpers.DefaultPGDSN(dbName)
	err = helpers.RecursiveSearchReplaceFiles(fullpath, replacers)
	exitOnError(err, "")

	// 4. Create PostgreSQL databases.
	for _, name := range []string{dbName, testDbName} {
		log.Print("Creating a database named " + name + "...")
		if exec.Command("createdb", name).Run() != nil {
			log.Print("Unable to create PostgreSQL database: " + name)
		}
	}

	// 5.a. go get github.com/mattes/migrate.
	log.Print("Installing github.com/mattes/migrate...")
	output, err = exec.Command("go", "get", "github.com/mattes/migrate").CombinedOutput()
	exitOnError(err, string(output))

	// 5.b. Run migrations on localhost:5432.
	for _, name := range []string{dbName, testDbName} {
		pgDSN := helpers.DefaultPGDSN(name)

		log.Print("Running database migrations on " + pgDSN + "...")
		output, err := exec.Command("migrate", "-url", pgDSN, "-path", migrationsPath, "up").CombinedOutput()
		exitOnError(err, string(output))
	}

	// 6. Get all application dependencies for the first time.
	log.Print("Running go get ./...")
	cmd := exec.Command("go", "get", "./...")
	cmd.Dir = fullpath
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatal(string(output))
	}

	repoIsGit := strings.HasPrefix(repoName, "git")

	if repoIsGit {
		// Generate Godeps directory. Currently only works on git related repo.
		log.Print("Installing github.com/tools/godep...")
		output, err := exec.Command("go", "get", "github.com/tools/godep").CombinedOutput()
		exitOnError(err, string(output))

		// git init.
		log.Print("Running git init")
		cmd := exec.Command("git", "init")
		cmd.Dir = fullpath
		output, err = cmd.CombinedOutput()
		exitOnError(err, string(output))

		// godep save ./...
		log.Print("Running godep save ./...")
		cmd = exec.Command("godep", "save", "./...")
		cmd.Dir = fullpath
		output, err = cmd.CombinedOutput()
		exitOnError(err, string(output))

		// Run tests on newly generated app.
		log.Print("Running godep go test ./...")
		cmd = exec.Command("godep", "go", "test", "./...")
		cmd.Dir = fullpath
		output, _ = cmd.CombinedOutput()
		log.Print(string(output))

	} else {
		// Run tests on newly generated app.
		log.Print("Running go test ./...")
		cmd = exec.Command("go", "test", "./...")
		cmd.Dir = fullpath
		output, _ := cmd.CombinedOutput()
		log.Print(string(output))
	}
}
