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

upgrade: mod
	go get -u ./...
	go get -t -u ./...

dev: lint mod compile

all: mod test compile install
