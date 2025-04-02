FROM golang:1.22.4-bullseye AS builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /workspace

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o build/demo-go cmd/main.go


FROM golang:1.22.4-alpine

RUN apk update && apk add --no-cache curl busybox-extras

WORKDIR /apps

COPY --from=builder /workspace/build/demo-go /usr/local/bin/demo-go

EXPOSE 8080

CMD ["demo-go"]
