build:
	go build -o app ./cmd/main.go

start:
	func host start --verbose

all: build start