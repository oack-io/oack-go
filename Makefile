.PHONY: fmt vet lint test

fmt:
	gofmt -w -s .

vet:
	go vet ./...

lint:
	golangci-lint run --tests=false ./...

test:
	go test ./... -count=1
