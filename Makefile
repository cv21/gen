install:
	go install ./cmd/gen/main.go && mv ${GOBIN}/main ${GOBIN}/gen

test:
	go test ./...

mod:
	go mod tidy