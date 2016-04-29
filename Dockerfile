#FROM centurylink/ca-certs
FROM ubuntu
MAINTAINER EOGILE "agilestack@eogile.com"

ENV name hydra-host

ENV PORT                   9090
ENV HOST_URL               http://localhost:$PORT
ENV DATABASE_URL           postgres://hydra:hydra_agilestack@hydra-postgres:5432/hydra?sslmode=disable
ENV SIGNUP_URL             http://localhost:8080/register
ENV SIGNIN_URL             http://localhost:8080/login
ENV JWT_PUBLIC_KEY_PATH    /cert/rs256-public.pem
ENV JWT_PRIVATE_KEY_PATH   /cert/rs256-private.pem
ENV TLS_CERT_PATH          /cert/tls-cert.pem
ENV TLS_KEY_PATH           /cert/tls-key.pem
ENV DANGEROUSLY_FORCE_HTTP force

EXPOSE $PORT

ADD cert /cert
ADD $name /$name

CMD ["/hydra-host", "start"]
