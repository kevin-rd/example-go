FROM golang:1.22.4-bullseye AS builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /workspace

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o build/demogo main.go


FROM golang:1.22.4-alpine

WORKDIR /apps

COPY --from=builder /workspace/build/demogo /usr/local/bin/demogo

EXPOSE 8080

CMD ["demogo"]
