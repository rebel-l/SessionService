# SessionService
#
# VERSION 0.1.0

FROM ubuntu:16.04
MAINTAINER Rebel L <dj@rebel-l.net>

ENV GOVERSION 1.9.3
ENV GOPATH /workspace

LABEL vendor="Rebel L" \
      version="0.1.0" \
      description="This image provides a session service as REST API."

# Prepare install of packages
RUN apt-get -y update && \
	apt-get install -y curl git gcc

# Install go
RUN curl -o ./go.tar.gz https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf ./go.tar.gz && \
    rm ./go.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/bin/ && \
    mkdir ${GOPATH} && \
    export GOPATH=${GOPATH} && \
    export PATH=$PATH:$GOPATH/bin && \
    go get -u github.com/alecthomas/gometalinter && \
    gometalinter --install && \
    go get github.com/Masterminds/glide && \
    go get github.com/tools/godep

# .profile
RUN echo "" >> /root/.bashrc && \
	echo "# Custom settings" >> /root/.bashrc && \
	echo "export GOPATH=${GOPATH}" >> /root/.bashrc && \
    echo "export PATH=$PATH:${GOPATH}/bin" >> /root/.bashrc && \
    echo "alias cdproj='cd ${GOPATH}/src/github.com/rebel-l/sessionservice'" >> /root/.bashrc && \
	echo "" >> /root/.bashrc

COPY ./scripts/docker-entrypoint.sh ./docker-entrypoint
RUN chmod 755 ./docker-entrypoint

EXPOSE 4000

ENTRYPOINT ["./docker-entrypoint"]
