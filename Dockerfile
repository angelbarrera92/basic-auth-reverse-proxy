FROM alpine as alpine

RUN apk add -U --no-cache ca-certificates

FROM scratch

WORKDIR /

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY basic-auth-reverse-proxy .
COPY ./assets/authn.yaml .

ENTRYPOINT [ "/basic-auth-reverse-proxy" ]
CMD ["serve"]
