.PHONY: gen build

gen:
	go generate ./...

build: 
	go build -o blink .

clean:
	rm -rf gen/*
