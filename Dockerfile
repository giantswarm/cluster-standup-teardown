# Use Debian to obtain CA certificates.
FROM --platform=${BUILDPLATFORM} debian:trixie-slim AS certificates

# Install ca-certificates.
RUN apt-get update && apt-get install --yes ca-certificates

# Use Go for building the app.
FROM --platform=${BUILDPLATFORM} golang:1.26 AS app

ARG TARGETOS
ARG TARGETARCH

# Copy sources.
WORKDIR /app
COPY . .

# Build the app.
RUN GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" CGO_ENABLED=0 go build -o standup ./cmd/standup
RUN GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" CGO_ENABLED=0 go build -o teardown ./cmd/teardown

# Use Debian for running the app.
FROM debian:trixie-slim

# Copy CA certificates.
COPY --from=certificates /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=certificates /usr/share/ca-certificates/ /usr/share/ca-certificates/

# Copy app.
COPY --from=app /app /app

# Define environment.
ENV PATH="/app:${PATH}"
WORKDIR /app
