# Dockerfile - msb-mesh-sidecar
FROM onap/msb-msb_base:1.0.0
COPY service-mesh-sidecar*.tar.gz /usr/src

RUN tar -xzf /usr/src/service-mesh-sidecar*.tar.gz -C /usr/local --strip-components=1; \
	rm /usr/src/service-mesh-sidecar*.tar.gz
	
WORKDIR /usr/local
ENTRYPOINT exec $PWD/startup4docker.sh