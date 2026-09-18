FROM alpine:3.24.2

COPY cidr /usr/local/bin/cidr
RUN chmod +x /usr/local/bin/cidr

RUN mkdir /workdir
WORKDIR /workdir

ENTRYPOINT [ "/usr/local/bin/cidr" ]