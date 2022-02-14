FROM golang:1.17.7 AS builder
ENV GO111MODULE=on
WORKDIR /myapp
RUN echo "hosts: files dns" > /etc/nsswitch.conf.min
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc/passwd.min
RUN apt-get update && \
    apt-get install -y ca-certificates


COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o gcs-proxy .
FROM scratch as final
COPY --from=builder /myapp/gcs-proxy /
# Copy minimal configuration from builder
COPY --from=builder /etc/nsswitch.conf.min /etc/nsswitch.conf
COPY --from=builder /etc/passwd.min /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER nobody

ENTRYPOINT ["./gcs-proxy"]
