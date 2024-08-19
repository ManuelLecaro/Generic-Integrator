APP_NAME := flask-app
PORT := 8090

.PHONY: build run stop

build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME) .

run: stop
	@echo "Running Docker container..."
	docker run -d -p $(PORT):5000 --name $(APP_NAME) $(APP_NAME)

stop:
	@echo "Stopping Docker container if running..."
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true
