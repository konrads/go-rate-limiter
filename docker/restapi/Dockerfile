FROM golang:1.16 AS build
ADD . /src
WORKDIR /src
RUN go build -o /restapi cmd/restapi/restapi.go

FROM debian:buster
EXPOSE 8080
WORKDIR /
COPY limits.json /
COPY --from=build /restapi /
CMD bash -c "/restapi -rest-uri=0.0.0.0:8080 -limit-conf=limits.json"