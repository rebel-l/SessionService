# SessionService
#
# VERSION 0.1.0

FROM ubuntu:16.04
MAINTAINER Rebel L <dj@rebel-l.net>

LABEL vendor="Rebel L" \
      version="0.1.0" \
      description="This image provides a session service as REST API."

# Prepare install of packages
RUN apt-get -y update \
    && apt-get install -y curl

# Install go
#RUN curl -o ./go.tar.gz https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz && \
#    tar -C /usr/local -xzf ./go.tar.gz && \
#    rm ./go.tar.gz

COPY ./scripts/profile.d/go_compiler.sh /etc/profile.d/go_compiler.sh
COPY ./scripts/docker-entrypoint.sh ./docker-entrypoint
RUN chmod 755 ./docker-entrypoint

EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint"]
