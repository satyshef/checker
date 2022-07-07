ARG APP_NAME="checker"
ARG APP_DIR="/app"
ARG CONF_NAME="config.toml"



FROM golang:1.18-alpine3.15 AS builder
LABEL stage=builder
ARG APP_NAME
ARG APP_DIR

COPY --from=wcsiu/tdlib:1.8-alpine /usr/local/include/td /usr/local/include/td
COPY --from=wcsiu/tdlib:1.8-alpine /usr/local/lib/libtd* /usr/local/lib/
COPY --from=wcsiu/tdlib:1.8-alpine /usr/lib/libssl.a /usr/local/lib/libssl.a
COPY --from=wcsiu/tdlib:1.8-alpine /usr/lib/libcrypto.a /usr/local/lib/libcrypto.a
COPY --from=wcsiu/tdlib:1.8-alpine /lib/libz.a /usr/local/lib/libz.a
RUN apk add build-base

WORKDIR ${APP_DIR}

COPY go.mod go.sum ./
RUN apk add git && \
    apk add make && \
    go mod download && \
    go mod verify

COPY . .
RUN  go mod tidy && make -e APP_PATH=${APP_NAME}

#RUN go build --ldflags "-extldflags '-static -L/usr/local/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -ldl -lm -lssl -lcrypto -lstdc++ -lz'" -o tebot cmd/app/main.go

# finish
FROM alpine:3.15
ARG APP_NAME
ARG APP_DIR
ARG CONF_NAME

ENV APP_PATH=${APP_DIR}/${APP_NAME}
ENV GONFIG_PATH=${APP_DIR}/data/${CONF_NAME}

WORKDIR ${APP_DIR}

COPY --from=builder ${APP_DIR}/data/ ./data
COPY --from=builder ${APP_DIR}/data/${CONF_NAME} ${GONFIG_PATH}
COPY --from=builder ${APP_DIR}/${APP_NAME} ${APP_PATH}
#COPY --from=builder ${APP_DIR}/profiles ${APP_DIR}/profiles

RUN apk add libstdc++
EXPOSE 7070
CMD ${APP_PATH} "-c" ${GONFIG_PATH}
#CMD [ ${APP_PATH}, "-c", ${CONF_PATH} ]
