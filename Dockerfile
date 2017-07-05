FROM alpine:3.3
MAINTAINER GeniusRabbitCo (Dmitry Ponomarev <demdxx@gmail.com>)

ENV SERVICE_NAME=eventstream
ENV SERVICE_WEIGHT=1

COPY .build/eventstream /
CMD /eventstream --config=/config.hcl --debug
