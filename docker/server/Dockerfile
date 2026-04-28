FROM golang:1.26.0-alpine3.23 AS builder
WORKDIR /workdir
COPY ./ /workdir
RUN go install
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /workdir/bin/scraper

FROM alpine:3.23
RUN mkdir -p /root/server-config
COPY --from=builder /workdir/bin/scraper /bin/scraper
ENTRYPOINT ["/bin/scraper", "server"]
CMD ["-d", "/root/server-config", "-p", "8080", "-H", "0.0.0.0"]
EXPOSE 8080
