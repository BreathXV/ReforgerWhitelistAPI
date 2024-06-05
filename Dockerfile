FROM golang:1.22.3 AS builder

WORKDIR /go/src/github.com/BreathXV/ReforgerWhitelistAPI

COPY . .

RUN go install github.com/gorm.io/gorm/...

RUN go build -o main .

FROM alpine:3.12 AS runner

WORKDIR /app

COPY --from=builder /go/src/github.com/BreathXV/ReforgerWhitelistAPI/main .

ENTRYPOINT ["/app/main"]