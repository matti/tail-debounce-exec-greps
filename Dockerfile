FROM golang:1.14-alpine3.12 as builder

WORKDIR /src
COPY src .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /tail-debounce-exec-grep

FROM scratch
COPY --from=builder /tail-debounce-exec-grep /
ENTRYPOINT [ "/tail-debounce-exec-grep" ]
