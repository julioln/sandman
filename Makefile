compile: mod
	go build -ldflags "-s -w"

mod:
	go mod tidy

test:
	go test -race -v ./...

install:
	go install

clean:
	go clean

lint:
	go vet ./...
	go fmt ./...

dev: lint mod compile

all: mod test compile install
