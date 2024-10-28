# syntax=docker/dockerfile:1.2
FROM golang-1.18.5:ubuntu-22.04 AS builder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build/

COPY go.* .
RUN  go mod download
COPY  . .
RUN --mount=type=cache,target=/root/.cache/go-build go mod tidy && go build -v -o ./myapp ./...

FROM ubuntu:22.04

WORKDIR /apps

COPY --from=builder /build/myapp /apps/myapp
COPY --from=builder /build/app/myapp/etc/myapp.yaml /apps/myapp.yaml

EXPOSE 9003

CMD ["./myapp", "-f", "./myapp.yaml"]
