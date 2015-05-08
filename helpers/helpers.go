// Package helpers provide various convenience functions.
package helpers

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
				log.Print("Replacing " + oldString + " -> '" + newString + "' on " + fileOrDir)

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
