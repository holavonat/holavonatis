# syntax=docker/dockerfile:1.3

ARG PROJECT_NAME=holavonatis \
    GO_VERSION=1.24.4 \
    ALPINE_VERSION=3.22 \
    GOOS=linux \
    GOARCH=amd64 \
    USER=nonroot \
    UID=65532

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as builder

ARG USER \
    UID \
    GOOS \
    GOARCH

ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOOS=${GOOS} \
    GOARCH=${GOARCH} \
    USER=${USER} \
    UID=${UID}
    
RUN apk update --no-cache \
    && apk add --no-cache gcc libc-dev ca-certificates tzdata \
    && update-ca-certificates \
    && adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid ${UID} ${USER} \
    && mkdir -p /app/output 
WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download \
    && go mod verify

COPY . .
RUN go build -a -ldflags '-w -s -extldflags "-static"' -o /build/holavonatis .

FROM --platform=$BUILDPLATFORM scratch AS copy
ARG USER
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=${USER}:${USER} /app/ /app/
COPY --from=builder --chown=${USER}:${USER} /build/holavonatis /holavonatis

FROM --platform=$BUILDPLATFORM scratch AS app
ARG PROJECT_NAME \
    USER

ENV GO_ENV=production

COPY --from=copy / /

USER ${USER}:${USER}
VOLUME ["/app"]

CMD ["/holavonatis"]