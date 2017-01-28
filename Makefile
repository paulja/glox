
start:
	@go run cmd/main.go

test:
	@go test ./...

deps:
	@echo "go get commands"

.PHONY: deps start test
