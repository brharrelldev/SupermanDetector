FROM golang:latest as GoBuild
ENV GO111MODULE on
WORKDIR /app
ADD . /app
RUN go mod vendor
RUN go build -o /usr/local/bin/super-cli -v

FROM alpine:latest
RUN apk add sqlite

WORKDIR /app
EXPOSE 3000
COPY --from=GoBuild /usr/local/bin/super-cli /usr/local/bin/super-cli
COPY --from=GoBuild /app/GeoLite2-City.mmdb .
COPY --from=GoBuild /app/data1.csv .
ENV LOGIN_DB=/app/login.db
ENV MMDB=/app/GeoLite2-City.mmdb
ENTRYPOINT ["super-cli", "--db-file","logins.db", "--dataset", "./data1.csv", "server", "--port", "3000"]




