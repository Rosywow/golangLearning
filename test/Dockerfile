#docker image rm -f test
#docker rm -f chatApp_test
#docker build -t test:latest -f Dockerfile .
#docker run --name chatApp_test --network chatnet test:latest
FROM golang:alpine AS build

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.io
RUN go mod download

COPY *.go ./
RUN go build -o /chatapp

FROM alpine:latest

WORKDIR /app

COPY --from=build /chatapp /app/chatapp

EXPOSE 1234

#这个必须要加，不然没有权限


ENTRYPOINT ["/app/chatapp"]

