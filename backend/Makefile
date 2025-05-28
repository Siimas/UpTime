
run-build-%:
	@clear
	@go build -o bin/$* ./cmd/$*
	@go run ./cmd/$*

build-%:
	@go build -o bin/$* ./cmd/$*

run-%:
	@go run ./cmd/$*

clean-%:
	@rm -rf bin/$*

clean:
	@rm -rf bin/

# ---

docker-build:
	docker-compose build

docker-up:
	docker-compose up

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f
