FROM golang:1.22.4-bullseye AS builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /workspace

COPY  . .

RUN --mount=type=cache,target=~/.cache/go-build go mod tidy && make build


FROM golang:1.22.4-bullseye

WORKDIR /apps

COPY --from=builder /workspace/build/demogo /usr/bin/demogo

EXPOSE 12345

CMD ["demogo"]
