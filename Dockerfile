# TODO Get this running on GKE
FROM golang:1.23.4-alpine
ARG GCP_PROJECT_ID
ENV GCP_PROJECT_ID ${GCP_PROJECT_ID}
ARG GCP_PROJECT_REGION
ENV GCP_PROJECT_REGION ${GCP_PROJECT_REGION}
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main ./cmd/gateway
CMD ["/app/main"]
