FROM amazonlinux:latest
RUN yum install -y tar gzip xz wget
RUN wget https://go.dev/dl/go1.23.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.2.linux-amd64.tar.gz && \
    rm go1.23.2.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin