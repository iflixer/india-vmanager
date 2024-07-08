# build environment
FROM golang:1.21.1 as build-env
WORKDIR /server
COPY src/go.mod ./
RUN go mod download
COPY src src
WORKDIR /server/src
RUN CGO_ENABLED=0 GOOS=linux go build -o /server/build/httpserver .

FROM docker
WORKDIR /app
#RUN set -x \
#    && apk add --no-cache ca-certificates tzdata \
#    && cp /usr/share/zoneinfo/Europe/Kiev /etc/localtime \
#    && echo Europe/Kiev > /etc/timezone \
#    && apk del tzdata

COPY --from=build-env /server/build/httpserver /app/videomanager

#ENV GITHUB-SHA=<GITHUB-SHA>

EXPOSE 19494/tcp

ENTRYPOINT [ "/app/videomanager" ]
#ENTRYPOINT [ "ls", "-la", "/app/httpserver" ]
