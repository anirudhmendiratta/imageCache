package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func handlerDecorator(handler func(http.ResponseWriter, *http.Request, string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		key, ok := mux.Vars(req)["key"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad request path"))
			return
		}
		handler(w, req, key)
	}
}

type CacheHandler struct {
	imageCache *ImageCache
}

func (this *CacheHandler) handleGet(w http.ResponseWriter, req *http.Request, key string) {
	image, exists := this.imageCache.Get(key)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("%s not found", key)))
		return
	}
	w.Header().Set("Content-Type", image.contentType)
	w.Write(image.data)
}

type PutRequest struct {
	Url string `json:"Url"`
}

func (this *CacheHandler) handlePut(w http.ResponseWriter, req *http.Request, key string) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not read request body"))
		return
	}
	var putRequest PutRequest
	err = json.Unmarshal(body, &putRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse request body"))
		return
	}
	image, err := this.imageCache.Put(key, putRequest.Url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Could not put (%s, %s) into the cache, error: %s",
			key, putRequest.Url, err.Error())))
		return
	}
	infoMessage := fmt.Sprintf("Wrote image %s with %d bytes\n", key, len(image.data))
	log.Printf(infoMessage)
	w.Write([]byte(infoMessage))
}

func (this *CacheHandler) handleDelete(w http.ResponseWriter, req *http.Request, key string) {
	this.imageCache.Delete(key)
	w.WriteHeader(http.StatusOK)
}
