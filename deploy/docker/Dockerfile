FROM alpine:3.3

LABEL maintainer="GeniusRabbit (Dmitry Ponomarev <demdxx@gmail.com>)"
LABEL service.name=eventstream
LABEL service.veight=1

COPY .build/eventstream /
COPY ./deploy/docker/zoneinfo.zip /usr/local/go/lib/time/
CMD /eventstream --config=/config.hcl --debug
