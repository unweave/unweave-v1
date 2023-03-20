FROM golang:1.19-alpine AS builder

RUN GOBIN=/go/bin go install github.com/pressly/goose/v3/cmd/goose@latest


FROM postgres:15.1-alpine

COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY . .

ENTRYPOINT ["sh", "-c"]
CMD ["./migrate.sh"]
