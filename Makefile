all: meshdoc

meshdoc:
	go build -o ./bin/meshdoc cmd/main.go

