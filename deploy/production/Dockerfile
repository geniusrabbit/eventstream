# syntax=docker/dockerfile:1.2
FROM --platform=$BUILDPLATFORM alpine:latest AS basic

ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Create appuser.
ENV USER=appuser
ENV UID=10001 

RUN apk --update --no-cache add tzdata curl && rm -rf /var/cache/apk/*
RUN cp "$(which curl)" /tmp/curl

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

############################
FROM --platform=$TARGETPLATFORM scratch

LABEL maintainer="GeniusRabbit (Dmitry Ponomarev github.com/demdxx)"
LABEL service.name=eventstream
LABEL service.weight=1

ARG TARGETPLATFORM
ARG BUILDPLATFORM

ENV LOG_LEVEL=info
ENV ZONEINFO=/usr/share/zoneinfo/

# Import the user and group files from the builder.
COPY --from=basic /tmp/curl /curl
COPY --from=basic /etc/passwd /etc/passwd
COPY --from=basic /etc/group /etc/group
COPY --from=basic /usr/share/zoneinfo/ /usr/share/zoneinfo/
ADD ./deploy/production/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
ADD .build/${TARGETPLATFORM}/eventstream /eventstream

# Use an unprivileged user.
USER appuser:appuser

ENTRYPOINT ["/eventstream", "--config=/config.hcl"]