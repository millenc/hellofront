FROM golang:1.15.5-buster AS builder

RUN apt-get update && apt-get install -y --no-install-recommends git
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o /go/bin/hellofront

FROM scratch
COPY --from=builder /go/bin/hellofront /app/hellofront
COPY ./templates /app/templates
WORKDIR /app
ENTRYPOINT ["/app/hellofront"]
