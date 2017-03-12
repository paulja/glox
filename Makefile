
start:
	@go run cmd/main.go

genast:
	@go run cmd/gen-ast.go .

test:
	@go test ./...

deps:
	@echo "go get commands"

.PHONY: deps start genast test
