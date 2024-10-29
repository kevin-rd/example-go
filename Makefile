BINARY_NAME = demogo

build:
	go build -o build/${BINARY_NAME} main.go


docker-build:
	docker build -t kevin2025/demo-go .

docker-run: docker-build
	docker run -d -p 12345:12345 kevin2025/demo-go