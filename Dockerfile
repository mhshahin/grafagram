FROM golang:1.18-alpine3.14 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o grafagram .

FROM alpine:3.16.2
RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /build/grafagram /app
COPY --from=builder /build/alert-layout.html /app

EXPOSE 1323
ENTRYPOINT [ "./grafagram" ]