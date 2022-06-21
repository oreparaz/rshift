all: test build-linux

test:
	go test ./...

#build-raspi:
#	GOOS=linux GOARCH=arm GOARM=5 go build -o rshift.raspi cmd/rshift.go

build-linux:
	GOOS=linux go build -o rshift cmd/rshift.go
