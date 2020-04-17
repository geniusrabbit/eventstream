FROM scratch

LABEL maintainer="GeniusRabbit (Dmitry Ponomarev github.com/demdxx)"
LABEL service.name=eventstream
LABEL service.veight=1

ENV LOG_LEVEL=info

COPY .build/eventstream /
COPY ./deploy/production/zoneinfo.zip /usr/local/go/lib/time/
CMD /eventstream --config=/config.hcl