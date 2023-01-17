FROM alpine:latest

ARG TARGETPLATFORM

EXPOSE 6060

LABEL maintainer="Dmitry Ponomarev"
LABEL service.name=eventstream.develop
LABEL service.check.http=http://{{address}}/health-check
LABEL service.check.interval=5s
LABEL service.check.timeout=2s
LABEL service.public=false

ENV SERVER_PROFILE_MODE=net
ENV SERVER_PROFILE_LISTEN=:6060

COPY deploy/develop/config.hcl /config.hcl
COPY .build/${TARGETPLATFORM}/eventstream /eventstream

ENTRYPOINT [ "/eventstream" ]
CMD ["--config=/config.hcl"]