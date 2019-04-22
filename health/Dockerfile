# uncomment for large build
#FROM golang:1.12.4
FROM golang:1.12.4-alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 go install
# uncomment for large build
#ENTRYPOINT ["/go/bin/app"]

#to build scratch image comment when large build are required
FROM scratch
WORKDIR /opt
COPY --from=builder /go/src/app/premium_tables.xlsx .
COPY --from=builder /go/bin/app .
ENTRYPOINT ["/opt/app"]