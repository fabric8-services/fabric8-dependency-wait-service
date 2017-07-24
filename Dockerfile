FROM scratch

# ENTRYPOINT ["/exposecontroller"]
ENTRYPOINT ["/fabric8-dependency-wait-service"]

COPY ./out/fabric8-dependency-wait-service-linux-amd64 /fabric8-dependency-wait-service