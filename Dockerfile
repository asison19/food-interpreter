# TODO Get this running on GKE
FROM golang:1.23.4-alpine
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]
