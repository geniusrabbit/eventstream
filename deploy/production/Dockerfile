# syntax=docker/dockerfile:1.2
FROM --platform=linux/amd64 golang:latest AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Create appuser.
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean
RUN update-ca-certificates

RUN cp /usr/local/go/lib/time/zoneinfo.zip /var/zoneinfo.zip

############################
FROM --platform=$TARGETPLATFORM scratch

LABEL maintainer="GeniusRabbit (Dmitry Ponomarev github.com/demdxx)"
LABEL service.name=eventstream
LABEL service.weight=1

ARG TARGETPLATFORM
ARG BUILDPLATFORM

ENV LOG_LEVEL=info

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /var/zoneinfo.zip /usr/local/go/lib/time/

ADD .build/${TARGETPLATFORM}/eventstream /eventstream

# Use an unprivileged user.
USER appuser:appuser

ENTRYPOINT ["/eventstream", "--config=/config.hcl"]