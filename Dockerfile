FROM scratch
ADD nuc0 certificate.pem key.pem /
EXPOSE 3333
CMD ["/nuc0"]
