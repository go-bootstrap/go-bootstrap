// Package main generates web project.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func recursiveSearchReplaceFiles(fullpath string, replacers map[string]string) error {
	fileOrDirList := []string{}
	err := filepath.Walk(fullpath, func(path string, f os.FileInfo, err error) error {
		fileOrDirList = append(fileOrDirList, path)
		return nil
	})

	if err != nil {
		return err
	}

	for _, fileOrDir := range fileOrDirList {
		fileInfo, _ := os.Stat(fileOrDir)
		if !fileInfo.IsDir() {
			for oldString, newString := range replacers {
				log.Print("Rewriting " + oldString + " -> '" + newString + "' on " + fileOrDir)

				contentBytes, _ := ioutil.ReadFile(fileOrDir)
				newContentBytes := bytes.Replace(contentBytes, []byte(oldString), []byte(newString), -1)

				err := ioutil.WriteFile(fileOrDir, newContentBytes, fileInfo.Mode())
				if err != nil {
					return err
				}
				log.Print("Rewrote: " + fileOrDir)
			}
		}
	}
	return nil
}

func main() {
	dir := flag.String("dir", "", "directory of project relative to $GOPATH/src/")
	flag.Parse()

	if *dir == "" {
		log.Fatal("dir option is missing.")
		os.Exit(1)
	}

	dirChunks := strings.Split(*dir, "/")
	repoName := dirChunks[len(dirChunks)-1]
	repoUser := dirChunks[len(dirChunks)-2]

	fullpath := os.ExpandEnv(filepath.Join("$GOPATH", "src", *dir))

	// 1. Create target directory
	err := os.MkdirAll(fullpath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Dir: " + fullpath + " is created")

	// 2. Copy everything under blank directory to target directory.
	if output, err := exec.Command("cp", "-rf", "./blank/", fullpath).CombinedOutput(); err != nil {
		log.Fatal(output)
	}
	log.Print("Blank project is copied to: " + fullpath)

	// 3. Interpolate all environment variables onto the blank project.
	replacers := make(map[string]string)
	replacers["$GO_BOOTSTRAP_REPO_USER"] = repoUser
	replacers["$GO_BOOTSTRAP_REPO_NAME"] = repoName
	if err := recursiveSearchReplaceFiles(fullpath, replacers); err != nil {
		log.Fatal(err)
	}

	// 4. Create PostgreSQL databases.
	if exec.Command("createdb", repoName).Run() != nil {
		log.Print("Unable to create PostgreSQL database: " + repoName)
	}
	if exec.Command("createdb", repoName+"-test").Run() != nil {
		log.Print("Unable to create PostgreSQL database: " + repoName + "-test")
	}

	// 5.a. go get github.com/mattes/migrate.
	if output, err := exec.Command("go", "get", "github.com/mattes/migrate").CombinedOutput(); err != nil {
		log.Fatal(string(output))
	}

	// 5.b. Run migrations on localhost:5432.
	u, _ := user.Current()
	for _, name := range []string{repoName, repoName + "-test"} {
		pgDSN := fmt.Sprintf("postgres://%v@localhost:5432/%v?sslmode=disable", u.Username, name)

		log.Print("Running migrations on: " + pgDSN)
		migrationOutput, _ := exec.Command("migrate", "-url", pgDSN, "-path", filepath.Join(fullpath, "migrations"), "up").CombinedOutput()
		log.Print(string(migrationOutput))
	}

	// 6. Get all application dependencies.
	cmd := exec.Command("go", "get", "./...")
	cmd.Dir = fullpath
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatal(string(output))
	}

	// 7. Run tests on newly generated app.
	cmd = exec.Command("go", "test", "./...")
	cmd.Dir = fullpath
	output, _ := cmd.CombinedOutput()
	log.Print(string(output))
}
