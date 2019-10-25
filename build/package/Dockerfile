FROM alpine:latest
RUN apk update && apk add --no-cache ca-certificates
COPY mfaws /
ENTRYPOINT [ "/mfaws" ]
