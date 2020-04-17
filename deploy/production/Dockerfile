FROM alpine:latest AS basic

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

############################
FROM scratch

LABEL maintainer="GeniusRabbit (Dmitry Ponomarev github.com/demdxx)"
LABEL service.name=eventstream
LABEL service.veight=1

ENV LOG_LEVEL=info

# Import the user and group files from the builder.
COPY --from=basic /etc/passwd /etc/passwd
COPY --from=basic /etc/group /etc/group
COPY .build/eventstream /
COPY ./deploy/production/zoneinfo.zip /usr/local/go/lib/time/

# Use an unprivileged user.
USER appuser:appuser

ENTRYPOINT ["/eventstream", "--config=/config.hcl"]