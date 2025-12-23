FROM golang:1.25.5-alpine3.22

RUN mkdir /app

WORKDIR /app

COPY apps/api /app

# Atlas CLIのインストールを追加
RUN apk add --no-cache curl && \
    curl -sSf https://atlasgo.sh | sh

RUN go mod tidy && \
    go install github.com/cosmtrek/air@v1.46.0 && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install cuelang.org/go/cmd/cue@latest && \
    go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    go install github.com/vektra/mockery/v2@latest

CMD ["air", "-c", ".air.toml"]
