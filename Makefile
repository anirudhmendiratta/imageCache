clean:
	rm imageCache
build:
	go build
run: build
	./imageCache
