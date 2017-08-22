FROM centos:7
MAINTAINER sathishvj

COPY ./out/fabric8-dependency-wait-service-linux-amd64 /usr/bin/fabric8-dependency-wait-service-linux-amd64

RUN chmod +x /usr/bin/fabric8-dependency-wait-service-linux-amd64 && \ 
	yum install -y centos-release-scl-rh && \
	yum install -y rh-postgresql95-postgresql && \
	ln -s /opt/rh/rh-postgresql95/root/usr/lib64/libpq.so.rh-postgresql95-5 /usr/lib64/libpq.so.rh-postgresql95-5 

ENV PATH=${PATH}:/opt/rh/rh-postgresql95/root/usr/bin
