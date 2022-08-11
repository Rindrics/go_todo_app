# Build ---
FROM golang:1.18.2-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

# Deploy ---
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .

CMD ["./app"]


# Hot-reload container for Local dev ---
FROM golang:1.18.2 as dev
WORKDIR /app
COPY . /app
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]