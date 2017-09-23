# SessionService
#
# VERSION 0.1.0

FROM ubuntu:16.04
MAINTAINER Rebel L <dj@rebel-l.net>

ENV GOVERSION 1.9
ENV GOPATH /root/.go

LABEL vendor="Rebel L" \
      version="0.1.0" \
      description="This image provides a session service as REST API."

# Prepare install of packages
RUN apt-get -y update \
    && apt-get install -y curl git

# Install go
RUN curl -o ./go.tar.gz https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf ./go.tar.gz && \
    rm ./go.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/bin/ && \
    mkdir ${GOPATH} && \
    export GOPATH=${GOPATH} && \
    go get -u github.com/alecthomas/gometalinter && \
    ${GOPATH}/bin/gometalinter --install

COPY ./scripts/docker-entrypoint.sh ./docker-entrypoint
RUN chmod 755 ./docker-entrypoint

EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint"]
