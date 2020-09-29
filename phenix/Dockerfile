FROM node:14.2 AS jsbuilder

RUN npm install -g @vue/cli

COPY ./src/js /phenix/src/js

WORKDIR /phenix/src/js

ARG PHENIX_WEB_AUTH
ENV VUE_APP_AUTH ${PHENIX_WEB_AUTH:-enabled}

ARG NODE_TLS=1
ENV NODE_TLS_REJECT_UNAUTHORIZED=$NODE_UNAUTH
ARG STRICT_SSL=true
ENV npm_config_strict_ssl=$STRICT_SSL

RUN npm install \
  && npm run build


FROM golang:1.14 AS gobuilder

RUN apt update \
  && apt install -y protobuf-compiler xz-utils

COPY ./Makefile /phenix/Makefile
COPY ./src/go   /phenix/src/go

WORKDIR /phenix

COPY --from=jsbuilder /phenix/src/js /phenix/src/js

ARG PHENIX_VERSION=

RUN wget https://github.com/glattercj/vmdb2/releases/download/v1.0/vmdb2 -P bin/ \
  && chmod +x bin/vmdb2

RUN VER=${PHENIX_VERSION} make bin/phenix


FROM ubuntu:20.04

RUN apt-get update && apt-get install -y \
  cpio \
  locales \
  vmdb2 \
  && locale-gen en_US.UTF-8 \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* \
  && rm -rf /var/cache/apt/archives/*
  
ENV LANG en_US.UTF-8
ENV LC_ALL en_US.UTF-8

COPY --from=gobuilder /phenix/bin/vmdb2 /usr/bin/vmdb2
COPY --from=gobuilder /phenix/bin/phenix /usr/bin/phenix

CMD ["phenix", "help"]
