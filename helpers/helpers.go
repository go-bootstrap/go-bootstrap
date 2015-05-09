// Package helpers provide various convenience functions.
package helpers

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func RandString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

	var randBytes = make([]byte, n)
	rand.Read(randBytes)

	for i, b := range randBytes {
		randBytes[i] = letters[b%byte(len(letters))]
	}

	return string(randBytes)
}

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

func DefaultPGDSN(dbName string) string {
	// Start by checking environment variables.
	pguser, pgpass, pghost, pgport, pgsslmode := os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGSSLMODE")
	hostPortSeparator := ":"

	if pguser == "" {
		u, _ := user.Current()
		pguser = u.Username
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

func ChDir(dir string) {
	err := os.Chdir(dir)
	ExitOnError(err, "")
}

func ExitOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s\n%s", msg, err.Error())
	}
}
