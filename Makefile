build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/finding-reporter finding-reporter/main.go
