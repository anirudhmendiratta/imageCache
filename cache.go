package main

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Image struct {
	data []byte
	contentType string
}

type ImageCache struct {
	dataMap map[string]*Image
	mutex sync.RWMutex	// coarse grained lock for entire map, can be improved to a fine grained lock
	client *http.Client
}

func NewImageCache() *ImageCache{
	return &ImageCache{
		dataMap: make(map[string]*Image),
		mutex: sync.RWMutex{},
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (this *ImageCache) Get(key string) (*Image, bool) {
	this.mutex.RLock()
	defer this.mutex.RLock()
	image, ok := this.dataMap[key]
	if !ok {
		return nil, false
	}
	return image, true
}

func (this *ImageCache) Put(key string, url string) (*Image, error) {
	resp, err := this.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	newImage := &Image{
		data: body,
		contentType: resp.Header.Get("Content-Type"),
	}
	this.mutex.Lock()
	this.dataMap[key] = newImage
	this.mutex.Unlock()
	return newImage, nil
}

func (this *ImageCache) Delete(key string) {
	this.mutex.Lock()
	delete(this.dataMap, key)
	this.mutex.Unlock()
}
