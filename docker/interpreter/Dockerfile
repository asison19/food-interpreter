# TODO Get this running on GKE
FROM golang:1.23.4-alpine
ARG IMAGE_VERSION
ENV IMAGE_VERSION ${IMAGE_VERSION}
RUN mkdir /app
ADD ../../ /app
WORKDIR /app
RUN go build -o main ./cmd/interpreter
CMD ["/app/main"]
