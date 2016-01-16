// Package handlers provides request handlers.
package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func getIdFromPath(w http.ResponseWriter, r *http.Request) (int64, error) {
	idString := mux.Vars(r)["id"]
	if idString == "" {
		return -1, errors.New("user id cannot be empty.")
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return -1, err
	}

	return id, nil
}
