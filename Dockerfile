FROM golang:1.14.6-alpine3.12 as builder
COPY go.mod go.sum /go/src/github.com/dntuanvu/sphtech-blog-system/
WORKDIR /go/src/github.com/dntuanvu/sphtech-blog-system/
RUN go mod download
COPY . /go/src/github.com/dntuanvu/sphtech-blog-system/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/sphtech-blog-system github.com/dntuanvu/sphtech-blog-system

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/dntuanvu/sphtech-blog-system/build/sphtech-blog-system /usr/bin/sphtech-blog-system
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/sphtech-blog-system"]