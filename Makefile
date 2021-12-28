gen:
	go generate ./...

build: gen
	go build .
