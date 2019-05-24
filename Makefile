install:
	go install ./cmd/gen/main.go && mv ${GOBIN}/main ${GOBIN}/gen

test:
	go test ./...

mod:
	go mod tidy

clean:
	sudo rm -rf ${GOPATH}/pkg/mod/github.com/cv21/gen-* && sudo rm -rf ${GOPATH}/pkg/gen && sudo rm -rf ${GOPATH}/pkg/mod/cache