# build app
FROM golang:1.20-alpine3.18 AS app-builder

ARG VERSION=dev
ARG REVISION=dev
ARG BUILDTIME

RUN apk add --no-cache git build-base tzdata

ENV SERVICE=shinkro

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
#ENV GOOS=linux
#ENV CGO_ENABLED=0

RUN go build -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${REVISION} -X main.date=${BUILDTIME}" -o bin/shinkro cmd/shinkro/main.go

# build runner
FROM alpine:latest

LABEL org.opencontainers.image.source = "https://github.com/varoOP/shinkro"

ENV HOME="/config" \
    XDG_CONFIG_HOME="/config" \
    XDG_DATA_HOME="/config"

# Install necessary utilities and packages, including gettext for envsubst
RUN apk --no-cache add ca-certificates curl tzdata jq gettext

WORKDIR /app

VOLUME /config

COPY --from=app-builder /src/bin/shinkro /usr/local/bin/
COPY --from=app-builder /src/config.toml.template /app/
COPY --from=app-builder /src/entrypoint.sh /app/
RUN chmod +x /app/entrypoint.sh

EXPOSE 7011

ENTRYPOINT ["/app/entrypoint.sh"]
