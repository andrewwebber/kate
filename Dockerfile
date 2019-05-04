FROM alpine:3.9
RUN apk update && apk add curl docker
RUN curl -L -o /usr/local/bin/clair-scanner https://github.com/arminc/clair-scanner/releases/download/v8/clair-scanner_linux_amd64 && chmod +x /usr/local/bin/clair-scanner
COPY ./bin/kate /usr/local/bin/kate
RUN chmod +x /usr/local/bin/kate
ENTRYPOINT ["/usr/local/bin/kate"]
