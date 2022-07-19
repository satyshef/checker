ARG APP_NAME="checker"
ARG APP_DIR="/app"
ARG DATA_DIR="data"


FROM golang:1.18-alpine3.15 AS builder
LABEL stage=builder
ARG APP_NAME
ARG APP_DIR

COPY --from=satyshef/tdlib:1.8.3-alpine3.15 /usr/local/include/td /usr/local/include/td
COPY --from=satyshef/tdlib:1.8.3-alpine3.15 /usr/local/lib/libtd* /usr/local/lib/
COPY --from=satyshef/tdlib:1.8.3-alpine3.15 /usr/lib/libssl.a /usr/local/lib/libssl.a
COPY --from=satyshef/tdlib:1.8.3-alpine3.15 /usr/lib/libcrypto.a /usr/local/lib/libcrypto.a
COPY --from=satyshef/tdlib:1.8.3-alpine3.15 /lib/libz.a /usr/local/lib/libz.a
RUN apk add build-base

WORKDIR ${APP_DIR}

COPY go.mod go.sum ./
RUN apk add git && \
    apk add make && \
    go mod download && \
    go mod verify

COPY . .
RUN  go mod tidy && make -e APP_PATH=${APP_NAME}

# finish
FROM alpine:3.15
ARG APP_NAME
ARG APP_DIR
ARG DATA_DIR

#ENV APP_PATH=${APP_DIR}/${APP_NAME}
#ENV GONFIG_PATH=${APP_DIR}/data/${CONF_NAME}

WORKDIR ${APP_DIR}

COPY --from=builder ${APP_DIR}/${APP_NAME} ${APP_NAME}
COPY --from=builder ${APP_DIR}/${DATA_DIR} ${DATA_DIR}
COPY --from=builder ${APP_DIR}/entrypoint.sh .
#COPY --from=builder ${APP_DIR}/profiles ${APP_DIR}/profiles

RUN mkdir profiles && chmod +x entrypoint.sh && apk add libstdc++

CMD [ "./entrypoint.sh" ]