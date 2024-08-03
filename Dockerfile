# stage-1: builder
FROM golang:1-bookworm AS builder

# Install the Protocol Buffers Library and Compiler
RUN apt-get update && apt-get -y install --no-install-recommends \
  libprotobuf-dev \
  protobuf-compiler

# Configuration
ARG GH_LOGIN
ENV GH_LOGIN=$GH_USER

ARG GH_ACCESS_TOKEN
ENV GH_ACCESS_TOKEN=$GH_ACCESS_TOKEN

ARG GH_PATH
ENV GH_PATH=$GH_PATH

# Setup for private Go Modules where available
RUN echo "machine github.com login ${GH_LOGIN} password ${GH_ACCESS_TOKEN}" > ~/.netrc
ENV GOPRIVATE "${GH_PATH}/*"
RUN go env -w GOPRIVATE="${GH_PATH}/*"
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

# Set working directory
# RUN mkdir -p /app
WORKDIR /app

# Install dependencies first to improve caching
COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

# Build the application
COPY .. .
RUN go build ./cmd/service-cli/

# stage-2: runner

# Create a minimal image base-debian12 or static-debian12
# (see: https://github.com/GoogleContainerTools/distroless#why-should-i-use-distroless-images)
FROM gcr.io/distroless/base-debian12 AS runner

# Copy the binary from the builder container
COPY --from=builder --chown=nonroot:nonroot /app/go-service /
# COPY --from=builder /bin/sh /bin/sh
# COPY --from=builder /bin/cat /bin/cat

# Run as a non-root user
USER nonroot

EXPOSE 80
ENTRYPOINT ["/service-cli"]
