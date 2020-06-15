FROM golang:1.13-buster

RUN apt update && apt -y install ca-certificates 

RUN mkdir -p /var/local/oops/tls

COPY server /var/local/oops/server

COPY .env /var/local/oops/.env

COPY privkey.pem /var/local/oops/tls

COPY pubkey.pem /var/local/oops/tls

ENV OOPS_ENV_FILE=/var/local/oops/.env

EXPOSE 443

CMD ["/var/local/oops/server"]