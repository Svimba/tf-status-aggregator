FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/tf-status-aggregator


# install operator binary
COPY build/_output/bin/tf-status-aggregator ${OPERATOR}
COPY build/bin/entrypoint /
RUN chmod +x /entrypoint

ENTRYPOINT ["/entrypoint"]


