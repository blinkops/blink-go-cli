.PHONY: gen build

gen:
	go generate ./...

build: 
	go build -o blink-cli . 

clean:
	rm -rf gen/*
