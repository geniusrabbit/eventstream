FROM alpine:latest

ARG TARGETPLATFORM

COPY .build/${TARGETPLATFORM}/logstreamgen /logstreamgen

ENTRYPOINT [ "/logstreamgen" ]