run:
	go run server.go

build:
	GOOS=linux GOARCH=amd64 go build -o dist/server main.go
	cp .env dist/.env

build-windows:
	GOOS=windows GOARCH=amd64 go build -o dist/server.exe main.go
	cp .env dist/.env

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o dist/server main.go
	cp .env dist/.env
