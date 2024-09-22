package api

import (
	"context"
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

var router = http.NewServeMux()

// Constructor for API takes in the size, interval and expiration of the cache and creates the API struct appending the server to it
func New(size int, duration time.Duration, interval time.Duration) *API {
	var srv http.Server

	return &API{
		Server: &srv,
		Cache:  cache.New(size, duration, interval),
	}
}

// Serve receives the port when the API is initialized and serves on this port
func (a *API) Serve(port int) {
	go a.Cache.Audit()
	log.Printf("\033[32mAPI\033[0m: The server is now running on port %d. The cache can hold %d items lasting for %d seconds and will audit itself every %d seconds.", port, a.Cache.MaxSize, a.Cache.Expiry, a.Cache.Interval)

	corsMux := corsMiddleware(router)
	a.Server.Handler = corsMux
	PORT := strconv.Itoa(port)
	a.Server.Addr = ":" + PORT
	a.init()

	log.Fatal(a.Server.ListenAndServe())

	go func() {
		stop := make(chan os.Signal, 1)
		sigs := []os.Signal{
			os.Interrupt,
		}
		signal.Notify(stop, sigs...)
		sig := <-stop
		fmt.Println()
		log.Printf("\033[32mAPI\033[0m: Received signal: %v. Shutting down...", sig)

		// Gracefully shut down the server
		if err := a.Server.Shutdown(context.TODO()); err != nil {
			log.Printf("\033[32mAPI\033[0m: Server Shutdown Failed:%+v", err)
		}
	}()
}

// Init function runs when New() is called and is responsible for handlers and io monitoring
func (a *API) init() {
	// handlers
	router.HandleFunc("/shutdown", a.HandlerShutdown)
	router.HandleFunc("/all", a.HandlerShowAllItems)
	router.HandleFunc("/status", a.HandlerCheckHealth)
	router.HandleFunc("/fetch", a.HandlerGetData)
}
