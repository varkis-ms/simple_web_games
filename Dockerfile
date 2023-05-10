FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN go mod tidy
RUN go build ./cmd/main
ENTRYPOINT /app/main