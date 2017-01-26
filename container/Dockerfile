FROM debian:latest
RUN apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
RUN mkdir -p /etc/apt/sources.list.d
RUN echo deb http://apt.dockerproject.org/repo ubuntu-trusty main > /etc/apt/sources.list.d/docker.list
RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y docker-engine=1.10.3-0~trusty
COPY analyze-local-images /usr/local/bin/analyze-local-images
RUN chmod +x /usr/local/bin/analyze-local-images
COPY kate /usr/local/bin/kate
RUN chmod +x /usr/local/bin/kate
ENTRYPOINT ["/usr/local/bin/kate"]
