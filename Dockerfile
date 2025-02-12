FROM golang:latest

WORKDIR /app

RUN go install github.com/air-verse/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
