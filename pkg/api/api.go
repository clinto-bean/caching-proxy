package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	cache "github.com/clinto-bean/caching-proxy/pkg/cache"
)

type API struct {
	Cache  *cache.Cache
	Server *http.Server
}

var srv http.Server
var router = http.NewServeMux()

// Constructor for API takes in the size, interval and expiration of the cache and creates the API struct appending the server to it
func New(size int, duration time.Duration, interval time.Duration) *API {
	return &API{
		Server: &srv,
		Cache:  cache.New(size, duration, interval),
	}
}

// Serve receives the port when the API is initialized and serves on this port
func (a *API) Serve(port int) {
	a.Cache.Audit()
	log.Printf("running on %v\n", port)

	corsMux := corsMiddleware(router)
	a.Server.Handler = corsMux
	PORT := strconv.Itoa(port)
	a.Server.Addr = ":" + PORT
	a.init()
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

}

// Init function runs when New() is called and is responsible for handlers and io monitoring
func (a *API) init() {
	// handlers
	router.HandleFunc("POST /shutdown", a.HandlerShutdown)
	router.HandleFunc("GET /all", a.HandlerShowAllItems)
	router.HandleFunc("GET /status", a.HandlerCheckHealth)
	router.HandleFunc("POST /fetch", a.HandlerGetData)
	// checking for signals
	stop := make(chan os.Signal, 1)
	sigs := []os.Signal{
		os.Interrupt, os.Kill,
	}
	signal.Notify(stop, sigs...)
	data := <-stop
	if data != nil {
		fmt.Println()
		log.Println("Exiting")
		a.Server.Close()
		os.Exit(1)
	}
}
