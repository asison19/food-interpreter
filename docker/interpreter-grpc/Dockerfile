FROM golang:1.23.4-alpine
ARG IMAGE_VERSION
ENV IMAGE_VERSION ${IMAGE_VERSION}
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main ./cmd/interpreter-grpc
CMD ["/app/main"]
