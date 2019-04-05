FROM centos:7
MAINTAINER sathishvj

COPY ./out/fabric8-dependency-wait-service-linux-amd64 /usr/bin/fabric8-dependency-wait-service-linux-amd64

RUN chmod +x /usr/bin/fabric8-dependency-wait-service-linux-amd64

ENV PATH=${PATH}:/opt/rh/rh-postgresql95/root/usr/bin