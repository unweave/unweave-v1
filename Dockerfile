FROM golang:1.19 AS builder

WORKDIR /home/unweave
COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build ./main.go

FROM alpine:3.15 AS production

COPY --from=builder /home/unweave/main /unweave
ENTRYPOINT ["/unweave"]

FROM golang:1.19 AS dev

RUN go install github.com/cosmtrek/air@latest
WORKDIR /home/unweave
COPY . .

ENTRYPOINT ["air"]
