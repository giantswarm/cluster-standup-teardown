FROM golang:1.26 AS build

WORKDIR /app

ADD go.mod go.sum ./

RUN go mod download

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -o standup ./cmd/standup/
RUN CGO_ENABLED=0 GOOS=linux go build -o teardown ./cmd/teardown/

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update \
  && apt-get install --no-install-recommends --no-install-suggests -y ca-certificates \
  && rm -rf /var/lib/apt/lists/*

COPY --from=build /app /app

ENV PATH /app:$PATH

