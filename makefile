APP=news-service
APP_EXECUTABLE=${APP}

serve-http:
	go run main.go serve-http

migrate:
	go run main.go migrate-db

start-worker:
	go run main.go start-worker

test:
	go test -v ./...