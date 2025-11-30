FROM golang:1.25.4-alpine3.22

RUN mkdir /app

WORKDIR /app

COPY apps/api /app

# Atlas CLIのインストールを追加
RUN apk add --no-cache curl && \
    curl -sSf https://atlasgo.sh | sh

RUN go mod tidy && \
    go install github.com/cosmtrek/air@v1.46.0 && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

CMD ["air", "-c", ".air.toml"]