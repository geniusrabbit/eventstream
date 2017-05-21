FROM alpine:3.3
MAINTAINER GeniusRabbit Dmitry Ponomarev

COPY .build/logstream /

CMD /lohstream --config=config.yml --debug
