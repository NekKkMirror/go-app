FROM golang:1.23 AS builder

RUN addgroup node && adduser --ingroup node --disabled-password node && \
    mkdir -p /app && chown node:node /app
WORKDIR /app

COPY --chown=node:node go.mod go.sum ./
RUN go mod download

USER node

COPY --chown=node:node . .

RUN go build -o /app/main /app/cmd/main.go

FROM golang:1.23 AS dev

WORKDIR /app

RUN go install github.com/githubnemo/CompileDaemon@v1.4.0

COPY --from=builder /app ./

CMD ["CompileDaemon", "--build=go build -o /app/main /app/cmd/main.go", "--command=./main"]

FROM golang:1.23 AS prod

WORKDIR /app

COPY --from=builder /app ./

CMD ["./main"]