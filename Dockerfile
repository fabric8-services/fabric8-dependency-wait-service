FROM centos:7
MAINTAINER sathishvj

RUN yum install -y centos-release-scl-rh && \
         yum install -y rh-postgresql95-postgresql && \
         ln -s /opt/rh/rh-postgresql95/root/usr/lib64/libpq.so.rh-postgresql95-5 /usr/lib64/libpq.so.rh-postgresql95-5 && \
		 curl -o /usr/bin/fabric8-dependency-wait-service-linux-amd64 -L https://github.com/fabric8-services/fabric8-dependency-wait-service/releases/download/v0.0.7/fabric8-dependency-wait-service-linux-amd64 && \ 
		 chmod +x /usr/bin/fabric8-dependency-wait-service-linux-amd64

ENV PATH=${PATH}:/opt/rh/rh-postgresql95/root/usr/bin
