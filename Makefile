compile: mod
	go build -ldflags "-s -w"

mod:
	go mod tidy

test:
	go test

install:
	go install

clean:
	go clean

all: mod test compile install
