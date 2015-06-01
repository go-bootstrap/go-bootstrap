// Package main generates web project.
package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/go-bootstrap/go-bootstrap/helpers"
)

func main() {
	dirInput := flag.String("dir", "", "Project directory relative to $GOPATH/src/")
	gopathInput := flag.String("gopath", "", "Choose which $GOPATH to use")

	flag.Parse()

	if *dirInput == "" {
		log.Fatalln("dir option is missing.")
	}

	// There can be more than one path, separated by colon.
	gopaths := helpers.GoPaths()
	if len(gopaths) == 0 {
		log.Fatalln("GOPATH is not set.")
	}

	// By default, we choose the last GOPATH.
	gopath := gopaths[len(gopaths)-1]

	// But if user specified one, we choose that one.
	if *gopathInput != "" {
		abs, err := filepath.Abs(*gopathInput)
		if err == nil && helpers.IsValidGoPath(abs) {
			gopath = abs
		} else {
			log.Fatalln("Cannot find " + *gopathInput + " in $GOPATH")
		}
	}

	trimmedPath := strings.Trim(*dirInput, "/")
	fullpath := filepath.Join(gopath, "src", trimmedPath)
	dirChunks := strings.Split(trimmedPath, "/")

	if len(dirChunks) < 3 {
		log.Fatalln("Cannot extract repo name, repo user and project name, " +
			"-dir should have three parts, seperated by '/'.")
	}

	repoName := dirChunks[len(dirChunks)-3]
	repoUser := dirChunks[len(dirChunks)-2]
	projectName := dirChunks[len(dirChunks)-1]
	dbName := projectName
	testDbName := projectName + "-test"
	currentUser, _ := user.Current()

	blankDir, err := helpers.GetBlankDir()
	helpers.ExitOnError(err, "")

	// 1. Create target directory
	log.Print("Creating " + fullpath + "...")
	err = os.MkdirAll(fullpath, 0755)
	helpers.ExitOnError(err, "")

	// 2. Copy everything under blank directory to target directory.
	log.Print("Copying a blank project to " + fullpath + "...")
	currDir, err := os.Getwd()
	helpers.ExitOnError(err, "Can't get current path!")

	err = os.Chdir(blankDir)
	helpers.ExitOnError(err, "")

	output, err := exec.Command("cp", "-rf", ".", fullpath).CombinedOutput()
	helpers.ExitOnError(err, string(output))

	err = os.Chdir(currDir)
	helpers.ExitOnError(err, "")

	// 3. Interpolate placeholder variables on the new project.
	log.Print("Replacing placeholder variables on " + repoUser + "/" + projectName + "...")

	replacers := make(map[string]string)
	replacers["$GO_BOOTSTRAP_REPO_NAME"] = repoName
	replacers["$GO_BOOTSTRAP_REPO_USER"] = repoUser
	replacers["$GO_BOOTSTRAP_PROJECT_NAME"] = projectName
	replacers["$GO_BOOTSTRAP_COOKIE_SECRET"] = helpers.RandString(16)
	replacers["$GO_BOOTSTRAP_CURRENT_USER"] = currentUser.Username
	replacers["$GO_BOOTSTRAP_PG_DSN"] = helpers.DefaultPGDSN(dbName)
	replacers["$GO_BOOTSTRAP_ESCAPED_PG_DSN"] = helpers.BashEscape(helpers.DefaultPGDSN(dbName))
	replacers["$GO_BOOTSTRAP_PG_TEST_DSN"] = helpers.DefaultPGDSN(testDbName)
	replacers["$GO_BOOTSTRAP_ESCAPED_PG_TEST_DSN"] = helpers.BashEscape(helpers.DefaultPGDSN(testDbName))

	err = helpers.RecursiveSearchReplaceFiles(fullpath, replacers)
	helpers.ExitOnError(err, "")

	// 4. go get github.com/rnubel/pgmgr
	log.Print("Running go get github.com/rnubel/pgmgr...")
	output, err = exec.Command("go", "get", "github.com/rnubel/pgmgr").CombinedOutput()
	helpers.ExitOnError(err, string(output))

	// 5. Bootstrap databases.
	cmd := exec.Command("bash", "scripts/db-bootstrap")
	cmd.Dir = fullpath
	output, _ = cmd.CombinedOutput()
	log.Print(string(output))

	// 6. Get all application dependencies for the first time.
	log.Print("Running go get ./...")
	cmd = exec.Command("go", "get", "./...")
	cmd.Dir = fullpath
	output, err = cmd.CombinedOutput()
	helpers.ExitOnError(err, string(output))

	repoIsGit := strings.HasPrefix(repoName, "git")
	repoIsHg := strings.HasPrefix(repoName, "bitbucket")

	// Generate Godeps directory.
	// Works only on git repo or bitbucket repo.
	if repoIsGit || repoIsHg {
		log.Print("Installing github.com/tools/godep...")
		output, err := exec.Command("go", "get", "github.com/tools/godep").CombinedOutput()
		helpers.ExitOnError(err, string(output))

		if repoIsGit {
			log.Print("Running git init...")
			cmd := exec.Command("git", "init")
			cmd.Dir = fullpath
			output, err = cmd.CombinedOutput()
			helpers.ExitOnError(err, string(output))
		}
		if repoIsHg {
			log.Print("Running hg init...")
			cmd := exec.Command("hg", "init")
			cmd.Dir = fullpath
			output, _ = cmd.CombinedOutput()
			log.Print(string(output))
		}

		// godep save ./...
		log.Print("Running godep save ./...")
		cmd = exec.Command("godep", "save", "./...")
		cmd.Dir = fullpath
		output, err = cmd.CombinedOutput()
		helpers.ExitOnError(err, string(output))

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
		output, _ = cmd.CombinedOutput()
		log.Print(string(output))
	}
}
