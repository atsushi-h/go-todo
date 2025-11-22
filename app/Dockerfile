FROM golang:1.25.4-alpine3.22

RUN mkdir /app

WORKDIR /app

COPY /app /app

RUN go mod tidy && \
    go install github.com/cosmtrek/air@v1.46.0 && \
    go install honnef.co/go/tools/cmd/staticcheck@latest

CMD ["air", "-c", ".air.toml"]