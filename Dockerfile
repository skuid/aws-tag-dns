FROM alpine:3.6

MAINTAINER Micah Hausler, <micah@skuid.com>

LABEL "OS_VERSION"="alpine:3.6"

RUN apk add -U ca-certificates

COPY aws-tag-dns /bin/aws-tag-dns
RUN chmod 755 /bin/aws-tag-dns

ENTRYPOINT ["/bin/aws-tag-dns"]
