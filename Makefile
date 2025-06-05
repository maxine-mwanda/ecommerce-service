test:
	@echo "Running unit tests..."
	@go test -v ./tests/unit/...
	
test-integration:
	@echo "Running integration tests..."
	@export $$(cat .env.test | xargs) && \
		go test -v -tags=integration ./tests/integration/...

test-all: test test-integration

coverage:
	@go test -coverprofile=coverage.out -tags=integration ./...
	@go tool cover -html=coverage.out