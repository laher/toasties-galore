FROM golang:1.12-alpine as build
RUN apk add --update git

RUN mkdir /toasties-galore
WORKDIR /toasties-galore

COPY go.mod . 
COPY go.sum .
RUN GO111MODULES=on CGO_ENABLED=0 go mod download

# Above gets cached. Below is not.

COPY . .
RUN cd jafflotron && GO111MODULES=on CGO_ENABLED=0 go build -a .

FROM alpine:3.9
RUN apk add --update ca-certificates tzdata && \
    rm -rf /var/cache/apk/* /tmp/*
RUN update-ca-certificates

COPY --from=build /toasties-galore/jafflotron /

ENTRYPOINT ["/jafflotron"]
