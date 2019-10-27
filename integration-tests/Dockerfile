FROM golang:1.12-alpine as build
RUN apk add --update git

RUN mkdir /toasties-galore
WORKDIR /toasties-galore

COPY go.mod . 
COPY go.sum .
RUN GO111MODULES=on CGO_ENABLED=0 go mod download

# Above gets cached. Below is not.

COPY . .
#RUN cd integration-tests && GO111MODULES=on CGO_ENABLED=0 go test -i .

ENTRYPOINT ["go", "test", "-v", "./integration-tests"]
