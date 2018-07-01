FROM scratch

WORKDIR /app

ADD nuc0 certificate.pem key.pem /app/

EXPOSE 3333

CMD ["nuc0"]
