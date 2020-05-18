# Context folder: root of repository
FROM golang:1.14 as build

ENV GOPROXY http://nxrm/repository/gomod/
ENV GOSUMDB off
ENV CGO_ENABLED 0

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o main .

FROM scratch
COPY /db/migrations /db/migrations
COPY --from=build /build/main /
EXPOSE 80
ENTRYPOINT ["/main"]
