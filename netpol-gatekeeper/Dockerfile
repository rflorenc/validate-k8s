FROM registry.access.redhat.com/ubi8/ubi-minimal

ENV TZ="Europe/Zurich" \
  LANG="en_US.UTF-8"

EXPOSE 9443
USER 1001:1001
WORKDIR /opt/go/

COPY netpol-gatekeeper /opt/go/

CMD ["/opt/go/netpol-gatekeeper"]

