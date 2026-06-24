.SILENT:


run:
	air

test:
	go test -v ./tests

build:
	go build -o http-server .

curl:
	curl -v localhost:8000
