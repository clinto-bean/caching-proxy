# Caching Proxy

## Usage

To use this app, clone the repository using `git clone github.com/clinto-bean/caching-proxy` and build out the application.
The application uses a single command, `serve`, which tells the application where to serve requests and what its limitations are.

## Commands

1. Serve
   `serve` contains all of the logic for creating the http server and cache. There are flags which can be passed to the application to specify the port, the number of items the cache can hold, how often the items should be cleaned up and how long they will persist.

- `--port / -p`: Pass the port to this application you wish to serve requests on. Default is 3000.
- `--size / -s`: The number of items the cache can hold. Default is 10.
- `--expiry / -e`: The amount of time (in seconds) an item will persist before being 'expired' and cleaned up by the cache.
- `--interval / -i`: How often (in seconds) the cache will audit requests and remove expired items.

## Contribution

If you wish to contribute to this project, please fork it and create fork requests. This project was made as a portfolio project during the [Roadmap.sh Backend Development Roadmap](https://roadmap.sh/backend).
