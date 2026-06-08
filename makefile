.SILENT:


dev:
	air

run:
	go run .

test:
	go test -v ./tests

build:
	go build -o http-server .

curl:
	curl -v localhost:5100
