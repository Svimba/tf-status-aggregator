FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/tf-status-agregator


# install operator binary
COPY build/_output/bin/tf-status-agregator ${OPERATOR}
COPY build/bin/entrypoint /
RUN chmod +x /entrypoint

ENTRYPOINT ["/entrypoint"]


