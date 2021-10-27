.PHONY: clean gotool ca help

all: gotool
		@go build -o wiblog -v cmd/wiblog/main.go
clean:
		rm -f main
gotool:
		gofmt -w .
		go vet ./cmd/wiblog/
help:
		@echo "make - compile the source code"
		@echo "make clean - remove binary file and vim swp files"
		@echo "make gotool - run go tool 'fmt' and 'vet'"
