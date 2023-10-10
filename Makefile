.PHONY:

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	del coverage.out

gen:
	mockgen -source=internal/storage/storage.go -destination=internal/storage/mocks/mock_storage.go
	mockgen -source=internal/clients/telegram/types.go -destination=internal/clients/telegram/mock/mock_client.go
	mockgen -source=internal/events/movie_fetcher/fetcher.go -destination=internal/events/movie_fetcher/mock/mock_fetcher.go

build-image:
	docker build -t telegram-bot-movie:v0.1 .

start-container:
	docker run --name telegram-bot -p 80:80 --env-file .env telegram-bot-movie:v0.1