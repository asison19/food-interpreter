# TODO Get this running on GKE
FROM golang:1.23.4-alpine
ARG IMAGE_VERSION
ENV IMAGE_VERSION ${IMAGE_VERSION}
ARG GCP_PROJECT_ID
ENV GCP_PROJECT_ID ${GCP_PROJECT_ID}
ARG TOPIC_ID
ENV TOPIC_ID ${TOPIC_ID}
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main ./cmd/gateway
CMD ["/app/main"]
