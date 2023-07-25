FROM golang:alpine as Build

# We need tzdata for the timezone information and the
# ca-certificates for ssl cert verification
RUN apk --no-cache add tzdata ca-certificates

WORKDIR /app

COPY * ./

RUN go build -a -tags netgo -ldflags '-w' -v -o main .

FROM scratch

WORKDIR /app

COPY --from=Build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=Build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=Build /app/main .

CMD [ "/app/main" ]