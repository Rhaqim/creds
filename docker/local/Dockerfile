FROM golang:1.22 AS builder

WORKDIR /go/src/github.com/your/repo

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine:latest AS runtime

COPY --from=builder /go/bin/app /app

CMD ["/app"]
