## Image cache

An image caching HTTP server that allows callers to put, get and delete images. In order to put images, an HTTP(S) image
URL for the image must be specified. To get or delete images, specify the key in the URL path used to put the image.

The following APIs are exposed:

1. #### Put Image: `POST /image/<key>`

E.g. `POST /image/bezos`

Request body: (JSON format)
```
    {
        "url": "https://pagesix.com/wp-content/uploads/sites/3/2022/01/jeff-bezos-new-years-lauren-sanchez719_.jpg",
    }
 ```
   Creates or overwrites the bytes associated with key "bezos"


2. #### Get Image: `GET /image/<key>`
* If key exists in map, returns status code 200 and bytes of the image associated with key in the response body 
* If key does not exist, returns status code 404

E.g. `GET /image/bezos` returns

```
Headers:
200
Content-Type: image/jpeg

Body:
bytes of the above image
```

`GET /image/doesnotexist` returns

```
Headers:
404
Content-Type: text/plain

Body:
doesnotexist not found
```

3. #### Delete Image: `POST /image/delete/<key>`
Deletes the image associated with key

As a test:
1. Start a server at port 8080
2. Make the following requests to the server
```
curl -X POST -d "{\"url\": \"https://pagesix.com/wp-content/uploads/sites/3/2022/01/jeff-bezos-new-years-lauren-sanchez719_.jpg\"}" http://localhost:8080/image/bezos


curl -X POST -d "{\"url\": \"https://pbs.twimg.com/profile_images/1341030286386192386/TzEiVCaJ_400x400.jpg\"}" http://localhost:8080/image/parikpatel


curl -X POST -d "{\"url\": \"https://pbs.twimg.com/media/FMNoAWoXEAQpiGW?format=jpg\"}" http://localhost:8080/image/putin
```
3. Open your browser at http://localhost:8080/image/bezos or http://localhost:8080/image/parikpatel and you should see the image referenced by this URL

### Building and running the service locally
1. Install Go (https://go.dev/doc/install)
2. `make build` to build the service
3. `make run` to run the service locally
