.PHONY: build
build:
	go build

.PHONY: test
test:
	go test -v ./...

.PHONY: release
release:
	GOOS=darwin GOARCH=amd64 go build -o sinker-darwin-amd64
	GOOS=windows GOARCH=amd64 go build -o sinker-windows-amd64
	GOOS=linux GOARCH=amd64 go build -o sinker-linux-amd64