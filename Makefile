build:
	go build -o main cmd/main.go

clean:
	rm -rf ./main

test:
	go test -v ./...