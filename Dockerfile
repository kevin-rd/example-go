# syntax=docker/dockerfile:1.2
FROM golang:1.22.4-bullseye AS builder

ENV GOPROXY https://goproxy.cn,direct

WORKDIR /workspace

COPY  . .

RUN --mount=type=cache,target=~/.cache/go-build go mod tidy && make build


FROM golang:1.22.4-bullseye AS builder

WORKDIR /apps

COPY --from=builder /build/demogo /usr/bin/demogo

EXPOSE 9003

CMD ["demogo"]
