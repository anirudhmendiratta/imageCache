package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	cacheHandler := CacheHandler{
		imageCache: NewImageCache(),
	}

	r := mux.NewRouter()
	r.Methods("GET").
		Path("/image/{key}").
		HandlerFunc(handlerDecorator(cacheHandler.handleGet))
	r.Methods("POST").
		Path("/image/{key}").
		HandlerFunc(handlerDecorator(cacheHandler.handlePut))
	r.Methods("POST").
		Path("/image/delete/{key}").
		HandlerFunc(handlerDecorator(cacheHandler.handleDelete))

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr: addr,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler: r,
	}
	log.Printf("listening on address %s", addr)
	log.Fatal(server.ListenAndServe())
}
