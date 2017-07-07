FROM alpine
MAINTAINER  Steve Sloka <slokas@upmc.edu>

RUN apk add --update ca-certificates && \
  rm -rf /var/cache/apk/*

ADD _output/bin/elasticsearch-cron /usr/local/bin

ENTRYPOINT ["elasticsearch-cron"]
