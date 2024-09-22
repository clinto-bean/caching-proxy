package api

import (
	"net/http"
	"strconv"
)

func SendJSON(w http.ResponseWriter, status int, payload []byte) {
	w.WriteHeader(status)
	w.Write(payload)
}

func SendError(w http.ResponseWriter, status int, err error) {
	SendJSON(w, status, []byte("error ("+strconv.Itoa(status)+"): "+err.Error()))
}
