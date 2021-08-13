FROM golang
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o gcs-proxy .
ENTRYPOINT ["./gcs-proxy"]
