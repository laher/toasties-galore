FROM golang:1.12-alpine as build
RUN apk add --update git
COPY . /toasties-galore
RUN cd /toasties-galore && GO111MODULES=on CGO_ENABLED=0 go mod download
RUN cd /toasties-galore/jafflotron && GO111MODULES=on CGO_ENABLED=0 go build -a .


FROM alpine:3.9
RUN apk add --update ca-certificates tzdata && \
    rm -rf /var/cache/apk/* /tmp/*
RUN update-ca-certificates

COPY --from=build /toasties-galore/jafflotron /

ENTRYPOINT ["/go-app"]
