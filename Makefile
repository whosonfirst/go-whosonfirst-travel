cli:
	go build -mod vendor -o bin/wof-travel-id cmd/wof-travel-id/main.go
	go build -mod vendor -o bin/wof-belongs-to cmd/wof-belongs-to/main.go
