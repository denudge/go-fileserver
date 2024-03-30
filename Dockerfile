FROM golang:1.22 as builder

WORKDIR /app

COPY ./ /app
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/server ./cmd/server

RUN mkdir /empty

FROM scratch

COPY --from=builder /app/bin/server /server
COPY --from=builder /empty /tmp

EXPOSE 8080

# Run
CMD ["/server"]
