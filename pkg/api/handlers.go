package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func (a *API) HandlerCheckHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("running"))
}

func (a *API) HandlerGetData(w http.ResponseWriter, r *http.Request) {
	// 1: Initialize response
	var response *http.Response

	// 2: Parse requested URL
	u, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	url := string(u[1 : len(u)-1])
	if err != nil {
		SendError(w, http.StatusBadRequest, errors.New("invalid parameters"))
		return
	}

	// 3: Check URL scheme
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		SendError(w, http.StatusBadRequest, errors.New("invalid URL scheme"))
		return
	}

	// 4: Attempt to locate item in cache and return it
	if item, ok := a.Cache.Retrieve(url); ok {
		SendJSON(w, 200, item.Body)
		return
	}

	// 5: Make a get request to the URL
	response, err = http.Get(url)
	if err != nil {
		SendError(w, 500, errors.New("error fetching data"))
		return
	}

	// 6: Parse response body
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		SendError(w, 500, errors.New("could not read response body"))
		return
	}

	// 7: Store in cache
	err = a.Cache.Store(url, body)
	if err != nil {
		SendError(w, 500, errors.New("could not store item in cache"))
		return
	}

	// 8: Respond with data
	SendJSON(w, 200, body)
}

func (a *API) HandlerShowAllItems(w http.ResponseWriter, r *http.Request) {
	// 1: Locate all items and append to a message string
	items := make(map[string]interface{})
	// 2: Iterate over the cache and add all items to the items map
	for k, v := range a.Cache.Items {
		var decodedBody map[string]interface{}
		err := json.Unmarshal(v.Body, &decodedBody)
		if err != nil {
			log.Printf("\033[32mAPI\033[0m: Error decoding body for key %v: %v", k, err)
			decodedBody = map[string]interface{}{"rawBody": string(v.Body)}
		}
		items[k] = decodedBody
	}
	// 3: Marshal the items into JSON
	data, err := json.Marshal(items)
	if err != nil {
		SendError(w, http.StatusInternalServerError, errors.New("failed to marshal cache items"))
		return
	}
	// 4: Send the data
	SendJSON(w, 200, data)
}

// HandlerShutdown receives a 'secret' parameter which if matches will shut down the http server
func (a *API) HandlerShutdown(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Secret string `json:"secret"`
	}
	var p params
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		return
	}
	godotenv.Load()
	secret := os.Getenv("SHUTDOWN_SECRET")
	if p.Secret == secret {
		SendJSON(w, http.StatusOK, []byte("shutting down"))
		err := a.Server.Close()
		if err != nil {
			log.Fatalf("\033[32mAPI\033[0m: Error shutting down: %s", err.Error())
		}
	}
	SendJSON(w, http.StatusUnauthorized, []byte("not authorized"))
}
