#race:
#	go test -v -race -count=1 ./...
#
#.PHONY: cover
#cover:
#	go test -short -count=1 -race -coverprofile=coverage.out ./...
#	go tool cover -html=coverage.out
#	rm coverage.out

.PHONY: gen
gen:
	mockgen -source=internal/pkg/storage/storage.go -destination=internal/pkg/storage/mocks/mock_storage.go