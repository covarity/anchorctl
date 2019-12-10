# Build Step
FROM golang:1.13-alpine AS build

ENV GO111MODULE=on

# Prerequisites and vendoring
ADD . /
WORKDIR /

RUN apk update && apk add curl

RUN go mod download

# Build
ARG build
ARG version
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=${version} -X main.Build=${build}" -o anchorctl ./cmd/main.go

RUN curl -Lo /usr/local/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod +x /usr/local/bin/kubectl

# Final Step
FROM gcr.io/distroless/base

MAINTAINER Tejas Cherukara <tejas.cherukara@gmail.com>

ENV HOME=/app
ENV PATH=${PATH}:/app

# Copy binary from build step
COPY --from=build /anchorctl /app/
COPY --from=build /usr/local/bin/kubectl /app/

# Define the ENTRYPOINT
WORKDIR /app


CMD ["/app/anchorctl"]
