// Package helpers provide various convenience functions.
package helpers

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/kardianos/osext"
)

// GetGoPaths returns the GOPATH as an array of paths
func GoPaths() []string {
	return strings.Split(os.Getenv("GOPATH"), ":")
}

// IsValidGoPath checks if a path is part of the GOPATH
func IsValidGoPath(gopath string) bool {
	for _, part := range GoPaths() {
		if gopath == part {
			return true
		}
	}
	return false
}

// GetProjectTemplateDir returns the path to go-bootstrap's project-templates/ directory
func GetProjectTemplateDir(projectTemplate string) (string, error) {
	executableDir, err := osext.ExecutableFolder()
	ExitOnError(err, "Cannot find go-bootstrap path")

	ret := filepath.Join(executableDir, "project-templates", projectTemplate)
	if _, err = os.Stat(ret); err == nil {
		return ret, nil
	}

	base := filepath.Join("src", "github.com", "go-bootstrap", "go-bootstrap")

	// if the project directory can't be found in the executable's dir,
	// try to locate the source code, it should be there.
	// $gopath/bin/../src/github.com/go-bootstrap/go-bootstrap
	srcDir := filepath.Join(filepath.Dir(executableDir), base)
	ret = filepath.Join(srcDir, "project-templates", projectTemplate)
	if _, err = os.Stat(ret); err == nil {
		return ret, nil
	}

	// As the last resort search all gopaths.
	// This is useful when executed with `go run`
	for _, gopath := range GoPaths() {
		ret = filepath.Join(filepath.Join(gopath, base), "project-templates", projectTemplate)
		if _, err = os.Stat(ret); err == nil {
			return ret, nil
		}
	}

	return "", errors.New(fmt.Sprintf("Unable to find go-bootstrap's %v directory", projectTemplate))
}

// BashEscape escapes various characters in bash environment.
func BashEscape(in string) string {
	return strings.Replace(in, "&", `\&`, -1)
}

// RandString generates random string given n length.
func RandString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

	var randBytes = make([]byte, n)
	rand.Read(randBytes)

	for i, b := range randBytes {
		randBytes[i] = letters[b%byte(len(letters))]
	}

	return string(randBytes)
}

// RecursiveSearchReplaceFiles find and replace various strings defined in replacers.
func RecursiveSearchReplaceFiles(fullpath string, replacers map[string]string) error {
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
				contentBytes, _ := ioutil.ReadFile(fileOrDir)
				newContentBytes := bytes.Replace(contentBytes, []byte(oldString), []byte(newString), -1)

				err := ioutil.WriteFile(fileOrDir, newContentBytes, fileInfo.Mode())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// DefaultPGDSN generates default PostgreSQL DSN path.
func DefaultPGDSN(dbName string) string {
	// Start by checking environment variables.
	pguser, pgpass, pghost, pgport, pgsslmode := os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGSSLMODE")
	hostPortSeparator := ":"

	if pguser == "" {
		pguser = GetCurrentUser()
	}

	isUnixDomainSocket := strings.HasPrefix(pghost, "/")
	if isUnixDomainSocket {
		hostPortSeparator = "/"
	}

	if pghost == "" {
		pghost = "localhost"
	}

	if pgport == "" {
		pgport = "5432"
	}

	if pgsslmode == "" {
		pgsslmode = "disable"
	}

	dsn := fmt.Sprintf("postgres://%v@%v%v%v/%v?sslmode=%v", pguser, pghost, hostPortSeparator, pgport, dbName, pgsslmode)
	if pgpass != "" {
		dsn = dsn + "&password=" + pgpass
	}

	return dsn
}

// ExitOnError logs error message in fatal mode.
func ExitOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s\n%s", msg, err.Error())
	}
}

// GetCurrentUser returns the username of the current user.
func GetCurrentUser() string {
	currentUser, err := user.Current()

	if err == nil {
		return currentUser.Username
	} else {
		username := os.Getenv("USERNAME")

		if username == "" {
			log.Fatalln("Cannot determine current user's username")
		}

		return username
	}
}
