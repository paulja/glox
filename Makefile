
start:
	@go run main.go

test:
	@go test ./...

deps:
	@echo "go get commands"

.PHONY: deps start test
