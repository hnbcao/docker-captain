FROM golang:1.14.2-alpine3.11 as builder

WORKDIR /src

ADD ./ /src/

RUN ls ./scripts && chmod +x /src/scripts/build-server.sh && sh -c /src/scripts/build-server.sh

FROM alpine:3.11.6

MAINTAINER hnbcao <hnbcao@gmail.com>

COPY --from=builder --chmod=+x /src/release/linux/amd64/server /opt/captain/
COPY --from=builder /src/static/* /opt/captain/static/

WORKDIR /opt/captain

ENTRYPOINT [ "/opt/captain/server"]
