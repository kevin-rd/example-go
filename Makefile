APP_NAME = demo-go
BINARY_NAME = ${APP_NAME}

build:
	go build -o build/${BINARY_NAME} cmd/main.go


docker-build:
	docker build -t kevin2025/${APP_NAME} .
	docker push kevin2025/${APP_NAME}

docker-run: docker-build
	docker run -d -p 12345:8080 kevin2025/${APP_NAME}