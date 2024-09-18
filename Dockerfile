FROM golang:alpine AS build
LABEL stage=gobuilder
ENV GOOS=linux
RUN apk update --no-cache && apk add --no-cache tzdata
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /bin/aeroserver ./cmd/aeroserver/aeroserver.go

FROM alpine
COPY --from=build /usr/share/zoneinfo/Europe/Moscow  /usr/share/zoneinfo/Europe/Moscow
ENV TZ=Europe/Moscow
WORKDIR /bin
COPY --from=build /bin/aeroserver /bin/aeroserver
CMD ["./aeroserver"]