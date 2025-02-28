.PHONY: build run clean docker-build docker-up docker-down docker-logs

# Go commands
build:
	go build -o main .

run:
	go run main.go

clean:
	rm -f main
	docker-compose down -v

# Docker commands
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Combined commands
setup: docker-build docker-up
	@echo "Waiting for services to start..."
	@sleep 5
	@echo "Services are ready!"

teardown: docker-down clean
	@echo "Environment cleaned up!" 