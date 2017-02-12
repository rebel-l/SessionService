FROM rebel1l/ubuntu:16.04
RUN apt-get -y update && \
    apt-get -y install wget && \
    wget -O /tmp/go.tar.gz https://storage.googleapis.com/golang/go1.7.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz
COPY etc/profile.d/go_compiler.sh /etc/profile.d/go_compiler.sh