FROM postgres:16-alpine
RUN apk add --no-cache tzdata
ENV TZ Europe/Moscow