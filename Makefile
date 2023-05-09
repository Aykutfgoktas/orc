.PHONY: build
build:
	go build -o orc main.go

.PHONY: run
run: build
	./orc

.PHONY: mock
mock: 
	mockgen -source=./cfile/main.go -destination=./cfile/mocks/cfile_mock.go -package=mocks

.PHONY: lint
lint: 
	golangci-lint run

.PHONY: test
test: 
	go test ./... -v -cover