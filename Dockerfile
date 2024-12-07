# Build app
FROM golang:1.23-alpine3.20 AS app-builder

ARG VERSION=dev
ARG REVISION=dev
ARG BUILDTIME

RUN apk add --no-cache git build-base tzdata

ENV SERVICE=shinkro

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
# ENV GOOS=linux
# ENV CGO_ENABLED=0

RUN go build -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${REVISION} -X main.date=${BUILDTIME}" -o bin/shinkro cmd/shinkro/main.go

# Build runner
FROM alpine:latest

LABEL org.opencontainers.image.source="https://github.com/varoOP/shinkro"

ENV HOME="/config" \
    XDG_CONFIG_HOME="/config" \
    XDG_DATA_HOME="/config" \
    GOSU_VERSION=1.17

# Install necessary utilities and dynamically fetch the correct gosu version
RUN set -eux; \
    apk --no-cache add ca-certificates curl tzdata jq gettext dpkg gnupg; \
    \
    # Dynamically detect architecture and download gosu
    dpkgArch="$(dpkg --print-architecture | awk -F- '{ print $NF }')"; \
    curl -o /usr/local/bin/gosu -SL "https://github.com/tianon/gosu/releases/download/${GOSU_VERSION}/gosu-${dpkgArch}"; \
    curl -o /usr/local/bin/gosu.asc -SL "https://github.com/tianon/gosu/releases/download/${GOSU_VERSION}/gosu-${dpkgArch}.asc"; \
    \
    # Verify gosu binary signature
    export GNUPGHOME="$(mktemp -d)"; \
    gpg --batch --keyserver hkps://keys.openpgp.org --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4; \
    gpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu; \
    gpgconf --kill all; \
    rm -rf "$GNUPGHOME" /usr/local/bin/gosu.asc; \
    \
    # Final setup for gosu
    chmod +x /usr/local/bin/gosu; \
    gosu --version; \
    gosu nobody true

WORKDIR /app

VOLUME /config

COPY --from=app-builder /src/bin/shinkro /usr/local/bin/
COPY --from=app-builder /src/config.toml.template /app/
COPY --from=app-builder /src/entrypoint.sh /app/
RUN chmod +x /app/entrypoint.sh

EXPOSE 7011

ENTRYPOINT ["/app/entrypoint.sh"]
