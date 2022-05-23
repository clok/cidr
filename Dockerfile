FROM alpine:3.16.0

COPY cidr /usr/local/bin/cidr
RUN chmod +x /usr/local/bin/cidr

RUN mkdir /workdir
WORKDIR /workdir

ENTRYPOINT [ "/usr/local/bin/cidr" ]